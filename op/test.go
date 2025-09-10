package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpTestOperation represents a test operation that checks if a value equals a specified value.
type OpTestOperation struct {
	BaseOp
	Value   interface{} `json:"value"`         // Expected value
	NotFlag bool        `json:"not,omitempty"` // Whether to negate the test
}

// NewOpTestOperation creates a new OpTestOperation operation.
func NewOpTestOperation(path []string, value interface{}) *OpTestOperation {
	return &OpTestOperation{
		BaseOp:  NewBaseOp(path),
		Value:   value,
		NotFlag: false,
	}
}

// Op returns the operation type.
func (o *OpTestOperation) Op() internal.OpType {
	return internal.OpTestType
}

// Code returns the operation code.
func (o *OpTestOperation) Code() int {
	return internal.OpTestCode
}

// Test performs the test operation.
func (o *OpTestOperation) Test(doc interface{}) (bool, error) {
	// Get target value
	target, err := getValue(doc, o.Path())
	if err != nil {
		// If path not found and we're negating, that's success
		if o.NotFlag {
			return true, nil
		}
		return false, err
	}

	// Compare values and apply negation using XOR logic
	result := deepEqual(target, o.Value)
	return result != o.NotFlag, nil
}

// Not returns whether this operation is negated.
func (o *OpTestOperation) Not() bool {
	return o.NotFlag
}

// Path returns the operation path.
func (o *OpTestOperation) Path() []string {
	return o.path
}

// Apply applies the test operation.
func (o *OpTestOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, err := getValue(doc, o.path)
	if err != nil {
		// If path not found and we're negating, that's success
		if o.NotFlag {
			return internal.OpResult[any]{Doc: doc, Old: nil}, nil
		}
		return internal.OpResult[any]{}, err
	}

	isEqual := deepEqual(value, o.Value)

	// Determine if test should pass using XOR logic
	shouldPass := isEqual != o.NotFlag

	if !shouldPass {
		// Test failed
		if o.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: expected not %v, but got %v", ErrTestOperationFailed, o.Value, value)
		}
		// Check if it's a string vs number comparison for specific error message
		if _, ok := o.Value.(string); ok {
			if _, ok := value.(float64); ok {
				return internal.OpResult[any]{}, ErrTestOperationNumberStringMismatch
			}
			return internal.OpResult[any]{}, ErrTestOperationStringNotEquivalent
		}
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected %v, got %v", ErrTestOperationFailed, o.Value, value)
	}

	// Test operations don't modify the document and return nil for old value
	return internal.OpResult[any]{Doc: doc, Old: nil}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpTestOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":    string(internal.OpTestType),
		"path":  formatPath(o.path),
		"value": o.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpTestOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestCode, o.path, o.Value}, nil
}

// Validate validates the test operation.
func (o *OpTestOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Short aliases for common use
var (
	// NewTest creates a new test operation
	NewTest = NewOpTestOperation
)
