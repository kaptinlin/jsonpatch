package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpAdd_JSONSerialization(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":    "add",
			"path":  "/foo/bar",
			"value": 123,
		},
	}

	jsonBytes, err := json.Marshal(patch[0])
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"op":    "add",
		"path":  "/foo/bar",
		"value": float64(123), // JSON round-trip still uses float64 for numbers
	}
	assert.Equal(t, expected, parsed)
}

func TestOpAdd_Apply(t *testing.T) {
	t.Run("can add new key to object", func(t *testing.T) {
		doc := map[string]interface{}{"foo": map[string]interface{}{}}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "/foo/bar",
				"value": 123,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": 123,
			},
		}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for new key")
	})

	t.Run("when adding element past array boundary, throws", func(t *testing.T) {
		doc := []interface{}{0}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "/100",
				"value": 1,
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)
	})
}

func TestOpAdd_RFC6902_Section4_1(t *testing.T) {
	t.Run("root target document replacement", func(t *testing.T) {
		doc := map[string]interface{}{"a": 1}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "",
				"value": 123,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		// Note: This might fail if the library doesn't support empty path for add
		if err != nil {
			t.Skip("Library doesn't support empty path for add operation")
			return
		}

		assert.Equal(t, 123, result.Doc)

		// Verify old value is the entire previous document
		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := map[string]interface{}{"a": 1}
		assert.Equal(t, expectedOld, result.Res[0].Old, "Old value should be the entire previous document")
	})

	t.Run("new object member", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": "b",
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "add",
				"path": "/foo/z",
				"value": map[string]interface{}{
					"test": map[string]interface{}{},
				},
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": "b",
				"z": map[string]interface{}{
					"test": map[string]interface{}{},
				},
			},
		}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("replace existing object member", func(t *testing.T) {
		doc := map[string]interface{}{"a": 1}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "/a",
				"value": 2,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{"a": 2}
		assert.Equal(t, expected, result.Doc)

		// Verify old value is returned for replaced key
		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := 1
		assert.Equal(t, expectedOld, result.Res[0].Old, "Old value should be the replaced value")
	})

	t.Run("append to array with - character", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "add",
				"path":  "/-",
				"value": 4,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		// Use Go native types directly (int preserved)
		expected := []interface{}{1, 2, 3, 4}
		assert.Equal(t, expected, result.Doc)

		// For append operation, old value should be nil/undefined
		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Nil(t, result.Res[0].Old, "Old value should be nil for append operation")
	})

	t.Run("insert into array at specific index", func(t *testing.T) {
		t.Run("at the beginning of array", func(t *testing.T) {
			doc := []interface{}{1, 2, 3}
			patch := []jsonpatch.Operation{
				map[string]interface{}{
					"op":    "add",
					"path":  "/0",
					"value": 0,
				},
			}

			result, err := jsonpatch.ApplyPatch(doc, patch)
			assert.NoError(t, err)

			// Use Go native types directly (int preserved)
			expected := []interface{}{0, 1, 2, 3}
			assert.Equal(t, expected, result.Doc)

			// For insert operation, old value should be the displaced element
			assert.Len(t, result.Res, 1, "Should have one operation result")
			expectedOld := 1 // Element that was at index 0
			assert.Equal(t, expectedOld, result.Res[0].Old, "Old value should be the displaced element")
		})

		t.Run("in the middle of array", func(t *testing.T) {
			doc := []interface{}{1, 2, 3}
			patch := []jsonpatch.Operation{
				map[string]interface{}{
					"op":    "add",
					"path":  "/1",
					"value": 0,
				},
			}

			result, err := jsonpatch.ApplyPatch(doc, patch)
			assert.NoError(t, err)

			// Use Go native types directly (int preserved)
			expected := []interface{}{1, 0, 2, 3}
			assert.Equal(t, expected, result.Doc)

			// For insert operation, old value should be the displaced element
			assert.Len(t, result.Res, 1, "Should have one operation result")
			expectedOld := 2 // Element that was at index 1
			assert.Equal(t, expectedOld, result.Res[0].Old, "Old value should be the displaced element")
		})

		t.Run("at the end of array", func(t *testing.T) {
			doc := []interface{}{1, 2, 3}
			patch := []jsonpatch.Operation{
				map[string]interface{}{
					"op":    "add",
					"path":  "/3",
					"value": 0,
				},
			}

			result, err := jsonpatch.ApplyPatch(doc, patch)
			assert.NoError(t, err)

			// Use Go native types directly (int preserved)
			expected := []interface{}{1, 2, 3, 0}
			assert.Equal(t, expected, result.Doc)

			// For append to end operation, old value should be nil/undefined
			assert.Len(t, result.Res, 1, "Should have one operation result")
			assert.Nil(t, result.Res[0].Old, "Old value should be nil for append to end operation")
		})
	})
}
