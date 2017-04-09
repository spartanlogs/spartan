package codecs

import "github.com/lfkeitel/spartan/event"

type DotCodec struct{}

func init() {
	register("dot", &DotCodec{})
}

func (c *DotCodec) Format(e *event.Event) []byte {
	return []byte{'.'}
}
