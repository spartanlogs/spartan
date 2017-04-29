package outputs

import "github.com/lfkeitel/spartan/event"
import "github.com/lfkeitel/spartan/utils"

// The End output is a special output that is used for internal purposes
// to terminate an output chain. It simply returns as a no-op.
type end struct{}

func newEndOutput(options utils.InterfaceMap) (Output, error) {
	return &end{}, nil
}

// Run terminates an Output chain by immediately returning.
func (*end) Run(in []*event.Event) {}

// SetNext is a no-op since End is meant to terminate an Output chain.
func (*end) SetNext(n Output) {}
