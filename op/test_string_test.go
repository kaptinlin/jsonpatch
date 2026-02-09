package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
)

func TestTestString_Apply(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		pos           float64
		expectedValue string
		expectError   bool
		expectedError error
	}{
		{
			name:          "string match success",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			pos:           0.0,
			expectedValue: "John",
			expectError:   false,
		},
		{
			name:          "empty string success",
			doc:           map[string]any{"description": ""},
			path:          []string{"description"},
			pos:           0.0,
			expectedValue: "",
			expectError:   false,
		},
		{
			name:          "string with special characters",
			doc:           map[string]any{"text": "Hello, World! 123"},
			path:          []string{"text"},
			pos:           7.0,
			expectedValue: "World",
			expectError:   false,
		},
		{
			name:          "string mismatch",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			pos:           0.0,
			expectedValue: "Jane",
			expectError:   true,
			expectedError: ErrSubstringMismatch,
		},
		{
			name:          "non-string value",
			doc:           map[string]any{"age": 25},
			path:          []string{"age"},
			pos:           0.0,
			expectedValue: "25",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "null value",
			doc:           map[string]any{"value": nil},
			path:          []string{"value"},
			pos:           0.0,
			expectedValue: "",
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "path not found",
			doc:           map[string]any{"name": "John"},
			path:          []string{"nonexistent"},
			pos:           0.0,
			expectedValue: "John",
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
			path:          []string{"user", "profile", "email"},
			pos:           5.0,
			expectedValue: "example",
			expectError:   false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:          []string{"items", "1"},
			pos:           0.0,
			expectedValue: "item2",
			expectError:   false,
		},
		{
			name:          "byte slice as string",
			doc:           map[string]any{"data": []byte("hello")},
			path:          []string{"data"},
			pos:           1.0,
			expectedValue: "ell",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStringOp := NewTestStringWithPos(tt.path, tt.expectedValue, tt.pos)
			result, err := testStringOp.Apply(tt.doc)

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

func TestToString(t *testing.T) {
	tests := []struct {
		name     string
		value    any
		expected string
		hasError bool
	}{
		{"string", "hello", "hello", false},
		{"byte slice", []byte("world"), "world", false},
		{"nil", nil, "", true},
		{"int", 42, "", true},
		{"bool", true, "", true},
		{"float", 3.14, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := toString(tt.value)
			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}
