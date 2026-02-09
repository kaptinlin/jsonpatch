package op

import (
	"fmt"
	"reflect"
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
func (o *TestTypeOperation) Op() internal.OpType {
	return internal.OpTestTypeType
}

// Code returns the operation code.
func (o *TestTypeOperation) Code() int {
	return internal.OpTestTypeCode
}

// getValueAndCheckType retrieves the value and checks if it matches any expected type.
func (o *TestTypeOperation) getValueAndCheckType(doc any) (any, string, bool, error) {
	val, err := getValue(doc, o.Path())
	if err != nil {
		return nil, "", false, err
	}

	actualType := getTypeNameWithIntegerSupport(val)
	typeMatches := o.checkTypeMatch(actualType)

	return val, actualType, typeMatches, nil
}

// checkTypeMatch checks if actualType matches any expected type.
func (o *TestTypeOperation) checkTypeMatch(actualType string) bool {
	for _, expectedType := range o.Types {
		if actualType == expectedType {
			return true
		}
		// Special case: if expected type is "number" and actual is "integer", it should match
		if expectedType == "number" && actualType == "integer" {
			return true
		}
	}
	return false
}

// Test evaluates the test type predicate condition.
func (o *TestTypeOperation) Test(doc any) (bool, error) {
	_, _, typeMatches, err := o.getValueAndCheckType(doc)
	if err != nil {
		// Path access error means the path doesn't exist
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return typeMatches, nil
}

// Apply applies the test type operation to the document.
func (o *TestTypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	_, actualType, typeMatches, err := o.getValueAndCheckType(doc)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	if !typeMatches {
		expectedTypesStr := strings.Join(o.Types, ", ")
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
func (o *TestTypeOperation) ToJSON() (internal.Operation, error) {
	if len(o.Types) == 1 {
		return internal.Operation{
			Op:   string(internal.OpTestTypeType),
			Path: formatPath(o.Path()),
			Type: o.Types[0],
		}, nil
	}
	return internal.Operation{
		Op:    string(internal.OpTestTypeType),
		Path:  formatPath(o.Path()),
		Value: o.Types,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (o *TestTypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestTypeCode, o.Path(), o.Types}, nil
}

// Validate validates the test type operation.
func (o *TestTypeOperation) Validate() error {
	if len(o.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(o.Types) == 0 {
		return ErrEmptyTypeList
	}
	for _, t := range o.Types {
		if !internal.IsValidJSONPatchType(t) {
			return ErrInvalidType
		}
	}
	return nil
}
