package mgmt

import (
	"fmt"
	"os/exec"
)

func ExecGo(path string, pameters ...string) error {
	goExec := exec.Command("go", pameters...)
	goExec.Dir = path

	fmt.Printf("%s execute:\n%s", path, goExec.String())
	output, err := goExec.CombinedOutput()
	if err != nil {
		fmt.Println(string(output))
		return err
	}
	fmt.Println(string(output))

	return nil
}
