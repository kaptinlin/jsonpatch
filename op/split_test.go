package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestSplit_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      any
		pos      float64
		props    any
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "split in middle",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello world"},
			pos:      5.0,
			props:    nil,
			expected: map[string]any{"text": []any{"hello", " world"}},
			oldValue: "hello world",
		},
		{
			name:     "split at start",
			path:     []string{"text"},
			doc:      map[string]any{"text": "world"},
			pos:      0.0,
			props:    nil,
			expected: map[string]any{"text": []any{"", "world"}},
			oldValue: "world",
		},
		{
			name:     "split at end",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello"},
			pos:      5.0,
			props:    nil,
			expected: map[string]any{"text": []any{"hello", ""}},
			oldValue: "hello",
		},
		{
			name:     "split unicode",
			path:     []string{"text"},
			doc:      map[string]any{"text": "你好世界"},
			pos:      2.0,
			props:    nil,
			expected: map[string]any{"text": []any{"你好", "世界"}},
			oldValue: "你好世界",
		},
		{
			name:     "split in nested",
			path:     []string{"user", "bio"},
			doc:      map[string]any{"user": map[string]any{"bio": "Go developer"}},
			pos:      2.0,
			props:    nil,
			expected: map[string]any{"user": map[string]any{"bio": []any{"Go", " developer"}}},
			oldValue: "Go developer",
		},
		{
			name:     "split in array element",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"foo", "hello world", "baz"}},
			pos:      5.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"foo", "hello", " world", "baz"}},
			oldValue: "hello world",
		},
		{
			name:     "split at root",
			path:     []string{},
			doc:      "abc",
			pos:      1.0,
			props:    nil,
			expected: []any{"a", "bc"},
			oldValue: "abc",
		},
		{
			name:     "split with props",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello world"},
			pos:      5.0,
			props:    map[string]any{"type": "split"},
			expected: map[string]any{"text": []any{map[string]any{"text": "hello", "type": "split"}, map[string]any{"text": " world", "type": "split"}}},
			oldValue: "hello world",
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]any{"text": "abc"},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:     "not a string",
			path:     []string{"num"},
			doc:      map[string]any{"num": 123},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"num": []any{float64(1), float64(122)}},
			oldValue: 123,
		},
		{
			name:     "root not a string",
			path:     []string{},
			doc:      123,
			pos:      1.0,
			props:    nil,
			expected: []any{float64(1), float64(122)},
			oldValue: 123,
		},
		{
			name:     "split position out of range",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      10.0,
			props:    nil,
			expected: map[string]any{"text": []any{"abc", ""}},
			oldValue: "abc",
		},
		{
			name:     "split negative position",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      -1.0,
			props:    nil,
			expected: map[string]any{"text": []any{"ab", "c"}},
			oldValue: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitOp := NewSplit(tt.path, tt.pos, tt.props)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := splitOp.Apply(docCopy)

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

func TestSplit_Constructor(t *testing.T) {
	path := []string{"user", "bio"}
	pos := 2.0
	props := map[string]any{"type": "split"}
	splitOp := NewSplit(path, pos, props)
	if diff := cmp.Diff(path, splitOp.Path()); diff != "" {
		t.Errorf("NewSplit() Path mismatch (-want +got):\n%s", diff)
	}
	if splitOp.Pos != pos {
		t.Errorf("NewSplit() Pos = %v, want %v", splitOp.Pos, pos)
	}
	if diff := cmp.Diff(props, splitOp.Props); diff != "" {
		t.Errorf("NewSplit() Props mismatch (-want +got):\n%s", diff)
	}
	if got := splitOp.Op(); got != internal.OpSplitType {
		t.Errorf("Op() = %v, want %v", got, internal.OpSplitType)
	}
	if got := splitOp.Code(); got != internal.OpSplitCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpSplitCode)
	}
}

func TestSplit_TypeScript_Compatibility(t *testing.T) {
	tests := []struct {
		name     string
		doc      any
		path     []string
		pos      float64
		props    any
		expected any
	}{
		{
			name:     "string without props",
			doc:      map[string]any{"text": "hello"},
			path:     []string{"text"},
			pos:      2,
			props:    nil,
			expected: map[string]any{"text": []any{"he", "llo"}},
		},
		{
			name:  "string with props",
			doc:   map[string]any{"text": "hello"},
			path:  []string{"text"},
			pos:   2,
			props: map[string]any{"bold": true},
			expected: map[string]any{"text": []any{
				map[string]any{"text": "he", "bold": true},
				map[string]any{"text": "llo", "bold": true},
			}},
		},
		{
			name:     "number",
			doc:      map[string]any{"num": 10},
			path:     []string{"num"},
			pos:      3,
			props:    nil,
			expected: map[string]any{"num": []any{3.0, 7.0}},
		},
		{
			name:     "root array element",
			doc:      []any{"hello world"},
			path:     []string{"0"},
			pos:      5,
			props:    nil,
			expected: []any{"hello", " world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			splitOp := NewSplit(tt.path, tt.pos, tt.props)
			result, err := splitOp.Apply(tt.doc)
			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, result.Doc); diff != "" {
				t.Errorf("Apply() Doc mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
