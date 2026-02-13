package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestStrDelOp(t *testing.T) {
	t.Parallel()
	t.Run("deletes characters from the beginning", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  0,
			Len:  7,
		}
		result := testutils.ApplyInternalOp(t, "Hello, world!", operation)
		assert.Equal(t, "world!", result, "result")
	})

	t.Run("deletes characters from the end", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  8,
		}
		result := testutils.ApplyInternalOp(t, "Hello, world!", operation)
		assert.Equal(t, "Hello", result, "result")
	})

	t.Run("deletes characters from the middle", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  10,
		}
		result := testutils.ApplyInternalOp(t, "Hello beautiful world", operation)
		if result != "Hello world" {
			assert.Equal(t, "Hello world", result, "result")
		}
	})

	t.Run("can delete multiple times", func(t *testing.T) {
		t.Parallel()
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
		if err != nil {
			t.Fatalf("ApplyPatch() first operation error: %v", err)
		}
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() second operation error: %v", err)
		}
		assert.Equal(t, "Helloworld", result2.Doc, "result")
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("deletes entire string", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "",
				Pos:  0,
				Len:  5,
			}
			result := testutils.ApplyInternalOp(t, "hello", operation)
			assert.Equal(t, "", result, "result")
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  0,
				Len:  7,
			}
			result := testutils.ApplyInternalOp(t, map[string]any{"msg": "Hello, world!"}, operation)
			expected := map[string]any{"msg": "world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  5,
				Len:  8,
			}
			result := testutils.ApplyInternalOp(t, map[string]any{"msg": "Hello, world!"}, operation)
			expected := map[string]any{"msg": "Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  5,
				Len:  10,
			}
			result := testutils.ApplyInternalOp(t, map[string]any{"msg": "Hello beautiful world"}, operation)
			expected := map[string]any{"msg": "Hello world"}
			assert.Equal(t, expected, result)
		})

		t.Run("negative position counts from end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/msg",
				Pos:  -1,
				Len:  1,
			}
			result := testutils.ApplyInternalOp(t, map[string]any{"msg": "Hello!"}, operation)
			expected := map[string]any{"msg": "Hello"}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("deletes characters from the beginning", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  0,
				Len:  7,
			}
			result := testutils.ApplyInternalOp(t, []any{"Hello, world!"}, operation)
			expected := []any{"world!"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  5,
				Len:  8,
			}
			result := testutils.ApplyInternalOp(t, []any{"Hello, world!"}, operation)
			expected := []any{"Hello"}
			assert.Equal(t, expected, result)
		})

		t.Run("deletes characters from the middle", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_del",
				Path: "/0",
				Pos:  5,
				Len:  10,
			}
			result := testutils.ApplyInternalOp(t, []any{"Hello beautiful world"}, operation)
			expected := []any{"Hello world"}
			assert.Equal(t, expected, result)
		})
	})
}
