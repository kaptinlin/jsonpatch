package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestExtend_Apply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		path       []string
		doc        any
		props      any
		deleteNull bool
		expected   any
		oldValue   any
		wantErr    bool
	}{
		{
			name:       "add new properties",
			path:       []string{"user"},
			doc:        map[string]any{"user": map[string]any{"name": "John"}},
			props:      map[string]any{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]any{"user": map[string]any{"name": "John", "age": 30, "city": "NYC"}},
			oldValue:   map[string]any{"name": "John"},
		},
		{
			name:       "update existing properties",
			path:       []string{"user"},
			doc:        map[string]any{"user": map[string]any{"name": "John", "age": 25}},
			props:      map[string]any{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]any{"user": map[string]any{"name": "John", "age": 30, "city": "NYC"}},
			oldValue:   map[string]any{"name": "John", "age": 25},
		},
		{
			name:       "delete null properties",
			path:       []string{"user"},
			doc:        map[string]any{"user": map[string]any{"name": "John", "age": 25, "city": "NYC"}},
			props:      map[string]any{"age": nil, "city": nil},
			deleteNull: true,
			expected:   map[string]any{"user": map[string]any{"name": "John"}},
			oldValue:   map[string]any{"name": "John", "age": 25, "city": "NYC"},
		},
		{
			name:       "keep null properties when deleteNull is false",
			path:       []string{"user"},
			doc:        map[string]any{"user": map[string]any{"name": "John", "age": 25}},
			props:      map[string]any{"age": nil, "city": nil},
			deleteNull: false,
			expected:   map[string]any{"user": map[string]any{"name": "John", "age": nil, "city": nil}},
			oldValue:   map[string]any{"name": "John", "age": 25},
		},
		{
			name:       "extend at root",
			path:       []string{},
			doc:        map[string]any{"name": "John"},
			props:      map[string]any{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]any{"name": "John", "age": 30, "city": "NYC"},
			oldValue:   map[string]any{"name": "John"},
		},
		{
			name:       "extend nested object",
			path:       []string{"user", "profile"},
			doc:        map[string]any{"user": map[string]any{"profile": map[string]any{"name": "John"}}},
			props:      map[string]any{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]any{"user": map[string]any{"profile": map[string]any{"name": "John", "age": 30, "city": "NYC"}}},
			oldValue:   map[string]any{"name": "John"},
		},
		{
			name:       "extend with complex properties",
			path:       []string{"config"},
			doc:        map[string]any{"config": map[string]any{"name": "app"}},
			props:      map[string]any{"settings": map[string]any{"theme": "dark"}, "enabled": true, "count": 42},
			deleteNull: false,
			expected:   map[string]any{"config": map[string]any{"name": "app", "settings": map[string]any{"theme": "dark"}, "enabled": true, "count": 42}},
			oldValue:   map[string]any{"name": "app"},
		},
		{
			name:       "path not found",
			path:       []string{"notfound"},
			doc:        map[string]any{"user": map[string]any{"name": "John"}},
			props:      map[string]any{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
		{
			name:       "not an object",
			path:       []string{"text"},
			doc:        map[string]any{"text": "abc"},
			props:      map[string]any{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
		{
			name:       "root not an object",
			path:       []string{},
			doc:        "abc",
			props:      map[string]any{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			extendOp := NewExtend(tt.path, tt.props.(map[string]any), tt.deleteNull)
			docCopy, err := DeepClone(tt.doc)
			if err != nil {
				t.Fatalf("DeepClone() error: %v", err)
			}

			result, err := extendOp.Apply(docCopy)

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

func TestExtend_Constructor(t *testing.T) {
	t.Parallel()
	path := []string{"user", "profile"}
	props := map[string]any{"age": 30, "city": "NYC"}
	deleteNull := true
	extendOp := NewExtend(path, props, deleteNull)
	if diff := cmp.Diff(path, extendOp.Path()); diff != "" {
		t.Errorf("NewExtend() Path mismatch (-want +got):\n%s", diff)
	}
	assert.Equal(t, props, extendOp.Properties)
	assert.Equal(t, deleteNull, extendOp.DeleteNull, "NewExtend() DeleteNull")
	if got := extendOp.Op(); got != internal.OpExtendType {
		assert.Equal(t, internal.OpExtendType, got, "Op()")
	}
	if got := extendOp.Code(); got != internal.OpExtendCode {
		assert.Equal(t, internal.OpExtendCode, got, "Code()")
	}
}
