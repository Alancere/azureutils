/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tspCmd represents the tsp command
var tspCmd = &cobra.Command{
	Use:   "tsp",
	Short: "about typespec (tsp) command",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("tsp called")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tspCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tspCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tspCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
