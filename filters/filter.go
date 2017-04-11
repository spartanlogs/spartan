package filters

import (
	"errors"

	"github.com/lfkeitel/spartan/event"
)

// A Filter is used to manipulate and parse an Event. A Filter receives a batch of events
// and must return a slice of events. Events may be removed, edited, or added. Once a filter
// is done processing, it should call the next Filter in line with its processed event batch.
type Filter interface {
	// SetNext sets the next Filter in line.
	SetNext(next Filter)

	// Run processes a batch.
	Run(batch []*event.Event) []*event.Event
}

type initFunc func(map[string]interface{}) (Filter, error)

var (
	registeredFilterInits map[string]initFunc

	// ErrFilterNotRegistered is returned when attempting to create an unregistered Filter.
	ErrFilterNotRegistered = errors.New("Filter doesn't exist")
)

func register(name string, init initFunc) {
	if registeredFilterInits == nil {
		registeredFilterInits = make(map[string]initFunc)
	}
	if _, exists := registeredFilterInits[name]; exists {
		panic("Duplicate registration of filter module: " + name)
	}
	registeredFilterInits[name] = init
}

// New creates an instance of Filter name with options. Options are dependent on the Filter.
func New(name string, options map[string]interface{}) (Filter, error) {
	init, exists := registeredFilterInits[name]
	if !exists {
		return nil, ErrFilterNotRegistered
	}
	return init(options)
}
