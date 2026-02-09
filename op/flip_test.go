package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpFlip_Apply(t *testing.T) {
	tests := []struct {
		name     string
		path     []string
		doc      any
		expected any
		oldValue any
		wantErr  bool
	}{
		{
			name:     "flip boolean true to false",
			path:     []string{"flag"},
			doc:      map[string]any{"flag": true},
			expected: map[string]any{"flag": false},
			oldValue: true,
		},
		{
			name:     "flip boolean false to true",
			path:     []string{"flag"},
			doc:      map[string]any{"flag": false},
			expected: map[string]any{"flag": true},
			oldValue: false,
		},
		{
			name:     "flip number 0 to true",
			path:     []string{"count"},
			doc:      map[string]any{"count": 0},
			expected: map[string]any{"count": true},
			oldValue: 0,
		},
		{
			name:     "flip number 5 to false",
			path:     []string{"count"},
			doc:      map[string]any{"count": 5},
			expected: map[string]any{"count": false},
			oldValue: 5,
		},
		{
			name:     "flip empty string to true",
			path:     []string{"text"},
			doc:      map[string]any{"text": ""},
			expected: map[string]any{"text": true},
			oldValue: "",
		},
		{
			name:     "flip non-empty string to false",
			path:     []string{"text"},
			doc:      map[string]any{"text": "hello"},
			expected: map[string]any{"text": false},
			oldValue: "hello",
		},
		{
			name:     "flip nil to true",
			path:     []string{"value"},
			doc:      map[string]any{"value": nil},
			expected: map[string]any{"value": true},
			oldValue: nil,
		},
		{
			name:     "flip empty array to false",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{}},
			expected: map[string]any{"items": false},
			oldValue: []any{},
		},
		{
			name:     "flip non-empty array to false",
			path:     []string{"items"},
			doc:      map[string]any{"items": []any{1, 2, 3}},
			expected: map[string]any{"items": false},
			oldValue: []any{1, 2, 3},
		},
		{
			name:     "flip empty map to false",
			path:     []string{"config"},
			doc:      map[string]any{"config": map[string]any{}},
			expected: map[string]any{"config": false},
			oldValue: map[string]any{},
		},
		{
			name:     "flip non-empty map to false",
			path:     []string{"config"},
			doc:      map[string]any{"config": map[string]any{"key": "value"}},
			expected: map[string]any{"config": false},
			oldValue: map[string]any{"key": "value"},
		},
		{
			name:     "flip nested path",
			path:     []string{"user", "active"},
			doc:      map[string]any{"user": map[string]any{"active": true}},
			expected: map[string]any{"user": map[string]any{"active": false}},
			oldValue: true,
		},
		{
			name:     "flip array element",
			path:     []string{"flags", "0"},
			doc:      map[string]any{"flags": []any{true, false, true}},
			expected: map[string]any{"flags": []any{false, false, true}},
			oldValue: true,
		},
		{
			name:     "flip root level boolean",
			path:     []string{},
			doc:      true,
			expected: false,
			oldValue: true,
		},
		{
			name:     "flip root level number",
			path:     []string{},
			doc:      42,
			expected: false,
			oldValue: 42,
		},
		{
			name:     "path not found creates true",
			path:     []string{"nonexistent"},
			doc:      map[string]any{"flag": true},
			expected: map[string]any{"flag": true, "nonexistent": true},
			oldValue: nil,
		},
		{
			name:    "invalid path for array",
			path:    []string{"items", "invalid"},
			doc:     map[string]any{"items": []any{1, 2, 3}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewFlip(tt.path)

			// Deep clone the document to avoid modifying the original
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

func TestOpFlip_Op(t *testing.T) {
	op := NewFlip([]string{"flag"})
	assert.Equal(t, internal.OpFlipType, op.Op())
}

func TestOpFlip_Code(t *testing.T) {
	op := NewFlip([]string{"flag"})
	assert.Equal(t, internal.OpFlipCode, op.Code())
}

func TestOpFlip_NewOpFlip(t *testing.T) {
	path := []string{"user", "active"}
	op := NewFlip(path)

	assert.Equal(t, path, op.Path())
	assert.Equal(t, internal.OpFlipType, op.Op())
	assert.Equal(t, internal.OpFlipCode, op.Code())
}

func TestOpFlip_ComplexTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected bool
	}{
		{"float64 zero", 0.0, true},
		{"float64 non-zero", 3.14, false},
		{"int8 zero", int8(0), true},
		{"int8 non-zero", int8(1), false},
		{"uint zero", uint(0), true},
		{"uint non-zero", uint(1), false},
		{"float32 zero", float32(0.0), true},
		{"float32 non-zero", float32(1.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewFlip([]string{"value"})
			doc := map[string]any{"value": tt.value}

			result, err := op.Apply(doc)
			require.NoError(t, err)

			// Check that the result is the expected boolean
			resultDoc := result.Doc.(map[string]any)
			assert.Equal(t, tt.expected, resultDoc["value"])
		})
	}
}
