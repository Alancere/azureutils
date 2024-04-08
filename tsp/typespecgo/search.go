package typespecgo

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// 有tspconfig.yaml并且有readme.go.md的目录
func SearchTSP(dir string) ([]string, error) {
	count := 0
	tspConfigs := make([]string, 0, 128)
	if err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() != "tspconfig.yaml" { // && strings.Contains(path, "resourcemanager")
			return nil
		}

		// 存在 mgmt
		readmeGOMDPath := filepath.Join(filepath.Dir(path), "../resource-manager/readme.go.md")
		_, err = os.Stat(readmeGOMDPath)
		if err != nil {
			return nil
		}

		fmt.Println(path)
		count++
		tspConfigs = append(tspConfigs, path)

		return nil
	}); err != nil {
		return nil, err
	}

	fmt.Println("tspconfig.yaml count:", count)
	return tspConfigs, nil
}
