package filters

import "github.com/lfkeitel/spartan/event"

type End struct{}

func (f *End) Run(in []*event.Event) []*event.Event {
	return in
}

func (f *End) SetNext(n Filter) {}
