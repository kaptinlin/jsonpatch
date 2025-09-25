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
	NotFlag    bool   `json:"not"`         // Whether to negate the result
}

type OpContainsOperation = ContainsOperation //nolint:revive // Backward compatibility alias

// NewOpContainsOperation creates a new OpContainsOperation operation.
func NewOpContainsOperation(path []string, substring string) *ContainsOperation {
	return &OpContainsOperation{
		BaseOp: NewBaseOp(path),
		Value:  substring,
	}
}

// NewOpContainsOperationWithIgnoreCase creates a new OpContainsOperation operation with ignore case option.
func NewOpContainsOperationWithIgnoreCase(path []string, substring string, ignoreCase bool) *ContainsOperation {
	return &OpContainsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      substring,
		IgnoreCase: ignoreCase,
	}
}

// NewOpContainsOperationWithFlags creates a new OpContainsOperation operation with full options.
func NewOpContainsOperationWithFlags(path []string, substring string, ignoreCase bool, notFlag bool) *ContainsOperation {
	return &OpContainsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      substring,
		IgnoreCase: ignoreCase,
		NotFlag:    notFlag,
	}
}

// Apply applies the contains test operation to the document.
func (op *ContainsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, actualValue, testValue, testString, err := op.getAndPrepareStrings(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	contains := strings.Contains(testValue, testString)

	// Apply negation if needed
	if op.NotFlag {
		contains = !contains
	}

	if !contains {
		if op.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: string %q contains %q", ErrStringMismatch, actualValue, op.Value)
		}
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

	contains := strings.Contains(testValue, testString)

	// Apply negation if needed
	if op.NotFlag {
		contains = !contains
	}

	return contains, nil
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

// Not returns the negation flag.
func (op *ContainsOperation) Not() bool {
	return op.NotFlag
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
		Not:        op.NotFlag,
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

// Path returns the path for the contains operation.
func (op *ContainsOperation) Path() []string {
	return op.path
}

// Short aliases for common use
var (
	// NewContains creates a new contains operation
	NewContains = NewOpContainsOperation
	// NewContainsWithIgnoreCase creates a new contains operation with ignore case
	NewContainsWithIgnoreCase = NewOpContainsOperationWithIgnoreCase
)
