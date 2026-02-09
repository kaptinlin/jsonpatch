package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpTestType_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		expectedType  string
		expectError   bool
		expectedError error
	}{
		{
			name:         "test string type success",
			doc:          map[string]any{"name": "John"},
			path:         []string{"name"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "test number type success",
			doc:          map[string]any{"age": 25.0},
			path:         []string{"age"},
			expectedType: "number",
			expectError:  false,
		},
		{
			name:         "test boolean type success",
			doc:          map[string]any{"active": true},
			path:         []string{"active"},
			expectedType: "boolean",
			expectError:  false,
		},
		{
			name:         "test array type success",
			doc:          map[string]any{"tags": []any{"tag1", "tag2"}},
			path:         []string{"tags"},
			expectedType: "array",
			expectError:  false,
		},
		{
			name:         "test object type success",
			doc:          map[string]any{"address": map[string]any{"city": "NYC"}},
			path:         []string{"address"},
			expectedType: "object",
			expectError:  false,
		},
		{
			name:         "test null type success",
			doc:          map[string]any{"value": nil},
			path:         []string{"value"},
			expectedType: "null",
			expectError:  false,
		},
		{
			name:          "test type mismatch",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			expectedType:  "number",
			expectError:   true,
			expectedError: ErrTypeMismatch,
		},
		{
			name:          "test path not found",
			doc:           map[string]any{"name": "John"},
			path:          []string{"nonexistent"},
			expectedType:  "string",
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "test nested path success",
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
			name: "test array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:         []string{"items", "1"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "test integer as number",
			doc:          map[string]any{"count": 42},
			path:         []string{"count"},
			expectedType: "number",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewTestType(tt.path, tt.expectedType)
			result, err := op.Apply(tt.doc)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.True(t, errors.Is(err, tt.expectedError), "Expected error %v, got %v", tt.expectedError, err)
				}
				// Check that result is empty when error occurs
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
