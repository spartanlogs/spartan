package outputs

import (
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/utils"
)

// The End output is a special output that is used for internal purposes
// to terminate an output chain. It simply returns as a no-op.
type end struct {
	BaseOutput
}

func newEndOutput(options utils.InterfaceMap) (Output, error) {
	return &end{}, nil
}

func (o *end) Run(batch []*event.Event) {}

func (o *end) LoadCodec(name string, options utils.InterfaceMap) error { return nil }
