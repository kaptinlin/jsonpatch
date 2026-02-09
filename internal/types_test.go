package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidJSONPatchType(t *testing.T) {
	valid := []string{"string", "number", "boolean", "object", "integer", "array", "null"}
	for _, v := range valid {
		assert.True(t, IsValidJSONPatchType(v), "expected %q to be valid", v)
	}

	invalid := []string{"", "unknown", "float", "map", "slice", "int"}
	for _, v := range invalid {
		assert.False(t, IsValidJSONPatchType(v), "expected %q to be invalid", v)
	}
}

func TestGetJSONPatchType(t *testing.T) {
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
		// Arrays.
		{"[]any", []any{1, 2}, JSONPatchTypeArray},
		{"[]string", []string{"a"}, JSONPatchTypeArray},
		{"[]int", []int{1, 2}, JSONPatchTypeArray},
		{"[]float64", []float64{1.1}, JSONPatchTypeArray},
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetJSONPatchType(tt.val))
		})
	}
}

func TestIsJSONPatchOperation(t *testing.T) {
	valid := []string{"add", "remove", "replace", "move", "copy", "test"}
	for _, op := range valid {
		assert.True(t, IsJSONPatchOperation(op), "expected %q to be valid", op)
	}

	invalid := []string{"", "inc", "flip", "and", "defined", "unknown"}
	for _, op := range invalid {
		assert.False(t, IsJSONPatchOperation(op), "expected %q to be invalid", op)
	}
}

func TestIsFirstOrderPredicateOperation(t *testing.T) {
	valid := []string{
		"test", "defined", "undefined", "test_type",
		"test_string", "test_string_len", "contains",
		"ends", "starts", "in", "less", "more", "matches",
	}
	for _, op := range valid {
		assert.True(t, IsFirstOrderPredicateOperation(op), "expected %q to be valid", op)
	}

	invalid := []string{"", "add", "and", "or", "not", "inc"}
	for _, op := range invalid {
		assert.False(t, IsFirstOrderPredicateOperation(op), "expected %q to be invalid", op)
	}
}

func TestIsSecondOrderPredicateOperation(t *testing.T) {
	valid := []string{"and", "or", "not"}
	for _, op := range valid {
		assert.True(t, IsSecondOrderPredicateOperation(op), "expected %q to be valid", op)
	}

	invalid := []string{"", "add", "test", "defined", "inc"}
	for _, op := range invalid {
		assert.False(t, IsSecondOrderPredicateOperation(op), "expected %q to be invalid", op)
	}
}

func TestIsPredicateOperation(t *testing.T) {
	firstOrder := []string{
		"test", "defined", "undefined", "test_type",
		"test_string", "test_string_len", "contains",
		"ends", "starts", "in", "less", "more", "matches",
	}
	for _, op := range firstOrder {
		assert.True(t, IsPredicateOperation(op), "expected first-order %q to be valid", op)
	}

	secondOrder := []string{"and", "or", "not"}
	for _, op := range secondOrder {
		assert.True(t, IsPredicateOperation(op), "expected second-order %q to be valid", op)
	}

	invalid := []string{"", "add", "remove", "inc", "flip", "unknown"}
	for _, op := range invalid {
		assert.False(t, IsPredicateOperation(op), "expected %q to be invalid", op)
	}
}

func TestIsJSONPatchExtendedOperation(t *testing.T) {
	valid := []string{"str_ins", "str_del", "flip", "inc", "split", "merge", "extend"}
	for _, op := range valid {
		assert.True(t, IsJSONPatchExtendedOperation(op), "expected %q to be valid", op)
	}

	invalid := []string{"", "add", "test", "and", "defined", "unknown"}
	for _, op := range invalid {
		assert.False(t, IsJSONPatchExtendedOperation(op), "expected %q to be invalid", op)
	}
}
