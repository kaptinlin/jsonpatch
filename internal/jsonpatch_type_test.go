package internal

import (
	"math"
	"testing"
)

func TestIsValidJSONPatchType(t *testing.T) {
	t.Parallel()
	valid := []string{
		"string", "number", "boolean",
		"object", "integer", "array", "null",
	}
	for _, v := range valid {
		if !IsValidJSONPatchType(v) {
			t.Errorf("IsValidJSONPatchType(%q) = false, want true", v)
		}
	}

	invalid := []string{"", "unknown", "float", "map", "slice", "int"}
	for _, v := range invalid {
		if IsValidJSONPatchType(v) {
			t.Errorf("IsValidJSONPatchType(%q) = true, want false", v)
		}
	}
}

func TestGetJSONPatchType(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		val  any
		want JSONPatchType
	}{
		// Null.
		{"nil", nil, JSONPatchTypeNull},
		{"unknown type", struct{}{}, JSONPatchTypeNull},
		// String.
		{"string", "hello", JSONPatchTypeString},
		{"empty string", "", JSONPatchTypeString},
		// Boolean.
		{"true", true, JSONPatchTypeBoolean},
		{"false", false, JSONPatchTypeBoolean},
		// Arrays (common types).
		{"[]any", []any{1, 2}, JSONPatchTypeArray},
		{"[]string", []string{"a"}, JSONPatchTypeArray},
		{"[]int", []int{1, 2}, JSONPatchTypeArray},
		{"[]float64", []float64{1.1}, JSONPatchTypeArray},
		// Arrays (reflect-based detection).
		{"[]bool", []bool{true, false}, JSONPatchTypeArray},
		{"[]uint", []uint{1, 2}, JSONPatchTypeArray},
		{"[]int64", []int64{1, 2}, JSONPatchTypeArray},
		{"[2]int", [2]int{1, 2}, JSONPatchTypeArray},
		// Object.
		{"map[string]any", map[string]any{"k": "v"}, JSONPatchTypeObject},
		// Integer types.
		{"int", 42, JSONPatchTypeInteger},
		{"int8", int8(1), JSONPatchTypeInteger},
		{"int16", int16(1), JSONPatchTypeInteger},
		{"int32", int32(1), JSONPatchTypeInteger},
		{"int64", int64(1), JSONPatchTypeInteger},
		{"uint", uint(1), JSONPatchTypeInteger},
		{"uint8", uint8(1), JSONPatchTypeInteger},
		{"uint16", uint16(1), JSONPatchTypeInteger},
		{"uint32", uint32(1), JSONPatchTypeInteger},
		{"uint64", uint64(1), JSONPatchTypeInteger},
		// Float as integer (whole numbers).
		{"float64 whole", float64(5), JSONPatchTypeInteger},
		{"float32 whole", float32(5), JSONPatchTypeInteger},
		// Float as number (fractional).
		{"float64 frac", 3.14, JSONPatchTypeNumber},
		{"float32 frac", float32(3.14), JSONPatchTypeNumber},
		// Float edge cases (NaN, Inf).
		{"float64 NaN", math.NaN(), JSONPatchTypeNumber},
		{"float64 +Inf", math.Inf(1), JSONPatchTypeNumber},
		{"float64 -Inf", math.Inf(-1), JSONPatchTypeNumber},
		{"float32 NaN", float32(math.NaN()), JSONPatchTypeNumber},
		{"float32 +Inf", float32(math.Inf(1)), JSONPatchTypeNumber},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := GetJSONPatchType(tt.val); got != tt.want {
				t.Errorf("GetJSONPatchType(%v) = %v, want %v", tt.val, got, tt.want)
			}
		})
	}
}
