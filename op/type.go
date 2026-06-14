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

// Test evaluates the type predicate condition.
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

// Validate validates the type operation.
func (tp *TypeOperation) Validate() error {
	if tp.TypeValue == "" {
		return ErrInvalidType
	}
	// Validate that the type is a known valid type
	if !internal.IsValidJSONPatchType(tp.TypeValue) {
		return ErrInvalidType
	}
	return nil
}
