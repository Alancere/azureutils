/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

// operationCmd represents the operation command
var operationCmd = &cobra.Command{
	Use:   "operation",
	Short: "goazure spec operation []",
	Long:  `Calculate the coverage of operation`,
	// Args: cobra.RangeArgs(2, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 2 {
			return compareOperationId(args[0], args[1])
		} else {
			fmt.Println("operationId liveTestPath(service/armservice) scenarioPath(scenarios)")
			return nil
		}
	},
}

func init() {
	SpecCmd.AddCommand(operationCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// operationCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// operationCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const (
	CALLOPERATION = "Call operation: "
)

func compareOperationId(liveTestPath, scenarioPath string) error {
	// read liveTestPath file operationId
	liveTest := make(map[string]bool)
	err := filepath.WalkDir(liveTestPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.Contains(info.Name(), "live_test.go") {
			fData, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			for _, line := range strings.Split(string(fData), "\n") {
				if strings.Contains(line, CALLOPERATION) {
					// save to map
					operationId := strings.TrimSuffix(line[strings.Index(line, CALLOPERATION)+len(CALLOPERATION):], "\")")
					if _, ok := liveTest[operationId]; !ok {
						liveTest[operationId] = true
					}
				}
			}
		}

		return nil
	})
	if err != nil {
		return err
	}
	// fmt.Println("live test operation coverage:", liveTest)

	// read scenarioPath file operationId
	scenario := make(map[string]bool)
	fData, err := os.ReadFile(filepath.Join(scenarioPath, "basic.yaml"))
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(fData), "\n") {
		if strings.Contains(line, "- operationId:") {
			_, after, _ := strings.Cut(line, "operationId: ")
			if _, ok := scenario[after]; !ok {
				scenario[after] = true
			}
		}

		if strings.Contains(line, "- step:") {
			_, after, _ := strings.Cut(line, "step: ")
			if _, ok := scenario[after]; !ok {
				scenario[after] = true
			}
		}
	}
	// fmt.Println("scenario operation coverage:", scenario)

	// 对比live test和scenario
	diff := make([]string, 0, len(liveTest))
	for k := range liveTest {
		if _, ok := scenario[k]; ok {
			scenario[k] = false
		}
	}

	// 输出live test没有覆盖到的operationId
	for k, v := range scenario {
		if v {
			diff = append(diff, k)
		}
	}

	// sort diff
	sort.Strings(diff)
	for _, v := range diff {
		fmt.Println(v)
	}

	fmt.Println("live test:", len(liveTest))
	fmt.Println("scenario:", len(scenario))
	fmt.Println("live test not cover operationId:", len(diff))
	fmt.Printf("coverage:%2.2f(%d/%d)", (float64(len(scenario)-len(diff)))/float64(len(scenario)), len(scenario)-len(diff), len(scenario))
	return nil
}
