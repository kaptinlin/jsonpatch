package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// applyOperationsExtend applies multiple operations to a document
func applyOperationsExtend(t *testing.T, doc interface{}, ops []internal.Operation) interface{} {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, ops, internal.ApplyPatchOptions{Mutate: true})
	require.NoError(t, err)
	return result.Doc
}

func TestExtendOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("can extend an object", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "",
					"props": map[string]interface{}{
						"a": "b",
						"c": 3,
					},
				},
			}
			result := applyOperationsExtend(t, map[string]interface{}{"foo": "bar"}, operations)
			expected := map[string]interface{}{
				"foo": "bar",
				"a":   "b",
				"c":   3,
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("can extend an object", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/0/lol",
					"props": map[string]interface{}{
						"b": 123,
					},
				},
			}
			doc := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"a": 1,
						},
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"a": 1,
							"b": 123,
						},
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can set null", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/0/lol",
					"props": map[string]interface{}{
						"b": 123,
						"c": nil,
						"a": nil,
					},
				},
			}
			doc := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"a": 1,
						},
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"a": nil,
							"b": 123,
							"c": nil,
						},
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can use null to delete a key", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/0/lol",
					"props": map[string]interface{}{
						"b": 123,
						"c": nil,
						"a": nil,
					},
					"deleteNull": true,
				},
			}
			doc := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"a": 1,
						},
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{
						"lol": map[string]interface{}{
							"b": 123,
						},
					},
				},
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("can extend an object", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/lol",
					"props": map[string]interface{}{
						"b": 123,
					},
				},
			}
			doc := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"a": 1,
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"a": 1,
						"b": 123,
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can set null", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/lol",
					"props": map[string]interface{}{
						"b": 123,
						"c": nil,
						"a": nil,
					},
				},
			}
			doc := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"a": 1,
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"a": nil,
						"b": 123,
						"c": nil,
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can use null to delete a key", func(t *testing.T) {
			operations := []internal.Operation{
				{
					"op":   "extend",
					"path": "/foo/lol",
					"props": map[string]interface{}{
						"b": 123,
						"c": nil,
						"a": nil,
					},
					"deleteNull": true,
				},
			}
			doc := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"a": 1,
					},
				},
			}
			result := applyOperationsExtend(t, doc, operations)
			expected := map[string]interface{}{
				"foo": map[string]interface{}{
					"lol": map[string]interface{}{
						"b": 123,
					},
				},
			}
			assert.Equal(t, expected, result)
		})
	})
}
