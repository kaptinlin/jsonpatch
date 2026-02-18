package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
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
		t.Fatalf("Not.Test() failed: %v", err)
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
		t.Fatalf("Not.Test() failed: %v", err)
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
		t.Fatalf("Not.Apply() failed: %v", err)
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
		t.Fatalf("len(Ops()) = %d, want 1", len(ops))
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
		t.Fatalf("ToJSON() failed: %v", err)
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
		t.Fatalf("ToCompact() failed: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want 3", len(compact))
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
		t.Errorf("Validate() = %v, want nil for valid operation", err)
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
