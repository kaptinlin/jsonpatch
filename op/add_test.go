package op

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestAdd_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	addOp := NewAdd([]string{"baz"}, "qux")

	result, err := addOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	if got := doc["foo"]; got != "bar" {
		assert.Equal(t, "bar", got, "doc[foo]")
	}
	if got := doc["baz"]; got != "qux" {
		assert.Equal(t, "qux", got, "doc[baz]")
	}

	resultDoc := result.Doc.(map[string]any)
	if got := resultDoc["foo"]; got != "bar" {
		assert.Equal(t, "bar", got, "resultDoc[foo]")
	}
	if got := resultDoc["baz"]; got != "qux" {
		assert.Equal(t, "qux", got, "resultDoc[baz]")
	}
	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}

	// Mutate behavior: result should point to the same underlying object
	if reflect.ValueOf(doc).Pointer() != reflect.ValueOf(resultDoc).Pointer() {
		assert.Fail(t, "Result should point to the same document object")
	}
}

func TestAdd_ReplaceExisting(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	addOp := NewAdd([]string{"foo"}, "new_value")

	result, err := addOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	resultDoc := result.Doc.(map[string]any)
	if got := resultDoc["foo"]; got != "new_value" {
		assert.Equal(t, "new_value", got, "resultDoc[foo]")
	}
	if got := result.Old; got != "bar" {
		assert.Equal(t, "bar", got, "result.Old")
	}
}

func TestAdd_InterfaceMethods(t *testing.T) {
	t.Parallel()
	addOp := NewAdd([]string{"foo"}, "bar")

	if got := addOp.Op(); got != internal.OpAddType {
		assert.Equal(t, internal.OpAddType, got, "Op()")
	}
	if got := addOp.Code(); got != internal.OpAddCode {
		assert.Equal(t, internal.OpAddCode, got, "Code()")
	}
	if diff := cmp.Diff([]string{"foo"}, addOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
}

func TestAdd_ToJSON(t *testing.T) {
	t.Parallel()
	addOp := NewAdd([]string{"foo", "bar"}, "baz")

	got, err := addOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	assert.Equal(t, "add", got.Op, "ToJSON().Op")
	assert.Equal(t, "/foo/bar", got.Path, "ToJSON().Path")
	assert.Equal(t, "baz", got.Value, "ToJSON().Value")
}

func TestAdd_ToCompact(t *testing.T) {
	t.Parallel()
	addOp := NewAdd([]string{"foo"}, "bar")

	compact, err := addOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	assert.Equal(t, internal.OpAddCode, compact[0], "compact[0]")
	assert.Equal(t, []string{"foo"}, compact[1])
	assert.Equal(t, "bar", compact[2], "compact[2]")
}

func TestAdd_Validate(t *testing.T) {
	t.Parallel()
	addOp := NewAdd([]string{"foo"}, "bar")
	if err := addOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	addOp = NewAdd([]string{}, "bar")
	err := addOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		assert.Equal(t, ErrPathEmpty, err, "Validate() error")
	}
}

func TestAdd_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"foo", "bar"}
	value := "baz"

	addOp := NewAdd(path, value)

	if diff := cmp.Diff(path, addOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	// Verify value indirectly through ToJSON since the value field is unexported
	got, err := addOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}
	assert.Equal(t, value, got.Value, "ToJSON().Value")
}
