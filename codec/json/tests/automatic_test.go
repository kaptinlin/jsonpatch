// Package tests contains automated test cases for JSON codec functionality.
package tests

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
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

			// Convert expected operation map to Operation struct for comparison
			expectedOp := mapToOperation(operation)

			// Verify round-trip consistency by comparing the struct fields
			assert.Equal(t, expectedOp.Op, encoded[0].Op)
			assert.Equal(t, expectedOp.Path, encoded[0].Path)
			assert.Equal(t, expectedOp.Value, encoded[0].Value)
			assert.Equal(t, expectedOp.From, encoded[0].From)
			assert.Equal(t, expectedOp.Inc, encoded[0].Inc)
			assert.Equal(t, expectedOp.Pos, encoded[0].Pos)
			assert.Equal(t, expectedOp.Str, encoded[0].Str)
			assert.Equal(t, expectedOp.Len, encoded[0].Len)
			assert.Equal(t, expectedOp.Not, encoded[0].Not)
			assert.Equal(t, expectedOp.Type, encoded[0].Type)
			assert.Equal(t, expectedOp.IgnoreCase, encoded[0].IgnoreCase)
			assert.Equal(t, expectedOp.Props, encoded[0].Props)
			assert.Equal(t, expectedOp.DeleteNull, encoded[0].DeleteNull)
			assert.Equal(t, expectedOp.OldValue, encoded[0].OldValue)
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

	// Convert expected operations to Operation structs for comparison
	expectedOps := make([]internal.Operation, len(originalOps))
	for i, opMap := range originalOps {
		expectedOps[i] = mapToOperation(opMap)
	}

	// Should be structurally identical
	assert.Equal(t, expectedOps, encoded)
}

// mapToOperation converts a map[string]interface{} to internal.Operation struct
func mapToOperation(opMap map[string]interface{}) internal.Operation {
	op := internal.Operation{}

	// Set basic fields first
	if val, ok := opMap["op"].(string); ok {
		op.Op = val
	}
	if val, ok := opMap["path"].(string); ok {
		op.Path = val
	}
	if val, ok := opMap["from"].(string); ok {
		op.From = val
	}
	if val, ok := opMap["inc"].(float64); ok {
		op.Inc = val
	}
	if val, ok := opMap["pos"].(float64); ok {
		op.Pos = int(val)
	}
	if val, ok := opMap["str"].(string); ok {
		op.Str = val
	}
	if val, ok := opMap["len"].(float64); ok {
		op.Len = int(val)
	}
	if val, ok := opMap["not"].(bool); ok {
		op.Not = val
	}
	if val, ok := opMap["ignore_case"].(bool); ok {
		op.IgnoreCase = val
	}
	if val, ok := opMap["props"].(map[string]interface{}); ok {
		op.Props = val
	}
	if val, ok := opMap["deleteNull"].(bool); ok {
		op.DeleteNull = val
	}
	if val, exists := opMap["oldValue"]; exists {
		op.OldValue = val
	}

	// Handle value and type fields with special logic for test_type operations
	if op.Op == "test_type" {
		// For test_type operations, type field logic takes precedence
		if typeField, exists := opMap["type"]; exists {
			if typeSlice, ok := typeField.([]interface{}); ok {
				typeStrings := make([]string, len(typeSlice))
				for i, t := range typeSlice {
					if typeStr, ok := t.(string); ok {
						typeStrings[i] = typeStr
					}
				}
				// Single type goes into Type field, multiple types go into Value field
				if len(typeStrings) == 1 {
					op.Type = typeStrings[0]
				} else {
					op.Value = typeStrings
				}
			} else if typeStringSlice, ok := typeField.([]string); ok {
				// Single type goes into Type field, multiple types go into Value field
				if len(typeStringSlice) == 1 {
					op.Type = typeStringSlice[0]
				} else {
					op.Value = typeStringSlice
				}
			} else if typeStr, ok := typeField.(string); ok {
				// Single type goes into the Type field
				op.Type = typeStr
			}
		}
	} else {
		// For non-test_type operations, handle value and type normally
		if val, exists := opMap["value"]; exists {
			op.Value = val
		}
		if typeStr, ok := opMap["type"].(string); ok {
			op.Type = typeStr
		}
	}

	// Handle apply field for compound operations
	if val, ok := opMap["apply"].([]interface{}); ok {
		applyOps := make([]internal.Operation, 0, len(val))
		for _, subOp := range val {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				applyOps = append(applyOps, mapToOperation(subOpMap))
			}
		}
		op.Apply = applyOps
	}

	return op
}
