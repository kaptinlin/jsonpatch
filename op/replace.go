package op

import (
	"github.com/kaptinlin/deepclone"

	"github.com/kaptinlin/jsonpatch/internal"
)

// ReplaceOperation represents a replace operation that replaces a value at a specified path.
type ReplaceOperation struct {
	BaseOp
	Value       any  `json:"value"`              // New value
	OldValue    any  `json:"oldValue,omitempty"` // The value that was replaced (optional)
	HasOldValue bool // Whether oldValue is explicitly set
}

// NewReplace creates a new replace operation.
func NewReplace(path []string, value any) *ReplaceOperation {
	return &ReplaceOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// NewReplaceWithOldValue creates a new replace operation with oldValue.
func NewReplaceWithOldValue(path []string, value any, oldValue any) *ReplaceOperation {
	return &ReplaceOperation{
		BaseOp:      NewBaseOp(path),
		Value:       value,
		OldValue:    oldValue,
		HasOldValue: true,
	}
}

// Op returns the operation type.
func (rp *ReplaceOperation) Op() internal.OpType {
	return internal.OpReplaceType
}

// Apply applies the replace operation to the document.
func (rp *ReplaceOperation) Apply(doc any) (internal.OpResult[any], error) {
	newValue := deepclone.Clone(rp.Value)

	newDoc, oldValue, err := replaceAtPath(doc, rp.path, newValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	return internal.OpResult[any]{Doc: newDoc, Old: oldValue}, nil
}

// Validate validates the replace operation.
func (rp *ReplaceOperation) Validate() error {
	// Empty path is valid root replacement.
	return nil
}
