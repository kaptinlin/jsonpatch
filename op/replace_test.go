package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestReplace_Basic(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{
			"nested": "value",
		},
	}

	replaceOp := NewReplace([]string{"foo"}, "new_value")
	result, err := replaceOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	if got := result.Old; got != "bar" {
		assert.Equal(t, "bar", got, "result.Old")
	}
	if got := modifiedDoc["foo"]; got != "new_value" {
		assert.Equal(t, "new_value", got, "modifiedDoc[foo]")
	}
	if got := modifiedDoc["baz"]; got != 123 {
		assert.Equal(t, 123, got, "modifiedDoc[baz]")
	}
}

func TestReplace_Nested(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": map[string]any{
			"bar": "baz",
			"qux": 123,
		},
	}

	replaceOp := NewReplace([]string{"foo", "bar"}, "new_nested_value")
	result, err := replaceOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	foo := modifiedDoc["foo"].(map[string]any)
	if got := result.Old; got != "baz" {
		assert.Equal(t, "baz", got, "result.Old")
	}
	if got := foo["bar"]; got != "new_nested_value" {
		assert.Equal(t, "new_nested_value", got, "foo[bar]")
	}
	if got := foo["qux"]; got != 123 {
		assert.Equal(t, 123, got, "foo[qux]")
	}
}

func TestReplace_Array(t *testing.T) {
	t.Parallel()
	doc := []any{
		"first",
		"second",
		"third",
	}

	replaceOp := NewReplace([]string{"1"}, "new_second")
	result, err := replaceOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedArray := result.Doc.([]any)
	if got := result.Old; got != "second" {
		assert.Equal(t, "second", got, "result.Old")
	}
	assert.Equal(t, "new_second", modifiedArray[1], "modifiedArray[1]")
	assert.Equal(t, "first", modifiedArray[0], "modifiedArray[0]")
	assert.Equal(t, "third", modifiedArray[2], "modifiedArray[2]")
}

func TestReplace_NonExistent(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{"qux"}, "new_value")
	_, err := replaceOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Apply() expected error for non-existent path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		assert.Equal(t, ErrPathNotFound, err, "Apply() error")
	}
}

func TestReplace_EmptyPath(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{}, "new_value")
	result, err := replaceOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}
	assert.Equal(t, "new_value", result.Doc, "result.Doc")
	assert.Equal(t, doc, result.Old)
}

func TestReplace_InterfaceMethods(t *testing.T) {
	t.Parallel()
	replaceOp := NewReplace([]string{"test"}, "value")

	if got := replaceOp.Op(); got != internal.OpReplaceType {
		assert.Equal(t, internal.OpReplaceType, got, "Op()")
	}
	if got := replaceOp.Code(); got != internal.OpReplaceCode {
		assert.Equal(t, internal.OpReplaceCode, got, "Code()")
	}
	assert.Equal(t, []string{"test"}, replaceOp.Path(), "Path()")
	assert.Equal(t, "value", replaceOp.Value, "Value")
}

func TestReplace_ToJSON(t *testing.T) {
	t.Parallel()
	replaceOp := NewReplace([]string{"test"}, "value")

	got, err := replaceOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	assert.Equal(t, "replace", got.Op, "ToJSON().Op")
	assert.Equal(t, "/test", got.Path, "ToJSON().Path")
	assert.Equal(t, "value", got.Value, "ToJSON().Value")
}

func TestReplace_ToCompact(t *testing.T) {
	t.Parallel()
	replaceOp := NewReplace([]string{"test"}, "value")

	compact, err := replaceOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	assert.Equal(t, internal.OpReplaceCode, compact[0], "compact[0]")
	assert.Equal(t, []string{"test"}, compact[1])
	assert.Equal(t, "value", compact[2], "compact[2]")
}

func TestReplace_Validate(t *testing.T) {
	t.Parallel()
	replaceOp := NewReplace([]string{"test"}, "value")
	if err := replaceOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	replaceOp = NewReplace([]string{}, "value")
	err := replaceOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		assert.Equal(t, ErrPathEmpty, err, "Validate() error")
	}
}
