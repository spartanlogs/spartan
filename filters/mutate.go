package filters

import (
	"errors"
	"fmt"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

var mutateActions = []string{"remove_field"}

type mutateConfig struct {
	fields []string
	action string
}

type MutateFilter struct {
	next   Filter
	config *mutateConfig
}

func NewMutateFilter(options map[string]interface{}) (*MutateFilter, error) {
	options = checkOptionsMap(options)
	f := &MutateFilter{config: &mutateConfig{}}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *MutateFilter) setConfig(options map[string]interface{}) error {
	if s, exists := options["fields"]; exists {
		switch s := s.(type) {
		case string:
			f.config.fields = []string{s}
		case []string:
			f.config.fields = s
		default:
			return errors.New("Fields must be a string or array of strings")
		}
	} else {
		return errors.New("Fields option required")
	}

	if s, exists := options["action"]; exists {
		f.config.action = s.(string)
		if !f.isValidAction(f.config.action) {
			return fmt.Errorf("%s is not a valid mutate action", f.config.action)
		}
	} else {
		return errors.New("Action option required")
	}

	return nil
}

func (f *MutateFilter) isValidAction(action string) bool {
	return utils.StringInSlice(action, mutateActions)
}

func (f *MutateFilter) SetNext(next Filter) {
	f.next = next
}

func (f *MutateFilter) Run(batch []*event.Event) []*event.Event {
	for _, event := range batch {
		switch f.config.action {
		case "remove_field":
			f.removeField(event)
		}
	}
	return f.next.Run(batch)
}

func (f *MutateFilter) removeField(e *event.Event) {
	for _, field := range f.config.fields {
		e.RemoveField(field)
	}
}
