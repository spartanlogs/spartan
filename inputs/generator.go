package inputs

import (
	"github.com/lfkeitel/spartan/config"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
	"gopkg.in/tomb.v2"
)

func init() {
	Register("generator", newGeneratorInput)
}

type generatorConfig struct {
	message string
	lines   []string
	count   int
}

var generatorConfigSchema = []config.Setting{
	{
		Name: "message",
		Type: config.String,
	},
	{
		Name: "count",
		Type: config.Int,
	},
	{
		Name:     "lines",
		Type:     config.Array,
		ElemType: &config.Setting{Type: config.String},
	},
}

type GeneratorInput struct {
	config *generatorConfig
	t      tomb.Tomb
	out    chan<- *event.Event
}

func newGeneratorInput(options utils.InterfaceMap) (Input, error) {
	i := &GeneratorInput{
		config: &generatorConfig{},
	}
	return i, i.setConfig(options)
}

func (i *GeneratorInput) setConfig(options utils.InterfaceMap) error {
	if err := config.VerifySettings(options, generatorConfigSchema); err != nil {
		return err
	}

	i.config.message = options.Get("message").(string)
	i.config.count = options.Get("count").(int)
	i.config.lines = options.Get("lines").([]string)

	if len(i.config.lines) == 0 {
		i.config.lines = []string{i.config.message}
	}

	return nil
}

func (i *GeneratorInput) Start(out chan<- *event.Event) error {
	i.out = out
	i.t.Go(i.run)
	return nil
}

func (i *GeneratorInput) Close() error {
	i.t.Kill(nil)
	return i.t.Wait()
}

func (i *GeneratorInput) run() error {
	count := 0
	index := 0
	lineLen := len(i.config.lines)
	for {
		if i.config.count > 0 && count >= i.config.count {
			break
		}

		select {
		case <-i.t.Dying():
			return nil
		default:
		}

		i.out <- event.New(i.config.lines[index])
		index = (index + 1) % lineLen
		count++
	}

	// Once we've generated the desired number of events, wait
	// for the input to be closed.
	<-i.t.Dying()
	return nil
}
