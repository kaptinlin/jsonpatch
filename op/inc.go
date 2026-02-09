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
func (o *IncOperation) Op() internal.OpType {
	return internal.OpIncType
}

// Code returns the operation code.
func (o *IncOperation) Code() int {
	return internal.OpIncCode
}

// Apply applies the increment operation to the document.
func (o *IncOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(o.path) == 0 {
		// Root level increment
		oldValue, ok := ToFloat64(doc)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
		result := oldValue + o.Inc
		return internal.OpResult[any]{Doc: result, Old: oldValue}, nil
	}

	parent, key, err := navigateToParent(doc, o.path)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Check if path exists and get current value
	var currentValue any
	var oldValue float64
	if pathExists(doc, o.path) {
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
	result := oldValue + o.Inc

	if err := updateParent(parent, key, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: currentValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *IncOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpIncType),
		Path: formatPath(o.path),
		Inc:  o.Inc,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *IncOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpIncCode, o.path, o.Inc}, nil
}

// Validate validates the increment operation.
func (o *IncOperation) Validate() error {
	return nil
}
