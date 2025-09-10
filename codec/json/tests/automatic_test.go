// Package tests contains automated test cases for JSON codec functionality.
package tests

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/stretchr/testify/assert"
)

// TestAutomaticCodec tests automatic encoding/decoding of operations.
// This test matches the TypeScript automatic.spec.ts
func TestAutomaticCodec(t *testing.T) {
	// Configure options for testing
	options := json.PatchOptions{}

	// Use all sample operations (equivalent to TypeScript's operations)
	for name, operation := range SampleOperations {
		t.Run(name, func(t *testing.T) {
			// Decode operation
			ops, err := json.Decode([]map[string]interface{}{operation}, options)
			if err != nil {
				t.Logf("Failed to decode operation %s: %v", name, err)
				t.Logf("Operation: %+v", operation)
				assert.NoError(t, err)
				return
			}
			assert.Len(t, ops, 1)

			// Encode back to JSON
			encoded, err := json.Encode(ops)
			if err != nil {
				t.Logf("Failed to encode operation %s: %v", name, err)
				assert.NoError(t, err)
				return
			}
			assert.Len(t, encoded, 1)

			// Verify round-trip consistency
			assert.Equal(t, operation, encoded[0])
		})
	}
}

func TestCodecRoundTrip(t *testing.T) {
	// Configure options for testing
	options := json.PatchOptions{}

	// Test that encoding and decoding preserves operation structure
	originalOps := []map[string]interface{}{
		{
			"op":    "add",
			"path":  "/foo",
			"value": "bar",
		},
		{
			"op":   "remove",
			"path": "/baz",
		},
		{
			"op":    "test",
			"path":  "/qux",
			"value": 42,
		},
	}

	// Decode
	decoded, err := json.Decode(originalOps, options)
	assert.NoError(t, err)

	// Encode back
	encoded, err := json.Encode(decoded)
	assert.NoError(t, err)

	// Should be identical
	assert.Equal(t, originalOps, encoded)
}
