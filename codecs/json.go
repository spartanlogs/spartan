package codecs

import (
	"encoding/json"
	"time"

	"github.com/lfkeitel/spartan/event"
)

type JsonCodec struct{}

func init() {
	register("json", &JsonCodec{})
}

func (c *JsonCodec) Format(e *event.Event) []byte {
	e.SetTimestamp(time.Unix(e.GetTimestamp().Unix(), 0))
	data := e.Squash()
	j, _ := json.Marshal(data)
	return j
}