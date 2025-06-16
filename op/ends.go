package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpEndsOperation represents a test operation that checks if a string value ends with a specific suffix.
type OpEndsOperation struct {
	BaseOp
	Value      string `json:"value"`       // Expected suffix
	IgnoreCase bool   `json:"ignore_case"` // Whether to ignore case
}

// NewOpEndsOperation creates a new OpEndsOperation operation.
func NewOpEndsOperation(path []string, suffix string) *OpEndsOperation {
	return &OpEndsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      suffix,
		IgnoreCase: false,
	}
}

// NewOpEndsOperationWithIgnoreCase creates a new OpEndsOperation operation with ignore case option.
func NewOpEndsOperationWithIgnoreCase(path []string, suffix string, ignoreCase bool) *OpEndsOperation {
	return &OpEndsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      suffix,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (op *OpEndsOperation) Op() internal.OpType {
	return internal.OpEndsType
}

// Code returns the operation code.
func (op *OpEndsOperation) Code() int {
	return internal.OpEndsCode
}

// Path returns the operation path.
func (op *OpEndsOperation) Path() []string {
	return op.path
}

// Test evaluates the ends predicate condition.
func (op *OpEndsOperation) Test(doc any) (bool, error) {
	_, str, err := op.getAndValidateString(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	if op.IgnoreCase {
		return strings.HasSuffix(strings.ToLower(str), strings.ToLower(op.Value)), nil
	}
	return strings.HasSuffix(str, op.Value), nil
}

// Apply applies the ends test operation to the document.
func (op *OpEndsOperation) Apply(doc any) (internal.OpResult[any], error) {
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

	if !hasSuffix {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not end with %q", ErrStringMismatch, str, op.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getAndValidateString retrieves and validates the string value at the path
func (op *OpEndsOperation) getAndValidateString(doc any) (interface{}, string, error) {
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
func (op *OpEndsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":    string(internal.OpEndsType),
		"path":  formatPath(op.Path()),
		"value": op.Value,
	}
	if op.IgnoreCase {
		result["ignore_case"] = op.IgnoreCase
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpEndsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpEndsCode, op.Path(), op.Value}, nil
}

// Validate validates the ends operation.
func (op *OpEndsOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
