package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// RemoveOperation represents a remove operation that removes a value at a specified path.
type RemoveOperation struct {
	BaseOp
	OldValue    interface{} `json:"oldValue,omitempty"` // The value that was removed (optional)
	HasOldValue bool        // Whether oldValue is explicitly set
}

// NewOpRemoveOperation creates a new OpRemoveOperation operation.
func NewOpRemoveOperation(path []string) *RemoveOperation {
	return &RemoveOperation{
		BaseOp:      NewBaseOp(path),
		OldValue:    nil,
		HasOldValue: false,
	}
}

// NewOpRemoveOperationWithOldValue creates a new OpRemoveOperation operation with oldValue.
func NewOpRemoveOperationWithOldValue(path []string, oldValue interface{}) *RemoveOperation {
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
		case map[string]interface{}:
			oldValue, exists := v[o.path[0]]
			if !exists {
				return internal.OpResult[any]{}, ErrPathNotFound
			}
			delete(v, o.path[0])
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		case []interface{}:
			index, err := parseArrayIndex(o.path[0])
			if err != nil {
				return internal.OpResult[any]{}, err
			}
			if index < 0 || index >= len(v) {
				return internal.OpResult[any]{}, ErrArrayIndexOutOfBounds
			}
			oldValue := v[index]
			// Create new array without the removed element
			newArray := make([]interface{}, 0, len(v)-1)
			newArray = append(newArray, v[:index]...)
			newArray = append(newArray, v[index+1:]...)
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
		case map[string]interface{}:
			if k, ok := key.(string); ok {
				if _, exists := p[k]; !exists {
					return internal.OpResult[any]{}, ErrPathNotFound
				}
			}
		case []interface{}:
			if k, ok := key.(int); ok {
				if k < 0 || k >= len(p) {
					return internal.OpResult[any]{}, ErrPathNotFound
				}
			}
		}
	}
	switch p := parent.(type) {
	case map[string]interface{}:
		if k, ok := key.(string); ok {
			delete(p, k)
		} else {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
	case []interface{}:
		if k, ok := key.(int); ok {
			if k < 0 || k >= len(p) {
				return internal.OpResult[any]{}, ErrIndexOutOfRange
			}
			// Create new slice without the removed element
			newSlice := make([]interface{}, 0, len(p)-1)
			newSlice = append(newSlice, p[:k]...)
			newSlice = append(newSlice, p[k+1:]...)
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
		"op":   string(internal.OpRemoveType),
		"path": formatPath(o.path),
	}

	if o.HasOldValue {
		result["oldValue"] = o.OldValue
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

// Short aliases for common use
var (
	// NewRemove creates a new remove operation
	NewRemove = NewOpRemoveOperation
	// NewRemoveWithOldValue creates a new remove operation with old value
	NewRemoveWithOldValue = NewOpRemoveOperationWithOldValue
)
