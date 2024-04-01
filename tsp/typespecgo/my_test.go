package typespecgo

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

/*
emit:
  - D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go # "@azure-tools/typespec-go"
  # - "@azure-tools/typespec-autorest"
options:
  "@azure-tools/typespec-go":
    module: github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sphere/armsphere
    module-version: 0.2.0
    generate-fakes: true
*/

func TestSearchTSP(t *testing.T) {
	dir := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs"
	_, err := SearchTSP(dir)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerateSDK(t *testing.T) {
	dir := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs"
	typespecgoEmit := "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go"
	// typespecgoEmit = "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go/azure-tools-typespec-go-0.1.0-dev.1.tgz" // 不能用这种方式

	configPaths, err := SearchTSP(dir)
	if err != nil {
		t.Fatal(err)
	}

	tspErrs := make([]error, 0)
	allErrMsg := make([]string, 0)
	tspCompilerLog, err := os.OpenFile(filepath.Join(dir, "tspCompiler.log"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	defer tspCompilerLog.Close()

	bufLog := bufio.NewWriter(tspCompilerLog)
	defer bufLog.Flush()

	for _, configPath := range configPaths {
		// filter
		if strings.Contains(configPath, "machinelearning\\Azure.AI.ChatProtocol") { // 没有go track2 config
			continue
		}

		// read readme.go.md
		readmeGOMD := readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/readme.go.md"))

		module := readmeGOMD["module"]
		moduleName := readmeGOMD["module-name"]
		module = strings.Replace(module.(string), "$(module-name)", moduleName.(string), -1)

		moduleVersion := "0.1.0" // default value, need from autorest.md get

		// update tspconfig
		tspConfig, err := NewTSPConfig(configPath)
		if err != nil {
			t.Fatal(err)
		}

		typespecgoOption := map[string]any{
			"module":             module,
			"module-version":     moduleVersion,
			"emitter-output-dir": fmt.Sprintf("{project-root}/go/%s", moduleName),
			// "generate-fakes": true,
		}

		tspConfig.OnlyEmit(typespecgoEmit)
		tspConfig.EditOptions(string(TypeSpec_GO), typespecgoOption, false)

		err = tspConfig.Write()
		if err != nil {
			t.Fatal(err)
		}

		output, tspErr := TSP(filepath.Dir(configPath), "compile", ".")
		if tspErr != nil {
			tspErrs = append(tspErrs, tspErr)
			// allErrMsg = append(allErrMsg, output)

			// write output in file
			if err = os.WriteFile(filepath.Join(configPath, "../error.log"), []byte(output), 0o777); err != nil {
				t.Fatal(err)
			}
		}
		allErrMsg = append(allErrMsg, output)
		bufLog.WriteString(output)
		bufLog.WriteByte('\n')
		// break

		///
		// go mod tidy and go vet ./...
		// gosdk := fmt.Sprintf("%s/go/%s", filepath.Dir(configPath), moduleName)
		if tspErr == nil {
			gosdk := filepath.Join(filepath.Dir(configPath), "go", moduleName.(string))
			if err = Go(gosdk, "mod", "tidy"); err != nil {
				log.Println("####go mod", err)
			}

			if err = Go(gosdk, "vet", "./..."); err != nil {
				log.Println("####go vet", err)
			}
		}
	}

	// fmt.Println(tspErrs)

	// write error msg to error.log
	fmt.Println("error count:", len(tspErrs))
	errMsg := ""
	for _, eMsg := range allErrMsg {
		errMsg = fmt.Sprintf("%s\n%s", errMsg, eMsg)
	}
	// if err = os.WriteFile(filepath.Join(dir, "tspCompiler.log"), []byte(errMsg), 0777); err != nil {
	//   t.Fatal(err)
	// }
}

func readmegomd(path string) map[string]any {
	result := map[string]any{}

	md, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	data := make([]byte, 0, 1024)
	for _, l := range strings.Split(string(md), "\n") {
		if strings.Contains(l, "module-name:") || strings.Contains(l, "module:") {
			data = append(data, []byte(l)...)
			data = append(data, byte('\n'))
		}
	}

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil
	}

	return result
}

func TestViper(t *testing.T) {
	config := viper.New()

	config.Set("emit", []string{"D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go"})
	config.Set("linter", map[string]any{
		"extends": []string{"@azure-tools/typespec-azure-resource-manager/all"},
	})

	err := config.SafeWriteConfigAs("viper.yaml")
	if err != nil {
		t.Fatal(err)
	}

	config.WriteConfig()
}
