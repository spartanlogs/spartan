package event

import "bytes"

type templateNode interface {
	evaluate(e *Event) string
}

type template struct {
	nodes []templateNode
}

func newTemplate() *template {
	return &template{
		nodes: make([]templateNode, 0, 1),
	}
}

func (t *template) add(node templateNode) {
	t.nodes = append(t.nodes, node)
}

func (t *template) size() int {
	return len(t.nodes)
}

func (t *template) evaluate(e *Event) string {
	var buf bytes.Buffer

	for _, node := range t.nodes {
		buf.WriteString(node.evaluate(e))
	}

	return buf.String()
}
