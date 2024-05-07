/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/Alancere/azureutils/mergego"
	"github.com/spf13/cobra"
)

// mergoCmd represents the mergo command
var mergoCmd = &cobra.Command{
	Use:   "mergo",
	Short: "mergo dir [outfile]",
	Long: `合并go package
	goazure mergo dir [outfile]`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("please input dir")
		}
		dir := args[0]
		outfile := ""
		if len(args) >= 2 {
			outfile = args[1]
		}
		if err := mergego.Merge(dir, outfile); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(mergoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// mergoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// mergoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
