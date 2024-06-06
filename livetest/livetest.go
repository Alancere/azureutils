package livetest

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// go.mod 包含 github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal 即当前arm service有live test
func AllLiveTestARMService(root string) ([]string, error) {

	livetests := make([]string, 0, 100)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() || !strings.Contains(path, "arm"){
			return nil
		}

		// skip internal and fake dir
		if d.Name() == "internal" || d.Name() == "fake" || d.Name() == ".vscode" {
			return nil
		}

		// read go.mod 
		data, err := os.ReadFile(filepath.Join(path, "go.mod"))
		if err != nil {
			log.Println("read go.mod err:", err)
			return nil
		}

		if strings.Contains(string(data), "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/internal") {
			livetests = append(livetests, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return livetests, nil
}

// 往指定目录写入一个文件
func WriteFile(path, fName string, data []byte) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Create(filepath.Join(path, fName))
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func GetTestdataPath(livetest string) (string, error) {

	testdataPath := ""

	err := filepath.WalkDir(livetest, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.Contains(path, "_live_test.go") {
			data,err := os.ReadFile(path)
			if err != nil {
				return err
			}

			// regexp sdk/resourcemanager/redis/armredis/testdata
			testdataPath = regexp.MustCompile(`sdk/resourcemanager/.*?/testdata`).FindString(string(data))
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return testdataPath, nil
}