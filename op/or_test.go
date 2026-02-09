package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOr_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.True(t, ok, "OR should pass when any operation passes")
}

func TestOr_AllFail(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "qux")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.False(t, ok, "OR should fail when all operations fail")
}

func TestOr_Empty(t *testing.T) {
	orOp := NewOr([]string{}, []any{})

	doc := map[string]any{"foo": "bar"}
	ok, err := orOp.Test(doc)
	require.NoError(t, err, "OR test should not fail")
	assert.False(t, ok, "Empty OR should return false")
}

func TestOr_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	result, err := orOp.Apply(doc)
	require.NoError(t, err, "OR apply should succeed when any operation passes")
	assert.True(t, deepEqual(result.Doc, doc), "Apply should return the original document")
}

func TestOr_Apply_Fails(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "qux")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	_, err := orOp.Apply(doc)
	assert.Error(t, err, "OR apply should fail when all operations fail")
	assert.ErrorIs(t, err, ErrOrTestFailed)
}

func TestOr_InterfaceMethods(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	assert.Equal(t, internal.OpOrType, orOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpOrCode, orOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"test"}, orOp.Path(), "Path() should return correct path")

	ops := orOp.Ops()
	assert.Len(t, ops, 2, "Ops() should return correct number of operations")
	assert.Equal(t, test1, ops[0], "First operation should match")
	assert.Equal(t, test2, ops[1], "Second operation should match")
}

func TestOr_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	got, err := orOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "or", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", got.Path, "JSON should contain correct formatted path")
	require.NotNil(t, got.Apply, "JSON should contain apply array")
	assert.Len(t, got.Apply, 2, "JSON should contain correct number of operations")
}

func TestOr_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	compact, err := orOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")

	assert.Equal(t, internal.OpOrCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
	assert.IsType(t, []any{}, compact[2], "Third element should be ops array")
}

func TestOr_Validate(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})
	err := orOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Empty OR is valid, just returns false
	orOp = NewOr([]string{"test"}, []any{})
	err = orOp.Validate()
	assert.NoError(t, err, "Empty operations are valid (though they return false)")
}
