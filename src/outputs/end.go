package outputs

import "github.com/lfkeitel/spartan/src/common"

type End struct{}

func (f *End) Run(in []*common.Event) {}

func (f *End) SetNext(n Output) {}
