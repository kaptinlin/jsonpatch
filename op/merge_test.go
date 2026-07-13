package op

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestMergeElementPathContract(t *testing.T) {
	t.Parallel()

	t.Run("root and nested arrays", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name string
			path []string
			doc  any
			want any
		}{
			{name: "root", path: []string{"1"}, doc: []any{"a", "b", "c"}, want: []any{"ab", "c"}},
			{name: "nested", path: []string{"nodes", "1"}, doc: map[string]any{"nodes": []any{"a", "b", "c"}}, want: map[string]any{"nodes": []any{"ab", "c"}}},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				result, err := NewMerge(tc.path, 99, nil).Apply(tc.doc)
				require.NoError(t, err)
				assert.Equal(t, tc.want, result.Doc)
				assert.Equal(t, []any{"a", "b"}, result.Old)
			})
		}
	})

	t.Run("pos is metadata", func(t *testing.T) {
		t.Parallel()
		for _, pos := range []float64{-7, 0, 3} {
			pos := pos
			t.Run(fmt.Sprintf("%g", pos), func(t *testing.T) {
				t.Parallel()
				result, err := NewMerge([]string{"1"}, pos, nil).Apply([]any{"a", "b", "c"})
				require.NoError(t, err)
				assert.Equal(t, []any{"ab", "c"}, result.Doc)
				require.NoError(t, NewMerge([]string{"1"}, pos, nil).Validate())
			})
		}
	})

	t.Run("invalid targets leave document unchanged", func(t *testing.T) {
		t.Parallel()
		tests := []struct {
			name    string
			path    []string
			doc     any
			wantErr error
		}{
			{name: "first element", path: []string{"nodes", "0"}, doc: map[string]any{"nodes": []any{"a", "b"}}, wantErr: ErrInvalidIndex},
			{name: "missing target", path: []string{"nodes", "2"}, doc: map[string]any{"nodes": []any{"a", "b"}}, wantErr: ErrIndexOutOfRange},
			{name: "path names array instead of second element", path: []string{"nodes"}, doc: map[string]any{"nodes": []any{"a", "b"}}, wantErr: ErrInvalidTarget},
			{name: "non-array parent", path: []string{"text", "1"}, doc: map[string]any{"text": "ab"}, wantErr: ErrInvalidTarget},
			{name: "non-canonical index", path: []string{"nodes", "01"}, doc: map[string]any{"nodes": []any{"a", "b"}}, wantErr: ErrPathNotFound},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				want, err := DeepClone(tc.doc)
				require.NoError(t, err)
				result, err := NewMerge(tc.path, 1, nil).Apply(tc.doc)
				require.ErrorIs(t, err, tc.wantErr)
				assert.Equal(t, internal.OpResult[any]{}, result)
				if diff := cmp.Diff(want, tc.doc); diff != "" {
					t.Errorf("failed merge mutated document (-want +got):\n%s", diff)
				}
			})
		}
	})
}

func TestMergeElementValueFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		one  any
		two  any
		want any
	}{
		{name: "string", one: "hello", two: " world", want: "hello world"},
		{name: "number", one: 2, two: 3, want: float64(5)},
		{
			name: "Slate text",
			one:  map[string]any{"text": "A", "bold": true},
			two:  map[string]any{"text": "da", "italic": true},
			want: map[string]any{"text": "Ada", "bold": true, "italic": true},
		},
		{
			name: "Slate element",
			one:  map[string]any{"type": "p", "children": []any{"A"}},
			two:  map[string]any{"align": "center", "children": []any{"B"}},
			want: map[string]any{"type": "p", "align": "center", "children": []any{"A", "B"}},
		},
		{name: "fallback pair", one: true, two: nil, want: []any{true, nil}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			result, err := NewMerge([]string{"1"}, 0, nil).Apply([]any{tc.one, tc.two, "tail"})
			require.NoError(t, err)
			assert.Equal(t, []any{tc.want, "tail"}, result.Doc)
			assert.Equal(t, []any{tc.one, tc.two}, result.Old)
		})
	}
}

