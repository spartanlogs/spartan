package filters

import "testing"

func TestPatternInterpolation(t *testing.T) {
	tests := []struct{ start, end string }{
		{
			start: "Hello",
			end:   "Hello",
		},
		{
			start: "%{HOUR}",
			end:   `(?:2[0123]|[01]?[0-9])`,
		},
		{
			start: "%{HOUR:theHour}",
			end:   `(?P<theHour>(?:2[0123]|[01]?[0-9]))`,
		},
		{
			start: "%{TIME}",
			end:   `(?:2[0123]|[01]?[0-9]):(?:[0-5][0-9])(?::(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?))?`,
		},
		{
			start: `^(?P<logdate>%{MONTHDAY}[-]%{MONTH}[-]%{YEAR} %{TIME}) client %{IP:clientip}#%{POSINT:clientport} \(%{GREEDYDATA:query}\): query: %{GREEDYDATA:target} IN %{GREEDYDATA:querytype} \(%{IP:dns}\)$`,
			end:   `^(?P<logdate>(?:(?:0[1-9])|(?:[12][0-9])|(?:3[01])|[1-9])[-]\b(?:[Jj]an(?:uary|uar)?|[Ff]eb(?:ruary|ruar)?|[Mm](?:a|ä)?r(?:ch|z)?|[Aa]pr(?:il)?|[Mm]a(?:y|i)?|[Jj]un(:e|i)?|[Jj]ul(?:y)?|[Aa]ug(?:ust)?|[Ss]ep(?:tember)?|[Oo](?:c|k)?t(?:ober)?|[Nn]ov(?:ember)?|[Dd]e(?:c|z)(?:ember)?)\b[-](?:\d\d){1,2} (?:2[0123]|[01]?[0-9]):(?:[0-5][0-9])(?::(?:(?:[0-5]?[0-9]|60)(?:[:.,][0-9]+)?))?) client (?P<clientip>(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9]))#(?P<clientport>\b(?:[1-9][0-9]*)\b) \((?P<query>.*)\): query: (?P<target>.*) IN (?P<querytype>.*) \((?P<dns>(?:(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|1[0-9][0-9]|[1-9]?[0-9]))\)$`,
		},
	}

	for i, test := range tests {
		got := interpolatePatterns(test.start)
		if got != test.end {
			t.Errorf("Pattern interpolation test %d. Expected '%s', got '%s'", i+1, test.end, got)
		}
	}
}
