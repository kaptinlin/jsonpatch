package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStrDelOp(t *testing.T) {
	t.Run("deletes characters from the beginning", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  0,
			Len:  7,
		}
		result := testutils.ApplyInternalOp(t, "Hello, world!", operation)
		assert.Equal(t, "world!", result)
	})

	t.Run("deletes characters from the end", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  8,
		}
		result := testutils.ApplyInternalOp(t, "Hello, world!", operation)
		assert.Equal(t, "Hello", result)
	})

	t.Run("deletes characters from the middle", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  10,
		}
		result := testutils.ApplyInternalOp(t, "Hello beautiful world", operation)
		assert.Equal(t, "Hello world", result)
	})

	t.Run("can delete multiple times", func(t *testing.T) {
		operations := []internal.Operation{
			{
				Op:   "str_del",
				Path: "",
				Pos:  5,
				Len:  10,
			},
			{
				Op:   "str_del",
				Path: "",
				Pos:  5,
				Len:  1,
			},
		}
		doc := "Hello beautiful world"
		result1, err := jsonpatch.ApplyPatch(doc, []internal.Operation{operations[0]}, internal.WithMutate(true))
		require.NoError(t, err)
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.WithMutate(true))
		require.NoError(t, err)
		assert.Equal(t, "Helloworld", result2.Doc)
	})

	t.Run("root", func(t *testing.T) {
		t.Run("deletes entire string", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "",
				Pos:  0,
				Len:  5,
			}
			result := testutils.ApplyInternalOp(t, "hello", operation)
			assert.Equal(t, "", result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  0,
				Len:  7,
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello, world!"}, operation)
			expected := map[string]interface{}{"msg": "world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  5,
				Len:  8,
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello, world!"}, operation)
			expected := map[string]interface{}{"msg": "Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  5,
				Len:  10,
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello beautiful world"}, operation)
			expected := map[string]interface{}{"msg": "Hello world"}
			assert.Equal(t, expected, result)
		})

		t.Run("negative position counts from end", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  -1,
				Len:  1,
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello!"}, operation)
			expected := map[string]interface{}{"msg": "Hello"}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  0,
				Len:  7,
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"Hello, world!"}, operation)
			expected := []interface{}{"world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  5,
				Len:  8,
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"Hello, world!"}, operation)
			expected := []interface{}{"Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  5,
				Len:  10,
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"Hello beautiful world"}, operation)
			expected := []interface{}{"Hello world"}
			assert.Equal(t, expected, result)
		})
	})
}
