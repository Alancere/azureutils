package livetest_test

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/common"
	"github.com/Alancere/azureutils/livetest"
)

func TestAllLiveTestARMService(t *testing.T) {
	sdkRoot := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	livetests, err := livetest.AllLiveTestARMService(sdkRoot)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("live test count:", len(livetests))
	for _, v := range livetests {
		fmt.Println(v)

		// skip armredis
		if strings.HasSuffix(v, "armredis") {
			continue
		}

		// 从 live_test.go 中获取 testdata 目录
		testdataPath, err := livetest.GetTestdataPath(v)
		if err != nil {
			t.Fatal(err)
		}
		if testdataPath == "" {
			log.Fatal("##get testdata path failed:", v)
			continue
		}

		// write utils_test.go 到当前目录下
		x := strings.Split(v, "\\")[len(strings.Split(v, "\\"))-1]
		err = livetest.WriteFile(v, "utils_test.go", []byte(fmt.Sprintf(utilsTestData, x, testdataPath)))
		if err != nil {
			t.Fatal(err)
		}

		// go.mod 文件中添加 replace
		err = addModReplace(filepath.Join(v, "go.mod"), "replace github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v2 v2.0.0 => ../../../../../dev/azure-sdk-for-go/sdk/resourcemanager/internal")
		if err != nil {
			t.Fatal(err)
		}

		common.Go(v, "mod", "tidy")
	}
}

var utilsTestData = `//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package %s_test

import (
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal/v3/testutil"
)

const (
	pathToPackage = "%s"
)

func TestMain(m *testing.M) {
	code := run(m)
	os.Exit(code)
}

func run(m *testing.M) int {
	f := testutil.StartProxy(pathToPackage)
	defer f()
	return m.Run()
}
`

// 在go.mod 文件末尾添加 replaces
func addModReplace(gomod string, replace ...string) error {
	f, err := os.OpenFile(gomod, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	for _, v := range replace {
		_, err = f.WriteString(v + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// 执行go test 并且执行 test-proxy update
func TestUpdateAssets(t *testing.T) {
	goTestFails := make([]string, 0, 10)
	updateAssetsFails := make([]string, 0, 10)

	sdkRoot := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	livetests, err := livetest.AllLiveTestARMService(sdkRoot)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("live test count:", len(livetests))
	for _, v := range livetests {
		fmt.Println(v)

		// go test -v .
		_, err := common.Go(v, "test", "-v", "--count=1", ".")
		if err != nil {
			goTestFails = append(goTestFails, v)
			continue
		}

		// test-proxy update assets
		_, err = common.TestProxy(v, "push", "-a", "assets.json")
		if err != nil {
			updateAssetsFails = append(updateAssetsFails, v)
			continue
		}
	}

	fmt.Println("go test fails counts:", len(goTestFails))
	for _, v := range goTestFails {
		fmt.Println(v)
	}

	fmt.Println("test-proxy update assets fails counts:", len(updateAssetsFails))
	for _, v := range updateAssetsFails {
		fmt.Println(v)
	}
}

// passed
func TestUseInternalV3(t *testing.T) {
	sdkRoot := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	livetests, err := livetest.AllLiveTestARMService(sdkRoot)
	if err != nil {
		t.Fatal(err)
	}

	for _, v := range livetests {
		// insert UsePipelineProxy: false to ci.yml
		data, err := os.ReadFile(filepath.Join(v, "ci.yml"))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "UsePipelineProxy: false") {
			ci, err := os.OpenFile(filepath.Join(v, "ci.yml"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			defer func() {
				ci.Close()
			}()
			if err != nil {
				t.Fatal(err)
			}
			ci.WriteString("    UsePipelineProxy: false")
			ci.WriteString("\n")
		}

		// replace pathToPackage to _live_test.go files
		// 从 live_test.go 中获取 testdata 目录
		testdataPath, err := livetest.GetTestdataPath(v)
		if err != nil {
			t.Fatal(err)
		}

		if testdataPath != "" {
			// 将当前目录下_live_test.go 文件中的 testdataPathValue 替换成 const pathToPackage
			err = filepath.Walk(v, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if strings.Contains(info.Name(), "_live_test.go") {
					data, err := os.ReadFile(path)
					if err != nil {
						return err
					}

					newData := bytes.Replace(data, []byte(fmt.Sprintf("\"%s\"", testdataPath)), []byte("pathToPackage"), -1)

					// 相同则不更新文件
					if slices.Equal(data, newData) {
						return nil
					}

					err = os.WriteFile(path, newData, 0o644)
					if err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		// go mod tidy
		common.Go(v, "mod", "tidy")
	}
}
