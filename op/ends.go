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
}

// NewEnds creates a new ends operation.
func NewEnds(path []string, suffix string) *EndsOperation {
	return &EndsOperation{
		BaseOp: NewBaseOp(path),
		Value:  suffix,
	}
}

// NewEndsWithIgnoreCase creates a new ends operation with ignore case option.
func NewEndsWithIgnoreCase(path []string, suffix string, ignoreCase bool) *EndsOperation {
	return &EndsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      suffix,
		IgnoreCase: ignoreCase,
	}
}

// Op returns the operation type.
func (e *EndsOperation) Op() internal.OpType {
	return internal.OpEndsType
}

// Code returns the operation code.
func (e *EndsOperation) Code() int {
	return internal.OpEndsCode
}

// Test evaluates the ends predicate condition.
func (e *EndsOperation) Test(doc any) (bool, error) {
	_, str, err := getAndValidateString(doc, e.Path())
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	if e.IgnoreCase {
		return strings.HasSuffix(strings.ToLower(str), strings.ToLower(e.Value)), nil
	}
	return strings.HasSuffix(str, e.Value), nil
}

// Apply applies the ends test operation to the document.
func (e *EndsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := getAndValidateString(doc, e.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	var hasSuffix bool
	if e.IgnoreCase {
		hasSuffix = strings.HasSuffix(strings.ToLower(str), strings.ToLower(e.Value))
	} else {
		hasSuffix = strings.HasSuffix(str, e.Value)
	}

	if !hasSuffix {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not end with %q", ErrStringMismatch, str, e.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (e *EndsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpEndsType),
		Path:  formatPath(e.Path()),
		Value: e.Value,
	}
	if e.IgnoreCase {
		result.IgnoreCase = e.IgnoreCase
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (e *EndsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpEndsCode, e.Path(), e.Value}, nil
}

// Validate validates the ends operation.
func (e *EndsOperation) Validate() error {
	if len(e.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
