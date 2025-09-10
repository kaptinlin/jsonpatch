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

// NewOpStartsOperation creates a new OpStartsOperation operation.
func NewOpStartsOperation(path []string, prefix string) *StartsOperation {
	return &StartsOperation{
		BaseOp:     NewBaseOp(path),
		Value:      prefix,
		IgnoreCase: false,
	}
}

// NewOpStartsOperationWithIgnoreCase creates a new OpStartsOperation operation with ignore case option.
func NewOpStartsOperationWithIgnoreCase(path []string, prefix string, ignoreCase bool) *StartsOperation {
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

// Path returns the operation path.
func (op *StartsOperation) Path() []string {
	return op.path
}

// Test evaluates the starts predicate condition.
func (op *StartsOperation) Test(doc any) (bool, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		// Path access error means the path doesn't exist, treat as non-string
		// For JSON Patch test operations, path not found or wrong type means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}

	// Convert to string or from byte slice
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return false, nil // Return false if not string or byte slice
	}

	if op.IgnoreCase {
		return strings.HasPrefix(strings.ToLower(str), strings.ToLower(op.Value)), nil
	}
	return strings.HasPrefix(str, op.Value), nil
}

// Apply applies the starts test operation to the document.
func (op *StartsOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if value is a string or convert byte slice to string
	var str string
	switch v := val.(type) {
	case string:
		str = v
	case []byte:
		str = string(v)
	default:
		return internal.OpResult[any]{}, ErrNotString
	}

	// Check if string starts with the prefix
	var hasPrefix bool
	if op.IgnoreCase {
		hasPrefix = strings.HasPrefix(strings.ToLower(str), strings.ToLower(op.Value))
	} else {
		hasPrefix = strings.HasPrefix(str, op.Value)
	}

	if !hasPrefix {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not start with %q", ErrStringMismatch, str, op.Value)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *StartsOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":    string(internal.OpStartsType),
		"path":  formatPath(op.Path()),
		"value": op.Value,
	}
	if op.IgnoreCase {
		result["ignore_case"] = op.IgnoreCase
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

// Short aliases for common use
var (
	// NewStarts creates a new starts operation
	NewStarts = NewOpStartsOperation
	// NewStartsWithIgnoreCase creates a new starts operation with ignore case
	NewStartsWithIgnoreCase = NewOpStartsOperationWithIgnoreCase
)
