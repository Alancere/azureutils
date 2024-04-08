package fetchmodule

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

/*
  使新release的package快速出现在pkg.go.dev上

  https://pkg.go.dev/about
*/

const (
	PKGGODEV = "https://pkg.go.dev"
	GOPROXY  = "https://proxy.golang.org"
)

// https://go.dev/fetch/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v2@v2.0.0
func Fetch(module string) error {
	url := fmt.Sprint(PKGGODEV + "/fetch/" + module)
	fmt.Println("fetch:", url)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		log.Printf("please click link:%s/%s", PKGGODEV, module)
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}

// "github.com/azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v2/@v/v2.info"
func Info(module string) error {
	module = strings.ToLower(module)
	before, after, _ := strings.Cut(module, "@")
	version, err := semver.NewVersion(after)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s/@v/v%d.info", GOPROXY, before, version.Major())
	fmt.Println("info:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Println(string(body))
	return nil
}

// GOPROXY=https://proxy.golang.org GO111MODULE=on go get example.com/my/module@v1.0.0
func GoGet(module string) error {
	os.Setenv("GOPROXY", "https://proxy.golang.org")
	os.Setenv("GO111MODULE", "auto")
	cmd := exec.Command("go", "get", module)
	// cmd.Dir = ""
	fmt.Printf("GO111MODULE=%s go get %s\n", os.Getenv("GO111MODULE"), module)

	output, err := cmd.CombinedOutput()
	log.Printf("Result of `%s` execution: \n%s", cmd.String(), string(output))
	if err != nil {
		return fmt.Errorf("failed to execute `%s` '%s': %+v", cmd.String(), string(output), err)
	}
	return nil
}

// new packages to add to pkg.go.dev
func NewPackages(module string) error {
	if err := Fetch(module); err != nil {
		return err
	}

	if err := Info(module); err != nil {
		return err
	}

	// if err := GoGet(module); err != nil {
	// 	return err
	// }

	return nil
}

// module-name@version
func Validate(module string) (bool, error) {
	module = strings.ToLower(module)

	b := regexp.MustCompile(`^[\w\d\.\-\/]+@v\d+\.\d+\.\d+[\-\w\d\.]?`).MatchString(module)
	if !b {
		return false, fmt.Errorf("invalid module name(module-name@version): %s", module)
	}
	return b, nil
}

// github.com/alancere/azureutils@v0.1.0
