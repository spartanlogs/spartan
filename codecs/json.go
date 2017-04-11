package codecs

import (
	"encoding/json"
	"time"

	"github.com/lfkeitel/spartan/event"
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
	e.SetTimestamp(time.Unix(e.GetTimestamp().Unix(), 0))
	data := e.Squash()
	j, _ := json.Marshal(data)
	return j
}

// Decode byte slice into an Event. CURRENTLY NOT IMPLEMENTED.
func (c *JSONCodec) Decode(data []byte) (*event.Event, error) {
	return nil, nil
}
