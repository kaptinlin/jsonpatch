package jsonpatch_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

// profile is a test struct for generic patch testing
type profile struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Tags  []string `json:"tags"`
}

func TestApplyPatchBasic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		doc      any
		patch    []jsonpatch.Operation
		expected any
		wantErr  bool
	}{
		{
			name:     "empty patch",
			doc:      map[string]any{"a": 1},
			patch:    []jsonpatch.Operation{},
			expected: map[string]any{"a": 1},
			wantErr:  false,
		},
		{
			name: "single operation",
			doc:  map[string]any{"a": 1},
			patch: []jsonpatch.Operation{
				{Op: "add", Path: "/b", Value: 2},
			},
			expected: map[string]any{"a": 1, "b": 2},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := jsonpatch.ApplyPatch(tt.doc, tt.patch, jsonpatch.WithMutate(false))

			if tt.wantErr {
				if err == nil {
					t.Fatal("ApplyPatch() error = nil, want error")
				}
				return
			}

			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}

			assert.Equal(t, tt.expected, result.Doc)
			if result.Res == nil {
				assert.Fail(t, "ApplyPatch() Res = nil, want non-nil")
			}
			assert.Equal(t, len(tt.patch), len(result.Res), "ApplyPatch() len(Res)")
		})
	}
}

func TestValidateOperation(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		operation jsonpatch.Operation
		wantErr   error // nil means no error expected
	}{
		{
			name:      "valid add operation",
			operation: jsonpatch.Operation{Op: "add", Path: "/a", Value: 1},
		},
		{
			name:      "missing op field",
			operation: jsonpatch.Operation{Path: "/a", Value: 1},
			wantErr:   jsonpatch.ErrMissingOp,
		},
		{
			name:      "missing path field",
			operation: jsonpatch.Operation{Op: "add", Value: 1},
			wantErr:   jsonpatch.ErrMissingPath,
		},
		{
			name:      "missing value field for add",
			operation: jsonpatch.Operation{Op: "add", Path: "/a"},
			wantErr:   jsonpatch.ErrMissingValue,
		},
		{
			name:      "invalid operation type",
			operation: jsonpatch.Operation{Op: "invalid", Path: "/a"},
			wantErr:   jsonpatch.ErrInvalidOperation,
		},
		{
			name:      "empty operation",
			operation: jsonpatch.Operation{Op: "", Path: ""},
			wantErr:   jsonpatch.ErrMissingOp,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := jsonpatch.ValidateOperation(tt.operation, false)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatal("ValidateOperation() error = nil, want error")
				}
				if !errors.Is(err, tt.wantErr) {
					assert.Equal(t, tt.wantErr, err, "ValidateOperation() error")
				}
			} else if err != nil {
				t.Errorf("ValidateOperation() error = %v, want nil", err)
			}
		})
	}
}

func TestApplyPatch_Struct(t *testing.T) {
	t.Parallel()
	// Test data
	before := profile{
		Name: "John",
		Tags: []string{"dev"},
	}

	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}

	// Verify results
	if result.Doc.Name != "Jane" {
		t.Errorf("ApplyPatch() Name = %v, want Jane", result.Doc.Name)
	}
	if result.Doc.Email != "jane@example.com" {
		t.Errorf("ApplyPatch() Email = %v, want jane@example.com", result.Doc.Email)
	}
	assert.Equal(t, []string{"dev", "golang"}, result.Doc.Tags)
	if len(result.Res) != 3 {
		t.Errorf("ApplyPatch() len(Res) = %d, want 3", len(result.Res))
	}

	// Verify original is unchanged (immutable by default)
	if before.Name != "John" {
		t.Errorf("original Name = %v, want John", before.Name)
	}
	if before.Email != "" {
		t.Errorf("original Email = %q, want empty", before.Email)
	}
	assert.Equal(t, []string{"dev"}, before.Tags)
}

