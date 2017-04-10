package filters

//go:generate go run ../cmd/patternsGen.go ../patterns grokPatterns.go
//go:generate gofmt -w grokPatterns.go

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// PatternExt is the extension to denote a grok pattern file
const PatternExt = ".p"

var (
	varInterpolatePattern = regexp.MustCompile(fmt.Sprintf("%%{%s}", grokPatterns["GROK_VARIABLE"]))
	varPattern            = regexp.MustCompile(grokPatterns["GROK_VARIABLE"])
)

func interpolatePatterns(s string) string {
	matches := varInterpolatePattern.FindAllStringSubmatch(s, -1)
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

// LoadPatterns will load pattern files from the path.
// Path may be a directory or file. If a directory,
// all files with a PatternExt extension will be loaded in.
// Subdirectories will be recursed.
func LoadPatterns(path string) error {
	return filepath.Walk(path, walkPatternDir)
}

func walkPatternDir(path string, info os.FileInfo, err error) error {
	if info.IsDir() || err != nil {
		return err
	}

	if !strings.HasSuffix(path, PatternExt) {
		return nil
	}

	return processPatternFile(path)
}

func processPatternFile(path string) error {
	path, _ = filepath.Abs(path)

	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	lineNum := 1
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || line[0] == '#' {
			lineNum++
			continue
		}

		parts := strings.SplitAfterN(line, " ", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Error in file %s on line %d", path, lineNum)
		}

		patternName := strings.TrimSpace(parts[0])
		if !varPattern.MatchString(patternName) {
			return fmt.Errorf("Invalid grok pattern name \"%s\"", patternName)
		}

		if _, exists := grokPatterns[patternName]; exists {
			return fmt.Errorf("Pattern %s already defined. File %s, line %d", patternName, path, lineNum)
		}

		patternRegex := strings.TrimSpace(parts[1])

		grokPatterns[patternName] = patternRegex
		lineNum++
	}

	return nil
}
