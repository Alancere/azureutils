package typespecgo

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/mergego"
)

func TestOneTSP(t *testing.T) {
	specDir := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs"
	sdkDir := "D:/Go/src/github.com/Azure/dev/azure-sdk-for-go"
	typespecgoEmit := "D:/Go/src/github.com/Azure/autorest.go/packages/typespec-go"

	// 用于设置是否使用autorest生成go sdk
	autorestGenerate := true
	autorestGenerate = false

	// tspconfig.yaml 绝对路径
	configPath := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs/specification/azurefleet/AzureFleet.Management/tspconfig.yaml"

	// read readme.go.md
	readmeGOMD := readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/readme.go.md"))
	if strings.Contains(configPath, "Workloads.SAPDiscoverySite.Management") {
		readmeGOMD = readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/Microsoft.Workloads/SAPDiscoverySites/readme.go.md"))
	} else if strings.Contains(configPath, "Workloads.SAPMonitor.Management") {
		readmeGOMD = readmegomd(filepath.Join(filepath.Dir(configPath), "../resource-manager/Microsoft.Workloads/monitors/readme.go.md"))
	}
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

	typespecgoOption := map[string]any{
		"module":                    module,
		"module-version":            moduleVersion,
		"emitter-output-dir":        fmt.Sprintf("{project-root}/go/%s", moduleName),
		"generate-fakes":            true,
		"head-as-boolean":           true, // head method
		"inject-spans":              true,
		"remove-unreferenced-types": true,
	}

	// typespce-go stutter
	if strings.Contains(configPath, "azurelargeinstance") {
		typespecgoOption["stutter"] = "azurelargeinstance"
	}

	if strings.Contains(configPath, "loadtestservice") {
		typespecgoOption["stutter"] = "loadtestservice"
	}

	if strings.Contains(configPath, "mongocluster") {
		typespecgoOption["stutter"] = "DocumentDB"
	}

	if strings.Contains(configPath, "mpcnetworkfunction") {
		typespecgoOption["stutter"] = "MobilePacketCore"
	}

	if strings.Contains(configPath, "playwrighttesting") {
		typespecgoOption["stutter"] = "AzurePlaywrightService"
	}

	if strings.Contains(configPath, "sphere") {
		typespecgoOption["stutter"] = "AzureSphere"
	}

	if strings.Contains(configPath, "codesigning") { // !!!
		typespecgoOption["stutter"] = "CodeSigning"
	}

	if strings.Contains(configPath, "azurefleet") {
		typespecgoOption["stutter"] = "AzureFleet"
	}

	tspConfig.OnlyEmit(typespecgoEmit)
	tspConfig.EditOptions(string(TypeSpec_GO), typespecgoOption, false)

	err = tspConfig.Write()
	if err != nil {
		t.Fatal(err)
	}

	// tsp install
	TSP(filepath.Dir(configPath), "install")

	tspCompileOpts := [2]string{"compile", "main.tsp"}
	if existClientTSP(filepath.Dir(configPath)) {
		tspCompileOpts[1] = "client.tsp"
	}
	output, tspErr := TSP(filepath.Dir(configPath), tspCompileOpts[0], tspCompileOpts[1])
	if tspErr != nil {
		// write output in file
		if err = os.WriteFile(filepath.Join(configPath, "../error.log"), []byte(output), 0o777); err != nil {
			t.Fatal(err)
		}
	}
	// break

	///
	// go mod tidy and go vet ./...
	// gosdk := fmt.Sprintf("%s/go/%s", filepath.Dir(configPath), moduleName)
	if tspErr == nil {
		tspsdk := filepath.Join(filepath.Dir(configPath), "go", moduleName.(string))
		if err = GoFmt(tspsdk, "-w", "."); err != nil {
			log.Println("##tsp##gofmt ", err)
		}

		if err = Go(tspsdk, "mod", "tidy"); err != nil {
			log.Println("##tsp##go mod", err)
		}

		if err = Go(tspsdk, "vet", "./..."); err != nil {
			log.Println("##tsp##go vet", err)
		} else {
			// merge go files
			if err = mergego.Merge(tspsdk, filepath.Join("D:/tmp/typespecp-diff", filepath.Base(tspsdk)+".go")); err != nil {
				log.Fatal(err)
			}

			// merge fake go files
			if err = mergego.Merge(filepath.Join(tspsdk, "fake"), filepath.Join("D:/tmp/typespecp-diff", filepath.Base(tspsdk)+"_fake.go")); err != nil {
				log.Fatal(err)
			}
		}

		// autorest generate go sdk
		if autorestGenerate {
			autorestsdk := filepath.Join(sdkDir, moduleName.(string))
			serviceName, armServiceName := armName(moduleName.(string))
			specName := filepath.Base(filepath.Join(filepath.Dir(configPath), "../"))
			autorestOps := []string{
				"release-v2",
				sdkDir,
				specDir,
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

			// 需要指定 readme.md 的 onboard service
			if strings.Contains(configPath, "Workloads.SAPDiscoverySite.Management") {
				autorestOps = append(autorestOps, "--specify-readme-path", "specification\\workloads\\resource-manager\\Microsoft.Workloads\\monitors\\readme.md")
			}

			output, err := readmePathExec(filepath.Join(sdkDir, "../"), autorestOps...)
			if err != nil {
				log.Println("##autorest##", err)
			} else {
				if err = Go(autorestsdk, "vet", "./..."); err != nil {
					log.Println("##autorest##go vet", err)
				} else {
					// merge go files
					if err = mergego.Merge(autorestsdk, filepath.Join("D:/tmp/autorest-diff", filepath.Base(tspsdk)+".go")); err != nil {
						log.Fatal(err)
					}

					// merge fake go files
					if err = mergego.Merge(filepath.Join(autorestsdk, "fake"), filepath.Join("D:/tmp/autorest-diff", filepath.Base(tspsdk)+"_fake.go")); err != nil {
						log.Fatal(err)
					}
				}
			}
			_ = output
		}
	}
}

// 等generator tool支持 --specify-readme-path 时 替换
func readmePathExec(dir string, args ...string) (string, error) {
	cmd := exec.Command("readmepath", args...)
	cmd.Dir = dir

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("###Command: %s\ngenerator %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)

	return output, err
}
