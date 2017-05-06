package event

import "fmt"

type epochNode struct{}

func newEpochNode() templateNode {
	return &epochNode{}
}

func (n *epochNode) evaluate(e *Event) string {
	return fmt.Sprintf("%s", e.GetTimestamp().Unix())
}
