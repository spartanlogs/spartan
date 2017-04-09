package outputs

import "github.com/lfkeitel/spartan/event"

type End struct{}

func (f *End) Run(in []*event.Event) {}

func (f *End) SetNext(n Output) {}
