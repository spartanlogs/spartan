package filters

import (
	"context"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

// The End filter is a special filter that is used for internal purposes
// to terminate a filter chain. It simply returns any batch it's given.
type end struct{}

func newEndFilter(options *utils.InterfaceMap) (Filter, error) {
	return &end{}, nil
}

// Run immediately returns the given batch.
func (*end) Run(ctx context.Context, batch []*event.Event) []*event.Event {
	return batch
}
