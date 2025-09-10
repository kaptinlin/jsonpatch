package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpNot_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Create a test operation that should pass
	testOp := NewTest([]string{"foo"}, "bar") // should pass

	// Create NOT operation
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	require.NoError(t, err, "NOT test should not fail")
	assert.False(t, ok, "NOT should return false when wrapped operation passes")
}

func TestOpNot_Negation(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Create a test operation that should fail
	testOp := NewTest([]string{"foo"}, "qux") // should fail

	// Create NOT operation
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	require.NoError(t, err, "NOT test should not fail")
	assert.True(t, ok, "NOT should return true when wrapped operation fails")
}

func TestOpNot_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Create a test operation that should fail
	testOp := NewTest([]string{"foo"}, "qux") // should fail

	// Create NOT operation
	notOp := NewNot(testOp)

	result, err := notOp.Apply(doc)
	require.NoError(t, err, "NOT apply should succeed when wrapped operation fails")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")
}

func TestOpNot_Apply_Fails(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Create a test operation that should pass
	testOp := NewTest([]string{"foo"}, "bar") // should pass

	// Create NOT operation
	notOp := NewNot(testOp)

	_, err := notOp.Apply(doc)
	assert.Error(t, err, "NOT apply should fail when wrapped operation passes")
	assert.Contains(t, err.Error(), "not test failed", "Error message should be descriptive")
}

func TestOpNot_InterfaceMethods(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	notOp := NewNot(testOp)

	// Test Op() method
	assert.Equal(t, internal.OpNotType, notOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpNotCode, notOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"foo"}, notOp.Path(), "Path() should return correct path")

	// Test Ops() method
	ops := notOp.Ops()
	assert.Len(t, ops, 1, "Ops() should return correct number of operations")
	assert.Equal(t, testOp, ops[0], "Operation should match")
}

func TestOpNot_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	json, err := notOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "not", json["op"], "JSON should contain correct op type")
	assert.Equal(t, "/foo", json["path"], "JSON should contain correct formatted path")
	assert.NotNil(t, json["apply"], "JSON should contain apply field")
}

func TestOpNot_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	compact, err := notOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")

	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpNotCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.NotNil(t, compact[2], "Third element should be the compact operand")
}

func TestOpNot_Validate(t *testing.T) {
	// Test valid operation
	testOp := NewTest([]string{"foo"}, "bar")

	notOp := NewNot(testOp)
	err := notOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty operations)
	notOp = &NotOperation{BaseOp: NewBaseOp([]string{"test"}), Operations: []interface{}{}}
	err = notOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "must have at least one operand", "Error message should mention missing operand")
}
