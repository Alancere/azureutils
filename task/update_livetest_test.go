package task_test

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/common"
	packages "github.com/Alancere/azureutils/package"
	"github.com/Alancere/azureutils/sdk"
	"github.com/stretchr/testify/assert"
)

// 1. 升级 resourcemanager/internal to v3.1.0
// 2. ci.yml 添加 UseFederatedAuth: true
func TestUpdateCred(t *testing.T) {
	mgmtPath := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"

	errors := make([]error, 0)

	liveTestMgmts := make([]string, 0)
	err := filepath.WalkDir(mgmtPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || d.Name() == ".vscode" || d.Name() == "fake" || !strings.HasPrefix(d.Name(), "arm") {
			return nil
		}

		b, err := sdk.ContainLiveTest(path)
		if err != nil {
			return err
		}

		if b {
			fmt.Println(path)
			liveTestMgmts = append(liveTestMgmts, path)
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, m := range liveTestMgmts {
		// upgrade resourcemanager/internal dep
		internalVersion := "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3@v3.1.0"
		err = packages.EditModRequired(m, internalVersion)
		if err != nil {
			errors = append(errors, err)
			continue
		}

		// add UseFederatedAuth: true to ci.yml
		ciPath := filepath.Join(m, "ci.yml")
		data, err := os.ReadFile(ciPath)
		if err != nil {
			errors = append(errors, fmt.Errorf("read ci.yml fail: %s\n%+v", m, err))
			continue
		}
		if strings.Contains(string(data), "UseFederatedAuth") {
			fmt.Println("already have UseFederatedAuth", m)
			continue
		}
		w, err := os.OpenFile(ciPath, os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			errors = append(errors, fmt.Errorf("open ci.yml fail: %s\n%+v", m, err))
			continue
		}
		defer w.Close()

		_, err = w.WriteString("    UseFederatedAuth: true\n")
		if err != nil {
			errors = append(errors, fmt.Errorf("write ci.yml fail: %s\n%+v", m, err))
			continue
		}

		// go mod tidy
		_, err = common.Go(m, "mod", "tidy")
		if err != nil {
			errors = append(errors, fmt.Errorf("go mod tidy fail: %s\n%+v", m, err))
			continue
		}

		// go vet
		_, err = common.Go(m, "vet", "./...")
		if err != nil {
			errors = append(errors, fmt.Errorf("go vet fail: %s\n%+v", m, err))
			continue
		}

		// go test
		_, err = common.Go(m, "test", "-v", ".")
		if err != nil {
			errors = append(errors, fmt.Errorf("go test fail: %s\n%+v", m, err))
			continue
		}
	}

	fmt.Println("errors count:", len(errors))
	for _, e := range errors {
		fmt.Println(e)
	}
}

// pass
func TestAppendFile(t *testing.T) {
	LocalPath := "D:\\Go\\src\\github.com\\Azure\\azure-sdk-for-go\\sdk\\resourcemanager\\apicenter\\armapicenter"
	w, err := os.OpenFile(filepath.Join(LocalPath, "ci.yml"), os.O_RDWR|os.O_APPEND, 0o666)
	assert.NoError(t, err)
	defer w.Close()

	_, err = w.WriteString("    UseFederatedAuth: true\n")
	assert.NoError(t, err)

	w.Sync()
}

func TestLoopDir(t *testing.T) {
	mgmtPath := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"

	liveTestMgmts := make([]string, 0)
	err := filepath.WalkDir(mgmtPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || d.Name() == ".vscode" || d.Name() == "fake" || !strings.HasPrefix(d.Name(), "arm") {
			return nil
		}

		b, err := sdk.ContainLiveTest(path)
		if err != nil {
			return err
		}

		if b {
			fmt.Println(path)
			liveTestMgmts = append(liveTestMgmts, path)
		}

		return nil
	})
	assert.NoError(t, err)

	fmt.Println("live test mgmt count:", len(liveTestMgmts)) // 70
}
