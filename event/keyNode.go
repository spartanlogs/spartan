package event

import "fmt"

type keyNode struct {
	key string
}

func newKeyNode(key string) templateNode {
	return &keyNode{key: key}
}

func (n *keyNode) evaluate(e *Event) string {
	field := e.GetField(n.key)
	if field == nil {
		return fmt.Sprintf("%%{%s}", n.key)
	}
	return fmt.Sprintf("%s", field)
}
