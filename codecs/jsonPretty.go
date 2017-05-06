package codecs

import (
	"encoding/json"
	"io"

	"github.com/spartanlogs/spartan/event"
)

// The JSONPrettyCodec encodes/decodes an event as formatted, pretty JSON.
type JSONPrettyCodec struct{}

func init() {
	register("json_pretty", newJSONPrettyCodec)
}

func newJSONPrettyCodec() (Codec, error) {
	return &JSONPrettyCodec{}, nil
}

// Encode Event as JSON object.
func (c *JSONPrettyCodec) Encode(e *event.Event) []byte {
	data := e.Data()
	j, _ := json.MarshalIndent(data, "", "  ")
	return j
}

// EncodeWriter reads events from in and writes them to w
func (c *JSONPrettyCodec) EncodeWriter(w io.Writer, in <-chan *event.Event) {}

// Decode byte slice into an Event. CURRENTLY NOT IMPLEMENTED.
func (c *JSONPrettyCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}

// DecodeReader reads from r and creates an event sent to out
func (c *JSONPrettyCodec) DecodeReader(r io.Reader, out chan<- *event.Event) {}
