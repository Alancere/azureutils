package typespecgo

import (
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"github.com/goccy/go-yaml/token"

	"github.com/iancoleman/orderedmap"
)

type NodeX struct {
	// tspconfit.yaml path
	Path string

	// tspconfig.yaml node
	File *ast.File
}

func NewNodeX(path string) (*NodeX, error) {
	file, err := parser.ParseFile(path, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &NodeX{
		Path: path,
		File: file,
	}, nil
}

// 通过路径获取节点
func Node(f *ast.File, path string) (ast.Node, error) {
	p, err := yaml.PathString(path)
	if err != nil {
		return nil, err
	}

	return p.FilterFile(f)
}

func MergeNode(f *ast.File, path string, value ast.Node) error {
	node, err := Node(f, path)
	if err != nil {
		return err
	}

	return ast.Merge(value, node)
}

// 将 @azure-tool/typespec-go 的配置添加到 ast.File 中
func AddGoOption(f *ast.File, mapValue *orderedmap.OrderedMap) error {
	typespecGoOption, err := Node(f, "$.options.@azure-tools/typespec-go")
	var tk *token.Token
	if err != nil {
		if yaml.IsNotFoundNodeError(err) {
			typespecGoOption = &ast.MappingNode{}
			tk = typespecGoOption.GetToken()
			// 是否需要获取到 $.options 的 token
			options, err := Node(f, "$.options")
			if err != nil {
				return err
			}
			options.GetToken()

			optionsX := options.(*ast.MappingNode)
			optionsX.GetToken()
		}
		return err
	} else {
		tk = typespecGoOption.GetToken()
	}

	for _, key := range mapValue.Keys() {
		value, _ := mapValue.Get(key)
		keyNode := ast.MappingKey(token.New(key, key, &token.Position{}))
		valueNode := ast.MappingValue(token.New(value.(string), value.(string), &token.Position{}), keyNode, &ast.IntegerNode{
			Value: value.(string),
		})
		ast.Mapping(tk, true, valueNode)
	}

	return nil
}
