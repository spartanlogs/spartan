package main

import (
	// Imports
	_ "github.com/spartanlogs/spartan-file-input"
	_ "github.com/spartanlogs/spartan-generator-input"

	// Filters
	_ "github.com/spartanlogs/spartan-date-filter"
	_ "github.com/spartanlogs/spartan-grok-filter"
	_ "github.com/spartanlogs/spartan-metrics-filter"
	_ "github.com/spartanlogs/spartan-remove-field-filter"

	// Outputs
	_ "github.com/spartanlogs/spartan-stdout-output"
)
