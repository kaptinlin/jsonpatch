package op

import (
	"fmt"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpContainsOperation represents a test operation that checks if a string value contains a substring.
type OpContainsOperation struct {
	BaseOp
	Value      string `json:"value"`       // Substring to search for
	IgnoreCase bool   `json:"ignore_case"` // Whether to ignore case when comparing
}

// ContainsOperation is a non-stuttering alias for OpContainsOperation.
type ContainsOperation = OpContainsOperation

// NewOpContainsOperation creates a new OpContainsOperation operation.
func NewOpContainsOperation(path []string, substring string) *OpContainsOperation {
	return &OpContainsOperation{
		BaseOp: NewBaseOp(path),
		Value:  substring,
	}
}

// NewOpContainsOperationWithIgnoreCase creates a new OpContainsOperation operation with ignore case option.
func NewOpContainsOperationWithIgnoreCase(path []string, substring string, ignoreCase bool) *OpContainsOperation {
	return &OpContainsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      substring,
		IgnoreCase: ignoreCase,
	}
}

// Apply applies the contains test operation to the document.
func (op *OpContainsOperation) Apply(doc any) (internal.OpResult[any], error) {
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
func (op *OpContainsOperation) Test(doc any) (bool, error) {
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
func (op *OpContainsOperation) getAndPrepareStrings(doc any) (interface{}, string, string, string, error) {
	value, err := getValue(doc, op.Path())
	if err != nil {
		return nil, "", "", "", ErrPathNotFound
	}

	actualValue, err := toString(value)
	if err != nil {
		return nil, "", "", "", ErrNotString
	}

	testValue := actualValue
	testString := op.Value
	if op.IgnoreCase {
		testValue = strings.ToLower(actualValue)
		testString = strings.ToLower(op.Value)
	}

	return value, actualValue, testValue, testString, nil
}

// Not returns false (contains operation doesn't support not modifier).
func (op *OpContainsOperation) Not() bool {
	return false
}

// Op returns the operation type.
func (op *OpContainsOperation) Op() internal.OpType {
	return internal.OpContainsType
}

// Code returns the operation code.
func (op *OpContainsOperation) Code() int {
	return internal.OpContainsCode
}

// ToJSON serializes the operation to JSON format.
func (op *OpContainsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":    string(internal.OpContainsType),
		"path":  formatPath(op.Path()),
		"value": op.Value,
	}
	if op.IgnoreCase {
		result["ignore_case"] = op.IgnoreCase
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpContainsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpContainsCode, op.Path(), op.Value}, nil
}

// Validate validates the contains operation.
func (op *OpContainsOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Path returns the path for the contains operation.
func (op *OpContainsOperation) Path() []string {
	return op.path
}

// Short aliases for common use
var (
	// NewContains creates a new contains operation
	NewContains = NewOpContainsOperation
	// NewContainsWithIgnoreCase creates a new contains operation with ignore case
	NewContainsWithIgnoreCase = NewOpContainsOperationWithIgnoreCase
)
