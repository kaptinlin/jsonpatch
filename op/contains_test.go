package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpContains_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           interface{}
		path          []string
		substring     string
		expectError   bool
		expectedError error
	}{
		{
			name:        "test contains substring success",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "World",
			expectError: false,
		},
		{
			name:        "test contains at beginning",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Hello",
			expectError: false,
		},
		{
			name:        "test contains at end",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "World!",
			expectError: false,
		},
		{
			name:        "test contains empty string",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "",
			expectError: false,
		},
		{
			name:        "test contains exact match",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Hello, World!",
			expectError: false,
		},
		{
			name:        "test contains case sensitive",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "hello",
			expectError: true,
		},
		{
			name:        "test contains substring not found",
			doc:         map[string]interface{}{"text": "Hello, World!"},
			path:        []string{"text"},
			substring:   "Python",
			expectError: true,
		},
		{
			name:          "test non-string value",
			doc:           map[string]interface{}{"age": 25},
			path:          []string{"age"},
			substring:     "25",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test null value",
			doc:           map[string]interface{}{"value": nil},
			path:          []string{"value"},
			substring:     "test",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test path not found",
			doc:           map[string]interface{}{"text": "Hello, World!"},
			path:          []string{"nonexistent"},
			substring:     "Hello",
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
			substring:   "@example.com",
			expectError: false,
		},
		{
			name: "test array index success",
			doc: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			path:        []string{"items", "1"},
			substring:   "item",
			expectError: false,
		},
		{
			name:        "test byte slice as string",
			doc:         map[string]interface{}{"data": []byte("hello world")},
			path:        []string{"data"},
			substring:   "world",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewContains(tt.path, tt.substring)
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
