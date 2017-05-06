package filters

import (
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/utils"
)

type filterWrapper struct {
	id   string
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

func (f *filterWrapper) GetID() string {
	return f.id
}

// SetNext sets the next Filter in line.
func (f *filterWrapper) SetNext(next FilterWrapper) {
	f.next = next
}

// Run processes a batch.
func (f *filterWrapper) Run(batch []*event.Event) []*event.Event {
	//fmt.Printf("Wrapper %s running...\n", f.id)

	if f.cmd == nil {
		return batch
	}
	batch = f.cmd.Filter(batch, f.matchFunc)

	if f.next == nil {
		return batch
	}
	return f.next.Run(batch)
}

func (f *filterWrapper) matchFunc(e *event.Event) {}
