package jsonpatch_test

import (
	"fmt"
	"reflect"
	"testing"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/op"
)

// profile is a test struct for generic patch testing
type profile struct {
	Name  string   `json:"name"`
	Email string   `json:"email,omitempty"`
	Tags  []string `json:"tags"`
}

func TestApplyBasic(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		doc      any
		patch    []jsoncodec.Operation
		expected any
		wantErr  bool
	}{
		{
			name:     "empty patch",
			doc:      map[string]any{"a": 1},
			patch:    []jsoncodec.Operation{},
			expected: map[string]any{"a": 1},
			wantErr:  false,
		},
		{
			name: "single operation",
			doc:  map[string]any{"a": 1},
			patch: []jsoncodec.Operation{
				{Op: "add", Path: "/b", Value: 2},
			},
			expected: map[string]any{"a": 1, "b": 2},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result, err := applyOperations(t, tt.doc, tt.patch)

			if tt.wantErr {
				if err == nil {
					require.FailNow(t, "Apply() error = nil, want error")
				}
				return
			}

			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}

			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, len(tt.patch), len(result.Steps), "Apply() len(Steps)")
		})
	}
}

func TestCompileOperationsBasic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation []jsoncodec.Operation
		wantErr   error
	}{
		{
			name:      "valid add operation",
			operation: []jsoncodec.Operation{{Op: "add", Path: "/a", Value: 1}},
		},
		{
			name:      "valid remove operation",
			operation: []jsoncodec.Operation{{Op: "remove", Path: "/a"}},
		},
		{
			name:      "valid inc operation",
			operation: []jsoncodec.Operation{{Op: "inc", Path: "/a", Inc: 1}},
		},
		{
			name:      "valid extend operation",
			operation: []jsoncodec.Operation{{Op: "extend", Path: "/a", Props: map[string]any{"ok": true}}},
		},
		{
			name:      "valid split operation",
			operation: []jsoncodec.Operation{{Op: "split", Path: "/a", Pos: 1}},
		},
		{
			name:      "missing op field",
			operation: []jsoncodec.Operation{{Path: "/a", Value: 1}},
			wantErr:   jsonpatch.ErrPayloadInvalid,
		},
		{
			name:      "root add operation",
			operation: []jsoncodec.Operation{{Op: "add", Path: "", Value: 1}},
		},
		{
			name:      "nil value field for add",
			operation: []jsoncodec.Operation{{Op: "add", Path: "/a", Value: nil}},
		},
		{
			name:      "invalid operation type",
			operation: []jsoncodec.Operation{{Op: "invalid", Path: "/a"}},
			wantErr:   jsonpatch.ErrPayloadInvalid,
		},
		{
			name:      "empty operation",
			operation: []jsoncodec.Operation{{Op: "", Path: ""}},
			wantErr:   jsonpatch.ErrPayloadInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := jsonpatch.CompileOperations(tt.operation, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))

			if tt.wantErr != nil {
				if err == nil {
					require.FailNow(t, "CompileOperations() error = nil, want error")
				}
				assert.ErrorIs(t, err, tt.wantErr)
			} else if err != nil {
				assert.Fail(t, fmt.Sprintf("CompileOperations() error = %v, want nil", err))
			}
		})
	}
}

func TestOperationInterfacesExposeCanonicalMethods(t *testing.T) {
	t.Parallel()

	add := op.NewAdd([]string{"items", "-"}, "x")
	require.NoError(t, add.Validate())
	assert.Equal(t, jsonpatch.OpAddType, add.Op())
	assert.Equal(t, int(compact.CodeAdd), add.Code())
	assert.Equal(t, []string{"items", "-"}, add.Path())

	applied, err := add.Apply(map[string]any{"items": []any{"a"}})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"items": []any{"a", "x"}}, applied.Doc)

	jsonOp, err := add.ToJSON()
	require.NoError(t, err)
	assert.Equal(t, jsoncodec.Operation{Op: "add", Path: "/items/-", Value: "x"}, jsonOp)

	compactOp, err := add.ToCompact()
	require.NoError(t, err)
	assert.Equal(t, []any{int(compact.CodeAdd), []string{"items", "-"}, "x"}, []any(compactOp))

	predicate := op.NewTestWithNot([]string{"active"}, true, true)
	ok, err := predicate.Test(map[string]any{"active": false})
	require.NoError(t, err)
	assert.True(t, ok)
	assert.True(t, predicate.Not())

	secondOrder := op.NewNot(op.NewDefined([]string{"active"}))
	assert.True(t, secondOrder.Not())
	assert.Len(t, secondOrder.Ops(), 1)
	assert.Equal(t, []string{"active"}, secondOrder.Ops()[0].Path())
}

