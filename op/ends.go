package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// EndsOperation represents a test operation that checks if a string value ends with a specific suffix.
type EndsOperation struct {
	BaseOp
	Value      string `json:"value"`       // Expected suffix
	IgnoreCase bool   `json:"ignore_case"` // Whether to ignore case
	NotFlag    bool   `json:"not"`         // Whether to negate the result
}

type OpEndsOperation = EndsOperation //nolint:revive // Backward compatibility alias

// NewOpEndsOperation creates a new OpEndsOperation operation.
func NewOpEndsOperation(path []string, suffix string) *EndsOperation {
	return &EndsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      suffix,
		IgnoreCase: false,
	}
}

// NewOpEndsOperationWithIgnoreCase creates a new OpEndsOperation operation with ignore case option.
func NewOpEndsOperationWithIgnoreCase(path []string, suffix string, ignoreCase bool) *EndsOperation {
	return &EndsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      suffix,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (op *EndsOperation) Op() internal.OpType {
	return internal.OpEndsType
}

// Code returns the operation code.
func (op *EndsOperation) Code() int {
	return internal.OpEndsCode
}

// Path returns the operation path.
func (op *EndsOperation) Path() []string {
	return op.path
}

// Test evaluates the ends predicate condition.
func (op *EndsOperation) Test(doc any) (bool, error) {
	_, str, err := op.getAndValidateString(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	var result bool
	if op.IgnoreCase {
		result = strings.HasSuffix(strings.ToLower(str), strings.ToLower(op.Value))
	} else {
		result = strings.HasSuffix(str, op.Value)
	}
	
	// Apply negation if needed
	return result != op.NotFlag, nil
}

// Apply applies the ends test operation to the document.
func (op *EndsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := op.getAndValidateString(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	var hasSuffix bool
	if op.IgnoreCase {
		hasSuffix = strings.HasSuffix(strings.ToLower(str), strings.ToLower(op.Value))
	} else {
		hasSuffix = strings.HasSuffix(str, op.Value)
	}

	// Apply negation if needed
	result := hasSuffix != op.NotFlag

	if !result {
		if op.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: string %q ends with %q (negated test failed)", ErrStringMismatch, str, op.Value)
		}
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not end with %q", ErrStringMismatch, str, op.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getAndValidateString retrieves and validates the string value at the path
func (op *EndsOperation) getAndValidateString(doc any) (interface{}, string, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return nil, "", ErrPathNotFound
	}

	// Convert to string
	str, err := toString(val)
	if err != nil {
		return nil, "", ErrNotString
	}

	return val, str, nil
}

// ToJSON serializes the operation to JSON format.
func (op *EndsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpEndsType),
		Path:  formatPath(op.Path()),
		Value: op.Value,
	}
	if op.IgnoreCase {
		result.IgnoreCase = op.IgnoreCase
	}
	if op.NotFlag {
		result.Not = op.NotFlag
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *EndsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpEndsCode, op.Path(), op.Value}, nil
}

// Not returns the negation flag for this operation.
func (op *EndsOperation) Not() bool {
	return op.NotFlag
}

// Validate validates the ends operation.
func (op *EndsOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Short aliases for common use
var (
	// NewEnds creates a new ends operation
	NewEnds = NewOpEndsOperation
	// NewEndsWithIgnoreCase creates a new ends operation with ignore case
	NewEndsWithIgnoreCase = NewOpEndsOperationWithIgnoreCase
)
