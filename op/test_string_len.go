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
func (op *TestStringLenOperation) Op() internal.OpType {
	return internal.OpTestStringLenType
}

// Code returns the operation code.
func (op *TestStringLenOperation) Code() int {
	return internal.OpTestStringLenCode
}

// Apply applies the test string length operation to the document.
func (op *TestStringLenOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get the value at the path
	value, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Convert value to string
	actualValue, err := toString(value)
	if err != nil {
		return internal.OpResult[any]{}, ErrNotString
	}

	// High-performance type conversion (single, boundary conversion)
	length := int(op.Length) // Already validated as safe integer
	// Check if the string length matches (>= comparison like TypeScript version)
	// Use XOR pattern: NotFlag XOR condition - if they're different, the test passes
	lengthMatches := len(actualValue) >= length
	shouldPass := lengthMatches != op.NotFlag
	if !shouldPass {
		// Test failed
		if op.NotFlag {
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
func (op *TestStringLenOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpTestStringLenType),
		Path: formatPath(op.Path()),
		Len:  int(op.Length),
	}
	if op.NotFlag {
		result.Not = op.NotFlag
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *TestStringLenOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringLenCode, op.Path(), op.Length}, nil
}

// Validate validates the test string length operation.
func (op *TestStringLenOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if op.Length < 0 {
		return ErrLengthNegative
	}
	return nil
}

// Test tests the string length condition on the document.
func (op *TestStringLenOperation) Test(doc any) (bool, error) {
	// Get the value at the path
	value, err := getValue(doc, op.Path())
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

	// High-performance type conversion (single, boundary conversion)
	length := int(op.Length) // Already validated as safe integer
	// Check if the string length matches (>= comparison like TypeScript version)
	lengthMatches := len(str) >= length
	// Use XOR pattern: NotFlag XOR condition - if they're different, the test passes
	return op.NotFlag != lengthMatches, nil
}

// Not returns whether this is a negation predicate.
func (op *TestStringLenOperation) Not() bool {
	return op.NotFlag
}

