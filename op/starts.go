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
		BaseOp: NewBaseOp(path),
		Value:  prefix,
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
func (o *StartsOperation) Op() internal.OpType {
	return internal.OpStartsType
}

// Code returns the operation code.
func (o *StartsOperation) Code() int {
	return internal.OpStartsCode
}

// Test evaluates the starts predicate condition.
func (o *StartsOperation) Test(doc any) (bool, error) {
	_, str, err := getAndValidateString(doc, o.Path())
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	if o.IgnoreCase {
		return strings.HasPrefix(strings.ToLower(str), strings.ToLower(o.Value)), nil
	}
	return strings.HasPrefix(str, o.Value), nil
}

// Apply applies the starts test operation to the document.
func (o *StartsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := getAndValidateString(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	var hasPrefix bool
	if o.IgnoreCase {
		hasPrefix = strings.HasPrefix(strings.ToLower(str), strings.ToLower(o.Value))
	} else {
		hasPrefix = strings.HasPrefix(str, o.Value)
	}

	if !hasPrefix {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not start with %q", ErrStringMismatch, str, o.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *StartsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpStartsType),
		Path:  formatPath(o.Path()),
		Value: o.Value,
	}
	if o.IgnoreCase {
		result.IgnoreCase = o.IgnoreCase
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *StartsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStartsCode, o.Path(), o.Value}, nil
}

// Validate validates the starts operation.
func (o *StartsOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
