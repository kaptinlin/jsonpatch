package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

// execTestType applies a single operation to a document
func execTestType(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, []internal.Operation{op}, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// execTestTypeWithError applies an operation expecting it to fail
func execTestTypeWithError(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	_, err := jsonpatch.ApplyPatch(doc, []internal.Operation{op}, internal.WithMutate(true))
	require.Error(t, err)
}

func TestTestTypeOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when target has correct type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"object"},
			}
			execTestType(t, map[string]interface{}{}, op)
		})

		t.Run("succeeds when target has correct type in list of types", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"number", "object"},
			}
			execTestType(t, map[string]interface{}{}, op)
		})

		t.Run("matches null as null type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"null"},
			}
			execTestType(t, nil, op)
		})

		t.Run("does not match null as object type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"object"},
			}
			execTestTypeWithError(t, nil, op)
		})

		t.Run("matches number as number type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"string", "number"},
			}
			execTestType(t, 123, op)
		})

		t.Run("does not match number as object and string types", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"string", "object"},
			}
			execTestTypeWithError(t, 123, op)
		})

		t.Run("matches float as number type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"string", "number"},
			}
			execTestType(t, 1.2, op)
		})

		t.Run("does not match float as integer", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"integer"},
			}
			execTestTypeWithError(t, 2.3, op)
		})

		t.Run("matches natural number as integer type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"integer"},
			}
			execTestType(t, 0, op)
		})

		t.Run("does not match array as object type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"object"},
			}
			execTestTypeWithError(t, []interface{}{1, 2, 3}, op)
		})

		t.Run("does not match array as null type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"null"},
			}
			execTestTypeWithError(t, []interface{}{1, 2, 3}, op)
		})

		t.Run("matches array as array type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"null", "object", "array"},
			}
			execTestType(t, []interface{}{1, 2, 3}, op)
		})

		t.Run("matches boolean as boolean type", func(t *testing.T) {
			op1 := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"boolean"},
			}
			execTestType(t, true, op1)

			op2 := internal.Operation{
				"op":   "test_type",
				"path": "",
				"type": []interface{}{"boolean"},
			}
			execTestType(t, false, op2)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("matches string with string type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "/a",
				"type": []interface{}{"string"},
			}
			execTestType(t, map[string]interface{}{"a": "asdf"}, op)
		})

		t.Run("does not match string as null type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "/a",
				"type": []interface{}{"null"},
			}
			execTestTypeWithError(t, map[string]interface{}{"a": "asdf"}, op)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("matches string with string type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "/a/0",
				"type": []interface{}{"string"},
			}
			execTestType(t, map[string]interface{}{"a": []interface{}{"asdf"}}, op)
		})

		t.Run("does not match string as null type", func(t *testing.T) {
			op := internal.Operation{
				"op":   "test_type",
				"path": "/a/0",
				"type": []interface{}{"null"},
			}
			execTestTypeWithError(t, map[string]interface{}{"a": []interface{}{"asdf"}}, op)
		})
	})
}
