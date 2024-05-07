package typespecgo

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/mergego"
	"github.com/goccy/go-yaml"
	"github.com/spf13/viper"
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
	mgmtTspCount := make([]string, 0)
	dataPlaneTspCount := make([]string, 0)
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
		versions := []string{"0.1.0", "1.0.0", "2.0.0", "30.0.0", "0.5.0-beta.1", "2.2.0-beta.2"}
		moduleVersion = versions[rand.Intn(6)]
		
		// tsp compile 之前把go目录和error.log删除
		gosdk := filepath.Join(filepath.Dir(configPath), "go", moduleName.(string))
		os.RemoveAll(gosdk)
		os.Remove(filepath.Join(configPath, "../error.log"))

		// update tspconfig
		tspConfig, err := NewTSPConfig(configPath)
		if err != nil {
			t.Fatal(err)
		}

		// 过滤 data-plane
		if v, ok := tspConfig.TypeSpecProjectSchema.Options["@azure-tools/typespec-autorest"]; ok {
			if pro, ok := v.(map[string]any)["azure-resource-provider-folder"]; ok {
				if strings.Contains(pro.(string), "data-plane") {
					dataPlaneTspCount = append(dataPlaneTspCount, configPath)
					continue
				} else if strings.Contains(pro.(string), "resource-manager") {
					mgmtTspCount = append(mgmtTspCount, configPath)
				}
			}
		} else {
			fmt.Println("not found @azure-tools/typespec-autorest option:", configPath)
		}

		typespecgoOption := map[string]any{
			"module":             module,
			"module-version":     moduleVersion,
			"emitter-output-dir": fmt.Sprintf("{project-root}/go/%s", moduleName),
			"generate-fakes":     true,
			"head-as-boolean":    true, // head method
		}

		tspConfig.OnlyEmit(typespecgoEmit)
		tspConfig.EditOptions(string(TypeSpec_GO), typespecgoOption, false)

		err = tspConfig.Write()
		if err != nil {
			t.Fatal(err)
		}

		// tsp install
		TSP(filepath.Dir(configPath), "install")

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
			if err = GoFmt(gosdk, "-w", "."); err != nil {
				log.Println("####gofmt ", err)
			}

			if err = Go(gosdk, "mod", "tidy"); err != nil {
				log.Println("####go mod", err)
			}

			if err = Go(gosdk, "vet", "./..."); err != nil {
				log.Println("####go vet", err)
			}else {
				// merge go files
				if err = mergego.Merge(gosdk, filepath.Join("D:/tmp/typespecp-diff", filepath.Base(gosdk) + ".go")); err != nil {
					log.Fatal(err)
				}

				// merge fake go files
				if err = mergego.Merge(filepath.Join(gosdk, "fake"), filepath.Join("D:/tmp/typespecp-diff", filepath.Base(gosdk) + "_fake.go")); err != nil {
					log.Fatal(err)
				}
			}
		}
		// break
	}

	// fmt.Println(tspErrs)

	// write error msg to error.log
	fmt.Println("tsp compiler error count:", len(tspErrs))
	errMsg := ""
	for _, eMsg := range allErrMsg {
		errMsg = fmt.Sprintf("%s\n%s", errMsg, eMsg)
	}
	// if err = os.WriteFile(filepath.Join(dir, "tspCompiler.log"), []byte(errMsg), 0777); err != nil {
	//   t.Fatal(err)
	// }

	fmt.Println("tsp resource management count:", len(mgmtTspCount))
	fmt.Println("tsp data-plane count:", len(dataPlaneTspCount))
}

