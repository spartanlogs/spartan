package codecs

import (
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The PlainCodec reads plaintext with no delimiting between events
type PlainCodec struct{}

func init() {
	register("plain", newPlainCodec)
}

func newPlainCodec() (Codec, error) {
	return &PlainCodec{}, nil
}

// Encode event as a simple message.
func (c *PlainCodec) Encode(e *event.Event) []byte {
	return []byte(e.String())
}

// EncodeWriter reads events from in and writes them to w
func (c *PlainCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode creates a new event with message set to data.
func (c *PlainCodec) Decode(data []byte) (*event.Event, error) {
	return event.New(string(data)), nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *PlainCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
