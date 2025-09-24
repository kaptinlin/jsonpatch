package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpDefined_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": map[string]interface{}{
			"qux": 123,
		},
	}

	// Test existing path
	definedOp := NewDefined([]string{"foo"})
	ok, err := definedOp.Test(doc)
	require.NoError(t, err, "Defined test should not fail")
	assert.True(t, ok, "Defined should return true for existing path")

	// Test non-existing path
	definedOp = NewDefined([]string{"qux"})
	ok, err = definedOp.Test(doc)
	require.NoError(t, err, "Defined test should not fail")
	assert.False(t, ok, "Defined should return false for non-existing path")

	// Test nested path
	definedOp = NewDefined([]string{"baz", "qux"})
	ok, err = definedOp.Test(doc)
	require.NoError(t, err, "Defined test should not fail")
	assert.True(t, ok, "Defined should return true for existing nested path")
}

func TestOpDefined_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test existing path
	definedOp := NewDefined([]string{"foo"})
	result, err := definedOp.Apply(doc)
	require.NoError(t, err, "Defined apply should succeed for existing path")
	assert.True(t, deepEqual(result.Doc, doc), "Apply should return the original document")

	// Test non-existing path
	definedOp = NewDefined([]string{"qux"})
	_, err = definedOp.Apply(doc)
	assert.Error(t, err, "Defined apply should fail for non-existing path")
	assert.Contains(t, err.Error(), "defined test failed", "Error message should be descriptive")
}

func TestOpDefined_InterfaceMethods(t *testing.T) {
	definedOp := NewDefined([]string{"test"})

	// Test Op() method
	assert.Equal(t, internal.OpDefinedType, definedOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpDefinedCode, definedOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, definedOp.Path(), "Path() should return correct path")
}

func TestOpDefined_ToJSON(t *testing.T) {
	definedOp := NewDefined([]string{"test"})

	operation, err := definedOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "defined", operation.Op, "Operation should contain correct op type")
	assert.Equal(t, "/test", operation.Path, "Operation should contain correct formatted path")
}

func TestOpDefined_ToCompact(t *testing.T) {
	definedOp := NewDefined([]string{"test"})

	// Test verbose format
	compact, err := definedOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 2, "Compact format should have 2 elements")
	assert.Equal(t, internal.OpDefinedCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")

	// Test non-verbose format
	compact, err = definedOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 2, "Compact format should have 2 elements")
}

func TestOpDefined_Validate(t *testing.T) {
	// Test valid operation
	definedOp := NewDefined([]string{"test"})
	err := definedOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test valid operation with empty path (root path is valid for defined)
	definedOp = NewDefined([]string{})
	err = definedOp.Validate()
	assert.NoError(t, err, "Empty path (root) should be valid for defined operation")
}
