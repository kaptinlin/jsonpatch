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
	undefinedOp := NewUndefined([]string{"qux"})
	ok, err := undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "Undefined should return true for non-existing path")

	// Test existing path
	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.False(t, ok, "Undefined should return false for existing path")

	// Test nested non-existing path
	undefinedOp = NewUndefined([]string{"baz", "quux"})
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "Undefined should return true for non-existing nested path")
}

func TestOpUndefined_Not(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test non-existing path (standard undefined behavior)
	undefinedOp := NewUndefined([]string{"qux"})
	ok, err := undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.True(t, ok, "undefined should return true for non-existing path")

	// Test existing path (standard undefined behavior)
	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	require.NoError(t, err, "Undefined test should not fail")
	assert.False(t, ok, "undefined should return false for existing path")
}

func TestOpUndefined_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test non-existing path
	undefinedOp := NewUndefined([]string{"qux"})
	result, err := undefinedOp.Apply(doc)
	require.NoError(t, err, "Undefined apply should succeed for non-existing path")
	assert.True(t, deepEqual(result.Doc, doc), "Apply should return the original document")

	// Test existing path
	undefinedOp = NewUndefined([]string{"foo"})
	_, err = undefinedOp.Apply(doc)
	assert.Error(t, err, "Undefined apply should fail for existing path")
	assert.ErrorIs(t, err, ErrUndefinedTestFailed)
}

func TestOpUndefined_InterfaceMethods(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"})

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
	undefinedOp := NewUndefined([]string{"test"})

	json, err := undefinedOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "undefined", json.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", json.Path, "JSON should contain correct formatted path")
	// Not field should not be present since undefined operation no longer supports direct negation
}

// TestOpUndefined_ToJSON_WithNot has been removed since undefined operation
// no longer supports direct negation. Use second-order predicate "not" for negation.

func TestOpUndefined_ToCompact(t *testing.T) {
	undefinedOp := NewUndefined([]string{"test"})

	// Test compact format (no longer includes 'not' parameter)
	compact, err := undefinedOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 2, "Compact format should have 2 elements")
	assert.Equal(t, internal.OpUndefinedCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
}

func TestOpUndefined_Validate(t *testing.T) {
	// Test valid operation
	undefinedOp := NewUndefined([]string{"test"})
	err := undefinedOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	undefinedOp = NewUndefined([]string{})
	err = undefinedOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrPathEmpty)
}
