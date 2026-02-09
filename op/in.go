package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// InOperation represents a test operation that checks if a value is present in a specified array.
type InOperation struct {
	BaseOp
	Value []any `json:"value"` // Array of values to check against
}

// NewIn creates a new in operation.
func NewIn(path []string, values []any) *InOperation {
	return &InOperation{
		BaseOp: NewBaseOp(path),
		Value:  values,
	}
}

// Op returns the operation type.
func (o *InOperation) Op() internal.OpType {
	return internal.OpInType
}

// Code returns the operation code.
func (o *InOperation) Code() int {
	return internal.OpInCode
}

// Test evaluates the in predicate condition.
func (o *InOperation) Test(doc any) (bool, error) {
	_, found, err := o.getValueAndCheckInArray(doc)
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return found, nil
}

// Apply applies the in test operation to the document.
func (o *InOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, found, err := o.getValueAndCheckInArray(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !found {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %v is not in array %v", ErrOperationFailed, value, o.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getValueAndCheckInArray retrieves the value and checks if it's in the array
func (o *InOperation) getValueAndCheckInArray(doc any) (any, bool, error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		return nil, false, ErrPathNotFound
	}

	// Check if the value is in the specified array
	for _, v := range o.Value {
		if deepEqual(val, v) {
			return val, true, nil
		}
	}

	return val, false, nil
}

// ToJSON serializes the operation to JSON format.
func (o *InOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpInType),
		Path:  formatPath(o.Path()),
		Value: o.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *InOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpInCode, o.Path(), o.Value}, nil
}

// Validate validates the in operation.
func (o *InOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(o.Value) == 0 {
		return ErrValuesArrayEmpty
	}
	return nil
}
