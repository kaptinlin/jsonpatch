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
		require.FailNow(t, fmt.Sprintf("Undefined.Test(doc, /qux) failed: %v", err))
	}
	if !ok {
		assert.Fail(t, "Undefined.Test(doc, /qux) = false, want true for non-existing path")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Undefined.Test(doc, /foo) failed: %v", err))
	}
	if ok {
		assert.Fail(t, "Undefined.Test(doc, /foo) = true, want false for existing path")
	}

	undefinedOp = NewUndefined([]string{"baz", "quux"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Undefined.Test(doc, /baz/quux) failed: %v", err))
	}
	if !ok {
		assert.Fail(t, "Undefined.Test(doc, /baz/quux) = false, want true for non-existing nested path")
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
		require.FailNow(t, fmt.Sprintf("Undefined.Test(doc, /qux) failed: %v", err))
	}
	if !ok {
		assert.Fail(t, "Undefined.Test(doc, /qux) = false, want true for non-existing path")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	ok, err = undefinedOp.Test(doc)
	if err != nil {
		require.FailNow(t, fmt.Sprintf("Undefined.Test(doc, /foo) failed: %v", err))
	}
	if ok {
		assert.Fail(t, "Undefined.Test(doc, /foo) = true, want false for existing path")
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
		require.FailNow(t, fmt.Sprintf("Undefined.Apply(doc, /qux) failed: %v", err))
	}
	if !deepEqual(result.Doc, doc) {
		assert.Fail(t, "Undefined.Apply(doc, /qux) did not return the original document")
	}

	undefinedOp = NewUndefined([]string{"foo"})
	_, err = undefinedOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Undefined.Apply(doc, /foo) succeeded, want error for existing path")
	}
	if !errors.Is(err, ErrUndefinedTestFailed) {
		assert.Equal(t, ErrUndefinedTestFailed, err, "Undefined.Apply(doc, /foo) error")
	}
}

func TestUndefined_InterfaceMethods(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	if got := undefinedOp.Op(); got != internal.OpUndefinedType {
		assert.Equal(t, internal.OpUndefinedType, got, "Op()")
	}
	if got := undefinedOp.Code(); got != internal.OpUndefinedCode {
		assert.Equal(t, internal.OpUndefinedCode, got, "Code()")
	}
	assert.Empty(t, cmp.Diff([]string{"test"}, undefinedOp.Path()), "Path() mismatch")
	if undefinedOp.Not() {
		assert.Fail(t, "Not() = true, want false for default operation")
	}
}

func TestUndefined_ToJSON(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	got, err := undefinedOp.ToJSON()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToJSON() failed: %v", err))
	}
	if got.Op != "undefined" {
		assert.Equal(t, "undefined", got.Op, "ToJSON().Op")
	}
	if got.Path != "/test" {
		assert.Equal(t, "/test", got.Path, "ToJSON().Path")
	}
}

// TestUndefined_ToJSON_WithNot has been removed since undefined operation
// no longer supports direct negation. Use second-order predicate "not" for negation.

func TestUndefined_ToCompact(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})

	compact, err := undefinedOp.ToCompact()
	if err != nil {
		require.FailNow(t, fmt.Sprintf("ToCompact() failed: %v", err))
	}
	if len(compact) != 2 {
		require.FailNow(t, fmt.Sprintf("len(ToCompact()) = %d, want 2", len(compact)))
	}
	assert.Equal(t, internal.OpUndefinedCode, compact[0], "ToCompact()[0]")
	assert.Equal(t, []string{"test"}, compact[1])
}

func TestUndefined_Validate(t *testing.T) {
	t.Parallel()
	undefinedOp := NewUndefined([]string{"test"})
	if err := undefinedOp.Validate(); err != nil {
		assert.Fail(t, fmt.Sprintf("Validate() = %v, want nil for valid operation", err))
	}

	// Empty path is valid per RFC 6902 and json-joy
	undefinedOp = NewUndefined([]string{})
	if err := undefinedOp.Validate(); err != nil {
		assert.Fail(t, fmt.Sprintf("Validate() unexpected error for empty path: %v", err))
	}
}
