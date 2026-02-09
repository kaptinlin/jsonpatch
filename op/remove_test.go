package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpRemove_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{
			"nested": "value",
		},
	}

	// Test removing a simple field
	removeOp := NewRemove([]string{"foo"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err, "Remove should succeed for existing field")

	// Check that the field was removed
	modifiedDoc := result.Doc.(map[string]any)
	assert.Equal(t, "bar", result.Old, "Old value should be returned")
	assert.NotContains(t, modifiedDoc, "foo", "Field should be removed")
	assert.Contains(t, modifiedDoc, "baz", "Other fields should remain")
	assert.Contains(t, modifiedDoc, "qux", "Other fields should remain")
}

func TestOpRemove_Nested(t *testing.T) {
	// Create a test document with nested structure
	doc := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
			"qux": 123,
		},
	}

	// Test removing a nested field
	removeOp := NewRemove([]string{"foo", "bar"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err, "Remove should succeed for existing nested field")

	// Check that the nested field was removed
	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	assert.Equal(t, "baz", result.Old, "Old value should be returned")
	assert.NotContains(t, foo, "bar", "Nested field should be removed")
	assert.Contains(t, foo, "qux", "Other nested fields should remain")
}

func TestOpRemove_Array(t *testing.T) {
	// Create a test document with array
	doc := []any{
		"first",
		"second",
		"third",
	}

	// Test removing an array element
	removeOp := NewRemove([]string{"1"})
	result, err := removeOp.Apply(doc)
	require.NoError(t, err, "Remove should succeed for existing array element")

	// Check that the element was removed
	modifiedArray := result.Doc.([]any)
	assert.Equal(t, "second", result.Old, "Old value should be returned")
	assert.Len(t, modifiedArray, 2, "Array should have one less element")
	assert.Equal(t, "first", modifiedArray[0], "First element should remain")
	assert.Equal(t, "third", modifiedArray[1], "Third element should become second")
}

func TestOpRemove_NonExistent(t *testing.T) {
	// Create a test document
	doc := map[string]any{
		"foo": "bar",
	}

	// Test removing a non-existent field
	removeOp := NewRemove([]string{"qux"})
	_, err := removeOp.Apply(doc)
	assert.Error(t, err, "Remove should fail for non-existent field")
	assert.ErrorIs(t, err, ErrPathNotFound)
}

func TestOpRemove_EmptyPath(t *testing.T) {
	// Create a test document
	doc := map[string]any{
		"foo": "bar",
	}

	// Test removing with empty path
	removeOp := NewRemove([]string{})
	_, err := removeOp.Apply(doc)
	assert.Error(t, err, "Remove should fail for empty path")
	assert.ErrorIs(t, err, ErrPathEmpty)
}

func TestOpRemove_InterfaceMethods(t *testing.T) {
	removeOp := NewRemove([]string{"test"})

	// Test Op() method
	assert.Equal(t, internal.OpRemoveType, removeOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpRemoveCode, removeOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, removeOp.Path(), "Path() should return correct path")
}

func TestOpRemove_ToJSON(t *testing.T) {
	removeOp := NewRemove([]string{"test"})

	json, err := removeOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "remove", json.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", json.Path, "JSON should contain correct formatted path")
}

func TestOpRemove_ToCompact(t *testing.T) {
	removeOp := NewRemove([]string{"test"})

	compact, err := removeOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 2, "Compact format should have 2 elements")
	assert.Equal(t, internal.OpRemoveCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
}

func TestOpRemove_Validate(t *testing.T) {
	// Test valid operation
	removeOp := NewRemove([]string{"test"})
	err := removeOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	removeOp = NewRemove([]string{})
	err = removeOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrPathEmpty)
}
