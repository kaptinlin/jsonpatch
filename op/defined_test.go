package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
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
		assert.Fail(t, "Defined.Test(doc, /foo) = false, want true for existing path")
	}

	definedOp = NewDefined([]string{"qux"})
	ok, err = definedOp.Test(doc)
	if err != nil {
		t.Fatalf("Defined.Test(doc, /qux) failed: %v", err)
	}
	if ok {
		assert.Fail(t, "Defined.Test(doc, /qux) = true, want false for non-existing path")
	}

	definedOp = NewDefined([]string{"baz", "qux"})
	ok, err = definedOp.Test(doc)
	if err != nil {
		t.Fatalf("Defined.Test(doc, /baz/qux) failed: %v", err)
	}
	if !ok {
		assert.Fail(t, "Defined.Test(doc, /baz/qux) = false, want true for existing nested path")
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
		assert.Fail(t, "Defined.Apply(doc, /foo) did not return the original document")
	}

	definedOp = NewDefined([]string{"qux"})
	_, err = definedOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Defined.Apply(doc, /qux) succeeded, want error for non-existing path")
	}
	if !errors.Is(err, ErrDefinedTestFailed) {
		assert.Equal(t, ErrDefinedTestFailed, err, "Defined.Apply(doc, /qux) error")
	}
}

func TestDefined_InterfaceMethods(t *testing.T) {
	t.Parallel()
	definedOp := NewDefined([]string{"test"})

	if got := definedOp.Op(); got != internal.OpDefinedType {
		assert.Equal(t, internal.OpDefinedType, got, "Op()")
	}
	if got := definedOp.Code(); got != internal.OpDefinedCode {
		assert.Equal(t, internal.OpDefinedCode, got, "Code()")
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
		assert.Equal(t, "defined", got.Op, "ToJSON().Op")
	}
	if got.Path != "/test" {
		assert.Equal(t, "/test", got.Path, "ToJSON().Path")
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
	assert.Equal(t, internal.OpDefinedCode, compact[0], "ToCompact()[0]")
	assert.Equal(t, []string{"test"}, compact[1])
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
