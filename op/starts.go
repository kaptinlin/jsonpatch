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
func (s *StartsOperation) Op() internal.OpType {
	return internal.OpStartsType
}

// Code returns the operation code.
func (s *StartsOperation) Code() int {
	return internal.OpStartsCode
}

// hasPrefix checks whether str starts with the expected prefix,
// respecting the IgnoreCase flag.
func (s *StartsOperation) hasPrefix(str string) bool {
	if s.IgnoreCase {
		return strings.HasPrefix(strings.ToLower(str), strings.ToLower(s.Value))
	}
	return strings.HasPrefix(str, s.Value)
}

// Test evaluates the starts predicate condition.
func (s *StartsOperation) Test(doc any) (bool, error) {
	_, str, err := validateString(doc, s.Path())
	if err != nil {
		//nolint:nilerr // intentional: path not found or wrong type means test fails
		return false, nil
	}
	return s.hasPrefix(str), nil
}

// Apply applies the starts test operation to the document.
func (s *StartsOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, str, err := validateString(doc, s.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	if !s.hasPrefix(str) {
		return internal.OpResult[any]{}, fmt.Errorf("%w: string %q does not start with %q", ErrStringMismatch, str, s.Value)
	}

	return internal.OpResult[any]{Doc: doc, Old: value}, nil
}

// ToJSON serializes the operation to JSON format.
func (s *StartsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpStartsType),
		Path:       formatPath(s.Path()),
		Value:      s.Value,
		IgnoreCase: s.IgnoreCase,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (s *StartsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStartsCode, s.Path(), s.Value}, nil
}

// Validate validates the starts operation.
func (s *StartsOperation) Validate() error {
	if len(s.Path()) == 0 {
		return ErrPathEmpty
	}
	return nil
}
