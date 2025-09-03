package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpStrIns_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      interface{}
		pos      float64
		str      string
		expected interface{}
		oldValue interface{}
		wantErr  bool
	}{
		{
			name:     "insert in middle",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "hello world"},
			pos:      5.0,
			str:      ", brave new",
			expected: map[string]interface{}{"text": "hello, brave new world"},
			oldValue: "hello world",
		},
		{
			name:     "insert at start",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "world"},
			pos:      0.0,
			str:      "hello ",
			expected: map[string]interface{}{"text": "hello world"},
			oldValue: "world",
		},
		{
			name:     "insert at end",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "hello"},
			pos:      5.0,
			str:      " world",
			expected: map[string]interface{}{"text": "hello world"},
			oldValue: "hello",
		},
		{
			name:     "insert unicode",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "你好世界"},
			pos:      2.0,
			str:      "，美丽的",
			expected: map[string]interface{}{"text": "你好，美丽的世界"},
			oldValue: "你好世界",
		},
		{
			name:     "insert in nested",
			path:     []string{"user", "bio"},
			doc:      map[string]interface{}{"user": map[string]interface{}{"bio": "Go dev"}},
			pos:      2.0,
			str:      "lang ",
			expected: map[string]interface{}{"user": map[string]interface{}{"bio": "Golang  dev"}},
			oldValue: "Go dev",
		},
		{
			name:     "insert in array element",
			path:     []string{"lines", "1"},
			doc:      map[string]interface{}{"lines": []interface{}{"foo", "bar", "baz"}},
			pos:      1.0,
			str:      "-insert-",
			expected: map[string]interface{}{"lines": []interface{}{"foo", "b-insert-ar", "baz"}},
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
			doc:     map[string]interface{}{"text": "abc"},
			pos:     1.0,
			str:     "-",
			wantErr: true,
		},
		{
			name:    "not a string",
			path:    []string{"num"},
			doc:     map[string]interface{}{"num": 123},
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
			doc:      map[string]interface{}{"text": "abc"},
			pos:      10.0,
			str:      "X",
			expected: map[string]interface{}{"text": "abcX"},
			oldValue: "abc",
		},
		{
			name:     "insert negative position",
			path:     []string{"text"},
			doc:      map[string]interface{}{"text": "abc"},
			pos:      -1.0,
			str:      "X",
			expected: map[string]interface{}{"text": "Xabc"},
			oldValue: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOpStrInsOperation(tt.path, tt.pos, tt.str)
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

func TestOpStrIns_Op(t *testing.T) {
	op := NewOpStrInsOperation([]string{"text"}, 1.0, "-")
	assert.Equal(t, internal.OpStrInsType, op.Op())
}

func TestOpStrIns_Code(t *testing.T) {
	op := NewOpStrInsOperation([]string{"text"}, 1.0, "-")
	assert.Equal(t, internal.OpStrInsCode, op.Code())
}

func TestOpStrIns_NewOpStrInsOperation(t *testing.T) {
	path := []string{"user", "bio"}
	pos := 2.0
	str := "abc"
	op := NewOpStrInsOperation(path, pos, str)
	assert.Equal(t, path, op.Path())
	assert.Equal(t, pos, op.Pos)
	assert.Equal(t, str, op.Str)
	assert.Equal(t, internal.OpStrInsType, op.Op())
	assert.Equal(t, internal.OpStrInsCode, op.Code())
}
