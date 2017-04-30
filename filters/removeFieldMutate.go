package filters

import (
	"github.com/lfkeitel/spartan/config"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	Register("remove_field", newRemoveFieldFilter)
}

type removeFieldConfig struct {
	fields []string
}

var removeFieldConfigSchema = []config.Setting{
	{
		Name:     "fields",
		Type:     config.Array,
		Required: true,
		ElemType: &config.Setting{Type: config.String},
	},
}

// A RemoveFieldFilter is used to perform several different actions on an Event.
// See the documentation for the Mutate filter for more information.
type RemoveFieldFilter struct {
	config *removeFieldConfig
}

func newRemoveFieldFilter(options utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	f := &RemoveFieldFilter{config: &removeFieldConfig{}}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *RemoveFieldFilter) setConfig(options utils.InterfaceMap) error {
	if err := config.VerifySettings(options, removeFieldConfigSchema); err != nil {
		return err
	}

	f.config.fields = options.Get("fields").([]string)

	return nil
}

// Filter processes a batch.
func (f *RemoveFieldFilter) Filter(batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for _, event := range batch {
		for _, field := range f.config.fields {
			event.DeleteField(field)
		}
	}
	return batch
}
