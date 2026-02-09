package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReplace_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{
			"nested": "value",
		},
	}

	replaceOp := NewReplace([]string{"foo"}, "new_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err)

	modifiedDoc := result.Doc.(map[string]any)
	assert.Equal(t, "bar", result.Old)
	assert.Equal(t, "new_value", modifiedDoc["foo"])
	assert.Equal(t, 123, modifiedDoc["baz"])
}

func TestReplace_Nested(t *testing.T) {
	doc := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
			"qux": 123,
		},
	}

	replaceOp := NewReplace([]string{"foo", "bar"}, "new_nested_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err)

	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	assert.Equal(t, "baz", result.Old)
	assert.Equal(t, "new_nested_value", foo["bar"])
	assert.Equal(t, 123, foo["qux"])
}

func TestReplace_Array(t *testing.T) {
	doc := []any{
		"first",
		"second",
		"third",
	}

	replaceOp := NewReplace([]string{"1"}, "new_second")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err)

	modifiedArray := result.Doc.([]any)
	assert.Equal(t, "second", result.Old)
	assert.Equal(t, "new_second", modifiedArray[1])
	assert.Equal(t, "first", modifiedArray[0])
	assert.Equal(t, "third", modifiedArray[2])
}

func TestReplace_NonExistent(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{"qux"}, "new_value")
	_, err := replaceOp.Apply(doc)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPathNotFound)
}

func TestReplace_EmptyPath(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{}, "new_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err)
	assert.Equal(t, "new_value", result.Doc)
	assert.Equal(t, doc, result.Old)
}

func TestReplace_InterfaceMethods(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	assert.Equal(t, internal.OpReplaceType, replaceOp.Op())
	assert.Equal(t, internal.OpReplaceCode, replaceOp.Code())
	assert.Equal(t, []string{"test"}, replaceOp.Path())
	assert.Equal(t, "value", replaceOp.Value)
}

func TestReplace_ToJSON(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	got, err := replaceOp.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, "replace", got.Op)
	assert.Equal(t, "/test", got.Path)
	assert.Equal(t, "value", got.Value)
}

func TestReplace_ToCompact(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	compact, err := replaceOp.ToCompact()
	require.NoError(t, err)
	require.Len(t, compact, 3)
	assert.Equal(t, internal.OpReplaceCode, compact[0])
	assert.Equal(t, []string{"test"}, compact[1])
	assert.Equal(t, "value", compact[2])
}

func TestReplace_Validate(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")
	err := replaceOp.Validate()
	assert.NoError(t, err)

	replaceOp = NewReplace([]string{}, "value")
	err = replaceOp.Validate()
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrPathEmpty)
}
