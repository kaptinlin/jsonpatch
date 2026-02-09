package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNot_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	require.NoError(t, err, "NOT test should not fail")
	assert.False(t, ok, "NOT should return false when wrapped operation passes")
}

func TestNot_Negation(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "qux")
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	require.NoError(t, err, "NOT test should not fail")
	assert.True(t, ok, "NOT should return true when wrapped operation fails")
}

func TestNot_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "qux")
	notOp := NewNot(testOp)

	result, err := notOp.Apply(doc)
	require.NoError(t, err, "NOT apply should succeed when wrapped operation fails")
	assert.Equal(t, doc, result.Doc, "Apply should return the original document")
}

func TestNot_Apply_Fails(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	_, err := notOp.Apply(doc)
	assert.Error(t, err, "NOT apply should fail when wrapped operation passes")
	assert.ErrorIs(t, err, ErrNotTestFailed)
}

func TestNot_InterfaceMethods(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	assert.Equal(t, internal.OpNotType, notOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpNotCode, notOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"foo"}, notOp.Path(), "Path() should return correct path")

	ops := notOp.Ops()
	assert.Len(t, ops, 1, "Ops() should return correct number of operations")
	assert.Equal(t, testOp, ops[0], "Operation should match")
}

func TestNot_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	got, err := notOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "not", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/foo", got.Path, "JSON should contain correct formatted path")
	assert.NotNil(t, got.Apply, "JSON should contain apply field")
}

func TestNot_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	compact, err := notOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")

	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpNotCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"foo"}, compact[1], "Second element should be path")
	assert.NotNil(t, compact[2], "Third element should be the compact operand")
}

func TestNot_Validate(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	notOp := NewNot(testOp)
	err := notOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	notOp = &NotOperation{BaseOp: NewBaseOp([]string{"test"}), Operations: []any{}}
	err = notOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrNotNoOperands)
}
