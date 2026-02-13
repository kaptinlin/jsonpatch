package op

import "github.com/kaptinlin/jsonpatch/internal"

// CopyOperation represents a copy operation that copies a value from one path to another.
type CopyOperation struct {
	BaseOp
}

// NewCopy creates a new copy operation.
func NewCopy(path, from []string) *CopyOperation {
	return &CopyOperation{
		BaseOp: NewBaseOpWithFrom(path, from),
	}
}

// Op returns the operation type.
func (c *CopyOperation) Op() internal.OpType {
	return internal.OpCopyType
}

// Code returns the operation code.
func (c *CopyOperation) Code() int {
	return internal.OpCopyCode
}

// Apply applies the copy operation.
func (c *CopyOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, err := getValue(doc, c.from)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Avoid deep copy for simple types
	var clonedValue any
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

	// Handle empty path (root replacement)
	if len(c.path) == 0 {
		return internal.OpResult[any]{Doc: clonedValue, Old: doc}, nil
	}

	var oldValue any
	if old, err := getValue(doc, c.path); err == nil {
		oldValue = old
	}

	err = insertValueAtPath(doc, c.path, clonedValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (c *CopyOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpCopyType),
		Path: formatPath(c.path),
		From: formatPath(c.from),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (c *CopyOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpCopyCode, c.path, c.from}, nil
}

// Validate validates the copy operation.
func (c *CopyOperation) Validate() error {
	if len(c.from) == 0 {
		return ErrFromPathEmpty
	}
	if pathEquals(c.path, c.from) {
		return ErrPathsIdentical
	}
	return nil
}
