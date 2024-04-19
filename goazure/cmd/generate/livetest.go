/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package generate

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// livetestCmd represents the livetest command
var livetestCmd = &cobra.Command{
	Use:   "livetest",
	Short: "",
	Long:  ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return fmt.Errorf("please input need generate livetest folder or file absolute path")
		}

		f, err := os.Stat(args[0])
		if err != nil {
			return err
		}

		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		execPath, err := cmd.Flags().GetString("execPath")
		if err != nil {
			return err
		}
		if execPath == "" {
			execPath, err = os.UserHomeDir()
			if err != nil {
				return err
			}
		}

		// D:\Go\src\github.com\Azure\azure-rest-api-specs\specification\apicenter\resource-manager\Microsoft.ApiCenter\stable\2024-03-01\apicenter.json
		// docker run --rm -v D:/Go/src/github.com/Azure/azure-rest-api-specs/specification:/swagger -w /swagger/.restler_output mcr.microsoft.com/restlerfuzzer/restler:v8.6.0 dotnet /RESTler/restler/Restler.dll compile /swagger/Compile/config.json
		dockerScript := `run --rm -v %s:/swagger -w /swagger/.restler_output mcr.microsoft.com/restlerfuzzer/restler dotnet /RESTler/restler/Restler.dll compile %s`
		// get repo path
		before, _, b := strings.Cut(args[0], "specification")
		if !b {
			return fmt.Errorf("please input correct path, exmple: ../azure-rest-api-specs/specification/....")
		}
		// compile
		dst := args[0]
		if f.IsDir() {
			s := fmt.Sprintf(dockerScript, before, "/swagger/Compile/config.json")
			if tag != "" {
				s = fmt.Sprintf("%s --tag %s", s, tag)
			}
			err = DockerCmd(execPath, s)
			if err != nil {
				return err
			}
		} else {
			s := fmt.Sprintf(dockerScript, before, "--api_spec "+args[0])
			if tag != "" {
				s = fmt.Sprintf("%s --tag %s", s, tag)
			}
			err = DockerCmd(execPath, s)
			if err != nil {
				return err
			}

			dst = filepath.Dir(dst)
		}

		// oav
		// oav generate-api-scenario static --readme specification/apicenter/resource-manager/readme.md --dependency specification/.restler_output/Compile/dependencies.json -o specification/apicenter/resource-manager/Microsoft.ApiCenter/stable/2024-03-01/scenarios --useExample
		oavScript := `generate-api-scenario static --readme %s --dependency %s -o %s --useExample`
		err = OavCmd(execPath, fmt.Sprintf(oavScript, filepath.Join(dst, "../../../readme.md"), filepath.Join(before, "specification", ".restler_output/Compile/dependencies.json"), filepath.Join(dst, "scenarios")))
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	GenerateCmd.AddCommand(livetestCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// livetestCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// livetestCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	livetestCmd.Flags().StringP("tag", "", "", "speficy readme.md tag")

	livetestCmd.Flags().StringP("execPath", "", "", "指定执行目录")
}

func DockerCmd(dir string, args ...string) error {
	goExec := exec.Command("docker", args...)
	goExec.Dir = dir

	fmt.Printf("%s execute:\n%s", dir, goExec.String())
	output, err := goExec.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("docker run error: %s\n%s", err, string(output))
	}

	return nil
}

func OavCmd(dir string, args ...string) error {
	goExec := exec.Command("oav", args...)
	goExec.Dir = dir

	fmt.Printf("%s execute:\n%s", dir, goExec.String())
	output, err := goExec.CombinedOutput()
	fmt.Println(string(output))
	if err != nil {
		return fmt.Errorf("oav run error: %s\n%s", err, string(output))
	}

	return nil
}
