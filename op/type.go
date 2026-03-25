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

// Test evaluates the type predicate condition.
// Matches json-joy behavior: uses typeof semantics, with special case
// only when expected type is "integer" and value is a whole number.
func (tp *TypeOperation) Test(doc any) (bool, error) {
	val, err := value(doc, tp.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}

	actualType := getTypeName(val)
	if actualType == tp.TypeValue {
		return true, nil
	}
	// Special case: "integer" matches whole numbers (json-joy behavior)
	if tp.TypeValue == "integer" && isWholeNumber(val) {
		return true, nil
	}
	return false, nil
}

// Apply applies the type operation to the document.
func (tp *TypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, err := value(doc, tp.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	actualType := getTypeName(val)
	matched := actualType == tp.TypeValue ||
		(tp.TypeValue == "integer" && isWholeNumber(val))
	if !matched {
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
