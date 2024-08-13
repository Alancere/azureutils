package demo_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/Masterminds/semver/v3"
	jsoniter "github.com/json-iterator/go"
)

func TestJsoniter(t *testing.T) {
	src := "D:\\Go\\src\\github.com\\Azure\\autorest.go\\packages\\typespec-go\\package.json"
	dst := "D:\\Go\\src\\github.com\\Azure\\azure-sdk-for-go\\eng\\package.json"
	dst = "D:\\Go\\src\\github.com\\Azure\\azure-sdk-for-go\\eng\\emitter-package.json"

	err := syncPackageJson(src, dst)
	if err != nil {
		t.Fatal(err)
	}
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

			fmt.Printf("replaced version:%s => %s\n", v1.String(), v2.String())
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

func TestDependencies(t *testing.T) {
	data := `{
    "main": "dist/src/index.js",
    "dependencies": {
      "@azure-tools/typespec-go": "0.3.0"
    },
    "devDependencies": {
      "@azure-tools/typespec-autorest": "0.44.1",
      "@azure-tools/typespec-azure-core": "0.44.0",
      "@azure-tools/typespec-azure-resource-manager": "0.44.0",
      "@azure-tools/typespec-azure-rulesets": "0.44.0",
      "@azure-tools/typespec-client-generator-core": "0.44.3",
      "@typespec/compiler": "0.58.1",
      "@typespec/http": "0.58.0",
      "@typespec/openapi": "0.58.0",
      "@typespec/rest": "0.58.0",
      "@typespec/versioning": "0.58.0"
    }
}`

	deps := subKeys([]byte(data), "devDependenciesX")
	fmt.Println(len(deps))
}