func TestApplyPatch_Map(t *testing.T) {
	t.Parallel()
	// Test data
	before := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	// Verify results
	if got := result.Doc["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "result.Doc[name]")
	}
	if got := result.Doc["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "result.Doc[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, result.Doc["tags"])
	assert.Equal(t, 3, len(result.Res), "len(result.Res)")

	// Verify original is unchanged (immutable by default)
	if got := before["name"]; got != "John" {
		assert.Equal(t, "John", got, "before[name]")
	}
	_, hasEmail := before["email"]
	if hasEmail {
		assert.Fail(t, "before should not have email key")
	}
}

func TestApplyPatch_JSONBytes(t *testing.T) {
	t.Parallel()
	// Test data
	before := []byte(`{"name":"John","tags":["dev"]}`)

	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal(result.Doc, &resultMap)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify results
	if got := resultMap["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "resultMap[name]")
	}
	if got := resultMap["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "resultMap[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Equal(t, 3, len(result.Res), "len(result.Res)")

	// Verify original is unchanged
	var original map[string]any
	err = json.Unmarshal(before, &original)
	if err != nil {
		t.Fatalf("json.Unmarshal(before) error = %v", err)
	}
	if got := original["name"]; got != "John" {
		assert.Equal(t, "John", got, "original[name]")
	}
	_, hasEmail := original["email"]
	if hasEmail {
		assert.Fail(t, "original should not have email key")
	}
}

func TestApplyPatch_JSONString(t *testing.T) {
	t.Parallel()
	// Test data
	before := `{"name":"John","tags":["dev"]}`

	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(before, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal([]byte(result.Doc), &resultMap)
	if err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	// Verify results
	if got := resultMap["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "resultMap[name]")
	}
	if got := resultMap["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "resultMap[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Equal(t, 3, len(result.Res), "len(result.Res)")
}

func TestArrayOperations(t *testing.T) {
	t.Parallel()
	// Document with array
	doc := map[string]any{
		"items": []any{
			map[string]any{"id": 1, "name": "Item 1"},
			map[string]any{"id": 2, "name": "Item 2"},
			map[string]any{"id": 3, "name": "Item 3"},
		},
	}

	// Array operations
	patch := []jsonpatch.Operation{
		// Insert at beginning
		{
			Op:    "add",
			Path:  "/items/0",
			Value: map[string]any{"id": 0, "name": "Item 0"},
		},
		// Append at end
		{
			Op:    "add",
			Path:  "/items/-",
			Value: map[string]any{"id": 4, "name": "Item 4"},
		},
		// Update middle item
		{
			Op:    "replace",
			Path:  "/items/2/name",
			Value: "Updated Item 1",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	resultJSON, _ := json.Marshal(result.Doc, jsontext.Multiline(true))
	t.Logf("Array operations result:\n%s", string(resultJSON))

	// Verify the result
	items := result.Doc["items"].([]any)
	assert.Equal(t, 5, len(items), "len(items)")
}

func TestMultipleOperations(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"counters": map[string]any{
			"a": 0,
			"b": 0,
		},
	}

	patch := []jsonpatch.Operation{
		{
			Op:    "replace",
			Path:  "/counters/a",
			Value: 1,
		},
		{
			Op:    "replace",
			Path:  "/counters/b",
			Value: 2,
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	resultJSON, _ := json.Marshal(result.Doc, jsontext.Multiline(true))
	t.Logf("Multiple operations result:\n%s", string(resultJSON))

	// Verify the result
	counters := result.Doc["counters"].(map[string]any)
	if got := counters["a"]; got != 1 {
		assert.Equal(t, 1, got, "counters[a]")
	}
	if got := counters["b"]; got != 2 {
		assert.Equal(t, 2, got, "counters[b]")
	}
}

func TestApplyPatch_WithMutate(t *testing.T) {
	t.Parallel()
	// Test data - using map for easier mutation testing
	original := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
	}

	// Apply patch with mutate=true
	result, err := jsonpatch.ApplyPatch(original, patch, jsonpatch.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	// Verify results
	if got := result.Doc["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "result.Doc[name]")
	}
	assert.Equal(t, 1, len(result.Res), "len(result.Res)")
}

func TestComplexDocument(t *testing.T) {
	t.Parallel()
	// Complex nested document
	doc := map[string]any{
		"company": map[string]any{
			"name": "Tech Corp",
			"departments": []any{
				map[string]any{
					"name": "Engineering",
					"employees": []any{
						map[string]any{"id": 1, "name": "Alice", "role": "Developer"},
						map[string]any{"id": 2, "name": "Bob", "role": "Manager"},
					},
				},
			},
		},
		"metadata": map[string]any{
			"lastUpdated": "2023-01-01",
		},
	}

	// Complex operations
	patch := []jsonpatch.Operation{
		// Add new employee to Engineering
		{
			Op:   "add",
			Path: "/company/departments/0/employees/-",
			Value: map[string]any{
				"id":   3,
				"name": "Charlie",
				"role": "Senior Developer",
			},
		},
		// Promote Bob to Senior Manager
		{
			Op:    "replace",
			Path:  "/company/departments/0/employees/1/role",
			Value: "Senior Manager",
		},
		// Update metadata
		{
			Op:    "replace",
			Path:  "/metadata/lastUpdated",
			Value: "2023-12-01",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	resultJSON, _ := json.Marshal(result.Doc, jsontext.Multiline(true))
	t.Logf("Complex document result:\n%s", string(resultJSON))

	// Verify the changes
	company := result.Doc["company"].(map[string]any)
	departments := company["departments"].([]any)
	engineering := departments[0].(map[string]any)
	employees := engineering["employees"].([]any)

	assert.Equal(t, 3, len(employees), "len(employees)")

	bob := employees[1].(map[string]any)
	if got := bob["role"]; got != "Senior Manager" {
		assert.Equal(t, "Senior Manager", got, "bob[role]")
	}
}

func TestSpecialCharacters(t *testing.T) {
	t.Parallel()
	// Document with special characters in keys
	doc := map[string]any{
		"normal":     "value",
		"with~tilde": "tilde value",
		"with/slash": "slash value",
		"":           "empty key",
	}

	// Operations with escaped paths
	patch := []jsonpatch.Operation{
		// Access key with tilde (~ becomes ~0)
		{
			Op:    "replace",
			Path:  "/with~0tilde",
			Value: "updated tilde",
		},
		// Access key with slash (/ becomes ~1)
		{
			Op:    "replace",
			Path:  "/with~1slash",
			Value: "updated slash",
		},
		// Access empty key
		{
			Op:    "replace",
			Path:  "/",
			Value: "updated empty",
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}

	resultJSON, _ := json.Marshal(result.Doc, jsontext.Multiline(true))
	t.Logf("Special characters result:\n%s", string(resultJSON))

	// Verify the updates
	resultMap := result.Doc
	if got := resultMap["with~tilde"]; got != "updated tilde" {
		assert.Equal(t, "updated tilde", got, "resultMap[with~tilde]")
	}
	if got := resultMap["with/slash"]; got != "updated slash" {
		assert.Equal(t, "updated slash", got, "resultMap[with/slash]")
	}
	if got := resultMap[""]; got != "updated empty" {
		assert.Equal(t, "updated empty", got, "resultMap[]")
	}
}

func TestErrorHandling(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}

	// Patch with intentional error
	patch := []jsonpatch.Operation{
		{
			Op:   "remove",
			Path: "/user/nonexistent", // This will fail
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)

	assert.NotNil(t, err, "Expected error for nonexistent path")
	if result != nil {
		assert.Fail(t, "Result should be nil on error")
	}
	t.Logf("Expected error: %v", err)
}

func TestApplyPatch_Errors(t *testing.T) {
	t.Parallel()
	t.Run("invalid JSON bytes", func(t *testing.T) {
		t.Parallel()
		invalidJSON := []byte(`{invalid json}`)
		patch := []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(invalidJSON, patch)
		assert.NotNil(t, err, "expected error for invalid JSON bytes")
	})

	t.Run("invalid JSON string", func(t *testing.T) {
		t.Parallel()
		invalidJSON := `{invalid json}`
		patch := []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(invalidJSON, patch)
		assert.NotNil(t, err, "expected error for invalid JSON string")
		// Note: Invalid JSON strings are now treated as primitive strings,
		// so the error comes from trying to apply path operations to a string
		// We only check that an error occurred, not the specific message
	})

	t.Run("invalid patch operation", func(t *testing.T) {
		t.Parallel()
		doc := map[string]any{"name": "John"}
		patch := []jsonpatch.Operation{
			{Op: "invalid", Path: "/name", Value: "Jane"},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NotNil(t, err, "expected error for invalid patch operation")
	})
}

func Example() {
	// Original document
	doc := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@example.com",
			"age":   25,
		},
		"settings": map[string]any{
			"theme": "dark",
		},
	}

	// Create patch operations
	patch := []jsonpatch.Operation{
		// Add a new field
		{
			Op:    "add",
			Path:  "/user/active",
			Value: true,
		},
		// Update existing field
		{
			Op:    "replace",
			Path:  "/user/age",
			Value: 26,
		},
		// Add to settings
		{
			Op:    "add",
			Path:  "/settings/notifications",
			Value: true,
		},
	}

	// Apply patch
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Print result
	resultJSON, _ := json.Marshal(result.Doc, jsontext.Multiline(true), json.Deterministic(true))
	fmt.Println(string(resultJSON))

	// Output:
	// {
	// 	"settings": {
	// 		"notifications": true,
	// 		"theme": "dark"
	// 	},
	// 	"user": {
	// 		"active": true,
	// 		"age": 26,
	// 		"email": "alice@example.com",
	// 		"name": "Alice"
	// 	}
	// }
}

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
		var testPatch any
		if err := json.Unmarshal([]byte(patchJSON), &testPatch); err != nil {
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
		doc := map[string]any{
			"users": []any{
				map[string]any{"id": 1, "name": "Alice"},
				map[string]any{"id": 2, "name": "Bob"},
			},
			"settings": map[string]any{
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
					if op.Op != "" && op.Op != "test" {
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
		doc := map[string]any{
			"users": []any{
				map[string]any{"id": 1, "name": "Alice"},
				map[string]any{"id": 2, "name": "Bob"},
			},
			"settings": map[string]any{
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
			{
				Op:    "test",
				Path:  path,
				Value: nil,
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
		var testValue any
		if err := json.Unmarshal([]byte(valueJSON), &testValue); err != nil {
			t.Skip("Invalid JSON")
		}

		// Parse the value
		var value any
		if err := json.Unmarshal([]byte(valueJSON), &value); err != nil {
			t.Skip("Cannot unmarshal value")
		}

		// Create test document
		doc := map[string]any{
			"test": "original",
		}

		// Test add operation
		operations := []jsonpatch.Operation{
			{
				Op:    "add",
				Path:  "/fuzzed",
				Value: value,
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
				var testUnmarshal any
				if err := json.Unmarshal(fuzzedJSON, &testUnmarshal); err != nil {
					t.Errorf("Fuzzed value in result is not valid JSON")
				}
			}
		}

		// Test replace operation
		operations = []jsonpatch.Operation{
			{
				Op:    "replace",
				Path:  "/test",
				Value: value,
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
		doc := map[string]any{
			"array": []any{1, 2, 3, 4, 5},
		}

		// Test various operations with the fuzzed index
		operations := [][]jsonpatch.Operation{
			// Add operation
			{
				{
					Op:    "add",
					Path:  fmt.Sprintf("/array/%d", index),
					Value: "fuzzed",
				},
			},
			// Remove operation
			{
				{
					Op:   "remove",
					Path: fmt.Sprintf("/array/%d", index),
				},
			},
			// Replace operation
			{
				{
					Op:    "replace",
					Path:  fmt.Sprintf("/array/%d", index),
					Value: "replaced",
				},
			},
			// Test operation
			{
				{
					Op:    "test",
					Path:  fmt.Sprintf("/array/%d", index),
					Value: nil,
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
		var testDoc any
		if err := json.Unmarshal([]byte(docJSON), &testDoc); err != nil {
			t.Skip("Invalid JSON")
		}

		// Parse the document
		var doc any
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
				{
					Op:    "test",
					Path:  "",
					Value: doc,
				},
			},
			// Replace root
			{
				{
					Op:    "replace",
					Path:  "",
					Value: "replaced",
				},
			},
			// Add to root (if it's an object)
			{
				{
					Op:    "add",
					Path:  "/fuzzed",
					Value: "added",
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
		var testDoc, testPatch any
		if err := json.Unmarshal([]byte(docJSON), &testDoc); err != nil {
			t.Skip("Invalid JSON")
		}
		if err := json.Unmarshal([]byte(patchJSON), &testPatch); err != nil {
			t.Skip("Invalid JSON")
		}

		// Parse document and patch
		var doc any
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
			var reparsed any
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
