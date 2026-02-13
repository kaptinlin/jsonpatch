package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestTestTypeOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when target has correct type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"object"},
			}
			testutils.ApplyInternalOp(t, map[string]any{}, op)
		})

		t.Run("succeeds when target has correct type in list of types", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"number", "object"},
			}
			testutils.ApplyInternalOp(t, map[string]any{}, op)
		})

		t.Run("matches null as null type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"null"},
			}
			testutils.ApplyInternalOp(t, nil, op)
		})

		t.Run("does not match null as object type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"object"},
			}
			testutils.ApplyInternalOpWithError(t, nil, op)
		})

		t.Run("matches number as number type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"string", "number"},
			}
			testutils.ApplyInternalOp(t, 123, op)
		})

		t.Run("does not match number as object and string types", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"string", "object"},
			}
			testutils.ApplyInternalOpWithError(t, 123, op)
		})

		t.Run("matches float as number type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"string", "number"},
			}
			testutils.ApplyInternalOp(t, 1.2, op)
		})

		t.Run("does not match float as integer", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"integer"},
			}
			testutils.ApplyInternalOpWithError(t, 2.3, op)
		})

		t.Run("matches natural number as integer type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"integer"},
			}
			testutils.ApplyInternalOp(t, 0, op)
		})

		t.Run("does not match array as object type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"object"},
			}
			testutils.ApplyInternalOpWithError(t, []any{1, 2, 3}, op)
		})

		t.Run("does not match array as null type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"null"},
			}
			testutils.ApplyInternalOpWithError(t, []any{1, 2, 3}, op)
		})

		t.Run("matches array as array type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"null", "object", "array"},
			}
			testutils.ApplyInternalOp(t, []any{1, 2, 3}, op)
		})

		t.Run("matches boolean as boolean type", func(t *testing.T) {
			t.Parallel()
			op1 := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"boolean"},
			}
			testutils.ApplyInternalOp(t, true, op1)

			op2 := internal.Operation{
				Op:   "test_type",
				Path: "",
				Type: []any{"boolean"},
			}
			testutils.ApplyInternalOp(t, false, op2)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("matches string with string type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "/a",
				Type: []any{"string"},
			}
			testutils.ApplyInternalOp(t, map[string]any{"a": "asdf"}, op)
		})

		t.Run("does not match string as null type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "/a",
				Type: []any{"null"},
			}
			testutils.ApplyInternalOpWithError(t, map[string]any{"a": "asdf"}, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("matches string with string type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "/a/0",
				Type: []any{"string"},
			}
			testutils.ApplyInternalOp(t, map[string]any{"a": []any{"asdf"}}, op)
		})

		t.Run("does not match string as null type", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_type",
				Path: "/a/0",
				Type: []any{"null"},
			}
			testutils.ApplyInternalOpWithError(t, map[string]any{"a": []any{"asdf"}}, op)
		})
	})
}
