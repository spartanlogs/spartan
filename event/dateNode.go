package event

type dateNode struct {
	format string
}

func newDateNode(format string) templateNode {
	return &dateNode{format: format}
}

func (n *dateNode) evaluate(e *Event) string {
	return e.GetTimestamp().UTC().Format(n.format)
}
