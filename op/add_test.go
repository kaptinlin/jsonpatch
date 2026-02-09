package op

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestAdd_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	addOp := NewAdd([]string{"baz"}, "qux")

	result, err := addOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	if got := doc["foo"]; got != "bar" {
		t.Errorf("doc[foo] = %v, want %v", got, "bar")
	}
	if got := doc["baz"]; got != "qux" {
		t.Errorf("doc[baz] = %v, want %v", got, "qux")
	}

	resultDoc := result.Doc.(map[string]any)
	if got := resultDoc["foo"]; got != "bar" {
		t.Errorf("resultDoc[foo] = %v, want %v", got, "bar")
	}
	if got := resultDoc["baz"]; got != "qux" {
		t.Errorf("resultDoc[baz] = %v, want %v", got, "qux")
	}
	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}

	// Mutate behavior: result should point to the same underlying object
	if reflect.ValueOf(doc).Pointer() != reflect.ValueOf(resultDoc).Pointer() {
		t.Error("Result should point to the same document object")
	}
}

func TestAdd_ReplaceExisting(t *testing.T) {
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
		t.Errorf("resultDoc[foo] = %v, want %v", got, "new_value")
	}
	if got := result.Old; got != "bar" {
		t.Errorf("result.Old = %v, want %v", got, "bar")
	}
}

func TestAdd_InterfaceMethods(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")

	if got := addOp.Op(); got != internal.OpAddType {
		t.Errorf("Op() = %v, want %v", got, internal.OpAddType)
	}
	if got := addOp.Code(); got != internal.OpAddCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpAddCode)
	}
	if diff := cmp.Diff([]string{"foo"}, addOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
}

func TestAdd_ToJSON(t *testing.T) {
	addOp := NewAdd([]string{"foo", "bar"}, "baz")

	got, err := addOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	if got.Op != "add" {
		t.Errorf("ToJSON().Op = %v, want %v", got.Op, "add")
	}
	if got.Path != "/foo/bar" {
		t.Errorf("ToJSON().Path = %v, want %v", got.Path, "/foo/bar")
	}
	if got.Value != "baz" {
		t.Errorf("ToJSON().Value = %v, want %v", got.Value, "baz")
	}
}

func TestAdd_ToCompact(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")

	compact, err := addOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	if compact[0] != internal.OpAddCode {
		t.Errorf("compact[0] = %v, want %v", compact[0], internal.OpAddCode)
	}
	if diff := cmp.Diff([]string{"foo"}, compact[1]); diff != "" {
		t.Errorf("compact[1] mismatch (-want +got):\n%s", diff)
	}
	if compact[2] != "bar" {
		t.Errorf("compact[2] = %v, want %v", compact[2], "bar")
	}
}

func TestAdd_Validate(t *testing.T) {
	addOp := NewAdd([]string{"foo"}, "bar")
	if err := addOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	addOp = NewAdd([]string{}, "bar")
	err := addOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathEmpty)
	}
}

func TestAdd_Constructor(t *testing.T) {
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
	if got.Value != value {
		t.Errorf("ToJSON().Value = %v, want %v", got.Value, value)
	}
}
