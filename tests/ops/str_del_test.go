package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationStrDel applies a single operation to a document
func applyOperationStrDel(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.ApplyPatchOptions{Mutate: true})
	require.NoError(t, err)
	return result.Doc
}

func TestStrDelOp(t *testing.T) {
	t.Run("deletes characters from the beginning", func(t *testing.T) {
		operation := internal.Operation{
			"op":   "str_del",
			"path": "",
			"pos":  0,
			"len":  7,
		}
		result := applyOperationStrDel(t, "Hello, world!", operation)
		assert.Equal(t, "world!", result)
	})

	t.Run("deletes characters from the end", func(t *testing.T) {
		operation := internal.Operation{
			"op":   "str_del",
			"path": "",
			"pos":  5,
			"len":  8,
		}
		result := applyOperationStrDel(t, "Hello, world!", operation)
		assert.Equal(t, "Hello", result)
	})

	t.Run("deletes characters from the middle", func(t *testing.T) {
		operation := internal.Operation{
			"op":   "str_del",
			"path": "",
			"pos":  5,
			"len":  10,
		}
		result := applyOperationStrDel(t, "Hello beautiful world", operation)
		assert.Equal(t, "Hello world", result)
	})

	t.Run("can delete multiple times", func(t *testing.T) {
		operations := []internal.Operation{
			{
				"op":   "str_del",
				"path": "",
				"pos":  5,
				"len":  10,
			},
			{
				"op":   "str_del",
				"path": "",
				"pos":  5,
				"len":  1,
			},
		}
		doc := "Hello beautiful world"
		result1, err := jsonpatch.ApplyPatch(doc, []internal.Operation{operations[0]}, internal.ApplyPatchOptions{Mutate: true})
		require.NoError(t, err)
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.ApplyPatchOptions{Mutate: true})
		require.NoError(t, err)
		assert.Equal(t, "Helloworld", result2.Doc)
	})

	t.Run("root", func(t *testing.T) {
		t.Run("deletes entire string", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "",
				"pos":  0,
				"len":  5,
			}
			result := applyOperationStrDel(t, "hello", operation)
			assert.Equal(t, "", result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/msg",
				"pos":  0,
				"len":  7,
			}
			result := applyOperationStrDel(t, map[string]interface{}{"msg": "Hello, world!"}, operation)
			expected := map[string]interface{}{"msg": "world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/msg",
				"pos":  5,
				"len":  8,
			}
			result := applyOperationStrDel(t, map[string]interface{}{"msg": "Hello, world!"}, operation)
			expected := map[string]interface{}{"msg": "Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/msg",
				"pos":  5,
				"len":  10,
			}
			result := applyOperationStrDel(t, map[string]interface{}{"msg": "Hello beautiful world"}, operation)
			expected := map[string]interface{}{"msg": "Hello world"}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/0",
				"pos":  0,
				"len":  7,
			}
			result := applyOperationStrDel(t, []interface{}{"Hello, world!"}, operation)
			expected := []interface{}{"world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/0",
				"pos":  5,
				"len":  8,
			}
			result := applyOperationStrDel(t, []interface{}{"Hello, world!"}, operation)
			expected := []interface{}{"Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "str_del",
				"path": "/0",
				"pos":  5,
				"len":  10,
			}
			result := applyOperationStrDel(t, []interface{}{"Hello beautiful world"}, operation)
			expected := []interface{}{"Hello world"}
			assert.Equal(t, expected, result)
		})
	})
}
