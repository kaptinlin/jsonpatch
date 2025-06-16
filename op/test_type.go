package op

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpTestTypeOperation represents a test operation that checks if a value is of a specific type.
type OpTestTypeOperation struct {
	BaseOp
	Types []string `json:"type"` // Expected type names
}

// NewOpTestTypeOperation creates a new OpTestTypeOperation operation.
func NewOpTestTypeOperation(path []string, expectedType string) *OpTestTypeOperation {
	return &OpTestTypeOperation{
		BaseOp: NewBaseOp(path),
		Types:  []string{expectedType},
	}
}

// NewOpTestTypeOperationMultiple creates a new OpTestTypeOperation operation with multiple internal.
func NewOpTestTypeOperationMultiple(path []string, expectedTypes []string) *OpTestTypeOperation {
	return &OpTestTypeOperation{
		BaseOp: NewBaseOp(path),
		Types:  expectedTypes,
	}
}

// Op returns the operation type.
func (op *OpTestTypeOperation) Op() internal.OpType {
	return internal.OpTestTypeType
}

// Code returns the operation code.
func (op *OpTestTypeOperation) Code() int {
	return internal.OpTestTypeCode
}

// Path returns the operation path.
func (op *OpTestTypeOperation) Path() []string {
	return op.path
}

// getValueAndCheckType retrieves the value and checks if it matches any expected type
func (op *OpTestTypeOperation) getValueAndCheckType(doc any) (interface{}, string, bool, error) {
	// Get target value
	val, err := getValue(doc, op.Path())
	if err != nil {
		// Path access error means the path doesn't exist
		return nil, "", false, err
	}

	// Get the actual type of the value
	actualType := getTypeNameWithIntegerSupport(val)

	// Check if the type matches any of the expected types
	typeMatches := op.checkTypeMatch(actualType)

	return val, actualType, typeMatches, nil
}

// checkTypeMatch checks if actualType matches any expected type
func (op *OpTestTypeOperation) checkTypeMatch(actualType string) bool {
	for _, expectedType := range op.Types {
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
func (op *OpTestTypeOperation) Test(doc any) (bool, error) {
	_, _, typeMatches, err := op.getValueAndCheckType(doc)
	if err != nil {
		// Path access error means the path doesn't exist
		// For JSON Patch test operations, path not found means test fails (returns false)
		// This is correct JSON Patch semantics - returning nil error with false result
		//nolint:nilerr // This is intentional behavior for test operations
		return false, nil
	}
	return typeMatches, nil
}

// Not returns false (test_type operation doesn't support not modifier).
func (op *OpTestTypeOperation) Not() bool {
	return false
}

// Apply applies the test type operation to the document.
func (op *OpTestTypeOperation) Apply(doc any) (internal.OpResult[any], error) {
	_, actualType, typeMatches, err := op.getValueAndCheckType(doc)
	if err != nil {
		return internal.OpResult[any]{}, ErrPathNotFound
	}

	if !typeMatches {
		expectedTypesStr := strings.Join(op.Types, ", ")
		return internal.OpResult[any]{}, fmt.Errorf("%w: expected type %s, got %s", ErrTypeMismatch, expectedTypesStr, actualType)
	}

	return internal.OpResult[any]{Doc: doc}, nil
}

// getTypeName returns the JSON type name of a value.
func getTypeName(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case bool:
		return "boolean"
	case float64, float32:
		return "number"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "number" // For backward compatibility, all integers are "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		// For other types, use reflection
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "number"
		case reflect.Float32, reflect.Float64:
			return "number"
		case reflect.String:
			return "string"
		case reflect.Bool:
			return "boolean"
		case reflect.Slice, reflect.Array:
			return "array"
		case reflect.Map:
			return "object"
		case reflect.Struct:
			return "object"
		case reflect.Ptr, reflect.Interface:
			if rt.Elem() != nil {
				return getTypeName(reflect.ValueOf(v).Elem().Interface())
			}
			return "object"
		case reflect.Invalid:
			return "null"
		case reflect.Uintptr:
			return "number"
		case reflect.Complex64, reflect.Complex128:
			return "number"
		case reflect.Chan, reflect.Func, reflect.UnsafePointer:
			return "object"
		default:
			return "object"
		}
	}
}

// getTypeNameWithIntegerSupport returns the JSON type name of a value, distinguishing integers from floats.
func getTypeNameWithIntegerSupport(value interface{}) string {
	if value == nil {
		return "null"
	}

	switch v := value.(type) {
	case bool:
		return "boolean"
	case float64:
		// Check if it's an integer (whole number)
		if v == float64(int64(v)) {
			return "integer"
		}
		return "number"
	case string:
		return "string"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		// For other types, use reflection
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return "integer"
		case reflect.Float32, reflect.Float64:
			return "number"
		case reflect.String:
			return "string"
		case reflect.Bool:
			return "boolean"
		case reflect.Slice, reflect.Array:
			return "array"
		case reflect.Map:
			return "object"
		case reflect.Struct:
			return "object"
		case reflect.Ptr, reflect.Interface:
			if rt.Elem() != nil {
				return getTypeNameWithIntegerSupport(reflect.ValueOf(v).Elem().Interface())
			}
			return "object"
		case reflect.Invalid:
			return "null"
		case reflect.Uintptr:
			return "integer"
		case reflect.Complex64, reflect.Complex128:
			return "number"
		case reflect.Chan, reflect.Func, reflect.UnsafePointer:
			return "object"
		default:
			return "object"
		}
	}
}

// ToJSON serializes the operation to JSON format.
func (op *OpTestTypeOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpTestTypeType),
		"path": formatPath(op.Path()),
		"type": op.Types,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpTestTypeOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpTestTypeCode, op.Path(), op.Types}, nil
}

// Validate validates the test type operation.
func (op *OpTestTypeOperation) Validate() error {
	if len(op.Path()) == 0 {
		return ErrPathEmpty
	}
	if len(op.Types) == 0 {
		return ErrEmptyTypeList
	}
	// Validate that all types are known valid types
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
		"object":  true,
		"array":   true,
		"null":    true,
		"integer": true, // Special type that's also valid
	}
	for _, t := range op.Types {
		if !validTypes[t] {
			return ErrInvalidType
		}
	}
	return nil
}
