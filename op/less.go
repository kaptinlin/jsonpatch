package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpLessOperation represents a test operation that checks if a numeric value is less than a specified value.
type OpLessOperation struct {
	BaseOp
	Value float64 `json:"value"` // Value to compare against
}

// NewOpLessOperation creates a new OpLessOperation operation.
func NewOpLessOperation(path []string, value float64) *OpLessOperation {
	return &OpLessOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// Apply applies the less test operation to the document.
func (op *OpLessOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, err := op.getAndValidateValue(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if actualValue >= op.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not less than %f", ErrComparisonFailed, actualValue, op.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getAndValidateValue retrieves and validates the numeric value at the path
func (op *OpLessOperation) getAndValidateValue(doc any) (interface{}, float64, error) {
	value, err := getValue(doc, op.Path())
	if err != nil {
		return nil, 0, ErrPathNotFound
	}
	actualValue, ok := toFloat64(value)
	if !ok {
		return nil, 0, ErrNotNumber
	}
	return value, actualValue, nil
}

// Op returns the operation type.
func (op *OpLessOperation) Op() internal.OpType {
	return internal.OpLessType
}

// Code returns the operation code.
func (op *OpLessOperation) Code() int {
	return internal.OpLessCode
}

// ToJSON serializes the operation to JSON format.
func (op *OpLessOperation) ToJSON() (internal.Operation, error) {
	// Convert float64 to int if it's a whole number
	var value interface{} = op.Value
	if op.Value == float64(int(op.Value)) {
		value = int(op.Value)
	}

	return internal.Operation{
		"op":    string(internal.OpLessType),
		"path":  formatPath(op.Path()),
		"value": value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpLessOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpLessCode, op.Path(), op.Value}, nil
}

// Validate validates the less operation.
func (op *OpLessOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

func (op *OpLessOperation) Path() []string {
	return op.path
}

// Test evaluates the less predicate condition.
func (op *OpLessOperation) Test(doc any) (bool, error) {
	_, actualValue, err := op.getAndValidateValue(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return actualValue < op.Value, nil
}

// Not returns false since this is not a NOT operation.
func (op *OpLessOperation) Not() bool {
	return false
}
