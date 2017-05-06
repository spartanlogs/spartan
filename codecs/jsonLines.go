package codecs

import (
	"encoding/json"
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The JSONLineCodec encodes/decodes an event as JSON.
type JSONLineCodec struct{}

func init() {
	register("json_lines", newJSONLineCodec)
}

func newJSONLineCodec() (Codec, error) {
	return &JSONLineCodec{}, nil
}

// Encode Event as JSON object.
func (c *JSONLineCodec) Encode(e *event.Event) []byte {
	data := e.Data()
	j, _ := json.Marshal(data)
	return append(j, '\n')
}

// EncodeWriter reads events from in and writes them to w
func (c *JSONLineCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode byte slice into an Event. CURRENTLY NOT IMPLEMENTED.
func (c *JSONLineCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *JSONLineCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
