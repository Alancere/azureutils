/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package generate

import (
	"github.com/spf13/cobra"
)

// GenerateCmd represents the generate command
var GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "autorest generate sdk code",
	Long: `github.com/Azure/azure-sdk-for-go/tools/generator 包装`,
	RunE: func(cmd *cobra.Command, args []string) error {

		return nil
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// GenerateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// GenerateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	GenerateCmd.PersistentFlags().StringP("spec", "s", "", "spec file path")
}
