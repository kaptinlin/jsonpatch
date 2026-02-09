package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// RemoveOperation represents a remove operation that removes a value at a specified path.
type RemoveOperation struct {
	BaseOp
	OldValue    any `json:"oldValue,omitempty"` // The value that was removed (optional)
	HasOldValue bool        // Whether oldValue is explicitly set
}

// NewRemove creates a new remove operation.
func NewRemove(path []string) *RemoveOperation {
	return &RemoveOperation{
		BaseOp:      NewBaseOp(path),
		OldValue:    nil,
		HasOldValue: false,
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
func (o *RemoveOperation) Op() internal.OpType {
	return internal.OpRemoveType
}

// Code returns the operation code.
func (o *RemoveOperation) Code() int {
	return internal.OpRemoveCode
}

// Path returns the operation path.
func (o *RemoveOperation) Path() []string {
	return o.path
}

// Apply applies the remove operation to the document.
func (o *RemoveOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(o.path) == 0 {
		return internal.OpResult[any]{}, ErrPathEmpty
	}
	if len(o.path) == 1 {
		switch v := doc.(type) {
		case map[string]any:
			oldValue, exists := v[o.path[0]]
			if !exists {
				return internal.OpResult[any]{}, ErrPathNotFound
			}
			delete(v, o.path[0])
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		case []any:
			index, err := parseArrayIndex(o.path[0])
			if err != nil {
				return internal.OpResult[any]{}, err
			}
			if index < 0 || index >= len(v) {
				return internal.OpResult[any]{}, ErrIndexOutOfRange
			}
			oldValue := v[index]
			// Optimize: pre-allocate exact size and use copy
			newArray := make([]any, len(v)-1)
			copy(newArray, v[:index])
			copy(newArray[index:], v[index+1:])
			return internal.OpResult[any]{Doc: newArray, Old: oldValue}, nil
		default:
			return internal.OpResult[any]{}, ErrCannotRemoveFromValue
		}
	}
	// Not root path, recursively delete
	parent, key, err := navigateToParent(doc, o.path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	oldValue := getValueFromParent(parent, key)

	// Check if the path actually exists
	if oldValue == nil {
		switch p := parent.(type) {
		case map[string]any:
			if k, ok := key.(string); ok {
				if _, exists := p[k]; !exists {
					return internal.OpResult[any]{}, ErrPathNotFound
				}
			}
		case []any:
			if k, ok := key.(int); ok {
				if k < 0 || k >= len(p) {
					return internal.OpResult[any]{}, ErrPathNotFound
				}
			}
		}
	}
	switch p := parent.(type) {
	case map[string]any:
		if k, ok := key.(string); ok {
			delete(p, k)
		} else {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
	case []any:
		if k, ok := key.(int); ok {
			if k < 0 || k >= len(p) {
				return internal.OpResult[any]{}, ErrIndexOutOfRange
			}
			// Optimize: pre-allocate exact size and use copy
			newSlice := make([]any, len(p)-1)
			copy(newSlice, p[:k])
			copy(newSlice[k:], p[k+1:])
			if err := setValueAtPath(doc, o.path[:len(o.path)-1], newSlice); err != nil {
				return internal.OpResult[any]{}, err
			}
		} else {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeSlice
		}
	default:
		return internal.OpResult[any]{}, ErrUnsupportedParentType
	}
	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *RemoveOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpRemoveType),
		Path: formatPath(o.path),
	}

	if o.HasOldValue {
		result.OldValue = o.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *RemoveOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpRemoveCode, o.path}, nil
}

// Validate validates the remove operation.
func (o *RemoveOperation) Validate() error {
	if len(o.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}

