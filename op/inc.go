package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpIncOperation represents an increment operation that increments a numeric value.
type OpIncOperation struct {
	BaseOp
	Inc float64 `json:"inc"` // Increment value
}

// NewOpIncOperation creates a new OpIncOperation operation.
func NewOpIncOperation(path []string, inc float64) *OpIncOperation {
	return &OpIncOperation{
		BaseOp: NewBaseOp(path),
		Inc:    inc,
	}
}

// Op returns the operation type.
func (op *OpIncOperation) Op() internal.OpType {
	return internal.OpIncType
}

// Code returns the operation code.
func (op *OpIncOperation) Code() int {
	return internal.OpIncCode
}

// Path returns the operation path.
func (op *OpIncOperation) Path() []string {
	return op.path
}

// Apply applies the increment operation to the document.
func (op *OpIncOperation) Apply(doc any) (internal.OpResult, error) {
	if len(op.path) == 0 {
		// Root level increment
		oldValue, ok := toFloat64(doc)
		if !ok {
			return internal.OpResult{}, ErrNotNumber
		}
		result := oldValue + op.Inc
		return internal.OpResult{Doc: result, Old: oldValue}, nil
	}

	if !pathExists(doc, op.path) {
		return internal.OpResult{}, ErrPathNotFound
	}

	parent, key, err := navigateToParent(doc, op.path)
	if err != nil {
		return internal.OpResult{}, ErrPathNotFound
	}

	currentValue := getValueFromParent(parent, key)
	oldValue, ok := toFloat64(currentValue)
	if !ok {
		return internal.OpResult{}, ErrNotNumber
	}
	result := oldValue + op.Inc

	if err := op.updateParent(parent, key, result); err != nil {
		return internal.OpResult{}, err
	}

	return internal.OpResult{Doc: doc, Old: oldValue}, nil
}

// updateParent updates the parent container with the new value
func (op *OpIncOperation) updateParent(parent interface{}, key interface{}, value interface{}) error {
	return updateParent(parent, key, value)
}

// ToJSON serializes the operation to JSON format.
func (op *OpIncOperation) ToJSON() (internal.Operation, error) {
	// Convert float64 to int if it's a whole number
	var inc interface{} = op.Inc
	if op.Inc == float64(int(op.Inc)) {
		inc = int(op.Inc)
	}

	return internal.Operation{
		"op":   string(internal.OpIncType),
		"path": formatPath(op.path),
		"inc":  inc,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpIncOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpIncCode, op.path, op.Inc}, nil
}

// Validate validates the increment operation.
func (op *OpIncOperation) Validate() error {
	// Empty path (root level) is allowed for inc operations
	return nil
}
