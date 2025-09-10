package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringLenOperation represents a test operation that checks if a string value has a specific length.
type TestStringLenOperation struct {
	BaseOp
	Length float64 `json:"len"` // Expected string length
	Not    bool    `json:"not"` // Whether to negate the result
}

// NewOpTestStringLenOperation creates a new OpTestStringLenOperation operation.
func NewOpTestStringLenOperation(path []string, expectedLength float64) *TestStringLenOperation {
	return &TestStringLenOperation{
		BaseOp: NewBaseOp(path),
		Length: expectedLength,
		Not:    false,
	}
}

// NewOpTestStringLenOperationWithNot creates a new OpTestStringLenOperation operation with not flag.
func NewOpTestStringLenOperationWithNot(path []string, expectedLength float64, not bool) *TestStringLenOperation {
	return &TestStringLenOperation{
		BaseOp: NewBaseOp(path),
		Length: expectedLength,
		Not:    not,
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

// Path returns the operation path.
func (op *TestStringLenOperation) Path() []string {
	return op.path
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
	lengthMatches := len(actualValue) >= length
	if op.Not {
		lengthMatches = !lengthMatches
	}

	if !lengthMatches {
		if op.Not {
			return internal.OpResult[any]{}, fmt.Errorf("%w: expected length NOT >= %d, but got %d", ErrStringLengthMismatch, length, len(actualValue))
		}
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
		"op":   string(internal.OpTestStringLenType),
		"path": formatPath(op.Path()),
		"len":  op.Length,
	}
	if op.Not {
		result["not"] = op.Not
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

// Short aliases for common use
var (
	// NewTestStringLen creates a new test string length operation
	NewTestStringLen = NewOpTestStringLenOperation
	// NewTestStringLenWithNot creates a new test string length operation with not flag
	NewTestStringLenWithNot = NewOpTestStringLenOperationWithNot
)
