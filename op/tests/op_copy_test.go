package jsonpatch_test

import (
	"encoding/json"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpCopy_JSONSerialization(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":   "copy",
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
		"op":   "copy",
		"path": "/foo/bar",
		"from": "/foo/baz",
	}
	assert.Equal(t, expected, parsed)
}

func TestOpCopy_Apply(t *testing.T) {
	t.Run("can add new key to object", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo":    map[string]interface{}{},
			"source": 123,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/foo/bar",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": 123, // ApplyPatch preserves Go native types
			},
			"source": 123,
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is nil when adding to new location
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new key")
	})

	t.Run("can replace existing key", func(t *testing.T) {
		doc := map[string]interface{}{
			"target": "old_value",
			"source": "new_value",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/target",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"target": "new_value", // Value copied from source
			"source": "new_value", // Source remains unchanged
		}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, "old_value", result.Res[0].Old, "Old value should be the replaced value")
	})

	t.Run("when adding element past array boundary, throws", func(t *testing.T) {
		doc := map[string]interface{}{
			"a":      []interface{}{0},
			"source": 123,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/a/100",
				"from": "/source",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
	})

	t.Run("recursive copy protection", func(t *testing.T) {
		doc := map[string]interface{}{
			"source": map[string]interface{}{
				"target": 0,
				"foo":    1,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/source/target",
				"from": "/source",
			},
		}

		// Note: TypeScript allows recursive copying with proper deep cloning
		// Go implementation may provide protection or handle it properly
		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})

		if err != nil {
			// Go provides safety protection
			assert.Contains(t, err.Error(), "cannot copy parent into child", "Error should mention recursive copy issue")
		} else {
			// Successful recursive copy with proper cloning
			expectedDoc := map[string]interface{}{
				"source": map[string]interface{}{
					"target": map[string]interface{}{
						"target": 0,
						"foo":    1,
					},
					"foo": 1,
				},
			}
			assert.Equal(t, expectedDoc, result.Doc)
			assert.Equal(t, 0, result.Res[0].Old, "Old value should be original target value")
		}
	})

	t.Run("deep copy verification", func(t *testing.T) {
		// Test that the copy is deep, not shallow
		doc := map[string]interface{}{
			"source": map[string]interface{}{
				"nested": map[string]interface{}{
					"value": 123,
				},
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/target",
				"from": "/source",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		// Verify the copy was made
		assert.Equal(t, doc["source"], result.Doc.(map[string]interface{})["target"])

		// Verify it's a deep copy by modifying the original and checking target is unchanged
		originalSource := doc["source"].(map[string]interface{})
		originalSource["nested"].(map[string]interface{})["value"] = 999

		// Target should still have the original value (proving deep copy)
		targetSource := result.Doc.(map[string]interface{})["target"].(map[string]interface{})
		assert.Equal(t, 123, targetSource["nested"].(map[string]interface{})["value"],
			"Target should maintain original value, proving deep copy")
	})

	t.Run("copy from array element", func(t *testing.T) {
		doc := map[string]interface{}{
			"arr": []interface{}{10, 20, 30},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/value",
				"from": "/arr/1",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"arr":   []interface{}{10, 20, 30},
			"value": 20, // Copied from arr[1]
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is nil for new key
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new key")
	})

	t.Run("copy to array element", func(t *testing.T) {
		doc := map[string]interface{}{
			"arr":   []interface{}{10, 20, 30},
			"value": 999,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/arr/1",
				"from": "/value",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"arr":   []interface{}{10, 999, 20, 30}, // 999 inserted at index 1
			"value": 999,
		}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is the displaced element
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 20, result.Res[0].Old, "Old value should be the displaced element")
	})
}

func TestOpCopy_RFC6902_Section4_5(t *testing.T) {
	t.Run("copies value at specified location to target location", func(t *testing.T) {
		doc := map[string]interface{}{
			"a": 1,
			"b": 2,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/a",
				"from": "/b",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		// ApplyPatch preserves Go native types
		expected := map[string]interface{}{
			"a": 2, // Value copied from b
			"b": 2, // Source unchanged
		}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 1, result.Res[0].Old, "Old value should be the original value at target location")
	})

	t.Run("from location must exist", func(t *testing.T) {
		doc := map[string]interface{}{
			"a": 1,
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/a",
				"from": "/nonexistent",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "path not found", "Error should indicate source path not found")
	})

	t.Run("copy from root to nested path", func(t *testing.T) {
		doc := "root_value"
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/nested/key",
				"from": "",
			},
		}

		// This should fail because we can't create nested path on a string root
		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err, "Should fail when trying to create nested path on non-object root")
	})

	t.Run("copy complex nested structures", func(t *testing.T) {
		doc := map[string]interface{}{
			"complex": map[string]interface{}{
				"level1": map[string]interface{}{
					"level2": []interface{}{
						map[string]interface{}{"id": 1, "data": "test1"},
						map[string]interface{}{"id": 2, "data": "test2"},
					},
				},
			},
		}

		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "copy",
				"path": "/backup",
				"from": "/complex/level1/level2/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expectedBackup := map[string]interface{}{"id": 1, "data": "test1"}
		actualDoc := result.Doc.(map[string]interface{})
		assert.Equal(t, expectedBackup, actualDoc["backup"])

		// Verify it's a deep copy
		originalItem := doc["complex"].(map[string]interface{})["level1"].(map[string]interface{})["level2"].([]interface{})[0].(map[string]interface{})
		originalItem["data"] = "modified"

		backupItem := actualDoc["backup"].(map[string]interface{})
		assert.Equal(t, "test1", backupItem["data"], "Backup should not be affected by original modification")
	})
}
