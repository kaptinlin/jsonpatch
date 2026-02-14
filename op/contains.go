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
func (co *ContainsOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, actualValue, testValue, testString, err := co.prepareStrings(doc)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !strings.Contains(testValue, testString) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not contain %q", ErrStringMismatch, actualValue, co.Value)
	}
	return internal.OpResult[any]{Doc: doc, Old: val}, nil
}

// Test performs the contains test operation.
func (co *ContainsOperation) Test(doc any) (bool, error) {
	_, _, testValue, testString, err := co.prepareStrings(doc)
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return strings.Contains(testValue, testString), nil
}

// prepareStrings retrieves the value, converts to string, and prepares
// comparison strings (lowercased when IgnoreCase is set).
// Returns: originalValue, actualString, comparableString, comparableTarget, error.
func (co *ContainsOperation) prepareStrings(doc any) (any, string, string, string, error) {
	val, err := value(doc, co.Path())
	if err != nil {
		return nil, "", "", "", ErrPathNotFound
	}

	actualValue, err := toString(val)
	if err != nil {
		return nil, "", "", "", ErrNotString
	}

	if !co.IgnoreCase {
		return val, actualValue, actualValue, co.Value, nil
	}

	return val, actualValue, strings.ToLower(actualValue), strings.ToLower(co.Value), nil
}

// Op returns the operation type.
func (co *ContainsOperation) Op() internal.OpType {
	return internal.OpContainsType
}

// Code returns the operation code.
func (co *ContainsOperation) Code() int {
	return internal.OpContainsCode
}

// ToJSON serializes the operation to JSON format.
func (co *ContainsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpContainsType),
		Path:       formatPath(co.Path()),
		Value:      co.Value,
		IgnoreCase: co.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (co *ContainsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpContainsCode, co.Path(), co.Value}, nil
}

// Validate validates the contains operation.
func (co *ContainsOperation) Validate() error {
	if len(co.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
