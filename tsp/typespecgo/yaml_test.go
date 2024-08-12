package typespecgo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Alancere/azureutils/tsp/typespecgo"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
)

var tspConfigPath = "D:/Go/src/github.com/Azure/azure-rest-api-specs/specification/vmware/Microsoft.AVS.Management/tspconfig.yaml"

func TestGoYaml(t *testing.T) {
	tspConfig, err := typespecgo.NewTSPConfig(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	node, err := yaml.ValueToNode(&tspConfig)
	if err != nil {
		t.Fatal(err)
	}
	_ = node

	mapingNode := node.(*ast.MappingNode)
	iter := mapingNode.MapRange()
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(key.String())
		fmt.Println(value.String())
	}
}

func TestXxx(t *testing.T) {
	data, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	var node ast.Node // niubi
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		t.Fatal(err)
	}

	// newData, err := yaml.Marshal(&v)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(newData))

	iter := node.(*ast.MappingNode).MapRange()
	for iter.Next() {
		// key := iter.Key()
		// value := iter.Value()
		// fmt.Println(key.String())
		// fmt.Println(value.String())

		// switch value.Type() {
		// case ast.MappingType:

		// case ast.SequenceType:
		// case ast.StringType:
		// }
	}
	mapNode := ast.Filter(ast.MappingType, node)
	for _, n := range mapNode {
		p := ast.Parent(node, n)
		fmt.Println("parent:", p.GetToken().Prev.Value)
		fmt.Println(n.String())
		fmt.Println(n.GetToken().Prev.Value)
		fmt.Println("path:::", n.GetPath())
		fmt.Println("##########################################")
	}
}

