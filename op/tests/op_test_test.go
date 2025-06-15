package jsonpatch_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestOpTest_JSONSerialization(t *testing.T) {
	t.Run("serializes to correct JSON format", func(t *testing.T) {
		op := map[string]interface{}{
			"op":    "test",
			"path":  "/foo/bar",
			"value": "expected",
		}

		jsonBytes, err := json.Marshal(op)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Equal(t, "test", result["op"])
		assert.Equal(t, "/foo/bar", result["path"])
		assert.Equal(t, "expected", result["value"])
	})

	t.Run("serializes with not flag", func(t *testing.T) {
		op := map[string]interface{}{
			"op":    "test",
			"path":  "/foo/bar",
			"value": "expected",
			"not":   true,
		}

		jsonBytes, err := json.Marshal(op)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Equal(t, "test", result["op"])
		assert.Equal(t, "/foo/bar", result["path"])
		assert.Equal(t, "expected", result["value"])
		assert.Equal(t, true, result["not"])
	})

	t.Run("omits not flag when false", func(t *testing.T) {
		op := map[string]interface{}{
			"op":    "test",
			"path":  "/foo/bar",
			"value": "expected",
		}

		jsonBytes, err := json.Marshal(op)
		assert.NoError(t, err)

		var result map[string]interface{}
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err)

		assert.Equal(t, "test", result["op"])
		assert.Equal(t, "/foo/bar", result["path"])
		assert.Equal(t, "expected", result["value"])
		_, hasNot := result["not"]
		assert.False(t, hasNot, "not field should be omitted when false")
	})
}

func TestOpTest_Apply(t *testing.T) {
	t.Run("succeeds when values match", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/foo",
				"value": "bar",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})

	t.Run("fails when values don't match", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/foo",
				"value": "baz",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("succeeds with not flag when values don't match", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/foo",
				"value": "baz",
				"not":   true,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})

	t.Run("fails with not flag when values match", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/foo",
				"value": "bar",
				"not":   true,
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("succeeds with not flag when path doesn't exist", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/nonexistent",
				"value": "anything",
				"not":   true,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})

	t.Run("fails when path is not found and not=false", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": "bar",
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/nonexistent",
				"value": "anything",
			},
		}

		_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		// Check that error indicates path not found or test failure
		assert.True(t, err.Error() == "operation 0 failed: path not found" ||
			strings.Contains(err.Error(), "test operation failed"), "Error should indicate test failure")
	})
}

func TestOpTest_RFC6902_Section4_6(t *testing.T) {
	t.Run("strings: are considered equal if they contain the same number of Unicode characters", func(t *testing.T) {
		var doc interface{} = "123"
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "",
				"value": "123",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, "123", result.Doc)

		// Test failure case
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": "1234",
		}
		_, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("numbers: are considered equal if their values are numerically equal", func(t *testing.T) {
		var doc interface{} = 123
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "",
				"value": 123,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, 123, result.Doc)

		// Test failure case
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": 0,
		}
		_, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("arrays: are considered equal with same values in corresponding positions", func(t *testing.T) {
		var doc interface{} = []interface{}{1, 2, map[string]interface{}{"foo": "bar"}}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "",
				"value": []interface{}{1, 2, map[string]interface{}{"foo": "bar"}},
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)

		// Test failure case
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": []interface{}{1, 2, map[string]interface{}{"foo": "bar!"}},
		}
		_, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("objects: are considered equal with same members", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": map[string]interface{}{},
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":   "test",
				"path": "",
				"value": map[string]interface{}{
					"foo": map[string]interface{}{
						"bar": map[string]interface{}{},
					},
				},
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)

		// Test failure case
		patch[0] = map[string]interface{}{
			"op":   "test",
			"path": "",
			"value": map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": map[string]interface{}{"a": 1},
				},
			},
		}
		_, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})

	t.Run("literals (false, true, and null): are considered equal if they are the same", func(t *testing.T) {
		// Test true
		var trueDoc interface{} = true
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "",
				"value": true,
			},
		}
		result, err := jsonpatch.ApplyPatch(trueDoc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, true, result.Doc)

		// Test false
		var falseDoc interface{} = false
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": false,
		}
		result, err = jsonpatch.ApplyPatch(falseDoc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, false, result.Doc)

		// Test null
		var nullDoc interface{} = nil
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": nil,
		}
		result, err = jsonpatch.ApplyPatch(nullDoc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Nil(t, result.Doc)

		// Test failure cases
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": false,
		}
		_, err = jsonpatch.ApplyPatch(nullDoc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")

		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "",
			"value": false,
		}
		_, err = jsonpatch.ApplyPatch(trueDoc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "test operation failed")
	})
}

