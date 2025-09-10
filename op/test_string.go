package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestStringOperation represents a test operation that checks if a value is a string and matches a pattern.
type TestStringOperation struct {
	BaseOp
	Str string  `json:"str"` // Expected string value
	Pos float64 `json:"pos"` // Position within string (optional)
}

// NewOpTestStringOperation creates a new OpTestStringOperation operation.
func NewOpTestStringOperation(path []string, expectedValue string) *TestStringOperation {
	return &TestStringOperation{
		BaseOp: NewBaseOp(path),
		Str:    expectedValue,
		Pos:    0, // Default position
	}
}

// NewOpTestStringOperationWithPos creates a new OpTestStringOperation operation with position.
func NewOpTestStringOperationWithPos(path []string, expectedValue string, pos float64) *TestStringOperation {
	return &TestStringOperation{
		BaseOp: NewBaseOp(path),
		Str:    expectedValue,
		Pos:    pos,
	}
}

// Op returns the operation type.
func (op *TestStringOperation) Op() internal.OpType {
	return internal.OpTestStringType
}

// Code returns the operation code.
func (op *TestStringOperation) Code() int {
	return internal.OpTestStringCode
}

// Path returns the operation path.
func (op *TestStringOperation) Path() []string {
	return op.path
}

// Test evaluates the test string predicate condition.
func (op *TestStringOperation) Test(doc any) (bool, error) {
	val, err := getValue(doc, op.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Convert to string or from byte slice
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return false, nil // Return false if not string or byte slice
	}

	return str == op.Str, nil
}

// Apply applies the test string operation to the document.
func (op *TestStringOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if value is a string or convert byte slice to string
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return internal.OpResult[any]{}, ErrNotString
	}

	// High-performance type conversion (single, boundary conversion)
	pos := int(op.Pos) // Already validated as safe integer
	// Check if substring matches at the specified position
	if pos < 0 || pos > len(str) {
		return internal.OpResult[any]{}, ErrPositionOutOfStringRange
	}

	endPos := pos + len(op.Str)
	if endPos > len(str) {
		return internal.OpResult[any]{}, ErrSubstringTooLong
	}

	substring := str[pos:endPos]
	if substring != op.Str {
		return internal.OpResult[any]{}, fmt.Errorf("%w at position %d", ErrSubstringMismatch, pos)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *TestStringOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":   string(internal.OpTestStringType),
		"path": formatPath(op.Path()),
		"str":  op.Str,
		"pos":  op.Pos,
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *TestStringOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringCode, op.Path(), op.Str}, nil
}

// Validate validates the test string operation.
func (op *TestStringOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Short aliases for common use
var (
	// NewTestString creates a new test string operation
	NewTestString = NewOpTestStringOperation
	// NewTestStringWithPos creates a new test string operation with position
	NewTestStringWithPos = NewOpTestStringOperationWithPos
)
