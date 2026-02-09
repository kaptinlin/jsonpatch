// Package op provides implementations for JSON Patch operations.
package op

import (
	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
)

// AddOperation represents an add operation that adds a value at a specified path.
type AddOperation struct {
	BaseOp
	Value any `json:"value"` // Value to add
}

// NewAdd creates a new AddOperation operation.
func NewAdd(path []string, value any) *AddOperation {
	return &AddOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// Op returns the operation type.
func (o *AddOperation) Op() internal.OpType {
	return internal.OpAddType
}

// Code returns the operation code.
func (o *AddOperation) Code() int {
	return internal.OpAddCode
}

// Path returns the operation path.
func (o *AddOperation) Path() []string {
	return o.path
}

// Apply applies the add operation.
func (o *AddOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Clone the new value to prevent external mutations
	newValue := deepclone.Clone(o.Value)

	// Handle empty path (root replacement) - only for truly empty path, not empty string key
	if len(o.path) == 0 {
		// Replace entire document
		return internal.OpResult[any]{Doc: newValue, Old: doc}, nil
	}

	newDoc, oldValue, err := addAtPath(doc, o.path, newValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	return internal.OpResult[any]{Doc: newDoc, Old: oldValue}, nil
}

// addAtPath recursively inserts value at the given path, returns new doc and old value if replaced.
func addAtPath(doc any, path []string, value any) (any, any, error) {
	switch v := doc.(type) {
	case map[string]any:
		return addToMap(v, path, value)
	case []any:
		return addToSlice(v, path, value)
	default:
		return nil, nil, ErrCannotAddToValue
	}
}

func addToMap(doc map[string]any, path []string, value any) (any, any, error) {
	key := path[0]

	if len(path) == 1 {
		oldValue := doc[key]
		doc[key] = value
		return doc, oldValue, nil
	}

	// Recursive case
	child, exists := doc[key]
	if !exists {
		// According to JSON Patch spec, missing objects are not created recursively
		return nil, nil, ErrCannotReplace
	}
	newChild, oldValue, err := addAtPath(child, path[1:], value)
	if err != nil {
		return nil, nil, err
	}
	doc[key] = newChild
	return doc, oldValue, nil
}

func addToSlice(doc []any, path []string, value any) (any, any, error) {
	key := path[0]

	if len(path) == 1 {
		if key == "-" {
			doc = append(doc, value)
			return doc, nil, nil
		}
		index, err := parseArrayIndex(key)
		if err != nil {
			return nil, nil, err
		}
		if index < 0 || index > len(doc) {
			return nil, nil, ErrIndexOutOfRange
		}

		// Get the displaced element (if any)
		var displacedElement any
		if index < len(doc) {
			displacedElement = doc[index]
		}

		// Optimize: pre-allocate correct size and use copy to avoid double allocation
		newV := make([]any, len(doc)+1)
		copy(newV, doc[:index])
		newV[index] = value
		copy(newV[index+1:], doc[index:])
		return newV, displacedElement, nil
	}

	// Recursive case
	index, err := parseArrayIndex(key)
	if err != nil {
		return nil, nil, err
	}
	if index < 0 || index >= len(doc) {
		return nil, nil, ErrIndexOutOfRange
	}
	newChild, oldValue, err := addAtPath(doc[index], path[1:], value)
	if err != nil {
		return nil, nil, err
	}
	doc[index] = newChild
	return doc, oldValue, nil
}

// ToJSON serializes the operation to JSON format.
func (o *AddOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:    string(internal.OpAddType),
		Path:  formatPath(o.path),
		Value: o.Value,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *AddOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpAddCode, o.path, o.Value}, nil
}

// Validate validates the add operation.
func (o *AddOperation) Validate() error {
	if len(o.path) == 0 {
		return ErrPathEmpty
	}
	// Note: value field is not validated here as it can be any value including nil
	// The value field presence is validated at the JSON parsing level
	return nil
}
