package mgmt

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestGoModtidy(t *testing.T) {
	dir := "D:/Go/src/github.com/Azure/azure-sdk-for-go-samples/sdk/resourcemanager"
	dir = "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	GoModTidy(dir)
}

func TestGoVet(t *testing.T) {
	dir := "D:/Go/src/github.com/Azure/azure-sdk-for-go-samples/sdk/resourcemanager"
	GoVet(dir)
}

/*
Updating ARM package dependencies
*/
func TestUpdateARMDependencies(t *testing.T) {
	mgmtPath := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"

	if err := filepath.WalkDir(mgmtPath, func(path string, d fs.DirEntry, err error) error {
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

		// go get -u ./...
		if err = ExecGo(path, "get", "-u", "./..."); err != nil {
			return err
		}

		// go mod tidy
		if err = ExecGo(path, "mod", "tidy"); err != nil {
			return err
		}

		version, err := AutorestMdVersion(filepath.Join(path, "autorest.md"))
		if err != nil {
			return err
		}
		v := version.IncPatch()
		version = &v

		if err = ReplaceVersion(path, version.String()); err != nil {
			return err
		}

		changelog := `### Feature
- Update packages to the latest azcore.
`
		releaseDate := ""
		additionalChangelog, err := AddChangelogToFile(changelog, version, path, releaseDate)
		if err != nil {
			return err
		}
		_ = additionalChangelog
		fmt.Println(additionalChangelog)

		return nil
	}); err != nil {
		t.Error(err)
	}
}
