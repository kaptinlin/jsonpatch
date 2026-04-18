package jsonpatch

import "regexp"

// CreateMatcherDefault creates a regex matcher from a pattern and case sensitivity flag.
// Returns a matcher that always returns false if pattern compilation fails.
// This aligns with json-joy's createMatcherDefault behavior.
func CreateMatcherDefault(pattern string, ignoreCase bool) RegexMatcher {
	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		// Return a matcher that always returns false if compilation fails
		return func(_ string) bool { return false }
	}

	return regex.MatchString
}
