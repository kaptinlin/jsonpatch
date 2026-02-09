package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestReplace_Basic(t *testing.T) {
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
		t.Errorf("result.Old = %v, want %v", got, "bar")
	}
	if got := modifiedDoc["foo"]; got != "new_value" {
		t.Errorf("modifiedDoc[foo] = %v, want %v", got, "new_value")
	}
	if got := modifiedDoc["baz"]; got != 123 {
		t.Errorf("modifiedDoc[baz] = %v, want %v", got, 123)
	}
}

func TestReplace_Nested(t *testing.T) {
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
		t.Errorf("result.Old = %v, want %v", got, "baz")
	}
	if got := foo["bar"]; got != "new_nested_value" {
		t.Errorf("foo[bar] = %v, want %v", got, "new_nested_value")
	}
	if got := foo["qux"]; got != 123 {
		t.Errorf("foo[qux] = %v, want %v", got, 123)
	}
}

func TestReplace_Array(t *testing.T) {
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
		t.Errorf("result.Old = %v, want %v", got, "second")
	}
	if modifiedArray[1] != "new_second" {
		t.Errorf("modifiedArray[1] = %v, want %v", modifiedArray[1], "new_second")
	}
	if modifiedArray[0] != "first" {
		t.Errorf("modifiedArray[0] = %v, want %v", modifiedArray[0], "first")
	}
	if modifiedArray[2] != "third" {
		t.Errorf("modifiedArray[2] = %v, want %v", modifiedArray[2], "third")
	}
}

func TestReplace_NonExistent(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{"qux"}, "new_value")
	_, err := replaceOp.Apply(doc)
	if err == nil {
		t.Error("Apply() expected error for non-existent path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		t.Errorf("Apply() error = %v, want %v", err, ErrPathNotFound)
	}
}

func TestReplace_EmptyPath(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	replaceOp := NewReplace([]string{}, "new_value")
	result, err := replaceOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}
	if result.Doc != "new_value" {
		t.Errorf("result.Doc = %v, want %v", result.Doc, "new_value")
	}
	if diff := cmp.Diff(doc, result.Old); diff != "" {
		t.Errorf("result.Old mismatch (-want +got):\n%s", diff)
	}
}

func TestReplace_InterfaceMethods(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	if got := replaceOp.Op(); got != internal.OpReplaceType {
		t.Errorf("Op() = %v, want %v", got, internal.OpReplaceType)
	}
	if got := replaceOp.Code(); got != internal.OpReplaceCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpReplaceCode)
	}
	if diff := cmp.Diff([]string{"test"}, replaceOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	if replaceOp.Value != "value" {
		t.Errorf("Value = %v, want %v", replaceOp.Value, "value")
	}
}

func TestReplace_ToJSON(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	got, err := replaceOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	if got.Op != "replace" {
		t.Errorf("ToJSON().Op = %v, want %v", got.Op, "replace")
	}
	if got.Path != "/test" {
		t.Errorf("ToJSON().Path = %v, want %v", got.Path, "/test")
	}
	if got.Value != "value" {
		t.Errorf("ToJSON().Value = %v, want %v", got.Value, "value")
	}
}

func TestReplace_ToCompact(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")

	compact, err := replaceOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	if compact[0] != internal.OpReplaceCode {
		t.Errorf("compact[0] = %v, want %v", compact[0], internal.OpReplaceCode)
	}
	if diff := cmp.Diff([]string{"test"}, compact[1]); diff != "" {
		t.Errorf("compact[1] mismatch (-want +got):\n%s", diff)
	}
	if compact[2] != "value" {
		t.Errorf("compact[2] = %v, want %v", compact[2], "value")
	}
}

func TestReplace_Validate(t *testing.T) {
	replaceOp := NewReplace([]string{"test"}, "value")
	if err := replaceOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	replaceOp = NewReplace([]string{}, "value")
	err := replaceOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathEmpty)
	}
}
