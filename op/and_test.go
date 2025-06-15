package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpAnd_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations that should both pass
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	// Create AND operation
	andOp := NewOpAndOperation([]string{}, []interface{}{test1, test2})

	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.True(t, ok, "AND should pass when all operations pass")
}

func TestOpAnd_OneFails(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations, one should pass, one should fail
	test1 := NewOpTestOperation([]string{"foo"}, "bar") // should pass
	test2 := NewOpTestOperation([]string{"baz"}, 456)   // should fail

	// Create AND operation
	andOp := NewOpAndOperation([]string{}, []interface{}{test1, test2})

	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.False(t, ok, "AND should fail when any operation fails")
}

func TestOpAnd_Empty(t *testing.T) {
	// Create AND operation with no sub-operations
	andOp := NewOpAndOperation([]string{}, []interface{}{})

	doc := map[string]interface{}{"foo": "bar"}
	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.False(t, ok, "Empty AND should return false")
}

func TestOpAnd_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations that should both pass
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	// Create AND operation
	andOp := NewOpAndOperation([]string{}, []interface{}{test1, test2})

	result, err := andOp.Apply(doc)
	require.NoError(t, err, "AND apply should succeed when all operations pass")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")
}

func TestOpAnd_Apply_Fails(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations, one should fail
	test1 := NewOpTestOperation([]string{"foo"}, "bar") // should pass
	test2 := NewOpTestOperation([]string{"baz"}, 456)   // should fail

	// Create AND operation
	andOp := NewOpAndOperation([]string{}, []interface{}{test1, test2})

	_, err := andOp.Apply(doc)
	assert.Error(t, err, "AND apply should fail when any operation fails")
	assert.Contains(t, err.Error(), "and test failed", "Error message should be descriptive")
}

func TestOpAnd_InterfaceMethods(t *testing.T) {
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	andOp := NewOpAndOperation([]string{"test"}, []interface{}{test1, test2})

	// Test Op() method
	assert.Equal(t, internal.OpAndType, andOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpAndCode, andOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, andOp.Path(), "Path() should return correct path")

	// Test Ops() method
	ops := andOp.Ops()
	assert.Len(t, ops, 2, "Ops() should return correct number of operations")
	assert.Equal(t, test1, ops[0], "First operation should match")
	assert.Equal(t, test2, ops[1], "Second operation should match")
}

func TestOpAnd_ToJSON(t *testing.T) {
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	andOp := NewOpAndOperation([]string{"test"}, []interface{}{test1, test2})

	json, err := andOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "and", json["op"], "JSON should contain correct op type")
	assert.Equal(t, "/test", json["path"], "JSON should contain correct formatted path")

	// Check apply array
	apply, ok := json["apply"].([]interface{})
	require.True(t, ok, "JSON should contain apply array")
	assert.Len(t, apply, 2, "JSON should contain correct number of operations")
}

func TestOpAnd_ToCompact(t *testing.T) {
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	andOp := NewOpAndOperation([]string{"test"}, []interface{}{test1, test2})

	compact, err := andOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")

	assert.Equal(t, internal.OpAndCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
	assert.IsType(t, []interface{}{}, compact[2], "Third element should be ops array")
}

func TestOpAnd_Validate(t *testing.T) {
	// Test valid operation
	test1 := NewOpTestOperation([]string{"foo"}, "bar")
	test2 := NewOpTestOperation([]string{"baz"}, 123)

	andOp := NewOpAndOperation([]string{"test"}, []interface{}{test1, test2})
	err := andOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty ops)
	andOp = NewOpAndOperation([]string{"test"}, []interface{}{})
	err = andOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "must have at least one operand", "Error message should mention missing operands")
}
