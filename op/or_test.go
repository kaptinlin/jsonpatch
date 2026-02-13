package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOr_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	ok, err := orOp.Test(doc)
	if err != nil {
		t.Fatalf("Or.Test() failed: %v", err)
	}
	if !ok {
		assert.Fail(t, "Or.Test() = false, want true when any operation passes")
	}
}

func TestOr_AllFail(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "qux")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	ok, err := orOp.Test(doc)
	if err != nil {
		t.Fatalf("Or.Test() failed: %v", err)
	}
	if ok {
		assert.Fail(t, "Or.Test() = true, want false when all operations fail")
	}
}

func TestOr_Empty(t *testing.T) {
	t.Parallel()
	orOp := NewOr([]string{}, []any{})

	doc := map[string]any{"foo": "bar"}
	ok, err := orOp.Test(doc)
	if err != nil {
		t.Fatalf("Or.Test() failed: %v", err)
	}
	if ok {
		assert.Fail(t, "Or.Test() = true, want false for empty OR")
	}
}

func TestOr_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	result, err := orOp.Apply(doc)
	if err != nil {
		t.Fatalf("Or.Apply() failed: %v", err)
	}
	if !deepEqual(result.Doc, doc) {
		assert.Fail(t, "Or.Apply() did not return the original document")
	}
}

func TestOr_Apply_Fails(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	test1 := NewTest([]string{"foo"}, "qux")
	test2 := NewTest([]string{"baz"}, 456)

	orOp := NewOr([]string{}, []any{test1, test2})

	_, err := orOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Or.Apply() succeeded, want error when all operations fail")
	}
	if !errors.Is(err, ErrOrTestFailed) {
		assert.Equal(t, ErrOrTestFailed, err, "Or.Apply() error")
	}
}

func TestOr_InterfaceMethods(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	if got := orOp.Op(); got != internal.OpOrType {
		assert.Equal(t, internal.OpOrType, got, "Op()")
	}
	if got := orOp.Code(); got != internal.OpOrCode {
		assert.Equal(t, internal.OpOrCode, got, "Code()")
	}
	if diff := cmp.Diff([]string{"test"}, orOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}

	ops := orOp.Ops()
	if len(ops) != 2 {
		t.Fatalf("len(Ops()) = %d, want 2", len(ops))
	}
	if ops[0] != test1 {
		assert.Fail(t, "Ops()[0] does not match first operation")
	}
	if ops[1] != test2 {
		assert.Fail(t, "Ops()[1] does not match second operation")
	}
}

func TestOr_ToJSON(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	got, err := orOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "or" {
		assert.Equal(t, "or", got.Op, "ToJSON().Op")
	}
	if got.Path != "/test" {
		assert.Equal(t, "/test", got.Path, "ToJSON().Path")
	}
	if got.Apply == nil {
		t.Fatal("ToJSON().Apply = nil, want non-nil")
	}
	if len(got.Apply) != 2 {
		t.Errorf("len(ToJSON().Apply) = %d, want 2", len(got.Apply))
	}
}

func TestOr_ToCompact(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	compact, err := orOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	assert.Equal(t, internal.OpOrCode, compact[0], "ToCompact()[0]")
	assert.Equal(t, []string{"test"}, compact[1])
	if _, ok := compact[2].([]any); !ok {
		t.Errorf("ToCompact()[2] type = %T, want []any", compact[2])
	}
}

func TestOr_Validate(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})
	if err := orOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	// Empty OR is valid, just returns false
	orOp = NewOr([]string{"test"}, []any{})
	if err := orOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for empty operations", err)
	}
}
