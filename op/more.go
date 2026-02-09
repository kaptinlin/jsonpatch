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

// Code returns the operation code.
func (mo *MoreOperation) Code() int {
	return internal.OpMoreCode
}

// Test evaluates the more predicate condition.
func (mo *MoreOperation) Test(doc any) (bool, error) {
	_, num, err := getNumericValue(doc, mo.Path())
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	return num > mo.Value, nil
}

// Apply applies the more operation.
func (mo *MoreOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, num, err := getNumericValue(doc, mo.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if num <= mo.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not greater than %f", ErrComparisonFailed, num, mo.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// ToJSON converts the operation to JSON representation.
func (mo *MoreOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpMoreType),
		Path:  formatPath(mo.Path()),
		Value: floatToJSONValue(mo.Value),
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (mo *MoreOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoreCode, mo.Path(), mo.Value}, nil
}

// Validate validates the more operation.
func (mo *MoreOperation) Validate() error {
	if len(mo.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
