package inputs

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

type Input interface {
	Start(chan<- *event.Event) error
	Close() error
}

type InitFunc func(map[string]interface{}) (Input, error)

var (
	registeredInputInits map[string]InitFunc

	ErrInputNotRegistered = errors.New("Input doesn't exist")
)

func register(name string, init InitFunc) {
	if registeredInputInits == nil {
		registeredInputInits = make(map[string]InitFunc)
	}
	if _, exists := registeredInputInits[name]; exists {
		panic("Duplicate registration of input module: " + name)
	}
	registeredInputInits[name] = init
}

func New(name string, options map[string]interface{}) (Input, error) {
	init, exists := registeredInputInits[name]
	if !exists {
		return nil, ErrInputNotRegistered
	}
	return init(options)
}
