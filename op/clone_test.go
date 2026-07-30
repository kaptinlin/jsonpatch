package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestCloneIsolatesOperationState(t *testing.T) {
	t.Parallel()

	t.Run("path", func(t *testing.T) {
		t.Parallel()

		operation := NewAdd([]string{"name"}, "Ada")
		cloned := mustClone[*AddOperation](t, operation)
		operation.path[0] = "changed"
		assertCloneEqual(t, []string{"name"}, cloned.Path())
	})

	t.Run("from", func(t *testing.T) {
		t.Parallel()

		operation := NewMove([]string{"to"}, []string{"from"})
		cloned := mustClone[*MoveOperation](t, operation)
		operation.from[0] = "changed"
		assertCloneEqual(t, []string{"from"}, cloned.From())
	})

	t.Run("value", func(t *testing.T) {
		t.Parallel()

		value := map[string]any{"name": "Ada"}
		operation := NewAdd(nil, value)
		cloned := mustClone[*AddOperation](t, operation)
		value["name"] = "changed"
		assertCloneEqual(t, map[string]any{"name": "Ada"}, cloned.Value)
	})

	t.Run("old value", func(t *testing.T) {
		t.Parallel()

		oldValue := map[string]any{"name": "Ada"}
		operation := NewReplaceWithOldValue(nil, "Grace", oldValue)
		cloned := mustClone[*ReplaceOperation](t, operation)
		oldValue["name"] = "changed"
		assertCloneEqual(t, map[string]any{"name": "Ada"}, cloned.OldValue)
	})

	t.Run("props", func(t *testing.T) {
		t.Parallel()

		props := map[string]any{"kind": "paragraph"}
		operation := NewMerge([]string{"1"}, 1, props)
		cloned := mustClone[*MergeOperation](t, operation)
		props["kind"] = "changed"
		assertCloneEqual(t, map[string]any{"kind": "paragraph"}, cloned.Props)
	})

	t.Run("composite predicate", func(t *testing.T) {
		t.Parallel()

		operand := NewContains([]string{"name"}, "Ada")
		operation := NewAnd(nil, operand)
		cloned := mustClone[*AndOperation](t, operation)
		operand.Value = "changed"

		clonedOperand, ok := cloned.Operations[0].(*ContainsOperation)
		require.True(t, ok)
		assert.Equal(t, "Ada", clonedOperand.Value)
	})

	t.Run("not operand", func(t *testing.T) {
		t.Parallel()

		operand := NewContains([]string{"name"}, "Ada")
		operation := NewNot(nil, operand)
		cloned := mustClone[*NotOperation](t, operation)
		operand.Value = "changed"

		clonedOperand, ok := cloned.Operand.(*ContainsOperation)
		require.True(t, ok)
		assert.Equal(t, "Ada", clonedOperand.Value)
	})
}

func mustClone[T any](t *testing.T, operation interface{ Clone() (internal.Op, error) }) T {
	t.Helper()

	cloned, err := operation.Clone()
	require.NoError(t, err)
	typed, ok := cloned.(T)
	require.True(t, ok)
	return typed
}

func assertCloneEqual(t *testing.T, want, got any) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Clone() mismatch (-want +got):\n%s", diff)
	}
}
