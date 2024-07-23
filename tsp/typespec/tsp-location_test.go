package typespec

import (
	"fmt"
	"testing"
)

func TestTspLocation_RepoLink(t *testing.T) {
	/*
		tsp-location:
			directory: specification/contosowidgetmanager/Contoso.WidgetManager
			commit: 431eb865a581da2cd7b9e953ae52cb146f31c2a6
			repo: Azure/azure-rest-api-specs
			additionalDirectories:
			- specification/contosowidgetmanager/Contoso.WidgetManager.Shared/
	*/
	tl := &TspLocation{
		Directory:             "specification/contosowidgetmanager/Contoso.WidgetManager",
		Commit:                "431eb865a581da2cd7b9e953ae52cb146f31c2a6",
		Repo:                  "Azure/azure-rest-api-specs",
		AdditionalDirectories: []string{"specification/contosowidgetmanager/Contoso.WidgetManager.Shared/"},
	}
	fmt.Println(tl.Link())

	err := tl.Validation()
	if err != nil {
		t.Error(err)
	}
}
