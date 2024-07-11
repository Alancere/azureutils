package autorest

import (
	"path/filepath"

	"github.com/Alancere/azureutils/common"
)

// 执行autorest cmd

func AutorestGo(outputFolder, autorestMDPath string, opts ...string) error {
	if autorestMDPath == "" {
		autorestMDPath = filepath.Join(outputFolder, "autorest.md")
	}

	opts = append(opts, autorestMDPath)
	_, err := common.AutorestCmd(outputFolder, opts...)
	if err == nil {
		return err
	}

	return nil
}
