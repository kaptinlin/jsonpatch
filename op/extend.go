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
	path := ex.Path()
	target, err := value(doc, path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetObj, ok := target.(map[string]any)
	if !ok {
		return internal.OpResult[any]{}, ErrNotObject
	}

	extendedObj := extendObject(targetObj, ex.Properties, ex.DeleteNull)
	if len(path) == 0 {
		return internal.OpResult[any]{Doc: extendedObj}, nil
	}

	if err := setValueAtPath(doc, path, extendedObj); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc}, nil
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
	if ex.Properties == nil {
		return ErrPropertiesNil
	}
	return nil
}

func extendObject(obj map[string]any, props map[string]any, deleteNull bool) map[string]any {
	result := maps.Clone(obj)

	for k, v := range props {
		if k == "__proto__" {
			continue
		}

		if deleteNull && v == nil {
			delete(result, k)
		} else {
			result[k] = v
		}
	}

	return result
}
