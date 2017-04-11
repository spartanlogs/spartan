package outputs

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

// An Output takes a batch of events and displays or transports them out of the system.
type Output interface {
	// SetNext sets the next Output in line.
	SetNext(next Output)

	// Run processes a batch.
	Run(batch []*event.Event)
}

type initFunc func(map[string]interface{}) (Output, error)

var (
	registeredOutputInits map[string]initFunc

	// ErrOutputNotRegistered is returned when attempting to create an unregistered Output.
	ErrOutputNotRegistered = errors.New("Output doesn't exist")
)

func register(name string, init initFunc) {
	if registeredOutputInits == nil {
		registeredOutputInits = make(map[string]initFunc)
	}
	if _, exists := registeredOutputInits[name]; exists {
		panic("Duplicate registration of output module: " + name)
	}
	registeredOutputInits[name] = init
}

// New creates an instance of Output name with options. Options are dependent on the Output.
func New(name string, options map[string]interface{}) (Output, error) {
	init, exists := registeredOutputInits[name]
	if !exists {
		return nil, ErrOutputNotRegistered
	}
	return init(options)
}
