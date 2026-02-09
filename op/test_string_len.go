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
		BaseOp:  NewBaseOp(path),
		Length:  expectedLength,
		NotFlag: false,
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
func (o *TestStringLenOperation) Op() internal.OpType {
	return internal.OpTestStringLenType
}

// Code returns the operation code.
func (o *TestStringLenOperation) Code() int {
	return internal.OpTestStringLenCode
}

// Apply applies the test string length operation to the document.
func (o *TestStringLenOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the value at the path
	value, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Convert value to string
	actualValue, err := toString(value)
	if err != nil {
		return internal.OpResult[any]{}, ErrNotString
	}

	length := int(o.Length)
	lengthMatches := len(actualValue) >= length
	shouldPass := lengthMatches != o.NotFlag
	if !shouldPass {
		// Test failed
		if o.NotFlag {
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
func (o *TestStringLenOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpTestStringLenType),
		Path: formatPath(o.Path()),
		Len:  int(o.Length),
	}
	if o.NotFlag {
		result.Not = o.NotFlag
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *TestStringLenOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringLenCode, o.Path(), o.Length}, nil
}

// Validate validates the test string length operation.
func (o *TestStringLenOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if o.Length < 0 {
		return ErrLengthNegative
	}
	return nil
}

// Test tests the string length condition on the document.
func (o *TestStringLenOperation) Test(doc any) (bool, error) {
	// Get the value at the path
	value, err := getValue(doc, o.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Convert value to string
	str, ok := extractString(value)
	if !ok {
		return false, nil // Return false if not string or byte slice
	}

	length := int(o.Length)
	lengthMatches := len(str) >= length
	return o.NotFlag != lengthMatches, nil
}

// Not returns whether this is a negation predicate.
func (o *TestStringLenOperation) Not() bool {
	return o.NotFlag
}

