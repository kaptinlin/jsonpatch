package ops_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestStrInsOp(t *testing.T) {
	t.Parallel()
	t.Run("inserts a string at the beginning", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  0,
			Str:  "Hello, ",
		}
		result := testutils.ApplyInternalOp(t, "world!", operation)
		if result != "Hello, world!" {
			t.Errorf("result = %v, want %v", result, "Hello, world!")
		}
	})

	t.Run("inserts a string at the end", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  5,
			Str:  ", world",
		}
		result := testutils.ApplyInternalOp(t, "Hello", operation)
		if result != "Hello, world" {
			t.Errorf("result = %v, want %v", result, "Hello, world")
		}
	})

	t.Run("inserts a string in the middle", func(t *testing.T) {
		t.Parallel()
		operation := internal.Operation{
			Op:   "str_ins",
			Path: "",
			Pos:  5,
			Str:  " beautiful",
		}
		result := testutils.ApplyInternalOp(t, "Hello world", operation)
		if result != "Hello beautiful world" {
			t.Errorf("result = %v, want %v", result, "Hello beautiful world")
		}
	})

	t.Run("can insert multiple times", func(t *testing.T) {
		t.Parallel()
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
		if err != nil {
			t.Fatalf("ApplyPatch() first operation error: %v", err)
		}
		result2, err := jsonpatch.ApplyPatch(result1.Doc, []internal.Operation{operations[1]}, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() second operation error: %v", err)
		}
		if result2.Doc != "Hello beautiful world bright" {
			t.Errorf("result = %v, want %v", result2.Doc, "Hello beautiful world bright")
		}
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("inserts into empty string", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "",
				Pos:  0,
				Str:  "hello",
			}
			result := testutils.ApplyInternalOp(t, "", operation)
			if result != "hello" {
				t.Errorf("result = %v, want %v", result, "hello")
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("inserts a string at the beginning", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  0,
				Str:  "Hello, ",
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "world!"}, operation)
			expected := map[string]interface{}{"msg": "Hello, world!"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("inserts a string at the end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  5,
				Str:  ", world",
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello"}, operation)
			expected := map[string]interface{}{"msg": "Hello, world"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("inserts a string in the middle", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  5,
				Str:  " beautiful",
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello world"}, operation)
			expected := map[string]interface{}{"msg": "Hello beautiful world"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("negative position counts from end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/msg",
				Pos:  -1,
				Str:  "!",
			}
			result := testutils.ApplyInternalOp(t, map[string]interface{}{"msg": "Hello"}, operation)
			expected := map[string]interface{}{"msg": "Hell!o"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("inserts a string at the beginning", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  0,
				Str:  "Hello, ",
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"world!"}, operation)
			expected := []interface{}{"Hello, world!"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("inserts a string at the end", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  5,
				Str:  ", world",
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"Hello"}, operation)
			expected := []interface{}{"Hello, world"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("inserts a string in the middle", func(t *testing.T) {
			t.Parallel()
			operation := internal.Operation{
				Op:   "str_ins",
				Path: "/0",
				Pos:  5,
				Str:  " beautiful",
			}
			result := testutils.ApplyInternalOp(t, []interface{}{"Hello world"}, operation)
			expected := []interface{}{"Hello beautiful world"}
			if diff := cmp.Diff(expected, result); diff != "" {
				t.Errorf("result mismatch (-want +got):\n%s", diff)
			}
		})
	})
}
