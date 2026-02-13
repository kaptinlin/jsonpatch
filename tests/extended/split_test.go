package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestSplitOp(t *testing.T) {
	t.Parallel()
	t.Run("Slate.js examples", func(t *testing.T) {
		t.Parallel()
		t.Run("split a single 'ab' paragraphs into two", func(t *testing.T) {
			t.Parallel()
			state := []any{
				map[string]any{
					"children": []any{
						map[string]any{
							"text": "ab",
						},
					},
				},
			}
			operations := []internal.Operation{
				{
					Op:   "split",
					Path: "/0/children/0",
					Pos:  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := []any{
				map[string]any{
					"children": []any{
						map[string]any{
							"text": "a",
						},
						map[string]any{
							"text": "b",
						},
					},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("split two element blocks into one", func(t *testing.T) {
			t.Parallel()
			state := []any{
				map[string]any{
					"children": []any{
						map[string]any{
							"text": "a",
						},
						map[string]any{
							"text": "b",
						},
					},
				},
			}
			operations := []internal.Operation{
				{
					Op:   "split",
					Path: "/0",
					Pos:  1,
				},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := []any{
				map[string]any{
					"children": []any{
						map[string]any{
							"text": "a",
						},
					},
				},
				map[string]any{
					"children": []any{
						map[string]any{
							"text": "b",
						},
					},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})
	})

	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("string", func(t *testing.T) {
			t.Parallel()
			t.Run("can split string in two", func(t *testing.T) {
				t.Parallel()
				var state any = "1234"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 2},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"12", "34"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two at pos=1", func(t *testing.T) {
				t.Parallel()
				var state any = "1234"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 1},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"1", "234"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two from beginning", func(t *testing.T) {
				t.Parallel()
				var state any = "1234"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 0},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"", "1234"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two from end", func(t *testing.T) {
				t.Parallel()
				var state any = "1234"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 4},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"1234", ""}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can split string in two when pos is greater than string length", func(t *testing.T) {
				t.Parallel()
				var state any = "12345"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 99999},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"12345", ""}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("takes characters from end if pos is negative", func(t *testing.T) {
				t.Parallel()
				var state any = "12345"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: -1},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"1234", "5"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("takes characters from end if pos is negative - 2", func(t *testing.T) {
				t.Parallel()
				var state any = "12345"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: -2},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"123", "45"}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("when negative pos overflows, first element is empty", func(t *testing.T) {
				t.Parallel()
				var state any = "12345"
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: -7},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{"", "12345"}
				assert.Equal(t, expected, result.Doc)
			})
		})

		t.Run("SlateTextNode", func(t *testing.T) {
			t.Parallel()
			t.Run("splits simple SlateTextNode", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"text": "foo bar",
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 3},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{"text": "foo"},
					map[string]any{"text": " bar"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("preserves text node attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 3},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{"text": "foo", "foo": "bar"},
					map[string]any{"text": " bar", "foo": "bar"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can add custom attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 3, Props: map[string]any{"baz": "qux"}},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{"text": "foo", "foo": "bar", "baz": "qux"},
					map[string]any{"text": " bar", "foo": "bar", "baz": "qux"},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("custom attributes can overwrite node attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"text": "foo bar",
					"foo":  "bar",
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 3, Props: map[string]any{"foo": "1"}},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{"text": "foo", "foo": "1"},
					map[string]any{"text": " bar", "foo": "1"},
				}
				assert.Equal(t, expected, result.Doc)
			})
		})

		t.Run("SlateElementNode", func(t *testing.T) {
			t.Parallel()
			t.Run("splits simple node", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"children": []any{
						map[string]any{"text": "foo"},
						map[string]any{"text": "bar"},
						map[string]any{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 1},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{
						"children": []any{
							map[string]any{"text": "foo"},
						},
					},
					map[string]any{
						"children": []any{
							map[string]any{"text": "bar"},
							map[string]any{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can provide custom attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"children": []any{
						map[string]any{"text": "foo"},
						map[string]any{"text": "bar"},
						map[string]any{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 2, Props: map[string]any{"f": 1}},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{
						"f": 1,
						"children": []any{
							map[string]any{"text": "foo"},
							map[string]any{"text": "bar"},
						},
					},
					map[string]any{
						"f": 1,
						"children": []any{
							map[string]any{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("carries over node attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"a": 1,
					"children": []any{
						map[string]any{"text": "foo"},
						map[string]any{"text": "bar"},
						map[string]any{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 2, Props: map[string]any{"f": 2}},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{
						"f": 2,
						"a": 1,
						"children": []any{
							map[string]any{"text": "foo"},
							map[string]any{"text": "bar"},
						},
					},
					map[string]any{
						"f": 2,
						"a": 1,
						"children": []any{
							map[string]any{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})

			t.Run("can overwrite node attributes", func(t *testing.T) {
				t.Parallel()
				var state any = map[string]any{
					"a": 1,
					"c": 3,
					"children": []any{
						map[string]any{"text": "foo"},
						map[string]any{"text": "bar"},
						map[string]any{"text": "baz"},
					},
				}
				operations := []internal.Operation{
					{Op: "split", Path: "", Pos: 2, Props: map[string]any{"f": 2, "a": 2}},
				}
				result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
				if err != nil {
					t.Fatalf("ApplyPatch() error: %v", err)
				}
				expected := []any{
					map[string]any{
						"f": 2,
						"a": 2,
						"c": 3,
						"children": []any{
							map[string]any{"text": "foo"},
							map[string]any{"text": "bar"},
						},
					},
					map[string]any{
						"f": 2,
						"a": 2,
						"c": 3,
						"children": []any{
							map[string]any{"text": "baz"},
						},
					},
				}
				assert.Equal(t, expected, result.Doc)
			})
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("can split string in two", func(t *testing.T) {
			t.Parallel()
			state := map[string]any{"foo": "ab"}
			operations := []internal.Operation{
				{Op: "split", Path: "/foo", Pos: 1},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := map[string]any{"foo": []any{"a", "b"}}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("if attribute are specified, wraps strings into nodes", func(t *testing.T) {
			t.Parallel()
			state := map[string]any{"foo": "ab"}
			operations := []internal.Operation{
				{Op: "split", Path: "/foo", Pos: 1, Props: map[string]any{"z": "x"}},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := map[string]any{
				"foo": []any{
					map[string]any{"text": "a", "z": "x"},
					map[string]any{"text": "b", "z": "x"},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("splits SlateTextNode", func(t *testing.T) {
			t.Parallel()
			state := map[string]any{"foo": map[string]any{"text": "777"}}
			operations := []internal.Operation{
				{Op: "split", Path: "/foo", Pos: 1, Props: map[string]any{"z": "x"}},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := map[string]any{
				"foo": []any{
					map[string]any{"text": "7", "z": "x"},
					map[string]any{"text": "77", "z": "x"},
				},
			}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("crates a tuple if target is a boolean value", func(t *testing.T) {
			t.Parallel()
			state := map[string]any{"foo": true}
			operations := []internal.Operation{
				{Op: "split", Path: "/foo", Pos: 1},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := map[string]any{"foo": []any{true, true}}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("divides number into two haves if target is a number", func(t *testing.T) {
			t.Parallel()
			state := map[string]any{"foo": 10}
			operations := []internal.Operation{
				{Op: "split", Path: "/foo", Pos: 9},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := map[string]any{"foo": []any{float64(9), float64(1)}}
			assert.Equal(t, expected, result.Doc)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("splits SlateElementNode into two", func(t *testing.T) {
			t.Parallel()
			state := []any{1, map[string]any{"children": []any{map[string]any{"text": "a"}, map[string]any{"text": "b"}}}, 2}
			operations := []internal.Operation{
				{Op: "split", Path: "/1", Pos: 0},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := []any{1, map[string]any{"children": []any{}}, map[string]any{"children": []any{map[string]any{"text": "a"}, map[string]any{"text": "b"}}}, 2}
			assert.Equal(t, expected, result.Doc)
		})

		t.Run("adds custom props and preserves node props", func(t *testing.T) {
			t.Parallel()
			state := []any{1, map[string]any{"foo": "bar", "children": []any{map[string]any{"text": "a"}, map[string]any{"text": "b"}}}, 2}
			operations := []internal.Operation{
				{Op: "split", Path: "/1", Pos: 0, Props: map[string]any{"a": "b"}},
			}
			result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error: %v", err)
			}
			expected := []any{
				1,
				map[string]any{"foo": "bar", "a": "b", "children": []any{}},
				map[string]any{"foo": "bar", "a": "b", "children": []any{map[string]any{"text": "a"}, map[string]any{"text": "b"}}},
				2,
			}
			assert.Equal(t, expected, result.Doc)
		})
	})
}
