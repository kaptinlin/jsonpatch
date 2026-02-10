package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestTestStringLenOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("positive", func(t *testing.T) {
			t.Parallel()
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  3,
				}
				testutils.ApplyInternalOp(t, "foo bar", op)
			})

			t.Run("succeeds when target length is equal to requested length", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  2,
				}
				testutils.ApplyInternalOp(t, "xo", op)
			})

			t.Run("throws when requested length is larger than target", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  9999,
				}
				testutils.ApplyInternalOpWithError(t, "asdf", op)
			})
		})

		t.Run("negative", func(t *testing.T) {
			t.Parallel()
			t.Run("throw when target is longer than requested", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  3, // "foo bar" has length 7, so >= 3 is true, with not=true it should fail
					Not:  true,
				}
				testutils.ApplyInternalOpWithError(t, "foo bar", op)
			})

			t.Run("throws when target length is equal to requested length", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  2,
					Not:  true,
				}
				testutils.ApplyInternalOpWithError(t, "xo", op)
			})

			t.Run("succeeds when requested length is larger than target", func(t *testing.T) {
				t.Parallel()
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "",
					Len:  9999,
					Not:  true,
				}
				testutils.ApplyInternalOp(t, "asdf", op)
			})
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("positive", func(t *testing.T) {
			t.Parallel()
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				t.Parallel()
				obj := map[string]interface{}{"a": "b"}
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "/a",
					Len:  1,
				}
				testutils.ApplyInternalOp(t, obj, op)

				op2 := internal.Operation{
					Op:   "test_string_len",
					Path: "/a",
					Len:  0,
				}
				testutils.ApplyInternalOp(t, obj, op2)
			})

			t.Run("throws when target is shorter than requested", func(t *testing.T) {
				t.Parallel()
				obj := map[string]interface{}{"a": "b"}
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "/a",
					Len:  99,
				}
				testutils.ApplyInternalOpWithError(t, obj, op)

				// This should succeed with not=true
				op2 := internal.Operation{
					Op:   "test_string_len",
					Path: "/a",
					Len:  99,
					Not:  true,
				}
				testutils.ApplyInternalOp(t, obj, op2)
			})
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("positive", func(t *testing.T) {
			t.Parallel()
			t.Run("succeeds when target is longer than requested", func(t *testing.T) {
				t.Parallel()
				obj := map[string]interface{}{"a": []interface{}{"b"}}
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "/a/0",
					Len:  1,
				}
				testutils.ApplyInternalOp(t, obj, op)

				op2 := internal.Operation{
					Op:   "test_string_len",
					Path: "/a/0",
					Len:  0,
				}
				testutils.ApplyInternalOp(t, obj, op2)
			})

			t.Run("throws when target is shorter than requested", func(t *testing.T) {
				t.Parallel()
				obj := map[string]interface{}{"a": []interface{}{"b"}}
				op := internal.Operation{
					Op:   "test_string_len",
					Path: "/a/0",
					Len:  99,
				}
				testutils.ApplyInternalOpWithError(t, obj, op)

				// This should succeed with not=true
				op2 := internal.Operation{
					Op:   "test_string_len",
					Path: "/a/0",
					Len:  99,
					Not:  true,
				}
				testutils.ApplyInternalOp(t, obj, op2)
			})
		})
	})
}
