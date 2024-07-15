package packages

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/tools/go/ast/astutil"
)

/*
使用新的版本号(newVersion)更新包(baseModule)的引用
*/
func UpdatePackageModule(newVersion string, packagePath, baseModule string, suffixs ...string) error {
	version, err := semver.NewVersion(newVersion)
	if err != nil {
		return err
	}

	return filepath.Walk(packagePath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Ext(info.Name()) != ".go" {
			return nil
		}

		hasSuffix := slices.ContainsFunc(suffixs, func(s string) bool { return strings.HasSuffix(info.Name(), s) })
		if len(suffixs) == 0 || hasSuffix {
			if err = ReplaceImport(path, baseModule, version.Major()); err != nil {
				return err
			}
		}

		return nil
	})
}

func ReplaceImport(sourceFile string, baseModule string, majorVersion uint64) error {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, sourceFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	rewrote := false
	for _, i := range f.Imports {
		if strings.HasPrefix(i.Path.Value, fmt.Sprintf("\"%s", baseModule)) {
			oldPath := importPath(i)
			after, _ := strings.CutPrefix(oldPath, baseModule)

			newPath := baseModule
			if after != "" {
				before, sub, _ := strings.Cut(strings.TrimLeft(after, "/"), "/")
				if majorVersion > 1 {
					newPath = fmt.Sprintf("%s/v%d", baseModule, majorVersion)
				}
				if !regexp.MustCompile(`^v\d+$`).MatchString(before) {
					newPath = fmt.Sprintf("%s/%s", newPath, before)
				}
				if sub != "" {
					newPath = fmt.Sprintf("%s/%s", newPath, sub)
				}
			} else {
				if majorVersion > 1 {
					newPath = fmt.Sprintf("%s/v%d", baseModule, majorVersion)
				}
			}

			if newPath != oldPath {
				rewrote = astutil.RewriteImport(fset, f, oldPath, newPath)
			}
		}
	}

	if rewrote {
		w, err := os.Create(sourceFile)
		if err != nil {
			return err
		}
		defer w.Close()

		return printer.Fprint(w, fset, f)
	}

	return nil
}

func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return ""
	}
	return t
}
