package filters

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	register("grok", newGrokFilter)
}

type grokConfig struct {
	field string
	regex []*regexp.Regexp
}

// A GrokFilter processes event fields based on give regex patterns.
// The first pattern to match is used for field data.
type GrokFilter struct {
	config *grokConfig
}

func newGrokFilter(options *utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	g := &GrokFilter{config: &grokConfig{}}
	if err := g.setConfig(options); err != nil {
		return nil, err
	}
	return g, nil
}

func (f *GrokFilter) setConfig(options *utils.InterfaceMap) error {
	if s, exists := options.GetOK("field"); exists {
		f.config.field = s.(string)
	} else {
		f.config.field = "message"
	}

	if s, exists := options.GetOK("patterns"); exists {
		var patterns []string
		switch sTyped := s.(type) {
		case []string:
			patterns = sTyped
		case string:
			patterns = []string{sTyped}
		default:
			return errors.New("regex must be string or array of strings")
		}

		for _, pattern := range patterns {
			r, err := regexp.Compile(interpolatePatterns(pattern))
			if err != nil {
				return fmt.Errorf("Regex failed to compile: %v", err)
			}
			f.config.regex = append(f.config.regex, r)
		}
	} else {
		return errors.New("Regex option required")
	}

	return nil
}

// Run processes a batch.
func (f *GrokFilter) Run(ctx context.Context, batch []*event.Event) []*event.Event {
	for _, event := range batch {
		if event == nil {
			continue
		}

		field := event.Get(f.config.field)
		if field == nil {
			fmt.Printf("Field %s doesn't exist\n", f.config.field)
			event.AddTag("_grokparsefailure")
			continue
		}

		fieldStr, ok := field.(string)
		if !ok {
			fmt.Printf("Field %s isn't a string\n", f.config.field)
			event.AddTag("_grokparsefailure")
			continue
		}

		matched := false
	regexLoop:
		for _, regex := range f.config.regex {
			matches := regex.FindAllStringSubmatch(fieldStr, -1)
			if len(matches) == 0 {
				continue
			}
			matched = true

			for i, group := range regex.SubexpNames() {
				if i == 0 {
					continue
				}
				event.Set(group, matches[0][i])
			}
			break regexLoop
		}

		if !matched {
			fmt.Printf("No matches")
			event.AddTag("_grokparsefailure")
		}
	}

	return batch
}
