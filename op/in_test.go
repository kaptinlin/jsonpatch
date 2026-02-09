package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestIn_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		values        []any
		expectError   bool
		expectedError error
	}{
		{
			name:        "value in array success",
			doc:         map[string]any{"status": "active"},
			path:        []string{"status"},
			values:      []any{"active", "inactive", "pending"},
			expectError: false,
		},
		{
			name:        "value not in array",
			doc:         map[string]any{"status": "deleted"},
			path:        []string{"status"},
			values:      []any{"active", "inactive", "pending"},
			expectError: true,
		},
		{
			name:        "number in array",
			doc:         map[string]any{"priority": 1},
			path:        []string{"priority"},
			values:      []any{1, 2, 3},
			expectError: false,
		},
		{
			name:        "boolean in array",
			doc:         map[string]any{"enabled": true},
			path:        []string{"enabled"},
			values:      []any{true, false},
			expectError: false,
		},
		{
			name:        "null value in array",
			doc:         map[string]any{"value": nil},
			path:        []string{"value"},
			values:      []any{nil, "test"},
			expectError: false,
		},
		{
			name:          "path not found",
			doc:           map[string]any{"status": "active"},
			path:          []string{"nonexistent"},
			values:        []any{"active", "inactive"},
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "nested path success",
			doc: map[string]any{
				"user": map[string]any{
					"role": "admin",
				},
			},
			path:        []string{"user", "role"},
			values:      []any{"admin", "user", "guest"},
			expectError: false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"apple", "banana", "cherry"},
			},
			path:        []string{"items", "1"},
			values:      []any{"banana", "orange"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inOp := NewIn(tt.path, tt.values)
			result, err := inOp.Apply(tt.doc)

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
