package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestEnds_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		suffix        string
		expectError   bool
		expectedError error
	}{
		{
			name:        "suffix success",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "World!",
			expectError: false,
		},
		{
			name:        "exact match",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "Hello, World!",
			expectError: false,
		},
		{
			name:        "empty string",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "",
			expectError: false,
		},
		{
			name:        "case sensitive",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "world!",
			expectError: true,
		},
		{
			name:        "suffix not found",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "Hello",
			expectError: true,
		},
		{
			name:          "non-string value",
			doc:           map[string]any{"age": 25},
			path:          []string{"age"},
			suffix:        "5",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "null value",
			doc:           map[string]any{"value": nil},
			path:          []string{"value"},
			suffix:        "test",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "path not found",
			doc:           map[string]any{"text": "Hello, World!"},
			path:          []string{"nonexistent"},
			suffix:        "World!",
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
			suffix:      ".com",
			expectError: false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:        []string{"items", "2"},
			suffix:      "3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endsOp := NewEnds(tt.path, tt.suffix)
			result, err := endsOp.Apply(tt.doc)

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
