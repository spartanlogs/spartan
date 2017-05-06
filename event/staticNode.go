package event

type staticNode struct {
	content string
}

func newStaticNode(content string) templateNode {
	return &staticNode{content: content}
}

func (n *staticNode) evaluate(e *Event) string { return n.content }
