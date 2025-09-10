package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpTypeOperation represents a type operation that checks if a value is of a specific type.
type OpTypeOperation struct {
	BaseOp
	TypeValue string `json:"value"` // Expected type name
}

// NewOpTypeOperation creates a new OpTypeOperation operation.
func NewOpTypeOperation(path []string, expectedType string) *OpTypeOperation {
	return &OpTypeOperation{
		BaseOp:    NewBaseOp(path),
		TypeValue: expectedType,
	}
}

// Op returns the operation type.
func (op *OpTypeOperation) Op() internal.OpType {
	return internal.OpTypeType
}

// Code returns the operation code.
func (op *OpTypeOperation) Code() int {
	return internal.OpTypeCode
}

// Path returns the operation path.
func (op *OpTypeOperation) Path() []string {
	return op.path
}

// Test evaluates the type predicate condition.
func (op *OpTypeOperation) Test(doc any) (bool, error) {
	val, err := getValue(doc, op.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Get the actual type of the value
	actualType := getTypeName(val)

	// Check if the type matches
	return actualType == op.TypeValue, nil
}

// Not returns false (type operation doesn't support not modifier).
func (op *OpTypeOperation) Not() bool {
	return false
}

// Apply applies the type operation to the document.
func (op *OpTypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Get the actual type of the value
	actualType := getTypeName(val)

	// Check if the type matches
	typeMatches := actualType == op.TypeValue

	// Special case: if expected type is "number" and actual is "integer", it should match
	if !typeMatches && op.TypeValue == "number" && actualType == "integer" {
		typeMatches = true
	}

	if !typeMatches {
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected type %s, got %s", ErrTypeMismatch, op.TypeValue, actualType)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *OpTypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":    string(internal.OpTypeType),
		"path":  formatPath(op.Path()),
		"value": op.TypeValue,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpTypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTypeCode, op.Path(), op.TypeValue}, nil
}

// Validate validates the type operation.
func (op *OpTypeOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if op.TypeValue == "" {
		return ErrInvalidType
	}
	// Validate that the type is a known valid type
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
		"object":  true,
		"array":   true,
		"null":    true,
		"integer": true, // Special type that's also valid
	}
	if !validTypes[op.TypeValue] {
		return ErrInvalidType
	}
	return nil
}

// Short aliases for common use
var (
	// NewType creates a new type operation
	NewType = NewOpTypeOperation
)
