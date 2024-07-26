package typespecgo

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type nodeX struct {
	parent *yaml.Node
	// key与value的索引
	index int

	// key   *yaml.Node
	// value *yaml.Node
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
	if n.Value().Kind == yaml.ScalarNode {
		return fmt.Sprintf("%s: %s", n.Key().Value, n.Value().Value)
	}

	buf := bytes.Buffer{}
	encode := yaml.NewEncoder(&buf)
	encode.SetIndent(2)
	err := encode.Encode(n.Value())
	if err != nil {
		return ""
	}

	// add two space
	lines := strings.Split(buf.String(), "\n")
	for i := range lines {
		lines[i] = fmt.Sprintf("  %s", lines[i])
	}

	return fmt.Sprintf("%s:\n%s", n.Key().Value, strings.Join(lines, "\n"))
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

					// key:   node.Content[i],
					// value: node.Content[i+1],
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

// 插入到的指定位置
func insertUpdateNode(node *yaml.Node, target string, values ...*yaml.Node) error {
	if node.IsZero() {
		return fmt.Errorf("node is empty")
	}

	x := getNode(node, target)
	if x == nil {
		return fmt.Errorf("%s not found", target)
	}

	if x.Value().Kind == yaml.MappingNode {
		x.Value().Content = append(x.Value().Content, values...)
		return nil
	} else if x.Value().Kind == yaml.SequenceNode {
		x.Value().Content = append(x.Value().Content, values...)
	} else if x.Value().Kind == yaml.ScalarNode {
		if len(values) != 1 {
			return fmt.Errorf("ScalarNode must only value")
		}
		x.Value().Value = values[0].Value
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

type TSPConfigX struct {
	node *yaml.Node

	TSPConfig
}

func (t *TSPConfigX) Unmarshal(data []byte) error {
	var node yaml.Node

	err := yaml.Unmarshal(data, &node)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, &t.TSPConfig.TypeSpecProjectSchema)
}

func (t *TSPConfigX) Marshal() ([]byte, error) {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	err := enc.Encode(t.node)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
