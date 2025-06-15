package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpTestStringOperation represents a test operation that checks if a value is a string and matches a pattern.
type OpTestStringOperation struct {
	BaseOp
	Str string `json:"str"` // Expected string value
	Pos int    `json:"pos"` // Position within string (optional)
}

// NewOpTestStringOperation creates a new OpTestStringOperation operation.
func NewOpTestStringOperation(path []string, expectedValue string) *OpTestStringOperation {
	return &OpTestStringOperation{
		BaseOp: NewBaseOp(path),
		Str:    expectedValue,
		Pos:    0, // Default position
	}
}

// NewOpTestStringOperationWithPos creates a new OpTestStringOperation operation with position.
func NewOpTestStringOperationWithPos(path []string, expectedValue string, pos int) *OpTestStringOperation {
	return &OpTestStringOperation{
		BaseOp: NewBaseOp(path),
		Str:    expectedValue,
		Pos:    pos,
	}
}

// Op returns the operation type.
func (op *OpTestStringOperation) Op() internal.OpType {
	return internal.OpTestStringType
}

// Code returns the operation code.
func (op *OpTestStringOperation) Code() int {
	return internal.OpTestStringCode
}

// Path returns the operation path.
func (op *OpTestStringOperation) Path() []string {
	return op.path
}

// Test evaluates the test string predicate condition.
func (op *OpTestStringOperation) Test(doc any) (bool, error) {
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
func (op *OpTestStringOperation) Apply(doc any) (internal.OpResult, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult{}, ErrPathNotFound
	}

	// Check if value is a string or convert byte slice to string
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return internal.OpResult{}, ErrNotString
	}

	// Check if substring matches at the specified position
	if op.Pos < 0 || op.Pos > len(str) {
		return internal.OpResult{}, ErrPositionOutOfStringRange
	}

	endPos := op.Pos + len(op.Str)
	if endPos > len(str) {
		return internal.OpResult{}, ErrSubstringTooLong
	}

	substring := str[op.Pos:endPos]
	if substring != op.Str {
		return internal.OpResult{}, fmt.Errorf("%w at position %d", ErrSubstringMismatch, op.Pos)
	}

	return internal.OpResult{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *OpTestStringOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":   string(internal.OpTestStringType),
		"path": formatPath(op.Path()),
		"str":  op.Str,
		"pos":  op.Pos,
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpTestStringOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestStringCode, op.Path(), op.Str}, nil
}

// Validate validates the test string operation.
func (op *OpTestStringOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
