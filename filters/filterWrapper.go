package filters

import (
	"context"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

type filterWrapper struct {
	cmd  Filter
	next FilterWrapper
}

func newFilterWrapper(cmd Filter, options utils.InterfaceMap) (*filterWrapper, error) {
	options = checkOptionsMap(options)
	fw := &filterWrapper{
		cmd: cmd,
	}
	if err := fw.setConfig(options); err != nil {
		return nil, err
	}
	return fw, nil
}

func (f *filterWrapper) setConfig(options utils.InterfaceMap) error {
	return nil
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
	return f.next.Run(ctx, f.cmd.Filter(ctx, batch, f.matchFunc))
}

func (f *filterWrapper) matchFunc(e *event.Event) {}
