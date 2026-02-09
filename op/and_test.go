package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
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
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if !ok {
		t.Error("And.Test() = false, want true when all operations pass")
	}
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
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if ok {
		t.Error("And.Test() = true, want false when any operation fails")
	}
}

func TestAnd_Empty(t *testing.T) {
	andOp := NewAnd([]string{}, []any{})

	doc := map[string]any{"foo": "bar"}
	ok, err := andOp.Test(doc)
	if err != nil {
		t.Fatalf("And.Test() failed: %v", err)
	}
	if !ok {
		t.Error("And.Test() = false, want true for empty AND (vacuous truth)")
	}
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
	if err != nil {
		t.Fatalf("And.Apply() failed: %v", err)
	}
	if diff := cmp.Diff(doc, result.Doc); diff != "" {
		t.Errorf("And.Apply() result.Doc mismatch (-want +got):\n%s", diff)
	}
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
	if err == nil {
		t.Error("And.Apply() succeeded, want error when any operation fails")
	}
	if !errors.Is(err, ErrAndTestFailed) {
		t.Errorf("And.Apply() error = %v, want %v", err, ErrAndTestFailed)
	}
}

func TestAnd_InterfaceMethods(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	if got := andOp.Op(); got != internal.OpAndType {
		t.Errorf("Op() = %v, want %v", got, internal.OpAndType)
	}
	if got := andOp.Code(); got != internal.OpAndCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpAndCode)
	}
	if diff := cmp.Diff([]string{"test"}, andOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}

	ops := andOp.Ops()
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

func TestAnd_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	got, err := andOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "and" {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, "and")
	}
	if got.Path != "/test" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/test")
	}
	if len(got.Apply) != 2 {
		t.Errorf("len(ToJSON().Apply) = %d, want 2", len(got.Apply))
	}
}

func TestAnd_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})

	compact, err := andOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	if compact[0] != internal.OpAndCode {
		t.Errorf("ToCompact()[0] = %v, want %v", compact[0], internal.OpAndCode)
	}
	if diff := cmp.Diff([]string{"test"}, compact[1]); diff != "" {
		t.Errorf("ToCompact()[1] mismatch (-want +got):\n%s", diff)
	}
	if _, ok := compact[2].([]any); !ok {
		t.Errorf("ToCompact()[2] type = %T, want []any", compact[2])
	}
}

func TestAnd_Validate(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	test2 := NewTest([]string{"baz"}, 123)

	andOp := NewAnd([]string{"test"}, []any{test1, test2})
	if err := andOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	// Vacuous truth allows empty operations
	andOp = NewAnd([]string{"test"}, []any{})
	if err := andOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for empty operations (vacuous truth)", err)
	}
}
