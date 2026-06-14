package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MoreOperation represents a "more" predicate operation that checks if a value is greater than a specified number.
type MoreOperation struct {
	BaseOp
	Value float64 `json:"value"` // The number to compare against
}

// NewMore creates a new more operation.
func NewMore(path []string, value float64) *MoreOperation {
	return &MoreOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// Op returns the operation type.
func (mo *MoreOperation) Op() internal.OpType {
	return internal.OpMoreType
}

// Test evaluates the more predicate condition.
func (mo *MoreOperation) Test(doc any) (bool, error) {
	_, num, err := numericValue(doc, mo.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return num > mo.Value, nil
}

// Apply applies the more operation.
func (mo *MoreOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, num, err := numericValue(doc, mo.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if num <= mo.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not greater than %f", ErrComparisonFailed, num, mo.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// Validate validates the more operation.
func (mo *MoreOperation) Validate() error {
	return nil
}
