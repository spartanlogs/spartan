package codecs

import (
	"fmt"

	"github.com/lfkeitel/spartan/event"
)

// A Codec is an object that can encode an Event to a byte slice
// and decode a byte slice into an Event.
type Codec interface {
	// Encode takes the given event and transforms it into a byte slice representation
	// depending on the codec itself.
	Encode(e *event.Event) []byte

	// Decode take a byte slice and attempts to turn it into an Event.
	Decode(data []byte) (*event.Event, error)
}

type codecInitFunc func() (Codec, error)

var registeredCodecs map[string]codecInitFunc

// register is an internal function for codecs to register their names
// and init functions.
func register(name string, c codecInitFunc) {
	if registeredCodecs == nil {
		registeredCodecs = make(map[string]codecInitFunc)
	}
	registeredCodecs[name] = c
}

// New will create an instance of the codec registered as name.
func New(name string) (Codec, error) {
	c, exists := registeredCodecs[name]
	if !exists {
		return nil, fmt.Errorf("Codec %s doesn't exist", name)
	}
	return c()
}
