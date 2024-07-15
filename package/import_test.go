package packages_test

import (
	"testing"

	packages "github.com/Alancere/azureutils/package"
	"github.com/stretchr/testify/assert"
)

func TestUpdatePackageModule(t *testing.T) {
	testDataPath := "./testdata"
	baseModule := "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"

	newVersion := "v5.0.0"
	err := packages.UpdatePackageModule(newVersion, testDataPath, baseModule, "armnetwork_test.go")
	assert.NoError(t, err)

	newVersion = "v1.0.0"
	err = packages.UpdatePackageModule(newVersion, testDataPath, baseModule, "armnetwork_test.go")
	assert.NoError(t, err)
}
