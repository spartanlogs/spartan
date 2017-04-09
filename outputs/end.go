package outputs

import "github.com/lfkeitel/spartan/common"

type End struct{}

func (f *End) Run(in []*common.Event) {}

func (f *End) SetNext(n Output) {}
