package filters

import (
	"errors"

	"github.com/spartanlogs/spartan/config/parser"
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/utils"
)

// A Filter is used to manipulate and parse an Event. A Filter receives a batch of events
// and must return a slice of events. Events may be removed, edited, or added. Once a filter
// is done processing, it should call the next Filter in line with its processed event batch.
type Filter interface {
	// Run processes a batch.
	Filter(batch []*event.Event, matched MatchFunc) []*event.Event
}

// A Flushable is a filter that wishes to flush data periodically. Flush is called ever 5 seconds.
// Filter can track their last flush time and flush less often by simply returning early.
type Flushable interface {
	Flush() []*event.Event
}

// A FilterWrapper wraps the execution of a filter to enforce deadlines and other external functions.
type FilterWrapper interface {
	// SetNext sets the next Filter in line.
	SetNext(next FilterWrapper)

	// Run processes a batch.
	Run(batch []*event.Event) []*event.Event
}

// InitFunc is registered with this package as an initializer for a Filter
type InitFunc func(utils.InterfaceMap) (Filter, error)

// MatchFunc is ran by a filter for a matching event
type MatchFunc func(*event.Event)

var (
	registeredFilterInits map[string]InitFunc

	// ErrFilterNotRegistered is returned when attempting to create an unregistered Filter.
	ErrFilterNotRegistered = errors.New("Filter doesn't exist")
)

// Register allows filters to register an init function with their name
func Register(name string, init InitFunc) {
	if registeredFilterInits == nil {
		registeredFilterInits = make(map[string]InitFunc)
	}
	if _, exists := registeredFilterInits[name]; exists {
		panic("Duplicate registration of filter module: " + name)
	}
	registeredFilterInits[name] = init
}

// New creates an instance of Filter name with options. Options are dependent on the Filter.
func New(name string, options utils.InterfaceMap) (Filter, error) {
	init, exists := registeredFilterInits[name]
	if !exists {
		return nil, ErrFilterNotRegistered
	}
	return init(options)
}

// GeneratePipeline creates an filter pipeline. The returned Filter is the starting
// point in the pipeline. All other filters have been chained together in their
// defined order. An error will be returned if a filter doesn't exist.
func GeneratePipeline(defs []*parser.PipelineDef, batchsize int, pipelineNum int) (*Controller, error) {
	controller := newController(withBatchSize(batchsize))
	filters := make([]Filter, len(defs))

	if len(defs) == 0 {
		fw, _ := newFilterWrapper(nil, nil)
		p := NewPipeline(fw)
		controller.pipelines = append(controller.pipelines, p)
		return controller, nil
	}

	// Generate filters
	for i, def := range defs {
		filter, err := New(def.Module, def.Options)
		if err != nil {
			return nil, err
		}
		filters[i] = filter

		if f, flushable := filter.(Flushable); flushable {
			controller.addFlusher(f)
		}
	}

	wrappers := make([]FilterWrapper, len(defs))

	// Wrap filters
	for i, filter := range filters {
		switch len(defs[i].Connections) {
		case 0, 1: // Normal next filter
			fw, err := newFilterWrapper(filter, defs[i].Options)
			if err != nil {
				return nil, err
			}
			fw.id = defs[i].Module
			wrappers[i] = fw
		case 3: // If statement
			return nil, utils.ErrNotImplemented
		}
	}

	// Connect wrappers
	for i, wrapper := range wrappers[:len(wrappers)-1] {
		wrapper.SetNext(wrappers[i+1])
	}

	// Generate the desired number of pipelines
	for i := 0; i < pipelineNum; i++ {
		p := NewPipeline(wrappers[0])
		controller.pipelines = append(controller.pipelines, p)
	}

	return controller, nil
}

// checkOptionsMap ensures an option map is never nil.
func checkOptionsMap(o utils.InterfaceMap) utils.InterfaceMap {
	if o == nil {
		o = utils.NewInterfaceMap()
	}
	return o
}
