/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package generate

import (
	"fmt"

	"github.com/spf13/cobra"
)

var tspConfigPath string

func init() {
	GenerateCmd.AddCommand(typespecCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// typespecCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// typespecCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	typespecCmd.Flags().StringVarP(&tspConfigPath, "tsp-config", "", "", "specify tspconfig.yaml file path")
}

// typespecCmd represents the typespec command
var typespecCmd = &cobra.Command{
	Use:   "typespec",
	Short: "generate go sdk from typespec",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		/*
			typespec ./sdkRepo servive armService --tsp-config(required) ./tspconfig.yaml
		*/
		fmt.Println("typespec called")
		if tspConfigPath == "" {
			return fmt.Errorf("please input tspconfig.yaml file path")
		}

		serviceName := args[1]
		armServiceName := fmt.Sprintf("arm%s", serviceName)
		if len(args) == 3 {
			armServiceName = args[2]
		}

		return nil
	},
}
