package outputs

import "github.com/lfkeitel/spartan/event"
import "github.com/lfkeitel/spartan/utils"

func init() {
	register("end", newEndOutput)
}

// The End output is a special output that is used for internal purposes
// to terminate an output chain. It simply returns as a no-op.
type End struct{}

func newEndOutput(options *utils.InterfaceMap) (Output, error) {
	return &End{}, nil
}

// Run terminates an Output chain by immediately returning.
func (f *End) Run(in []*event.Event) {}

// SetNext is a no-op since End is meant to terminate an Output chain.
func (f *End) SetNext(n Output) {}
