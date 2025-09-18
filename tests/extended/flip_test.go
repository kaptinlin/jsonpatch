package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationsFlip applies multiple operations to a document
func applyOperationsFlip(t *testing.T, doc interface{}, ops []internal.Operation) interface{} {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

func TestFlipOp(t *testing.T) {
	t.Run("casts values and them flips them", func(t *testing.T) {
		doc := map[string]interface{}{
			"val1": true,
			"val2": false,
			"val3": 1,
			"val4": 0,
		}
		operations := []internal.Operation{
			{"op": "flip", "path": "/val1"},
			{"op": "flip", "path": "/val2"},
			{"op": "flip", "path": "/val3"},
			{"op": "flip", "path": "/val4"},
		}
		result := applyOperationsFlip(t, doc, operations)
		expected := map[string]interface{}{
			"val1": false,
			"val2": true,
			"val3": false,
			"val4": true,
		}
		assert.Equal(t, expected, result)
	})

	t.Run("root", func(t *testing.T) {
		t.Run("flips true to false", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "",
			}
			result := applyOperationsFlip(t, true, []internal.Operation{operation})
			assert.Equal(t, false, result)
		})

		t.Run("flips false to true", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "",
			}
			result := applyOperationsFlip(t, false, []internal.Operation{operation})
			assert.Equal(t, true, result)
		})

		t.Run("flips truthy number to false", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "",
			}
			result := applyOperationsFlip(t, 123, []internal.Operation{operation})
			assert.Equal(t, false, result)
		})

		t.Run("flips zero to true", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "",
			}
			result := applyOperationsFlip(t, 0, []internal.Operation{operation})
			assert.Equal(t, true, result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("flips true to false", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "/foo",
			}
			result := applyOperationsFlip(t, map[string]interface{}{"foo": true}, []internal.Operation{operation})
			expected := map[string]interface{}{"foo": false}
			assert.Equal(t, expected, result)
		})

		t.Run("flips false to true", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "flip",
				"path": "/foo",
			}
			result := applyOperationsFlip(t, map[string]interface{}{"foo": false}, []internal.Operation{operation})
			expected := map[string]interface{}{"foo": true}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("flips true to false and back", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "flip",
					"path": "/0",
				},
				{
					"op":   "flip",
					"path": "/1",
				},
			}
			result := applyOperationsFlip(t, []interface{}{true, false}, operations)
			expected := []interface{}{false, true}
			assert.Equal(t, expected, result)
		})
	})
}