func TestMerge_Apply(t *testing.T) {
	t.Parallel()
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
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"hello", " world", "!"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"hello world", "!"}},
			oldValue: []any{"hello", " world"},
		},
		{
			name:     "merge at end",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"hello", " world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"hello world"}},
			oldValue: []any{"hello", " world"},
		},
		{
			name:     "merge non-strings",
			path:     []string{"items", "1"},
			doc:      map[string]any{"items": []any{1, 2, 3}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"items": []any{float64(3), 3}},
			oldValue: []any{1, 2},
		},
		{
			name:     "merge mixed types",
			path:     []string{"items", "1"},
			doc:      map[string]any{"items": []any{"hello", 123, "world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"items": []any{[]any{"hello", 123}, "world"}},
			oldValue: []any{"hello", 123},
		},
		{
			name:     "merge in nested",
			path:     []string{"user", "tags", "1"},
			doc:      map[string]any{"user": map[string]any{"tags": []any{"go", "lang", "dev"}}},
			pos:      1.0,
			props:    nil,
			expected: map[string]any{"user": map[string]any{"tags": []any{"golang", "dev"}}},
			oldValue: []any{"go", "lang"},
		},
		{
			name:     "merge at root",
			path:     []string{"1"},
			doc:      []any{"a", "b", "c"},
			pos:      1.0,
			props:    nil,
			expected: []any{"ab", "c"},
			oldValue: []any{"a", "b"},
		},
		{
			name:     "merge with props",
			path:     []string{"lines", "1"},
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
			name:     "merge negative split metadata",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"a", "b"}},
			pos:      -1.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"ab"}},
			oldValue: []any{"a", "b"},
		},
		{
			name:     "merge zero split metadata",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"a", "b"}},
			pos:      0.0,
			props:    nil,
			expected: map[string]any{"lines": []any{"ab"}},
			oldValue: []any{"a", "b"},
		},
		{
			name:    "single element array has no second element",
			path:    []string{"lines", "1"},
			doc:     map[string]any{"lines": []any{"a"}},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var props map[string]any
			if tt.props != nil {
				props = tt.props.(map[string]any)
			}
			mergeOp := NewMerge(tt.path, tt.pos, props)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("DeepClone() error: %v", err))
			}

			result, err := mergeOp.Apply(docCopy)

			if tt.wantErr {
				if err == nil {
					assert.Fail(t, "Apply() expected error, got nil")
				}
				return
			}

			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
			}
			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, tt.oldValue, result.Old)
		})
	}
}

func TestMerge_RejectsMissingElement(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		doc  map[string]any
		path []string
	}{
		{
			name: "past end",
			doc:  map[string]any{"items": []any{"a", "b"}},
			path: []string{"items", "3"},
		},
		{
			name: "empty array",
			doc:  map[string]any{"items": []any{}},
			path: []string{"items", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			want, err := DeepClone(tt.doc)
			require.NoError(t, err)
			result, err := NewMerge(tt.path, 0, nil).Apply(tt.doc)

			require.ErrorIs(t, err, ErrIndexOutOfRange)
			assert.Equal(t, internal.OpResult[any]{}, result)
			assert.Equal(t, want, tt.doc)
		})
	}
}

func TestMerge_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"user", "tags", "1"}
	pos := 1.0
	props := map[string]any{"type": "merge"}
	mergeOp := NewMerge(path, pos, props)
	assert.Equal(t, path, mergeOp.Path(), "Path()")
	assert.Equal(t, pos, mergeOp.Pos, "NewMerge() Pos")
	assert.Equal(t, props, mergeOp.Props)
	if got := mergeOp.Op(); got != internal.OpMergeType {
		assert.Equal(t, internal.OpMergeType, got, "Op()")
	}
	if got := mergeOp.Code(); got != internal.OpMergeCode {
		assert.Equal(t, internal.OpMergeCode, got, "Code()")
	}
}

func TestMerge_SlateNodesIgnoreMetadataProps(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		doc      []any
		props    map[string]any
		expected []any
	}{
		{
			name: "text nodes",
			doc: []any{
				map[string]any{"text": "hello", "bold": true},
				map[string]any{"text": " world", "italic": true},
			},
			props: map[string]any{"mark": "merged"},
			expected: []any{
				map[string]any{"text": "hello world", "bold": true, "italic": true},
			},
		},
		{
			name: "element nodes",
			doc: []any{
				map[string]any{"type": "paragraph", "children": []any{map[string]any{"text": "a"}}},
				map[string]any{"align": "center", "children": []any{map[string]any{"text": "b"}}},
			},
			props: map[string]any{"type": "merged"},
			expected: []any{
				map[string]any{
					"type":     "paragraph",
					"align":    "center",
					"children": []any{map[string]any{"text": "a"}, map[string]any{"text": "b"}},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mergeOp := NewMerge([]string{"1"}, 1, tt.props)

			result, err := mergeOp.Apply(tt.doc)

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, tt.doc, result.Old)
		})
	}
}
