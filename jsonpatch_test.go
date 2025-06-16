package jsonpatch_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// =============================================================================
// TEST TYPES AND HELPERS
// =============================================================================

// profile is a test struct for generic patch testing
type profile struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Tags  []string `json:"tags"`
}

// =============================================================================
// BASIC FUNCTIONALITY TESTS
// =============================================================================

// TestApplyPatchBasic tests basic ApplyPatch functionality
func TestApplyPatchBasic(t *testing.T) {
	tests := []struct {
		name     string
		doc      interface{}
		patch    []jsonpatch.Operation
		expected interface{}
		wantErr  bool
	}{
		{
			name:     "empty patch",
			doc:      map[string]interface{}{"a": 1},
			patch:    []jsonpatch.Operation{},
			expected: map[string]interface{}{"a": 1},
			wantErr:  false,
		},
		{
			name: "single operation",
			doc:  map[string]interface{}{"a": 1},
			patch: []jsonpatch.Operation{
				map[string]interface{}{"op": "add", "path": "/b", "value": 2},
			},
			expected: map[string]interface{}{"a": 1, "b": 2},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := jsonpatch.ApplyPatch(tt.doc, tt.patch, jsonpatch.WithMutate(false))

			if tt.wantErr {
				require.Error(t, err, "Expected an error but got none")
				return
			}

			require.NoError(t, err, "Unexpected error: %v", err)

			assert.Equal(t, tt.expected, result.Doc, "Result document should match expected")
			assert.NotNil(t, result.Res, "Result operations should not be nil")
			assert.Len(t, result.Res, len(tt.patch), "Number of operation results should match patch length")
		})
	}
}

// TestValidateOperation tests operation validation
func TestValidateOperation(t *testing.T) {
	tests := []struct {
		name      string
		operation jsonpatch.Operation
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid add operation",
			operation: map[string]interface{}{"op": "add", "path": "/a", "value": 1},
			wantErr:   false,
		},
		{
			name:      "missing op field",
			operation: map[string]interface{}{"path": "/a", "value": 1},
			wantErr:   true,
			errMsg:    "missing required field 'op'",
		},
		{
			name:      "missing path field",
			operation: map[string]interface{}{"op": "add", "value": 1},
			wantErr:   true,
			errMsg:    "missing required field 'path'",
		},
		{
			name:      "missing value field for add",
			operation: map[string]interface{}{"op": "add", "path": "/a"},
			wantErr:   true,
			errMsg:    "missing required field 'value'",
		},
		{
			name:      "invalid operation type",
			operation: map[string]interface{}{"op": "invalid", "path": "/a"},
			wantErr:   true,
			errMsg:    "unknown operation 'invalid'",
		},
		{
			name:      "nil operation",
			operation: nil,
			wantErr:   true,
			errMsg:    "invalid operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := jsonpatch.ValidateOperation(tt.operation, false)

			if tt.wantErr {
				require.Error(t, err, "Expected validation error")
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err, "Validation should pass")
			}
		})
	}
}

// =============================================================================
// GENERIC TYPE TESTS
// =============================================================================

