package fetchmodule_test

import (
	"testing"

	fetchmodule "github.com/Alancere/azureutils/fetch_module"
)

func TestPkgGoDev(t *testing.T) {
	module := "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/redisenterprise/armredisenterprise/v2@v2.0.0"

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
