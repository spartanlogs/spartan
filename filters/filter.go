package filters

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

type Filter interface {
	SetNext(Filter)
	Run([]*event.Event) []*event.Event
}

type InitFunc func(map[string]interface{}) (Filter, error)

var (
	registeredFilterInits map[string]InitFunc

	ErrFilterNotRegistered = errors.New("Filter doesn't exist")
)

func register(name string, init InitFunc) {
	if registeredFilterInits == nil {
		registeredFilterInits = make(map[string]InitFunc)
	}
	if _, exists := registeredFilterInits[name]; exists {
		panic("Duplicate registration of filter module: " + name)
	}
	registeredFilterInits[name] = init
}

func New(name string, options map[string]interface{}) (Filter, error) {
	init, exists := registeredFilterInits[name]
	if !exists {
		return nil, ErrFilterNotRegistered
	}
	return init(options)
}
