package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestStrIns_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      any
		pos      float64
		str      string
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "insert in middle",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello world"},
			pos:      5.0,
			str:      ", brave new",
			expected: map[string]any{"text": "hello, brave new world"},
			oldValue: "hello world",
		},
		{
			name:     "insert at start",
			path:     []string{"text"},
			doc:      map[string]any{"text": "world"},
			pos:      0.0,
			str:      "hello ",
			expected: map[string]any{"text": "hello world"},
			oldValue: "world",
		},
		{
			name:     "insert at end",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello"},
			pos:      5.0,
			str:      " world",
			expected: map[string]any{"text": "hello world"},
			oldValue: "hello",
		},
		{
			name:     "insert unicode",
			path:     []string{"text"},
			doc:      map[string]any{"text": "你好世界"},
			pos:      2.0,
			str:      "，美丽的",
			expected: map[string]any{"text": "你好，美丽的世界"},
			oldValue: "你好世界",
		},
		{
			name:     "insert in nested",
			path:     []string{"user", "bio"},
			doc:      map[string]any{"user": map[string]any{"bio": "Go dev"}},
			pos:      2.0,
			str:      "lang ",
			expected: map[string]any{"user": map[string]any{"bio": "Golang  dev"}},
			oldValue: "Go dev",
		},
		{
			name:     "insert in array element",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"foo", "bar", "baz"}},
			pos:      1.0,
			str:      "-insert-",
			expected: map[string]any{"lines": []any{"foo", "b-insert-ar", "baz"}},
			oldValue: "bar",
		},
		{
			name:     "insert at root",
			path:     []string{},
			doc:      "abc",
			pos:      1.0,
			str:      "-",
			expected: "a-bc",
			oldValue: "abc",
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]any{"text": "abc"},
			pos:     1.0,
			str:     "-",
			wantErr: true,
		},
		{
			name:    "not a string",
			path:    []string{"num"},
			doc:     map[string]any{"num": 123},
			pos:     1.0,
			str:     "-",
			wantErr: true,
		},
		{
			name:    "root not a string",
			path:    []string{},
			doc:     123,
			pos:     1.0,
			str:     "-",
			wantErr: true,
		},
		{
			name:     "insert position out of range",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      10.0,
			str:      "X",
			expected: map[string]any{"text": "abcX"},
			oldValue: "abc",
		},
		{
			name:     "insert negative position",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      -1.0,
			str:      "X",
			expected: map[string]any{"text": "abXc"},
			oldValue: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strInsOp := NewStrIns(tt.path, tt.pos, tt.str)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := strInsOp.Apply(docCopy)

			if tt.wantErr {
				if err == nil {
					t.Error("Apply() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, result.Doc); diff != "" {
				t.Errorf("Apply() Doc mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.oldValue, result.Old); diff != "" {
				t.Errorf("Apply() Old mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestStrIns_Constructor(t *testing.T) {
	path := []string{"user", "bio"}
	pos := 2.0
	str := "abc"
	strInsOp := NewStrIns(path, pos, str)
	if diff := cmp.Diff(path, strInsOp.Path()); diff != "" {
		t.Errorf("NewStrIns() Path mismatch (-want +got):\n%s", diff)
	}
	if strInsOp.Pos != pos {
		t.Errorf("NewStrIns() Pos = %v, want %v", strInsOp.Pos, pos)
	}
	if strInsOp.Str != str {
		t.Errorf("NewStrIns() Str = %q, want %q", strInsOp.Str, str)
	}
	if got := strInsOp.Op(); got != internal.OpStrInsType {
		t.Errorf("Op() = %v, want %v", got, internal.OpStrInsType)
	}
	if got := strInsOp.Code(); got != internal.OpStrInsCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpStrInsCode)
	}
}
