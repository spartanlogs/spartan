package filters

//go:generate go run ../cmd/patternsGen.go ../patterns grokPatterns.go
//go:generate gofmt -w grokPatterns.go

import (
	"fmt"
	"regexp"
	"strings"
)

var varPattern = regexp.MustCompile(fmt.Sprintf("%%{%s}", grokPatterns["GROK_VARIABLE"]))

func interpolatePatterns(s string) string {
	matches := varPattern.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s
	}

	for _, match := range matches {
		var r string
		if match[2] != "" {
			r = fmt.Sprintf(`(?P<%s>%s)`, match[2], grokPatterns[match[1]])
		} else {
			r = grokPatterns[match[1]]
		}

		s = strings.Replace(s, match[0], r, 1)
	}
	return interpolatePatterns(s)
}
