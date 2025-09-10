package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpUndefined_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"qux": 123,
		},
	}

	// Test non-existing path
	undefinedOp := NewUndefined([]string{"qux"}, false)
	ok, err := undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "Undefined should return true for non-existing path")

	// Test existing path
	undefinedOp = NewUndefined([]string{"foo"}, false)
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.False(t, ok, "Undefined should return false for existing path")

	// Test nested non-existing path
	undefinedOp = NewUndefined([]string{"baz", "quux"}, false)
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "Undefined should return true for non-existing nested path")
}

func TestOpUndefined_Not(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test non-existing path with not=true
	undefinedOp := NewUndefined([]string{"qux"}, true)
	ok, err := undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.False(t, ok, "NOT undefined should return false for non-existing path")

	// Test existing path with not=true
	undefinedOp = NewUndefined([]string{"foo"}, true)
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "NOT undefined should return true for existing path")
}

func TestOpUndefined_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test non-existing path
	undefinedOp := NewUndefined([]string{"qux"}, false)
	result, err := undefinedOp.Apply(doc)
	require.NoError(t, err, "Undefined apply should succeed for non-existing path")
	assert.True(t, deepEqual(result.Doc, doc), "Apply should return the original document")

	// Test existing path
	undefinedOp = NewUndefined([]string{"foo"}, false)
	_, err = undefinedOp.Apply(doc)
	assert.Error(t, err, "Undefined apply should fail for existing path")
	assert.Contains(t, err.Error(), "undefined test failed", "Error message should be descriptive")
}

func TestOpUndefined_InterfaceMethods(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"}, false)

	// Test Op() method
	assert.Equal(t, internal.OpUndefinedType, undefinedOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpUndefinedCode, undefinedOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, undefinedOp.Path(), "Path() should return correct path")

	// Test Not() method
	assert.False(t, undefinedOp.Not(), "Not() should return false for default operation")
}

func TestOpUndefined_ToJSON(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"}, false)

	json, err := undefinedOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "undefined", json["op"], "JSON should contain correct op type")
	assert.Equal(t, "/test", json["path"], "JSON should contain correct formatted path")
	assert.Nil(t, json["not"], "JSON should not contain not field when false")
}

func TestOpUndefined_ToJSON_WithNot(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"}, true)

	json, err := undefinedOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, true, json["not"], "JSON should contain not field when true")
}

func TestOpUndefined_ToCompact(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"}, false)

	// Test verbose format
	compact, err := undefinedOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpUndefinedCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
	assert.Equal(t, false, compact[2], "Third element should be not flag")

	// Test non-verbose format
	compact, err = undefinedOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
}

func TestOpUndefined_Validate(t *testing.T) {
	// Test valid operation
	undefinedOp := NewUndefined([]string{"test"}, false)
	err := undefinedOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	undefinedOp = NewUndefined([]string{}, false)
	err = undefinedOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "path cannot be empty", "Error message should mention empty path")
}
