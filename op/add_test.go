package op

import (
	"reflect"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpAdd_Basic(t *testing.T) {
	// Test adding to object
	doc := map[string]interface{}{
		"foo": "bar",
	}

	op := NewAdd([]string{"baz"}, "qux")

	result, err := op.Apply(doc)
	require.NoError(t, err, "Add operation should succeed")

	// Check that the operation works directly on the document (mutate behavior)
	assert.Equal(t, "bar", doc["foo"], "Original document should preserve existing fields")
	assert.Equal(t, "qux", doc["baz"], "Original document should now have new field (mutate behavior)")

	// Check the result points to the same document
	resultDoc := result.Doc.(map[string]interface{})
	assert.Equal(t, "bar", resultDoc["foo"], "Result should preserve existing fields")
	assert.Equal(t, "qux", resultDoc["baz"], "Result should contain new field")
	assert.Nil(t, result.Old, "Old value should be nil for new field")

	// Verify it's the same object (mutate behavior)
	assert.True(t, reflect.ValueOf(doc).Pointer() == reflect.ValueOf(resultDoc).Pointer(), "Result should point to the same document object")
}

func TestOpAdd_ReplaceExisting(t *testing.T) {
	// Test replacing existing field
	doc := map[string]interface{}{
		"foo": "bar",
	}

	op := NewAdd([]string{"foo"}, "new_value")

	result, err := op.Apply(doc)
	require.NoError(t, err, "Add operation should succeed")

	// Check the result
	resultDoc := result.Doc.(map[string]interface{})
	assert.Equal(t, "new_value", resultDoc["foo"], "Result should contain new value")
	assert.Equal(t, "bar", result.Old, "Old value should be preserved")
}

func TestOpAdd_InterfaceMethods(t *testing.T) {
	op := NewAdd([]string{"foo"}, "bar")

	// Test Op() method
	assert.Equal(t, internal.OpAddType, op.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpAddCode, op.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"foo"}, op.Path(), "Path() should return correct path")
}

func TestOpAdd_ToJSON(t *testing.T) {
	op := NewAdd([]string{"foo", "bar"}, "baz")

	json, err := op.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "add", json.Op, "JSON should contain correct op type")
	assert.Equal(t, "/foo/bar", json.Path, "JSON should contain correct formatted path")
	assert.Equal(t, "baz", json.Value, "JSON should contain correct value")
}

func TestOpAdd_ToCompact(t *testing.T) {
	op := NewAdd([]string{"foo"}, "bar")

	// Test verbose format
	compact, err := op.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpAddCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.Equal(t, "bar", compact[2], "Third element should be value")

	// Test non-verbose format
	compact, err = op.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
}

func TestOpAdd_Validate(t *testing.T) {
	// Test valid operation
	op := NewAdd([]string{"foo"}, "bar")
	err := op.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	op = NewAdd([]string{}, "bar")
	err = op.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "path cannot be empty", "Error message should mention empty path")
}

func TestOpAdd_Constructor(t *testing.T) {
	path := []string{"foo", "bar"}
	value := "baz"

	op := NewAdd(path, value)

	assert.Equal(t, path, op.Path(), "Constructor should set correct path")
	// Note: we can't directly access the value field, but we can test it through ToJSON
	json, err := op.ToJSON()
	require.NoError(t, err)
	assert.Equal(t, value, json.Value, "Constructor should set correct value")
}
