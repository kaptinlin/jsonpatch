package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringOperation represents a test operation that checks if a value is a string and matches a pattern.
type TestStringOperation struct {
	BaseOp
	Str        string `json:"str"`                   // Expected string value
	Pos        int    `json:"pos"`                   // Position within string
	NotFlag    bool   `json:"not,omitempty"`         // Whether to negate the result
	IgnoreCase bool   `json:"ignore_case,omitempty"` // Whether to ignore case
}

// NewTestString creates a new test string operation.
func NewTestString(path []string, str string, pos float64, not bool, ignoreCase bool) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:     NewBaseOp(path),
		Str:        str,
		Pos:        int(pos),
		NotFlag:    not,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (ts *TestStringOperation) Op() internal.OpType {
	return internal.OpTestStringType
}

// Code returns the operation code.
func (ts *TestStringOperation) Code() int {
	return internal.OpTestStringCode
}

// Not returns whether this operation is a negation predicate.
func (ts *TestStringOperation) Not() bool {
	return ts.NotFlag
}

// Test evaluates the test string predicate condition.
func (ts *TestStringOperation) Test(doc any) (bool, error) {
	val, err := value(doc, ts.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}

	str, ok := extractString(val)
	if !ok {
		return false, nil
	}

	// Match json-joy behavior: clamp positions to string bounds
	length := len(str)
	start := min(ts.Pos, length)
	end := min(ts.Pos+len(ts.Str), length)

	substring := str[start:end]
	var match bool
	if ts.IgnoreCase {
		match = strings.EqualFold(substring, ts.Str)
	} else {
		match = substring == ts.Str
	}
	return ts.NotFlag != match, nil
}

// Apply applies the test string operation to the document.
func (ts *TestStringOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, err := value(doc, ts.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	str, ok := extractString(val)
	if !ok {
		return internal.OpResult[any]{}, ErrNotString
	}

	pos := ts.Pos
	if pos < 0 || pos > len(str) {
		return internal.OpResult[any]{}, ErrPositionOutOfStringRange
	}

	endPos := pos + len(ts.Str)
	if endPos > len(str) {
		return internal.OpResult[any]{}, ErrSubstringTooLong
	}

	substring := str[pos:endPos]
	var matches bool
	if ts.IgnoreCase {
		matches = strings.EqualFold(substring, ts.Str)
	} else {
		matches = substring == ts.Str
	}

	// Apply negation logic using XOR (same as Test method)
	shouldPass := matches != ts.NotFlag
	if !shouldPass {
		if ts.NotFlag {
			// When Not is true and test fails, it means the string DID match when we expected it not to
			return internal.OpResult[any]{}, fmt.Errorf("%w: string matched %q at position %d when NOT expected", ErrSubstringMismatch, ts.Str, pos)
		}
		// When Not is false and test fails, it means the string didn't match when we expected it to
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected %q at position %d, got %q", ErrSubstringMismatch, ts.Str, pos, substring)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (ts *TestStringOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpTestStringType),
		Path:       formatPath(ts.Path()),
		Str:        ts.Str,
		Pos:        ts.Pos,
		Not:        ts.NotFlag,
		IgnoreCase: ts.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ts *TestStringOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringCode, ts.Path(), ts.Str}, nil
}

// Validate validates the test string operation.
func (ts *TestStringOperation) Validate() error {
	if len(ts.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
