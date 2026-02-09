package op

import (
	"reflect"

	"github.com/kaptinlin/jsonpatch/internal"
)

// FlipOperation represents a flip operation that inverts boolean values
// or converts other types to boolean and then inverts them.
type FlipOperation struct {
	BaseOp
}

// NewFlip creates a new flip operation.
func NewFlip(path []string) *FlipOperation {
	return &FlipOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (f *FlipOperation) Op() internal.OpType {
	return internal.OpFlipType
}

// Code returns the operation code.
func (f *FlipOperation) Code() int {
	return internal.OpFlipCode
}

// Apply applies the flip operation to the document.
func (f *FlipOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(f.Path()) == 0 {
		flipped := flipValue(doc)
		return internal.OpResult[any]{Doc: flipped, Old: doc}, nil
	}

	value, err := getValue(doc, f.Path())
	var oldValue any
	if err != nil {
		// Path doesn't exist: treat as undefined (false), flip to true
		value = nil
	}
	oldValue = value

	flipped := flipValue(value)

	err = setValueAtPath(doc, f.Path(), flipped)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// flipValue inverts a value's boolean interpretation.
func flipValue(value any) any {
	if v, ok := value.(bool); ok {
		return !v
	}
	return !toBool(value)
}

// toBool converts a value to its boolean interpretation following JavaScript semantics.
func toBool(value any) bool {
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
		return toBoolReflect(value)
	}
}

// toBoolReflect handles boolean conversion for complex types via reflection.
func toBoolReflect(value any) bool {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.Struct:
		return true
	case reflect.Chan:
		return val.Len() > 0
	case reflect.Pointer, reflect.Interface, reflect.Func, reflect.UnsafePointer:
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
	}
	return true
}

// ToJSON serializes the operation to JSON format.
func (f *FlipOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpFlipType),
		Path: formatPath(f.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (f *FlipOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpFlipCode, f.Path()}, nil
}

// Validate validates the flip operation.
func (f *FlipOperation) Validate() error {
	return nil
}
