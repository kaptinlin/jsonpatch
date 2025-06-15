package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

// applyOperationType applies a single operation to a document
func applyOperationType(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.ApplyPatchOptions{Mutate: true})
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithErrorType applies an operation expecting it to fail
func applyOperationWithErrorType(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.ApplyPatchOptions{Mutate: true})
	require.Error(t, err)
}

func TestTypeOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when type matches", func(t *testing.T) {
			tests := []struct {
				value interface{}
				typ   string
			}{
				{1, "number"},
				{"hello", "string"},
				{true, "boolean"},
				{nil, "null"},
				{[]interface{}{}, "array"},
				{map[string]interface{}{}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "",
					"value": test.typ,
				}
				applyOperationType(t, test.value, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			tests := []struct {
				value interface{}
				typ   string
			}{
				{1, "string"},
				{"hello", "number"},
				{true, "null"},
				{nil, "boolean"},
				{[]interface{}{}, "object"},
				{map[string]interface{}{}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "",
					"value": test.typ,
				}
				applyOperationWithErrorType(t, test.value, op)
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when type matches", func(t *testing.T) {
			tests := []struct {
				obj map[string]interface{}
				typ string
			}{
				{map[string]interface{}{"foo": 1}, "number"},
				{map[string]interface{}{"foo": "hello"}, "string"},
				{map[string]interface{}{"foo": true}, "boolean"},
				{map[string]interface{}{"foo": nil}, "null"},
				{map[string]interface{}{"foo": []interface{}{}}, "array"},
				{map[string]interface{}{"foo": map[string]interface{}{}}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "/foo",
					"value": test.typ,
				}
				applyOperationType(t, test.obj, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			tests := []struct {
				obj map[string]interface{}
				typ string
			}{
				{map[string]interface{}{"foo": 1}, "string"},
				{map[string]interface{}{"foo": "hello"}, "number"},
				{map[string]interface{}{"foo": true}, "null"},
				{map[string]interface{}{"foo": nil}, "boolean"},
				{map[string]interface{}{"foo": []interface{}{}}, "object"},
				{map[string]interface{}{"foo": map[string]interface{}{}}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "/foo",
					"value": test.typ,
				}
				applyOperationWithErrorType(t, test.obj, op)
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when type matches", func(t *testing.T) {
			tests := []struct {
				arr []interface{}
				typ string
			}{
				{[]interface{}{1}, "number"},
				{[]interface{}{"hello"}, "string"},
				{[]interface{}{true}, "boolean"},
				{[]interface{}{nil}, "null"},
				{[]interface{}{[]interface{}{}}, "array"},
				{[]interface{}{map[string]interface{}{}}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "/0",
					"value": test.typ,
				}
				applyOperationType(t, test.arr, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			tests := []struct {
				arr []interface{}
				typ string
			}{
				{[]interface{}{1}, "string"},
				{[]interface{}{"hello"}, "number"},
				{[]interface{}{true}, "null"},
				{[]interface{}{nil}, "boolean"},
				{[]interface{}{[]interface{}{}}, "object"},
				{[]interface{}{map[string]interface{}{}}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					"op":    "type",
					"path":  "/0",
					"value": test.typ,
				}
				applyOperationWithErrorType(t, test.arr, op)
			}
		})
	})
}
