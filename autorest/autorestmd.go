package autoret

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/russross/blackfriday/v2"
	"gopkg.in/yaml.v3"
)

var EXT = []string{"yaml", "yml", "go", "json"}

type AutorestMD struct {
	AzureARM      string   `yaml:"azure-arm,omitempty"`
	Require       []string `yaml:"require,omitempty"`
	LicenseHeader string   `yaml:"license-header,omitempty"`
	ModuleVersion string   `yaml:"module-version,omitempty"`
	Tag           string   `yaml:"tag,omitempty"`

	Directive any `yaml:"directive,omitempty"`
}

func ReadAutorestMD(md string) (*AutorestMD, error) {
	data, err := os.ReadFile(md)
	if err != nil {
		return nil, err
	}

	autorestMD := AutorestMD{}
	result := make([]byte, 0)

	x := blackfriday.New()
	n := x.Parse(data)
	n.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			if node.Type.String() == "Code" {
				fmt.Printf("%s\n", string(node.Literal))

				split := strings.Split(string(node.Literal), "\n")
				if slices.Contains(EXT, strings.TrimSpace(split[0])) {
					split = split[1:]
				}

				result = []byte(strings.Join(split, "\n"))
				return blackfriday.SkipChildren
			}
		}
		return blackfriday.GoToNext
	})

	if err = yaml.Unmarshal(result, &autorestMD); err != nil {
		return nil, err
	}

	return &autorestMD, nil
}
