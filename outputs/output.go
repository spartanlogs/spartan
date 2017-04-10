package outputs

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

type Output interface {
	SetNext(Output)
	Run([]*event.Event)
}

type InitFunc func(map[string]interface{}) (Output, error)

var (
	registeredOutputInits map[string]InitFunc

	ErrOutputNotRegistered = errors.New("Output doesn't exist")
)

func register(name string, init InitFunc) {
	if registeredOutputInits == nil {
		registeredOutputInits = make(map[string]InitFunc)
	}
	if _, exists := registeredOutputInits[name]; exists {
		panic("Duplicate registration of output module: " + name)
	}
	registeredOutputInits[name] = init
}

func New(name string, options map[string]interface{}) (Output, error) {
	init, exists := registeredOutputInits[name]
	if !exists {
		return nil, ErrOutputNotRegistered
	}
	return init(options)
}
