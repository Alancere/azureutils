package spec_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/Alancere/azureutils/spec"
)

func TestCheckAzureSpecs(t *testing.T) {
	fmt.Printf("### check %s ###", spec.PublicRepository)
	err := spec.NewPullRequest().Check(spec.PublicRepository)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("### check %s ###", spec.PrivateRepository)
	err = spec.NewPullRequest().Check(spec.PrivateRepository)
	if err != nil {
		log.Fatal(err)
	}
}
