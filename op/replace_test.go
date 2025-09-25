package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpReplace_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]interface{}{
			"nested": "value",
		},
	}

	// Test replacing a simple field
	replaceOp := NewReplace([]string{"foo"}, "new_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err, "Replace should succeed for existing field")

	// Check that the field was replaced
	modifiedDoc := result.Doc.(map[string]interface{})
	assert.Equal(t, "bar", result.Old, "Old value should be returned")
	assert.Equal(t, "new_value", modifiedDoc["foo"], "Field should be replaced")
	assert.Equal(t, 123, modifiedDoc["baz"], "Other fields should remain unchanged")
}

func TestOpReplace_Nested(t *testing.T) {
	// Create a test document with nested structure
	doc := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
			"qux": 123,
		},
	}

	// Test replacing a nested field
	replaceOp := NewReplace([]string{"foo", "bar"}, "new_nested_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err, "Replace should succeed for existing nested field")

	// Check that the nested field was replaced
	modifiedDoc := result.Doc.(map[string]interface{})
	foo := modifiedDoc["foo"].(map[string]interface{})
	assert.Equal(t, "baz", result.Old, "Old value should be returned")
	assert.Equal(t, "new_nested_value", foo["bar"], "Nested field should be replaced")
	assert.Equal(t, 123, foo["qux"], "Other nested fields should remain unchanged")
}

func TestOpReplace_Array(t *testing.T) {
	// Create a test document with array
	doc := []interface{}{
		"first",
		"second",
		"third",
	}

	// Test replacing an array element
	replaceOp := NewReplace([]string{"1"}, "new_second")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err, "Replace should succeed for existing array element")

	// Check that the element was replaced
	modifiedArray := result.Doc.([]interface{})
	assert.Equal(t, "second", result.Old, "Old value should be returned")
	assert.Equal(t, "new_second", modifiedArray[1], "Array element should be replaced")
	assert.Equal(t, "first", modifiedArray[0], "Other elements should remain unchanged")
	assert.Equal(t, "third", modifiedArray[2], "Other elements should remain unchanged")
}

func TestOpReplace_NonExistent(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test replacing a non-existent field
	replaceOp := NewReplace([]string{"qux"}, "new_value")
	_, err := replaceOp.Apply(doc)
	assert.Error(t, err, "Replace should fail for non-existent field")
	assert.Contains(t, err.Error(), "NOT_FOUND", "Error message should be descriptive")
}

func TestOpReplace_EmptyPath(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test replacing with empty path
	replaceOp := NewReplace([]string{}, "new_value")
	result, err := replaceOp.Apply(doc)
	require.NoError(t, err, "Replace should succeed for empty path (root replacement)")
	assert.Equal(t, "new_value", result.Doc, "Document should be replaced entirely")
	assert.Equal(t, doc, result.Old, "Old value should be the entire document")
}

func TestOpReplace_InterfaceMethods(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	// Test Op() method
	assert.Equal(t, internal.OpReplaceType, replaceOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpReplaceCode, replaceOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, replaceOp.Path(), "Path() should return correct path")

	// Test Value() method
	assert.Equal(t, "value", replaceOp.Value, "Value() should return correct value")
}

func TestOpReplace_ToJSON(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	json, err := replaceOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "replace", json.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", json.Path, "JSON should contain correct formatted path")
	assert.Equal(t, "value", json.Value, "JSON should contain correct value")
}

func TestOpReplace_ToCompact(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	// Test verbose format
	compact, err := replaceOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpReplaceCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
	assert.Equal(t, "value", compact[2], "Third element should be value")

	// Test non-verbose format
	compact, err = replaceOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
}

func TestOpReplace_Validate(t *testing.T) {
	// Test valid operation
	replaceOp := NewReplace([]string{"test"}, "value")
	err := replaceOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	replaceOp = NewReplace([]string{}, "value")
	err = replaceOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "OP_PATH_INVALID", "Error message should mention empty path")
}
