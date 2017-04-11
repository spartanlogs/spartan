package inputs

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

// An Input generates events to be processed.
type Input interface {
	// Start creates a go routine for the Input to run in. out is an Event
	// channel where the Input will send generated Events.
	Start(out chan<- *event.Event) error

	// Close allows graceful shutdown of an Input.
	Close() error
}

type initFunc func(map[string]interface{}) (Input, error)

var (
	registeredInputInits map[string]initFunc

	// ErrInputNotRegistered is returned when attempting to create an unregistered Input.
	ErrInputNotRegistered = errors.New("Input doesn't exist")
)

func register(name string, init initFunc) {
	if registeredInputInits == nil {
		registeredInputInits = make(map[string]initFunc)
	}
	if _, exists := registeredInputInits[name]; exists {
		panic("Duplicate registration of input module: " + name)
	}
	registeredInputInits[name] = init
}

// New creates an instance of Input name with options. Options are dependent on the Input.
func New(name string, options map[string]interface{}) (Input, error) {
	init, exists := registeredInputInits[name]
	if !exists {
		return nil, ErrInputNotRegistered
	}
	return init(options)
}
