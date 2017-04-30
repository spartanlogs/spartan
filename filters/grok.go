package filters

import (
	"fmt"
	"regexp"

	"github.com/lfkeitel/spartan/config"
	"github.com/lfkeitel/spartan/event"
	"github.com/lfkeitel/spartan/utils"
)

func init() {
	Register("grok", newGrokFilter)
}

type grokConfig struct {
	field              string
	regex              []*regexp.Regexp
	ignoreMissingField bool
}

var grokConfigSchema = []config.Setting{
	{
		Name:    "field",
		Type:    config.String,
		Default: "message",
	},
	{
		Name:     "patterns",
		Type:     config.Array,
		Required: true,
		ElemType: &config.Setting{Type: config.String},
	},
	{
		Name:    "ignore_missing",
		Type:    config.Bool,
		Default: false,
	},
}

// A GrokFilter processes event fields based on give regex patterns.
// The first pattern to match is used for field data.
type GrokFilter struct {
	config *grokConfig
}

func newGrokFilter(options utils.InterfaceMap) (Filter, error) {
	options = checkOptionsMap(options)
	g := &GrokFilter{config: &grokConfig{}}
	if err := g.setConfig(options); err != nil {
		return nil, err
	}
	return g, nil
}

func (f *GrokFilter) setConfig(options utils.InterfaceMap) error {
	if err := config.VerifySettings(options, grokConfigSchema); err != nil {
		return err
	}

	f.config.field = options.Get("field").(string)
	f.config.ignoreMissingField = options.Get("ignore_missing").(bool)

	patterns := options.Get("patterns").([]string)
	f.config.regex = make([]*regexp.Regexp, len(patterns))
	var err error
	for i, pattern := range patterns {
		f.config.regex[i], err = regexp.Compile(interpolatePatterns(pattern))
		if err != nil {
			return fmt.Errorf("Regex failed to compile: %v", err)
		}
	}

	return nil
}

// Filter processes a batch.
func (f *GrokFilter) Filter(batch []*event.Event, matchedFunc MatchFunc) []*event.Event {
	for _, event := range batch {
		if event == nil {
			continue
		}

		field := event.GetField(f.config.field)
		if field == nil {
			if !f.config.ignoreMissingField {
				fmt.Printf("Field %s doesn't exist\n", f.config.field)
				event.AddTag("_grokparsefailure")
			}
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
				if i == 0 || group == "" {
					continue
				}
				event.SetField(group, matches[0][i])
			}
			break regexLoop
		}

		if matched {
			matchedFunc(event)
		} else {
			event.AddTag("_grokparsefailure")
		}
	}

	return batch
}
