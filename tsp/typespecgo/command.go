package typespecgo

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

func Go(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir

	output, err := cmd.CombinedOutput()
	fmt.Printf("###Command: %s\ngo %s\n%s", cmd.Dir, strings.Join(args, " "), string(output))
	if err != nil {
		return err
	}

	return nil
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