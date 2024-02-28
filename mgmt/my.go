package mgmt

import (
	"io/fs"
	"os"
	"path/filepath"
)

/*
	通过指定文件夹，对其中的所有go.mod执行 go mod tidy
*/

func GoModTidy(dir string) error {
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || d.Name() == ".vscode" {
			return nil
		}

		gomodPath := filepath.Join(path, "go.mod")
		_, err = os.Stat(gomodPath)
		if err != nil {
			return nil
		}

		// if err = RemoveGoModRequire(gomodPath); err != nil {
		// 	return err
		// }

		// go get -u ./...
		if err = ExecGo(path, "get", "-u", "./..."); err != nil {
			return err
		}

		// go mod tidy
		if err = ExecGo(path, "mod", "tidy"); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func GoVet(dir string) error {
	if err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || d.Name() == ".vscode" {
			return nil
		}

		gomodPath := filepath.Join(path, "go.mod")
		_, err = os.Stat(gomodPath)
		if err != nil {
			return nil
		}

		if err = ExecGo(path, "vet", "./..."); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
