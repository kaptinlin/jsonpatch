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
	matcherErr error
}

// createMatcherDefault is the default regex matcher factory using Go's regexp package.
// It creates a RegexMatcher function from a pattern and ignoreCase flag.
func createMatcherDefault(pattern string, ignoreCase bool) (internal.RegexMatcher, error) {
	regexPattern := pattern
	if ignoreCase {
		regexPattern = "(?i)" + pattern
	}

	re, err := regexp.Compile(regexPattern)
	if err != nil {
		return nil, err
	}

	return re.MatchString, nil
}

// NewMatches creates a new matches operation.
// If createMatcher is nil, uses the default Go regexp implementation.
func NewMatches(path []string, pattern string, ignoreCase bool, createMatcher internal.CreateRegexMatcher) *MatchesOperation {
	if createMatcher == nil {
		createMatcher = createMatcherDefault
	}
	matcher, err := createMatcher(pattern, ignoreCase)
	if err != nil {
		err = fmt.Errorf("%w: %w", ErrRegexPattern, err)
	} else if matcher == nil {
		err = fmt.Errorf("%w: matcher factory returned nil", ErrRegexPattern)
	}

	return &MatchesOperation{
		BaseOp:     NewBaseOp(path),
		Pattern:    pattern,
		IgnoreCase: ignoreCase,
		matcher:    matcher,
		matcherErr: err,
	}
}

// Op returns the operation type.
func (ma *MatchesOperation) Op() internal.OpType {
	return internal.OpMatchesType
}

// Test evaluates the matches predicate condition.
func (ma *MatchesOperation) Test(doc any) (bool, error) {
	if err := ma.matcherError(); err != nil {
		return false, err
	}
	_, str, err := validateString(doc, ma.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return ma.matcher(str), nil
}

// Apply applies the matches operation.
func (ma *MatchesOperation) Apply(doc any) (internal.OpResult[any], error) {
	if err := ma.matcherError(); err != nil {
		return internal.OpResult[any]{}, err
	}
	val, str, err := validateString(doc, ma.Path())
	if err != nil {
		return internal.OpResult[any]{}, predicateFailure(err)
	}

	if !ma.matcher(str) {
		return internal.OpResult[any]{}, predicateFailure(fmt.Errorf("%w: string '%s' does not match pattern", ErrStringMismatch, str))
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// Validate validates the matches operation.
func (ma *MatchesOperation) Validate() error {
	if ma.Pattern == "" {
		return ErrPatternEmpty
	}
	return ma.matcherError()
}

func (ma *MatchesOperation) matcherError() error {
	if ma.matcherErr != nil {
		return ma.matcherErr
	}
	if ma.matcher == nil {
		return fmt.Errorf("%w: matcher is nil", ErrRegexPattern)
	}
	return nil
}
