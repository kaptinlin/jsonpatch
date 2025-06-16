package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpEnds_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           interface{}
		path          []string
		suffix        string
		expectError   bool
		expectedError error
	}{
		{
			name:        "test ends with suffix success",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "World!",
			expectError: false,
		},
		{
			name:        "test ends with exact match",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "Hello, World!",
			expectError: false,
		},
		{
			name:        "test ends with empty string",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "",
			expectError: false,
		},
		{
			name:        "test ends case sensitive",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "world!",
			expectError: true,
		},
		{
			name:        "test ends suffix not found",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			suffix:      "Hello",
			expectError: true,
		},
		{
			name:          "test non-string value",
			doc:           map[string]interface{}{"age": 25},
			path:          []string{"age"},
			suffix:        "5",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test null value",
			doc:           map[string]interface{}{"value": nil},
			path:          []string{"value"},
			suffix:        "test",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test path not found",
			doc:           map[string]interface{}{"text": "Hello, World!"},
			path:          []string{"nonexistent"},
			suffix:        "World!",
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "test nested path success",
			doc: map[string]interface{}{
				"user": map[string]interface{}{
					"profile": map[string]interface{}{
						"email": "john@example.com",
					},
				},
			},
			path:        []string{"user", "profile", "email"},
			suffix:      ".com",
			expectError: false,
		},
		{
			name: "test array index success",
			doc: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			path:        []string{"items", "2"},
			suffix:      "3",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOpEndsOperation(tt.path, tt.suffix)
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
