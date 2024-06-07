package generate

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// generate go sdk from typespec

func generate() {
	// generate release-v2 args --TyepSpec
}

// 解析 tspconfig.yaml
// 提供两种路径(local or http path)
func ParsePath(tspConfigPath string) (data []byte, err error) {
	data = make([]byte, 0, 1024)
	// 判断路径类型
	if strings.HasPrefix(tspConfigPath, "http") {
		// http path
		resp, err := http.Get(tspConfigPath)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// 读取resp body数据
		data, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	} else {
		// local path
		data, err = os.ReadFile(tspConfigPath)
		if err != nil {
			return nil, err
		}
	}

	return
}

// 判断是否为preview版本
// 通过读取tspSpecPath文件下的所有.tsp文件，判断是否有preview版本
// 有则返回true，否则返回false
func currentTSPIsPreview(tspSpecPath string) (bool, error) {

	isPreview := false
	err := filepath.Walk(tspSpecPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(info.Name(), ".tsp") {
			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// if strings.Contains(string(data), "preview") {}
			// 需要确认怎么从tsp中确定要release的版本
			if strings.Contains(string(data), "enum Versions {") {
				lines := strings.Split(string(data), "\n")
				for i, line := range lines {
					if line ==  "}" {
						if strings.Contains(lines[i-1], "preview") {
							isPreview = true
							return filepath.SkipAll
						}
					}
				}
				
			}
		}

		return nil
	})
	if err != nil && !errors.Is(err, filepath.SkipDir) {
		return false, err
	}

	return isPreview, nil
}
