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
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithError applies an operation expecting it to fail
func applyOperationWithError(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.Error(t, err)
}

func TestTestOp(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("should test against root on json document of type object and return true", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: map[string]interface{}{"hello": "world"},
			}
			result := applyOperation(t, obj, op)
			assert.Equal(t, obj, result)
		})

		t.Run("should test against root on json document of type object and return false", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 1,
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
				Op:    "test",
				Path:  "",
				Value: 1,
			}
			applyOperationWithError(t, obj, op)
		})

		t.Run("should throw against root", func(t *testing.T) {
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 2,
				Not:   false,
			}
			applyOperationWithError(t, 1, op)
		})

		t.Run("should throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 2,
				Not:   false,
			}
			applyOperationWithError(t, obj, op)
		})

		t.Run("should not throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 1,
				Not:   false,
			}
			applyOperation(t, obj, op)
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("should test against root", func(t *testing.T) {
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 2,
				Not:   true,
			}
			applyOperation(t, 1, op)
		})

		t.Run("should not throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 2,
				Not:   true,
			}
			applyOperation(t, obj, op)
		})

		t.Run("should throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 1,
				Not:   true,
			}
			applyOperationWithError(t, obj, op)
		})
	})
}
