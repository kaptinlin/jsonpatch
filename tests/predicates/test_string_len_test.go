package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

// applyOperationTestStringLen applies a single operation to a document
func applyOperationTestStringLen(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithErrorTestStringLen applies an operation expecting it to fail
func applyOperationWithErrorTestStringLen(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.Error(t, err)
}

func TestTestStringLenOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  3,
				}
				applyOperationTestStringLen(t, "foo bar", op)
			})

			t.Run("succeeds when target length is equal to requested length", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  2,
				}
				applyOperationTestStringLen(t, "xo", op)
			})

			t.Run("throws when requested length is larger than target", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  9999,
				}
				applyOperationWithErrorTestStringLen(t, "asdf", op)
			})
		})

		t.Run("negative", func(t *testing.T) {
			t.Run("throw when target is longer than requested", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  3, // "foo bar" has length 7, so >= 3 is true, with not=true it should fail
					Not: true,
				}
				applyOperationWithErrorTestStringLen(t, "foo bar", op)
			})

			t.Run("throws when target length is equal to requested length", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  2,
					Not: true,
				}
				applyOperationWithErrorTestStringLen(t, "xo", op)
			})

			t.Run("succeeds when requested length is larger than target", func(t *testing.T) {
				op := internal.Operation{
					Op: "test_string_len",
					Path: "",
					Len:  9999,
					Not: true,
				}
				applyOperationTestStringLen(t, "asdf", op)
			})
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				obj := map[string]interface{}{"a": "b"}
				op := internal.Operation{
					Op: "test_string_len",
					Path: "/a",
					Len:  1,
				}
				applyOperationTestStringLen(t, obj, op)

				op2 := internal.Operation{
					Op: "test_string_len",
					Path: "/a",
					Len:  0,
				}
				applyOperationTestStringLen(t, obj, op2)
			})

			t.Run("throws when target is shorter than requested", func(t *testing.T) {
				obj := map[string]interface{}{"a": "b"}
				op := internal.Operation{
					Op: "test_string_len",
					Path: "/a",
					Len:  99,
				}
				applyOperationWithErrorTestStringLen(t, obj, op)

				// This should succeed with not=true
				op2 := internal.Operation{
					Op: "test_string_len",
					Path: "/a",
					Len:  99,
					Not: true,
				}
				applyOperationTestStringLen(t, obj, op2)
			})
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("positive", func(t *testing.T) {
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				obj := map[string]interface{}{"a": []interface{}{"b"}}
				op := internal.Operation{
					Op: "test_string_len",
					Path: "/a/0",
					Len:  1,
				}
				applyOperationTestStringLen(t, obj, op)

				op2 := internal.Operation{
					Op: "test_string_len",
					Path: "/a/0",
					Len:  0,
				}
				applyOperationTestStringLen(t, obj, op2)
			})

			t.Run("throws when target is shorter than requested", func(t *testing.T) {
				obj := map[string]interface{}{"a": []interface{}{"b"}}
				op := internal.Operation{
					Op: "test_string_len",
					Path: "/a/0",
					Len:  99,
				}
				applyOperationWithErrorTestStringLen(t, obj, op)

				// This should succeed with not=true
				op2 := internal.Operation{
					Op: "test_string_len",
					Path: "/a/0",
					Len:  99,
					Not: true,
				}
				applyOperationTestStringLen(t, obj, op2)
			})
		})
	})
}
