package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpOr_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations, one should pass, one should fail
	test1 := NewTest([]string{"foo"}, "bar") // should pass
	test2 := NewTest([]string{"baz"}, 456)   // should fail

	// Create OR operation
	orOp := NewOr([]string{}, []interface{}{test1, test2})

	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.True(t, ok, "OR should pass when any operation passes")
}

func TestOpOr_AllFail(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations that should both fail
	test1 := NewTest([]string{"foo"}, "qux") // should fail
	test2 := NewTest([]string{"baz"}, 456)   // should fail

	// Create OR operation
	orOp := NewOr([]string{}, []interface{}{test1, test2})

	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.False(t, ok, "OR should fail when all operations fail")
}

func TestOpOr_Empty(t *testing.T) {
	// Create OR operation with no sub-operations
	orOp := NewOr([]string{}, []interface{}{})

	doc := map[string]interface{}{"foo": "bar"}
	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.False(t, ok, "Empty OR should return false")
}

func TestOpOr_Apply(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations, one should pass
	test1 := NewTest([]string{"foo"}, "bar") // should pass
	test2 := NewTest([]string{"baz"}, 456)   // should fail

	// Create OR operation
	orOp := NewOr([]string{}, []interface{}{test1, test2})

	result, err := orOp.Apply(doc)
	require.NoError(t, err, "OR apply should succeed when any operation passes")
	assert.True(t, deepEqual(result.Doc, doc), "Apply should return the original document")
}

func TestOpOr_Apply_Fails(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	// Create two test operations that should both fail
	test1 := NewTest([]string{"foo"}, "qux") // should fail
	test2 := NewTest([]string{"baz"}, 456)   // should fail

	// Create OR operation
	orOp := NewOr([]string{}, []interface{}{test1, test2})

	_, err := orOp.Apply(doc)
	assert.Error(t, err, "OR apply should fail when all operations fail")
	assert.ErrorIs(t, err, ErrOrTestFailed)
}

func TestOpOr_InterfaceMethods(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []interface{}{test1, test2})

	// Test Op() method
	assert.Equal(t, internal.OpOrType, orOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpOrCode, orOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"test"}, orOp.Path(), "Path() should return correct path")

	// Test Ops() method
	ops := orOp.Ops()
	assert.Len(t, ops, 2, "Ops() should return correct number of operations")
	assert.Equal(t, test1, ops[0], "First operation should match")
	assert.Equal(t, test2, ops[1], "Second operation should match")
}

func TestOpOr_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []interface{}{test1, test2})

	json, err := orOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")
	jsonMap := json

	assert.Equal(t, "or", jsonMap.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", jsonMap.Path, "JSON should contain correct formatted path")

	// Check apply array
	apply := jsonMap.Apply
	require.NotNil(t, apply, "JSON should contain apply array")
	assert.Len(t, apply, 2, "JSON should contain correct number of operations")
}

func TestOpOr_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []interface{}{test1, test2})

	compact, err := orOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	compactArr := compact

	assert.Equal(t, internal.OpOrCode, compactArr[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compactArr[1], "Second element should be path")
	assert.IsType(t, []interface{}{}, compactArr[2], "Third element should be ops array")
}

func TestOpOr_Validate(t *testing.T) {
	// Test valid operation
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []interface{}{test1, test2})
	err := orOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test valid empty operations (empty OR is valid, just returns false)
	orOp = NewOr([]string{"test"}, []interface{}{})
	err = orOp.Validate()
	assert.NoError(t, err, "Empty operations are valid (though they return false)")
}