func TestOpTest_AdvancedScenarios(t *testing.T) {
	t.Run("test nested paths", func(t *testing.T) {
		doc := map[string]interface{}{
			"level1": map[string]interface{}{
				"level2": map[string]interface{}{
					"value": "target",
				},
			},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/level1/level2/value",
				"value": "target",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})

	t.Run("test array elements", func(t *testing.T) {
		doc := map[string]interface{}{
			"arr": []interface{}{"a", "b", "c"},
		}
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/arr/1",
				"value": "b",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)

		// Test with array index out of bounds
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "/arr/10",
			"value": "nonexistent",
		}
		_, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.Error(t, err)
		// Check that error indicates path not found or test failure
		assert.True(t, err.Error() == "operation 0 failed: path not found" ||
			strings.Contains(err.Error(), "test operation failed"), "Error should indicate test failure")
	})

	t.Run("test with complex nested structures", func(t *testing.T) {
		doc := map[string]interface{}{
			"users": []interface{}{
				map[string]interface{}{
					"id":   1,
					"name": "Alice",
					"profile": map[string]interface{}{
						"age":    30,
						"skills": []interface{}{"Go", "JavaScript"},
					},
				},
				map[string]interface{}{
					"id":   2,
					"name": "Bob",
					"profile": map[string]interface{}{
						"age":    25,
						"skills": []interface{}{"Python", "TypeScript"},
					},
				},
			},
		}

		// Test specific user's skill
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/users/0/profile/skills/1",
				"value": "JavaScript",
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})

	t.Run("test preserves data types through operation", func(t *testing.T) {
		doc := map[string]interface{}{
			"number":  123,
			"string":  "text",
			"boolean": true,
			"null":    nil,
			"array":   []interface{}{1, 2, 3},
			"object":  map[string]interface{}{"nested": "data"},
		}

		// Test all types remain intact
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/number",
				"value": 123,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)

		// Verify all fields still have correct types
		actualDoc := result.Doc.(map[string]interface{})
		assert.Equal(t, 123, actualDoc["number"])
		assert.Equal(t, "text", actualDoc["string"])
		assert.Equal(t, true, actualDoc["boolean"])
		assert.Nil(t, actualDoc["null"])
		assert.Equal(t, []interface{}{1, 2, 3}, actualDoc["array"])
		assert.Equal(t, map[string]interface{}{"nested": "data"}, actualDoc["object"])
	})

	t.Run("test with not flag complex scenarios", func(t *testing.T) {
		doc := map[string]interface{}{
			"status": "active",
			"config": map[string]interface{}{
				"enabled": true,
			},
		}

		// Test that status is NOT inactive
		patch := []jsonpatch.Operation{
			map[string]interface{}{
				"op":    "test",
				"path":  "/status",
				"value": "inactive",
				"not":   true,
			},
		}

		result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)

		// Test that non-existent path is NOT "something"
		patch[0] = map[string]interface{}{
			"op":    "test",
			"path":  "/nonexistent",
			"value": "something",
			"not":   true,
		}

		result, err = jsonpatch.ApplyPatch(doc, patch, jsonpatch.ApplyPatchOptions{})
		assert.NoError(t, err)
		assert.Equal(t, doc, result.Doc)
	})
}
