/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Alancere/azureutils/tsp/utils"
	"github.com/spf13/cobra"
)

// swagCmd represents the swag command
var swagCmd = &cobra.Command{
	Use:   "swag",
	Short: "compare swagger file",
	Long: `go azure tsp swag originDir compiledFile
	or goazure tsp swag --originDir originDir --compiledFile compiledFile`,
	RunE: func(cmd *cobra.Command, args []string) error {
		l := len(args)
		if l == 1 {
			outputDir = args[0]
		}else if l == 2 {
			originDir = args[0]
			compiledFile = args[1]
		}else if l == 3 {
			originDir = args[0]
			compiledFile = args[1]
			outputDir = args[2]
		}

		if outputDir == "" {
			outputDir = "."
		}else {
			// create output dir
			if err := os.Mkdir(outputDir, 0o666); err != nil {
				return  err
			}
		}

		if originDir == "" || compiledFile == "" {
			return fmt.Errorf("please input args: originDir and compileFile")
		}
		mergePath := filepath.Join(outputDir, "merge.json")
		if err := utils.MergeJson(originDir, mergePath); err != nil {
			return err
		}

		formatPath := filepath.Join(outputDir, "format.json")
		if err := utils.FormatJson(compiledFile, formatPath); err != nil {
			return err
		}

		if err := utils.ComparePath(mergePath, formatPath, filepath.Join(outputDir, "paths.md")); err != nil {
			return err
		}

		return nil
	},
}

var (
	originDir string
	compiledFile string
	outputDir string
)

func init() {
	tspCmd.AddCommand(swagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// swagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// swagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	swagCmd.Flags().StringVarP(&originDir, "originDir", "", "", "origin dir")
	swagCmd.Flags().StringVarP(&compiledFile, "compiledFile", "", "", "compile file")
	// swagCmd.Flags().StringVarP(&outputDir, "outputDir", "", "", "output dir, option")
}
