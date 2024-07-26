/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tsp

import (
	"github.com/Alancere/azureutils/tsp/typespec"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// linkCmd represents the link command
var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "goazure tsp link: return the link of the tsp-location.yaml",
	Long:  ``,
	Args:  cobra.RangeArgs(1, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tspLocationPath := args[0]
		tspLocation, err := typespec.ParseTspLocation(tspLocationPath)
		if err != nil {
			return err
		}
		color.Green("Validation of the link")
		color.Blue(tspLocation.Link().String())
		return tspLocation.Validation()
	},
}

func init() {
	TspCmd.AddCommand(linkCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// linkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// linkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
