package op

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestNot_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Not.Test() failed: %v", err))
	}
	if ok {
		assert.Fail(t, "Not.Test() = true, want false when wrapped operation passes")
	}
}

func TestNot_Negation(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "qux")
	notOp := NewNot(testOp)

	ok, err := notOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Not.Test() failed: %v", err))
	}
	if !ok {
		assert.Fail(t, "Not.Test() = false, want true when wrapped operation fails")
	}
}

func TestNot_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "qux")
	notOp := NewNot(testOp)

	result, err := notOp.Apply(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Not.Apply() failed: %v", err))
	}
	assert.Equal(t, doc, result.Doc)
}

func TestNot_Apply_Fails(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	_, err := notOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Not.Apply() succeeded, want error when wrapped operation passes")
	}
	if !errors.Is(err, ErrNotTestFailed) {
		assert.Equal(t, ErrNotTestFailed, err, "Not.Apply() error")
	}
}

func TestNot_InterfaceMethods(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	if got := notOp.Op(); got != internal.OpNotType {
		assert.Equal(t, internal.OpNotType, got, "Op()")
	}
	if got := notOp.Code(); got != internal.OpNotCode {
		assert.Equal(t, internal.OpNotCode, got, "Code()")
	}
	assert.Equal(t, []string{"foo"}, notOp.Path(), "Path()")

	ops := notOp.Ops()
	if len(ops) != 1 {
		require.FailNow(t, fmt.Sprintf("len(Ops()) = %d, want 1", len(ops)))
	}
	if ops[0] != testOp {
		assert.Fail(t, "Ops()[0] does not match wrapped operation")
	}
}

func TestNot_ToJSON(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	got, err := notOp.ToJSON()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToJSON() failed: %v", err))
	}
	if got.Op != "not" {
		assert.Equal(t, "not", got.Op, "ToJSON().Op")
	}
	if got.Path != "/foo" {
		assert.Equal(t, "/foo", got.Path, "ToJSON().Path")
	}
	if got.Apply == nil {
		assert.Fail(t, "ToJSON().Apply = nil, want non-nil")
	}
}

func TestNot_ToCompact(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	compact, err := notOp.ToCompact()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToCompact() failed: %v", err))
	}
	if len(compact) != 3 {
		require.FailNow(t, fmt.Sprintf("len(ToCompact()) = %d, want 3", len(compact)))
	}
	assert.Equal(t, internal.OpNotCode, compact[0], "ToCompact()[0]")
	assert.Equal(t, []string{"foo"}, compact[1])
	if compact[2] == nil {
		assert.Fail(t, "ToCompact()[2] = nil, want non-nil compact operand")
	}
}

func TestNot_Validate(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	notOp := NewNot(testOp)
	if err := notOp.Validate(); err != nil {
		assert.Fail(t, fmt.Sprintf("Validate() = %v, want nil for valid operation", err))
	}

	notOp = &NotOperation{BaseOp: NewBaseOp([]string{"test"}), Operations: []any{}}
	err := notOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() = nil, want error for empty operations")
	}
	if !errors.Is(err, ErrNotNoOperands) {
		assert.Equal(t, ErrNotNoOperands, err, "Validate() error")
	}
}

func TestNot_MultipleOperandsContract(t *testing.T) {
	t.Parallel()

	first := NewTest([]string{"name"}, "Grace")
	second := NewTest([]string{"role"}, "owner")
	notOp := NewNotMultiple([]string{"profile"}, []any{first, second})
	assert.True(t, notOp.Not())

	ops := notOp.Ops()
	require.Len(t, ops, 2)
	assert.Same(t, first, ops[0])
	assert.Same(t, second, ops[1])

	matched, err := notOp.Test(map[string]any{"name": "Ada", "role": "admin"})
	assert.NoError(t, err)
	assert.True(t, matched)

	_, err = notOp.Apply(map[string]any{"name": "Grace", "role": "admin"})
	assert.ErrorIs(t, err, ErrNotTestFailed)

	jsonOp, err := notOp.ToJSON()
	assert.NoError(t, err)
	wantJSON := internal.Operation{
		Op:   "not",
		Path: "/profile",
		Apply: []internal.Operation{
			{Op: "test", Path: "/name", Value: "Grace"},
			{Op: "test", Path: "/role", Value: "owner"},
		},
	}
	if diff := cmp.Diff(wantJSON, jsonOp); diff != "" {
		t.Errorf("ToJSON() mismatch (-want +got):\n%s", diff)
	}

	compactOp, err := notOp.ToCompact()
	assert.NoError(t, err)
	wantCompact := internal.CompactOperation{
		internal.OpNotCode,
		[]string{"profile"},
		[]any{
			internal.CompactOperation{internal.OpTestCode, []string{"name"}, "Grace"},
			internal.CompactOperation{internal.OpTestCode, []string{"role"}, "owner"},
		},
	}
	if diff := cmp.Diff(wantCompact, compactOp); diff != "" {
		t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
	}

	invalid := NewNotMultiple([]string{"profile"}, []any{"not a predicate"})
	assert.ErrorIs(t, invalid.Validate(), ErrInvalidPredicateInNot)
	_, err = invalid.Test(map[string]any{})
	assert.ErrorIs(t, err, ErrInvalidPredicateInNot)
	_, err = invalid.ToJSON()
	assert.ErrorIs(t, err, ErrInvalidPredicateInNot)
	_, err = invalid.ToCompact()
	assert.ErrorIs(t, err, ErrInvalidPredicateInNot)
}
