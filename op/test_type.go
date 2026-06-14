package op

import (
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// TestTypeOperation represents a test operation that checks if a value is of a specific type.
type TestTypeOperation struct {
	BaseOp
	Types []string `json:"type"` // Expected type names
}

// NewTestType creates a new test type operation.
func NewTestType(path []string, expectedType string) *TestTypeOperation {
	return &TestTypeOperation{
		BaseOp: NewBaseOp(path),
		Types:  []string{expectedType},
	}
}

// NewTestTypeMultiple creates a new test type operation with multiple types.
func NewTestTypeMultiple(path []string, expectedTypes []string) *TestTypeOperation {
	return &TestTypeOperation{
		BaseOp: NewBaseOp(path),
		Types:  expectedTypes,
	}
}

// Op returns the operation type.
func (tt *TestTypeOperation) Op() internal.OpType {
	return internal.OpTestTypeType
}

// getValueAndCheckType retrieves the value and checks if it matches any expected type.
func (tt *TestTypeOperation) getValueAndCheckType(doc any) (any, string, bool, error) {
	val, err := value(doc, tt.Path())
	if err != nil {
		return nil, "", false, err
	}

	actualType := getTypeName(val)
	typeMatches := tt.checkTypeMatch(actualType, val)

	return val, actualType, typeMatches, nil
}

func (tt *TestTypeOperation) checkTypeMatch(actualType string, val any) bool {
	return slices.ContainsFunc(tt.Types, func(expectedType string) bool {
		if actualType == expectedType {
			return true
		}
		if expectedType == "integer" && isWholeNumber(val) {
			return true
		}
		return false
	})
}

// Test evaluates the test type predicate condition.
func (tt *TestTypeOperation) Test(doc any) (bool, error) {
	_, _, typeMatches, err := tt.getValueAndCheckType(doc)
	if err != nil {
		//nolint:nilerr // intentional: path not found means test fails
		return false, nil
	}
	return typeMatches, nil
}

// Apply applies the test type operation to the document.
func (tt *TestTypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	_, actualType, typeMatches, err := tt.getValueAndCheckType(doc)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	if !typeMatches {
		expectedTypesStr := strings.Join(tt.Types, ", ")
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected type %s, got %s", ErrTypeMismatch, expectedTypesStr, actualType)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// getTypeName returns the JSON type name of a value.
// All numeric Go values map to the JSON "number" type.
func getTypeName(val any) string {
	if val == nil {
		return "null"
	}

	switch val.(type) {
	case bool:
		return "boolean"
	case float64, float32:
		return "number"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return getTypeNameViaReflection(val)
	}
}

// isWholeNumber reports whether val is an integer or whole-number float.
func isWholeNumber(val any) bool {
	switch v := val.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float64:
		return v == float64(int64(v))
	case float32:
		return v == float32(int32(v))
	default:
		rt := reflect.TypeOf(val)
		if rt == nil {
			return false
		}
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return true
		case reflect.Float32, reflect.Float64:
			f := reflect.ValueOf(val).Float()
			return f == float64(int64(f))
		default:
			return false
		}
	}
}

// getTypeNameViaReflection handles type detection for non-standard types using reflection.
func getTypeNameViaReflection(v any) string {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return "number"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map, reflect.Struct:
		return "object"
	case reflect.Pointer, reflect.Interface:
		if rt.Elem() != nil {
			return getTypeName(reflect.ValueOf(v).Elem().Interface())
		}
		return "object"
	case reflect.Invalid:
		return "null"
	case reflect.Complex64, reflect.Complex128:
		return "number"
	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		return "object"
	default:
		return "object"
	}
}

// Validate validates the test type operation.
func (tt *TestTypeOperation) Validate() error {
	if len(tt.Types) == 0 {
		return ErrEmptyTypeList
	}
	if slices.ContainsFunc(tt.Types, func(t string) bool {
		return !internal.IsValidJSONPatchType(t)
	}) {
		return ErrInvalidType
	}
	return nil
}
