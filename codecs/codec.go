package codecs

import (
	"fmt"

	"github.com/lfkeitel/spartan/event"
)

type Codec interface {
	Format(*event.Event) []byte
}

var registeredCodecs map[string]Codec

func register(name string, c Codec) {
	if registeredCodecs == nil {
		registeredCodecs = make(map[string]Codec)
	}
	registeredCodecs[name] = c
}

func New(name string) (Codec, error) {
	c, exists := registeredCodecs[name]
	if !exists {
		return nil, fmt.Errorf("Codec %s doesn't exist", name)
	}
	return c, nil
}
