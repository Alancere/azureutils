package sdk

import (
	"os"
	"path/filepath"
)

const recordFile = "assets.json"

// 如何当前目录下有assets.json就认为有live test
func ContainLiveTest(packagePath string) (bool, error) {
	_, err := os.Stat(filepath.Join(packagePath, recordFile))
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
