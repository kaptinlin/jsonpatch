package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpStarts_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		prefix        string
		expectError   bool
		expectedError error
	}{
		{
			name:        "test starts with prefix success",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			prefix:      "Hello",
			expectError: false,
		},
		{
			name:        "test starts with exact match",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			prefix:      "Hello, World!",
			expectError: false,
		},
		{
			name:        "test starts with empty string",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			prefix:      "",
			expectError: false,
		},
		{
			name:        "test starts case sensitive",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			prefix:      "hello",
			expectError: true,
		},
		{
			name:        "test starts prefix not found",
			doc:         map[string]any{"text": "Hello, World!"},
			path:        []string{"text"},
			prefix:      "World",
			expectError: true,
		},
		{
			name:          "test non-string value",
			doc:           map[string]any{"age": 25},
			path:          []string{"age"},
			prefix:        "2",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test null value",
			doc:           map[string]any{"value": nil},
			path:          []string{"value"},
			prefix:        "test",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test path not found",
			doc:           map[string]any{"text": "Hello, World!"},
			path:          []string{"nonexistent"},
			prefix:        "Hello",
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
			path:        []string{"user", "profile", "email"},
			prefix:      "john",
			expectError: false,
		},
		{
			name: "test array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:        []string{"items", "1"},
			prefix:      "item",
			expectError: false,
		},
		{
			name:        "test byte slice as string",
			doc:         map[string]any{"data": []byte("hello world")},
			path:        []string{"data"},
			prefix:      "hello",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewStarts(tt.path, tt.prefix)
			result, err := op.Apply(tt.doc)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.True(t, errors.Is(err, tt.expectedError))
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
