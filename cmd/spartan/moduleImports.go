package main

import (
	// Imports
	_ "github.com/spartanlogs/spartan-input-file"
	_ "github.com/spartanlogs/spartan-input-generator"

	// Filters
	_ "github.com/spartanlogs/spartan-filter-date"
	_ "github.com/spartanlogs/spartan-filter-grok"
	_ "github.com/spartanlogs/spartan-filter-metrics"
	_ "github.com/spartanlogs/spartan-filter-remove-field"

	// Outputs
	_ "github.com/spartanlogs/spartan-output-stdout"
)
