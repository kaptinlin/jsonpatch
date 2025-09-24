// Package op provides implementations for JSON Patch operations.
package op

import (
	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
)

// AddOperation represents an add operation that adds a value at a specified path.
type AddOperation struct {
	BaseOp
	Value interface{} `json:"value"` // Value to add
}

// NewAdd creates a new AddOperation operation.
func NewAdd(path []string, value interface{}) *AddOperation {
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
func addAtPath(doc interface{}, path []string, value interface{}) (interface{}, interface{}, error) {
	if len(path) == 1 {
		key := path[0]
		switch v := doc.(type) {
		case map[string]interface{}:
			oldValue := v[key]
			v[key] = value
			return doc, oldValue, nil
		case []interface{}:
			if key == "-" {
				v = append(v, value)
				return v, nil, nil
			}
			index, err := parseArrayIndex(key)
			if err != nil {
				return nil, nil, err
			}
			if index < 0 || index > len(v) {
				return nil, nil, ErrArrayIndexOutOfBounds
			}

			// Get the displaced element (if any)
			var displacedElement interface{}
			if index < len(v) {
				displacedElement = v[index]
			}

			v = append(v[:index], append([]interface{}{value}, v[index:]...)...)
			return v, displacedElement, nil
		default:
			return nil, nil, ErrCannotAddToValue
		}
	}
	// Recursive case
	key := path[0]
	rest := path[1:]
	switch v := doc.(type) {
	case map[string]interface{}:
		child, exists := v[key]
		if !exists {
			// According to JSON Patch spec, missing objects are not created recursively
			return nil, nil, ErrPathMissingRecursive
		}
		newChild, oldValue, err := addAtPath(child, rest, value)
		if err != nil {
			return nil, nil, err
		}
		v[key] = newChild
		return doc, oldValue, nil
	case []interface{}:
		index, err := parseArrayIndex(key)
		if err != nil {
			return nil, nil, err
		}
		if index < 0 || index >= len(v) {
			return nil, nil, ErrArrayIndexOutOfBounds
		}
		newChild, oldValue, err := addAtPath(v[index], rest, value)
		if err != nil {
			return nil, nil, err
		}
		v[index] = newChild
		return v, oldValue, nil
	default:
		return nil, nil, ErrCannotAddToValue
	}
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
