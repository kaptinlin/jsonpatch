package op

import (
	"reflect"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdd_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	addOp := NewAdd([]string{"baz"}, "qux")

	result, err := addOp.Apply(doc)
	require.NoError(t, err, "Add operation should succeed")

	assert.Equal(t, "bar", doc["foo"], "Original document should preserve existing fields")
	assert.Equal(t, "qux", doc["baz"], "Original document should now have new field (mutate behavior)")

	resultDoc := result.Doc.(map[string]any)
	assert.Equal(t, "bar", resultDoc["foo"], "Result should preserve existing fields")
	assert.Equal(t, "qux", resultDoc["baz"], "Result should contain new field")
	assert.Nil(t, result.Old, "Old value should be nil for new field")

	// Mutate behavior: result should point to the same underlying object
	assert.True(t, reflect.ValueOf(doc).Pointer() == reflect.ValueOf(resultDoc).Pointer(), "Result should point to the same document object")
}

func TestAdd_ReplaceExisting(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	addOp := NewAdd([]string{"foo"}, "new_value")

	result, err := addOp.Apply(doc)
	require.NoError(t, err, "Add operation should succeed")

	resultDoc := result.Doc.(map[string]any)
	assert.Equal(t, "new_value", resultDoc["foo"], "Result should contain new value")
	assert.Equal(t, "bar", result.Old, "Old value should be preserved")
}

func TestAdd_InterfaceMethods(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")

	assert.Equal(t, internal.OpAddType, addOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpAddCode, addOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"foo"}, addOp.Path(), "Path() should return correct path")
}

func TestAdd_ToJSON(t *testing.T) {
	addOp := NewAdd([]string{"foo", "bar"}, "baz")

	got, err := addOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "add", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/foo/bar", got.Path, "JSON should contain correct formatted path")
	assert.Equal(t, "baz", got.Value, "JSON should contain correct value")
}

func TestAdd_ToCompact(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")

	compact, err := addOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpAddCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.Equal(t, "bar", compact[2], "Third element should be value")
}

func TestAdd_Validate(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")
	err := addOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	addOp = NewAdd([]string{}, "bar")
	err = addOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrPathEmpty)
}

func TestAdd_Constructor(t *testing.T) {
	path := []string{"foo", "bar"}
	value := "baz"

	addOp := NewAdd(path, value)

	assert.Equal(t, path, addOp.Path(), "Constructor should set correct path")
	// Verify value indirectly through ToJSON since the value field is unexported
	got, err := addOp.ToJSON()
	require.NoError(t, err)
	assert.Equal(t, value, got.Value, "Constructor should set correct value")
}
