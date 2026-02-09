package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// CopyOperation represents a copy operation that copies a value from one path to another.
type CopyOperation struct {
	BaseOp
	FromPath []string `json:"from"` // Source path
}

// NewCopy creates a new copy operation.
func NewCopy(path, from []string) *CopyOperation {
	return &CopyOperation{
		BaseOp:   NewBaseOpWithFrom(path, from),
		FromPath: from,
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

// From returns the source path.
func (c *CopyOperation) From() []string {
	return c.FromPath
}

// Apply applies the copy operation.
func (c *CopyOperation) Apply(doc any) (internal.OpResult[any], error) {
	value, err := getValue(doc, c.FromPath)
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
	if len(c.Path()) == 0 {
		// Copy to root - replace entire document
		return internal.OpResult[any]{Doc: clonedValue, Old: doc}, nil
	}

	var oldValue any
	if old, err := getValue(doc, c.Path()); err == nil {
		oldValue = old
	}
	// If the value is not found, oldValue remains nil, which is correct behavior

	// Set value to target path
	err = insertValueAtPath(doc, c.Path(), clonedValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (c *CopyOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpCopyType),
		Path: formatPath(c.Path()),
		From: formatPath(c.FromPath),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (c *CopyOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpCopyCode, c.Path(), c.FromPath}, nil
}

// Validate validates the copy operation.
func (c *CopyOperation) Validate() error {
	// Empty path is valid for copy (copies to root)
	// Only from path cannot be empty
	if len(c.FromPath) == 0 {
		return ErrFromPathEmpty
	}
	// Check that path and from are not the same
	if pathEquals(c.Path(), c.FromPath) {
		return ErrPathsIdentical
	}
	return nil
}
