package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TypeOperation represents a type operation that checks if a value is of a specific type.
type TypeOperation struct {
	BaseOp
	TypeValue string `json:"value"` // Expected type name
}

// NewType creates a new type operation.
func NewType(path []string, expectedType string) *TypeOperation {
	return &TypeOperation{
		BaseOp:    NewBaseOp(path),
		TypeValue: expectedType,
	}
}

// Op returns the operation type.
func (op *TypeOperation) Op() internal.OpType {
	return internal.OpTypeType
}

// Code returns the operation code.
func (op *TypeOperation) Code() int {
	return internal.OpTypeCode
}

// Test evaluates the type predicate condition.
func (op *TypeOperation) Test(doc any) (bool, error) {
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

// Apply applies the type operation to the document.
func (op *TypeOperation) Apply(doc any) (internal.OpResult[any], error) {
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
func (op *TypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpTypeType),
		Path:  formatPath(op.Path()),
		Value: op.TypeValue,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *TypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTypeCode, op.Path(), op.TypeValue}, nil
}

// Validate validates the type operation.
func (op *TypeOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if op.TypeValue == "" {
		return ErrInvalidType
	}
	// Validate that the type is a known valid type
	if !internal.IsValidJSONPatchType(op.TypeValue) {
		return ErrInvalidType
	}
	return nil
}

