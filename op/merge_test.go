package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpMerge_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      interface{}
		pos      float64
		props    interface{}
		expected interface{}
		oldValue interface{}
		wantErr  bool
	}{
		{
			name:     "merge strings",
			path:     []string{"lines", "1"}, // Path to array element, not array itself
			doc:      map[string]interface{}{"lines": []interface{}{"hello", " world", "!"}},
			pos:      1.0, // Merge position in the array (element 1 with element 0)
			props:    nil,
			expected: map[string]interface{}{"lines": []interface{}{"hello world", "!"}},
			oldValue: []interface{}{"hello", " world"},
		},
		{
			name:     "merge at end",
			path:     []string{"lines", "1"}, // Path to array element
			doc:      map[string]interface{}{"lines": []interface{}{"hello", " world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]interface{}{"lines": []interface{}{"hello world"}},
			oldValue: []interface{}{"hello", " world"},
		},
		{
			name:     "merge non-strings",
			path:     []string{"items", "1"},
			doc:      map[string]interface{}{"items": []interface{}{1, 2, 3}},
			pos:      1.0,
			props:    nil,
			expected: map[string]interface{}{"items": []interface{}{float64(3), 3}},
			oldValue: []interface{}{1, 2},
		},
		{
			name:     "merge mixed types",
			path:     []string{"items", "1"},
			doc:      map[string]interface{}{"items": []interface{}{"hello", 123, "world"}},
			pos:      1.0,
			props:    nil,
			expected: map[string]interface{}{"items": []interface{}{[]interface{}{"hello", 123}, "world"}},
			oldValue: []interface{}{"hello", 123},
		},
		{
			name:     "merge in nested",
			path:     []string{"user", "tags", "1"}, // Path to array element
			doc:      map[string]interface{}{"user": map[string]interface{}{"tags": []interface{}{"go", "lang", "dev"}}},
			pos:      1.0,
			props:    nil,
			expected: map[string]interface{}{"user": map[string]interface{}{"tags": []interface{}{"golang", "dev"}}},
			oldValue: []interface{}{"go", "lang"},
		},
		{
			name:     "merge at root",
			path:     []string{},
			doc:      []interface{}{"a", "b", "c"},
			pos:      1.0,
			props:    nil,
			expected: []interface{}{"ab", "c"},
			oldValue: []interface{}{"a", "b"},
		},
		{
			name:     "merge with props",
			path:     []string{"lines", "1"},
			doc:      map[string]interface{}{"lines": []interface{}{"hello", " world"}},
			pos:      1.0,
			props:    map[string]interface{}{"type": "merge"},
			expected: map[string]interface{}{"lines": []interface{}{"hello world"}},
			oldValue: []interface{}{"hello", " world"},
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]interface{}{"lines": []interface{}{"a", "b"}},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "not an array",
			path:    []string{"text"},
			doc:     map[string]interface{}{"text": "abc"},
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
			doc:     map[string]interface{}{"lines": []interface{}{"a", "b"}},
			pos:     2.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "merge negative position",
			path:    []string{"lines"},
			doc:     map[string]interface{}{"lines": []interface{}{"a", "b"}},
			pos:     -1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "merge position zero (invalid)",
			path:    []string{"lines"},
			doc:     map[string]interface{}{"lines": []interface{}{"a", "b"}},
			pos:     0.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:    "single element array",
			path:    []string{"lines"},
			doc:     map[string]interface{}{"lines": []interface{}{"a"}},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var props map[string]interface{}
			if tt.props != nil {
				props = tt.props.(map[string]interface{})
			}
			op := NewMerge(tt.path, tt.pos, props)
			docCopy, err := DeepClone(tt.doc)
			require.NoError(t, err)

			result, err := op.Apply(docCopy)

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

func TestOpMerge_Op(t *testing.T) {
	op := NewMerge([]string{"lines"}, 0, nil)
	assert.Equal(t, internal.OpMergeType, op.Op())
}

func TestOpMerge_Code(t *testing.T) {
	op := NewMerge([]string{"lines"}, 0, nil)
	assert.Equal(t, internal.OpMergeCode, op.Code())
}

func TestOpMerge_NewOpMerge(t *testing.T) {
	path := []string{"user", "tags"}
	pos := 1.0
	props := map[string]interface{}{"type": "merge"}
	op := NewMerge(path, pos, props)
	assert.Equal(t, path, op.Path())
	assert.Equal(t, pos, op.Pos)
	assert.Equal(t, props, op.Props)
	assert.Equal(t, internal.OpMergeType, op.Op())
	assert.Equal(t, internal.OpMergeCode, op.Code())
}

func TestOpMerge_TypeScript_Compatibility(t *testing.T) {
	// Test cases based on TypeScript reference implementation
	tests := []struct {
		name     string
		doc      interface{}
		path     []string
		pos      float64
		expected interface{}
	}{
		{
			name:     "merge two text nodes in array",
			doc:      []interface{}{map[string]interface{}{"text": "foo"}, map[string]interface{}{"text": "bar"}},
			path:     []string{"1"}, // Path to second element
			pos:      1,             // Merge with previous element
			expected: []interface{}{map[string]interface{}{"text": "foobar"}},
		},
		{
			name:     "merge in nested array",
			doc:      map[string]interface{}{"foo": []interface{}{map[string]interface{}{"children": []interface{}{"1", "2"}}, map[string]interface{}{"children": []interface{}{"3", "4"}}}},
			path:     []string{"foo", "1"},
			pos:      1,
			expected: map[string]interface{}{"foo": []interface{}{map[string]interface{}{"children": []interface{}{"1", "2", "3", "4"}}}},
		},
		{
			name:     "merge simple strings",
			doc:      []interface{}{"hello", " world"},
			path:     []string{"1"},
			pos:      1,
			expected: []interface{}{"hello world"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mergeOp := NewMerge(tt.path, tt.pos, nil)
			result, err := mergeOp.Apply(tt.doc)
			require.NoError(t, err, "Merge operation should work")
			assert.Equal(t, tt.expected, result.Doc, "Result should match TypeScript behavior")
		})
	}
}
