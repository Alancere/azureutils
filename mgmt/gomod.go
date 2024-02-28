package mgmt

import (
	"os"

	"golang.org/x/mod/modfile"
)

func RemoveGoModRequire(gomodFile string) error {
	data, err := os.ReadFile(gomodFile)
	if err != nil {
		return err
	}

	file, err := modfile.Parse(gomodFile, data, nil)
	if err != nil {
		return err
	}

	file.SetRequire(nil)

	newData, err := file.Format()
	if err != nil {
		return err
	}

	return os.WriteFile(gomodFile, newData, os.ModePerm)
}
