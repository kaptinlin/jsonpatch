package op

import (
	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
)

// ReplaceOperation represents a replace operation that replaces a value at a specified path.
type ReplaceOperation struct {
	BaseOp
	Value    interface{} `json:"value"`              // New value
	OldValue interface{} `json:"oldValue,omitempty"` // The value that was replaced (optional)
}

// NewOpReplaceOperation creates a new OpReplaceOperation operation.
func NewOpReplaceOperation(path []string, value interface{}) *ReplaceOperation {
	return &ReplaceOperation{
		BaseOp:   NewBaseOp(path),
		Value:    value,
		OldValue: nil,
	}
}

// NewOpReplaceOperationWithOldValue creates a new OpReplaceOperation operation with oldValue.
func NewOpReplaceOperationWithOldValue(path []string, value interface{}, oldValue interface{}) *ReplaceOperation {
	return &ReplaceOperation{
		BaseOp:   NewBaseOp(path),
		Value:    value,
		OldValue: oldValue,
	}
}

// Op returns the operation type.
func (o *ReplaceOperation) Op() internal.OpType {
	return internal.OpReplaceType
}

// Code returns the operation code.
func (o *ReplaceOperation) Code() int {
	return internal.OpReplaceCode
}

// Apply applies the replace operation to the document.
func (o *ReplaceOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Clone the new value to prevent external mutations
	newValue := deepclone.Clone(o.Value)

	if len(o.path) == 0 {
		// Replace entire document
		oldValue := doc
		return internal.OpResult[any]{Doc: newValue, Old: oldValue}, nil
	}
	if len(o.path) == 1 && o.path[0] == "" {
		// Special case: path "/" refers to the key "" in the root object
		switch v := doc.(type) {
		case map[string]interface{}:
			oldValue := v[""]
			v[""] = newValue
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		default:
			return internal.OpResult[any]{}, ErrCannotReplace
		}
	}

	// Optimize: directly check type and get value in type switch
	parent, key, err := navigateToParent(doc, o.path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Optimize: directly check type and get value in type switch
	switch p := parent.(type) {
	case map[string]interface{}:
		k, ok := key.(string)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
		// Optimize: merge existence check and value retrieval
		if oldValue, exists := p[k]; exists {
			p[k] = newValue
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		}
		return internal.OpResult[any]{}, ErrPathNotFound

	case []interface{}:
		k, ok := key.(int)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeSlice
		}
		// Optimize: merge boundary check and value retrieval
		if k >= 0 && k < len(p) {
			oldValue := p[k]
			p[k] = newValue
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		}
		return internal.OpResult[any]{}, ErrPathNotFound

	default:
		return internal.OpResult[any]{}, ErrUnsupportedParentType
	}
}

// ToJSON serializes the operation to JSON format.
func (o *ReplaceOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpReplaceType),
		Path:  formatPath(o.path),
		Value: o.Value,
	}

	if o.OldValue != nil {
		result.OldValue = o.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *ReplaceOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpReplaceCode, o.path, o.Value}, nil
}

// Validate validates the replace operation.
func (o *ReplaceOperation) Validate() error {
	if len(o.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}

// Short aliases for common use
var (
	// NewReplace creates a new replace operation
	NewReplace = NewOpReplaceOperation
	// NewReplaceWithOldValue creates a new replace operation with old value
	NewReplaceWithOldValue = NewOpReplaceOperationWithOldValue
)
