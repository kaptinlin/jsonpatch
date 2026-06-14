package op

import (
	"slices"

	"github.com/kaptinlin/deepclone"

	"github.com/kaptinlin/jsonpatch/internal"
)

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

// Apply applies the copy operation.
func (c *CopyOperation) Apply(doc any) (internal.OpResult[any], error) {
	val, err := value(doc, c.from)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Deep clone the value uniformly to ensure immutability
	clonedValue := deepclone.Clone(val)

	// Handle empty path (root replacement)
	if len(c.path) == 0 {
		return internal.OpResult[any]{Doc: clonedValue, Old: doc}, nil
	}

	newDoc, oldValue, err := addAtPath(doc, c.path, clonedValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: newDoc, Old: oldValue}, nil
}

// Validate validates the copy operation.
func (c *CopyOperation) Validate() error {
	if slices.Equal(c.path, c.from) {
		return ErrPathsIdentical
	}
	return nil
}
