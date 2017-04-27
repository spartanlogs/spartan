package filters

import (
	"context"
	"errors"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	register("remove_field", newRemoveFieldFilter)
}

type removeFieldConfig struct {
	fields []string
	action string
}

// A RemoveFieldFilter is used to perform several different actions on an Event.
// See the documentation for the Mutate filter for more information.
type RemoveFieldFilter struct {
	config *removeFieldConfig
}

func newRemoveFieldFilter(options *utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	f := &RemoveFieldFilter{config: &removeFieldConfig{}}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *RemoveFieldFilter) setConfig(options *utils.InterfaceMap) error {
	if s, exists := options.GetOK("fields"); exists {
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

	return nil
}

// Filter processes a batch.
func (f *RemoveFieldFilter) Filter(ctx context.Context, batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for _, event := range batch {
		for _, field := range f.config.fields {
			event.RemoveField(field)
		}
	}
	return batch
}
