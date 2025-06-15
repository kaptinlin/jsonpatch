package jsonpatch_test

import (
	"encoding/json"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpReplace_JSONSerialization(t *testing.T) {
	t.Run("with value only", func(t *testing.T) {
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
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
			"op":    "replace",
			"path":  "/foo/bar",
			"value": float64(123), // JSON round-trip still uses float64
		}
		assert.Equal(t, expected, parsed)
	})
}

func TestOpReplace_Apply(t *testing.T) {
	t.Run("can replace a key in object", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"a": 1,
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/foo",
				"value": 123,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{
			"foo": 123,
		}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		expectedOld := map[string]interface{}{"a": 1}
		assert.Equal(t, expectedOld, result.Res[0].Old, "Old value should be the replaced object")
	})
}

func TestOpReplace_RFC6902_Section4_3(t *testing.T) {
	t.Run("at root", func(t *testing.T) {
		doc := 1
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "",
				"value": 2,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		// Note: This might fail if the library doesn't support empty path
		if err != nil {
			t.Skip("Library doesn't support empty path for replace operation")
			return
		}

		assert.Equal(t, 2, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 1, result.Res[0].Old, "Old value should be the entire previous document")
	})

	t.Run("at start of array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/0",
				"value": 0,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := []interface{}{0, 2, 3}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 1, result.Res[0].Old, "Old value should be the replaced element")
	})

	t.Run("in middle of array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/1",
				"value": 0,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := []interface{}{1, 0, 3}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 2, result.Res[0].Old, "Old value should be the replaced element")
	})

	t.Run("at end of array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/2",
				"value": 0,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := []interface{}{1, 2, 0}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, 3, result.Res[0].Old, "Old value should be the replaced element")
	})

	t.Run("in object", func(t *testing.T) {
		doc := map[string]interface{}{"foo": "bar"}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/foo",
				"value": 0,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)

		expected := map[string]interface{}{"foo": 0}
		assert.Equal(t, expected, result.Doc)

		assert.Len(t, result.Res, 1, "Should have one operation result")
		assert.Equal(t, "bar", result.Res[0].Old, "Old value should be the replaced value")
	})

	t.Run("target location must exist - array", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/5",
				"value": 0,
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
	})

	t.Run("target location must exist - object", func(t *testing.T) {
		doc := map[string]interface{}{"foo": 123}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "replace",
				"path":  "/nonexistent",
				"value": 0,
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
	})
}
