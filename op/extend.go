package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// ExtendOperation represents an object extend operation.
// path: target path
// props: properties to add/update
// deleteNull: whether to delete properties with null values
// Only supports object type fields.
type ExtendOperation struct {
	BaseOp
	Properties map[string]interface{} `json:"props"`      // Properties to add
	DeleteNull bool                   `json:"deleteNull"` // Whether to delete null properties
}

// OpExtendOperation is a backward-compatible alias for ExtendOperation.
type OpExtendOperation = ExtendOperation

// NewOpExtendOperation creates a new object extend operation.
func NewOpExtendOperation(path []string, properties map[string]interface{}, deleteNull bool) *ExtendOperation {
	return &ExtendOperation{
		BaseOp:     NewBaseOp(path),
		Properties: properties,
		DeleteNull: deleteNull,
	}
}

// Op returns the operation type.
func (op *ExtendOperation) Op() internal.OpType {
	return internal.OpExtendType
}

// Code returns the operation code.
func (op *ExtendOperation) Code() int {
	return internal.OpExtendCode
}

// Apply applies the object extend operation.
func (op *ExtendOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level extend specially
	if len(op.Path()) == 0 {
		// Target is the root document
		targetObj, ok := doc.(map[string]interface{})
		if !ok {
			return internal.OpResult[any]{}, ErrNotObject
		}

		// Create a copy and extend it
		original := make(map[string]interface{})
		for k, v := range targetObj {
			original[k] = v
		}

		// Use objExtend to properly handle the extension with deleteNull
		extendedObj := objExtend(targetObj, op.Properties, op.DeleteNull)

		return internal.OpResult[any]{Doc: extendedObj, Old: original}, nil
	}

	// Get the target object
	target, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Check if target is an object
	targetObj, ok := target.(map[string]interface{})
	if !ok {
		return internal.OpResult[any]{}, ErrNotObject
	}

	// Use objExtend to properly handle the extension with deleteNull
	extendedObj := objExtend(targetObj, op.Properties, op.DeleteNull)

	// Set the extended object back
	err = setValueAtPath(doc, op.Path(), extendedObj)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// ToJSON serializes the operation to JSON format.
func (op *ExtendOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":    string(internal.OpExtendType),
		"path":  formatPath(op.Path()),
		"props": op.Properties,
	}
	if op.DeleteNull {
		result["deleteNull"] = op.DeleteNull
	}
	return result, nil
}

// ToCompact serializes the operation to compact format.
func (op *ExtendOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpExtendCode, op.Path(), op.Properties}, nil
}

// Validate validates the extend operation.
func (op *ExtendOperation) Validate() error {
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

// Short aliases for common use
var (
	// NewExtend creates a new extend operation
	NewExtend = NewOpExtendOperation
)
