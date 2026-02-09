package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// IncOperation represents an increment operation that increments a numeric value.
type IncOperation struct {
	BaseOp
	Inc float64 `json:"inc"` // Increment value
}

// NewInc creates a new increment operation.
func NewInc(path []string, inc float64) *IncOperation {
	return &IncOperation{
		BaseOp: NewBaseOp(path),
		Inc:    inc,
	}
}

// Op returns the operation type.
func (op *IncOperation) Op() internal.OpType {
	return internal.OpIncType
}

// Code returns the operation code.
func (op *IncOperation) Code() int {
	return internal.OpIncCode
}

// Apply applies the increment operation to the document.
func (op *IncOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(op.path) == 0 {
		// Root level increment
		oldValue, ok := ToFloat64(doc)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
		result := oldValue + op.Inc
		return internal.OpResult[any]{Doc: result, Old: oldValue}, nil
	}

	parent, key, err := navigateToParent(doc, op.path)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if path exists and get current value
	var currentValue any
	var oldValue float64
	if pathExists(doc, op.path) {
		currentValue = getValueFromParent(parent, key)
		var ok bool
		oldValue, ok = ToFloat64(currentValue)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
	} else {
		// Path doesn't exist, treat as undefined (which becomes 0 in JavaScript)
		currentValue = nil
		oldValue = 0
	}
	result := oldValue + op.Inc

	if err := updateParent(parent, key, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: currentValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *IncOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpIncType),
		Path: formatPath(op.path),
		Inc:  op.Inc,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *IncOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpIncCode, op.path, op.Inc}, nil
}

// Validate validates the increment operation.
func (op *IncOperation) Validate() error {
	return nil
}

