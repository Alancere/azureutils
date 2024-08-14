/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cadlranch

import (
	"github.com/spf13/cobra"

	cadl_ranch "github.com/Alancere/azureutils/cadl-ranch"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "generate cadl-ranch test",
	Long:  `goazure cadlranch test []`,
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cadlRanchSpecsPath, err := cmd.Flags().GetString("cadl-ranch-specs")
		if err != nil {
			return err
		}

		if len(args) == 2 {
			cadlRanchSpecsPath = args[1]
		}

		return cadl_ranch.GenerateCadlRanchTest(args[0], cadlRanchSpecsPath)
	},
}

func init() {
	CadlRanchCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	testCmd.Flags().StringP("cadl-ranch-specs", "t", "", "cadl-ranch-specs path")
}
