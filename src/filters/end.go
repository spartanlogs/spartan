package filters

import "github.com/lfkeitel/spartan/src/common"

type End struct{}

func (f *End) Run(in []*common.Event) []*common.Event {
	return in
}

func (f *End) SetNext(n Filter) {}
