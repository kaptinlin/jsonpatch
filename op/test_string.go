package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringOperation represents a test operation that checks if a value is a string and matches a pattern.
type TestStringOperation struct {
	BaseOp
	Str     string `json:"str"`           // Expected string value
	Pos     int    `json:"pos"`           // Position within string
	NotFlag bool   `json:"not,omitempty"` // Whether to negate the result
}

// NewTestString creates a new test string operation.
func NewTestString(path []string, str string, pos int, not bool) *TestStringOperation {
	return &TestStringOperation{
		BaseOp:  NewBaseOp(path),
		Str:     str,
		Pos:     pos,
		NotFlag: not,
	}
}

// Op returns the operation type.
func (ts *TestStringOperation) Op() internal.OpType {
	return internal.OpTestStringType
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
	substring, err := codePointSubstring(str, ts.Pos, codePointCount(ts.Str))
	if err != nil {
		return false, err
	}
	return ts.NotFlag != (substring == ts.Str), nil
}

// Apply applies the test string operation to the document.
func (ts *TestStringOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, err := value(doc, ts.Path())
	if err != nil {
		return internal.OpResult[any]{}, predicateFailure(ErrPathNotFound)
	}

	str, ok := extractString(val)
	if !ok {
		return internal.OpResult[any]{}, predicateFailure(ErrNotString)
	}

	pos := ts.Pos
	substring, err := codePointSubstring(str, pos, codePointCount(ts.Str))
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	shouldPass := (substring == ts.Str) != ts.NotFlag
	if !shouldPass {
		if ts.NotFlag {
			// When Not is true and test fails, it means the string DID match when we expected it not to
			return internal.OpResult[any]{}, predicateFailure(fmt.Errorf("%w: string matched %q at position %d when NOT expected", ErrSubstringMismatch, ts.Str, pos))
		}
		// When Not is false and test fails, it means the string didn't match when we expected it to
		return internal.OpResult[any]{}, predicateFailure(fmt.Errorf("%w: expected %q at position %d, got %q", ErrSubstringMismatch, ts.Str, pos, substring))
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// Validate validates the test string operation.
func (ts *TestStringOperation) Validate() error {
	return nil
}
