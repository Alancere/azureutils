/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tsp

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Alancere/azureutils/common"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Generating lock file. goazure tsp lock [repoRoot]",
	Long:  `exec: sdk repo root(azure-sdk-for-go).tsp-client --generate-lock-file`,
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		repoRoot := args[0]
		return generateLockFile(repoRoot, repoRoot)
	},
}

func init() {
	TspCmd.AddCommand(lockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const (
	npm_package      = "package.json"
	npm_package_lock = "package-lock.json"

	emitter_package      = "emitter-package.json"
	emitter_package_lock = "emitter-package-lock.json"
)

func generateLockFile(rootUrl, repoRoot string) error {
	fmt.Println("Generating lock file...")
	args := []string{"install"}
	froceInstall := os.Getenv("TSPCLIENT_FORCE_INSTALL")
	if froceInstall != "" {
		force, err := strconv.ParseBool(froceInstall)
		if err != nil {
			return err
		}
		if force {
			args = append(args, "--force")
		}
	}

	// create temp directory
	tempRoot := createTempDirectory(rootUrl)

	// copy emitter-package to package.json
	data, err := os.ReadFile(filepath.Join(repoRoot, "eng", emitter_package))
	if err != nil {
		return err
	}

	if err = os.WriteFile(filepath.Join(tempRoot, npm_package), data, os.ModePerm); err != nil {
		return err
	}

	// npm install
	if err = common.Npm(tempRoot, args...); err != nil {
		return err
	}

	// lock file
	lockPath := filepath.Join(tempRoot, npm_package_lock)
	lockFile, err := os.Stat(lockPath)
	if err != nil {
		return err
	}
	if !lockFile.IsDir() {
		// copy package-lock.json to emitter-package-lock.json
		data, err := os.ReadFile(lockPath)
		if err != nil {
			return err
		}
		if err = os.WriteFile(filepath.Join(repoRoot, "eng", emitter_package_lock), data, os.ModePerm); err != nil {
			return err
		}
	}

	// remove temp directory
	if err = os.RemoveAll(tempRoot); err != nil {
		return err
	}

	fmt.Println("Lock file generated in", filepath.Join(repoRoot, "eng", emitter_package_lock))
	return nil
}

func createTempDirectory(outputDir string) string {
	tempRoot := filepath.Join(outputDir, "TempTypeSpecFiles")
	os.MkdirAll(tempRoot, os.ModePerm)
	fmt.Println("Creating temporary working directory", tempRoot)
	return tempRoot
}
