package json

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

// Simple test data for benchmarking - using struct-based Operations
var testOperations = []internal.Operation{
	{Op: "add", Path: "/name", Value: "John"},
	{Op: "replace", Path: "/age", Value: 30},
	{Op: "remove", Path: "/temp"},
	{Op: "test", Path: "/active", Value: true},
}

// Legacy test data as maps for compatibility testing
var testOperationMaps = []map[string]any{
	{"op": "add", "path": "/name", "value": "John"},
	{"op": "replace", "path": "/age", "value": 30},
	{"op": "remove", "path": "/temp"},
	{"op": "test", "path": "/active", "value": true},
}

var testPatch = `[
	{"op": "add", "path": "/name", "value": "John"},
	{"op": "replace", "path": "/age", "value": 30},
	{"op": "remove", "path": "/temp"},
	{"op": "test", "path": "/active", "value": true}
]`

func BenchmarkDecode(b *testing.B) {
	options := PatchOptions{}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, err := DecodeOperations(testOperations, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	options := PatchOptions{}
	ops, err := DecodeOperations(testOperations, options)
	if err != nil {
		b.Fatalf("setup DecodeOperations: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, err := Encode(ops)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDecodeJSON(b *testing.B) {
	data := []byte(testPatch)
	options := PatchOptions{}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, err := DecodeJSON(data, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeJSON(b *testing.B) {
	options := PatchOptions{}
	ops, err := DecodeOperations(testOperations, options)
	if err != nil {
		b.Fatalf("setup DecodeOperations: %v", err)
	}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		_, err := EncodeJSON(ops)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	options := PatchOptions{}
	data := []byte(testPatch)
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ops, err := DecodeJSON(data, options)
		if err != nil {
			b.Fatal(err)
		}
		_, err = EncodeJSON(ops)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test basic functionality
func TestBasicDecodeEncode(t *testing.T) {
	options := PatchOptions{}

	// Test decode using struct-based operations
	ops, err := DecodeOperations(testOperations, options)
	require.NoError(t, err)
	if len(ops) != 4 {
		t.Fatalf("DecodeOperations returned %d ops, want 4", len(ops))
	}

	// Test encode
	encoded, err := Encode(ops)
	require.NoError(t, err)
	if len(encoded) != 4 {
		t.Fatalf("Encode returned %d ops, want 4", len(encoded))
	}

	// Verify roundtrip matches original structure
	if got, want := encoded[0].Op, "add"; got != want {
		assert.Equal(t, want, got, "encoded[0].Op")
	}
	if got, want := encoded[0].Path, "/name"; got != want {
		assert.Equal(t, want, got, "encoded[0].Path")
	}
	assert.Equal(t, any("John"), encoded[0].Value)
}

func TestJSONDecodeEncode(t *testing.T) {
	options := PatchOptions{}

	// Test JSON decode
	ops, err := DecodeJSON([]byte(testPatch), options)
	require.NoError(t, err)
	if len(ops) != 4 {
		t.Fatalf("DecodeJSON returned %d ops, want 4", len(ops))
	}

	// Test JSON encode
	data, err := EncodeJSON(ops)
	require.NoError(t, err)

	// Verify it's valid JSON and roundtrip works
	var decoded []map[string]any
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("json.Unmarshal: %v", err)
	}
	if len(decoded) != 4 {
		t.Fatalf("json.Unmarshal returned %d ops, want 4", len(decoded))
	}
}

// Test all operation types work with API
func TestAllOperationTypes(t *testing.T) {
	options := PatchOptions{}

	allOps := []internal.Operation{
		{Op: "add", Path: "/add", Value: "test"},
		{Op: "remove", Path: "/remove"},
		{Op: "replace", Path: "/replace", Value: "new"},
		{Op: "test", Path: "/test", Value: "expected"},
		{Op: "move", Path: "/move", From: "/source"},
		{Op: "copy", Path: "/copy", From: "/source"},
		{Op: "defined", Path: "/defined"},
		{Op: "undefined", Path: "/undefined"},
		{Op: "contains", Path: "/text", Value: "search"},
		{Op: "starts", Path: "/text", Value: "prefix"},
		{Op: "ends", Path: "/text", Value: "suffix"},
		{Op: "more", Path: "/number", Value: 10.0},
		{Op: "less", Path: "/number", Value: 100.0},
		{Op: "in", Path: "/value", Value: []any{"a", "b", "c"}},
		{Op: "type", Path: "/type", Value: "string"},
		{Op: "test_type", Path: "/test_type", Type: "number"},
		{Op: "test_string", Path: "/test_string", Str: "hello"},
		{Op: "test_string_len", Path: "/test_string_len", Len: 5},
		{Op: "inc", Path: "/counter", Inc: 1.0},
		{Op: "flip", Path: "/boolean"},
		{Op: "str_ins", Path: "/string", Pos: 1, Str: "prefix"},
		{Op: "str_del", Path: "/string", Pos: 1, Len: 5},
		{Op: "split", Path: "/array", Pos: 1, Props: map[string]any{"key": "value"}},
		{Op: "merge", Path: "/object", Props: map[string]any{"key": "value"}},
		{Op: "extend", Path: "/object", Props: map[string]any{"newKey": "newValue"}},
	}

	// Test all operations can be decoded
	ops, err := DecodeOperations(allOps, options)
	require.NoError(t, err)
	if len(ops) != len(allOps) {
		t.Fatalf("DecodeOperations returned %d ops, want %d", len(ops), len(allOps))
	}

	// Test all operations can be encoded back
	encoded, err := Encode(ops)
	require.NoError(t, err)
	if len(encoded) != len(allOps) {
		t.Fatalf("Encode returned %d ops, want %d", len(encoded), len(allOps))
	}
}

// TestOperationToMapNestedType verifies that the type field is preserved
// when converting nested operations (e.g., test_type inside and/or/not).
func TestOperationToMapNestedType(t *testing.T) {
	options := PatchOptions{}

	// test_type with string type inside "and" composite operation
	andOp := []internal.Operation{
		{
			Op:   "and",
			Path: "/data",
			Apply: []internal.Operation{
				{Op: "test_type", Path: "/value", Type: "number"},
				{Op: "defined", Path: "/value"},
			},
		},
	}

	ops, err := DecodeOperations(andOp, options)
	require.NoError(t, err)

	encoded, err := Encode(ops)
	require.NoError(t, err)
	if got, want := encoded[0].Op, "and"; got != want {
		assert.Equal(t, want, got, "Op")
	}
}

// TestOperationToMapNestedTypeList verifies that multi-type test_type
// operations are preserved through nested operationToMap conversion.
func TestOperationToMapNestedTypeList(t *testing.T) {
	options := PatchOptions{}

	andOp := []internal.Operation{
		{
			Op:   "and",
			Path: "/data",
			Apply: []internal.Operation{
				{
					Op:   "test_type",
					Path: "/value",
					Type: []any{"number", "string"},
				},
			},
		},
	}

	ops, err := DecodeOperations(andOp, options)
	require.NoError(t, err)

	encoded, err := Encode(ops)
	require.NoError(t, err)
	if got, want := encoded[0].Op, "and"; got != want {
		assert.Equal(t, want, got, "Op")
	}
}

// Benchmark comparing struct-based vs map-based operations
func BenchmarkDecodeStructVsMap(b *testing.B) {
	options := PatchOptions{}

	b.Run("Struct", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, err := DecodeOperations(testOperations, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for b.Loop() {
			_, err := Decode(testOperationMaps, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark the complete roundtrip with struct-based operations
func BenchmarkStructRoundTrip(b *testing.B) {
	options := PatchOptions{}
	b.ResetTimer()
	b.ReportAllocs()
	for b.Loop() {
		ops, err := DecodeOperations(testOperations, options)
		if err != nil {
			b.Fatal(err)
		}
		_, err = Encode(ops)
		if err != nil {
			b.Fatal(err)
		}
	}
}
