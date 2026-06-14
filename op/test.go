package op

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestOperation represents a test operation that checks if a value equals a specified value.
type TestOperation struct {
	BaseOp
	Value   any  `json:"value"`         // Expected value
	NotFlag bool `json:"not,omitempty"` // Whether to negate the test
}

// NewTest creates a new test operation.
func NewTest(path []string, value any) *TestOperation {
	return &TestOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// NewTestWithNot creates a new test operation with not flag.
func NewTestWithNot(path []string, value any, not bool) *TestOperation {
	return &TestOperation{
		BaseOp:  NewBaseOp(path),
		Value:   value,
		NotFlag: not,
	}
}

// Op returns the operation type.
func (t *TestOperation) Op() internal.OpType {
	return internal.OpTestType
}

// Test performs the test operation.
func (t *TestOperation) Test(doc any) (bool, error) {
	// Get target value
	target, err := value(doc, t.Path())
	if err != nil {
		//nolint:nilerr // This is intentional JSON Patch behavior - path not found is not an error
		return t.NotFlag, nil
	}

	// Compare values and apply negation using XOR logic
	result := deepEqual(target, t.Value)
	return result != t.NotFlag, nil
}

// Not returns whether this operation is negated.
func (t *TestOperation) Not() bool {
	return t.NotFlag
}

// Apply applies the test operation.
func (t *TestOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, err := value(doc, t.path)
	if err != nil {
		shouldPass := t.NotFlag
		if !shouldPass {
			return internal.OpResult[any]{}, fmt.Errorf("%w: path not found", ErrTestOperationFailed)
		}
		return internal.OpResult[any]{Doc: doc, Old: nil}, nil
	}

	isEqual := deepEqual(value, t.Value)

	// Determine if test should pass using XOR logic
	shouldPass := isEqual != t.NotFlag

	if !shouldPass {
		// Test failed
		if t.NotFlag {
			return internal.OpResult[any]{}, fmt.Errorf("%w: expected not %v, but got %v", ErrTestOperationFailed, t.Value, value)
		}
		// Check if it's a string vs number comparison for specific error message
		if _, ok := t.Value.(string); ok {
			if _, ok := value.(float64); ok {
				return internal.OpResult[any]{}, ErrTestOperationNumberStringMismatch
			}
			return internal.OpResult[any]{}, ErrTestOperationStringNotEquivalent
		}
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected %v, got %v", ErrTestOperationFailed, t.Value, value)
	}

	// Test operations don't modify the document and return nil for old value
	return internal.OpResult[any]{Doc: doc, Old: nil}, nil
}

// Validate validates the test operation.
func (t *TestOperation) Validate() error {
	// Empty path is valid root-document testing.
	return nil
}
