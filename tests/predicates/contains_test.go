package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestContainsOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			op := internal.Operation{
				Op:    "contains",
				Path:  "",
				Value: "oo b",
			}
			testutils.ApplyInternalOp(t, "foo bar", op)
		})

		t.Run("succeeds when matches start of the string", func(t *testing.T) {
			op := internal.Operation{
				Op:    "contains",
				Path:  "",
				Value: "foo",
			}
			testutils.ApplyInternalOp(t, "foo bar", op)
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				Op:         "contains",
				Path:       "",
				Value:      "oo B",
				IgnoreCase: true,
			}
			testutils.ApplyInternalOp(t, "foo bar", op)
		})

		t.Run("throws when case does not match", func(t *testing.T) {
			op := internal.Operation{
				Op:    "contains",
				Path:  "",
				Value: "oo B",
			}
			testutils.ApplyInternalOpWithError(t, "foo bar", op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			op := internal.Operation{
				Op:    "contains",
				Path:  "",
				Value: "oo 0",
			}
			testutils.ApplyInternalOpWithError(t, "foo bar", op)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			obj := map[string]interface{}{"foo": "foo bar"}
			op := internal.Operation{
				Op:    "contains",
				Path:  "/foo",
				Value: "oo b",
			}
			testutils.ApplyInternalOp(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			obj := map[string]interface{}{"foo": "foo bar"}
			op := internal.Operation{
				Op:    "contains",
				Path:  "/foo",
				Value: "oo 0",
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			arr := []interface{}{"foo bar"}
			op := internal.Operation{
				Op:    "contains",
				Path:  "/0",
				Value: "oo b",
			}
			testutils.ApplyInternalOp(t, arr, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			arr := []interface{}{"foo bar"}
			op := internal.Operation{
				Op:    "contains",
				Path:  "/0",
				Value: "oo 0",
			}
			testutils.ApplyInternalOpWithError(t, arr, op)
		})
	})
}
