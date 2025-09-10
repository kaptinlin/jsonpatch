package op

import (
	"fmt"
	"reflect"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpInOperation represents a test operation that checks if a value is present in a specified array.
type OpInOperation struct {
	BaseOp
	Values []interface{} `json:"values"` // Array of values to check against
}

// NewOpInOperation creates a new OpInOperation operation.
func NewOpInOperation(path []string, values []interface{}) *OpInOperation {
	return &OpInOperation{
		BaseOp: NewBaseOp(path),
		Values: values,
	}
}

// Op returns the operation type.
func (op *OpInOperation) Op() internal.OpType {
	return internal.OpInType
}

// Code returns the operation code.
func (op *OpInOperation) Code() int {
	return internal.OpInCode
}

// Path returns the operation path.
func (op *OpInOperation) Path() []string {
	return op.path
}

// Test evaluates the in predicate condition.
func (op *OpInOperation) Test(doc any) (bool, error) {
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
func (op *OpInOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, found, err := op.getValueAndCheckInArray(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !found {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %v is not in array %v", ErrOperationFailed, value, op.Values)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getValueAndCheckInArray retrieves the value and checks if it's in the array
func (op *OpInOperation) getValueAndCheckInArray(doc any) (interface{}, bool, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return nil, false, ErrPathNotFound
	}

	// Check if the value is in the specified array
	for _, v := range op.Values {
		if reflect.DeepEqual(val, v) {
			return val, true, nil
		}
	}

	return val, false, nil
}

// ToJSON serializes the operation to JSON format.
func (op *OpInOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":    string(internal.OpInType),
		"path":  formatPath(op.Path()),
		"value": op.Values,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpInOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpInCode, op.Path(), op.Values}, nil
}

// Validate validates the in operation.
func (op *OpInOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(op.Values) == 0 {
		return ErrValuesArrayEmpty
	}
	return nil
}

// Not returns false since this is not a NOT operation.
func (op *OpInOperation) Not() bool {
	return false
}

// Short aliases for common use
var (
	// NewIn creates a new in operation
	NewIn = NewOpInOperation
)
