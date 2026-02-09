package op

import (
	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
)

// ReplaceOperation represents a replace operation that replaces a value at a specified path.
type ReplaceOperation struct {
	BaseOp
	Value    any `json:"value"`              // New value
	OldValue any `json:"oldValue,omitempty"` // The value that was replaced (optional)
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
		BaseOp:   NewBaseOp(path),
		Value:    value,
		OldValue: oldValue,
	}
}

// Op returns the operation type.
func (rp *ReplaceOperation) Op() internal.OpType {
	return internal.OpReplaceType
}

// Code returns the operation code.
func (rp *ReplaceOperation) Code() int {
	return internal.OpReplaceCode
}

// Apply applies the replace operation to the document.
func (rp *ReplaceOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Clone the new value to prevent external mutations
	newValue := deepclone.Clone(rp.Value)

	if len(rp.path) == 0 {
		// Replace entire document
		oldValue := doc
		return internal.OpResult[any]{Doc: newValue, Old: oldValue}, nil
	}
	if len(rp.path) == 1 && rp.path[0] == "" {
		// Special case: path "/" refers to the key "" in the root object
		switch v := doc.(type) {
		case map[string]any:
			oldValue := v[""]
			v[""] = newValue
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		default:
			return internal.OpResult[any]{}, ErrCannotReplace
		}
	}

	parent, key, err := navigateToParent(doc, rp.path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	switch p := parent.(type) {
	case map[string]any:
		k, ok := key.(string)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeMap
		}
		if oldValue, exists := p[k]; exists {
			p[k] = newValue
			return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
		}
		return internal.OpResult[any]{}, ErrPathNotFound

	case []any:
		k, ok := key.(int)
		if !ok {
			return internal.OpResult[any]{}, ErrInvalidKeyTypeSlice
		}
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
func (rp *ReplaceOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:    string(internal.OpReplaceType),
		Path:  formatPath(rp.path),
		Value: rp.Value,
	}

	if rp.OldValue != nil {
		result.OldValue = rp.OldValue
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (rp *ReplaceOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpReplaceCode, rp.path, rp.Value}, nil
}

// Validate validates the replace operation.
func (rp *ReplaceOperation) Validate() error {
	if len(rp.path) == 0 {
		return ErrPathEmpty
	}
	return nil
}
