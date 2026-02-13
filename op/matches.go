package op

import (
	"fmt"
	"regexp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MatchesOperation represents a "matches" predicate operation that checks if a string matches a regex pattern.
type MatchesOperation struct {
	BaseOp
	Pattern    string                `json:"value"`       // The regex pattern string
	IgnoreCase bool                  `json:"ignore_case"` // Case insensitive flag
	matcher    internal.RegexMatcher // Compiled regex matcher function
}

// createMatcherDefault is the default regex matcher factory using Go's regexp package.
// It creates a RegexMatcher function from a pattern and ignoreCase flag.
// If the pattern is invalid, returns a matcher that always returns false.
// This aligns with json-joy's createMatcherDefault behavior.
func createMatcherDefault(pattern string, ignoreCase bool) internal.RegexMatcher {
	var regexPattern string
	if ignoreCase {
		regexPattern = "(?i)" + pattern
	} else {
		regexPattern = pattern
	}

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		// Return a matcher that always returns false if compilation fails
		return func(_ string) bool { return false }
	}

	return re.MatchString
}

// NewMatches creates a new matches operation.
// If createMatcher is nil, uses the default Go regexp implementation.
// This aligns with json-joy's OpMatches constructor pattern.
func NewMatches(path []string, pattern string, ignoreCase bool, createMatcher internal.CreateRegexMatcher) *MatchesOperation {
	if createMatcher == nil {
		createMatcher = createMatcherDefault
	}

	return &MatchesOperation{
		BaseOp:     NewBaseOp(path),
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
		matcher:    createMatcher(pattern, ignoreCase),
	}
}

// Op returns the operation type.
func (ma *MatchesOperation) Op() internal.OpType {
	return internal.OpMatchesType
}

// Code returns the operation code.
func (ma *MatchesOperation) Code() int {
	return internal.OpMatchesCode
}

// Test evaluates the matches predicate condition.
func (ma *MatchesOperation) Test(doc any) (bool, error) {
	_, str, err := getAndValidateString(doc, ma.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return ma.matcher(str), nil
}

// Apply applies the matches operation.
func (ma *MatchesOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, str, err := getAndValidateString(doc, ma.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !ma.matcher(str) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string '%s' does not match pattern", ErrStringMismatch, str)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// ToJSON converts the operation to JSON representation.
func (ma *MatchesOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpMatchesType),
		Path:       formatPath(ma.Path()),
		Value:      ma.Pattern,
		IgnoreCase: ma.IgnoreCase,
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (ma *MatchesOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMatchesCode, ma.Path(), ma.Pattern, ma.IgnoreCase}, nil
}

// Validate validates the matches operation.
func (ma *MatchesOperation) Validate() error {
	if len(ma.Path()) == 0 {
		return ErrPathEmpty
	}
	if ma.Pattern == "" {
		return ErrPatternEmpty
	}
	return nil
}
