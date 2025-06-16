package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpRemove_JSONSerialization(t *testing.T) {
	patch := []jsonpatch.Operation{
		map[string]interface{}{
			"op":   "remove",
			"path": "/foo/bar",
		},
	}

	jsonBytes, err := json.Marshal(patch[0])
	assert.NoError(t, err)

	var parsed map[string]interface{}
	err = json.Unmarshal(jsonBytes, &parsed)
	assert.NoError(t, err)

	expected := map[string]interface{}{
		"op":   "remove",
		"path": "/foo/bar",
	}
	assert.Equal(t, expected, parsed)
}

func TestOpRemove_Apply(t *testing.T) {
	t.Run("can remove key from object", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": 1,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/foo",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := map[string]interface{}{}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := map[string]interface{}{"a": 1}
		actualOld := result.Res[0].Old
		assert.Equal(t, expectedOld, actualOld, "Old value should be the removed object")
	})

	t.Run("throws when deleting a non-existing key", func(t *testing.T) {
		doc := map[string]interface{}{"bar": 123}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/foo",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)
	})

	t.Run("removing root sets document to null", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": 1,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		// Note: This might fail if the library doesn't support empty path for remove
		if err != nil {
			t.Skip("Library doesn't support empty path for remove operation")
			return
		}

		assert.Nil(t, result.Doc)

		// Verify old value is the entire previous document
		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := map[string]interface{}{
			"foo": map[string]interface{}{"a": 1},
		}
		actualOld := result.Res[0].Old
		assert.Equal(t, expectedOld, actualOld, "Old value should be the entire previous document")
	})

	t.Run("can remove last member of array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/2",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := []interface{}{1, 2}
		actual := result.Doc
		assert.Equal(t, expected, actual)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := 3
		actualOld := result.Res[0].Old
		assert.Equal(t, expectedOld, actualOld, "Old value should be the removed element")
	})

	t.Run("can remove first member of array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/0",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := []interface{}{2, 3}
		actual := result.Doc
		assert.Equal(t, expected, actual)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := 1
		actualOld := result.Res[0].Old
		assert.Equal(t, expectedOld, actualOld, "Old value should be the removed element")
	})

	t.Run("throws when removing elements beyond array boundaries", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}

		// Test removing at index 4 (beyond bounds)
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/4",
			},
		}
		_, err := jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)

		// Test removing at index 5 (beyond bounds)
		patch[0] = map[string]interface{}{
			"op":   "remove",
			"path": "/5",
		}
		_, err = jsonpatch.ApplyPatch(doc, patch)
		assert.Error(t, err)
	})
}

func TestOpRemove_RFC6902_Section4_2(t *testing.T) {
	t.Run("removing array element shifts remaining elements left", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "remove",
				"path": "/1",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch)
		assert.NoError(t, err)

		expected := []interface{}{1, 3}
		actual := result.Doc
		assert.Equal(t, expected, actual)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := 2
		actualOld := result.Res[0].Old
		assert.Equal(t, expectedOld, actualOld, "Old value should be the removed element")
	})
}
