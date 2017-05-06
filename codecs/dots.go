package codecs

import (
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The DotCodec "dot" converts an event to a single dot.
type DotCodec struct{}

func init() {
	register("dots", newDotCodec)
}

func newDotCodec() (Codec, error) {
	return &DotCodec{}, nil
}

// Encode event as a dot.
func (c *DotCodec) Encode(e *event.Event) []byte {
	return []byte{'.'}
}

// EncodeWriter reads events from in and writes them to w
func (c *DotCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode is effectivly a no-op. An event can't be decoded from a dot.
// Like seriously. A dot?
func (c *DotCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *DotCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
