package cadlranch

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Alancere/azureutils/common"
	"github.com/iancoleman/strcase"
)

const (
	TSP_COMPILE_JS = "tspcompile.js"
	CADL_RANCH_JS  = "cadl-ranch.js"

	TSP_COMPILE_PATH = ".scripts/tspcompile.js"
	CADL_RANCH_PATH  = ".scripts/cadl-ranch.js"
)

// 添加一个resourcemanager的cadl-ranch test
/*
 1. 指定 typespec-go path
 2. 指定要生成的 cadl-ranch-specs path
*/
func GenerateCadlRanchTest(typepsecGoPath, cadlRanchSpecsPath string) error {
	if err := ValidationCadlRanchTestPath(cadlRanchSpecsPath); err != nil {
		return err
	}

	var err error
	typepsecGoPath, err = filepath.Abs(typepsecGoPath)
	if err != nil {
		return err
	}

	// read and write tspcompile.js
	tspCompilePath := filepath.Join(typepsecGoPath, TSP_COMPILE_PATH)
	tspCompileJsData, err := os.ReadFile(tspCompilePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(tspCompileJsData), "\n")
	tName := cadlRanchTestName(cadlRanchSpecsPath)
	var strBuf strings.Builder
	for i, line := range lines {
		strBuf.WriteString(line)
		strBuf.WriteString("\n")
		if strings.Contains(line, "const cadlRanch = {") {
			strBuf.WriteString(fmt.Sprintf("  '%s': ['%s'],", tName, cadlRanchSpecsPath))
			strBuf.WriteString("\n")
			strBuf.WriteString(strings.Join(lines[i+1:], "\n"))
			if err := os.WriteFile(tspCompilePath, []byte(strBuf.String()), 0o644); err != nil {
				return err
			}
		}
	}

	// run tspcompile.js --filter
	if err := common.Node(typepsecGoPath, tspCompilePath, "--filter", tName); err != nil {
		return err
	}

	return nil
}

func cadlRanchTestName(cadlRanchSpecsPath string) string {
	_, f := filepath.Split(cadlRanchSpecsPath)
	return strcase.ToSnake(f)
}

// https://github.com/Azure/cadl-ranch/tree/main/packages/cadl-ranch-specs/http/azure/resource-manager/models
func ValidationCadlRanchTestPath(cadlRanchSpecsPath string) error {
	url := fmt.Sprintf("https://github.com/Azure/cadl-ranch/tree/main/packages/cadl-ranch-specs/http/%s", strings.ReplaceAll(strings.Trim(cadlRanchSpecsPath, "/"), "\\", "/"))
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("cadl-ranch-specs path %s not found", url)
	}

	return nil
}
