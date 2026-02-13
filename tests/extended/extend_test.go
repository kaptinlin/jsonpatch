package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestExtendOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("can extend an object", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "",
					Props: map[string]any{
						"a": "b",
						"c": 3,
					},
				},
			}
			result := testutils.ApplyInternalOps(t, map[string]any{"foo": "bar"}, operations)
			expected := map[string]any{
				"foo": "bar",
				"a":   "b",
				"c":   3,
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("can extend an object", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/0/lol",
					Props: map[string]any{
						"b": 123,
					},
				},
			}
			doc := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
							"a": 1,
						},
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
							"a": 1,
							"b": 123,
						},
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can set null", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/0/lol",
					Props: map[string]any{
						"b": 123,
						"c": nil,
						"a": nil,
					},
				},
			}
			doc := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
							"a": 1,
						},
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
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
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/0/lol",
					Props: map[string]any{
						"b": 123,
						"c": nil,
						"a": nil,
					},
					DeleteNull: true,
				},
			}
			doc := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
							"a": 1,
						},
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": []any{
					map[string]any{
						"lol": map[string]any{
							"b": 123,
						},
					},
				},
			}
			assert.Equal(t, expected, result)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("can extend an object", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/lol",
					Props: map[string]any{
						"b": 123,
					},
				},
			}
			doc := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"a": 1,
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"a": 1,
						"b": 123,
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can set null", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/lol",
					Props: map[string]any{
						"b": 123,
						"c": nil,
						"a": nil,
					},
				},
			}
			doc := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"a": 1,
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"a": nil,
						"b": 123,
						"c": nil,
					},
				},
			}
			assert.Equal(t, expected, result)
		})

		t.Run("can use null to delete a key", func(t *testing.T) {
			t.Parallel()
			operations := []internal.Operation{
				{
					Op:   "extend",
					Path: "/foo/lol",
					Props: map[string]any{
						"b": 123,
						"c": nil,
						"a": nil,
					},
					DeleteNull: true,
				},
			}
			doc := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"a": 1,
					},
				},
			}
			result := testutils.ApplyInternalOps(t, doc, operations)
			expected := map[string]any{
				"foo": map[string]any{
					"lol": map[string]any{
						"b": 123,
					},
				},
			}
			assert.Equal(t, expected, result)
		})
	})
}
