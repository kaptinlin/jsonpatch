package tests

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
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
			ops, err := json.Decode([]map[string]any{operation}, options)
			if err != nil {
				t.Logf("Failed to decode operation %s: %v", name, err)
				t.Logf("Operation: %+v", operation)
				if err != nil {
					t.Errorf("Decode() error = %v", err)
				}
				return
			}
			if len(ops) != 1 {
				t.Errorf("Decode() returned %d ops, want 1", len(ops))
			}

			// Encode back to JSON
			encoded, err := json.Encode(ops)
			if err != nil {
				t.Logf("Failed to encode operation %s: %v", name, err)
				if err != nil {
					t.Errorf("Encode() error = %v", err)
				}
				return
			}
			if len(encoded) != 1 {
				t.Errorf("Encode() returned %d ops, want 1", len(encoded))
			}

			// Convert expected operation map to Operation struct for comparison
			expectedOp := mapToOperation(operation)

			// Verify round-trip consistency by comparing the struct fields
			if got, want := encoded[0].Op, expectedOp.Op; got != want {
				t.Errorf("Op = %v, want %v", got, want)
			}
			if got, want := encoded[0].Path, expectedOp.Path; got != want {
				t.Errorf("Path = %v, want %v", got, want)
			}
			if diff := cmp.Diff(expectedOp.Value, encoded[0].Value); diff != "" {
				t.Errorf("Value mismatch (-want +got):\n%s", diff)
			}
			if got, want := encoded[0].From, expectedOp.From; got != want {
				t.Errorf("From = %v, want %v", got, want)
			}
			if got, want := encoded[0].Inc, expectedOp.Inc; got != want {
				t.Errorf("Inc = %v, want %v", got, want)
			}
			if got, want := encoded[0].Pos, expectedOp.Pos; got != want {
				t.Errorf("Pos = %v, want %v", got, want)
			}
			if got, want := encoded[0].Str, expectedOp.Str; got != want {
				t.Errorf("Str = %v, want %v", got, want)
			}
			if got, want := encoded[0].Len, expectedOp.Len; got != want {
				t.Errorf("Len = %v, want %v", got, want)
			}
			if got, want := encoded[0].Not, expectedOp.Not; got != want {
				t.Errorf("Not = %v, want %v", got, want)
			}
			if diff := cmp.Diff(expectedOp.Type, encoded[0].Type); diff != "" {
				t.Errorf("Type mismatch (-want +got):\n%s", diff)
			}
			if got, want := encoded[0].IgnoreCase, expectedOp.IgnoreCase; got != want {
				t.Errorf("IgnoreCase = %v, want %v", got, want)
			}
			if diff := cmp.Diff(expectedOp.Props, encoded[0].Props); diff != "" {
				t.Errorf("Props mismatch (-want +got):\n%s", diff)
			}
			if got, want := encoded[0].DeleteNull, expectedOp.DeleteNull; got != want {
				t.Errorf("DeleteNull = %v, want %v", got, want)
			}
			if diff := cmp.Diff(expectedOp.OldValue, encoded[0].OldValue); diff != "" {
				t.Errorf("OldValue mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCodecRoundTrip(t *testing.T) {
	// Configure options for testing
	options := json.PatchOptions{}

	// Test that encoding and decoding preserves operation structure
	originalOps := []map[string]any{
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
	if err != nil {
		t.Errorf("Decode() error = %v", err)
	}

	// Encode back
	encoded, err := json.Encode(decoded)
	if err != nil {
		t.Errorf("Encode() error = %v", err)
	}

	// Convert expected operations to Operation structs for comparison
	expectedOps := make([]internal.Operation, len(originalOps))
	for i, opMap := range originalOps {
		expectedOps[i] = mapToOperation(opMap)
	}

	// Should be structurally identical
	if diff := cmp.Diff(expectedOps, encoded); diff != "" {
		t.Errorf("round-trip mismatch (-want +got):\n%s", diff)
	}
}

// mapToOperation converts a map[string]any to internal.Operation struct
func mapToOperation(opMap map[string]any) internal.Operation {
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
	if val, ok := opMap["props"].(map[string]any); ok {
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
			if typeSlice, ok := typeField.([]any); ok {
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
	if val, ok := opMap["apply"].([]any); ok {
		applyOps := make([]internal.Operation, 0, len(val))
		for _, subOp := range val {
			if subOpMap, ok := subOp.(map[string]any); ok {
				applyOps = append(applyOps, mapToOperation(subOpMap))
			}
		}
		op.Apply = applyOps
	}

	return op
}
