package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestStrDel_Apply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		path     []string
		doc      any
		pos      float64
		length   float64
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "delete in middle",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello, brave new world"},
			pos:      5.0,
			length:   12.0,
			expected: map[string]any{"text": "helloworld"},
			oldValue: "hello, brave new world",
		},
		{
			name:     "delete at start",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello world"},
			pos:      0.0,
			length:   6.0,
			expected: map[string]any{"text": "world"},
			oldValue: "hello world",
		},
		{
			name:     "delete at end",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello world"},
			pos:      5.0,
			length:   6.0,
			expected: map[string]any{"text": "hello"},
			oldValue: "hello world",
		},
		{
			name:     "delete unicode",
			path:     []string{"text"},
			doc:      map[string]any{"text": "你好，美丽的世界"},
			pos:      2.0,
			length:   4.0,
			expected: map[string]any{"text": "你好世界"},
			oldValue: "你好，美丽的世界",
		},
		{
			name:     "delete in nested",
			path:     []string{"user", "bio"},
			doc:      map[string]any{"user": map[string]any{"bio": "Golang dev"}},
			pos:      2.0,
			length:   4.0,
			expected: map[string]any{"user": map[string]any{"bio": "Go dev"}},
			oldValue: "Golang dev",
		},
		{
			name:     "delete in array element",
			path:     []string{"lines", "1"},
			doc:      map[string]any{"lines": []any{"foo", "b-insert-ar", "baz"}},
			pos:      1.0,
			length:   8.0,
			expected: map[string]any{"lines": []any{"foo", "bar", "baz"}},
			oldValue: "b-insert-ar",
		},
		{
			name:     "delete at root",
			path:     []string{},
			doc:      "a-bc",
			pos:      1.0,
			length:   1.0,
			expected: "abc",
			oldValue: "a-bc",
		},
		{
			name:    "path not found",
			path:    []string{"notfound"},
			doc:     map[string]any{"text": "abc"},
			pos:     1.0,
			length:  1.0,
			wantErr: true,
		},
		{
			name:    "not a string",
			path:    []string{"num"},
			doc:     map[string]any{"num": 123},
			pos:     1.0,
			length:  1.0,
			wantErr: true,
		},
		{
			name:    "root not a string",
			path:    []string{},
			doc:     123,
			pos:     1.0,
			length:  1.0,
			wantErr: true,
		},
		{
			name:     "delete position out of range",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      10.0,
			length:   1.0,
			expected: map[string]any{"text": "abc"},
			oldValue: "abc",
		},
		{
			name:     "delete negative position",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      -1.0,
			length:   1.0,
			expected: map[string]any{"text": "ab"},
			oldValue: "abc",
		},
		{
			name:     "delete negative length",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      1.0,
			length:   -1.0,
			expected: map[string]any{"text": "abc"},
			oldValue: "abc",
		},
		{
			name:     "delete length out of range",
			path:     []string{"text"},
			doc:      map[string]any{"text": "abc"},
			pos:      1.0,
			length:   10.0,
			expected: map[string]any{"text": "a"},
			oldValue: "abc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			strDelOp := NewStrDel(tt.path, tt.pos, tt.length)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := strDelOp.Apply(docCopy)

			if tt.wantErr {
				if err == nil {
					assert.Fail(t, "Apply() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("Apply() unexpected error: %v", err)
			}
			assert.Equal(t, tt.expected, result.Doc)
			assert.Equal(t, tt.oldValue, result.Old)
		})
	}
}

func TestStrDel_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"user", "bio"}
	pos := 2.0
	length := 3.0
	strDelOp := NewStrDel(path, pos, length)
	if diff := cmp.Diff(path, strDelOp.Path()); diff != "" {
		t.Errorf("NewStrDel() Path mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, pos, strDelOp.Pos, "NewStrDel() Pos")
	assert.Equal(t, length, strDelOp.Len, "NewStrDel() Len")
	if got := strDelOp.Op(); got != internal.OpStrDelType {
		assert.Equal(t, internal.OpStrDelType, got, "Op()")
	}
	if got := strDelOp.Code(); got != internal.OpStrDelCode {
		assert.Equal(t, internal.OpStrDelCode, got, "Code()")
	}
}
