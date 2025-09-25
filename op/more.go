package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// MoreOperation represents a "more" predicate operation that checks if a value is greater than a specified number.
type MoreOperation struct {
	BaseOp
	Value   float64 `json:"value"` // The number to compare against
	NotFlag bool    `json:"not"`   // Whether to negate the result
}

// NewOpMoreOperation creates a new more operation.
func NewOpMoreOperation(path []string, value float64) *MoreOperation {
	return &MoreOperation{
		BaseOp:  NewBaseOp(path),
		Value:   value,
		NotFlag: false,
	}
}

// NewOpMoreOperationWithFlags creates a new more operation with full options.
func NewOpMoreOperationWithFlags(path []string, value float64, notFlag bool) *MoreOperation {
	return &MoreOperation{
		BaseOp:  NewBaseOp(path),
		Value:   value,
		NotFlag: notFlag,
	}
}

// Op returns the operation type.
func (o *MoreOperation) Op() internal.OpType {
	return internal.OpMoreType
}

// Code returns the operation code.
func (o *MoreOperation) Code() int {
	return internal.OpMoreCode
}

// Test evaluates the more predicate condition.
func (o *MoreOperation) Test(doc interface{}) (bool, error) {
	_, num, err := o.getAndValidateValue(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	
	greater := num > o.Value
	
	// Apply negation if needed
	if o.NotFlag {
		greater = !greater
	}
	
	return greater, nil
}

// Apply applies the more operation.
func (o *MoreOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, num, err := o.getAndValidateValue(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	greater := num > o.Value
	
	// Apply negation if needed
	if o.NotFlag {
		greater = !greater
	}

	if !greater {
		if o.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is greater than %f", ErrComparisonFailed, num, o.Value)
		}
		return internal.OpResult[any]{}, fmt.Errorf("%w: value %f is not greater than %f", ErrComparisonFailed, num, o.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// getAndValidateValue retrieves and validates the numeric value at the path
func (o *MoreOperation) getAndValidateValue(doc interface{}) (interface{}, float64, error) {
	value, err := getValue(doc, o.Path())
	if err != nil {
		return nil, 0, ErrPathNotFound
	}
	actualValue, ok := ToFloat64(value)
	if !ok {
		return nil, 0, ErrNotNumber
	}
	return value, actualValue, nil
}

// ToJSON converts the operation to JSON representation.
func (o *MoreOperation) ToJSON() (internal.Operation, error) {
	// Convert float64 to int if it's a whole number
	var value interface{} = o.Value
	if o.Value == float64(int(o.Value)) {
		value = int(o.Value)
	}

	return internal.Operation{
		Op:    string(internal.OpMoreType),
		Path:  formatPath(o.Path()),
		Value: value,
		Not:   o.NotFlag,
	}, nil
}

// ToCompact converts the operation to compact array representation.
func (o *MoreOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpMoreCode, o.Path(), o.Value}, nil
}

// Not returns the negation flag.
func (o *MoreOperation) Not() bool {
	return o.NotFlag
}

// Validate validates the more operation.
func (o *MoreOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Path returns the path for the more operation.
func (o *MoreOperation) Path() []string {
	return o.path
}

// Short aliases for common use
var (
	// NewMore creates a new more operation
	NewMore = NewOpMoreOperation
)
