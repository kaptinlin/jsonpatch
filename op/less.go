package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// LessOperation represents a test operation that checks if a numeric value is less than a specified value.
type LessOperation struct {
	BaseOp
	Value float64 `json:"value"` // Value to compare against
}

// NewLess creates a new less operation.
func NewLess(path []string, value float64) *LessOperation {
	return &LessOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// Op returns the operation type.
func (o *LessOperation) Op() internal.OpType {
	return internal.OpLessType
}

// Code returns the operation code.
func (o *LessOperation) Code() int {
	return internal.OpLessCode
}

// Test evaluates the less predicate condition.
func (o *LessOperation) Test(doc any) (bool, error) {
	_, actualValue, err := getNumericValue(doc, o.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}
	return actualValue < o.Value, nil
}

// Apply applies the less test operation to the document.
func (o *LessOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, err := getNumericValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if actualValue >= o.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not less than %f", ErrComparisonFailed, actualValue, o.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *LessOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpLessType),
		Path:  formatPath(o.Path()),
		Value: floatToJSONValue(o.Value),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *LessOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpLessCode, o.Path(), o.Value}, nil
}

// Validate validates the less operation.
func (o *LessOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
