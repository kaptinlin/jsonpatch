package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpSplit_Apply(t *testing.T) {
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
			name:     "split in middle",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "hello world"},
			pos:      5.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"hello", " world"}},
			oldValue: "hello world",
		},
		{
			name:     "split at start",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "world"},
			pos:      0.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"", "world"}},
			oldValue: "world",
		},
		{
			name:     "split at end",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "hello"},
			pos:      5.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"hello", ""}},
			oldValue: "hello",
		},
		{
			name:     "split unicode",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "你好世界"},
			pos:      2.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"你好", "世界"}},
			oldValue: "你好世界",
		},
		{
			name:     "split in nested",
			path:     []string{"user", "bio"},
			doc:      map[string]interface{}{"user": map[string]interface{}{"bio": "Go developer"}},
			pos:      2.0,
			props:    nil,
			expected: map[string]interface{}{"user": map[string]interface{}{"bio": []interface{}{"Go", " developer"}}},
			oldValue: "Go developer",
		},
		{
			name:     "split in array element",
			path:     []string{"lines", "1"},
			doc:      map[string]interface{}{"lines": []interface{}{"foo", "hello world", "baz"}},
			pos:      5.0,
			props:    nil,
			expected: map[string]interface{}{"lines": []interface{}{"foo", "hello", " world", "baz"}},
			oldValue: "hello world",
		},
		{
			name:     "split at root",
			path:     []string{},
			doc:      "abc",
			pos:      1.0,
			props:    nil,
			expected: []interface{}{"a", "bc"},
			oldValue: "abc",
		},
		{
			name:     "split with props",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "hello world"},
			pos:      5.0,
			props:    map[string]interface{}{"type": "split"},
			expected: map[string]interface{}{"text": []interface{}{map[string]interface{}{"text": "hello", "type": "split"}, map[string]interface{}{"text": " world", "type": "split"}}},
			oldValue: "hello world",
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]interface{}{"text": "abc"},
			pos:     1.0,
			props:   nil,
			wantErr: true,
		},
		{
			name:     "not a string",
			path:     []string{"num"},
			doc:      map[string]interface{}{"num": 123},
			pos:      1.0,
			props:    nil,
			expected: map[string]interface{}{"num": []interface{}{float64(1), float64(122)}},
			oldValue: 123,
		},
		{
			name:     "root not a string",
			path:     []string{},
			doc:      123,
			pos:      1.0,
			props:    nil,
			expected: []interface{}{float64(1), float64(122)},
			oldValue: 123,
		},
		{
			name:     "split position out of range",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "abc"},
			pos:      10.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"abc", ""}},
			oldValue: "abc",
		},
		{
			name:     "split negative position",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "abc"},
			pos:      -1.0,
			props:    nil,
			expected: map[string]interface{}{"text": []interface{}{"ab", "c"}},
			oldValue: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewSplit(tt.path, tt.pos, tt.props)
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

func TestOpSplit_Op(t *testing.T) {
	op := NewSplit([]string{"text"}, 1, nil)
	assert.Equal(t, internal.OpSplitType, op.Op())
}

func TestOpSplit_Code(t *testing.T) {
	op := NewSplit([]string{"text"}, 1, nil)
	assert.Equal(t, internal.OpSplitCode, op.Code())
}

func TestOpSplit_NewOpSplit(t *testing.T) {
	path := []string{"user", "bio"}
	pos := 2.0
	props := map[string]interface{}{"type": "split"}
	op := NewSplit(path, pos, props)
	assert.Equal(t, path, op.Path())
	assert.Equal(t, pos, op.Pos)
	assert.Equal(t, props, op.Props)
	assert.Equal(t, internal.OpSplitType, op.Op())
	assert.Equal(t, internal.OpSplitCode, op.Code())
}
