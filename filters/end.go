package filters

import "github.com/lfkeitel/spartan/event"
import "github.com/lfkeitel/spartan/utils"

func init() {
	register("end", newEndFilter)
}

// The End filter is a special filter that is used for internal purposes
// to terminate a filter chain. It simply returns any batch it's given.
type End struct{}

func newEndFilter(options *utils.InterfaceMap) (Filter, error) {
	return &End{}, nil
}

// Run immediately returns the given batch.
func (f *End) Run(batch []*event.Event) []*event.Event {
	return batch
}

// SetNext is a no-op since End terminates a filter pipeline.
func (f *End) SetNext(n Filter) {}
