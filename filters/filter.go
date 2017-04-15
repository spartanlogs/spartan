package filters

import (
	"errors"

	"github.com/lfkeitel/spartan/config/parser"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
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

type initFunc func(*utils.InterfaceMap) (Filter, error)

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
func New(name string, options *utils.InterfaceMap) (Filter, error) {
	init, exists := registeredFilterInits[name]
	if !exists {
		return nil, ErrFilterNotRegistered
	}
	return init(options)
}

// GeneratePipeline creates an filter pipeline. The returned Filter is the starting
// point in the pipeline. All other filters have been chained together in their
// defined order. An error will be returned if a filter doesn't exist.
func GeneratePipeline(defs []*parser.PipelineDef) (Filter, error) {
	filters := make([]Filter, len(defs))

	// Generate filters
	for i, def := range defs {
		filter, err := New(def.Module, def.Options)
		if err != nil {
			return nil, err
		}
		filters[i] = filter
	}

	// Connect filters
	for i, filter := range filters {
		switch len(defs[i].Connections) {
		case 0: // End of a pipeline
			filter.SetNext(&End{})
		case 1: // Normal next filter
			filter.SetNext(filters[defs[i].Connections[0]])
		case 3: // If statement
			return nil, utils.ErrNotImplemented
		}
	}

	return filters[0], nil
}