func TestParse(t *testing.T) {
	f, err := parser.ParseFile(tspConfigPath, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	path, err := yaml.PathString("$.options.@azure-tools/typespec-autorest")
	if err != nil {
		t.Fatal(err)
	}
	node, err := path.FilterFile(f)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(node.String())

	m := node.(*ast.MappingNode)
	// m.Values = append(m.Values, &ast.MappingValueNode{})
	m.Merge(&ast.MappingNode{})
}

func rootNode(path string) *ast.Node {
	return nil
}

func TestParseByte(t *testing.T) {
	data := `options:
  "@azure-tools/typespec-go":
    service-dir: sdk
    package-dir: resourcemanager/mongocluster/armmongocluster
    module: github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}
    generate-fakes: true
    head-as-boolean: true
    inject-spans: true
    remove-unreferenced-types: true`
	data = `linter:
  extends:
    - "@azure-tools/typespec-azure-rulesets/resource-manager"
options:
  "@azure-tools/typespec-go":
    service-dir: sdk
    package-dir: resourcemanager/mongocluster/armmongocluster
    module: github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}
    generate-fakes: true
    head-as-boolean: true
    inject-spans: true
    remove-unreferenced-types: true`

	x, err := parser.ParseBytes([]byte(data), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	for _, doc := range x.Docs {
		fmt.Println(doc.String())
	}

	dataX, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	var node ast.Node // niubi
	err = yaml.Unmarshal(dataX, &node)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(node.Type(), x.Docs[0].Type())

	/*
		node, x.Docs[0]
		cannot merge MappingValue into Mapping
	*/
	err = ast.Merge(node, x.Docs[0])
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(node.String())
}

func TestParse2(t *testing.T) {
	data := `linter:
  extends:
    - "@azure-tools/typespec-azure-rulesets/resource-manager"
options:
  "@azure-tools/typespec-go":
    service-dir: sdk
    package-dir: resourcemanager/mongocluster/armmongocluster
    module: github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}
    generate-fakes: true
    head-as-boolean: true
    inject-spans: true
    remove-unreferenced-types: true`

	x, err := parser.ParseBytes([]byte(data), parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}

	f, err := parser.ParseFile(tspConfigPath, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	node, err := typespecgo.Node(f, "$.options")
	if err != nil {
		t.Fatal(err)
	}

	// *****************
	p, err := yaml.PathString("$.options")
	if err != nil {
		t.Fatal(err)
	}
	err = p.MergeFromFile(f, x)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(f.String())
	// **********************

	// dataX, err := os.ReadFile(tspConfigPath)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// var node ast.Node // niubi
	// err = yaml.Unmarshal(dataX, &node)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// mapValueNode := node.(*ast.MappingValueNode)
	// err = ast.Merge(mapValueNode.Value, x.Docs[0])
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(mapValueNode.String())

	dstMappingNode, ok := node.(*ast.MappingNode)
	if !ok {
		mappingValueNode, ok := node.(*ast.MappingValueNode)
		if ok {
			dstMappingNode = ast.Mapping(mappingValueNode.GetToken(), true, mappingValueNode)
		} else {
			t.Fatal("invalid node")
		}
	}
	mapVal := x.Docs[0].Body.(*ast.MappingValueNode)
	mapping := ast.Mapping(mapVal.GetToken(), true, mapVal)
	err = ast.Merge(dstMappingNode, mapping)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(dstMappingNode.String())
}

// 使用MapSlice是不错的选择保持了顺序
func TestOption(t *testing.T) {
	tspConfigPath := "D:/Go/src/github.com/Azure/azure-rest-api-specs/specification/vmware/Microsoft.AVS.Management/tspconfig.yaml"

	// tspConfig, err := typespecgo.NewTSPConfig(tspConfigPath)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	data, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}
	tspConfig := typespecgo.TSPConfig{}
	comment := yaml.CommentMap{}

	err = yaml.UnmarshalWithOptions(data, &(tspConfig.TypeSpecProjectSchema), yaml.CommentToMap(comment), yaml.UseOrderedMap())
	if err != nil {
		t.Fatal(err)
	}

	dataX, err := yaml.MarshalWithOptions(
		tspConfig.TypeSpecProjectSchema,

		// yaml.IndentSequence(true),
		// yaml.Flow(true),
		yaml.WithComment(comment),
		yaml.UseLiteralStyleIfMultiline(true),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(dataX))
	// err = os.WriteFile(tspConfigPath, data, 0o666)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	var orderMap yaml.MapSlice
	err = yaml.UnmarshalWithOptions(data, &orderMap, yaml.CommentToMap(comment), yaml.UseOrderedMap())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(orderMap)

	toMap := orderMap.ToMap()
	options, ok := toMap["options"]
	if !ok {
		t.Fatal("not found options")
	}

	optionsMap, ok := options.(yaml.MapSlice)
	if !ok {
		t.Fatal("options not map")
	}

	typespecgoOptionTemplate := `
@azure-tools/typespec-go:
  "generate-fakes":            true
  "head-as-boolean":           true
  "inject-spans":              true
  "remove-unreferenced-types": true
  "service-dir":        "sdk"
  "package-dir":        resourcemanager/%s/%s
  "module":             "github.com/Azure/azure-sdk-for-go/{service-dir}/{package-dir}"
  "examples-directory": "./examples"
`

	serviceName := "AVS"
	armServiceName := "AVS"
	typespecgoOption := fmt.Sprintf(typespecgoOptionTemplate, serviceName, armServiceName)
	mapItem := yaml.MapItem{}
	err = yaml.UnmarshalWithOptions([]byte(typespecgoOption), &mapItem, yaml.UseOrderedMap())
	if err != nil {
		t.Fatal(err)
	}

	optionsMap = append(optionsMap, mapItem)

	// orderMap["options"] = optionsMap
	fmt.Println(orderMap)
	for i, v := range orderMap {
		if v.Key.(string) == "options" {
			orderMap[i].Value = optionsMap
			break
		}
	}

	newData, err := yaml.MarshalWithOptions(orderMap, yaml.WithComment(comment))
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(tspConfigPath, newData, 0o666)
	if err != nil {
		t.Fatal(err)
	}
}
