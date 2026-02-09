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

// Apply applies the less test operation to the document.
func (op *LessOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, err := getNumericValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if actualValue >= op.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not less than %f", ErrComparisonFailed, actualValue, op.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// Op returns the operation type.
func (op *LessOperation) Op() internal.OpType {
	return internal.OpLessType
}

// Code returns the operation code.
func (op *LessOperation) Code() int {
	return internal.OpLessCode
}

// ToJSON serializes the operation to JSON format.
func (op *LessOperation) ToJSON() (internal.Operation, error) {
	// Convert float64 to int if it's a whole number
	var value interface{} = op.Value
	if op.Value == float64(int(op.Value)) {
		value = int(op.Value)
	}

	return internal.Operation{
		Op:    string(internal.OpLessType),
		Path:  formatPath(op.Path()),
		Value: value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *LessOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpLessCode, op.Path(), op.Value}, nil
}

// Validate validates the less operation.
func (op *LessOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Test evaluates the less predicate condition.
func (op *LessOperation) Test(doc any) (bool, error) {
	_, actualValue, err := getNumericValue(doc, op.Path())
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return actualValue < op.Value, nil
}

