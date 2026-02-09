package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestMerge_Apply(t *testing.T) {
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
			name:     "merge strings",
			path:     []string{"lines"},
			doc:      map[string]any{"lines": []any{"hello", " world", "!"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"hello world", "!"}},
			oldValue: []any{"hello", " world"},
		},
		{
			name:     "merge at end",
			path:     []string{"lines"},
			doc:      map[string]any{"lines": []any{"hello", " world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"hello world"}},
			oldValue: []any{"hello", " world"},
		},
		{
			name:     "merge non-strings",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{1, 2, 3}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"items": []any{float64(3), 3}},
			oldValue: []any{1, 2},
		},
		{
			name:     "merge mixed types",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{"hello", 123, "world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"items": []any{[]any{"hello", 123}, "world"}},
			oldValue: []any{"hello", 123},
		},
		{
			name:     "merge in nested",
			path:     []string{"user", "tags"},
			doc:      map[string]any{"user": map[string]any{"tags": []any{"go", "lang", "dev"}}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"user": map[string]any{"tags": []any{"golang", "dev"}}},
			oldValue: []any{"go", "lang"},
		},
		{
			name:     "merge at root",
			path:     []string{},
			doc:      []any{"a", "b", "c"},
			pos:      1.0,
			props:    nil,
			expected: []any{"ab", "c"},
			oldValue: []any{"a", "b"},
		},
		{
			name:     "merge with props",
			path:     []string{"lines"},
			doc:      map[string]any{"lines": []any{"hello", " world"}},
			pos:      1.0,
			props:    map[string]any{"type": "merge"},
			expected: map[string]any{"lines": []any{"hello world"}},
			oldValue: []any{"hello", " world"},
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]any{"lines": []any{"a", "b"}},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "not an array",
			path:    []string{"text"},
			doc:     map[string]any{"text": "abc"},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "root not an array",
			path:    []string{},
			doc:     "abc",
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "merge position out of range",
			path:    []string{"lines"},
			doc:     map[string]any{"lines": []any{"a", "b"}},
			pos:     2.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "merge negative position",
			path:    []string{"lines"},
			doc:     map[string]any{"lines": []any{"a", "b"}},
			pos:     -1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "merge position zero (invalid)",
			path:    []string{"lines"},
			doc:     map[string]any{"lines": []any{"a", "b"}},
			pos:     0.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "single element array",
			path:    []string{"lines"},
			doc:     map[string]any{"lines": []any{"a"}},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var props map[string]any
			if tt.props != nil {
				props = tt.props.(map[string]any)
			}
			mergeOp := NewMerge(tt.path, tt.pos, props)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := mergeOp.Apply(docCopy)

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

func TestMerge_Constructor(t *testing.T) {
	path := []string{"user", "tags"}
	pos := 1.0
	props := map[string]any{"type": "merge"}
	mergeOp := NewMerge(path, pos, props)
	if diff := cmp.Diff(path, mergeOp.Path()); diff != "" {
		t.Errorf("NewMerge() Path mismatch (-want +got):\n%s", diff)
	}
	if mergeOp.Pos != pos {
		t.Errorf("NewMerge() Pos = %v, want %v", mergeOp.Pos, pos)
	}
	if diff := cmp.Diff(props, mergeOp.Props); diff != "" {
		t.Errorf("NewMerge() Props mismatch (-want +got):\n%s", diff)
	}
	if got := mergeOp.Op(); got != internal.OpMergeType {
		t.Errorf("Op() = %v, want %v", got, internal.OpMergeType)
	}
	if got := mergeOp.Code(); got != internal.OpMergeCode {
		t.Errorf("Code() = %v, want %v", got, internal.OpMergeCode)
	}
}
