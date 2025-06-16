package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestOpTestType_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           interface{}
		path          []string
		expectedType  string
		expectError   bool
		expectedError error
	}{
		{
			name:         "test string type success",
			doc:          map[string]interface{}{"name": "John"},
			path:         []string{"name"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "test number type success",
			doc:          map[string]interface{}{"age": 25.0},
			path:         []string{"age"},
			expectedType: "number",
			expectError:  false,
		},
		{
			name:         "test boolean type success",
			doc:          map[string]interface{}{"active": true},
			path:         []string{"active"},
			expectedType: "boolean",
			expectError:  false,
		},
		{
			name:         "test array type success",
			doc:          map[string]interface{}{"tags": []interface{}{"tag1", "tag2"}},
			path:         []string{"tags"},
			expectedType: "array",
			expectError:  false,
		},
		{
			name:         "test object type success",
			doc:          map[string]interface{}{"address": map[string]interface{}{"city": "NYC"}},
			path:         []string{"address"},
			expectedType: "object",
			expectError:  false,
		},
		{
			name:         "test null type success",
			doc:          map[string]interface{}{"value": nil},
			path:         []string{"value"},
			expectedType: "null",
			expectError:  false,
		},
		{
			name:          "test type mismatch",
			doc:           map[string]interface{}{"name": "John"},
			path:          []string{"name"},
			expectedType:  "number",
			expectError:   true,
			expectedError: ErrTypeMismatch,
		},
		{
			name:          "test path not found",
			doc:           map[string]interface{}{"name": "John"},
			path:          []string{"nonexistent"},
			expectedType:  "string",
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
			path:         []string{"user", "profile", "email"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name: "test array index success",
			doc: map[string]interface{}{
				"items": []interface{}{"item1", "item2", "item3"},
			},
			path:         []string{"items", "1"},
			expectedType: "string",
			expectError:  false,
		},
		{
			name:         "test integer as number",
			doc:          map[string]interface{}{"count": 42},
			path:         []string{"count"},
			expectedType: "number",
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewOpTestTypeOperation(tt.path, tt.expectedType)
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

func TestGetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"null", nil, "null"},
		{"string", "hello", "string"},
		{"number float64", 3.14, "number"},
		{"number int", 42, "number"},
		{"boolean", true, "boolean"},
		{"array", []interface{}{1, 2, 3}, "array"},
		{"object", map[string]interface{}{"key": "value"}, "object"},
		{"int32", int32(42), "number"},
		{"uint64", uint64(42), "number"},
		{"float32", float32(3.14), "number"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getTypeName(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}
