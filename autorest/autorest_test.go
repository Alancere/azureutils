package autoret_test

import (
	"testing"

	autoret "github.com/Alancere/azureutils/autorest"
)

func TestReadAutorsetMD(t *testing.T) {
	mdPath := "autorest.md"
	_, err := autoret.ReadAutorestMD(mdPath)
	if err != nil {
		t.Errorf("Error reading autorest.md: %v", err)
	}
}
