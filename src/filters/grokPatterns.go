package filters

import (
	"fmt"
	"regexp"
	"strings"
)

var pattern = regexp.MustCompile(`%{([a-zA-Z_]+)(?:\:(.*?))?}`)

var patterns = map[string]string{
	"MONTHDAY":   `(?:(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9])`,
	"MONTH":      `\b(?:[Jj]an(?:uary|uar)?|[Ff]eb(?:ruary|ruar)?|[Mm](?:a|ä)?r(?:ch|z)?|[Aa]pr(?:il)?|[Mm]a(?:y|i)?|[Jj]un(?:e|i)?|[Jj]ul(?:y)?|[Aa]ug(?:ust)?|[Ss]ep(?:tember)?|[Oo](?:c|k)?t(?:ober)?|[Nn]ov(?:ember)?|[Dd]e(?:c|z)(?:ember)?)\b`,
	"YEAR":       `\d{2,4}`,
	"HOUR":       `(?:2[0123]|[01]?[0-9])`,
	"MINUTE":     `(?:[0-5][0-9])`,
	"SECOND":     `(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?)`,
	"TIME":       `%{HOUR}:%{MINUTE}(?::%{SECOND})?`,
	"IP":         `(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])`,
	"POSINT":     `\b(?:[1-9][0-9]*)\b`,
	"GREEDYDATA": `.*`,
}

func interpolatePatterns(s string) string {
	matches := pattern.FindAllStringSubmatch(s, -1)
	if len(matches) == 0 {
		return s
	}

	for _, match := range matches {
		var r string
		if match[2] != "" {
			r = fmt.Sprintf(`(?P<%s>%s)`, match[2], patterns[match[1]])
		} else {
			r = patterns[match[1]]
		}

		s = strings.Replace(s, match[0], r, 1)
	}
	return interpolatePatterns(s)
}
