package op

import (
	"reflect"

	"github.com/kaptinlin/jsonpatch/internal"
)

// FlipOperation represents a flip operation that inverts boolean values
// or converts other types to boolean and then inverts them
type FlipOperation struct {
	BaseOp
}

// NewFlip creates a new flip operation.
func NewFlip(path []string) *FlipOperation {
	return &FlipOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type
func (o *FlipOperation) Op() internal.OpType {
	return internal.OpFlipType
}

// Code returns the operation code
func (o *FlipOperation) Code() int {
	return internal.OpFlipCode
}

// Apply applies the flip operation to the document
func (o *FlipOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level flip specially
	if len(o.Path()) == 0 {
		// Target is the root document
		flippedValue := o.flipValue(doc)
		return internal.OpResult[any]{Doc: flippedValue, Old: doc}, nil
	}

	// Get the value to flip, or undefined if it doesn't exist
	value, err := getValue(doc, o.Path())
	var oldValue any
	if err != nil {
		// If path doesn't exist, treat as undefined (which becomes false, then flips to true)
		value = nil
		oldValue = nil
	} else {
		oldValue = value
	}

	// Flip the value
	flippedValue := o.flipValue(value)

	// Set the flipped value back
	err = setValueAtPath(doc, o.Path(), flippedValue)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// flipValue flips boolean values or converts other types to boolean and flip
func (o *FlipOperation) flipValue(value any) any {
	switch v := value.(type) {
	case bool:
		return !v
	default:
		// Convert to boolean first, then flip
		return !o.toBool(value)
	}
}

// toBool converts values to boolean for flip operation
func (o *FlipOperation) toBool(value any) bool {
	switch v := value.(type) {
	case nil:
		return false
	case bool:
		return v
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case float32:
		return v != 0.0
	case float64:
		return v != 0.0
	case string:
		return v != ""
	default:
		// For complex types, arrays and objects are always truthy in JavaScript
		val := reflect.ValueOf(value)
		switch val.Kind() {
		case reflect.Array, reflect.Slice:
			// Arrays are always truthy in JavaScript, even if empty
			return true
		case reflect.Map:
			// Objects are always truthy in JavaScript, even if empty
			return true
		case reflect.Chan:
			return val.Len() > 0
		case reflect.Pointer, reflect.Interface:
			return !val.IsNil()
		case reflect.Invalid:
			return false
		case reflect.Bool:
			return val.Bool()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return val.Int() != 0
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return val.Uint() != 0
		case reflect.Float32, reflect.Float64:
			return val.Float() != 0.0
		case reflect.Complex64, reflect.Complex128:
			return val.Complex() != 0
		case reflect.String:
			return val.String() != ""
		case reflect.Struct:
			// Structs are always considered truthy unless they're nil pointers
			return true
		case reflect.Func:
			return !val.IsNil()
		case reflect.UnsafePointer:
			return !val.IsNil()
		default:
			return true // Other types are considered truthy
		}
	}
}

// ToJSON serializes the operation to JSON format.
func (o *FlipOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpFlipType),
		Path: formatPath(o.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *FlipOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpFlipCode, o.Path()}, nil
}

// Validate validates the flip operation.
func (o *FlipOperation) Validate() error {
	// Empty path is valid for flip operation (root level)
	return nil
}

