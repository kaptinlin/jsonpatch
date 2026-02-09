package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// ContainsOperation represents a contains operation that tests if a string contains a substring.
type ContainsOperation struct {
	BaseOp
	Value      string `json:"value"`       // Substring to search for
	IgnoreCase bool   `json:"ignore_case"` // Whether to ignore case when comparing
}

// NewContains creates a new contains operation.
func NewContains(path []string, substring string) *ContainsOperation {
	return &ContainsOperation{
		BaseOp: NewBaseOp(path),
		Value:  substring,
	}
}

// NewContainsWithIgnoreCase creates a new contains operation with ignore case option.
func NewContainsWithIgnoreCase(path []string, substring string, ignoreCase bool) *ContainsOperation {
	return &ContainsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      substring,
		IgnoreCase: ignoreCase,
	}
}

// Apply applies the contains test operation to the document.
func (o *ContainsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, testValue, testString, err := o.getAndPrepareStrings(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !strings.Contains(testValue, testString) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not contain %q", ErrStringMismatch, actualValue, o.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// Test performs the contains test operation.
func (o *ContainsOperation) Test(doc any) (bool, error) {
	_, _, testValue, testString, err := o.getAndPrepareStrings(doc)
	if err != nil {
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	return strings.Contains(testValue, testString), nil
}

// getAndPrepareStrings retrieves the value, converts to string, and prepares test strings
// Optimized to avoid unnecessary allocations in case-sensitive operations
func (o *ContainsOperation) getAndPrepareStrings(doc any) (any, string, string, string, error) {
	value, err := getValue(doc, o.Path())
	if err != nil {
		return nil, "", "", "", ErrPathNotFound
	}

	actualValue, err := toString(value)
	if err != nil {
		return nil, "", "", "", ErrNotString
	}

	// Fast path: case-sensitive comparison (most common case)
	if !o.IgnoreCase {
		return value, actualValue, actualValue, o.Value, nil
	}

	// Slower path: case-insensitive comparison
	testValue := strings.ToLower(actualValue)
	testString := strings.ToLower(o.Value)

	return value, actualValue, testValue, testString, nil
}

// Op returns the operation type.
func (o *ContainsOperation) Op() internal.OpType {
	return internal.OpContainsType
}

// Code returns the operation code.
func (o *ContainsOperation) Code() int {
	return internal.OpContainsCode
}

// ToJSON serializes the operation to JSON format.
func (o *ContainsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpContainsType),
		Path:       formatPath(o.Path()),
		Value:      o.Value,
		IgnoreCase: o.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *ContainsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpContainsCode, o.Path(), o.Value}, nil
}

// Validate validates the contains operation.
func (o *ContainsOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
