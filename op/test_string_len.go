package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringLenOperation represents a test operation that checks if a string value has a specific length.
type TestStringLenOperation struct {
	BaseOp
	Length  float64 `json:"len"` // Expected string length
	NotFlag bool    `json:"not"` // Whether to negate the result
}

// NewTestStringLen creates a new test string length operation.
func NewTestStringLen(path []string, expectedLength float64) *TestStringLenOperation {
	return &TestStringLenOperation{
		BaseOp: NewBaseOp(path),
		Length: expectedLength,
	}
}

// NewTestStringLenWithNot creates a new test string length operation with not flag.
func NewTestStringLenWithNot(path []string, expectedLength float64, not bool) *TestStringLenOperation {
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

// Code returns the operation code.
func (tl *TestStringLenOperation) Code() int {
	return internal.OpTestStringLenCode
}

// Not returns whether this is a negation predicate.
func (tl *TestStringLenOperation) Not() bool {
	return tl.NotFlag
}

// Test tests the string length condition on the document.
func (tl *TestStringLenOperation) Test(doc any) (bool, error) {
	value, err := getValue(doc, tl.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}

	str, ok := extractString(value)
	if !ok {
		return false, nil
	}

	lengthMatches := len(str) >= int(tl.Length)
	return tl.NotFlag != lengthMatches, nil
}

// Apply applies the test string length operation to the document.
func (tl *TestStringLenOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the value at the path
	value, err := getValue(doc, tl.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Convert value to string
	actualValue, err := toString(value)
	if err != nil {
		return internal.OpResult[any]{}, ErrNotString
	}

	length := int(tl.Length)
	lengthMatches := len(actualValue) >= length
	shouldPass := lengthMatches != tl.NotFlag
	if !shouldPass {
		// Test failed
		if tl.NotFlag {
			// When Not is true and test fails, it means the length DID match when we expected it not to
			return internal.OpResult[any]{}, fmt.Errorf("%w: string length %d matched condition (>= %d) when NOT expected", ErrStringLengthMismatch, len(actualValue), length)
		}
		// When Not is false and test fails, it means the length didn't match when we expected it to
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected length >= %d, got %d", ErrStringLengthMismatch, length, len(actualValue))
	}

	// Test operations don't modify the document
	return internal.OpResult[any]{
		Doc: doc,
		Old: value,
	}, nil
}

// ToJSON serializes the operation to JSON format.
func (tl *TestStringLenOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpTestStringLenType),
		Path: formatPath(tl.Path()),
		Len:  int(tl.Length),
		Not:  tl.NotFlag,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tl *TestStringLenOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringLenCode, tl.Path(), tl.Length}, nil
}

// Validate validates the test string length operation.
func (tl *TestStringLenOperation) Validate() error {
	if len(tl.Path()) == 0 {
		return ErrPathEmpty
	}
	if tl.Length < 0 {
		return ErrLengthNegative
	}
	return nil
}
