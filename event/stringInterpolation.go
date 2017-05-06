package event

import (
	"regexp"
	"strings"
	"sync"
)

var (
	sprintfRegex    = regexp.MustCompile(`%{(.*?)}`)
	logstashSprintf = regexp.MustCompile(`(?:\[.*?\])+`)
)

type stringInterpreter struct {
	lock      sync.Mutex
	templates map[string]templateNode
}

func newStringInterpreter() *stringInterpreter {
	return &stringInterpreter{
		templates: make(map[string]templateNode),
	}
}

func (i *stringInterpreter) evaluate(e *Event, format string) string {
	i.lock.Lock()
	template, exists := i.templates[format]
	if exists {
		i.lock.Unlock()
		return template.evaluate(e)
	}

	template = i.compile(format)
	i.templates[format] = template
	i.lock.Unlock()

	return template.evaluate(e)
}

func (i *stringInterpreter) compile(format string) templateNode {
	if !strings.Contains(format, "%") {
		return newStaticNode(format)
	}

	t := newTemplate()
	indicies := sprintfRegex.FindAllStringIndex(format, -1)

	pos := 0
	for _, loc := range indicies {
		if loc[0] != pos {
			t.add(newStaticNode(format[pos:loc[0]]))
		}

		pattern := format[loc[0]+2 : loc[1]-1]
		if pattern == "+%s" {
			t.add(newEpochNode())
			pos = loc[1]
			continue
		}
		if pattern[0] == '+' {
			t.add(newDateNode(pattern[1:]))
			pos = loc[1]
			continue
		}

		// Don't do a regex match unless we need to
		if pattern[0] == '[' && logstashSprintf.MatchString(pattern) {
			pattern = strings.Replace(pattern, "][", ".", -1)
			pattern = pattern[1 : len(pattern)-1]
		}
		t.add(newKeyNode(pattern))
		pos = loc[1]
	}

	return t
}
