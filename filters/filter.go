package filters

import (
	"context"
	"errors"

	"github.com/lfkeitel/spartan/config/parser"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

// A Filter is used to manipulate and parse an Event. A Filter receives a batch of events
// and must return a slice of events. Events may be removed, edited, or added. Once a filter
// is done processing, it should call the next Filter in line with its processed event batch.
type Filter interface {
	// Run processes a batch.
	Run(ctx context.Context, batch []*event.Event) []*event.Event
}

// A FilterWrapper wraps the execution of a filter to enforce deadlines and other external functions.
type FilterWrapper interface {
	// SetNext sets the next Filter in line.
	SetNext(next FilterWrapper)

	// Run processes a batch.
	Run(ctx context.Context, batch []*event.Event) []*event.Event
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
func GeneratePipeline(defs []*parser.PipelineDef) (FilterWrapper, error) {
	filters := make([]Filter, len(defs))

	if len(defs) == 0 {
		return newFilterWrapper(nil), nil
	}

	// Generate filters
	for i, def := range defs {
		filter, err := New(def.Module, def.Options)
		if err != nil {
			return nil, err
		}
		filters[i] = filter
	}

	wrappers := make([]FilterWrapper, len(defs))

	// Wrap filters
	for i, filter := range filters {
		switch len(defs[i].Connections) {
		case 0: // End of a pipeline
			wrappers[i] = newFilterWrapper(nil)
		case 1: // Normal next filter
			wrappers[i] = newFilterWrapper(filter)
		case 3: // If statement
			return nil, utils.ErrNotImplemented
		}
	}

	// Connect wrappers
	for i, wrapper := range wrappers {
		if i < len(wrappers)-1 {
			wrapper.SetNext(wrappers[i+1])
		}
	}

	return wrappers[0], nil
}

// checkOptionsMap ensures an option map is never nil.
func checkOptionsMap(o *utils.InterfaceMap) *utils.InterfaceMap {
	if o == nil {
		o = utils.NewInterfaceMap()
	}
	return o
}
