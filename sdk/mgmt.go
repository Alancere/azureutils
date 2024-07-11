package sdk

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alancere/azureutils/autorest"
	"golang.org/x/mod/modfile"
)

type MgmtSDK struct {
	// e: redis/armredis
	Name    string
	ArmName string

	SpecName string // specification 第一层目录
	SpecPath string // 到readme.md的路径

	LocalPath string

	AutoRest *autorest.AutoRestMarkdown

	Module *Module

	GoMod *modfile.File
}

func NewMgmtSDK(localPath string) (*MgmtSDK, error) {
	sdk := new(MgmtSDK)
	sdk.LocalPath = localPath

	// go.mod
	gomod := filepath.Join(localPath, "go.mod")
	data, err := os.ReadFile(gomod)
	if err != nil {
		return nil, err
	}
	p, err := modfile.Parse(gomod, data, nil)
	if err != nil {
		return nil, err
	}
	sdk.GoMod = p

	// is deprecated
	if p.Module.Deprecated != "" {
		return sdk, nil
	}

	// autorest.md
	autorest, err := autorest.ReadAutoRestMarkdown(filepath.Join(localPath, "autorest.md"))
	if err != nil {
		return nil, err
	}
	sdk.AutoRest = autorest
	sdk.SpecName, sdk.SpecPath = autorest.GetSpec()

	// constants.go
	module, err := ReadModule(filepath.Join(localPath, "constants.go"))
	if err != nil {
		return nil, err
	}
	sdk.Module = module
	sdk.Name, sdk.ArmName = module.GetSDKAndArm()

	return sdk, nil
}

// constants.go
type Module struct {
	ModuleName    string
	ModuleVersion string
}

func ReadModule(constantPath string) (*Module, error) {
	data, err := os.ReadFile(constantPath)
	if err != nil {
		return nil, err
	}

	module := new(Module)
	for _, line := range strings.Split(string(data), "\n") {
		if module.ModuleName != "" && module.ModuleVersion != "" {
			break
		}

		if strings.Contains(line, "moduleName") {
			m, err := cutConstant(line)
			if err != nil {
				return nil, err
			}
			module.ModuleName = m
		}

		if strings.Contains(line, "moduleVersion") {
			m, err := cutConstant(line)
			if err != nil {
				return nil, err
			}
			module.ModuleVersion = m
		}
	}

	return module, nil
}

func cutConstant(c string) (string, error) {
	_, after, b := strings.Cut(c, "\"")
	if !b {
		return "", errors.New("cut moduleName fail")
	}

	return strings.Trim(strings.TrimSpace(after), "\""), nil
}

// response sdkName and armName
func (m Module) GetSDKAndArm() (string, string) {
	splits := strings.Split(m.ModuleName, "/")
	return splits[len(splits)-2], splits[len(splits)-1]
}

var DeprecatedMgmtSDK = make([]*MgmtSDK, 0, 10)

func GetAllMgmtSDK(src string) ([]*MgmtSDK, error) {
	allMgmt := make([]*MgmtSDK, 0, 200)

	err := filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			return nil
		}

		if strings.HasPrefix(d.Name(), "arm") {
			m, err := NewMgmtSDK(path)
			if err != nil {
				return err
			}
			// skip deprecated
			if m.IsDeprecated() {
				DeprecatedMgmtSDK = append(DeprecatedMgmtSDK, m)
				return nil
			}
			allMgmt = append(allMgmt, m)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return allMgmt, nil
}

// 判断go.mod是否 Deprecated
func (ms MgmtSDK) IsDeprecated() bool {
	return ms.GoMod.Module.Deprecated != ""
}
