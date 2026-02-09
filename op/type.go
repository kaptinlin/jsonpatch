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
func (o *TypeOperation) Op() internal.OpType {
	return internal.OpTypeType
}

// Code returns the operation code.
func (o *TypeOperation) Code() int {
	return internal.OpTypeCode
}

// Test evaluates the type predicate condition.
func (o *TypeOperation) Test(doc any) (bool, error) {
	val, err := getValue(doc, o.Path())
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Get the actual type of the value
	actualType := getTypeName(val)

	// Check if the type matches
	return actualType == o.TypeValue, nil
}

// Apply applies the type operation to the document.
func (o *TypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Get the actual type of the value
	actualType := getTypeName(val)

	// Check if the type matches
	typeMatches := actualType == o.TypeValue

	// Special case: if expected type is "number" and actual is "integer", it should match
	if !typeMatches && o.TypeValue == "number" && actualType == "integer" {
		typeMatches = true
	}

	if !typeMatches {
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected type %s, got %s", ErrTypeMismatch, o.TypeValue, actualType)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *TypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpTypeType),
		Path:  formatPath(o.Path()),
		Value: o.TypeValue,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *TypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTypeCode, o.Path(), o.TypeValue}, nil
}

// Validate validates the type operation.
func (o *TypeOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if o.TypeValue == "" {
		return ErrInvalidType
	}
	// Validate that the type is a known valid type
	if !internal.IsValidJSONPatchType(o.TypeValue) {
		return ErrInvalidType
	}
	return nil
}

