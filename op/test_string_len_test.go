package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpTestStringLen_Apply(t *testing.T) {
	tests := []struct {
		name           string
		doc            interface{}
		path           []string
		expectedLength int
		expectError    bool
		expectedError  error
	}{
		{
			name:           "test string length match success",
			doc:            map[string]interface{}{"name": "John"},
			path:           []string{"name"},
			expectedLength: 4,
			expectError:    false,
		},
		{
			name:           "test empty string length",
			doc:            map[string]interface{}{"description": ""},
			path:           []string{"description"},
			expectedLength: 0,
			expectError:    false,
		},
		{
			name:           "test long string length",
			doc:            map[string]interface{}{"text": "Hello, World! 123"},
			path:           []string{"text"},
			expectedLength: 17,
			expectError:    false,
		},
		{
			name:           "test unicode string length",
			doc:            map[string]interface{}{"text": "你好世界"},
			path:           []string{"text"},
			expectedLength: 12, // 4 Chinese characters = 12 bytes in UTF-8
			expectError:    false,
		},
		{
			name:           "test string length mismatch",
			doc:            map[string]interface{}{"name": "John"},
			path:           []string{"name"},
			expectedLength: 5,
			expectError:    true,
			expectedError:  ErrStringLengthMismatch,
		},
		{
			name:           "test non-string value",
			doc:            map[string]interface{}{"age": 25},
			path:           []string{"age"},
			expectedLength: 2,
			expectError:    true,
			expectedError:  ErrNotString,
		},
		{
			name:           "test null value",
			doc:            map[string]interface{}{"value": nil},
			path:           []string{"value"},
			expectedLength: 0,
			expectError:    true,
			expectedError:  ErrNotString,
		},
		{
			name:           "test path not found",
			doc:            map[string]interface{}{"name": "John"},
			path:           []string{"nonexistent"},
			expectedLength: 4,
			expectError:    true,
			expectedError:  ErrPathNotFound,
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
			path:           []string{"user", "profile", "email"},
			expectedLength: 16,
			expectError:    false,
		},
		{
			name: "test array index success",
			doc: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			path:           []string{"items", "1"},
			expectedLength: 5,
			expectError:    false,
		},
		{
			name:           "test byte slice as string",
			doc:            map[string]interface{}{"data": []byte("hello")},
			path:           []string{"data"},
			expectedLength: 5,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOpTestStringLenOperation(tt.path, tt.expectedLength)
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
