package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationIn applies a single operation to a document
func applyOperationIn(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

// applyOperationWithErrorIn applies an operation expecting it to fail
func applyOperationWithErrorIn(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	require.Error(t, err)
}

func TestInOp(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("should test against root (on a json document of type object) - and return true", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op: "in",
				Path: "",
				Value: []interface{}{
					1,
					map[string]interface{}{
						"hello": "world",
					},
				},
			}
			result := applyOperationIn(t, obj, op)
			assert.Equal(t, obj, result)
		})

		t.Run("should test against root (on a json document of type object) - and return false", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op: "in",
				Path: "",
				Value: []interface{}{1},
			}
			applyOperationWithErrorIn(t, obj, op)
		})

		t.Run("should test against root (on a json document of type array) - and return false", func(t *testing.T) {
			obj := []interface{}{
				map[string]interface{}{
					"hello": "world",
				},
			}
			op := internal.Operation{
				Op: "in",
				Path: "",
				Value: []interface{}{1},
			}
			applyOperationWithErrorIn(t, obj, op)
		})
	})
}
