package codecs

import (
	"encoding/json"
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The JSONCodec encodes/decodes an event as JSON.
type JSONCodec struct{}

func init() {
	register("json", newJSONCodec)
}

func newJSONCodec() (Codec, error) {
	return &JSONCodec{}, nil
}

// Encode Event as JSON object.
func (c *JSONCodec) Encode(e *event.Event) []byte {
	data := e.Data()
	j, _ := json.Marshal(data)
	return j
}

// EncodeWriter reads events from in and writes them to w
func (c *JSONCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode byte slice into an Event. CURRENTLY NOT IMPLEMENTED.
func (c *JSONCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *JSONCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
