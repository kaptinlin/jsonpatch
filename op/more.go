package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpMoreOperation represents a "more" predicate operation that checks if a value is greater than a specified number.
type OpMoreOperation struct {
	PredicateOpBase
	Value float64 // The number to compare against
}

// NewOpMoreOperation creates a new more operation.
func NewOpMoreOperation(path []string, value float64) *OpMoreOperation {
	return &OpMoreOperation{
		PredicateOpBase: PredicateOpBase{
			BaseOp: BaseOp{path: path},
		},
		Value: value,
	}
}

// Op returns the operation type.
func (o *OpMoreOperation) Op() internal.OpType {
	return internal.OpMoreType
}

// Code returns the operation code.
func (o *OpMoreOperation) Code() int {
	return internal.OpMoreCode
}

// Test evaluates the more predicate condition.
func (o *OpMoreOperation) Test(doc interface{}) (bool, error) {
	_, num, err := o.getAndValidateValue(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return num > o.Value, nil
}

// Apply applies the more operation.
func (o *OpMoreOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, num, err := o.getAndValidateValue(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if num <= o.Value {
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not greater than %f", ErrComparisonFailed, num, o.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// getAndValidateValue retrieves and validates the numeric value at the path
func (o *OpMoreOperation) getAndValidateValue(doc interface{}) (interface{}, float64, error) {
	// Get target value
	val, err := getValue(doc, o.Path())
	if err != nil {
		return nil, 0, ErrPathNotFound
	}

	// Convert to float64 for comparison
	var num float64
	switch v := val.(type) {
	case int:
		num = float64(v)
	case int32:
		num = float64(v)
	case int64:
		num = float64(v)
	case float32:
		num = float64(v)
	case float64:
		num = v
	default:
		return nil, 0, ErrNotNumber
	}

	return val, num, nil
}

// ToJSON converts the operation to JSON representation.
func (o *OpMoreOperation) ToJSON() (internal.Operation, error) {
	// Convert float64 to int if it's a whole number
	var value interface{} = o.Value
	if o.Value == float64(int(o.Value)) {
		value = int(o.Value)
	}

	return internal.Operation{
		"op":    string(internal.OpMoreType),
		"path":  formatPath(o.Path()),
		"value": value,
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (o *OpMoreOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoreCode, o.Path(), o.Value}, nil
}

// Validate validates the more operation.
func (o *OpMoreOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Path returns the path for the more operation.
func (o *OpMoreOperation) Path() []string {
	return o.path
}

// Short aliases for common use
var (
	// NewMore creates a new more operation
	NewMore = NewOpMoreOperation
)
