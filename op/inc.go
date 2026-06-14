package op

import "github.com/kaptinlin/jsonpatch/internal"

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
func (ic *IncOperation) Op() internal.OpType {
	return internal.OpIncType
}

// Apply applies the increment operation to the document.
func (ic *IncOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(ic.path) == 0 {
		// Root level increment
		oldValue, ok := ToFloat64(doc)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
		result := oldValue + ic.Inc
		return internal.OpResult[any]{Doc: result, Old: oldValue}, nil
	}

	parent, key, err := navigateToParent(doc, ic.path)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	// Missing final targets are treated as 0, creating the field with the incremented value.
	var currentValue any
	var oldValue float64
	if pathExists(doc, ic.path) {
		currentValue = valueFromParent(parent, key)
		var ok bool
		oldValue, ok = ToFloat64(currentValue)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
	}
	result := oldValue + ic.Inc

	if err := updateParent(parent, key, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: currentValue}, nil
}

// Validate validates the increment operation.
func (ic *IncOperation) Validate() error {
	return nil
}
