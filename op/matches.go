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
	NotFlag    bool                  `json:"not"`         // Whether to negate the result
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

// NewOpMatchesOperation creates a new matches operation.
// If createMatcher is nil, uses the default Go regexp implementation.
// This aligns with json-joy's OpMatches constructor pattern.
func NewOpMatchesOperation(path []string, pattern string, ignoreCase bool, notFlag bool, createMatcher internal.CreateRegexMatcher) *MatchesOperation {
	if createMatcher == nil {
		createMatcher = createMatcherDefault
	}

	return &MatchesOperation{
		BaseOp:     NewBaseOp(path),
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
		NotFlag:    notFlag,
		matcher:    createMatcher(pattern, ignoreCase),
	}
}

// Op returns the operation type.
func (o *MatchesOperation) Op() internal.OpType {
	return internal.OpMatchesType
}

// Code returns the operation code.
func (o *MatchesOperation) Code() int {
	return internal.OpMatchesCode
}

// Test evaluates the matches predicate condition.
func (o *MatchesOperation) Test(doc interface{}) (bool, error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Convert to string
	str, err := toString(val)
	if err != nil {
		// For JSON Patch test operations, wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	matches := o.matcher(str)

	// Apply negation if needed
	if o.NotFlag {
		matches = !matches
	}

	return matches, nil
}

// Apply applies the matches operation.
func (o *MatchesOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Convert to string
	str, err := toString(val)
	if err != nil {
		return internal.OpResult[any]{}, ErrNotString
	}

	matches := o.matcher(str)

	// Apply negation if needed
	if o.NotFlag {
		matches = !matches
	}

	if !matches {
		if o.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: string '%s' matches pattern", ErrStringMismatch, str)
		}
		return internal.OpResult[any]{}, fmt.Errorf("%w: string '%s' does not match pattern", ErrStringMismatch, str)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// ToJSON converts the operation to JSON representation.
func (o *MatchesOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:         string(internal.OpMatchesType),
		Path:       formatPath(o.Path()),
		Value:      o.Pattern,
		IgnoreCase: o.IgnoreCase,
		Not:        o.NotFlag,
	}

	return result, nil
}

// ToCompact converts the operation to compact array representation.
func (o *MatchesOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMatchesCode, o.Path(), o.Pattern, o.IgnoreCase}, nil
}

// Not returns the negation flag.
func (o *MatchesOperation) Not() bool {
	return o.NotFlag
}

// Validate validates the matches operation.
func (o *MatchesOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if o.Pattern == "" {
		return ErrPatternEmpty
	}
	return nil
}

// Path returns the path for the matches operation.
func (o *MatchesOperation) Path() []string {
	return o.path
}

// Short aliases for common use
var (
	// NewMatches creates a new matches operation
	NewMatches = NewOpMatchesOperation
)
