package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringLenOperation represents a test operation that checks if a string value has a specific length.
type TestStringLenOperation struct {
	BaseOp
	Length  int  `json:"len"` // Expected string length
	NotFlag bool `json:"not"` // Whether to negate the result
}

// NewTestStringLen creates a new test string length operation.
func NewTestStringLen(path []string, expectedLength int) *TestStringLenOperation {
	return &TestStringLenOperation{
		BaseOp: NewBaseOp(path),
		Length: expectedLength,
	}
}

// NewTestStringLenWithNot creates a new test string length operation with not flag.
func NewTestStringLenWithNot(path []string, expectedLength int, not bool) *TestStringLenOperation {
	return &TestStringLenOperation{
		BaseOp:  NewBaseOp(path),
		Length:  expectedLength,
		NotFlag: not,
	}
}

// Op returns the operation type.
func (tl *TestStringLenOperation) Op() internal.OpType {
	return internal.OpTestStringLenType
}

// Not returns whether this is a negation predicate.
func (tl *TestStringLenOperation) Not() bool {
	return tl.NotFlag
}

// Test tests the string length condition on the document.
func (tl *TestStringLenOperation) Test(doc any) (bool, error) {
	value, err := value(doc, tl.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}

	str, ok := extractString(value)
	if !ok {
		return false, nil
	}

	lengthMatches := codePointCount(str) >= tl.Length
	return tl.NotFlag != lengthMatches, nil
}

// Apply applies the test string length operation to the document.
func (tl *TestStringLenOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the value at the path
	value, err := value(doc, tl.Path())
	if err != nil {
		return internal.OpResult[any]{}, predicateFailure(ErrPathNotFound)
	}

	// Convert value to string
	actualValue, err := toString(value)
	if err != nil {
		return internal.OpResult[any]{}, predicateFailure(ErrNotString)
	}

	length := tl.Length
	actualLength := codePointCount(actualValue)
	lengthMatches := actualLength >= length
	shouldPass := lengthMatches != tl.NotFlag
	if !shouldPass {
		// Test failed
		if tl.NotFlag {
			// When Not is true and test fails, it means the length DID match when we expected it not to
			return internal.OpResult[any]{}, predicateFailure(fmt.Errorf("%w: string length %d matched condition (>= %d) when NOT expected", ErrStringLengthMismatch, actualLength, length))
		}
		// When Not is false and test fails, it means the length didn't match when we expected it to
		return internal.OpResult[any]{}, predicateFailure(fmt.Errorf("%w: expected length >= %d, got %d", ErrStringLengthMismatch, length, actualLength))
	}

	// Test operations don't modify the document
	return internal.OpResult[any]{
		Doc: doc,
		Old: value,
	}, nil
}

// Validate validates the test string length operation.
func (tl *TestStringLenOperation) Validate() error {
	if tl.Length < 0 {
		return ErrLengthNegative
	}
	return nil
}
