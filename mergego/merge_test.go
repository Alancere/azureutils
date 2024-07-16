package mergego_test

import (
	"testing"

	"github.com/Alancere/azureutils/mergego"
)

func TestMerge(t *testing.T) {
	dir := "./mypackage" // replace with your package directory
	dir = "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs/specification/apicenter/ApiCenter.Management/go/sdk/resourcemanager/apicenter/armapicenter"
	outfile := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs/specification/apicenter/ApiCenter.Management/go/sdk/resourcemanager/apicenter/armapicenter/merged.go" // output file

	err := mergego.Merge(dir, outfile, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMergeFake(t *testing.T) {
	dir := "./mypackage" // replace with your package directory
	dir = "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs/specification/apicenter/ApiCenter.Management/go/sdk/resourcemanager/apicenter/armapicenter/fake"
	outfile := "D:/Go/src/github.com/Azure/dev/azure-rest-api-specs/specification/apicenter/ApiCenter.Management/go/sdk/resourcemanager/apicenter/armapicenter/fake.go" // output file

	err := mergego.Merge(dir, outfile, false)
	if err != nil {
		t.Fatal(err)
	}
}
