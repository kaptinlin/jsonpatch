package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestNot_Basic(t *testing.T) {
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
		t.Error("Not.Test() = true, want false when wrapped operation passes")
	}
}

func TestNot_Negation(t *testing.T) {
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
		t.Error("Not.Test() = false, want true when wrapped operation fails")
	}
}

func TestNot_Apply(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "qux")
	notOp := NewNot(testOp)

	result, err := notOp.Apply(doc)
	if err != nil {
		t.Fatalf("Not.Apply() failed: %v", err)
	}
	if diff := cmp.Diff(doc, result.Doc); diff != "" {
		t.Errorf("Not.Apply() result.Doc mismatch (-want +got):\n%s", diff)
	}
}

func TestNot_Apply_Fails(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	_, err := notOp.Apply(doc)
	if err == nil {
		t.Error("Not.Apply() succeeded, want error when wrapped operation passes")
	}
	if !errors.Is(err, ErrNotTestFailed) {
		t.Errorf("Not.Apply() error = %v, want %v", err, ErrNotTestFailed)
	}
}

func TestNot_InterfaceMethods(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(testOp)

	if got := notOp.Op(); got != internal.OpNotType {
		t.Errorf("Op() = %v, want %v", got, internal.OpNotType)
	}
	if got := notOp.Code(); got != internal.OpNotCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpNotCode)
	}
	if diff := cmp.Diff([]string{"foo"}, notOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}

	ops := notOp.Ops()
	if len(ops) != 1 {
		t.Fatalf("len(Ops()) = %d, want 1", len(ops))
	}
	if ops[0] != testOp {
		t.Error("Ops()[0] does not match wrapped operation")
	}
}

func TestNot_ToJSON(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	got, err := notOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "not" {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, "not")
	}
	if got.Path != "/foo" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/foo")
	}
	if got.Apply == nil {
		t.Error("ToJSON().Apply = nil, want non-nil")
	}
}

func TestNot_ToCompact(t *testing.T) {
	test1 := NewTest([]string{"foo"}, "bar")
	notOp := NewNot(test1)

	compact, err := notOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want 3", len(compact))
	}
	if compact[0] != internal.OpNotCode {
		t.Errorf("ToCompact()[0] = %v, want %v", compact[0], internal.OpNotCode)
	}
	if diff := cmp.Diff([]string{"foo"}, compact[1]); diff != "" {
		t.Errorf("ToCompact()[1] mismatch (-want +got):\n%s", diff)
	}
	if compact[2] == nil {
		t.Error("ToCompact()[2] = nil, want non-nil compact operand")
	}
}

func TestNot_Validate(t *testing.T) {
	testOp := NewTest([]string{"foo"}, "bar")

	notOp := NewNot(testOp)
	if err := notOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	notOp = &NotOperation{BaseOp: NewBaseOp([]string{"test"}), Operations: []any{}}
	err := notOp.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error for empty operations")
	}
	if !errors.Is(err, ErrNotNoOperands) {
		t.Errorf("Validate() error = %v, want %v", err, ErrNotNoOperands)
	}
}
