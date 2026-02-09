package op

import (
	"maps"

	"github.com/kaptinlin/jsonpatch/internal"
)

// ExtendOperation represents an object extend operation.
// path: target path
// props: properties to add/update
// deleteNull: whether to delete properties with null values
// Only supports object type fields.
type ExtendOperation struct {
	BaseOp
	Properties map[string]any `json:"props"`      // Properties to add
	DeleteNull bool           `json:"deleteNull"` // Whether to delete null properties
}

// NewExtend creates a new object extend operation.
func NewExtend(path []string, properties map[string]any, deleteNull bool) *ExtendOperation {
	return &ExtendOperation{
		BaseOp:     NewBaseOp(path),
		Properties: properties,
		DeleteNull: deleteNull,
	}
}

// Op returns the operation type.
func (ex *ExtendOperation) Op() internal.OpType {
	return internal.OpExtendType
}

// Code returns the operation code.
func (ex *ExtendOperation) Code() int {
	return internal.OpExtendCode
}

// Apply applies the object extend operation.
func (ex *ExtendOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level extend specially
	if len(ex.Path()) == 0 {
		// Target is the root document
		targetObj, ok := doc.(map[string]any)
		if !ok {
			return internal.OpResult[any]{}, ErrNotObject
		}

		// Create a copy and extend it
		original := maps.Clone(targetObj)

		// Use objExtend to properly handle the extension with deleteNull
		extendedObj := objExtend(targetObj, ex.Properties, ex.DeleteNull)

		return internal.OpResult[any]{Doc: extendedObj, Old: original}, nil
	}

	// Get the target object
	target, err := getValue(doc, ex.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Check if target is an object
	targetObj, ok := target.(map[string]any)
	if !ok {
		return internal.OpResult[any]{}, ErrNotObject
	}

	// Use objExtend to properly handle the extension with deleteNull
	extendedObj := objExtend(targetObj, ex.Properties, ex.DeleteNull)

	// Set the extended object back
	err = setValueAtPath(doc, ex.Path(), extendedObj)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// ToJSON serializes the operation to JSON format.
func (ex *ExtendOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:         string(internal.OpExtendType),
		Path:       formatPath(ex.Path()),
		Props:      ex.Properties,
		DeleteNull: ex.DeleteNull,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (ex *ExtendOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpExtendCode, ex.Path(), ex.Properties}, nil
}

// Validate validates the extend operation.
func (ex *ExtendOperation) Validate() error {
	// Empty path is valid for extend operation (root level)
	if ex.Properties == nil {
		return ErrPropertiesNil
	}
	return nil
}

// objExtend extends object obj with properties from props.
// If deleteNull is true, properties with null values are deleted.
func objExtend(obj map[string]any, props map[string]any, deleteNull bool) map[string]any {
	// Create a copy of the original object
	result := maps.Clone(obj)

	// Add/update properties from props
	for k, v := range props {
		// Security check: prevent __proto__ pollution
		if k == "__proto__" {
			continue // Skip __proto__ keys for security
		}

		if deleteNull && v == nil {
			delete(result, k)
		} else {
			result[k] = v
		}
	}

	return result
}
