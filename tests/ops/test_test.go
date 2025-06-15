package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperation applies a single operation to a document
func applyOperation(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.ApplyPatchOptions{Mutate: true})
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithError applies an operation expecting it to fail
func applyOperationWithError(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.ApplyPatchOptions{Mutate: true})
	require.Error(t, err)
}

func TestTestOp(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("should test against root on json document of type object and return true", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				"op":    "test",
				"path":  "",
				"value": map[string]interface{}{"hello": "world"},
			}
			result := applyOperation(t, obj, op)
			assert.Equal(t, obj, result)
		})

		t.Run("should test against root on json document of type object and return false", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				"op":    "test",
				"path":  "",
				"value": 1,
			}
			applyOperationWithError(t, obj, op)
		})

		t.Run("should test against root on json document of type array and return false", func(t *testing.T) {
			obj := []interface{}{
				map[string]interface{}{
					"hello": "world",
				},
			}
			op := internal.Operation{
				"op":    "test",
				"path":  "",
				"value": 1,
			}
			applyOperationWithError(t, obj, op)
		})

		t.Run("should throw against root", func(t *testing.T) {
			op := internal.Operation{
				"op":    "test",
				"path":  "",
				"value": 2,
				"not":   false,
			}
			applyOperationWithError(t, 1, op)
		})

		t.Run("should throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				"op":    "test",
				"path":  "/foo",
				"value": 2,
				"not":   false,
			}
			applyOperationWithError(t, obj, op)
		})

		t.Run("should not throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				"op":    "test",
				"path":  "/foo",
				"value": 1,
				"not":   false,
			}
			applyOperation(t, obj, op)
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("should test against root", func(t *testing.T) {
			op := internal.Operation{
				"op":    "test",
				"path":  "",
				"value": 2,
				"not":   true,
			}
			applyOperation(t, 1, op)
		})

		t.Run("should not throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				"op":    "test",
				"path":  "/foo",
				"value": 2,
				"not":   true,
			}
			applyOperation(t, obj, op)
		})

		t.Run("should throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				"op":    "test",
				"path":  "/foo",
				"value": 1,
				"not":   true,
			}
			applyOperationWithError(t, obj, op)
		})
	})
}
