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
func (l *LessOperation) Op() internal.OpType {
	return internal.OpLessType
}

// Code returns the operation code.
func (l *LessOperation) Code() int {
	return internal.OpLessCode
}

// Test evaluates the less predicate condition.
func (l *LessOperation) Test(doc any) (bool, error) {
	_, actualValue, err := numericValue(doc, l.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}
	return actualValue < l.Value, nil
}

// Apply applies the less test operation to the document.
func (l *LessOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, err := numericValue(doc, l.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if actualValue >= l.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not less than %f", ErrComparisonFailed, actualValue, l.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (l *LessOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpLessType),
		Path:  formatPath(l.Path()),
		Value: floatToJSONValue(l.Value),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (l *LessOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpLessCode, l.Path(), l.Value}, nil
}

// Validate validates the less operation.
func (l *LessOperation) Validate() error {
	if len(l.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
