package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestTestString(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_string",
				Path: "",
				Pos:  1,
				Str:  "oo b",
			}
			testutils.ApplyInternalOp(t, "foo bar", op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "test_string",
				Path: "",
				Pos:  3,
				Str:  "oo",
			}
			testutils.ApplyInternalOpWithError(t, "foo bar", op)

			// This should succeed
			op2 := internal.Operation{
				Op:   "test_string",
				Path: "",
				Pos:  4,
				Str:  "bar",
			}
			testutils.ApplyInternalOp(t, "foo bar", op2)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			t.Parallel()
			obj := map[string]interface{}{"a": "b", "test": "foo bar"}
			op := internal.Operation{
				Op:   "test_string",
				Path: "/test",
				Pos:  1,
				Str:  "oo b",
			}
			testutils.ApplyInternalOp(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			t.Parallel()
			obj := map[string]interface{}{"test": "foo bar"}
			op := internal.Operation{
				Op:   "test_string",
				Path: "/test",
				Pos:  3,
				Str:  "oo",
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			t.Parallel()
			obj := map[string]interface{}{"a": "b", "test": []interface{}{"foo bar"}}
			op := internal.Operation{
				Op:   "test_string",
				Path: "/test/0",
				Pos:  1,
				Str:  "oo b",
			}
			testutils.ApplyInternalOp(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			t.Parallel()
			obj := map[string]interface{}{"test": []interface{}{"foo bar"}}
			op := internal.Operation{
				Op:   "test_string",
				Path: "/test/0",
				Pos:  3,
				Str:  "oo",
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})
}
