package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpTest_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	// Test basic equality
	op := NewTest([]string{"foo"}, "bar")

	ok, err := op.Test(doc)
	require.NoError(t, err, "Test should not fail for valid path")
	assert.True(t, ok, "Test should pass for equal values")

	// Test inequality
	op = NewTest([]string{"foo"}, "qux")
	ok, err = op.Test(doc)
	require.NoError(t, err, "Test should not fail for valid path")
	assert.False(t, ok, "Test should fail for different values")
}

func TestOpTest_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	// Test successful apply
	op := NewTest([]string{"foo"}, "bar")
	result, err := op.Apply(doc)
	require.NoError(t, err, "Apply should succeed for matching values")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")

	// Test failed apply
	op = NewTest([]string{"foo"}, "qux")
	_, err = op.Apply(doc)
	assert.Error(t, err, "Apply should fail for non-matching values")
}

func TestOpTest_ToJSON(t *testing.T) {
	op := NewTest([]string{"foo"}, "bar")

	json, err := op.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "test", json.Op, "JSON should contain correct op type")
	assert.Equal(t, "/foo", json.Path, "JSON should contain correct path")
	assert.Equal(t, "bar", json.Value, "JSON should contain correct value")
}

func TestOpTest_ToCompact(t *testing.T) {
	op := NewTest([]string{"foo"}, "bar")

	// Test compact format
	compact, err := op.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpTestCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.Equal(t, "bar", compact[2], "Third element should be value")
}

func TestOpTest_Validate(t *testing.T) {
	// Test valid operation
	op := NewTest([]string{"foo"}, "bar")
	err := op.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	op = NewTest([]string{}, "bar")
	err = op.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
}

func TestOpTest_InterfaceMethods(t *testing.T) {
	op := NewTest([]string{"foo"}, "bar")

	// Test Op() method
	assert.Equal(t, internal.OpTestType, op.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpTestCode, op.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"foo"}, op.Path(), "Path() should return correct path")
}
