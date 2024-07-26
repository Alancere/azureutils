package typespecgo_test

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/Alancere/azureutils/tsp/typespecgo"
	"gopkg.in/yaml.v3"
)

func TestUnMarshalTspConfig(t *testing.T) {
	tspConfigPath := "D:/Go/src/github.com/Azure/azure-rest-api-specs/specification/vmware/Microsoft.AVS.Management/tspconfig.yaml"

	tspConfig := typespecgo.TSPConfig{}

	data, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(data, &tspConfig.TypeSpecProjectSchema)
	if err != nil {
		t.Fatal(err)
	}

	// data, err = yaml.Marshal(tspConfig.TypeSpecProjectSchema)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// save to file
	// f, err := os.Open(tspConfigPath)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// defer f.Close()

	encode := yaml.NewEncoder(os.Stdout)
	encode.SetIndent(2)
	err = encode.Encode(tspConfig.TypeSpecProjectSchema)
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(string(data))
}

func TestNode(t *testing.T) {
	tspConfigPath := "D:/Go/src/github.com/Azure/azure-rest-api-specs/specification/vmware/Microsoft.AVS.Management/tspconfig.yaml"
	data, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		t.Fatal(err)
	}

	// outputData, err := yaml.Marshal(&node)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println(string(outputData))

	//
	// t.Log(node.IsZero())
	// t.Log(node.LongTag())
	// t.Log(node.ShortTag())
	// node.SetString("test")

	editGoEmitOption(&node)

	// 遍历寻找@azure-tools/typespec-go并添加flavor: azure
	found := addFlavorAzure(node.Content[0], "@azure-tools/typespec-go", "flavor", "azure")
	if !found {
		t.Fatal("@azure-tools/typespec-go not found")
	}

	f, err := os.OpenFile(tspConfigPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	encode := yaml.NewEncoder(f)
	encode.SetIndent(2)
	err = encode.Encode(&node)
	if err != nil {
		t.Fatal(err)
	}
}

func addFlavorAzure(node *yaml.Node, path, key, value string) bool {
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == "options" { // 需要先找到options
				conents := node.Content[i+1].Content
				for j := 0; j < len(node.Content[i+1].Content); j += 2 {
					if conents[j].Value == path {
						fmt.Println(conents[j].Value)
						// 在找到的节点下添加flavor: azure
						conents[j+1].Content = append(conents[j+1].Content, &yaml.Node{Kind: yaml.ScalarNode, Value: key}, &yaml.Node{Kind: yaml.ScalarNode, Value: value})
						return true
					}
				}
			}
		}
	}
	return false
}

func editGoEmitOption(node *yaml.Node) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}
	if node.Kind == yaml.DocumentNode {
		node = node.Content[0]
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == "options" { // 需要先找到options
				options := node.Content[i+1].Content

				existTypeSpecGo := false
				for j := 0; j < len(options); j += 2 {
					if options[j].Value == "@azure-tools/typespec-go" {
						existTypeSpecGo = true
						options[j+1].Content = append(options[j+1].Content,
							&yaml.Node{Kind: yaml.ScalarNode, Value: "flavor"},
							&yaml.Node{Kind: yaml.ScalarNode, Value: "azure"},
						)
						return nil
					}
				}

				if !existTypeSpecGo {
				}
			}
		}
	}

	return fmt.Errorf("node is not a mapping node")
}

func getNode(node *yaml.Node, target string) *nodeX {
	if node.IsZero() {
		return nil
	}
	if node.Kind == yaml.DocumentNode {
		node = node.Content[0]
	}

	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			if node.Content[i].Value == target {
				return &nodeX{
					parent: node,
					index:  i,

					key:   node.Content[i],
					value: node.Content[i+1],
				}
			}
		}

		for j := 1; j < len(node.Content); j += 2 {
			if node.Content[j].Kind == yaml.MappingNode {
				n := getNode(node.Content[j], target)
				if n != nil {
					return n
				}
			}
		}
	}

	return nil
}

type nodeX struct {
	parent *yaml.Node
	// key与value的索引
	index int

	key   *yaml.Node
	value *yaml.Node
}

func (n *nodeX) Key() *yaml.Node {
	if n.parent.IsZero() {
		return nil
	}

	return n.parent.Content[n.index]
}

func (n *nodeX) Value() *yaml.Node {
	if n.parent.IsZero() {
		return nil
	}

	return n.parent.Content[n.index+1]
}

func (n *nodeX) String() string {
	if n.value.Kind == yaml.ScalarNode {
		return fmt.Sprintf("%s: %s", n.key.Value, n.value.Value)
	}

	buf := bytes.Buffer{}
	encode := yaml.NewEncoder(&buf)
	encode.SetIndent(2)
	err := encode.Encode(&n.value)
	if err != nil {
		return ""
	}

	// add two space
	lines := strings.Split(buf.String(), "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("  %s", lines[i])
	}

	return fmt.Sprintf("%s:\n%s", n.key.Value, strings.Join(lines, "\n"))
}

// 插入到的指定位置
func insertMapNode(node *yaml.Node, target string, values ...*yaml.Node) error {
	if node.IsZero() {
		return fmt.Errorf("node is empty")
	}

	x := getNode(node, target)
	if x == nil {
		return fmt.Errorf("%s not found", target)
	}

	if x.value.Kind == yaml.MappingNode {
		x.value.Content = append(x.value.Content, values...)
		return nil
	} else if x.value.Kind == yaml.SequenceNode {
		x.value.Content = append(x.value.Content, values...)
	} else if x.value.Kind == yaml.ScalarNode {
		if len(values) != 1 {
			return fmt.Errorf("ScalarNode must only value")
		}
		x.value.Value = values[0].Value
		return nil
	}

	return fmt.Errorf("target is not a mapping node")
}

func removeNode(node *yaml.Node, target string) {
	if node.IsZero() {
		return
	}

	x := getNode(node, target)
	if x == nil {
		return
	}

	x.parent.Content = append(x.parent.Content[:x.index], x.parent.Content[x.index+2:]...)
}

func TestNodeX(t *testing.T) {
	tspConfigPath := "D:/Go/src/github.com/Azure/azure-rest-api-specs/specification/vmware/Microsoft.AVS.Management/tspconfig.yaml"
	data, err := os.ReadFile(tspConfigPath)
	if err != nil {
		t.Fatal(err)
	}

	var node yaml.Node
	err = yaml.Unmarshal(data, &node)
	if err != nil {
		t.Fatal(err)
	}

	// removeNode(&node, "options")
	x := getNode(&node, "options")
	if x == nil {
		t.Fatal(err)
	}
	// removeNode(&node, "options")
	fmt.Println(x.String())

	x = getNode(&node, "emit")
	fmt.Println(x.String())

	// insertMapNode(&node, "emit", nil, nil)
	insertMapNode(&node, "use-read-only-status-schema", &yaml.Node{Kind: yaml.ScalarNode, Value: "XXX"})
	x = getNode(&node, "use-read-only-status-schema")
	fmt.Println(x.String())

	x = getNode(&node, "@azure-tools/typespec-go")
	if x == nil {
		t.Fatal(err)
	}
	insertMapNode(&node, "@azure-tools/typespec-go", &yaml.Node{Kind: yaml.ScalarNode, Value: "flavor"}, &yaml.Node{Kind: yaml.ScalarNode, Value: "azure"})
	fmt.Println(x.String())

	x = getNode(&node, "@azure-tools/typespec-ts")
	if x == nil {
		t.Log("not found")
	}
	fmt.Println(x.String())
}

