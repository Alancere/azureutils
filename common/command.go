package common

import (
	"fmt"
	"os/exec"
	"strings"
)

func TSP(dir string, args ...string) (string, error) {
	cmd := exec.Command("tsp", args...)
	cmd.Dir = dir

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("###Command: %s\ntsp %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)
	if err != nil {
		return output, err
	}

	return output, nil
}

func Go(dir string, args ...string) (string, error) {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("###Command: %s\ngo %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)

	return output, err
}

func GoFmt(dir string, args ...string) error {
	cmd := exec.Command("gofmt", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	fmt.Printf("###Command: %s\ngofmt %s\n%s", cmd.Dir, strings.Join(args, " "), string(output))
	if err != nil {
		return err
	}

	return nil
}

func GoFumpt(dir string, args ...string) error {
	cmd := exec.Command("gofumpt", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	fmt.Printf("###Command: %s\ngofumpt %s\n%s", cmd.Dir, strings.Join(args, " "), string(output))
	if err != nil {
		return err
	}

	return nil
}

func AutorestCmd(workspace string, args ...string) (string, error) {
	cmd := exec.Command("autorest", args...)
	cmd.Dir = workspace

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("### %s\nautorest %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)

	return output, err
}

func GoImports(dir string, args ...string) error {
	cmd := exec.Command("goimports", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	fmt.Printf("###Command: %s\ngoimports %s\n%s", cmd.Dir, strings.Join(args, " "), string(output))
	if err != nil {
		return err
	}

	return nil
}

func Generate(dir string, args ...string) (string, error) {
	// cmd := exec.Command("generator", args...)
	cmd := exec.Command("tsp-sdk", args...)
	cmd.Dir = dir

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("###Command: %s\ngenerator %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)

	return output, err
}

func TestProxy(dir string, args ...string) (string, error) {
	cmd := exec.Command("test-proxy", args...)
	cmd.Dir = dir

	combinedOutput, err := cmd.CombinedOutput()
	output := fmt.Sprintf("###Command: %s\ntest-proxy %s\n%s", cmd.Dir, strings.Join(args, " "), string(combinedOutput))
	fmt.Println(output)

	return output, err
}
