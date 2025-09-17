package op

import (
	"fmt"
	"reflect"

	"github.com/kaptinlin/jsonpatch/internal"
)

// InOperation represents a test operation that checks if a value is present in a specified array.
type InOperation struct {
	BaseOp
	Value []interface{} `json:"value"` // Array of values to check against
}

// NewOpInOperation creates a new OpInOperation operation.
func NewOpInOperation(path []string, values []interface{}) *InOperation {
	return &InOperation{
		BaseOp: NewBaseOp(path),
		Value:  values,
	}
}

// Op returns the operation type.
func (op *InOperation) Op() internal.OpType {
	return internal.OpInType
}

// Code returns the operation code.
func (op *InOperation) Code() int {
	return internal.OpInCode
}

// Path returns the operation path.
func (op *InOperation) Path() []string {
	return op.path
}

// Test evaluates the in predicate condition.
func (op *InOperation) Test(doc any) (bool, error) {
	_, found, err := op.getValueAndCheckInArray(doc)
	if err != nil {
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return found, nil
}

// Apply applies the in test operation to the document.
func (op *InOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, found, err := op.getValueAndCheckInArray(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !found {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %v is not in array %v", ErrOperationFailed, value, op.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getValueAndCheckInArray retrieves the value and checks if it's in the array
func (op *InOperation) getValueAndCheckInArray(doc any) (interface{}, bool, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return nil, false, ErrPathNotFound
	}

	// Check if the value is in the specified array
	for _, v := range op.Value {
		if reflect.DeepEqual(val, v) {
			return val, true, nil
		}
	}

	return val, false, nil
}

// ToJSON serializes the operation to JSON format.
func (op *InOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":    string(internal.OpInType),
		"path":  formatPath(op.Path()),
		"value": op.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *InOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpInCode, op.Path(), op.Value}, nil
}

// Validate validates the in operation.
func (op *InOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(op.Value) == 0 {
		return ErrValuesArrayEmpty
	}
	return nil
}

// Not returns false since this is not a NOT operation.
func (op *InOperation) Not() bool {
	return false
}

// Short aliases for common use
var (
	// NewIn creates a new in operation
	NewIn = NewOpInOperation
)
