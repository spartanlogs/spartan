package inputs

import (
	"errors"

	"github.com/spartanlogs/spartan/config/parser"
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/utils"
)

// An Input generates events to be processed.
type Input interface {
	// Start creates a go routine for the Input to run in. out is an Event
	// channel where the Input will send generated Events.
	Start(out chan<- *event.Event) error

	// Close allows graceful shutdown of an Input.
	Close() error
}

type initFunc func(utils.InterfaceMap) (Input, error)

var (
	registeredInputInits map[string]initFunc

	// ErrInputNotRegistered is returned when attempting to create an unregistered Input.
	ErrInputNotRegistered = errors.New("Input doesn't exist")
)

// Register allows an input to register an init function with their name
func Register(name string, init initFunc) {
	if registeredInputInits == nil {
		registeredInputInits = make(map[string]initFunc)
	}
	if _, exists := registeredInputInits[name]; exists {
		panic("Duplicate registration of input module: " + name)
	}
	registeredInputInits[name] = init
}

// New creates an instance of Input name with options. Options are dependent on the Input.
func New(name string, options utils.InterfaceMap) (Input, error) {
	init, exists := registeredInputInits[name]
	if !exists {
		return nil, ErrInputNotRegistered
	}
	return init(options)
}

// CreateFromDefs initalizes the inputs defined in the give slice.
// An error will return if an input module doesn't exist.
func CreateFromDefs(defs []*parser.InputDef) ([]Input, error) {
	inputs := make([]Input, len(defs))

	for i, def := range defs {
		input, err := New(def.Module, def.Options)
		if err != nil {
			return nil, err
		}
		inputs[i] = input
	}

	return inputs, nil
}
