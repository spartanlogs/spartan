package outputs

import (
	"errors"

	"github.com/lfkeitel/spartan/config/parser"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

// An Output takes a batch of events and displays or transports them out of the system.
type Output interface {
	// SetNext sets the next Output in line.
	SetNext(next Output)

	// Run processes a batch.
	Run(batch []*event.Event)
}

type initFunc func(*utils.InterfaceMap) (Output, error)

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
func New(name string, options *utils.InterfaceMap) (Output, error) {
	init, exists := registeredOutputInits[name]
	if !exists {
		return nil, ErrOutputNotRegistered
	}
	return init(options)
}

// GeneratePipeline creates an output pipeline. The returned Output is the starting
// point in the pipeline. All other outputs have been chained together in their
// defined order. An error will be returned if an output doesn't exist.
func GeneratePipeline(defs []*parser.PipelineDef) (Output, error) {
	outputs := make([]Output, len(defs))

	if len(defs) == 0 {
		return &end{}, nil
	}

	// Generate outputs
	for i, def := range defs {
		output, err := New(def.Module, def.Options)
		if err != nil {
			return nil, err
		}
		outputs[i] = output
	}

	// Connect outputs
	for i, output := range outputs {
		switch len(defs[i].Connections) {
		case 0: // End of a pipeline
			output.SetNext(&end{})
		case 1: // Normal next output
			output.SetNext(outputs[defs[i].Connections[0]])
		case 3: // If statement
			return nil, utils.ErrNotImplemented
		}
	}

	return outputs[0], nil
}

// checkOptionsMap ensures an option map is never nil.
func checkOptionsMap(o *utils.InterfaceMap) *utils.InterfaceMap {
	if o == nil {
		o = utils.NewInterfaceMap()
	}
	return o
}
