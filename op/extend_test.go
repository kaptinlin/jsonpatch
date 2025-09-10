package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpExtend_Apply(t *testing.T) {
	tests := []struct {
		name       string
		path       []string
		doc        interface{}
		props      interface{}
		deleteNull bool
		expected   interface{}
		oldValue   interface{}
		wantErr    bool
	}{
		{
			name:       "add new properties",
			path:       []string{"user"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"name": "John"}},
			props:      map[string]interface{}{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": 30, "city": "NYC"}},
			oldValue:   map[string]interface{}{"name": "John"},
		},
		{
			name:       "update existing properties",
			path:       []string{"user"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": 25}},
			props:      map[string]interface{}{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": 30, "city": "NYC"}},
			oldValue:   map[string]interface{}{"name": "John", "age": 25},
		},
		{
			name:       "delete null properties",
			path:       []string{"user"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": 25, "city": "NYC"}},
			props:      map[string]interface{}{"age": nil, "city": nil},
			deleteNull: true,
			expected:   map[string]interface{}{"user": map[string]interface{}{"name": "John"}},
			oldValue:   map[string]interface{}{"name": "John", "age": 25, "city": "NYC"},
		},
		{
			name:       "keep null properties when deleteNull is false",
			path:       []string{"user"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": 25}},
			props:      map[string]interface{}{"age": nil, "city": nil},
			deleteNull: false,
			expected:   map[string]interface{}{"user": map[string]interface{}{"name": "John", "age": nil, "city": nil}},
			oldValue:   map[string]interface{}{"name": "John", "age": 25},
		},
		{
			name:       "extend at root",
			path:       []string{},
			doc:        map[string]interface{}{"name": "John"},
			props:      map[string]interface{}{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]interface{}{"name": "John", "age": 30, "city": "NYC"},
			oldValue:   map[string]interface{}{"name": "John"},
		},
		{
			name:       "extend nested object",
			path:       []string{"user", "profile"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"profile": map[string]interface{}{"name": "John"}}},
			props:      map[string]interface{}{"age": 30, "city": "NYC"},
			deleteNull: false,
			expected:   map[string]interface{}{"user": map[string]interface{}{"profile": map[string]interface{}{"name": "John", "age": 30, "city": "NYC"}}},
			oldValue:   map[string]interface{}{"name": "John"},
		},
		{
			name:       "extend with complex properties",
			path:       []string{"config"},
			doc:        map[string]interface{}{"config": map[string]interface{}{"name": "app"}},
			props:      map[string]interface{}{"settings": map[string]interface{}{"theme": "dark"}, "enabled": true, "count": 42},
			deleteNull: false,
			expected:   map[string]interface{}{"config": map[string]interface{}{"name": "app", "settings": map[string]interface{}{"theme": "dark"}, "enabled": true, "count": 42}},
			oldValue:   map[string]interface{}{"name": "app"},
		},
		{
			name:       "path not found",
			path:       []string{"notfound"},
			doc:        map[string]interface{}{"user": map[string]interface{}{"name": "John"}},
			props:      map[string]interface{}{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
		{
			name:       "not an object",
			path:       []string{"text"},
			doc:        map[string]interface{}{"text": "abc"},
			props:      map[string]interface{}{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
		{
			name:       "root not an object",
			path:       []string{},
			doc:        "abc",
			props:      map[string]interface{}{"age": 30},
			deleteNull: false,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewExtend(tt.path, tt.props.(map[string]interface{}), tt.deleteNull)
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

func TestOpExtend_Op(t *testing.T) {
	op := NewExtend([]string{"user"}, map[string]interface{}{"age": 30}, false)
	assert.Equal(t, internal.OpExtendType, op.Op())
}

func TestOpExtend_Code(t *testing.T) {
	op := NewExtend([]string{"user"}, map[string]interface{}{"age": 30}, false)
	assert.Equal(t, internal.OpExtendCode, op.Code())
}

func TestOpExtend_NewOpExtend(t *testing.T) {
	path := []string{"user", "profile"}
	props := map[string]interface{}{"age": 30, "city": "NYC"}
	deleteNull := true
	op := NewExtend(path, props, deleteNull)
	assert.Equal(t, path, op.Path())
	assert.Equal(t, props, op.Properties)
	assert.Equal(t, deleteNull, op.DeleteNull)
	assert.Equal(t, internal.OpExtendType, op.Op())
	assert.Equal(t, internal.OpExtendCode, op.Code())
}
