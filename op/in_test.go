package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
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
					t.Error("Apply() succeeded, want error")
				}
				if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
					t.Errorf("Apply() error = %v, want %v", err, tt.expectedError)
				}
				if diff := cmp.Diff(internal.OpResult[any]{}, result); diff != "" {
					t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
				}
			} else {
				if err != nil {
					t.Errorf("Apply() failed: %v", err)
				}
				if result.Doc == nil {
					t.Error("Apply() result.Doc = nil, want non-nil")
				}
				if diff := cmp.Diff(tt.doc, result.Doc); diff != "" {
					t.Errorf("Apply() result.Doc mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
