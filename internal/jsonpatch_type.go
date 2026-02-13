package internal

import (
	"math"
	"reflect"
)

// JSONPatchType represents valid JSON types for type-checking operations.
type JSONPatchType string

// Valid JSON Patch types.
const (
	JSONPatchTypeString  JSONPatchType = "string"
	JSONPatchTypeNumber  JSONPatchType = "number"
	JSONPatchTypeBoolean JSONPatchType = "boolean"
	JSONPatchTypeObject  JSONPatchType = "object"
	JSONPatchTypeInteger JSONPatchType = "integer"
	JSONPatchTypeArray   JSONPatchType = "array"
	JSONPatchTypeNull    JSONPatchType = "null"
)

// IsValidJSONPatchType reports whether typeStr is a valid JSON Patch
// type name.
func IsValidJSONPatchType(typeStr string) bool {
	switch JSONPatchType(typeStr) {
	case JSONPatchTypeString, JSONPatchTypeNumber,
		JSONPatchTypeBoolean, JSONPatchTypeObject,
		JSONPatchTypeInteger, JSONPatchTypeArray,
		JSONPatchTypeNull:
		return true
	default:
		return false
	}
}

// GetJSONPatchType returns the JSON Patch type for a Go value.
// Whole-number floats return "integer"; unknown types return "null".
func GetJSONPatchType(value any) JSONPatchType {
	if value == nil {
		return JSONPatchTypeNull
	}

	switch v := value.(type) {
	case string:
		return JSONPatchTypeString
	case bool:
		return JSONPatchTypeBoolean
	case map[string]any:
		return JSONPatchTypeObject
	case []any:
		return JSONPatchTypeArray
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return JSONPatchTypeInteger
	case float32:
		return classifyFloat(float64(v))
	case float64:
		return classifyFloat(v)
	default:
		if isSliceKind(value) {
			return JSONPatchTypeArray
		}
		return JSONPatchTypeNull
	}
}

// classifyFloat returns "integer" for whole-number floats, "number" otherwise.
func classifyFloat(f float64) JSONPatchType {
	if !math.IsNaN(f) && !math.IsInf(f, 0) && math.Trunc(f) == f {
		return JSONPatchTypeInteger
	}
	return JSONPatchTypeNumber
}

// isSliceKind reports whether value is a slice or array using reflection.
func isSliceKind(value any) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}
