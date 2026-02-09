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
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return JSONPatchTypeInteger
	case float32:
		f := float64(v)
		if !math.IsNaN(f) && !math.IsInf(f, 0) && math.Trunc(f) == f {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber
	case float64:
		if !math.IsNaN(v) && !math.IsInf(v, 0) && math.Trunc(v) == v {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber
	default:
		if isSliceKind(value) {
			return JSONPatchTypeArray
		}
		return JSONPatchTypeNull
	}
}

// isSliceKind reports whether value is a slice or array using reflection.
func isSliceKind(value any) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}
