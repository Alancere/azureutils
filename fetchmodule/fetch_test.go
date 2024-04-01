package fetchmodule_test

import (
	"fmt"
	"testing"

	fetchmodule "github.com/Alancere/azureutils/fetchmodule"
)

func TestPkgGoDev(t *testing.T) {
	module := "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v2@v2.0.0"
	module = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/storageactions/armstorageactions@v0.1.0"
	// 1.
	// if err := fetchmodule.Fetch(module); err != nil {
	// 	t.Fatal(err)
	// }

	// if err := fetchmodule.Info(module); err != nil {
	// 	t.Fatal(err)
	// }

	// if err := fetchmodule.GoGet(module); err != nil {
	// 	t.Fatal(err)
	// }

	if err := fetchmodule.NewPackages(module); err != nil {
		t.Fatal(err)
	}
}

func TestValidate(t *testing.T) {
	module := "github.com/alancere/azureutils"
	// module = "github.com/azure/azure-sdk-for-go"
	b, err := fetchmodule.Validate(module)
	fmt.Println(b, err)

	module = "github.com/azure/azure-sdk-for-go@v0.1.0"
	b, err = fetchmodule.Validate(module)
	fmt.Println(b, err)

	module = "github.com/azure/azure-sdk-for-go@v0.1.0-beta.1"
	b, err = fetchmodule.Validate(module)
	fmt.Println(b, err)
}
