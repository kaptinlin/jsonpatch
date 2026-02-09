package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			require.NoError(t, err)

			result, err := mergeOp.Apply(docCopy)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, tt.oldValue, result.Old)
		})
	}
}

func TestMerge_Constructor(t *testing.T) {
	path := []string{"user", "tags"}
	pos := 1.0
	props := map[string]any{"type": "merge"}
	mergeOp := NewMerge(path, pos, props)
	assert.Equal(t, path, mergeOp.Path())
	assert.Equal(t, pos, mergeOp.Pos)
	assert.Equal(t, props, mergeOp.Props)
	assert.Equal(t, internal.OpMergeType, mergeOp.Op())
	assert.Equal(t, internal.OpMergeCode, mergeOp.Code())
}
