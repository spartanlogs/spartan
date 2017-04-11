package codecs

import "github.com/lfkeitel/spartan/event"

// The DotCodec "dot" converts an event to a single dot.
type DotCodec struct{}

func init() {
	register("dot", newDotCodec)
}

func newDotCodec() (Codec, error) {
	return &DotCodec{}, nil
}

// Encode event as a dot.
func (c *DotCodec) Encode(e *event.Event) []byte {
	return []byte{'.'}
}

// Decode is effectivly a no-op. An event can't be decoded from a dot.
// Like seriously. A dot?
func (c *DotCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}
