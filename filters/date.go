package filters

import (
	"context"
	"errors"
	"time"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	register("date", newDateFilter)
}

type dateConfig struct {
	field    string
	patterns []string
	timezone string
}

// The DateFilter is used to set the canonical @timestamp field of an Event.
// A field is tested against an array of date patterns and if on matches,
// the resulting parsed time is set as the Events timestamp.
type DateFilter struct {
	config *dateConfig
}

func newDateFilter(options *utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	f := &DateFilter{config: &dateConfig{}}
	if err := f.setConfig(options); err != nil {
		return nil, err
	}
	return f, nil
}

func (f *DateFilter) setConfig(options *utils.InterfaceMap) error {
	if s, exists := options.GetOK("field"); exists {
		f.config.field = s.(string)
	} else {
		return errors.New("Field option required")
	}

	if s, exists := options.GetOK("patterns"); exists {
		switch s := s.(type) {
		case string:
			f.config.patterns = []string{s}
		case []string:
			f.config.patterns = s
		default:
			return errors.New("Patterns must be a string or array of strings")
		}
	} else {
		return errors.New("Patterns option required")
	}

	if s, exists := options.GetOK("timezone"); exists {
		f.config.timezone = s.(string)

		_, err := time.LoadLocation(f.config.timezone)
		if err != nil {
			return errors.New("Invalid timezone")
		}
	} else {
		f.config.timezone = "UTC"
	}

	return nil
}

// Filter processes a batch.
func (f *DateFilter) Filter(ctx context.Context, batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for _, event := range batch {
		field := event.Get(f.config.field)
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
