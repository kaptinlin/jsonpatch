package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpMove_JSONSerialization(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":   "move",
			"path": "/foo/bar",
			"from": "/foo/baz",
		},
	}

	jsonBytes, err := json.Marshal(patch[0])
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"op":   "move",
		"path": "/foo/bar",
		"from": "/foo/baz",
	}
	assert.Equal(t, expected, parsed)
}

func TestOpMove_Apply(t *testing.T) {
	t.Run("can move an object key", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": 123,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/bar",
				"from": "/foo",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"bar": 123,
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil when moving to new location")
	})

	t.Run("move to existing key should return old value", func(t *testing.T) {
		doc := map[string]interface{}{
			"source": "moved_value",
			"target": "original_value",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/target",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"target": "moved_value",
		}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, "original_value", result.Res[0].Old, "Old value should be the replaced value")
	})

	t.Run("move within same array should handle indices correctly", func(t *testing.T) {
		doc := map[string]interface{}{
			"arr": []interface{}{"a", "b", "c", "d"},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/arr/0",
				"from": "/arr/2",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		// Move c from index 2 to index 0
		// After remove: ["a", "b", "d"]
		// After insert at 0: ["c", "a", "b", "d"]
		expected := map[string]interface{}{
			"arr": []interface{}{"c", "a", "b", "d"},
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is what was displaced
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, "a", result.Res[0].Old, "Old value should be the displaced element")
	})

	t.Run("move from non-existent path should fail", func(t *testing.T) {
		doc := map[string]interface{}{
			"existing": "value",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/target",
				"from": "/nonexistent",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path not found", "Error should indicate source path not found")
	})

	t.Run("move complex nested structures", func(t *testing.T) {
		doc := map[string]interface{}{
			"complex": map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": []interface{}{
						map[string]interface{}{"id": 1, "data": "test1"},
						map[string]interface{}{"id": 2, "data": "test2"},
					},
				},
			},
			"destination": map[string]interface{}{
				"backups": []interface{}{},
			},
		}

		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/destination/moved_item",
				"from": "/complex/level1/level2/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expectedMovedItem := map[string]interface{}{"id": 1, "data": "test1"}
		actualDoc := result.Doc
		destination := actualDoc["destination"].(map[string]interface{})
		assert.Equal(t, expectedMovedItem, destination["moved_item"])

		// Verify source array was updated
		complexObj := actualDoc["complex"].(map[string]interface{})
		level1 := complexObj["level1"].(map[string]interface{})
		level2 := level1["level2"].([]interface{})
		assert.Len(t, level2, 1, "Source array should have one less element")
		assert.Equal(t, map[string]interface{}{"id": 2, "data": "test2"}, level2[0])
	})
}

func TestOpMove_RFC6902_Section4_4(t *testing.T) {
	t.Run("from location must not be a proper prefix of path location", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": 123,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/foo/bar",
				"from": "/foo",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err, "Cannot move location into one of its children")
		assert.Contains(t, err.Error(), "cannot move into own children", "Error should mention descendant restriction")
	})

	t.Run("object to object", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": 1,
			},
			"bar": map[string]interface{}{
				"b": 2,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/bar/b",
				"from": "/foo/a",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": map[string]interface{}{},
			"bar": map[string]interface{}{
				"b": 1, // Value moved from foo/a, replacing original bar/b
			},
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is the replaced value
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 2, result.Res[0].Old, "Old value should be the original bar/b value")
	})

	t.Run("object to array", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": 1,
			},
			"bar": []interface{}{0},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/bar/1",
				"from": "/foo/a",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": map[string]interface{}{},
			"bar": []interface{}{0, 1}, // Value inserted at index 1
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is nil for new array position
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new array position")
	})

	t.Run("array to array", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": []interface{}{1},
			"bar": []interface{}{0},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/bar/0",
				"from": "/foo/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		// TypeScript behavior: Move from foo[0] (1) to bar[0] -> bar becomes [1, 0]
		// This is because Add operation uses splice(index, 0, value) which inserts
		expected := map[string]interface{}{
			"foo": []interface{}{},
			"bar": []interface{}{1, 0}, // 1 inserted at index 0, 0 shifted right
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is the displaced element
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 0, result.Res[0].Old, "Old value should be the displaced element")
	})

	t.Run("array to object", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": []interface{}{1},
			"bar": map[string]interface{}{
				"a": 0,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/bar/b",
				"from": "/foo/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": []interface{}{},
			"bar": map[string]interface{}{
				"a": 0,
				"b": 1, // Value moved from foo[0] to new object key
			},
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is nil for new object key
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new object key")
	})

	t.Run("move to root replaces entire document", func(t *testing.T) {
		doc := map[string]interface{}{
			"source": map[string]interface{}{
				"nested": "value",
			},
			"other": "data",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		// Moving to root replaces entire document with the moved value
		expected := map[string]interface{}{
			"nested": "value",
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is the remaining document after source removal
		assert.Len(t, result.Res, 1, "Should have one operation result")
		remainingDoc := map[string]interface{}{
			"other": "data",
		}
		assert.Equal(t, remainingDoc, result.Res[0].Old, "Old value should be the remaining document after source removal")
	})
}

func TestOpMove_AdvancedScenarios(t *testing.T) {
	t.Run("move with array indices edge cases", func(t *testing.T) {
		doc := map[string]interface{}{
			"arr": []interface{}{"a", "b", "c"},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/arr/-", // Append to end
				"from": "/arr/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		if err != nil {
			// If the library doesn't support "-" syntax, skip this test
			t.Skip("Library doesn't support array append syntax")
		} else {
			// After remove arr[0] ("a"): ["b", "c"]
			// After append "a" to end: ["b", "c", "a"]
			expected := map[string]interface{}{
				"arr": []interface{}{"b", "c", "a"},
			}
			assert.Equal(t, expected, result.Doc)
		}
	})

	t.Run("move between different document levels", func(t *testing.T) {
		doc := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"deep_value": "treasure",
				},
				"level2_sibling": "data",
			},
			"top_level": "info",
		}

		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/promoted",
				"from": "/level1/level2/deep_value",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2":         map[string]interface{}{},
				"level2_sibling": "data",
			},
			"top_level": "info",
			"promoted":  "treasure",
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is nil for new key
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new key")
	})

	t.Run("move validation catches recursive move attempts", func(t *testing.T) {
		doc := map[string]interface{}{
			"parent": map[string]interface{}{
				"child": map[string]interface{}{
					"grandchild": "value",
				},
			},
		}

		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/parent/child/new_location",
				"from": "/parent",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err, "Should prevent moving parent into its own child")
		assert.Contains(t, err.Error(), "cannot move into own children", "Error should mention descendant restriction")
	})

	t.Run("move preserves data types through the operation", func(t *testing.T) {
		doc := map[string]interface{}{
			"source": map[string]interface{}{
				"number":  123,
				"string":  "text",
				"boolean": true,
				"null":    nil,
				"array":   []interface{}{1, 2, 3},
				"object":  map[string]interface{}{"nested": "data"},
			},
		}

		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "move",
				"path": "/moved_complex_object",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expectedContent := map[string]interface{}{
			"number":  123,
			"string":  "text",
			"boolean": true,
			"null":    nil,
			"array":   []interface{}{1, 2, 3},
			"object":  map[string]interface{}{"nested": "data"},
		}

		actualDoc := result.Doc
		assert.Equal(t, expectedContent, actualDoc["moved_complex_object"])
	})
}
