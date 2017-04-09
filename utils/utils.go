package utils

// StringInSlice checks haystack for the presense of needle
func StringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
