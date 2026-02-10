package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
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
		t.Error("Or.Test() = false, want true when any operation passes")
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
		t.Error("Or.Test() = true, want false when all operations fail")
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
		t.Error("Or.Test() = true, want false for empty OR")
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
		t.Error("Or.Apply() did not return the original document")
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
		t.Error("Or.Apply() succeeded, want error when all operations fail")
	}
	if !errors.Is(err, ErrOrTestFailed) {
		t.Errorf("Or.Apply() error = %v, want %v", err, ErrOrTestFailed)
	}
}

func TestOr_InterfaceMethods(t *testing.T) {
	t.Parallel()
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	orOp := NewOr([]string{"test"}, []any{test1, test2})

	if got := orOp.Op(); got != internal.OpOrType {
		t.Errorf("Op() = %v, want %v", got, internal.OpOrType)
	}
	if got := orOp.Code(); got != internal.OpOrCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpOrCode)
	}
	if diff := cmp.Diff([]string{"test"}, orOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}

	ops := orOp.Ops()
	if len(ops) != 2 {
		t.Fatalf("len(Ops()) = %d, want 2", len(ops))
	}
	if ops[0] != test1 {
		t.Error("Ops()[0] does not match first operation")
	}
	if ops[1] != test2 {
		t.Error("Ops()[1] does not match second operation")
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
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, "or")
	}
	if got.Path != "/test" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/test")
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
	if compact[0] != internal.OpOrCode {
		t.Errorf("ToCompact()[0] = %v, want %v", compact[0], internal.OpOrCode)
	}
	if diff := cmp.Diff([]string{"test"}, compact[1]); diff != "" {
		t.Errorf("ToCompact()[1] mismatch (-want +got):\n%s", diff)
	}
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
