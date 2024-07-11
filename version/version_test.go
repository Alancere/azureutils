package version_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Alancere/azureutils/version"
)

func TestIsBetaVersion(t *testing.T) {
	versions := []string{
		"1.0.0",
		"0.1.0",
		"0.1.1-beta.1",
		"1.2.3-alpha.1",
		"1.2.3-alpha.1+build.1",
		"2.0.0-build.1",
		"2.1.0",
	}

	for _, v := range versions {
		isBeta, err := version.IsBetaVersion(v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s is beta: %t\n", v, isBeta)
	}
}

func TestTempFile(t *testing.T) {
	tempFile, err := os.CreateTemp("", "tspconfig-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	fmt.Println(tempFile.Name())
}
