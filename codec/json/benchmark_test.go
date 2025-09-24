// Package json implements performance benchmarks for JSON codec optimizations.
package json

import (
	"testing"

	"github.com/go-json-experiment/json"
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
var testOperationMaps = []map[string]interface{}{
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

	for i := 0; i < b.N; i++ {
		_, err := DecodeOperations(testOperations, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	options := PatchOptions{}
	ops, err := DecodeOperations(testOperations, options)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
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

	for i := 0; i < b.N; i++ {
		_, err := DecodeJSON(data, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncodeJSON(b *testing.B) {
	options := PatchOptions{}
	ops, err := DecodeOperations(testOperations, options)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
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

	for i := 0; i < b.N; i++ {
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
	require.Len(t, ops, 4)

	// Test encode
	encoded, err := Encode(ops)
	require.NoError(t, err)
	require.Len(t, encoded, 4)

	// Verify roundtrip matches original structure
	require.Equal(t, "add", encoded[0].Op)
	require.Equal(t, "/name", encoded[0].Path)
	require.Equal(t, "John", encoded[0].Value)
}

func TestJSONDecodeEncode(t *testing.T) {
	options := PatchOptions{}

	// Test JSON decode
	ops, err := DecodeJSON([]byte(testPatch), options)
	require.NoError(t, err)
	require.Len(t, ops, 4)

	// Test JSON encode
	data, err := EncodeJSON(ops)
	require.NoError(t, err)

	// Verify it's valid JSON and roundtrip works
	var decoded []map[string]interface{}
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	require.Len(t, decoded, 4)
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
		{Op: "in", Path: "/value", Value: []interface{}{"a", "b", "c"}},
		{Op: "type", Path: "/type", Value: "string"},
		{Op: "test_type", Path: "/test_type", Type: "number"},
		{Op: "test_string", Path: "/test_string", Str: "hello"},
		{Op: "test_string_len", Path: "/test_string_len", Len: 5},
		{Op: "inc", Path: "/counter", Inc: 1.0},
		{Op: "flip", Path: "/boolean"},
		{Op: "str_ins", Path: "/string", Pos: 1, Str: "prefix"},
		{Op: "str_del", Path: "/string", Pos: 1, Len: 5},
		{Op: "split", Path: "/array", Pos: 1, Props: map[string]interface{}{"key": "value"}},
		{Op: "merge", Path: "/object", Props: map[string]interface{}{"key": "value"}},
		{Op: "extend", Path: "/object", Props: map[string]interface{}{"newKey": "newValue"}},
	}

	// Test all operations can be decoded
	ops, err := DecodeOperations(allOps, options)
	require.NoError(t, err)
	require.Len(t, ops, len(allOps))

	// Test all operations can be encoded back
	encoded, err := Encode(ops)
	require.NoError(t, err)
	require.Len(t, encoded, len(allOps))
}

// Benchmark comparing struct-based vs map-based operations
func BenchmarkDecodeStructVsMap(b *testing.B) {
	options := PatchOptions{}

	b.Run("Struct", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := DecodeOperations(testOperations, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Map", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
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

	for i := 0; i < b.N; i++ {
		// Decode struct-based operations
		ops, err := DecodeOperations(testOperations, options)
		if err != nil {
			b.Fatal(err)
		}

		// Encode back to struct format
		_, err = Encode(ops)
		if err != nil {
			b.Fatal(err)
		}
	}
}
