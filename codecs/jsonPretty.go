package codecs

import (
	"encoding/json"

	"github.com/lfkeitel/spartan/event"
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

// Decode byte slice into an Event. CURRENTLY NOT IMPLEMENTED.
func (c *JSONPrettyCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}
