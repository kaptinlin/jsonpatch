package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestTest_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
	}

	testOp := NewTest([]string{"foo"}, "bar")

	ok, err := testOp.Test(doc)
	if err != nil {
		t.Fatalf("Test() unexpected error: %v", err)
	}
	if !ok {
		t.Error("Test() = false, want true for equal values")
	}

	testOp = NewTest([]string{"foo"}, "qux")
	ok, err = testOp.Test(doc)
	if err != nil {
		t.Fatalf("Test() unexpected error: %v", err)
	}
	if ok {
		t.Error("Test() = true, want false for different values")
	}
}

func TestTest_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	testOp := NewTest([]string{"foo"}, "bar")
	result, err := testOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}
	if diff := cmp.Diff(doc, result.Doc); diff != "" {
		t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
	}

	testOp = NewTest([]string{"foo"}, "qux")
	_, err = testOp.Apply(doc)
	if err == nil {
		t.Error("Apply() expected error for non-matching values")
	}
}

func TestTest_ToJSON(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	got, err := testOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	if got.Op != "test" {
		t.Errorf("ToJSON().Op = %v, want %v", got.Op, "test")
	}
	if got.Path != "/foo" {
		t.Errorf("ToJSON().Path = %v, want %v", got.Path, "/foo")
	}
	if got.Value != "bar" {
		t.Errorf("ToJSON().Value = %v, want %v", got.Value, "bar")
	}
}

func TestTest_ToCompact(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	compact, err := testOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	if compact[0] != internal.OpTestCode {
		t.Errorf("compact[0] = %v, want %v", compact[0], internal.OpTestCode)
	}
	if diff := cmp.Diff([]string{"foo"}, compact[1]); diff != "" {
		t.Errorf("compact[1] mismatch (-want +got):\n%s", diff)
	}
	if compact[2] != "bar" {
		t.Errorf("compact[2] = %v, want %v", compact[2], "bar")
	}
}

func TestTest_Validate(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")
	if err := testOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	testOp = NewTest([]string{}, "bar")
	if err := testOp.Validate(); err == nil {
		t.Error("Validate() expected error for empty path")
	}
}

func TestTest_InterfaceMethods(t *testing.T) {
	t.Parallel()
	testOp := NewTest([]string{"foo"}, "bar")

	if got := testOp.Op(); got != internal.OpTestType {
		t.Errorf("Op() = %v, want %v", got, internal.OpTestType)
	}
	if got := testOp.Code(); got != internal.OpTestCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpTestCode)
	}
	if diff := cmp.Diff([]string{"foo"}, testOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
}
