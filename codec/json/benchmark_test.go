// Package json implements performance benchmarks for JSON codec optimizations.
package json

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/stretchr/testify/require"
)

// Simple test data for benchmarking
var testOperations = []map[string]interface{}{
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
		_, err := Decode(testOperations, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	options := PatchOptions{}
	ops, err := Decode(testOperations, options)
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
	ops, err := Decode(testOperations, options)
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

	// Test decode
	ops, err := Decode(testOperations, options)
	require.NoError(t, err)
	require.Len(t, ops, 4)

	// Test encode
	encoded, err := Encode(ops)
	require.NoError(t, err)
	require.Len(t, encoded, 4)

	// Verify roundtrip matches original structure
	require.Equal(t, "add", encoded[0]["op"])
	require.Equal(t, "/name", encoded[0]["path"])
	require.Equal(t, "John", encoded[0]["value"])
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

	allOps := []map[string]interface{}{
		{"op": "add", "path": "/add", "value": "test"},
		{"op": "remove", "path": "/remove"},
		{"op": "replace", "path": "/replace", "value": "new"},
		{"op": "test", "path": "/test", "value": "expected"},
		{"op": "move", "path": "/move", "from": "/source"},
		{"op": "copy", "path": "/copy", "from": "/source"},
		{"op": "defined", "path": "/defined"},
		{"op": "undefined", "path": "/undefined"},
		{"op": "contains", "path": "/text", "value": "search"},
		{"op": "starts", "path": "/text", "value": "prefix"},
		{"op": "ends", "path": "/text", "value": "suffix"},
		{"op": "more", "path": "/number", "value": 10.0},
		{"op": "less", "path": "/number", "value": 100.0},
		{"op": "in", "path": "/value", "value": []interface{}{"a", "b", "c"}},
		{"op": "type", "path": "/type", "value": "string"},
		{"op": "test_type", "path": "/test_type", "type": "number"},
		{"op": "test_string", "path": "/test_string", "str": "hello"},
		{"op": "test_string_len", "path": "/test_string_len", "len": 5.0},
		{"op": "inc", "path": "/counter", "inc": 1.0},
		{"op": "flip", "path": "/boolean"},
		{"op": "str_ins", "path": "/string", "pos": 0.0, "str": "prefix"},
		{"op": "str_del", "path": "/string", "pos": 0.0, "len": 5.0},
		{"op": "split", "path": "/array", "pos": 1.0, "props": map[string]interface{}{"key": "value"}},
		{"op": "merge", "path": "/object", "props": map[string]interface{}{"key": "value"}},
		{"op": "extend", "path": "/object", "props": map[string]interface{}{"newKey": "newValue"}},
	}

	// Test all operations can be decoded
	ops, err := Decode(allOps, options)
	require.NoError(t, err)
	require.Len(t, ops, len(allOps))

	// Test all operations can be encoded back
	encoded, err := Encode(ops)
	require.NoError(t, err)
	require.Len(t, encoded, len(allOps))
}
