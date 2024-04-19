package autorest

import (
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var EXT = []string{"yaml", "yml", "go", "json"}

type AutoRestMarkdown struct {
	AzureARM      string   `yaml:"azure-arm,omitempty"`
	Require       []string `yaml:"require,omitempty"`
	LicenseHeader string   `yaml:"license-header,omitempty"`
	ModuleVersion string   `yaml:"module-version,omitempty"`
	Tag           string   `yaml:"tag,omitempty"`

	Directive any `yaml:"directive,omitempty"`
}

func ReadAutoRestMarkdown(md string) (*AutoRestMarkdown, error) {
	data, err := os.ReadFile(md)
	if err != nil {
		return nil, err
	}

	armd := AutoRestMarkdown{}
	r := regexp.MustCompile("``` yaml\n([\\s\\S]*)\n```")
	matches := r.FindSubmatch(data)
	if err = yaml.Unmarshal(matches[1], &armd); err != nil {
		return nil, err
	}

	return &armd, nil
}

func (a AutoRestMarkdown) GetSpec() (string, string) {
	_, spec, _ := strings.Cut(strings.ReplaceAll(a.Require[0], "\\", "/"), "specification/")

	return spec[:strings.Index(spec, "/")], spec[:strings.LastIndex(spec, "/")]
}
