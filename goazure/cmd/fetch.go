/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Alancere/azureutils/fetchmodule"
	"github.com/spf13/cobra"
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch release from pkg.go.dev",
	Long:  `release时执行，目的是为了快速将released的package出现在pkg.go.dev上`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("fetch called")

		module := os.Args[len(os.Args)-1]
		b, err := fetchmodule.Validate(module)
		if !b {
			return err
		}

		firstRelease, _ := cmd.Flags().GetBool("first")
		if firstRelease {
			before, _, _ := strings.Cut(module, "@")
			err = fetchmodule.Fetch(before)
			if err != nil {
				return err
			}
			return fetchmodule.Info(module)
		} else {
			return fetchmodule.NewPackages(module)
		}
	},
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fetchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// fetchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	fetchCmd.Flags().BoolP("first", "f", false, "fetch first release")
}
