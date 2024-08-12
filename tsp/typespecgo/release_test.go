package typespecgo_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/common"
	"github.com/Alancere/azureutils/mergego"
	"github.com/Alancere/azureutils/tsp/typespecgo"
	"github.com/goccy/go-yaml"
)

/*
branch:(sdkrepo specRepo)
	tsp_mgmt_0723 tsp_spec_0723
*/

// test generator release-v2 tsp
func TestGenerateTool_Support_TSP(t *testing.T) {
	// 存放通过typespec-go生成的结果
	rootDir := "D:/Go/src/github.com/Azure"
	specDir := "D:/Go/src/github.com/Azure/azure-rest-api-specs"
	sdkDir := "D:/Go/src/github.com/Azure/azure-sdk-for-go"

	// 存放通过autorset.go生成的结果
	debugSdkDir := "D:/Go/src/github.com/Azure/debug/azure-sdk-for-go"
	debugSpecDir := "D:/Go/src/github.com/Azure/debug/azure-rest-api-specs"
	autorestGenerate := false // 是否通过autorest.go生成
	autorestGenerate = true

	genertorErrs := make([]error, 0)
	goVetErrs := make([]error, 0)

	configPaths, err := typespecgo.SearchTSP(specDir)
	if err != nil {
		t.Fatal(err)
	}

	tspErrs := make([]error, 0)
	allErrMsg := make([]string, 0)
	mgmtTspCount := make([]string, 0)
	dataPlaneTspCount := make([]string, 0)

	successedTSP := make([]string, 0)

	configPaths = []string{
		"D:\\Go\\src\\github.com\\Azure\\azure-rest-api-specs\\specification\\healthdataaiservices\\HealthDataAIServices.Management\\tspconfig.yaml",
		// "D:\\Go\\src\\github.com\\Azure\\azure-rest-api-specs\\specification\\mongocluster\\DocumentDB.MongoCluster.Management\\tspconfig.yaml",
	}

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
			"service-dir":               fmt.Sprintf("sdk/resourcemanager/%s", serviceName),
			"package-dir":               armServiceName,
			"module":                    "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}",
			"examples-directory":        "{project-root}/examples",
			"fix-const-stuttering":      true,
			"flavor":                    "azure",
			"generate-examples":         true,
			"generate-fakes":            true,
			"head-as-boolean":           true,
			"inject-spans":              true,
			"remove-unreferenced-types": true,
		}

		// typespce-go stutter
		// stutter(configPath, typespecgoOption)

		// tspConfig.OnlyEmit(typespecgoEmit)
		tspConfig.EditOptions(string(typespecgo.TypeSpec_GO), typespecgoOption, false)
		err = tspConfig.Write()
		if err != nil {
			t.Fatal(err)
		}

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
			"--tsp-client-option", "--debug,--save-inputs",
		)
		// go mod tidy and go vet ./...
		if err != nil {
			log.Printf("tsp-sdk error###%s %s\n", tspConfig.Path, output)
			genertorErrs = append(genertorErrs, err)
			continue
		}

		tspsdk := filepath.Join(sdkDir, "sdk", "resourcemanager", serviceName, armServiceName)
		output, err = common.Go(tspsdk, "vet", "./...")
		if err != nil {
			log.Println("govet###", output)
			goVetErrs = append(goVetErrs, err)
			continue
		} else {
			// merge go files
			if err = mergego.Merge(tspsdk, filepath.Join("D:/tmp/typespec-X-diff", filepath.Base(tspsdk)+".go"), false); err != nil {
				log.Fatal(err)
			}

			// merge fake go files
			if err = mergego.Merge(filepath.Join(tspsdk, "fake"), filepath.Join("D:/tmp/typespec-X-diff", filepath.Base(tspsdk)+"_fake.go"), false); err != nil {
				log.Fatal(err)
			}
		}
		_ = output

		successedTSP = append(successedTSP, configPath)
		// break

		// autorest generate go sdk
		if autorestGenerate {
			autorestsdk := filepath.Join(debugSdkDir, moduleName.(string))
			serviceName, armServiceName := armName(moduleName.(string))
			specName := filepath.Base(filepath.Join(filepath.Dir(configPath), "../"))
			autorestOps := []string{
				"release-v2",
				debugSdkDir,
				debugSpecDir,
				serviceName, armServiceName,
				fmt.Sprintf("--spec-rp-name=%s", specName),
				"--skip-generate-example",
				"--skip-create-branch",
			}
			defaultTag, err := readmemd(filepath.Join(filepath.Dir(configPath), "../resource-manager/readme.md"))
			if err != nil {
				log.Println(err)
			}
			if defaultTag != "" {
				autorestOps = append(autorestOps, fmt.Sprintf("--package-config=%s", strings.TrimSpace(defaultTag)))
			}
			output, err := common.Generate(filepath.Join(debugSdkDir, "../"), autorestOps...)
			if err != nil {
				log.Println("##autorest##", err)
			} else {
				if _, err = common.Go(autorestsdk, "vet", "./..."); err != nil {
					log.Println("##autorest##go vet", err)
				} else {
					// merge go files
					if err = mergego.Merge(autorestsdk, filepath.Join("D:/tmp/autorest-X-diff", filepath.Base(autorestsdk)+".go"), false); err != nil {
						log.Fatal(err)
					}

					// merge fake go files
					if err = mergego.Merge(filepath.Join(autorestsdk, "fake"), filepath.Join("D:/tmp/autorest-X-diff", filepath.Base(autorestsdk)+"_fake.go"), false); err != nil {
						log.Fatal(err)
					}
				}
			}
			_ = output
		}
	}

	fmt.Println("successed tsp count:", len(successedTSP))
	for _, s := range successedTSP {
		fmt.Println(s)
	}

	// write error msg to error.log
	fmt.Println("tsp compiler error count:", len(tspErrs))
	errMsg := ""
	for _, eMsg := range allErrMsg {
		errMsg = fmt.Sprintf("%s\n%s", errMsg, eMsg)
	}

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

func readmemd(path string) (string, error) {
	md, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	for _, l := range strings.Split(string(md), "\n") {
		// get first tag
		if strings.Contains(l, "tag:") {
			return l, nil
		}
	}

	return "", nil
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
