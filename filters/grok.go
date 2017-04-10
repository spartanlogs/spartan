package filters

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lfkeitel/spartan/event"
)

func init() {
	register("grok", newGrokFilter)
}

type grokConfig struct {
	field string
	regex *regexp.Regexp
}

type GrokFilter struct {
	next   Filter
	config *grokConfig
}

func newGrokFilter(options map[string]interface{}) (Filter, error) {
	options = checkOptionsMap(options)
	g := &GrokFilter{config: &grokConfig{}}
	if err := g.setConfig(options); err != nil {
		return nil, err
	}
	return g, nil
}

func (f *GrokFilter) setConfig(options map[string]interface{}) error {
	if s, exists := options["field"]; exists {
		f.config.field = s.(string)
	} else {
		f.config.field = "message"
	}

	if s, exists := options["regex"]; exists {
		r, err := regexp.Compile(interpolatePatterns(s.(string)))
		if err != nil {
			return fmt.Errorf("Regex failed to compile: %v", err)
		}
		f.config.regex = r
	} else {
		return errors.New("Regex option required")
	}

	return nil
}

func (f *GrokFilter) SetNext(next Filter) {
	f.next = next
}

func (f *GrokFilter) Run(batch []*event.Event) []*event.Event {
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

		matches := f.config.regex.FindAllStringSubmatch(fieldStr, -1)
		if len(matches) == 0 {
			fmt.Printf("No matches")
			event.AddTag("_grokparsefailure")
			continue
		}

		for i, group := range f.config.regex.SubexpNames() {
			if i == 0 {
				continue
			}
			event.Set(group, matches[0][i])
		}
	}

	return f.next.Run(batch)
}
