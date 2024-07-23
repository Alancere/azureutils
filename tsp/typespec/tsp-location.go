package typespec

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

// tsp-location.yaml
type TspLocation struct {
	Directory             string   `yaml:"directory" validate:"required"`
	Commit                string   `yaml:"commit" validate:"required"`
	Repo                  string   `yaml:"repo" validate:"required"`
	AdditionalDirectories []string `yaml:"additionalDirectories" validate:"omitempty"`
}

func ParseTspLocation(tspLocationPath string) (*TspLocation, error) {
	var tl TspLocation

	data, err := os.ReadFile(tspLocationPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &tl)
	if err != nil {
		return nil, err
	}

	if err = validator.New(validator.WithRequiredStructEnabled()).Struct(tl); err != nil {
		return nil, err
	}

	return &tl, nil
}

// 快速获取tsp-location.yaml指向的github link
func (tl *TspLocation) Link() *Link {
	link := new(Link)
	linkPrefx := fmt.Sprintf("https://github.com/%s/tree/%s", strings.Trim(tl.Repo, "/"), tl.Commit)
	link.link = fmt.Sprintf("%s/%s", linkPrefx, tl.Directory)
	if len(tl.AdditionalDirectories) > 0 {
		for _, ad := range tl.AdditionalDirectories {
			link.additionalLinks = append(link.additionalLinks, fmt.Sprintf("%s/%s", linkPrefx, strings.Trim(ad, "/")))
		}
	}

	return link
}

func (tl *TspLocation) Validation() error {
	// 判断这些是否是有效link
	link := tl.Link()
	for _, link := range link.Links() {
		resp, err := http.Get(link)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s is not valid", link)
		}
	}

	return nil
}
