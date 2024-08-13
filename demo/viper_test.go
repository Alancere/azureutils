package demo_test

import (
	"fmt"
	"testing"

	"github.com/spf13/viper"
)

func TestViperParseYaml(t *testing.T) {
	tspConfigPath := "D:\\Go\\src\\github.com\\Azure\\azure-rest-api-specs\\specification\\mongocluster\\DocumentDB.MongoCluster.Management"
	viper.SetConfigName("tspconfig")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(tspConfigPath)

	if err := viper.ReadInConfig(); err != nil {
		t.Errorf("Failed to read configuration file: %s", err)
	}

	options := viper.GetStringMap("options")
	fmt.Println(options)

	typespecGoOption := viper.GetStringMap("options.@azure-tools/typespec-go")
	fmt.Println(typespecGoOption)

	viper.Set("options.@azure-tools/typespec-go", "test")
	if err := viper.WriteConfig(); err != nil {
		t.Errorf("Failed to write configuration file: %s", err)
	}
}
