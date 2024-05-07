package mergego

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"

	"github.com/Alancere/azureutils/common"
)

func MergeGo(dir string, outfile string) error {
    fset := token.NewFileSet()
    pkgs, err := parser.ParseDir(fset, dir, nil, 0)
    if err != nil {
        return err
    }

	merged := ast.MergePackageFiles(pkgs["armapicenter"], ast.FilterImportDuplicates)

    // Separate import declarations and other declarations
    var importDecls []ast.Decl
    var otherDecls []ast.Decl
    for _, decl := range merged.Decls {
        if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.IMPORT {
            importDecls = append(importDecls, decl)
        } else {
            otherDecls = append(otherDecls, decl)
        }
    }

    // Write import declarations first, then other declarations
    merged.Decls = append(importDecls, otherDecls...)

	f, err := os.Create(outfile)
    if err != nil {
        return err
    }
    defer f.Close()

    err = format.Node(f, fset, merged)
    if err != nil {
        return err
    }

    return nil
}

func Merge(dir string, outfile string) error {
    if outfile == "" {
        outfile = filepath.Join(dir, "merged.go")
    }

    if err := MergeGo(dir, outfile); err != nil {
        return err
    }

    // goimports
    common.GoImports(dir, "-w", outfile)
    
    // gofumpt
    common.GoFumpt(dir, "-w", outfile)

    return nil
}