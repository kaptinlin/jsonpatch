package op

import "github.com/kaptinlin/jsonpatch/internal"

// RemoveOperation represents a remove operation that removes a value at a specified path.
type RemoveOperation struct {
	BaseOp
	OldValue    any  `json:"oldValue,omitempty"` // The value that was removed (optional)
	HasOldValue bool // Whether oldValue is explicitly set
}

// NewRemove creates a new remove operation.
func NewRemove(path []string) *RemoveOperation {
	return &RemoveOperation{
		BaseOp: NewBaseOp(path),
	}
}

// NewRemoveWithOldValue creates a new remove operation with oldValue.
func NewRemoveWithOldValue(path []string, oldValue any) *RemoveOperation {
	return &RemoveOperation{
		BaseOp:      NewBaseOp(path),
		OldValue:    oldValue,
		HasOldValue: true,
	}
}

// Op returns the operation type.
func (r *RemoveOperation) Op() internal.OpType {
	return internal.OpRemoveType
}

// Apply applies the remove operation to the document.
func (r *RemoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	newDoc, oldValue, err := removeAtPath(doc, r.path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	return internal.OpResult[any]{Doc: newDoc, Old: oldValue}, nil
}

// Validate validates the remove operation.
func (r *RemoveOperation) Validate() error {
	// Empty path is valid root removal.
	return nil
}
