package codecs

import (
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The LineCodec reads plaintext with delimiter of \n
type LineCodec struct{}

func init() {
	register("line", newLineCodec)
}

func newLineCodec() (Codec, error) {
	return &LineCodec{}, nil
}

// Encode event as a simple message.
func (c *LineCodec) Encode(e *event.Event) []byte {
	return []byte(e.String() + "\n")
}

// EncodeWriter reads events from in and writes them to w
func (c *LineCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode creates a new event with message set to data.
func (c *LineCodec) Decode(data []byte) (*event.Event, error) {
	return event.New(string(data)), nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *LineCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
