package outputs

import (
	"errors"
	"fmt"

	"github.com/spartanlogs/spartan/codecs"
	"github.com/spartanlogs/spartan/config/parser"
	"github.com/spartanlogs/spartan/event"
	"github.com/spartanlogs/spartan/utils"
)

// An Output takes a batch of events and displays or transports them out of the system.
type Output interface {
	// SetNext sets the next Output in line.
	SetNext(next Output)

	// Run processes a batch.
	Run(batch []*event.Event)

	// LoadCodec will create a codec object for the Output. The output can disable
	// codec creation by redeclaring this function and simply returning nil.
	LoadCodec(name string, options utils.InterfaceMap) error
}

// BaseOutput can be embeded in other output plugins to provide utility functionality.
// The Run method must be defined by the child struct. The LoadCodec method can be
// "disabled" by redefining it on the struct and returning nil.
type BaseOutput struct {
	Codec codecs.Codec
	Next  Output
}

// SetNext sets the next Output in line.
func (o *BaseOutput) SetNext(next Output) { o.Next = next }

// LoadCodec creates a codec object for the output
func (o *BaseOutput) LoadCodec(name string, options utils.InterfaceMap) error {
	if name == "" {
		name = "plain"
	}

	codec, err := codecs.New(name, options)
	if err != nil {
		return err
	}
	o.Codec = codec
	return nil
}

type initFunc func(utils.InterfaceMap) (Output, error)

var (
	registeredOutputInits map[string]initFunc

	// ErrOutputNotRegistered is returned when attempting to create an unregistered Output.
	ErrOutputNotRegistered = errors.New("Output doesn't exist")
)

// Register allows an output to register an init function with their name
func Register(name string, init initFunc) {
	if registeredOutputInits == nil {
		registeredOutputInits = make(map[string]initFunc)
	}
	if _, exists := registeredOutputInits[name]; exists {
		panic("Duplicate registration of output module: " + name)
	}
	registeredOutputInits[name] = init
}

// New creates an instance of Output name with options. Options are dependent on the Output.
func New(name string, options utils.InterfaceMap) (Output, error) {
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

		// Create and set the codec for the plugin
		codecOption := def.Options.Get("codec")
		codecName := ""
		if codecOption == nil {
			codecName = "plain"
		} else if cn, ok := codecOption.(string); ok {
			codecName = cn
		} else {
			return nil, fmt.Errorf("invalid codec setting in %s plugin", def.Module)
		}

		if err := output.LoadCodec(codecName, def.CodecOptions); err != nil {
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
