package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestDefined_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": map[string]any{
			"qux": 123,
		},
	}

	definedOp := NewDefined([]string{"foo"})
	ok, err := definedOp.Test(doc)
	if err != nil {
		t.Fatalf("Defined.Test(doc, /foo) failed: %v", err)
	}
	if !ok {
		t.Error("Defined.Test(doc, /foo) = false, want true for existing path")
	}

	definedOp = NewDefined([]string{"qux"})
	ok, err = definedOp.Test(doc)
	if err != nil {
		t.Fatalf("Defined.Test(doc, /qux) failed: %v", err)
	}
	if ok {
		t.Error("Defined.Test(doc, /qux) = true, want false for non-existing path")
	}

	definedOp = NewDefined([]string{"baz", "qux"})
	ok, err = definedOp.Test(doc)
	if err != nil {
		t.Fatalf("Defined.Test(doc, /baz/qux) failed: %v", err)
	}
	if !ok {
		t.Error("Defined.Test(doc, /baz/qux) = false, want true for existing nested path")
	}
}

func TestDefined_Apply(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	definedOp := NewDefined([]string{"foo"})
	result, err := definedOp.Apply(doc)
	if err != nil {
		t.Fatalf("Defined.Apply(doc, /foo) failed: %v", err)
	}
	if !deepEqual(result.Doc, doc) {
		t.Error("Defined.Apply(doc, /foo) did not return the original document")
	}

	definedOp = NewDefined([]string{"qux"})
	_, err = definedOp.Apply(doc)
	if err == nil {
		t.Error("Defined.Apply(doc, /qux) succeeded, want error for non-existing path")
	}
	if !errors.Is(err, ErrDefinedTestFailed) {
		t.Errorf("Defined.Apply(doc, /qux) error = %v, want %v", err, ErrDefinedTestFailed)
	}
}

func TestDefined_InterfaceMethods(t *testing.T) {
	t.Parallel()
	definedOp := NewDefined([]string{"test"})

	if got := definedOp.Op(); got != internal.OpDefinedType {
		t.Errorf("Op() = %v, want %v", got, internal.OpDefinedType)
	}
	if got := definedOp.Code(); got != internal.OpDefinedCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpDefinedCode)
	}
	if diff := cmp.Diff([]string{"test"}, definedOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
}

func TestDefined_ToJSON(t *testing.T) {
	t.Parallel()
	definedOp := NewDefined([]string{"test"})

	got, err := definedOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() failed: %v", err)
	}
	if got.Op != "defined" {
		t.Errorf("ToJSON().Op = %q, want %q", got.Op, "defined")
	}
	if got.Path != "/test" {
		t.Errorf("ToJSON().Path = %q, want %q", got.Path, "/test")
	}
}

func TestDefined_ToCompact(t *testing.T) {
	t.Parallel()
	definedOp := NewDefined([]string{"test"})

	compact, err := definedOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() failed: %v", err)
	}
	if len(compact) != 2 {
		t.Fatalf("len(ToCompact()) = %d, want 2", len(compact))
	}
	if compact[0] != internal.OpDefinedCode {
		t.Errorf("ToCompact()[0] = %v, want %v", compact[0], internal.OpDefinedCode)
	}
	if diff := cmp.Diff([]string{"test"}, compact[1]); diff != "" {
		t.Errorf("ToCompact()[1] mismatch (-want +got):\n%s", diff)
	}
}

func TestDefined_Validate(t *testing.T) {
	t.Parallel()
	definedOp := NewDefined([]string{"test"})
	if err := definedOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for valid operation", err)
	}

	definedOp = NewDefined([]string{})
	if err := definedOp.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil for empty path (root)", err)
	}
}