func TestApply_Struct(t *testing.T) {
	t.Parallel()
	// Test data
	before := profile{
		Name: "John",
		Tags: []string{"dev"},
	}

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := applyOperations(t, before, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
	}

	// Verify results
	if result.Doc.Name != "Jane" {
		assert.Fail(t, fmt.Sprintf("Apply() Name = %v, want Jane", result.Doc.Name))
	}
	if result.Doc.Email != "jane@example.com" {
		assert.Fail(t, fmt.Sprintf("Apply() Email = %v, want jane@example.com", result.Doc.Email))
	}
	assert.Equal(t, []string{"dev", "golang"}, result.Doc.Tags)
	if len(result.Steps) != 3 {
		assert.Fail(t, fmt.Sprintf("Apply() len(Steps) = %d, want 3", len(result.Steps)))
	}

	// Verify original is unchanged (immutable by default)
	if before.Name != "John" {
		assert.Fail(t, fmt.Sprintf("original Name = %v, want John", before.Name))
	}
	if before.Email != "" {
		assert.Fail(t, fmt.Sprintf("original Email = %q, want empty", before.Email))
	}
	assert.Equal(t, []string{"dev"}, before.Tags)
}

func TestApply_Map(t *testing.T) {
	t.Parallel()
	// Test data
	before := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := applyOperations(t, before, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
	}

	// Verify results
	if got := result.Doc["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "result.Doc[name]")
	}
	if got := result.Doc["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "result.Doc[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, result.Doc["tags"])
	assert.Equal(t, 3, len(result.Steps), "len(result.Steps)")

	// Verify original is unchanged (immutable by default)
	if got := before["name"]; got != "John" {
		assert.Equal(t, "John", got, "before[name]")
	}
	_, hasEmail := before["email"]
	if hasEmail {
		assert.Fail(t, "before should not have email key")
	}
}

func TestApply_JSONBytes(t *testing.T) {
	t.Parallel()
	// Test data
	before := []byte(`{"name":"John","tags":["dev"]}`)

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := applyOperations(t, before, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
	}

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal(result.Doc, &resultMap)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("json.Unmarshal() error = %v", err))
	}

	// Verify results
	if got := resultMap["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "resultMap[name]")
	}
	if got := resultMap["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "resultMap[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Equal(t, 3, len(result.Steps), "len(result.Steps)")

	// Verify original is unchanged
	var original map[string]any
	err = json.Unmarshal(before, &original)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("json.Unmarshal(before) error = %v", err))
	}
	if got := original["name"]; got != "John" {
		assert.Equal(t, "John", got, "original[name]")
	}
	_, hasEmail := original["email"]
	if hasEmail {
		assert.Fail(t, "original should not have email key")
	}
}

