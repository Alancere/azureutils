package sdk_test

import (
	"fmt"
	"testing"

	"github.com/Alancere/azureutils/sdk"
)

func TestGetAllMgmtSDK(t *testing.T) {
	localPath := "D:/Go/src/github.com/Azure/azure-sdk-for-go/sdk/resourcemanager"
	mgmts, err := sdk.GetAllMgmtSDK(localPath)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("mgmt count:", len(mgmts))
	for _, m := range mgmts {
		fmt.Println(m.LocalPath)
	}

	fmt.Println("depcated count:", len(sdk.DeprecatedMgmtSDK))
	for _, m := range sdk.DeprecatedMgmtSDK {
		fmt.Println(m.LocalPath)
	}

	fmt.Println("all count:(mgmtsdk + deprecated)", len(mgmts)+len(sdk.DeprecatedMgmtSDK))
}
