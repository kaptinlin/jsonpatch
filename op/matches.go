package op

import (
	"fmt"
	"regexp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MatchesOperation represents a "matches" predicate operation that checks if a string matches a regex pattern.
type MatchesOperation struct {
	BaseOp
	Pattern    string         // The regex pattern string
	IgnoreCase bool           // Case insensitive flag
	matcher    *regexp.Regexp // Compiled regex matcher
}

// NewOpMatchesOperation creates a new matches operation.
func NewOpMatchesOperation(path []string, pattern string, ignoreCase bool) (*MatchesOperation, error) {
	// Compile the regex pattern
	var regexPattern string
	if ignoreCase {
		regexPattern = "(?i)" + pattern
	} else {
		regexPattern = pattern
	}

	matcher, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRegexPattern, err)
	}

	return &MatchesOperation{
		BaseOp:     NewBaseOp(path),
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
		matcher:    matcher,
	}, nil
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

	return o.matcher.MatchString(str), nil
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

	if !o.matcher.MatchString(str) {
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
	}

	return result, nil
}

// ToCompact converts the operation to compact array representation.
func (o *MatchesOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMatchesCode, o.Path(), o.Pattern, o.IgnoreCase}, nil
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
