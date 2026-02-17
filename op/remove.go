package op

import (
	"slices"

	"github.com/kaptinlin/jsonpatch/internal"
)

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

// Code returns the operation code.
func (r *RemoveOperation) Code() int {
	return internal.OpRemoveCode
}

// Apply applies the remove operation to the document.
func (r *RemoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(r.path) == 0 {
		return internal.OpResult[any]{}, ErrPathEmpty
	}
	if len(r.path) == 1 {
		switch v := doc.(type) {
		case map[string]any:
			oldValue, exists := v[r.path[0]]
			if !exists {
				return internal.OpResult[any]{}, ErrPathNotFound
			}
			delete(v, r.path[0])
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		case []any:
			index, err := parseArrayIndex(r.path[0])
			if err != nil {
				return internal.OpResult[any]{}, err
			}
			if index < 0 || index >= len(v) {
				return internal.OpResult[any]{}, ErrIndexOutOfRange
			}
			oldValue := v[index]
			// Use slices.Delete (Go 1.21+) for efficient element removal
			newArray := slices.Delete(slices.Clone(v), index, index+1)
			return internal.OpResult[any]{Doc: newArray, Old: oldValue}, nil
		default:
			return internal.OpResult[any]{}, ErrCannotRemoveFromValue
		}
	}
	// Not root path, recursively delete
	parent, key, err := navigateToParent(doc, r.path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	switch p := parent.(type) {
	case map[string]any:
		k, ok := key.(string)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
		oldValue, exists := p[k]
		if !exists {
			return internal.OpResult[any]{}, ErrPathNotFound
		}
		delete(p, k)
		return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
	case []any:
		k, ok := key.(int)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeSlice
		}
		if k < 0 || k >= len(p) {
			return internal.OpResult[any]{}, ErrIndexOutOfRange
		}
		oldValue := p[k]
		// Use slices.Delete (Go 1.21+) for efficient element removal
		newSlice := slices.Delete(slices.Clone(p), k, k+1)
		if err := setValueAtPath(doc, r.path[:len(r.path)-1], newSlice); err != nil {
			return internal.OpResult[any]{}, err
		}
		return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
	default:
		return internal.OpResult[any]{}, ErrUnsupportedParentType
	}
}

// ToJSON serializes the operation to JSON format.
func (r *RemoveOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpRemoveType),
		Path: formatPath(r.path),
	}

	if r.HasOldValue {
		result.OldValue = r.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (r *RemoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpRemoveCode, r.path}, nil
}

// Validate validates the remove operation.
func (r *RemoveOperation) Validate() error {
	if len(r.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}
