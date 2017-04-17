package filters

import (
	"context"

	"github.com/lfkeitel/spartan/event"
)

type filterWrapper struct {
	cmd  Filter
	next Filter
}

func newFilterWrapper(cmd Filter) *filterWrapper {
	return &filterWrapper{
		cmd: cmd,
	}
}

// SetNext sets the next Filter in line.
func (f *filterWrapper) SetNext(next FilterWrapper) {
	f.next = next
}

// Run processes a batch.
func (f *filterWrapper) Run(ctx context.Context, batch []*event.Event) []*event.Event {
	// It's the end of the pipeline as we know it.
	if f.cmd == nil {
		return batch
	}

	// For now just pass along the data
	return f.next.Run(ctx, f.cmd.Run(ctx, batch))
}
