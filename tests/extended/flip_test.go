package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestFlipOp(t *testing.T) {
	t.Run("casts values and them flips them", func(t *testing.T) {
		doc := map[string]interface{}{
			"val1": true,
			"val2": false,
			"val3": 1,
			"val4": 0,
		}
		operations := []internal.Operation{
			{Op: "flip", Path: "/val1"},
			{Op: "flip", Path: "/val2"},
			{Op: "flip", Path: "/val3"},
			{Op: "flip", Path: "/val4"},
		}
		result := testutils.ApplyInternalOps(t, doc, operations)
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
				Op:   "flip",
				Path: "",
			}
			result := testutils.ApplyInternalOps(t, true, []internal.Operation{operation})
			assert.Equal(t, false, result)
		})

		t.Run("flips false to true", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "",
			}
			result := testutils.ApplyInternalOps(t, false, []internal.Operation{operation})
			assert.Equal(t, true, result)
		})

		t.Run("flips truthy number to false", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "",
			}
			result := testutils.ApplyInternalOps(t, 123, []internal.Operation{operation})
			assert.Equal(t, false, result)
		})

		t.Run("flips zero to true", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "",
			}
			result := testutils.ApplyInternalOps(t, 0, []internal.Operation{operation})
			assert.Equal(t, true, result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("flips true to false", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "/foo",
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"foo": true}, []internal.Operation{operation})
			expected := map[string]interface{}{"foo": false}
			assert.Equal(t, expected, result)
		})

		t.Run("flips false to true", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "/foo",
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{"foo": false}, []internal.Operation{operation})
			expected := map[string]interface{}{"foo": true}
			assert.Equal(t, expected, result)
		})

		t.Run("treats empty arrays and objects as truthy", func(t *testing.T) {
			operations := []internal.Operation{
				{Op: "flip", Path: "/empty_array"},
				{Op: "flip", Path: "/empty_object"},
			}
			doc := map[string]interface{}{
				"empty_array":  []interface{}{},
				"empty_object": map[string]interface{}{},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]interface{}{
				"empty_array":  false, // empty array is truthy -> false
				"empty_object": false, // empty object is truthy -> false
			}
			assert.Equal(t, expected, result)
		})

		t.Run("creates value when path doesn't exist", func(t *testing.T) {
			operation := internal.Operation{
				Op:   "flip",
				Path: "/newfield",
			}
			result := testutils.ApplyInternalOps(t, map[string]interface{}{}, []internal.Operation{operation})
			expected := map[string]interface{}{"newfield": true}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("flips true to false and back", func(t *testing.T) {
			operations := []internal.Operation{
				{
					Op:   "flip",
					Path: "/0",
				},
				{
					Op:   "flip",
					Path: "/1",
				},
			}
			result := testutils.ApplyInternalOps(t, []interface{}{true, false}, operations)
			expected := []interface{}{false, true}
			assert.Equal(t, expected, result)
		})
	})
}
