package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpIn_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           interface{}
		path          []string
		values        []interface{}
		expectError   bool
		expectedError error
	}{
		{
			name:        "test value in array success",
			doc:         map[string]interface{}{"status": "active"},
			path:        []string{"status"},
			values:      []interface{}{"active", "inactive", "pending"},
			expectError: false,
		},
		{
			name:        "test value not in array",
			doc:         map[string]interface{}{"status": "deleted"},
			path:        []string{"status"},
			values:      []interface{}{"active", "inactive", "pending"},
			expectError: true,
		},
		{
			name:        "test number in array",
			doc:         map[string]interface{}{"priority": 1},
			path:        []string{"priority"},
			values:      []interface{}{1, 2, 3},
			expectError: false,
		},
		{
			name:        "test boolean in array",
			doc:         map[string]interface{}{"enabled": true},
			path:        []string{"enabled"},
			values:      []interface{}{true, false},
			expectError: false,
		},
		{
			name:        "test null value in array",
			doc:         map[string]interface{}{"value": nil},
			path:        []string{"value"},
			values:      []interface{}{nil, "test"},
			expectError: false,
		},
		{
			name:          "test path not found",
			doc:           map[string]interface{}{"status": "active"},
			path:          []string{"nonexistent"},
			values:        []interface{}{"active", "inactive"},
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "test nested path success",
			doc: map[string]interface{}{
				"user": map[string]interface{}{
					"role": "admin",
				},
			},
			path:        []string{"user", "role"},
			values:      []interface{}{"admin", "user", "guest"},
			expectError: false,
		},
		{
			name: "test array index success",
			doc: map[string]interface{}{
				"items": []interface{}{"apple", "banana", "cherry"},
			},
			path:        []string{"items", "1"},
			values:      []interface{}{"banana", "orange"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOpInOperation(tt.path, tt.values)
			result, err := op.Apply(tt.doc)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.True(t, errors.Is(err, tt.expectedError))
				}
				// Check that result is empty when error occurs
				assert.Equal(t, internal.OpResult{}, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.doc, result.Doc)
			}
		})
	}
}
