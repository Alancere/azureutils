/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package spec

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// scenarioCmd represents the scenario command
var scenarioCmd = &cobra.Command{
	Use:   "scenario",
	Short: "goazure spec scenario []",
	Long:  `Generate basic scenarios files`,
	// Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("scenario called")
		if len(args) == 1 {
			return generateScenariosFile(os.Args[0])
		} else {
			fmt.Println("please input path(github.com/Azure/azure-rest-api-specs/specification/service/resource-manager/Microsoft/stable/version)....")
			return nil
		}
	},
}

func init() {
	SpecCmd.AddCommand(scenarioCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// scenarioCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scenarioCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const SCENARIOSFORMAT = `# yaml-language-server: $schema=https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/documentation/api-scenario/references/v1.2/schema.json
scope: ResourceGroup

scenarios:
  - scenario: 
    description: 
    steps:

`

func generateScenariosFile(versionPath string) error {
	count := 0
	if err := filepath.Walk(versionPath, func(path string, info fs.FileInfo, err error) error {
		// read all file in the current directory
		if !info.IsDir() && !strings.Contains(path, "examples") && !strings.Contains(path, "scenarios") {
			count++
			fmt.Println(info.Name(), path)

			data, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			jsonObject := make(map[string]any)
			if err = json.Unmarshal(data, &jsonObject); err != nil {
				return err
			}

			if _, ok := jsonObject["paths"]; ok {
				scenarioPath := filepath.Join(versionPath, "scenarios")
				if _, err = os.Open(scenarioPath); err != nil {
					// create scenarios dirctory
					if err = os.Mkdir(scenarioPath, 0o666); err != nil {
						return err
					}
				}

				// generate api scenarios file
				b, _, _ := strings.Cut(info.Name(), ".json")
				scenarioFile, err := os.Create(filepath.Join(scenarioPath, fmt.Sprintf("%s.yaml", b)))
				if err != nil {
					return err
				}
				scenarioFile.WriteString(SCENARIOSFORMAT)
				scenarioFile.Close()
			}
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println("count:", count)
	return nil
}
