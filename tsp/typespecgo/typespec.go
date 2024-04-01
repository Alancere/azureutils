package typespecgo

import (
	"os"
	"regexp"
	"slices"

	"gopkg.in/yaml.v3"
)

/*
tspconfig schema: https://typespec.io/docs/handbook/configuration#schema
*/
/*
model TypeSpecProjectSchema {
  extends?: string;
  parameters?: Record<{default: string}>
  "environment-variables"?: Record<{default: string}>
  "warn-as-error"?: boolean;
  "output-dir"?: boolean;
  "trace"?: string | string[];
  imports?: string;
  emit?: string[];
  options?: Record<unknown>;
  linter?: LinterConfig;
}

model LinterConfig {
  extends?: RuleRef[];
  enable?: Record<RuleRef, boolean>;
  disable?: Record<RuleRef, string>;
}
*/

type TSPConfig struct {
	Path string

	TypeSpecProjectSchema
}

// https://typespec.io/docs/handbook/configuration#schema
type TypeSpecProjectSchema struct {
	Extends              string         `yaml:"extends,omitempty"`
	Parameters           map[string]any `yaml:"parameters,omitempty"`
	EnvironmentVariables map[string]any `yaml:"environment-variables,omitempty"`
	WarnAsError          bool           `yaml:"warn-as-error,omitempty"`
	OutPutDir            string         `yaml:"output-dir,omitempty"` // 不应该是bool
	Trace                []string       `yaml:"trace,omitempty"`
	Imports              string         `yaml:"imports,omitempty"`
	Emit                 []string       `yaml:"emit,omitempty"`
	Options              map[string]any `yaml:"options,omitempty"`
	Linter               LinterConfig   `yaml:"linter,omitempty"`
}

// <library name>:<rule/ruleset name>
type LinterConfig struct {
	Extends []RuleRef          `yaml:"extends,omitempty"`
	Enable  map[RuleRef]bool   `yaml:"enable,omitempty"`
	Disable map[RuleRef]string `yaml:"disable,omitempty"`
}

type TypeSpecAzureTools string

const (
	TypeSpec_GO       TypeSpecAzureTools = "@azure-tools/typespec-go"
	TypeSpec_AUTOREST TypeSpecAzureTools = "@azure-tools/typespec-autorest"
	TypeSpec_CSHARP   TypeSpecAzureTools = "@azure-tools/typespec-csharp"
	TypeSpec_PYTHON   TypeSpecAzureTools = "@azure-tools/typespec-python"
	TypeSpec_TS       TypeSpecAzureTools = "@azure-tools/typespec-ts"
)

type RuleRef string

func (r RuleRef) Validate() bool {
	return regexp.MustCompile(`.*/.*`).MatchString(string(r))
}

func NewTSPConfig(tspconfigYaml string) (*TSPConfig, error) {
	tspConfig := TSPConfig{}

	data, err := os.ReadFile(tspconfigYaml)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(data, &(tspConfig.TypeSpecProjectSchema))
	if err != nil {
		return nil, err
	}

	tspConfig.Path = tspconfigYaml
	return &tspConfig, err
}

func (tc *TSPConfig) EditEmit(emits []string) {
	if tc.Emit == nil {
		tc.Emit = emits
		return
	}

	tc.Emit = append(tc.Emit, emits...)
	tc.Emit = slices.Compact(tc.Emit)
}

func (tc *TSPConfig) OnlyEmit(emit string) {
	tc.Emit = []string{emit}
}

func (tc *TSPConfig) EditOptions(emit string, option any, append bool) {
	if tc.Options == nil {
		tc.Options = make(map[string]any)
	}

	if _, ok := tc.Options[emit]; ok {
		if append {
			op1 := tc.Options
			for k, v := range option.(map[string]any) {
				op1[k] = v
			}
			tc.Options = op1
		} else {
			tc.Options[emit] = option
		}
	} else {
		tc.Options[emit] = option
	}
}

func (tc *TSPConfig) Write() error {
	data, err := yaml.Marshal(tc.TypeSpecProjectSchema)
	if err != nil {
		return err
	}

	return os.WriteFile(tc.Path, data, 0o666)

	// f,err := os.OpenFile(tc.Path, os.O_CREATE | os.O_RDWR, 0777)
	// if err != nil {
	// 	return err
	// }

	// w := bufio.NewWriter(f)

	// x := yaml.NewEncoder(w)
	// x.SetIndent(2)
	// err = x.Encode(tc.TypeSpecProjectSchema)
	// if err != nil {
	// 	return err
	// }

	// return w.Flush()
}

// func (tsps TypeSpecProjectSchema) MarshalYAML() (interface{}, error) {
// 	return &struct{
// 		Extends *yaml.Node `yaml:"extends,omitempty"`
// 		Parameters map[string]any `yaml:"parameters,omitempty"`
// 		EnvironmentVariables map[string]any `yaml:"environment-variables,omitempty"`
// 		WarnAsError bool `yaml:"warn-as-error,omitempty"`
// 		OutPutDir string `yaml:"output-dir,omitempty"` // 不应该是bool
// 		Trace []string `yaml:"trace,omitempty"`
// 		Imports string `yaml:"imports,omitempty"`
// 		Emit []string `yaml:"emit,omitempty"`
// 		Options map[string]any `yaml:"options,omitempty"`
// 		Linter LinterConfig `yaml:"linter,omitempty"`
// 	}{
// 		Extends: &yaml.Node{Kind: yaml.ScalarNode, Value: tsps.Extends, Style: },
// 	}, nil
// }
