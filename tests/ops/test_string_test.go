package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

// applyOperationTestString applies a single operation to a document
func applyOperationTestString(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithErrorTestString applies an operation expecting it to fail
func applyOperationWithErrorTestString(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.Error(t, err)
}

func TestTestString(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_string",
				"path": "",
				"pos":  1,
				"str":  "oo b",
			}
			applyOperationTestString(t, "foo bar", op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_string",
				"path": "",
				"pos":  3,
				"str":  "oo",
			}
			applyOperationWithErrorTestString(t, "foo bar", op)

			// This should succeed
			op2 := internal.Operation{
				"op":   "test_string",
				"path": "",
				"pos":  4,
				"str":  "bar",
			}
			applyOperationTestString(t, "foo bar", op2)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			obj := map[string]interface{}{"a": "b", "test": "foo bar"}
			op := internal.Operation{
				"op":   "test_string",
				"path": "/test",
				"pos":  1,
				"str":  "oo b",
			}
			applyOperationTestString(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			obj := map[string]interface{}{"test": "foo bar"}
			op := internal.Operation{
				"op":   "test_string",
				"path": "/test",
				"pos":  3,
				"str":  "oo",
			}
			applyOperationWithErrorTestString(t, obj, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			obj := map[string]interface{}{"a": "b", "test": []interface{}{"foo bar"}}
			op := internal.Operation{
				"op":   "test_string",
				"path": "/test/0",
				"pos":  1,
				"str":  "oo b",
			}
			applyOperationTestString(t, obj, op)
		})

		t.Run("throws when matches substring incorrectly", func(t *testing.T) {
			obj := map[string]interface{}{"test": []interface{}{"foo bar"}}
			op := internal.Operation{
				"op":   "test_string",
				"path": "/test/0",
				"pos":  3,
				"str":  "oo",
			}
			applyOperationWithErrorTestString(t, obj, op)
		})
	})
}
