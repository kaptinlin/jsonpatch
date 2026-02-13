package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestIn_Apply(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			inOp := NewIn(tt.path, tt.values)
			result, err := inOp.Apply(tt.doc)

			if tt.expectError {
				if err == nil {
					assert.Fail(t, "Apply() succeeded, want error")
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					assert.Equal(t, tt.expectedError, err, "Apply() error")
				}
				assert.Equal(t, internal.OpResult[any]{}, result)
			} else {
				if err != nil {
					t.Errorf("Apply() failed: %v", err)
				}
				if result.Doc == nil {
					assert.Fail(t, "Apply() result.Doc = nil, want non-nil")
				}
				assert.Equal(t, tt.doc, result.Doc)
			}
		})
	}
}
