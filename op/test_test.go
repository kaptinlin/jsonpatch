package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTest_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	testOp := NewTest([]string{"foo"}, "bar")

	ok, err := testOp.Test(doc)
	require.NoError(t, err, "Test should not fail for valid path")
	assert.True(t, ok, "Test should pass for equal values")

	testOp = NewTest([]string{"foo"}, "qux")
	ok, err = testOp.Test(doc)
	require.NoError(t, err, "Test should not fail for valid path")
	assert.False(t, ok, "Test should fail for different values")
}

func TestTest_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	result, err := testOp.Apply(doc)
	require.NoError(t, err, "Apply should succeed for matching values")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")

	testOp = NewTest([]string{"foo"}, "qux")
	_, err = testOp.Apply(doc)
	assert.Error(t, err, "Apply should fail for non-matching values")
}

func TestTest_ToJSON(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	got, err := testOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "test", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/foo", got.Path, "JSON should contain correct path")
	assert.Equal(t, "bar", got.Value, "JSON should contain correct value")
}

func TestTest_ToCompact(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	compact, err := testOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpTestCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.Equal(t, "bar", compact[2], "Third element should be value")
}

func TestTest_Validate(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")
	err := testOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	testOp = NewTest([]string{}, "bar")
	err = testOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
}

func TestTest_InterfaceMethods(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	assert.Equal(t, internal.OpTestType, testOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpTestCode, testOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"foo"}, testOp.Path(), "Path() should return correct path")
}
