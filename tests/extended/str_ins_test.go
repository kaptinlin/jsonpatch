package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationStrIns applies a single operation to a document
func applyOperationStrIns(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

func TestStrInsOp(t *testing.T) {
	t.Run("inserts a string at the beginning", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  0,
			Str:  "Hello, ",
		}
		result := applyOperationStrIns(t, "world!", operation)
		assert.Equal(t, "Hello, world!", result)
	})

	t.Run("inserts a string at the end", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  5,
			Str:  ", world",
		}
		result := applyOperationStrIns(t, "Hello", operation)
		assert.Equal(t, "Hello, world", result)
	})

	t.Run("inserts a string in the middle", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  5,
			Str:  " beautiful",
		}
		result := applyOperationStrIns(t, "Hello world", operation)
		assert.Equal(t, "Hello beautiful world", result)
	})

	t.Run("can insert multiple times", func(t *testing.T) {
		operations := []internal.Operation{
			{
				Op:   "str_ins",
				Path: "",
				Pos:  5,
				Str:  " beautiful",
			},
			{
				Op:   "str_ins",
				Path: "",
				Pos:  21,
				Str:  " bright",
			},
		}
		doc := "Hello world"
		result1, err := jsonpatch.ApplyPatch(doc, []internal.Operation{operations[0]}, internal.WithMutate(true))
		require.NoError(t, err)
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.WithMutate(true))
		require.NoError(t, err)
		assert.Equal(t, "Hello beautiful world bright", result2.Doc)
	})

	t.Run("root", func(t *testing.T) {
		t.Run("inserts into empty string", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "",
				Pos:  0,
				Str:  "hello",
			}
			result := applyOperationStrIns(t, "", operation)
			assert.Equal(t, "hello", result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("inserts a string at the beginning", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  0,
				Str:  "Hello, ",
			}
			result := applyOperationStrIns(t, map[string]interface{}{"msg": "world!"}, operation)
			expected := map[string]interface{}{"msg": "Hello, world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("inserts a string at the end", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  5,
				Str:  ", world",
			}
			result := applyOperationStrIns(t, map[string]interface{}{"msg": "Hello"}, operation)
			expected := map[string]interface{}{"msg": "Hello, world"}
			assert.Equal(t, expected, result)
		})

		t.Run("inserts a string in the middle", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  5,
				Str:  " beautiful",
			}
			result := applyOperationStrIns(t, map[string]interface{}{"msg": "Hello world"}, operation)
			expected := map[string]interface{}{"msg": "Hello beautiful world"}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("inserts a string at the beginning", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  0,
				Str:  "Hello, ",
			}
			result := applyOperationStrIns(t, []interface{}{"world!"}, operation)
			expected := []interface{}{"Hello, world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("inserts a string at the end", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  5,
				Str:  ", world",
			}
			result := applyOperationStrIns(t, []interface{}{"Hello"}, operation)
			expected := []interface{}{"Hello, world"}
			assert.Equal(t, expected, result)
		})

		t.Run("inserts a string in the middle", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  5,
				Str:  " beautiful",
			}
			result := applyOperationStrIns(t, []interface{}{"Hello world"}, operation)
			expected := []interface{}{"Hello beautiful world"}
			assert.Equal(t, expected, result)
		})
	})
}
