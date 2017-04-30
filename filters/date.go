package filters

import (
	"time"

	"github.com/lfkeitel/spartan/config"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	Register("date", newDateFilter)
}

type dateConfig struct {
	field    string
	patterns []string
	timezone string
}

var dateConfigSchema = []config.Setting{
	{
		Name:    "field",
		Type:    config.String,
		Default: "message",
	},
	{
		Name:    "timezone",
		Type:    config.String,
		Default: "UTC",
	},
	{
		Name:     "patterns",
		Type:     config.Array,
		Required: true,
		ElemType: &config.Setting{Type: config.String},
	},
}

// The DateFilter is used to set the canonical @timestamp field of an Event.
// A field is tested against an array of date patterns and if on matches,
// the resulting parsed time is set as the Events timestamp.
type DateFilter struct {
	config *dateConfig
}

func newDateFilter(options utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	f := &DateFilter{config: &dateConfig{}}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *DateFilter) setConfig(options utils.InterfaceMap) error {
	if err := config.VerifySettings(options, dateConfigSchema); err != nil {
		return err
	}

	f.config.field = options.Get("field").(string)
	f.config.timezone = options.Get("timezone").(string)
	f.config.patterns = options.Get("patterns").([]string)

	return nil
}

// Filter processes a batch.
func (f *DateFilter) Filter(batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for _, event := range batch {
		field := event.GetField(f.config.field)
		if field == nil {
			continue
		}

		fieldStr, ok := field.(string)
		if !ok {
			continue
		}

		loc, _ := time.LoadLocation(f.config.timezone)

		matched := false
		for _, p := range f.config.patterns {
			newTime, err := time.ParseInLocation(p, fieldStr, loc)
			if err != nil {
				continue
			}
			event.SetTimestamp(newTime)
			matched = true
			break
		}

		if !matched {
			continue
		}

		matchedFunc(event)
	}
	return batch
}
