package packages_test

import (
	"fmt"
	"testing"

	packages "github.com/Alancere/azureutils/package"
	"github.com/stretchr/testify/assert"
)

func TestGoModValidate(t *testing.T) {
	sdkPath := "D:\\Go\\src\\github.com\\Azure\\azure-sdk-for-go\\sdk\\resourcemanager\\network\\armnetwork"
	err := packages.GoModValidate(sdkPath)
	fmt.Println(err, packages.ErrInvalidModule)
	assert.ErrorAs(t, err, &packages.ErrInvalidModule)
	assert.ErrorIs(t, err, packages.ErrInvalidModule)
}
