package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMergeOp(t *testing.T) {
	t.Run("can merge two nodes in an array", func(t *testing.T) {
		state := []interface{}{
			map[string]interface{}{"text": "foo"},
			map[string]interface{}{"text": "bar"},
		}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		require.NoError(t, err)
		expected := []interface{}{
			map[string]interface{}{"text": "foobar"},
		}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("cannot target first array element when merging", func(t *testing.T) {
		state := []interface{}{
			map[string]interface{}{"text": "foo"},
			map[string]interface{}{"text": "bar"},
		}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  0,
			},
		}
		_, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		require.Error(t, err)
	})

	t.Run("can merge slate element nodes", func(t *testing.T) {
		state := map[string]interface{}{
			"foo": []interface{}{
				map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "1"}, map[string]interface{}{"text": "2"}}},
				map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "1"}, map[string]interface{}{"text": "2"}}},
				map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "3"}, map[string]interface{}{"text": "4"}}},
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
		require.NoError(t, err)
		expected := map[string]interface{}{
			"foo": []interface{}{
				map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "1"}, map[string]interface{}{"text": "2"}}},
				map[string]interface{}{"children": []interface{}{map[string]interface{}{"text": "1"}, map[string]interface{}{"text": "2"}, map[string]interface{}{"text": "3"}, map[string]interface{}{"text": "4"}}},
			},
		}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("cannot merge root", func(t *testing.T) {
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		_, err := jsonpatch.ApplyPatch(123, operations, internal.WithMutate(true))
		require.Error(t, err)
	})

	t.Run("can merge strings", func(t *testing.T) {
		state := []interface{}{"hello", " world"}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		require.NoError(t, err)
		expected := []interface{}{"hello world"}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("can merge numbers", func(t *testing.T) {
		state := []interface{}{5, 3}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		require.NoError(t, err)
		expected := []interface{}{float64(8)}
		assert.Equal(t, expected, result.Doc)
	})

	t.Run("returns array for non-mergeable types", func(t *testing.T) {
		state := []interface{}{true, false}
		operations := []internal.Operation{
			{
				Op:   "merge",
				Path: "",
				Pos:  1,
			},
		}
		result, err := jsonpatch.ApplyPatch(state, operations, internal.WithMutate(true))
		require.NoError(t, err)
		expected := []interface{}{[]interface{}{true, false}}
		assert.Equal(t, expected, result.Doc)
	})
}
