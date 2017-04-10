package outputs

import "github.com/lfkeitel/spartan/event"

func init() {
	register("end", newEndOutput)
}

type End struct{}

func newEndOutput(options map[string]interface{}) (Output, error) {
	return &End{}, nil
}

func (f *End) Run(in []*event.Event) {}

func (f *End) SetNext(n Output) {}
