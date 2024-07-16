package mergego

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alancere/azureutils/common"
)

func Merge(dir string, outfile string, mergeTest bool) error {
	if outfile == "" {
		outfile = filepath.Join(dir, "merged.go")
	}

	_, err := os.Stat(filepath.Dir(outfile))
	if err != nil {
		err := os.MkdirAll(filepath.Dir(outfile), os.ModePerm)
		if err != nil {
			return err
		}
	}

	if err := MergeGo(dir, outfile); err != nil {
		return err
	}

	outputTestFile := ""
	if mergeTest {
		outputTestFile = strings.TrimSuffix(outfile, ".go") + "_test.go"
		if err := MergeTest(dir, outputTestFile); err != nil {
			return err
		}
	}

	// goimports
	common.GoImports(dir, "-w", outfile)

	// gofumpt
	common.GoFumpt(dir, "-w", outfile)

	if outputTestFile != "" {
		// goimports
		common.GoImports(dir, "-w", outputTestFile)

		// gofumpt
		common.GoFumpt(dir, "-w", outputTestFile)
	}

	return nil
}

func MergeGo(dir string, outfile string) error {
	filter := func(info os.FileInfo) bool {
		// Skip test files
		if strings.HasSuffix(info.Name(), "_test.go") {
			return false
		}

		if info.Name() == "build.go" {
			return false
		}

		return true
	}

	return MergeX(dir, outfile, filter)
}

func MergeTest(dir string, outfile string) error {
	filter := func(info fs.FileInfo) bool {
		// skip live test
		if strings.HasSuffix(info.Name(), "_live_test.go") || strings.HasSuffix(info.Name(), "utils_test.go"){
			return false
		}
		return strings.HasSuffix(info.Name(), "_test.go")
	}
	return MergeX(dir, outfile, filter)
}

func MergeX(dir string, outfile string, filter func(fs.FileInfo) bool) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, filter, parser.ParseComments)
	if err != nil {
		return err
	}

	merged := &ast.File{}
	for k := range pkgs {
		merged = ast.MergePackageFiles(pkgs[k], ast.FilterImportDuplicates) // |ast.FilterUnassociatedComments)
	}

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
