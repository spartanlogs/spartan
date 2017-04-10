package filters

import "github.com/lfkeitel/spartan/event"

func init() {
	register("end", newEndFilter)
}

type End struct{}

func newEndFilter(options map[string]interface{}) (Filter, error) {
	return &End{}, nil
}

func (f *End) Run(in []*event.Event) []*event.Event {
	return in
}

func (f *End) SetNext(n Filter) {}
