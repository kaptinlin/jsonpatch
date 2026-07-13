package jsonpatch

import "regexp"

// CreateMatcherDefault creates a regex matcher from a pattern and case sensitivity flag.
// It returns a regexp syntax error when pattern is invalid.
func CreateMatcherDefault(pattern string, ignoreCase bool) (RegexMatcher, error) {
	flags := ""
	if ignoreCase {
		flags = "(?i)"
	}

	regex, err := regexp.Compile(flags + pattern)
	if err != nil {
		return nil, err
	}

	return regex.MatchString, nil
}
