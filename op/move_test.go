package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestMove_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{
			"nested": "value",
		},
	}

	moveOp := NewMove([]string{"qux", "moved"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}
	if _, ok := modifiedDoc["foo"]; ok {
		t.Error("modifiedDoc contains key \"foo\" after move")
	}
	if got := modifiedDoc["qux"].(map[string]any)["moved"]; got != "bar" {
		t.Errorf("modifiedDoc[qux][moved] = %v, want %v", got, "bar")
	}
	if got := modifiedDoc["baz"]; got != 123 {
		t.Errorf("modifiedDoc[baz] = %v, want %v", got, 123)
	}
}

func TestMove_Array(t *testing.T) {
	doc := map[string]any{
		"items": []any{
			"first",
			"second",
			"third",
		},
		"target": map[string]any{},
	}

	moveOp := NewMove([]string{"target", "moved"}, []string{"items", "1"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	modifiedDoc := result.Doc.(map[string]any)
	items := modifiedDoc["items"].([]any)
	target := modifiedDoc["target"].(map[string]any)

	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}
	if len(items) != 2 {
		t.Fatalf("len(items) = %d, want %d", len(items), 2)
	}
	if items[0] != "first" {
		t.Errorf("items[0] = %v, want %v", items[0], "first")
	}
	if items[1] != "third" {
		t.Errorf("items[1] = %v, want %v", items[1], "third")
	}
	if target["moved"] != "second" {
		t.Errorf("target[moved] = %v, want %v", target["moved"], "second")
	}
}

func TestMove_FromNonExistent(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	moveOp := NewMove([]string{"target"}, []string{"qux"})
	_, err := moveOp.Apply(doc)
	if err == nil {
		t.Error("Apply() expected error for non-existent from path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		t.Errorf("Apply() error = %v, want %v", err, ErrPathNotFound)
	}
}

func TestMove_SamePath(t *testing.T) {
	doc := map[string]any{"foo": 1}
	moveOp := NewMove([]string{"foo"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}
	if diff := cmp.Diff(doc, result.Doc); diff != "" {
		t.Errorf("result.Doc mismatch (-want +got):\n%s", diff)
	}
	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}
}

func TestMove_RootArray(t *testing.T) {
	doc := []any{"first", "second", "third"}
	moveOp := NewMove([]string{"0"}, []string{"2"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	resultArray := result.Doc.([]any)
	want := []any{"third", "first", "second"}
	if diff := cmp.Diff(want, resultArray); diff != "" {
		t.Errorf("result.Doc mismatch (-want +got):\n%s", diff)
	}
	if result.Old != "first" {
		t.Errorf("result.Old = %v, want %v", result.Old, "first")
	}
}

func TestMove_EmptyPath(t *testing.T) {
	moveOp := NewMove([]string{}, []string{"foo"})
	err := moveOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathEmpty)
	}
}

func TestMove_EmptyFrom(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{})
	err := moveOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty from path")
	}
	if !errors.Is(err, ErrFromPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrFromPathEmpty)
	}
}

func TestMove_InterfaceMethods(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	if got := moveOp.Op(); got != internal.OpMoveType {
		t.Errorf("Op() = %v, want %v", got, internal.OpMoveType)
	}
	if got := moveOp.Code(); got != internal.OpMoveCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpMoveCode)
	}
	if diff := cmp.Diff([]string{"target"}, moveOp.Path()); diff != "" {
		t.Errorf("Path() mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"source"}, moveOp.From()); diff != "" {
		t.Errorf("From() mismatch (-want +got):\n%s", diff)
	}
	if !moveOp.HasFrom() {
		t.Error("HasFrom() = false, want true")
	}
}

func TestMove_ToJSON(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	got, err := moveOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	if got.Op != "move" {
		t.Errorf("ToJSON().Op = %v, want %v", got.Op, "move")
	}
	if got.Path != "/target" {
		t.Errorf("ToJSON().Path = %v, want %v", got.Path, "/target")
	}
	if got.From != "/source" {
		t.Errorf("ToJSON().From = %v, want %v", got.From, "/source")
	}
}

func TestMove_ToCompact(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	compact, err := moveOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	if compact[0] != internal.OpMoveCode {
		t.Errorf("compact[0] = %v, want %v", compact[0], internal.OpMoveCode)
	}
	if diff := cmp.Diff([]string{"target"}, compact[1]); diff != "" {
		t.Errorf("compact[1] mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]string{"source"}, compact[2]); diff != "" {
		t.Errorf("compact[2] mismatch (-want +got):\n%s", diff)
	}
}

func TestMove_Validate(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})
	if err := moveOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	moveOp = NewMove([]string{}, []string{"source"})
	err := moveOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathEmpty)
	}

	moveOp = NewMove([]string{"target"}, []string{})
	err = moveOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for empty from path")
	}
	if !errors.Is(err, ErrFromPathEmpty) {
		t.Errorf("Validate() error = %v, want %v", err, ErrFromPathEmpty)
	}

	moveOp = NewMove([]string{"same"}, []string{"same"})
	err = moveOp.Validate()
	if err == nil {
		t.Error("Validate() expected error for identical paths")
	}
	if !errors.Is(err, ErrPathsIdentical) {
		t.Errorf("Validate() error = %v, want %v", err, ErrPathsIdentical)
	}
}

func TestMove_RFC6902_RemoveAddPattern(t *testing.T) {
	tests := []struct {
		name     string
		doc      map[string]any
		from     []string
		path     []string
		expected map[string]any
	}{
		{
			name: "object property to array element",
			doc: map[string]any{
				"baz": []any{map[string]any{"qux": "hello"}},
				"bar": 1,
			},
			from: []string{"baz", "0", "qux"},
			path: []string{"baz", "1"},
			expected: map[string]any{
				"baz": []any{map[string]any{}, "hello"},
				"bar": 1,
			},
		},
		{
			name: "array element to front",
			doc: map[string]any{
				"users": []any{
					map[string]any{"name": "Alice"},
					map[string]any{"name": "Bob"},
				},
			},
			from: []string{"users", "1"},
			path: []string{"users", "0"},
			expected: map[string]any{
				"users": []any{
					map[string]any{"name": "Bob"},
					map[string]any{"name": "Alice"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moveOp := NewMove(tt.path, tt.from)
			result, err := moveOp.Apply(tt.doc)
			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, result.Doc); diff != "" {
				t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
