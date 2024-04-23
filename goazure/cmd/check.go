/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Alancere/azureutils/spec"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "check azure specs",
	Long: ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("### check %s ###", spec.PublicRepository)
		err := spec.NewPullRequest().Check(spec.PublicRepository)
		if err != nil {
			return err
		}
	
		fmt.Printf("### check %s ###", spec.PrivateRepository)
		err = spec.NewPullRequest().Check(spec.PrivateRepository)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
