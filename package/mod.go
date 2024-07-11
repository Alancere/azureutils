package packages

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Alancere/azureutils/common"
	"golang.org/x/mod/modfile"
)

// 获取go.mod文件信息
func ModFile(goModPath string) (*modfile.File, error) {
	if filepath.Base(goModPath) != "go.mod" {
		return nil, fmt.Errorf("go.mod file path is invalid: %s", goModPath)
	}

	data, err := os.ReadFile(goModPath)
	if err != nil {
		return nil, err
	}

	modFile, err := modfile.Parse(goModPath, data, nil)
	if err != nil {
		return nil, err
	}

	return modFile, nil
}

var ErrInvalidModule = fmt.Errorf("go.mod does not allow to include other versions of module")

// 校验go.mod
func GoModValidate(packagePath string) error {
	goModPath := filepath.Join(packagePath, "go.mod")
	modFile, err := ModFile(goModPath)
	if err != nil {
		return err
	}

	baseModule := getBaseModule(modFile.Module)
	if baseModule == "" {
		return nil
	}

	// module 不建议引用 其他般的的module，如当前是v5,则不建议引用v4及以下的版本
	for _, require := range modFile.Require {
		if strings.Contains(require.Mod.Path, baseModule) {
			return errors.Join(ErrInvalidModule, fmt.Errorf(": %s", require.Mod.Path))
		}
	}

	return nil
}

func getBaseModule(mod *modfile.Module) string {
	if mod == nil || mod.Mod.Path == "" {
		return ""
	}

	parts := strings.Split(mod.Mod.Path, "/")
	if regexp.MustCompile(`^v\d+$`).MatchString(parts[len(parts)-1]) {
		return strings.Join(parts[:len(parts)-1], "/")
	}

	return mod.Mod.Path
}

func EditModule(packagePath, newModule string) error {
	_, err := common.Go(packagePath, "mod", "edit", "-module", newModule)
	return err
}
