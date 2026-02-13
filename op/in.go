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
func (in *InOperation) Op() internal.OpType {
	return internal.OpInType
}

// Code returns the operation code.
func (in *InOperation) Code() int {
	return internal.OpInCode
}

// Test evaluates the in predicate condition.
func (in *InOperation) Test(doc any) (bool, error) {
	_, found, err := in.getValueAndCheckInArray(doc)
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}
	return found, nil
}

// Apply applies the in test operation to the document.
func (in *InOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, found, err := in.getValueAndCheckInArray(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !found {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %v is not in array %v", ErrOperationFailed, value, in.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getValueAndCheckInArray retrieves the value and checks if it's in the array
func (in *InOperation) getValueAndCheckInArray(doc any) (any, bool, error) {
	// Get target value
	val, err := getValue(doc, in.Path())
	if err != nil {
		return nil, false, ErrPathNotFound
	}

	// Check if the value is in the specified array
	for _, v := range in.Value {
		if deepEqual(val, v) {
			return val, true, nil
		}
	}

	return val, false, nil
}

// ToJSON serializes the operation to JSON format.
func (in *InOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpInType),
		Path:  formatPath(in.Path()),
		Value: in.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (in *InOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpInCode, in.Path(), in.Value}, nil
}

// Validate validates the in operation.
func (in *InOperation) Validate() error {
	if len(in.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(in.Value) == 0 {
		return ErrValuesArrayEmpty
	}
	return nil
}
