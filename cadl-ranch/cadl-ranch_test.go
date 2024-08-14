package cadlranch_test

import (
	"testing"

	cadlranch "github.com/Alancere/azureutils/cadl-ranch"
)

func TestGenerateCadlRanchTest(t *testing.T) {
	typespecGoPath := "D:\\Go\\src\\github.com\\Azure\\autorest.go\\packages\\typespec-go"
	cadlRanchTest := "azure/resource-manager/models/resources"
	cadlRanchTest = "azure/resource-manager/models/common-types/managed-identity"
	if err := cadlranch.GenerateCadlRanchTest(typespecGoPath, cadlRanchTest); err != nil {
		t.Fatal(err)
	}
}

func TestValidationCadlRanchTestPath(t *testing.T) {
	cadlRanchTest := "azure/resource-manager/models/resources"
	err := cadlranch.ValidationCadlRanchTestPath(cadlRanchTest)
	if err != nil {
		t.Fatal(err)
	}

	cadlRanchTest = "azure/resource-manager/models/common-types/managed-identity"
	err = cadlranch.ValidationCadlRanchTestPath(cadlRanchTest)
	if err != nil {
		t.Fatal(err)
	}

	cadlRanchTest = "azure/resource-manager/models/common-types/managed-identityX"
	err = cadlranch.ValidationCadlRanchTestPath(cadlRanchTest)
	if err != nil {
		t.Fatal(err)
	}
}
