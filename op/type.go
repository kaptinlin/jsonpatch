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
func (tp *TypeOperation) Op() internal.OpType {
	return internal.OpTypeType
}

// Code returns the operation code.
func (tp *TypeOperation) Code() int {
	return internal.OpTypeCode
}

// typeMatches checks if the actual type matches the expected type,
// with a special case where "number" also matches "integer".
func (tp *TypeOperation) typeMatches(actualType string) bool {
	if actualType == tp.TypeValue {
		return true
	}
	return tp.TypeValue == "number" && actualType == "integer"
}

// Test evaluates the type predicate condition.
func (tp *TypeOperation) Test(doc any) (bool, error) {
	val, err := value(doc, tp.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}

	return tp.typeMatches(getTypeName(val)), nil
}

// Apply applies the type operation to the document.
func (tp *TypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, err := value(doc, tp.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	actualType := getTypeName(val)
	if !tp.typeMatches(actualType) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected type %s, got %s", ErrTypeMismatch, tp.TypeValue, actualType)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (tp *TypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpTypeType),
		Path:  formatPath(tp.Path()),
		Value: tp.TypeValue,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tp *TypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTypeCode, tp.Path(), tp.TypeValue}, nil
}

// Validate validates the type operation.
func (tp *TypeOperation) Validate() error {
	if len(tp.Path()) == 0 {
		return ErrPathEmpty
	}
	if tp.TypeValue == "" {
		return ErrInvalidType
	}
	// Validate that the type is a known valid type
	if !internal.IsValidJSONPatchType(tp.TypeValue) {
		return ErrInvalidType
	}
	return nil
}
