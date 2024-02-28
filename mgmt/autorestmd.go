package mgmt

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
)

const (
	ModuleVersionPrefix = "module-version: "
	ChangelogMD         = "CHANGELOG.md"
)

var (
	ModuleVersionRegex           = regexp.MustCompile(`moduleVersion\s*=\s*\".*v.+"`)
	ChangelogPosWithPreviewRegex = regexp.MustCompile(`##\s*(?P<version>.+)\s*\((\d{4}-\d{2}-\d{2}|Unreleased)\)`)
)

/*
从 autorest.md 中获取 module-version
*/
func AutorestMdVersion(md string) (*semver.Version, error) {
	mdFile, err := os.ReadFile(md)
	if err != nil {
		return nil, err
	}

	for _, v := range strings.Split(string(mdFile), "\n") {
		if strings.Contains(v, ModuleVersionPrefix) {
			mv, _ := strings.CutPrefix(v, ModuleVersionPrefix)
			version, err := semver.NewVersion(strings.TrimSpace(mv))
			if err != nil {
				return nil, err
			}
			return version, nil
		}
	}

	return nil, fmt.Errorf("%s 中不存在 %s", md, ModuleVersionPrefix)
}

// replace version: use `module-version: ` prefix to locate version in autorest.md file, use version = "v*.*.*" regrex to locate version in constants.go file
func ReplaceVersion(packageRootPath string, newVersion string) error {
	path := filepath.Join(packageRootPath, "autorest.md")
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(b), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, ModuleVersionPrefix) {
			lines[i] = line[:len(ModuleVersionPrefix)] + newVersion
			break
		}
	}

	if err = os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return err
	}

	path = filepath.Join(packageRootPath, "constants.go")
	if b, err = os.ReadFile(path); err != nil {
		return err
	}
	contents := ModuleVersionRegex.ReplaceAllString(string(b), "moduleVersion = \"v"+newVersion+"\"")

	return os.WriteFile(path, []byte(contents), 0644)
}

// add new changelog md to changelog file
func AddChangelogToFile(additionalChangelog string, version *semver.Version, packageRootPath, releaseDate string) (string, error) {
	path := filepath.Join(packageRootPath, ChangelogMD)
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	oldChangelog := string(b)
	newChangelog := "# Release History\n\n"
	matchResults := ChangelogPosWithPreviewRegex.FindAllStringSubmatchIndex(oldChangelog, -1)
	// additionalChangelog := changelog.ToCompactMarkdown()
	if releaseDate == "" {
		releaseDate = time.Now().Format("2006-01-02")
	}

	for _, matchResult := range matchResults {
		newChangelog = newChangelog + "## " + version.String() + " (" + releaseDate + ")\r\n" + additionalChangelog + "\r\n\r\n" + oldChangelog[matchResult[0]:]
		break
	}

	err = os.WriteFile(path, []byte(newChangelog), 0644)
	if err != nil {
		return "", err
	}
	return additionalChangelog, nil
}
