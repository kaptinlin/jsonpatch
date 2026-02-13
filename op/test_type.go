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

// Code returns the operation code.
func (tt *TestTypeOperation) Code() int {
	return internal.OpTestTypeCode
}

// getValueAndCheckType retrieves the value and checks if it matches any expected type.
func (tt *TestTypeOperation) getValueAndCheckType(doc any) (any, string, bool, error) {
	val, err := getValue(doc, tt.Path())
	if err != nil {
		return nil, "", false, err
	}

	actualType := getTypeNameWithIntegerSupport(val)
	typeMatches := tt.checkTypeMatch(actualType)

	return val, actualType, typeMatches, nil
}

// checkTypeMatch checks if actualType matches any expected type.
func (tt *TestTypeOperation) checkTypeMatch(actualType string) bool {
	return slices.ContainsFunc(tt.Types, func(expectedType string) bool {
		return actualType == expectedType ||
			// Special case: if expected type is "number" and actual is "integer", it should match
			(expectedType == "number" && actualType == "integer")
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
// All integer types are reported as "number" for backward compatibility.
func getTypeName(value any) string {
	return getTypeNameInternal(value, false)
}

// getTypeNameWithIntegerSupport returns the JSON type name of a value,
// distinguishing integers from floats (returns "integer" for whole numbers).
func getTypeNameWithIntegerSupport(value any) string {
	return getTypeNameInternal(value, true)
}

// getTypeNameInternal is the shared implementation for type name detection.
// When distinguishInteger is true, integer types return "integer" instead of "number".
func getTypeNameInternal(value any, distinguishInteger bool) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case bool:
		return "boolean"
	case float64:
		if distinguishInteger && v == float64(int64(v)) {
			return "integer"
		}
		return "number"
	case float32:
		return "number"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if distinguishInteger {
			return "integer"
		}
		return "number"
	case string:
		return "string"
	case []any:
		return "array"
	case map[string]any:
		return "object"
	default:
		return getTypeNameViaReflection(v, distinguishInteger)
	}
}

// getTypeNameViaReflection handles type detection for non-standard types using reflection.
func getTypeNameViaReflection(v any, distinguishInteger bool) string {
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		if distinguishInteger {
			return "integer"
		}
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
			return getTypeNameInternal(reflect.ValueOf(v).Elem().Interface(), distinguishInteger)
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

// ToJSON serializes the operation to JSON format.
func (tt *TestTypeOperation) ToJSON() (internal.Operation, error) {
	if len(tt.Types) == 1 {
		return internal.Operation{
			Op:   string(internal.OpTestTypeType),
			Path: formatPath(tt.Path()),
			Type: tt.Types[0],
		}, nil
	}
	return internal.Operation{
		Op:    string(internal.OpTestTypeType),
		Path:  formatPath(tt.Path()),
		Value: tt.Types,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (tt *TestTypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestTypeCode, tt.Path(), tt.Types}, nil
}

// Validate validates the test type operation.
func (tt *TestTypeOperation) Validate() error {
	if len(tt.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(tt.Types) == 0 {
		return ErrEmptyTypeList
	}
	for _, t := range tt.Types {
		if !internal.IsValidJSONPatchType(t) {
			return ErrInvalidType
		}
	}
	return nil
}
