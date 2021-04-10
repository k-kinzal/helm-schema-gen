package helper

import (
	"gopkg.in/yaml.v3"
	"strings"
)

type visitor interface {
	Visit(node *yaml.Node) (w visitor)
}

type inspector func(*yaml.Node) bool

func (f inspector) Visit(node *yaml.Node) visitor {
	if f(node) {
		return f
	}
	return nil
}


func walk(v visitor, node *yaml.Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch node.Kind {
	case yaml.DocumentNode, yaml.SequenceNode, yaml.MappingNode:
		for _, child := range node.Content {
			walk(v, child)
		}
	}
}

func InspectYAML(node *yaml.Node, f func(*yaml.Node) bool) {
	walk(inspector(f), node)
}

func UnmarshalWithUncommentYAML(in []byte, out interface{}) (err error) {
	root := yaml.Node{}
	if err := yaml.Unmarshal(in, &root); err != nil {
		return err
	}

	lines := splitLine(string(in))
	InspectYAML(&root, func(node *yaml.Node) bool {
		root := yaml.Node{}
		if err := yaml.Unmarshal([]byte(uncommentYAML(node.HeadComment)), &root); err == nil && len(root.Content) > 0 {
			if  root.Content[0].Kind == yaml.SequenceNode || root.Content[0].Kind == yaml.MappingNode {
				from := node.Line - 2 - lineCount(node.HeadComment)
				end := node.Line - 2
				for i := from + 1; i <= end; i++ {
					lines[i] = uncommentYAML(lines[i])
				}
				if root.Content[0].Kind == yaml.SequenceNode && indent(lines[from]) <= indent(lines[from + 1]) {
					lines[from] = strings.Replace(lines[from], "[]", "", 1)
				}
				if root.Content[0].Kind == yaml.MappingNode && indent(lines[from]) < indent(lines[from + 1]) {
					lines[from] = strings.Replace(lines[from], "{}", "", 1)
				}
			}
		}
		if err := yaml.Unmarshal([]byte(uncommentYAML(node.FootComment)), &root); err == nil && len(root.Content) > 0 {
			if root.Content[0].Kind == yaml.SequenceNode || root.Content[0].Kind == yaml.MappingNode {
				from := node.Line - 1
				end := from + lineCount(node.FootComment)
				for i := from + 1; i <= end; i++ {
					lines[i] = uncommentYAML(lines[i])
				}
				if root.Content[0].Kind == yaml.SequenceNode && indent(lines[from]) <= indent(lines[from + 1]) {
					lines[from] = strings.Replace(lines[from], "[]", "", 1)
				}
				if root.Content[0].Kind == yaml.MappingNode && indent(lines[from]) < indent(lines[from + 1]) {
					lines[from] = strings.Replace(lines[from], "{}", "", 1)
				}
			}
		}
		return true
	})

	if err := yaml.Unmarshal([]byte(strings.Join(lines, "\n")), &out); err != nil {
		return err
	}

	return nil
}