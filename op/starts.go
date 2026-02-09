package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StartsOperation represents a test operation that checks if a string value starts with a specific prefix.
type StartsOperation struct {
	BaseOp
	Value      string `json:"value"`       // Expected prefix
	IgnoreCase bool   `json:"ignore_case"` // Whether to ignore case
}

// NewStarts creates a new starts operation.
func NewStarts(path []string, prefix string) *StartsOperation {
	return &StartsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      prefix,
		IgnoreCase: false,
	}
}

// NewStartsWithIgnoreCase creates a new starts operation with ignore case option.
func NewStartsWithIgnoreCase(path []string, prefix string, ignoreCase bool) *StartsOperation {
	return &StartsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      prefix,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (op *StartsOperation) Op() internal.OpType {
	return internal.OpStartsType
}

// Code returns the operation code.
func (op *StartsOperation) Code() int {
	return internal.OpStartsCode
}

// Test evaluates the starts predicate condition.
func (op *StartsOperation) Test(doc any) (bool, error) {
	_, str, err := op.getAndValidateString(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	if op.IgnoreCase {
		return strings.HasPrefix(strings.ToLower(str), strings.ToLower(op.Value)), nil
	}
	return strings.HasPrefix(str, op.Value), nil
}

// Apply applies the starts test operation to the document.
func (op *StartsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := op.getAndValidateString(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	var hasPrefix bool
	if op.IgnoreCase {
		hasPrefix = strings.HasPrefix(strings.ToLower(str), strings.ToLower(op.Value))
	} else {
		hasPrefix = strings.HasPrefix(str, op.Value)
	}

	if !hasPrefix {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not start with %q", ErrStringMismatch, str, op.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// getAndValidateString retrieves and validates the string value at the path
func (op *StartsOperation) getAndValidateString(doc any) (interface{}, string, error) {
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
func (op *StartsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpStartsType),
		Path:  formatPath(op.Path()),
		Value: op.Value,
	}
	if op.IgnoreCase {
		result.IgnoreCase = op.IgnoreCase
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *StartsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStartsCode, op.Path(), op.Value}, nil
}

// Validate validates the starts operation.
func (op *StartsOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

