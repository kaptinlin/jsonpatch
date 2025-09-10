package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpCopyOperation represents a copy operation that copies a value from one path to another.
type OpCopyOperation struct {
	BaseOp
	FromPath []string `json:"from"` // Source path
}

// CopyOperation is a non-stuttering alias for OpCopyOperation.
type CopyOperation = OpCopyOperation

// NewOpCopyOperation creates a new OpCopyOperation operation.
func NewOpCopyOperation(path, from []string) *OpCopyOperation {
	return &OpCopyOperation{
		BaseOp:   NewBaseOpWithFrom(path, from),
		FromPath: from,
	}
}

// Op returns the operation type.
func (o *OpCopyOperation) Op() internal.OpType {
	return internal.OpCopyType
}

// Code returns the operation code.
func (o *OpCopyOperation) Code() int {
	return internal.OpCopyCode
}

// From returns the source path.
func (o *OpCopyOperation) From() []string {
	return o.FromPath
}

// Apply applies the copy operation.
func (o *OpCopyOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, err := getValue(doc, o.FromPath)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Optimize: avoid unnecessary deep copy for simple types
	var clonedValue interface{}
	switch v := value.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, nil:
		clonedValue = v // Simple types are assigned directly
	default:
		// Deep copy for complex types
		clonedValue, err = DeepClone(v)
		if err != nil {
			return internal.OpResult[any]{}, err
		}
	}

	// Optimize: inline old value retrieval, reducing function calls
	var oldValue interface{}
	if len(o.Path()) > 0 {
		// Directly try to get value, more efficient than checking existence first
		if old, err := getValue(doc, o.Path()); err == nil {
			oldValue = old
		}
		// If the value is not found, oldValue remains nil, which is correct behavior
	}

	// Set value to target path
	err = insertValueAtPath(doc, o.Path(), clonedValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpCopyOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpCopyType),
		"path": formatPath(o.Path()),
		"from": formatPath(o.FromPath),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpCopyOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpCopyCode, o.Path(), o.FromPath}, nil
}

// Validate validates the copy operation.
func (o *OpCopyOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(o.FromPath) == 0 {
		return ErrFromPathEmpty
	}
	// Check that path and from are not the same
	if pathEquals(o.Path(), o.FromPath) {
		return ErrPathsIdentical
	}
	return nil
}

// Short aliases for common use
var (
	// NewCopy creates a new copy operation
	NewCopy = NewOpCopyOperation
)
