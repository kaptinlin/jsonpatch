package ops_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
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
		if result != "world!" {
			t.Errorf("result = %v, want %v", result, "world!")
		}
	})

	t.Run("deletes characters from the end", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  8,
		}
		result := testutils.ApplyInternalOp(t, "Hello, world!", operation)
		if result != "Hello" {
			t.Errorf("result = %v, want %v", result, "Hello")
		}
	})

	t.Run("deletes characters from the middle", func(t *testing.T) {
		operation := internal.Operation{
			Op:   "str_del",
			Path: "",
			Pos:  5,
			Len:  10,
		}
		result := testutils.ApplyInternalOp(t, "Hello beautiful world", operation)
		if result != "Hello world" {
			t.Errorf("result = %v, want %v", result, "Hello world")
		}
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
		if err != nil {
			t.Fatalf("ApplyPatch() first operation error: %v", err)
		}
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() second operation error: %v", err)
		}
		if result2.Doc != "Helloworld" {
			t.Errorf("result = %v, want %v", result2.Doc, "Helloworld")
		}
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
			if result != "" {
				t.Errorf("result = %v, want %v", result, "")
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
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
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	})
}
