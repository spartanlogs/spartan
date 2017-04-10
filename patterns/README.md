# Grok Patterns

This directory contains all the core patterns in Sparta. These patterns are used by the Grok filter
with the `%{PATTERN_NAME}` syntax. Only files with a `.p` extension are read. Patterns in this directory
are compiled into the binary with a Go script using `make generate`. Directories are recursed. Patterns
are checked for valid pattern names. Regex validation happens at runtime.

The base set of patterns was taken from the Logstash project [here](https://github.com/logstash-plugins/logstash-patterns-core/tree/master/patterns).
Many patterns had to be modified to work with Go's regex engine which doesn't do backtracking
and thus doesn't allow several syntaxes that other engines do.

## Making changes

When a change is made to any pattern or file, run `make generate` to build the Go source file
with the new patterns to be compiled in.

## File format

Each pattern is a single line consistaning of a pattern name, followed by a space, and then the
rest of the line is taken as the regex pattern. Pattern names may contain upper or lowercase English
letters, numbers, or underscores. For example:

```
POSINT \b(?:[1-9][0-9]*)\b
GREEDYDATA .*
ISO8601_TIMEZONE (?:Z|[+-]%{HOUR}(?::?%{MINUTE}))
DATE_US %{MONTHNUM}[/-]%{MONTHDAY}[/-]%{YEAR}
```

Comments are allowed by starting a line with a pound sign. Comments must be on a separate line:

```
# This is a comment
POSINT \b(?:[1-9][0-9]*)\b # This is part of the pattern
GREEDYDATA .*
```
