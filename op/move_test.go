package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestMove_Basic(t *testing.T) {
	t.Parallel()
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
		assert.Equal(t, "bar", got, "modifiedDoc[qux][moved]")
	}
	if got := modifiedDoc["baz"]; got != 123 {
		assert.Equal(t, 123, got, "modifiedDoc[baz]")
	}
}

func TestMove_Array(t *testing.T) {
	t.Parallel()
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
	assert.Equal(t, "first", items[0], "items[0]")
	assert.Equal(t, "third", items[1], "items[1]")
	assert.Equal(t, "second", target["moved"], "target[moved]")
}

func TestMove_FromNonExistent(t *testing.T) {
	t.Parallel()
	doc := map[string]any{
		"foo": "bar",
	}

	moveOp := NewMove([]string{"target"}, []string{"qux"})
	_, err := moveOp.Apply(doc)
	if err == nil {
		assert.Fail(t, "Apply() expected error for non-existent from path")
	}
	if !errors.Is(err, ErrPathNotFound) {
		assert.Equal(t, ErrPathNotFound, err, "Apply() error")
	}
}

func TestMove_SamePath(t *testing.T) {
	t.Parallel()
	doc := map[string]any{"foo": 1}
	moveOp := NewMove([]string{"foo"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}
	assert.Equal(t, doc, result.Doc)
	if result.Old != nil {
		t.Errorf("result.Old = %v, want nil", result.Old)
	}
}

func TestMove_RootArray(t *testing.T) {
	t.Parallel()
	doc := []any{"first", "second", "third"}
	moveOp := NewMove([]string{"0"}, []string{"2"})
	result, err := moveOp.Apply(doc)
	if err != nil {
		t.Fatalf("Apply() unexpected error: %v", err)
	}

	resultArray := result.Doc.([]any)
	want := []any{"third", "first", "second"}
	assert.Equal(t, want, resultArray)
	assert.Equal(t, "first", result.Old, "result.Old")
}

func TestMove_EmptyPath(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{}, []string{"foo"})
	err := moveOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		assert.Equal(t, ErrPathEmpty, err, "Validate() error")
	}
}

func TestMove_EmptyFrom(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{"target"}, []string{})
	err := moveOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty from path")
	}
	if !errors.Is(err, ErrFromPathEmpty) {
		assert.Equal(t, ErrFromPathEmpty, err, "Validate() error")
	}
}

func TestMove_InterfaceMethods(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{"target"}, []string{"source"})

	if got := moveOp.Op(); got != internal.OpMoveType {
		assert.Equal(t, internal.OpMoveType, got, "Op()")
	}
	if got := moveOp.Code(); got != internal.OpMoveCode {
		assert.Equal(t, internal.OpMoveCode, got, "Code()")
	}
	assert.Equal(t, []string{"target"}, moveOp.Path(), "Path()")
	assert.Equal(t, []string{"source"}, moveOp.From(), "From()")
	if !moveOp.HasFrom() {
		assert.Fail(t, "HasFrom() = false, want true")
	}
}

func TestMove_ToJSON(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{"target"}, []string{"source"})

	got, err := moveOp.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() unexpected error: %v", err)
	}

	assert.Equal(t, "move", got.Op, "ToJSON().Op")
	assert.Equal(t, "/target", got.Path, "ToJSON().Path")
	assert.Equal(t, "/source", got.From, "ToJSON().From")
}

func TestMove_ToCompact(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{"target"}, []string{"source"})

	compact, err := moveOp.ToCompact()
	if err != nil {
		t.Fatalf("ToCompact() unexpected error: %v", err)
	}
	if len(compact) != 3 {
		t.Fatalf("len(ToCompact()) = %d, want %d", len(compact), 3)
	}
	assert.Equal(t, internal.OpMoveCode, compact[0], "compact[0]")
	assert.Equal(t, []string{"target"}, compact[1])
	assert.Equal(t, []string{"source"}, compact[2])
}

func TestMove_Validate(t *testing.T) {
	t.Parallel()
	moveOp := NewMove([]string{"target"}, []string{"source"})
	if err := moveOp.Validate(); err != nil {
		t.Errorf("Validate() unexpected error: %v", err)
	}

	moveOp = NewMove([]string{}, []string{"source"})
	err := moveOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty path")
	}
	if !errors.Is(err, ErrPathEmpty) {
		assert.Equal(t, ErrPathEmpty, err, "Validate() error")
	}

	moveOp = NewMove([]string{"target"}, []string{})
	err = moveOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for empty from path")
	}
	if !errors.Is(err, ErrFromPathEmpty) {
		assert.Equal(t, ErrFromPathEmpty, err, "Validate() error")
	}

	moveOp = NewMove([]string{"same"}, []string{"same"})
	err = moveOp.Validate()
	if err == nil {
		assert.Fail(t, "Validate() expected error for identical paths")
	}
	if !errors.Is(err, ErrPathsIdentical) {
		assert.Equal(t, ErrPathsIdentical, err, "Validate() error")
	}
}

func TestMove_RFC6902_RemoveAddPattern(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			moveOp := NewMove(tt.path, tt.from)
			result, err := moveOp.Apply(tt.doc)
			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			assert.Equal(t, tt.expected, result.Doc)
		})
	}
}
