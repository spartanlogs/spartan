package utils

import "errors"

// ErrNotImplemented is used in placeholders to indicate a function
// or feature that isn't implemented yet.
var ErrNotImplemented = errors.New("Feature not implemented")

// StringInSlice checks haystack for the presense of needle
func StringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