// TestApplyPatch_Struct tests applying patches to struct types
func TestApplyPatch_Struct(t *testing.T) {
	// Test data
	before := profile{
		Name: "John",
		Tags: []string{"dev"},
	}

	patch := []jsonpatch.Operation{
		map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		map[string]interface{}{"op": "add", "path": "/tags/-", "value": "golang"},
		map[string]interface{}{"op": "add", "path": "/email", "value": "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "Jane", result.Doc.Name)
	assert.Equal(t, "jane@example.com", result.Doc.Email)
	assert.Equal(t, []string{"dev", "golang"}, result.Doc.Tags)
	assert.Len(t, result.Res, 3, "Should have 3 operation results")

	// Verify original is unchanged (immutable by default)
	assert.Equal(t, "John", before.Name)
	assert.Empty(t, before.Email)
	assert.Equal(t, []string{"dev"}, before.Tags)
}

// TestApplyPatch_Map tests applying patches to map types
func TestApplyPatch_Map(t *testing.T) {
	// Test data
	before := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsonpatch.Operation{
		map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		map[string]interface{}{"op": "add", "path": "/tags/-", "value": "golang"},
		map[string]interface{}{"op": "add", "path": "/email", "value": "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "Jane", result.Doc["name"])
	assert.Equal(t, "jane@example.com", result.Doc["email"])
	assert.Equal(t, []any{"dev", "golang"}, result.Doc["tags"])
	assert.Len(t, result.Res, 3, "Should have 3 operation results")

	// Verify original is unchanged (immutable by default)
	assert.Equal(t, "John", before["name"])
	_, hasEmail := before["email"]
	assert.False(t, hasEmail)
}

// TestApplyPatch_JSONBytes tests applying patches to []byte containing JSON
func TestApplyPatch_JSONBytes(t *testing.T) {
	// Test data
	before := []byte(`{"name":"John","tags":["dev"]}`)

	patch := []jsonpatch.Operation{
		map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		map[string]interface{}{"op": "add", "path": "/tags/-", "value": "golang"},
		map[string]interface{}{"op": "add", "path": "/email", "value": "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	require.NoError(t, err)

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal(result.Doc, &resultMap)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "Jane", resultMap["name"])
	assert.Equal(t, "jane@example.com", resultMap["email"])
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Len(t, result.Res, 3, "Should have 3 operation results")

	// Verify original is unchanged
	var original map[string]any
	err = json.Unmarshal(before, &original)
	require.NoError(t, err)
	assert.Equal(t, "John", original["name"])
	_, hasEmail := original["email"]
	assert.False(t, hasEmail)
}

// TestApplyPatch_JSONString tests applying patches to JSON strings
func TestApplyPatch_JSONString(t *testing.T) {
	// Test data
	before := `{"name":"John","tags":["dev"]}`

	patch := []jsonpatch.Operation{
		map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		map[string]interface{}{"op": "add", "path": "/tags/-", "value": "golang"},
		map[string]interface{}{"op": "add", "path": "/email", "value": "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	require.NoError(t, err)

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal([]byte(result.Doc), &resultMap)
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "Jane", resultMap["name"])
	assert.Equal(t, "jane@example.com", resultMap["email"])
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Len(t, result.Res, 3, "Should have 3 operation results")
}

// =============================================================================
// ARRAY OPERATIONS TESTS
// =============================================================================

// TestArrayOperations demonstrates array manipulation with JSON Patch
func TestArrayOperations(t *testing.T) {
	// Document with array
	doc := map[string]interface{}{
		"items": []interface{}{
			map[string]interface{}{"id": 1, "name": "Item 1"},
			map[string]interface{}{"id": 2, "name": "Item 2"},
			map[string]interface{}{"id": 3, "name": "Item 3"},
		},
	}

	// Array operations
	patch := []jsonpatch.Operation{
		// Insert at beginning
		map[string]interface{}{
			"op":    "add",
			"path":  "/items/0",
			"value": map[string]interface{}{"id": 0, "name": "Item 0"},
		},
		// Append at end
		map[string]interface{}{
			"op":    "add",
			"path":  "/items/-",
			"value": map[string]interface{}{"id": 4, "name": "Item 4"},
		},
		// Update middle item
		map[string]interface{}{
			"op":    "replace",
			"path":  "/items/2/name",
			"value": "Updated Item 1",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	require.NoError(t, err)

	resultJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	t.Logf("Array operations result:\n%s", string(resultJSON))

	// Verify the result
	items := result.Doc["items"].([]interface{})
	assert.Len(t, items, 5, "Expected 5 items after operations")
}

// TestMultipleOperations demonstrates applying multiple operations
func TestMultipleOperations(t *testing.T) {
	doc := map[string]interface{}{
		"counters": map[string]interface{}{
			"a": 0,
			"b": 0,
		},
	}

	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":    "replace",
			"path":  "/counters/a",
			"value": 1,
		},
		map[string]interface{}{
			"op":    "replace",
			"path":  "/counters/b",
			"value": 2,
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	require.NoError(t, err)

	resultJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	t.Logf("Multiple operations result:\n%s", string(resultJSON))

	// Verify the result
	counters := result.Doc["counters"].(map[string]interface{})
	assert.Equal(t, 1, counters["a"], "Counter a should be updated")
	assert.Equal(t, 2, counters["b"], "Counter b should be updated")
}

// =============================================================================
// OPTIONS TESTS
// =============================================================================

// TestApplyPatch_WithMutate tests the mutate option
func TestApplyPatch_WithMutate(t *testing.T) {
	// Test data - using map for easier mutation testing
	original := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsonpatch.Operation{
		map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
	}

	// Apply patch with mutate=true
	result, err := jsonpatch.ApplyPatch(original, patch, jsonpatch.WithMutate(true))
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "Jane", result.Doc["name"])
	assert.Len(t, result.Res, 1, "Should have 1 operation result")
}

// =============================================================================
// COMPLEX DOCUMENT TESTS
// =============================================================================

// TestComplexDocument demonstrates complex document operations
func TestComplexDocument(t *testing.T) {
	// Complex nested document
	doc := map[string]interface{}{
		"company": map[string]interface{}{
			"name": "Tech Corp",
			"departments": []interface{}{
				map[string]interface{}{
					"name": "Engineering",
					"employees": []interface{}{
						map[string]interface{}{"id": 1, "name": "Alice", "role": "Developer"},
						map[string]interface{}{"id": 2, "name": "Bob", "role": "Manager"},
					},
				},
			},
		},
		"metadata": map[string]interface{}{
			"lastUpdated": "2023-01-01",
		},
	}

	// Complex operations
	patch := []jsonpatch.Operation{
		// Add new employee to Engineering
		map[string]interface{}{
			"op":   "add",
			"path": "/company/departments/0/employees/-",
			"value": map[string]interface{}{
				"id":   3,
				"name": "Charlie",
				"role": "Senior Developer",
			},
		},
		// Promote Bob to Senior Manager
		map[string]interface{}{
			"op":    "replace",
			"path":  "/company/departments/0/employees/1/role",
			"value": "Senior Manager",
		},
		// Update metadata
		map[string]interface{}{
			"op":    "replace",
			"path":  "/metadata/lastUpdated",
			"value": "2023-12-01",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	require.NoError(t, err)

	resultJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	t.Logf("Complex document result:\n%s", string(resultJSON))

	// Verify the changes
	company := result.Doc["company"].(map[string]interface{})
	departments := company["departments"].([]interface{})
	engineering := departments[0].(map[string]interface{})
	employees := engineering["employees"].([]interface{})

	assert.Len(t, employees, 3, "Should have 3 employees after adding Charlie")

	bob := employees[1].(map[string]interface{})
	assert.Equal(t, "Senior Manager", bob["role"], "Bob should be promoted to Senior Manager")
}

// =============================================================================
// SPECIAL CHARACTERS TESTS
// =============================================================================

// TestSpecialCharacters demonstrates handling special characters in paths
func TestSpecialCharacters(t *testing.T) {
	// Document with special characters in keys
	doc := map[string]interface{}{
		"normal":     "value",
		"with~tilde": "tilde value",
		"with/slash": "slash value",
		"":           "empty key",
	}

	// Operations with escaped paths
	patch := []jsonpatch.Operation{
		// Access key with tilde (~ becomes ~0)
		map[string]interface{}{
			"op":    "replace",
			"path":  "/with~0tilde",
			"value": "updated tilde",
		},
		// Access key with slash (/ becomes ~1)
		map[string]interface{}{
			"op":    "replace",
			"path":  "/with~1slash",
			"value": "updated slash",
		},
		// Access empty key
		map[string]interface{}{
			"op":    "replace",
			"path":  "/",
			"value": "updated empty",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	require.NoError(t, err)

	resultJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	t.Logf("Special characters result:\n%s", string(resultJSON))

	// Verify the updates
	resultMap := result.Doc
	assert.Equal(t, "updated tilde", resultMap["with~tilde"], "Tilde key should be updated")
	assert.Equal(t, "updated slash", resultMap["with/slash"], "Slash key should be updated")
	assert.Equal(t, "updated empty", resultMap[""], "Empty key should be updated")
}

// =============================================================================
// ERROR HANDLING TESTS
// =============================================================================

// TestErrorHandling demonstrates proper error handling
func TestErrorHandling(t *testing.T) {
	doc := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
	}

	// Patch with intentional error
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":   "remove",
			"path": "/user/nonexistent", // This will fail
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)

	assert.Error(t, err, "Expected error for nonexistent path")
	assert.Nil(t, result, "Result should be nil on error")
	t.Logf("Expected error: %v", err)
}

// TestApplyPatch_Errors tests error handling for different input types
func TestApplyPatch_Errors(t *testing.T) {
	t.Run("invalid JSON bytes", func(t *testing.T) {
		invalidJSON := []byte(`{invalid json}`)
		patch := []jsonpatch.Operation{
			map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(invalidJSON, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse JSON bytes")
	})

	t.Run("invalid JSON string", func(t *testing.T) {
		invalidJSON := `{invalid json}`
		patch := []jsonpatch.Operation{
			map[string]interface{}{"op": "replace", "path": "/name", "value": "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(invalidJSON, patch)
		assert.Error(t, err)
		// Note: Invalid JSON strings are now treated as primitive strings,
		// so the error comes from trying to apply path operations to a string
		assert.Contains(t, err.Error(), "invalid path")
	})

	t.Run("invalid patch operation", func(t *testing.T) {
		doc := map[string]any{"name": "John"}
		patch := []jsonpatch.Operation{
			map[string]interface{}{"op": "invalid", "path": "/name", "value": "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)
	})
}

// =============================================================================
// EXAMPLE TESTS
// =============================================================================

// Example demonstrates basic JSON Patch operations
func Example() {
	// Original document
	doc := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
			"age":   25,
		},
		"settings": map[string]interface{}{
			"theme": "dark",
		},
	}

	// Create patch operations
	patch := []jsonpatch.Operation{
		// Add a new field
		map[string]interface{}{
			"op":    "add",
			"path":  "/user/active",
			"value": true,
		},
		// Update existing field
		map[string]interface{}{
			"op":    "replace",
			"path":  "/user/age",
			"value": 26,
		},
		// Add to settings
		map[string]interface{}{
			"op":    "add",
			"path":  "/settings/notifications",
			"value": true,
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print result
	resultJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Println(string(resultJSON))

	// Output:
	// {
	//   "settings": {
	//     "notifications": true,
	//     "theme": "dark"
	//   },
	//   "user": {
	//     "active": true,
	//     "age": 26,
	//     "email": "alice@example.com",
	//     "name": "Alice"
	//   }
	// }
}

// =============================================================================
// FUZZ TESTS
// =============================================================================

// FuzzOperationSequence performs fuzz testing on operation sequences
func FuzzOperationSequence(f *testing.F) {
	// Seed with some basic operation sequences
	seeds := []string{
		`[{"op":"add","path":"/test","value":"hello"}]`,
		`[{"op":"remove","path":"/test"}]`,
		`[{"op":"replace","path":"/test","value":"world"}]`,
		`[{"op":"test","path":"/test","value":"hello"}]`,
		`[{"op":"copy","from":"/test","path":"/copy"}]`,
		`[{"op":"move","from":"/test","path":"/moved"}]`,
		`[{"op":"add","path":"/users/-","value":{"id":1}}]`,
		`[{"op":"add","path":"/users/0","value":{"id":2}}]`,
		`[{"op":"remove","path":"/users/0"}]`,
		`[{"op":"replace","path":"/users/0/name","value":"Alice"}]`,
		// Complex sequences
		`[{"op":"add","path":"/a","value":1},{"op":"add","path":"/b","value":2}]`,
		`[{"op":"add","path":"/test","value":"hello"},{"op":"test","path":"/test","value":"hello"}]`,
		`[{"op":"add","path":"/src","value":"data"},{"op":"copy","from":"/src","path":"/dst"}]`,
		`[{"op":"add","path":"/src","value":"data"},{"op":"move","from":"/src","path":"/dst"}]`,
		// Edge cases
		`[{"op":"add","path":"","value":{"root":true}}]`,
		`[{"op":"test","path":"","value":null}]`,
		`[{"op":"add","path":"/~0","value":"tilde"}]`,
		`[{"op":"add","path":"/~1","value":"slash"}]`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, patchJSON string) {
		// Skip obviously invalid JSON
		if !json.Valid([]byte(patchJSON)) {
			t.Skip("Invalid JSON")
		}

		// Parse patch operations
		var operations []jsonpatch.Operation
		if err := json.Unmarshal([]byte(patchJSON), &operations); err != nil {
			t.Skip("Cannot unmarshal operations")
		}

		// Skip empty operations
		if len(operations) == 0 {
			t.Skip("Empty operations")
		}

		// Create a test document
		doc := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{"id": 1, "name": "Alice"},
				map[string]interface{}{"id": 2, "name": "Bob"},
			},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
			"test": "value",
		}

		// Apply patch - should not panic
		result, err := jsonpatch.ApplyPatch(doc, operations)

		// If successful, verify result is valid JSON
		if err == nil && result != nil {
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				t.Errorf("Result is not valid JSON: %v", jsonErr)
			}

			// Verify original document is unchanged (immutability)
			if reflect.DeepEqual(doc, result.Doc) && len(operations) > 0 {
				// Check if operations actually modify the document
				hasModifyingOp := false
				for _, op := range operations {
					if opType, exists := op["op"]; exists && opType != "test" {
						hasModifyingOp = true
						break
					}
				}
				if hasModifyingOp {
					t.Logf("Warning: Document unchanged after modifying operations")
				}
			}
		}

		// Test that errors are properly structured
		if err != nil {
			// Verify error has proper structure
			if err.Error() == "" {
				t.Errorf("Error missing message")
			}
		}
	})
}

// FuzzJSONPointerPaths performs fuzz testing on JSON Pointer paths
func FuzzJSONPointerPaths(f *testing.F) {
	// Seed with various path patterns
	seeds := []string{
		"",
		"/",
		"/test",
		"/users/0",
		"/users/-",
		"/users/0/name",
		"/settings/theme",
		"/~0",
		"/~1",
		"/~0~1",
		"/a/b/c/d/e",
		"/123",
		"/0",
		"/-1",
		"/very/deep/nested/path/structure",
		"/with spaces",
		"/with\"quotes",
		"/with\nnewlines",
		"/with\ttabs",
		"/unicode🚀",
		"/empty/",
		"//double//slash",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, path string) {
		// Create test document
		doc := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{"id": 1, "name": "Alice"},
				map[string]interface{}{"id": 2, "name": "Bob"},
			},
			"settings": map[string]interface{}{
				"theme": "dark",
				"lang":  "en",
			},
			"test":        "value",
			"~0":          "tilde",
			"~1":          "slash",
			"with spaces": "spaces",
			"unicode🚀":    "rocket",
			"123":         "number",
			"":            "empty",
		}

		operations := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  path,
				"value": nil,
			},
		}

		// Apply patch - should not panic
		_, err := jsonpatch.ApplyPatch(doc, operations)

		// We don't care about the result, just that it doesn't panic
		// and errors are properly structured
		if err != nil {
			if err.Error() == "" {
				t.Errorf("Error missing message for path: %q", path)
			}
		}
	})
}

// FuzzOperationValues performs fuzz testing on operation values
func FuzzOperationValues(f *testing.F) {
	// Seed with various value types
	seeds := []string{
		`null`,
		`true`,
		`false`,
		`0`,
		`-1`,
		`123`,
		`-456`,
		`3.14`,
		`-2.71`,
		`""`,
		`"hello"`,
		`"with\"quotes"`,
		`"with\nnewlines"`,
		`"unicode🚀"`,
		`[]`,
		`[1,2,3]`,
		`["a","b","c"]`,
		`{}`,
		`{"key":"value"}`,
		`{"nested":{"deep":true}}`,
		`{"array":[1,2,3],"object":{"key":"value"}}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, valueJSON string) {
		// Skip obviously invalid JSON
		if !json.Valid([]byte(valueJSON)) {
			t.Skip("Invalid JSON")
		}

		// Parse the value
		var value interface{}
		if err := json.Unmarshal([]byte(valueJSON), &value); err != nil {
			t.Skip("Cannot unmarshal value")
		}

		// Create test document
		doc := map[string]interface{}{
			"test": "original",
		}

		// Test add operation
		operations := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "/fuzzed",
				"value": value,
			},
		}

		// Apply patch - should not panic
		result, err := jsonpatch.ApplyPatch(doc, operations)

		if err == nil && result != nil {
			// Verify the added value can be serialized back to JSON
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				t.Errorf("Result with fuzzed value is not valid JSON: %v", jsonErr)
			}

			// Verify the value was actually added
			resultMap := result.Doc
			if fuzzedValue, exists := resultMap["fuzzed"]; exists {
				// The value should be equivalent (though not necessarily identical due to JSON round-trip)
				fuzzedJSON, _ := json.Marshal(fuzzedValue)
				if !json.Valid(fuzzedJSON) {
					t.Errorf("Fuzzed value in result is not valid JSON")
				}
			}
		}

		// Test replace operation
		operations = []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/test",
				"value": value,
			},
		}

		// Apply patch - should not panic
		result, err = jsonpatch.ApplyPatch(doc, operations)

		if err == nil && result != nil {
			// Verify result is valid JSON
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				t.Errorf("Result with replaced value is not valid JSON: %v", jsonErr)
			}
		}
	})
}

// FuzzArrayIndices performs fuzz testing on array indices
func FuzzArrayIndices(f *testing.F) {
	// Seed with various array indices
	seeds := []int{
		-1000, -100, -10, -1, 0, 1, 10, 100, 1000, 10000,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, index int) {
		// Create test document with array
		doc := map[string]interface{}{
			"array": []interface{}{1, 2, 3, 4, 5},
		}

		// Test various operations with the fuzzed index
		operations := [][]jsonpatch.Operation{
			// Add operation
			{
				map[string]interface{}{
					"op":    "add",
					"path":  fmt.Sprintf("/array/%d", index),
					"value": "fuzzed",
				},
			},
			// Remove operation
			{
				map[string]interface{}{
					"op":   "remove",
					"path": fmt.Sprintf("/array/%d", index),
				},
			},
			// Replace operation
			{
				map[string]interface{}{
					"op":    "replace",
					"path":  fmt.Sprintf("/array/%d", index),
					"value": "replaced",
				},
			},
			// Test operation
			{
				map[string]interface{}{
					"op":    "test",
					"path":  fmt.Sprintf("/array/%d", index),
					"value": nil,
				},
			},
		}

		for _, ops := range operations {
			// Apply patch - should not panic
			_, err := jsonpatch.ApplyPatch(doc, ops)

			// We expect many of these to fail (invalid indices), but they shouldn't panic
			if err != nil {
				if err.Error() == "" {
					t.Errorf("Error missing message for index: %d", index)
				}
			}
		}
	})
}

// FuzzComplexDocuments performs fuzz testing on complex document structures
func FuzzComplexDocuments(f *testing.F) {
	// Seed with various document structures
	seeds := []string{
		`{}`,
		`[]`,
		`null`,
		`{"key":"value"}`,
		`[1,2,3]`,
		`{"array":[1,2,3],"object":{"nested":true}}`,
		`{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}]}`,
		`{"deeply":{"nested":{"structure":{"with":{"many":{"levels":true}}}}}}`,
		`{"mixed":[{"type":"object"},["nested","array"],42,true,null]}`,
		`{"empty":{},"emptyArray":[],"null":null,"boolean":true,"number":42}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, docJSON string) {
		// Skip obviously invalid JSON
		if !json.Valid([]byte(docJSON)) {
			t.Skip("Invalid JSON")
		}

		// Parse the document
		var doc interface{}
		if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
			t.Skip("Cannot unmarshal document")
		}

		// Skip null documents as they cause issues with reflection
		if doc == nil {
			t.Skip("Null document")
		}

		// Test basic operations on the fuzzed document
		operations := [][]jsonpatch.Operation{
			// Test root
			{
				map[string]interface{}{
					"op":    "test",
					"path":  "",
					"value": doc,
				},
			},
			// Replace root
			{
				map[string]interface{}{
					"op":    "replace",
					"path":  "",
					"value": "replaced",
				},
			},
			// Add to root (if it's an object)
			{
				map[string]interface{}{
					"op":    "add",
					"path":  "/fuzzed",
					"value": "added",
				},
			},
		}

		for _, ops := range operations {
			// Apply patch - should not panic
			result, err := jsonpatch.ApplyPatch(doc, ops)

			if err == nil && result != nil {
				// Verify result is valid JSON
				_, jsonErr := json.Marshal(result.Doc)
				if jsonErr != nil {
					t.Errorf("Result is not valid JSON: %v", jsonErr)
				}
			}

			// Verify errors are properly structured
			if err != nil {
				if err.Error() == "" {
					t.Errorf("Error missing message")
				}
			}
		}
	})
}

// FuzzEdgeCases performs fuzz testing on edge cases and special conditions
func FuzzEdgeCases(f *testing.F) {
	// Seed with edge case scenarios
	seeds := []struct {
		doc   string
		patch string
	}{
		// Empty documents
		{`{}`, `[{"op":"add","path":"/test","value":"hello"}]`},
		{`[]`, `[{"op":"add","path":"/-","value":"hello"}]`},
		{`null`, `[{"op":"replace","path":"","value":"hello"}]`},

		// Special characters in keys
		{`{"~0":"tilde","~1":"slash"}`, `[{"op":"test","path":"/~00","value":"tilde"}]`},
		{`{"":"empty"}`, `[{"op":"test","path":"/","value":"empty"}]`},
		{`{" ":"space"}`, `[{"op":"test","path":"/ ","value":"space"}]`},

		// Deeply nested structures
		{`{"a":{"b":{"c":{"d":{"e":"deep"}}}}}`, `[{"op":"test","path":"/a/b/c/d/e","value":"deep"}]`},

		// Large arrays
		{`{"array":[0,1,2,3,4,5,6,7,8,9]}`, `[{"op":"add","path":"/array/-","value":10}]`},

		// Mixed types
		{`{"mixed":[1,"two",true,null,{"five":5}]}`, `[{"op":"test","path":"/mixed/4/five","value":5}]`},

		// Circular references in operations
		{`{"src":"data"}`, `[{"op":"copy","from":"/src","path":"/src"}]`},
		{`{"a":"data"}`, `[{"op":"move","from":"/a","path":"/a"}]`},
	}

	for _, seed := range seeds {
		f.Add(seed.doc, seed.patch)
	}

	f.Fuzz(func(t *testing.T, docJSON, patchJSON string) {
		// Skip obviously invalid JSON
		if !json.Valid([]byte(docJSON)) || !json.Valid([]byte(patchJSON)) {
			t.Skip("Invalid JSON")
		}

		// Parse document and patch
		var doc interface{}
		var operations []jsonpatch.Operation

		if err := json.Unmarshal([]byte(docJSON), &doc); err != nil {
			t.Skip("Cannot unmarshal document")
		}

		// Skip null documents as they cause issues with reflection
		if doc == nil {
			t.Skip("Null document")
		}

		if err := json.Unmarshal([]byte(patchJSON), &operations); err != nil {
			t.Skip("Cannot unmarshal operations")
		}

		if len(operations) == 0 {
			t.Skip("Empty operations")
		}

		// Apply patch - should not panic
		result, err := jsonpatch.ApplyPatch(doc, operations)

		// Verify no panics occurred and results are consistent
		if err == nil && result != nil {
			// Result should be valid JSON
			resultJSON, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				t.Errorf("Result is not valid JSON: %v", jsonErr)
			}

			// Result should be parseable back
			var reparsed interface{}
			if json.Unmarshal(resultJSON, &reparsed) != nil {
				t.Errorf("Result cannot be reparsed from JSON")
			}
		}

		// Errors should be properly structured
		if err != nil {
			if err.Error() == "" {
				t.Errorf("Error missing message")
			}
		}
	})
}
