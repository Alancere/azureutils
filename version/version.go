package version

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

func IsBetaVersion(v string) (bool, error) {
	newVersion, err := semver.NewVersion(v)
	if err != nil {
		return false, err
	}

	fmt.Println(newVersion.Major(), newVersion.Minor(), newVersion.Patch(), newVersion.Prerelease())

	if newVersion.Major() == 0 || newVersion.Prerelease() != "" {
		return true, nil
	}

	return false, nil
}
