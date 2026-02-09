package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestTestType_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		expectedType  string
		expectError   bool
		expectedError error
	}{
		{
			name:         "string type success",
			doc:          map[string]any{"name": "John"},
			path:         []string{"name"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "number type success",
			doc:          map[string]any{"age": 25.0},
			path:         []string{"age"},
			expectedType: "number",
			expectError:  false,
		},
		{
			name:         "boolean type success",
			doc:          map[string]any{"active": true},
			path:         []string{"active"},
			expectedType: "boolean",
			expectError:  false,
		},
		{
			name:         "array type success",
			doc:          map[string]any{"tags": []any{"tag1", "tag2"}},
			path:         []string{"tags"},
			expectedType: "array",
			expectError:  false,
		},
		{
			name:         "object type success",
			doc:          map[string]any{"address": map[string]any{"city": "NYC"}},
			path:         []string{"address"},
			expectedType: "object",
			expectError:  false,
		},
		{
			name:         "null type success",
			doc:          map[string]any{"value": nil},
			path:         []string{"value"},
			expectedType: "null",
			expectError:  false,
		},
		{
			name:          "type mismatch",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			expectedType:  "number",
			expectError:   true,
			expectedError: ErrTypeMismatch,
		},
		{
			name:          "path not found",
			doc:           map[string]any{"name": "John"},
			path:          []string{"nonexistent"},
			expectedType:  "string",
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "nested path success",
			doc: map[string]any{
				"user": map[string]any{
					"profile": map[string]any{
						"email": "john@example.com",
					},
				},
			},
			path:         []string{"user", "profile", "email"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:         []string{"items", "1"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "integer as number",
			doc:          map[string]any{"count": 42},
			path:         []string{"count"},
			expectedType: "number",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeOp := NewTestType(tt.path, tt.expectedType)
			result, err := typeOp.Apply(tt.doc)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
				assert.Equal(t, internal.OpResult[any]{}, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.doc, result.Doc)
			}
		})
	}
}

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
	}{
		{"null", nil, "null"},
		{"string", "hello", "string"},
		{"number float64", 3.14, "number"},
		{"number int", 42, "number"},
		{"boolean", true, "boolean"},
		{"array", []any{1, 2, 3}, "array"},
		{"object", map[string]any{"key": "value"}, "object"},
		{"int32", int32(42), "number"},
		{"uint64", uint64(42), "number"},
		{"float32", float32(3.14), "number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTypeName(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
