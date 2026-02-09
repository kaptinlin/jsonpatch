package op

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRemove_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{"nested": "value"},
	}

	removeOp := NewRemove([]string{"foo"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err)

	modifiedDoc := result.Doc.(map[string]any)
	assert.Equal(t, "bar", result.Old)
	assert.NotContains(t, modifiedDoc, "foo")
	assert.Contains(t, modifiedDoc, "baz")
	assert.Contains(t, modifiedDoc, "qux")
}

func TestRemove_Nested(t *testing.T) {
	doc := map[string]any{
		"foo": map[string]any{"bar": "baz", "qux": 123},
	}

	removeOp := NewRemove([]string{"foo", "bar"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err)

	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	assert.Equal(t, "baz", result.Old)
	assert.NotContains(t, foo, "bar")
	assert.Contains(t, foo, "qux")
}

func TestRemove_Array(t *testing.T) {
	doc := []any{"first", "second", "third"}

	removeOp := NewRemove([]string{"1"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err)

	modifiedArray := result.Doc.([]any)
	assert.Equal(t, "second", result.Old)
	assert.Len(t, modifiedArray, 2)
	assert.Equal(t, "first", modifiedArray[0])
	assert.Equal(t, "third", modifiedArray[1])
}

func TestRemove_NonExistent(t *testing.T) {
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{"qux"})
	_, err := removeOp.Apply(doc)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPathNotFound)
}

func TestRemove_EmptyPath(t *testing.T) {
	doc := map[string]any{"foo": "bar"}

	removeOp := NewRemove([]string{})
	_, err := removeOp.Apply(doc)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPathEmpty)
}
