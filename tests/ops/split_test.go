package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitOp(t *testing.T) {
	t.Run("Slate.js examples", func(t *testing.T) {
		t.Run("split a single 'ab' paragraphs into two", func(t *testing.T) {
			state := []interface{}{
				map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{
							"text": "ab",
						},
					},
				},
			}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/0/children/0",
					"pos":  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := []interface{}{
				map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{
							"text": "a",
						},
						map[string]interface{}{
							"text": "b",
						},
					},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("split two element blocks into one", func(t *testing.T) {
			state := []interface{}{
				map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{
							"text": "a",
						},
						map[string]interface{}{
							"text": "b",
						},
					},
				},
			}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/0",
					"pos":  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := []interface{}{
				map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{
							"text": "a",
						},
					},
				},
				map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{
							"text": "b",
						},
					},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})
	})

	t.Run("root", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			t.Run("can split string in two", func(t *testing.T) {
				state := "1234"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  2,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"12", "34"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two at pos=1", func(t *testing.T) {
				state := "1234"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  1,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"1", "234"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two from beginning", func(t *testing.T) {
				state := "1234"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  0,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"", "1234"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two from end", func(t *testing.T) {
				state := "1234"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  4,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"1234", ""}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two when pos is greater than string length", func(t *testing.T) {
				state := "12345"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  99999,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"12345", ""}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("takes characters from end if pos is negative", func(t *testing.T) {
				state := "12345"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  -1,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"1234", "5"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("takes characters from end if pos is negative - 2", func(t *testing.T) {
				state := "12345"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  -2,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"123", "45"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("when negative pos overflows, first element is empty", func(t *testing.T) {
				state := "12345"
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  -7,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{"", "12345"}
				assert.Equal(t, expected, result.Doc)
			})
		})

		t.Run("SlateTextNode", func(t *testing.T) {
			t.Run("splits simple SlateTextNode", func(t *testing.T) {
				state := map[string]interface{}{
					"text": "foo bar",
				}
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  3,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{"text": "foo"},
					map[string]interface{}{"text": " bar"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("preserves text node attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  3,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{"text": "foo", "foo": "bar"},
					map[string]interface{}{"text": " bar", "foo": "bar"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can add custom attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{
						"op":    "split",
						"path":  "",
						"pos":   3,
						"props": map[string]interface{}{"baz": "qux"},
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{"text": "foo", "foo": "bar", "baz": "qux"},
					map[string]interface{}{"text": " bar", "foo": "bar", "baz": "qux"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("custom attributes can overwrite node attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{
						"op":    "split",
						"path":  "",
						"pos":   3,
						"props": map[string]interface{}{"foo": "1"},
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{"text": "foo", "foo": "1"},
					map[string]interface{}{"text": " bar", "foo": "1"},
				}
				assert.Equal(t, expected, result.Doc)
			})
		})

		t.Run("SlateElementNode", func(t *testing.T) {
			t.Run("splits simple node", func(t *testing.T) {
				state := map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{"text": "foo"},
						map[string]interface{}{"text": "bar"},
						map[string]interface{}{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{
						"op":   "split",
						"path": "",
						"pos":  1,
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{
						"children": []interface{}{
							map[string]interface{}{"text": "foo"},
						},
					},
					map[string]interface{}{
						"children": []interface{}{
							map[string]interface{}{"text": "bar"},
							map[string]interface{}{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can provide custom attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"children": []interface{}{
						map[string]interface{}{"text": "foo"},
						map[string]interface{}{"text": "bar"},
						map[string]interface{}{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{
						"op":    "split",
						"path":  "",
						"pos":   2,
						"props": map[string]interface{}{"f": 1},
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{
						"f": 1,
						"children": []interface{}{
							map[string]interface{}{"text": "foo"},
							map[string]interface{}{"text": "bar"},
						},
					},
					map[string]interface{}{
						"f": 1,
						"children": []interface{}{
							map[string]interface{}{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("carries over node attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"a": 1,
					"children": []interface{}{
						map[string]interface{}{"text": "foo"},
						map[string]interface{}{"text": "bar"},
						map[string]interface{}{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{
						"op":    "split",
						"path":  "",
						"pos":   2,
						"props": map[string]interface{}{"f": 2},
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{
						"f": 2,
						"a": 1,
						"children": []interface{}{
							map[string]interface{}{"text": "foo"},
							map[string]interface{}{"text": "bar"},
						},
					},
					map[string]interface{}{
						"f": 2,
						"a": 1,
						"children": []interface{}{
							map[string]interface{}{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can overwrite node attributes", func(t *testing.T) {
				state := map[string]interface{}{
					"a": 1,
					"c": 3,
					"children": []interface{}{
						map[string]interface{}{"text": "foo"},
						map[string]interface{}{"text": "bar"},
						map[string]interface{}{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{
						"op":    "split",
						"path":  "",
						"pos":   2,
						"props": map[string]interface{}{"f": 2, "a": 2},
					},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
				require.NoError(t, err)
				expected := []interface{}{
					map[string]interface{}{
						"f": 2,
						"a": 2,
						"c": 3,
						"children": []interface{}{
							map[string]interface{}{"text": "foo"},
							map[string]interface{}{"text": "bar"},
						},
					},
					map[string]interface{}{
						"f": 2,
						"a": 2,
						"c": 3,
						"children": []interface{}{
							map[string]interface{}{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("can split string in two", func(t *testing.T) {
			state := map[string]interface{}{"foo": "ab"}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/foo",
					"pos":  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := map[string]interface{}{"foo": []interface{}{"a", "b"}}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("if attribute are specified, wraps strings into nodes", func(t *testing.T) {
			state := map[string]interface{}{"foo": "ab"}
			operations := []internal.Operation{
				{
					"op":    "split",
					"path":  "/foo",
					"pos":   1,
					"props": map[string]interface{}{"z": "x"},
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{"text": "a", "z": "x"},
					map[string]interface{}{"text": "b", "z": "x"},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("splits SlateTextNode", func(t *testing.T) {
			state := map[string]interface{}{"foo": map[string]interface{}{"text": "777"}}
			operations := []internal.Operation{
				{
					"op":    "split",
					"path":  "/foo",
					"pos":   1,
					"props": map[string]interface{}{"z": "x"},
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := map[string]interface{}{
				"foo": []interface{}{
					map[string]interface{}{"text": "7", "z": "x"},
					map[string]interface{}{"text": "77", "z": "x"},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("crates a tuple if target is a boolean value", func(t *testing.T) {
			state := map[string]interface{}{"foo": true}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/foo",
					"pos":  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := map[string]interface{}{"foo": []interface{}{true, true}}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("divides number into two haves if target is a number", func(t *testing.T) {
			state := map[string]interface{}{"foo": 10}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/foo",
					"pos":  9,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := map[string]interface{}{"foo": []interface{}{float64(9), float64(1)}}
			assert.Equal(t, expected, result.Doc)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("splits SlateElementNode into two", func(t *testing.T) {
			state := []interface{}{1, map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "a"}, map[string]interface{}{"text": "b"}}}, 2}
			operations := []internal.Operation{
				{
					"op":   "split",
					"path": "/1",
					"pos":  0,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := []interface{}{1, map[string]interface{}{"children": []interface{}{}}, map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "a"}, map[string]interface{}{"text": "b"}}}, 2}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("adds custom props and preserves node props", func(t *testing.T) {
			state := []interface{}{1, map[string]interface{}{"foo": "bar", "children": []interface{}{map[string]interface{}{"text": "a"}, map[string]interface{}{"text": "b"}}}, 2}
			operations := []internal.Operation{
				{
					"op":    "split",
					"path":  "/1",
					"pos":   0,
					"props": map[string]interface{}{"a": "b"},
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
			expected := []interface{}{
				1,
				map[string]interface{}{"foo": "bar", "a": "b", "children": []interface{}{}},
				map[string]interface{}{"foo": "bar", "a": "b", "children": []interface{}{map[string]interface{}{"text": "a"}, map[string]interface{}{"text": "b"}}},
				2,
			}
			assert.Equal(t, expected, result.Doc)
		})
	})
}
