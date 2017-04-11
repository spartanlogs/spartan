package outputs

import (
	"fmt"

	"github.com/lfkeitel/spartan/codecs"
	"github.com/lfkeitel/spartan/event"
)

func init() {
	register("stdout", newStdOutOutput)
}

type stdOutConfig struct {
	codec codecs.Codec
}

// StdOutOutput prints events to StdOut.
type StdOutOutput struct {
	config *stdOutConfig
	next   Output
}

func newStdOutOutput(options map[string]interface{}) (Output, error) {
	options = checkOptionsMap(options)
	o := &StdOutOutput{config: &stdOutConfig{}}
	if err := o.setConfig(options); err != nil {
		return nil, err
	}
	return o, nil
}

func (o *StdOutOutput) setConfig(options map[string]interface{}) error {
	if s, exists := options["codec"]; exists {
		c, err := codecs.New(s.(string))
		if err != nil {
			return err
		}
		o.config.codec = c
	} else {
		c, _ := codecs.New("json")
		o.config.codec = c
	}

	return nil
}

// SetNext sets the next Output in line.
func (o *StdOutOutput) SetNext(next Output) {
	o.next = next
}

// Run processes a batch.
func (o *StdOutOutput) Run(batch []*event.Event) {
	for _, event := range batch {
		if event != nil {
			fmt.Printf("%s\n", o.config.codec.Encode(event))
		}
	}
	o.next.Run(batch)
}
