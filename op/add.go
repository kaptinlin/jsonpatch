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

// NewAdd creates a new add operation.
func NewAdd(path []string, value any) *AddOperation {
	return &AddOperation{
		BaseOp: NewBaseOp(path),
		Value:  value,
	}
}

// Op returns the operation type.
func (a *AddOperation) Op() internal.OpType {
	return internal.OpAddType
}

// Apply applies the add operation.
func (a *AddOperation) Apply(doc any) (internal.OpResult[any], error) {
	newValue := deepclone.Clone(a.Value)

	newDoc, oldValue, err := addAtPath(doc, a.path, newValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}
	return internal.OpResult[any]{Doc: newDoc, Old: oldValue}, nil
}

// Validate validates the add operation.
func (a *AddOperation) Validate() error {
	// Empty path is valid root replacement.
	// Value field is not validated here as it can be any value including nil.
	// The value field presence is validated at the JSON parsing level.
	return nil
}
