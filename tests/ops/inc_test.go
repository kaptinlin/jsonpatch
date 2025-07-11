package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationsInc applies multiple operations to a document
func applyOperationsInc(t *testing.T, doc interface{}, ops []internal.Operation) interface{} {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	require.NoError(t, err)
	return result.Doc
}

func TestIncOp(t *testing.T) {
	t.Run("casts values and then increments them", func(t *testing.T) {
		doc := map[string]interface{}{
			"val1": true,
			"val2": false,
			"val3": 1,
			"val4": 0,
		}
		operations := []internal.Operation{
			{"op": "inc", "path": "/val1", "inc": 1},
			{"op": "inc", "path": "/val2", "inc": 1},
			{"op": "inc", "path": "/val3", "inc": 1},
			{"op": "inc", "path": "/val4", "inc": 1},
		}
		result := applyOperationsInc(t, doc, operations)
		expected := map[string]interface{}{
			"val1": float64(2),
			"val2": float64(1),
			"val3": float64(2),
			"val4": float64(1),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("can use arbitrary increment value and can decrement", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": 1,
		}
		operations := []internal.Operation{
			{"op": "inc", "path": "/foo", "inc": 10},
			{"op": "inc", "path": "/foo", "inc": -3},
		}
		result := applyOperationsInc(t, doc, operations)
		expected := map[string]interface{}{
			"foo": float64(8),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("increment can be a floating point number", func(t *testing.T) {
		doc := map[string]interface{}{
			"foo": 1,
		}
		operations := []internal.Operation{
			{"op": "inc", "path": "/foo", "inc": 0.1},
		}
		result := applyOperationsInc(t, doc, operations)
		expected := map[string]interface{}{
			"foo": 1.1,
		}
		assert.Equal(t, expected, result)
	})

	t.Run("root", func(t *testing.T) {
		t.Run("increments from 0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "",
				"inc":  5,
			}
			result := applyOperationsInc(t, 0, []internal.Operation{operation})
			assert.Equal(t, float64(5), result)
		})

		t.Run("increments from -0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "",
				"inc":  5,
			}
			result := applyOperationsInc(t, -0, []internal.Operation{operation})
			assert.Equal(t, float64(5), result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("increments from 0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "/lala",
				"inc":  5,
			}
			result := applyOperationsInc(t, map[string]interface{}{"lala": 0}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(5)}
			assert.Equal(t, expected, result)
		})

		t.Run("increments from -0 to 5", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "/lala",
				"inc":  5,
			}
			result := applyOperationsInc(t, map[string]interface{}{"lala": -0}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(5)}
			assert.Equal(t, expected, result)
		})

		t.Run("casts string to number", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "/lala",
				"inc":  5,
			}
			result := applyOperationsInc(t, map[string]interface{}{"lala": "4"}, []internal.Operation{operation})
			expected := map[string]interface{}{"lala": float64(9)}
			assert.Equal(t, expected, result)
		})

		t.Run("can increment twice", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "inc",
					"path": "/lala",
					"inc":  1,
				},
				{
					"op":   "inc",
					"path": "/lala",
					"inc":  2,
				},
			}
			result := applyOperationsInc(t, map[string]interface{}{"lala": 0}, operations)
			expected := map[string]interface{}{"lala": float64(3)}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("increments from 0 to -3", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "/0",
				"inc":  -3,
			}
			result := applyOperationsInc(t, []interface{}{0}, []internal.Operation{operation})
			expected := []interface{}{float64(-3)}
			assert.Equal(t, expected, result)
		})

		t.Run("increments from -0 to -3", func(t *testing.T) {
			operation := internal.Operation{
				"op":   "inc",
				"path": "/0",
				"inc":  -3,
			}
			result := applyOperationsInc(t, []interface{}{-0}, []internal.Operation{operation})
			expected := []interface{}{float64(-3)}
			assert.Equal(t, expected, result)
		})
	})
}
