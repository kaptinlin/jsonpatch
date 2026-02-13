package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestMergeOp(t *testing.T) {
	t.Parallel()
	t.Run("can merge two nodes in an array", func(t *testing.T) {
		t.Parallel()
		state := []any{
			map[string]any{"text": "foo"},
			map[string]any{"text": "bar"},
		}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() error: %v", err)
		}
		expected := []any{
			map[string]any{"text": "foobar"},
		}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("cannot target first array element when merging", func(t *testing.T) {
		t.Parallel()
		state := []any{
			map[string]any{"text": "foo"},
			map[string]any{"text": "bar"},
		}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  0,
			},
		}
		_, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err == nil {
			t.Fatal("ApplyPatch() error = nil, want error")
		}
	})

	t.Run("can merge slate element nodes", func(t *testing.T) {
		t.Parallel()
		state := map[string]any{
			"foo": []any{
				map[string]any{"children": []any{map[string]any{"text": "1"}, map[string]any{"text": "2"}}},
				map[string]any{"children": []any{map[string]any{"text": "1"}, map[string]any{"text": "2"}}},
				map[string]any{"children": []any{map[string]any{"text": "3"}, map[string]any{"text": "4"}}},
			},
		}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "/foo",
				Pos:  2,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() error: %v", err)
		}
		expected := map[string]any{
			"foo": []any{
				map[string]any{"children": []any{map[string]any{"text": "1"}, map[string]any{"text": "2"}}},
				map[string]any{"children": []any{map[string]any{"text": "1"}, map[string]any{"text": "2"}, map[string]any{"text": "3"}, map[string]any{"text": "4"}}},
			},
		}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("cannot merge root", func(t *testing.T) {
		t.Parallel()
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		_, err := jsonpatch.ApplyPatch(123, operations, internal.WithMutate(true))
		if err == nil {
			t.Fatal("ApplyPatch() error = nil, want error")
		}
	})

	t.Run("can merge strings", func(t *testing.T) {
		t.Parallel()
		state := []any{"hello", " world"}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() error: %v", err)
		}
		expected := []any{"hello world"}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("can merge numbers", func(t *testing.T) {
		t.Parallel()
		state := []any{5, 3}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() error: %v", err)
		}
		expected := []any{float64(8)}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("returns array for non-mergeable types", func(t *testing.T) {
		t.Parallel()
		state := []any{true, false}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		if err != nil {
			t.Fatalf("ApplyPatch() error: %v", err)
		}
		expected := []any{[]any{true, false}}
		assert.Equal(t, expected, result.Doc)
	})
}