func TestApply_JSONString(t *testing.T) {
	t.Parallel()
	// Test data
	before := `{"name":"John","tags":["dev"]}`

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/tags/-", Value: "golang"},
		{Op: "add", Path: "/email", Value: "jane@example.com"},
	}

	// Apply patch
	result, err := applyOperations(t, jsonpatch.JSONText(before), patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
	}

	// Parse result to verify
	var resultMap map[string]any
	err = json.Unmarshal([]byte(result.Doc), &resultMap)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("json.Unmarshal() error = %v", err))
	}

	// Verify results
	if got := resultMap["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "resultMap[name]")
	}
	if got := resultMap["email"]; got != "jane@example.com" {
		assert.Equal(t, "jane@example.com", got, "resultMap[email]")
	}
	assert.Equal(t, []any{"dev", "golang"}, resultMap["tags"])
	assert.Equal(t, 3, len(result.Steps), "len(result.Steps)")
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
	patch := []jsoncodec.Operation{
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
	result, err := applyOperations(t, doc, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
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

	patch := []jsoncodec.Operation{
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
	result, err := applyOperations(t, doc, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
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

func TestApplyInPlace(t *testing.T) {
	t.Parallel()
	// Test data - using map for easier mutation testing
	original := map[string]any{
		"name": "John",
		"tags": []any{"dev"},
	}

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
	}

	// Apply patch in place.
	result, err := applyOperationsInPlace(t, original, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v", err))
	}

	// Verify results
	if got := result.Doc["name"]; got != "Jane" {
		assert.Equal(t, "Jane", got, "result.Doc[name]")
	}
	assert.Equal(t, "Jane", original["name"], "original[name]")
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
	patch := []jsoncodec.Operation{
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
	result, err := applyOperations(t, doc, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
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
	patch := []jsoncodec.Operation{
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
	result, err := applyOperations(t, doc, patch)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Apply() error = %v", err))
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
	patch := []jsoncodec.Operation{
		{
			Op:   "remove",
			Path: "/user/nonexistent", // This will fail
		},
	}

	// Apply patch
	result, err := applyOperations(t, doc, patch)

	assert.NotNil(t, err, "Expected error for nonexistent path")
	if result != nil {
		assert.Fail(t, "Result should be nil on error")
	}
	t.Logf("Expected error: %v", err)
}

func TestApply_Errors(t *testing.T) {
	t.Parallel()
	t.Run("invalid JSON bytes", func(t *testing.T) {
		t.Parallel()
		invalidJSON := []byte(`{invalid json}`)
		patch := []jsoncodec.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
		}

		_, err := applyOperations(t, invalidJSON, patch)
		assert.NotNil(t, err, "expected error for invalid JSON bytes")
	})

	t.Run("invalid JSON text", func(t *testing.T) {
		t.Parallel()
		invalidJSON := jsonpatch.JSONText(`{invalid json}`)
		patch := []jsoncodec.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
		}

		_, err := applyOperations(t, invalidJSON, patch)
		assert.NotNil(t, err, "expected error for invalid JSON text")
	})

	t.Run("invalid patch operation", func(t *testing.T) {
		t.Parallel()
		doc := map[string]any{"name": "John"}
		patch := []jsoncodec.Operation{
			{Op: "invalid", Path: "/name", Value: "Jane"},
		}

		_, err := applyOperations(t, doc, patch)
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
	patch := []jsoncodec.Operation{
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

	compiled, err := jsonpatch.CompileOperations(patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	result, err := jsonpatch.Apply(compiled, doc)
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
		var operations []jsoncodec.Operation
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

		patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
		if err != nil {
			return
		}

		// Apply patch - should not panic
		result, err := jsonpatch.Apply(patch, doc)

		// If successful, verify result is valid JSON
		if err == nil && result != nil {
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				assert.Fail(t, fmt.Sprintf("Result is not valid JSON: %v", jsonErr))
			}

			// Verify original document is unchanged (immutability)
			if reflect.DeepEqual(doc, result.Doc) && len(operations) > 0 {
				// Check if operations actually modify the document
				hasModifyingOp := false
				for i := range operations {
					if operations[i].Op != "" && operations[i].Op != "test" {
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
				assert.Fail(t, "Error missing message")
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

		operations := []jsoncodec.Operation{
			{
				Op:    "test",
				Path:  path,
				Value: nil,
			},
		}

		patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
		if err != nil {
			return
		}

		// Apply patch - should not panic
		_, err = jsonpatch.Apply(patch, doc)

		// We don't care about the result, just that it doesn't panic
		// and errors are properly structured
		if err != nil {
			if err.Error() == "" {
				assert.Fail(t, fmt.Sprintf("Error missing message for path: %q", path))
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
		operations := []jsoncodec.Operation{
			{
				Op:    "add",
				Path:  "/fuzzed",
				Value: value,
			},
		}

		// Apply patch - should not panic
		result, err := applyOperations(t, doc, operations)

		if err == nil && result != nil {
			// Verify the added value can be serialized back to JSON
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				assert.Fail(t, fmt.Sprintf("Result with fuzzed value is not valid JSON: %v", jsonErr))
			}

			// Verify the value was actually added
			resultMap := result.Doc
			if fuzzedValue, exists := resultMap["fuzzed"]; exists {
				// The value should be equivalent (though not necessarily identical due to JSON round-trip)
				fuzzedJSON, _ := json.Marshal(fuzzedValue)
				var testUnmarshal any
				if err := json.Unmarshal(fuzzedJSON, &testUnmarshal); err != nil {
					assert.Fail(t, "Fuzzed value in result is not valid JSON")
				}
			}
		}

		// Test replace operation
		operations = []jsoncodec.Operation{
			{
				Op:    "replace",
				Path:  "/test",
				Value: value,
			},
		}

		// Apply patch - should not panic
		result, err = applyOperations(t, doc, operations)

		if err == nil && result != nil {
			// Verify result is valid JSON
			_, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				assert.Fail(t, fmt.Sprintf("Result with replaced value is not valid JSON: %v", jsonErr))
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
		operations := [][]jsoncodec.Operation{
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
			_, err := applyOperations(t, doc, ops)

			// We expect many of these to fail (invalid indices), but they shouldn't panic
			if err != nil {
				if err.Error() == "" {
					assert.Fail(t, fmt.Sprintf("Error missing message for index: %d", index))
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
		operations := [][]jsoncodec.Operation{
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
			result, err := applyOperations(t, doc, ops)

			if err == nil && result != nil {
				// Verify result is valid JSON
				_, jsonErr := json.Marshal(result.Doc)
				if jsonErr != nil {
					assert.Fail(t, fmt.Sprintf("Result is not valid JSON: %v", jsonErr))
				}
			}

			// Verify errors are properly structured
			if err != nil {
				if err.Error() == "" {
					assert.Fail(t, "Error missing message")
				}
			}
		}
	})
}

func TestPublicHelpersAndOps(t *testing.T) {
	t.Parallel()

	t.Run("create matcher default honors ignore case and invalid patterns", func(t *testing.T) {
		t.Parallel()

		matcher := jsonpatch.CreateMatcherDefault("^hello$", true)
		assert.True(t, matcher("HELLO"))
		assert.False(t, matcher("HELLO!"))

		invalid := jsonpatch.CreateMatcherDefault("(", false)
		assert.False(t, invalid("anything"))
	})

	t.Run("compile operations covers public error paths", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name    string
			ops     []jsoncodec.Operation
			opts    []jsonpatch.CompileOption
			wantErr error
		}{
			{name: "nil patch", ops: nil, opts: allCapabilities()},
			{name: "empty patch", ops: []jsoncodec.Operation{}, opts: allCapabilities()},
			{
				name:    "matches disabled",
				ops:     []jsoncodec.Operation{{Op: "matches", Path: "/name", Value: "^a"}},
				wantErr: jsonpatch.ErrUnsupportedCapability,
			},
			{
				name: "matches allowed",
				ops:  []jsoncodec.Operation{{Op: "matches", Path: "/name", Value: "^a"}},
				opts: allCapabilities(),
			},
			{
				name: "copy from root",
				ops:  []jsoncodec.Operation{{Op: "copy", Path: "/name", From: ""}},
				opts: allCapabilities(),
			},
			{
				name:    "move into own children",
				ops:     []jsoncodec.Operation{{Op: "move", Path: "/a/b", From: "/a"}},
				opts:    allCapabilities(),
				wantErr: jsonpatch.ErrPayloadInvalid,
			},
			{
				name: "test type list validates",
				ops:  []jsoncodec.Operation{{Op: "test_type", Path: "/kind", Type: []any{"string", "number"}}},
				opts: allCapabilities(),
			},
			{
				name:    "test type invalid member",
				ops:     []jsoncodec.Operation{{Op: "test_type", Path: "/kind", Type: []any{"string", 1}}},
				opts:    allCapabilities(),
				wantErr: jsonpatch.ErrPayloadInvalid,
			},
			{
				name:    "in requires array",
				ops:     []jsoncodec.Operation{{Op: "in", Path: "/kind", Value: "admin"}},
				opts:    allCapabilities(),
				wantErr: jsonpatch.ErrPayloadInvalid,
			},
			{
				name:    "merge position must be positive",
				ops:     []jsoncodec.Operation{{Op: "merge", Path: "/items/1"}},
				opts:    allCapabilities(),
				wantErr: jsonpatch.ErrPayloadInvalid,
			},
			{
				name:    "composite requires operands",
				ops:     []jsoncodec.Operation{{Op: "and", Path: "/user", Apply: []jsoncodec.Operation{}}},
				opts:    allCapabilities(),
				wantErr: jsonpatch.ErrPayloadInvalid,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				_, err := jsonpatch.CompileOperations(tt.ops, tt.opts...)
				if tt.wantErr == nil {
					require.NoError(t, err)
					return
				}
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
			})
		}
	})

	t.Run("apply op preserves document types", func(t *testing.T) {
		t.Parallel()

		t.Run("map document", func(t *testing.T) {
			t.Parallel()

			doc := map[string]any{"name": "Ada"}
			result, err := applyOps(t, doc, op.NewReplace([]string{"name"}, "Grace"))
			require.NoError(t, err)
			assert.Equal(t, map[string]any{"name": "Grace"}, result.Doc)
			require.Len(t, result.Steps, 1)
			assert.Equal(t, "Ada", result.Steps[0].Old())
		})

		t.Run("json bytes document", func(t *testing.T) {
			t.Parallel()

			result, err := applyOps(t, []byte(`{"name":"Ada"}`), op.NewAdd([]string{"role"}, "admin"))
			require.NoError(t, err)
			assert.JSONEq(t, `{"name":"Ada","role":"admin"}`, string(result.Doc))
		})

		t.Run("json string document", func(t *testing.T) {
			t.Parallel()

			result, err := applyOps(t, jsonpatch.JSONText(`{"name":"Ada"}`), op.NewAdd([]string{"role"}, "admin"))
			require.NoError(t, err)
			assert.JSONEq(t, `{"name":"Ada","role":"admin"}`, string(result.Doc))
		})

		t.Run("primitive document preserves convertible result", func(t *testing.T) {
			t.Parallel()

			result, err := applyOps(t, 10, op.NewInc(nil, 5))
			require.NoError(t, err)
			assert.Equal(t, 15, result.Doc)
		})

		t.Run("struct document", func(t *testing.T) {
			t.Parallel()

			before := profile{Name: "Ada", Tags: []string{"go"}}
			result, err := applyOps(t, before, op.NewReplace([]string{"name"}, "Grace"))
			require.NoError(t, err)
			assert.Equal(t, profile{Name: "Grace", Tags: []string{"go"}}, result.Doc)
			require.Len(t, result.Steps, 1)
			assert.Equal(t, "Ada", result.Steps[0].Old())
		})
	})

	t.Run("apply ops preserves nil interface results", func(t *testing.T) {
		t.Parallel()

		result, err := applyOps[any](t, map[string]any{"name": "Ada"}, op.NewReplace(nil, nil))
		require.NoError(t, err)
		assert.Nil(t, result.Doc)
		require.Len(t, result.Steps, 1)
		assert.Equal(t, map[string]any{"name": "Ada"}, result.Steps[0].Old())
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
		var operations []jsoncodec.Operation

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
		result, err := applyOperations(t, doc, operations)

		// Verify no panics occurred and results are consistent
		if err == nil && result != nil {
			// Result should be valid JSON
			resultJSON, jsonErr := json.Marshal(result.Doc)
			if jsonErr != nil {
				assert.Fail(t, fmt.Sprintf("Result is not valid JSON: %v", jsonErr))
			}

			// Result should be parseable back
			var reparsed any
			if json.Unmarshal(resultJSON, &reparsed) != nil {
				assert.Fail(t, "Result cannot be reparsed from JSON")
			}
		}

		// Errors should be properly structured
		if err != nil {
			if err.Error() == "" {
				assert.Fail(t, "Error missing message")
			}
		}
	})
}
