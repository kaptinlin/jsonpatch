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
func (op *ContainsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, testValue, testString, err := op.getAndPrepareStrings(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !strings.Contains(testValue, testString) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not contain %q", ErrStringMismatch, actualValue, op.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// Test performs the contains test operation.
func (op *ContainsOperation) Test(doc any) (bool, error) {
	_, _, testValue, testString, err := op.getAndPrepareStrings(doc)
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
func (op *ContainsOperation) getAndPrepareStrings(doc any) (interface{}, string, string, string, error) {
	value, err := getValue(doc, op.Path())
	if err != nil {
		return nil, "", "", "", ErrPathNotFound
	}

	actualValue, err := toString(value)
	if err != nil {
		return nil, "", "", "", ErrNotString
	}

	// Fast path: case-sensitive comparison (most common case)
	if !op.IgnoreCase {
		return value, actualValue, actualValue, op.Value, nil
	}

	// Slower path: case-insensitive comparison
	testValue := strings.ToLower(actualValue)
	testString := strings.ToLower(op.Value)

	return value, actualValue, testValue, testString, nil
}

// Op returns the operation type.
func (op *ContainsOperation) Op() internal.OpType {
	return internal.OpContainsType
}

// Code returns the operation code.
func (op *ContainsOperation) Code() int {
	return internal.OpContainsCode
}

// ToJSON serializes the operation to JSON format.
func (op *ContainsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpContainsType),
		Path:       formatPath(op.Path()),
		Value:      op.Value,
		IgnoreCase: op.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *ContainsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpContainsCode, op.Path(), op.Value}, nil
}

// Validate validates the contains operation.
func (op *ContainsOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

