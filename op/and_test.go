package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAnd_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{}, []any{test1, test2})

	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.True(t, ok, "AND should pass when all operations pass")
}

func TestAnd_OneFails(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	andOp := NewAnd([]string{}, []any{test1, test2})

	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.False(t, ok, "AND should fail when any operation fails")
}

func TestAnd_Empty(t *testing.T) {
	andOp := NewAnd([]string{}, []any{})

	doc := map[string]any{"foo": "bar"}
	ok, err := andOp.Test(doc)
	require.NoError(t, err, "AND test should not fail")
	assert.True(t, ok, "Empty AND should return true (vacuous truth)")
}

func TestAnd_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{}, []any{test1, test2})

	result, err := andOp.Apply(doc)
	require.NoError(t, err, "AND apply should succeed when all operations pass")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")
}

func TestAnd_Apply_Fails(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	andOp := NewAnd([]string{}, []any{test1, test2})

	_, err := andOp.Apply(doc)
	assert.Error(t, err, "AND apply should fail when any operation fails")
	assert.ErrorIs(t, err, ErrAndTestFailed)
}

func TestAnd_InterfaceMethods(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	assert.Equal(t, internal.OpAndType, andOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpAndCode, andOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"test"}, andOp.Path(), "Path() should return correct path")

	ops := andOp.Ops()
	assert.Len(t, ops, 2, "Ops() should return correct number of operations")
	assert.Equal(t, test1, ops[0], "First operation should match")
	assert.Equal(t, test2, ops[1], "Second operation should match")
}

func TestAnd_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	got, err := andOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "and", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/test", got.Path, "JSON should contain correct formatted path")
	assert.Len(t, got.Apply, 2, "JSON should contain correct number of operations")
}

func TestAnd_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	compact, err := andOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")

	assert.Equal(t, internal.OpAndCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"test"}, compact[1], "Second element should be path")
	assert.IsType(t, []any{}, compact[2], "Third element should be ops array")
}

func TestAnd_Validate(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})
	err := andOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Vacuous truth allows empty operations
	andOp = NewAnd([]string{"test"}, []any{})
	err = andOp.Validate()
	assert.NoError(t, err, "Empty operations are valid (vacuous truth)")
}
