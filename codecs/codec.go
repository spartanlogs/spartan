package codecs

import (
	"errors"

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

var (
	registeredCodecInits map[string]codecInitFunc

	// ErrCodecNotRegistered is returned when attempting to create an unregistered Codec.
	ErrCodecNotRegistered = errors.New("Codec doesn't exist")
)

// register is an internal function for codecs to register their names
// and init functions.
func register(name string, c codecInitFunc) {
	if registeredCodecInits == nil {
		registeredCodecInits = make(map[string]codecInitFunc)
	}
	if _, exists := registeredCodecInits[name]; exists {
		panic("Duplicate registration of filter module: " + name)
	}
	registeredCodecInits[name] = c
}

// New will create an instance of the codec registered as name.
func New(name string) (Codec, error) {
	c, exists := registeredCodecInits[name]
	if !exists {
		return nil, ErrCodecNotRegistered
	}
	return c()
}
