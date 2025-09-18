package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

// applyOperation applies a single operation to a document
func applyOperationContains(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithErrorContains applies an operation expecting it to fail
func applyOperationWithErrorContains(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.Error(t, err)
}

func TestContainsOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			op := internal.Operation{
				"op":    "contains",
				"path":  "",
				"value": "oo b",
			}
			applyOperationContains(t, "foo bar", op)
		})

		t.Run("succeeds when matches start of the string", func(t *testing.T) {
			op := internal.Operation{
				"op":    "contains",
				"path":  "",
				"value": "foo",
			}
			applyOperationContains(t, "foo bar", op)
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				"op":          "contains",
				"path":        "",
				"value":       "oo B",
				"ignore_case": true,
			}
			applyOperationContains(t, "foo bar", op)
		})

		t.Run("throws when case does not match", func(t *testing.T) {
			op := internal.Operation{
				"op":    "contains",
				"path":  "",
				"value": "oo B",
			}
			applyOperationWithErrorContains(t, "foo bar", op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			op := internal.Operation{
				"op":    "contains",
				"path":  "",
				"value": "oo 0",
			}
			applyOperationWithErrorContains(t, "foo bar", op)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			obj := map[string]interface{}{"foo": "foo bar"}
			op := internal.Operation{
				"op":    "contains",
				"path":  "/foo",
				"value": "oo b",
			}
			applyOperationContains(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			obj := map[string]interface{}{"foo": "foo bar"}
			op := internal.Operation{
				"op":    "contains",
				"path":  "/foo",
				"value": "oo 0",
			}
			applyOperationWithErrorContains(t, obj, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			arr := []interface{}{"foo bar"}
			op := internal.Operation{
				"op":    "contains",
				"path":  "/0",
				"value": "oo b",
			}
			applyOperationContains(t, arr, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			arr := []interface{}{"foo bar"}
			op := internal.Operation{
				"op":    "contains",
				"path":  "/0",
				"value": "oo 0",
			}
			applyOperationWithErrorContains(t, arr, op)
		})
	})
}
