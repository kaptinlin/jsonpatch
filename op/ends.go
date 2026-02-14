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

// hasSuffix checks whether str ends with the expected suffix,
// respecting the IgnoreCase flag.
func (e *EndsOperation) hasSuffix(str string) bool {
	if e.IgnoreCase {
		return strings.HasSuffix(strings.ToLower(str), strings.ToLower(e.Value))
	}
	return strings.HasSuffix(str, e.Value)
}

// Test evaluates the ends predicate condition.
func (e *EndsOperation) Test(doc any) (bool, error) {
	_, str, err := validateString(doc, e.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return e.hasSuffix(str), nil
}

// Apply applies the ends test operation to the document.
func (e *EndsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := validateString(doc, e.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !e.hasSuffix(str) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not end with %q", ErrStringMismatch, str, e.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (e *EndsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpEndsType),
		Path:       formatPath(e.Path()),
		Value:      e.Value,
		IgnoreCase: e.IgnoreCase,
	}, nil
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
