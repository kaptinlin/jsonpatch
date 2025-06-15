package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpExtendOperation represents an object extend operation.
// path: target path
// props: properties to add/update
// deleteNull: whether to delete properties with null values
// Only supports object type fields.
type OpExtendOperation struct {
	BaseOp
	Properties map[string]interface{} `json:"props"`      // Properties to add
	DeleteNull bool                   `json:"deleteNull"` // Whether to delete null properties
}

// NewOpExtendOperation creates a new object extend operation.
func NewOpExtendOperation(path []string, properties map[string]interface{}, deleteNull bool) *OpExtendOperation {
	return &OpExtendOperation{
		BaseOp:     NewBaseOp(path),
		Properties: properties,
		DeleteNull: deleteNull,
	}
}

// Op returns the operation type.
func (op *OpExtendOperation) Op() internal.OpType {
	return internal.OpExtendType
}

// Code returns the operation code.
func (op *OpExtendOperation) Code() int {
	return internal.OpExtendCode
}

// Apply applies the object extend operation.
func (o *OpExtendOperation) Apply(doc any) (internal.OpResult, error) {
	// Handle root level extend specially
	if len(o.Path()) == 0 {
		// Target is the root document
		targetObj, ok := doc.(map[string]interface{})
		if !ok {
			return internal.OpResult{}, ErrNotObject
		}

		// Create a copy and extend it
		original := make(map[string]interface{})
		for k, v := range targetObj {
			original[k] = v
		}

		// Use objExtend to properly handle the extension with deleteNull
		extendedObj := objExtend(targetObj, o.Properties, o.DeleteNull)

		return internal.OpResult{Doc: extendedObj, Old: original}, nil
	}

	// Get the target object
	target, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult{}, err
	}

	// Check if target is an object
	targetObj, ok := target.(map[string]interface{})
	if !ok {
		return internal.OpResult{}, ErrNotObject
	}

	// Use objExtend to properly handle the extension with deleteNull
	extendedObj := objExtend(targetObj, o.Properties, o.DeleteNull)

	// Set the extended object back
	err = setValueAtPath(doc, o.Path(), extendedObj)
	if err != nil {
		return internal.OpResult{}, err
	}

	return internal.OpResult{Doc: doc, Old: target}, nil
}

// ToJSON serializes the operation to JSON format.
func (o *OpExtendOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":    string(internal.OpExtendType),
		"path":  formatPath(o.Path()),
		"props": o.Properties,
	}
	if o.DeleteNull {
		result["deleteNull"] = o.DeleteNull
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpExtendOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpExtendCode, o.Path(), o.Properties}, nil
}

// Validate validates the extend operation.
func (op *OpExtendOperation) Validate() error {
	// Empty path is valid for extend operation (root level)
	if op.Properties == nil {
		return ErrPropertiesNil
	}
	return nil
}

// objExtend extends object obj with properties from props.
// If deleteNull is true, properties with null values are deleted.
func objExtend(obj map[string]interface{}, props map[string]interface{}, deleteNull bool) map[string]interface{} {
	// Create a copy of the original object
	result := make(map[string]interface{})
	for k, v := range obj {
		result[k] = v
	}

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
