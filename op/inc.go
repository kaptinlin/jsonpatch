package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// IncOperation represents an increment operation that increments a numeric value.
type IncOperation struct {
	BaseOp
	Inc float64 `json:"inc"` // Increment value
}

// NewOpIncOperation creates a new OpIncOperation operation.
func NewOpIncOperation(path []string, inc float64) *IncOperation {
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

// Path returns the operation path.
func (op *IncOperation) Path() []string {
	return op.path
}

// Apply applies the increment operation to the document.
func (op *IncOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(op.path) == 0 {
		// Root level increment
		oldValue, ok := toFloat64(doc)
		if !ok {
			return internal.OpResult[any]{}, ErrNotNumber
		}
		result := oldValue + op.Inc
		return internal.OpResult[any]{Doc: result, Old: oldValue}, nil
	}

	if !pathExists(doc, op.path) {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	parent, key, err := navigateToParent(doc, op.path)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	currentValue := getValueFromParent(parent, key)
	oldValue, ok := toFloat64(currentValue)
	if !ok {
		return internal.OpResult[any]{}, ErrNotNumber
	}
	result := oldValue + op.Inc

	if err := op.updateParent(parent, key, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// updateParent updates the parent container with the new value
func (op *IncOperation) updateParent(parent interface{}, key interface{}, value interface{}) error {
	return updateParent(parent, key, value)
}

// ToJSON serializes the operation to JSON format.
func (op *IncOperation) ToJSON() (internal.Operation, error) {
	// 遵循 float64 统一化原则：始终输出 float64
	return internal.Operation{
		"op":   string(internal.OpIncType),
		"path": formatPath(op.path),
		"inc":  op.Inc, // 统一的 float64 输出
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *IncOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpIncCode, op.path, op.Inc}, nil
}

// Validate validates the increment operation.
func (op *IncOperation) Validate() error {
	// Empty path (root level) is allowed for inc operations
	return nil
}

// Short aliases for common use
var (
	// NewInc creates a new inc operation
	NewInc = NewOpIncOperation
)
