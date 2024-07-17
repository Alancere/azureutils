package typespecgo_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/common"
	"github.com/Alancere/azureutils/tsp/typespecgo"
	"github.com/goccy/go-yaml"
)

// test generator release-v2 tsp
func TestGenerateTool_Support_TSP(t *testing.T) {
	specDir := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs"
	sdkDir := "D:/Go/src/github.com/Azure/dev/azure-sdk-for-go"
	// typespecgoEmit := "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go"
	// typespecgoEmit = "@azure-tools/typespec-go"
	// typespecgoEmit = "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go/azure-tools-typespec-go-0.1.0-dev.1.tgz" // 不能用这种方式

	specDir = "D:/Go/src/github.com/Azure/azure-rest-api-specs"
	sdkDir = "D:/Go/src/github.com/Azure/azure-sdk-for-go"
	rootDir := "D:/Go/src/github.com/Azure"

	genertorErrs := make([]error, 0)
	goVetErrs := make([]error, 0)

	// 用于设置是否使用autorest生成go sdk
	// autorestGenerate := true
	// autorestGenerate = false

	configPaths, err := typespecgo.SearchTSP(specDir)
	if err != nil {
		t.Fatal(err)
	}

	tspErrs := make([]error, 0)
	allErrMsg := make([]string, 0)
	mgmtTspCount := make([]string, 0)
	dataPlaneTspCount := make([]string, 0)
	// tspCompilerLog, err := os.OpenFile(filepath.Join(specDir, "tspCompiler.log"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer tspCompilerLog.Close()

	// bufLog := bufio.NewWriter(tspCompilerLog)
	// defer bufLog.Flush()

	for _, configPath := range configPaths {
		// filter
		if strings.Contains(configPath, "machinelearning\\Azure.AI.ChatProtocol") { // 没有go track2 config
			continue
		}

		if strings.Contains(configPath, "azurestackhci\\Operations.Management") { // azurestackhci的子service
			continue
		}

		if strings.Contains(configPath, "monitor\\Microsoft.Monitor") { // monitor 改动较大
			continue
		}

		// read readme.go.md
		readmeGOMD := readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/readme.go.md"))

		// deep readme.md
		deepReamd(configPath, readmeGOMD)

		module := readmeGOMD["module"]
		moduleName := readmeGOMD["module-name"]
		module = strings.Replace(module.(string), "$(module-name)", moduleName.(string), -1)
		// moduleVersion := "0.1.0" // default value, need from autorest.md get
		// versions := []string{"0.1.0", "1.0.0", "2.0.0", "30.0.0", "0.5.0-beta.1", "2.2.0-beta.2"}
		// moduleVersion = versions[rand.Intn(6)]

		// tsp compile 之前把go目录和error.log删除
		// gosdk := filepath.Join(filepath.Dir(configPath), "go", moduleName.(string))
		// os.RemoveAll(gosdk)
		// os.Remove(filepath.Join(configPath, "../error.log"))

		// update tspconfig
		tspConfig, err := typespecgo.NewTSPConfig(configPath)
		if err != nil {
			t.Fatal(err)
		}

		// 需要额外设置 additionalDirectories
		// Oracle, azurestackhci
		if strings.Contains(configPath, "Oracle.Database") {
			if tspConfig.Parameters == nil {
				tspConfig.Parameters = map[string]any{}
			}
			tspConfig.Parameters["dependencies"] = map[string]any{
				"default":               "",
				"additionalDirectories": []string{"specification/oracle/models"},
			}
		} else if strings.Contains(configPath, "AzureStackHCI.StackHCIVM.Management") {
			if tspConfig.Parameters == nil {
				tspConfig.Parameters = map[string]any{}
			}
			tspConfig.Parameters["dependencies"] = map[string]any{
				"default":               "",
				"additionalDirectories": []string{"specification/azurestackhci/Operations.Management"}, // Operations.Management应该也要生成package.json
			}
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

		fmt.Println("exec ###", configPath)

		serviceName, armServiceName := armName(moduleName.(string))
		typespecgoOption := map[string]any{
			// "module": module,
			// "module-version":            moduleVersion,
			// "emitter-output-dir":        fmt.Sprintf("{project-root}/go/%s", moduleName),
			"generate-fakes":            true,
			"head-as-boolean":           true, // head method
			"inject-spans":              true,
			"remove-unreferenced-types": true,

			"service-dir": "sdk",
			"package-dir": fmt.Sprintf("resourcemanager/%s/%s", serviceName, armServiceName),
			"module": "{service-dir}/{package-dir}",
			"examples-directory": "./examples",
		}

		// typespce-go stutter
		// stutter(configPath, typespecgoOption)

		// tspConfig.OnlyEmit(typespecgoEmit)
		tspConfig.EditOptions(string(typespecgo.TypeSpec_GO), typespecgoOption, false)
		err = tspConfig.Write()
		if err != nil {
			t.Fatal(err)
		}

		// tsp install
		// typespecgo.TSP(filepath.Dir(configPath), "install")

		// exec genertor release-v2 --tsp
		output, err := common.Generate(
			rootDir,
			"release-v2",
			sdkDir,
			specDir,
			serviceName, armServiceName,
			"--tsp-config",
			tspConfig.Path,
			"--skip-create-branch=true",
			"--tsp-client-option", "--debug",
		)
		// go mod tidy and go vet ./...
		if err != nil {
			log.Printf("tsp-sdk###%s %s\n", tspConfig.Path, output)
			genertorErrs = append(genertorErrs, err)
			continue
		}

		tspsdk := filepath.Join(sdkDir, "sdk", "resourcemanager", serviceName, armServiceName)
		output, err = common.Go(tspsdk, "vet", "./...")
		if err != nil {
			log.Println("govet###", output)
			goVetErrs = append(goVetErrs, err)
			continue
		}
		_ = output

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

	fmt.Println("genertor tsp errors:", len(genertorErrs))
	for _, err := range genertorErrs {
		fmt.Println(err)
	}

	fmt.Println("go vet errors:", len(goVetErrs))
	for _, err := range goVetErrs {
		fmt.Println(err)
	}
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

func deepReamd(configPath string, readmeGOMD map[string]any) {
	if strings.Contains(configPath, "Workloads.SAPDiscoverySite.Management") {
		readmeGOMD = readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/Microsoft.Workloads/SAPDiscoverySites/readme.go.md"))
	} else if strings.Contains(configPath, "Workloads.SAPMonitor.Management") {
		readmeGOMD = readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/Microsoft.Workloads/monitors/readme.go.md"))
	}
}

func armName(module string) (string, string) {
	a, _ := strings.CutPrefix(module, "sdk/resourcemanager/")
	b, af, _ := strings.Cut(a, "/")
	return b, af
}