func TestGeneratePrivateSDK(t *testing.T) {
	dir := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs-pr"
	typespecgoEmit := "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go"
	// typespecgoEmit = "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go/azure-tools-typespec-go-0.1.0-dev.1.tgz" // 不能用这种方式

	configPaths, err := SearchTSP(dir)
	if err != nil {
		t.Fatal(err)
	}

	tspErrs := make([]error, 0)
	allErrMsg := make([]string, 0)
	mgmtTspCount := make([]string, 0)
	dataPlaneTspCount := make([]string, 0)
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

		// 移除tsp中的import "@azure-tools/typespec-providerhub" deprecated
		// if !strings.Contains(configPath, "Community.Management") {
		// 	continue
		// }
		err = removeImport(filepath.Dir(configPath))
		if err != nil {
			t.Fatal(err)
		}
		// return

		// read readme.go.md
		readmeGOMD := readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/readme.go.md"))
		module := readmeGOMD["module"]
		moduleName := readmeGOMD["module-name"]
		module = strings.Replace(module.(string), "$(module-name)", moduleName.(string), -1)
		moduleVersion := "0.1.0" // default value, need from autorest.md get
		// 随机设置module version
		versions := []string{"0.1.0", "1.0.0", "2.0.0", "30.0.0", "0.5.0-beta.1", "2.2.0-beta.2"}
		moduleVersion = versions[rand.Intn(6)]

		// tsp compile 之前把go目录和error.log删除
		gosdk := filepath.Join(filepath.Dir(configPath), "go", moduleName.(string))
		os.RemoveAll(gosdk)
		os.Remove(filepath.Join(configPath, "../error.log"))

		// update tspconfig
		tspConfig, err := NewTSPConfig(configPath)
		if err != nil {
			t.Fatal(err)
		}

		// 过滤 data-plane
		if v, ok := tspConfig.TypeSpecProjectSchema.Options["@azure-tools/typespec-autorest"]; ok {
			if pro, ok := v.(map[string]any)["azure-resource-provider-folder"]; ok {
				if strings.Contains(pro.(string), "data-plane") {
					dataPlaneTspCount = append(dataPlaneTspCount, configPath)
					continue
				} else if strings.Contains(pro.(string), "resource-manager") {
					mgmtTspCount = append(mgmtTspCount, configPath)
				}
			}
		} else {
			fmt.Println("not found @azure-tools/typespec-autorest option:", configPath)
		}

		typespecgoOption := map[string]any{
			"module":             module,
			"module-version":     moduleVersion,
			"emitter-output-dir": fmt.Sprintf("{project-root}/go/%s", moduleName),
			"generate-fakes":     true,
			"head-as-boolean":    true, // head method
		}

		tspConfig.OnlyEmit(typespecgoEmit)
		tspConfig.EditOptions(string(TypeSpec_GO), typespecgoOption, false)

		err = tspConfig.Write()
		if err != nil {
			t.Fatal(err)
		}

		// tsp install
		TSP(filepath.Dir(configPath), "install")

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
			if err = GoFmt(gosdk, "-w", "."); err != nil {
				log.Println("####gofmt ", err)
			}

			if err = Go(gosdk, "mod", "tidy"); err != nil {
				log.Println("####go mod", err)
			}

			if err = Go(gosdk, "vet", "./..."); err != nil {
				log.Println("####go vet", err)
			}else {
				// merge go files
				if err = mergego.Merge(gosdk, filepath.Join("D:/tmp/typespecp-diff-pr", filepath.Base(gosdk) + ".go")); err != nil {
					log.Fatal(err)
				}

				// merge fake go files
				if err = mergego.Merge(filepath.Join(gosdk, "fake"), filepath.Join("D:/tmp/typespecp-diff-pr", filepath.Base(gosdk) + "_fake.go")); err != nil {
					log.Fatal(err)
				}
			}
		}
		// break
	}

	// fmt.Println(tspErrs)

	// write error msg to error.log
	fmt.Println("tsp compiler error count:", len(tspErrs))
	errMsg := ""
	for _, eMsg := range allErrMsg {
		errMsg = fmt.Sprintf("%s\n%s", errMsg, eMsg)
	}
	// if err = os.WriteFile(filepath.Join(dir, "tspCompiler.log"), []byte(errMsg), 0777); err != nil {
	//   t.Fatal(err)
	// }

	fmt.Println("tsp resource management count:", len(mgmtTspCount))
	fmt.Println("tsp data-plane count:", len(dataPlaneTspCount))
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

func removeImport(tspPath string) error {
	return filepath.Walk(tspPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !strings.Contains(info.Name(), ".tsp") {
			return nil
		}
		// fmt.Println(path)

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if !strings.Contains(string(data), "\"@azure-tools/typespec-providerhub\"") {
			return nil
		}

		lines := strings.Split(string(data), "\n")
		newLines := make([]string, 0, len(lines))
		for i, l := range lines {
			if strings.Contains(l, "\"@azure-tools/typespec-providerhub\"") {
				newLines = append(lines[:i], lines[i+1:]...)
				break
			}
		}

		if len(newLines) == len(lines) {
			return nil
		}
		if err := os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0o777); err != nil {
			return err
		}

		return nil
	})
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
