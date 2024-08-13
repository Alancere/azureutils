/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package tsp

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Alancere/azureutils/common"
	"github.com/Masterminds/semver/v3"
	"github.com/fatih/color"
	jsoniter "github.com/json-iterator/go"
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

		up, err := cmd.Flags().GetBool("upgrade")
		if err != nil {
			return err
		}

		syncGoOptionPath, err := cmd.Flags().GetString("sync")
		if err != nil {
			return err
		}

		return generateLockFile(repoRoot, up, syncGoOptionPath)
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
	lockCmd.Flags().BoolP("upgrade", "u", false, "Upgrade npm packages to latest version(upgrade fails when sync exists)")
	lockCmd.Flags().StringP("sync", "s", "", "Sync the lock file to the @azure-tools/typespec-go/package.json")
}

const (
	npm_package      = "package.json"
	npm_package_lock = "package-lock.json"

	emitter_package      = "emitter-package.json"
	emitter_package_lock = "emitter-package-lock.json"
)

func generateLockFile(repoRoot string, upgrade bool, goOptionPath string) error {
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

	repoRoot, err := filepath.Abs(repoRoot)
	if err != nil {
		return err
	}

	// create temp directory
	tempRoot := createTempDirectory(repoRoot)
	defer os.RemoveAll(tempRoot)

	// copy emitter-package to package.json
	if err = copyFile(filepath.Join(repoRoot, "eng", emitter_package), filepath.Join(tempRoot, npm_package)); err != nil {
		return err
	}

	if len(goOptionPath) > 0 {
		if err = syncPackageJson(goOptionPath, filepath.Join(tempRoot, npm_package)); err != nil {
			return err
		}
	} else if upgrade {
		// npm-check-updates -u
		if err = common.Npx(tempRoot, "npm-check-updates", "-u"); err != nil {
			return err
		}
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
		if err = copyFile(lockPath, filepath.Join(repoRoot, "eng", emitter_package_lock)); err != nil {
			return err
		}

		// copy package.json to emitter-package.json
		if len(goOptionPath) > 0 || upgrade {
			if err = copyFile(filepath.Join(tempRoot, npm_package), filepath.Join(repoRoot, "eng", emitter_package)); err != nil {
				return err
			}
		}
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

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err = os.WriteFile(dst, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

func syncPackageJson(src, dst string) error {
	srcData, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	dstData, err := os.ReadFile(dst)
	if err != nil {
		return err
	}

	syncDeps := typespec_go_package_deps(srcData)
	if len(syncDeps) == 0 {
		return err
	}

	newData := make([]byte, len(dstData))
	copy(newData, dstData)

	dstDeps := subKeys(dstData, "devDependencies")
	for _, dep := range dstDeps {
		d := json.Get(dstData, "devDependencies", dep).ToString()
		v1, err := semver.NewConstraint(d)
		if err != nil {
			return err
		}

		if syncDep, ok := syncDeps[dep]; ok {
			v2, err := semver.NewConstraint(syncDep)
			if err != nil {
				return err
			}

			if v1.String() == v2.String() {
				continue
			}

			color.Green("replaced version:%s --> %s\n", v1.String(), v2.String())
			newData = bytes.ReplaceAll(newData, []byte(fmt.Sprintf("\"%s\": \"%s\"", dep, v1.String())), []byte(fmt.Sprintf("\"%s\": \"%s\"", dep, v2.String())))
		}
	}

	if !bytes.Equal(newData, dstData) {
		return os.WriteFile(dst, newData, 0o644)
	}

	return nil
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func subKeys(data []byte, key string) []string {
	k := json.Get(data, key)
	if k.Size() == 0 {
		return nil
	}

	return k.Keys()
}

func typespec_go_package_deps(data []byte) map[string]string {
	deps := subKeys(data, "dependencies")

	m := make(map[string]string, len(deps))
	for _, dep := range deps {
		m[dep] = json.Get(data, "dependencies", dep).ToString()
	}

	devDeps := subKeys(data, "devDependencies")
	for _, dep := range devDeps {
		m[dep] = json.Get(data, "devDependencies", dep).ToString()
	}

	return m
}
