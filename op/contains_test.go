package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestContains_Apply(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		doc           any
		path          []string
		substring     string
		expectError   bool
		expectedError error
	}{
		{
			name:        "substring success",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "World",
			expectError: false,
		},
		{
			name:        "at beginning",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Hello",
			expectError: false,
		},
		{
			name:        "at end",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "World!",
			expectError: false,
		},
		{
			name:        "empty string",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "",
			expectError: false,
		},
		{
			name:        "exact match",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Hello, World!",
			expectError: false,
		},
		{
			name:        "case sensitive",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "hello",
			expectError: true,
		},
		{
			name:        "substring not found",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Python",
			expectError: true,
		},
		{
			name:          "non-string value",
			doc:           map[string]any{"age": 25},
			path:          []string{"age"},
			substring:     "25",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "null value",
			doc:           map[string]any{"value": nil},
			path:          []string{"value"},
			substring:     "test",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "path not found",
			doc:           map[string]any{"text": "Hello, World!"},
			path:          []string{"nonexistent"},
			substring:     "Hello",
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
			path:        []string{"user", "profile", "email"},
			substring:   "@example.com",
			expectError: false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:        []string{"items", "1"},
			substring:   "item",
			expectError: false,
		},
		{
			name:        "byte slice as string",
			doc:         map[string]any{"data": []byte("hello world")},
			path:        []string{"data"},
			substring:   "world",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			containsOp := NewContains(tt.path, tt.substring)
			result, err := containsOp.Apply(tt.doc)

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
