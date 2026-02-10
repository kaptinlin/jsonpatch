package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestUndefined_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": 123,
		},
	}

	undefinedOp := NewUndefined([]string{"qux"})
	ok, err := undefinedOp.Test(doc)
	if err != nil {
		t.Fatalf("Undefined.Test(doc, /qux) failed: %v", err)
	}
	if !ok {
		t.Error("Undefined.Test(doc, /qux) = false, want true for non-existing path")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		t.Fatalf("Undefined.Test(doc, /foo) failed: %v", err)
	}
	if ok {
		t.Error("Undefined.Test(doc, /foo) = true, want false for existing path")
	}

	undefinedOp = NewUndefined([]string{"baz", "quux"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		t.Fatalf("Undefined.Test(doc, /baz/quux) failed: %v", err)
	}
	if !ok {
		t.Error("Undefined.Test(doc, /baz/quux) = false, want true for non-existing nested path")
	}
}

func TestUndefined_Not(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	undefinedOp := NewUndefined([]string{"qux"})
	ok, err := undefinedOp.Test(doc)
	if err != nil {
		t.Fatalf("Undefined.Test(doc, /qux) failed: %v", err)
	}
	if !ok {
		t.Error("Undefined.Test(doc, /qux) = false, want true for non-existing path")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		t.Fatalf("Undefined.Test(doc, /foo) failed: %v", err)
	}
	if ok {
		t.Error("Undefined.Test(doc, /foo) = true, want false for existing path")
	}
}

func TestUndefined_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	undefinedOp := NewUndefined([]string{"qux"})
	result, err := undefinedOp.Apply(doc)
	if err != nil {
		t.Fatalf("Undefined.Apply(doc, /qux) failed: %v", err)
	}
	if !deepEqual(result.Doc, doc) {
		t.Error("Undefined.Apply(doc, /qux) did not return the original document")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	_, err = undefinedOp.Apply(doc)
	if err == nil {
		t.Error("Undefined.Apply(doc, /foo) succeeded, want error for existing path")
	}
	if !errors.Is(err, ErrUndefinedTestFailed) {
		t.Errorf("Undefined.Apply(doc, /foo) error = %v, want %v", err, ErrUndefinedTestFailed)
	}
}

func TestUndefined_InterfaceMethods(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	if got := undefinedOp.Op(); got != internal.OpUndefinedType {
		t.Errorf("Op() = %v, want %v", got, internal.OpUndefinedType)
	}
	if got := undefinedOp.Code(); got != internal.OpUndefinedCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpUndefinedCode)
	}
	if diff := cmp.Diff([]string{"test"}, undefinedOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	if undefinedOp.Not() {
		t.Error("Not() = true, want false for default operation")
	}
}

func TestUndefined_ToJSON(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	got, err := undefinedOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "undefined" {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, "undefined")
	}
	if got.Path != "/test" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/test")
	}
}

// TestUndefined_ToJSON_WithNot has been removed since undefined operation
// no longer supports direct negation. Use second-order predicate "not" for negation.

func TestUndefined_ToCompact(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	compact, err := undefinedOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	if len(compact) != 2 {
		t.Fatalf("len(ToCompact()) = %d, want 2", len(compact))
	}
	if compact[0] != internal.OpUndefinedCode {
		t.Errorf("ToCompact()[0] = %v, want %v", compact[0], internal.OpUndefinedCode)
	}
	if diff := cmp.Diff([]string{"test"}, compact[1]); diff != "" {
		t.Errorf("ToCompact()[1] mismatch (-want +got):\n%s", diff)
	}
}

func TestUndefined_Validate(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})
	if err := undefinedOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	undefinedOp = NewUndefined([]string{})
	err := undefinedOp.Validate()
	if err == nil {
		t.Error("Validate() = nil, want error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathEmpty)
	}
}
