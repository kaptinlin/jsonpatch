package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
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
			testStringOp := NewTestString(tt.path, tt.expectedValue, tt.pos, false, false)
			result, err := testStringOp.Apply(tt.doc)

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
				if err == nil {
					t.Error("toString() succeeded, want error")
				}
			} else {
				if err != nil {
					t.Errorf("toString() failed: %v", err)
				}
				if result != tt.expected {
					t.Errorf("toString() = %q, want %q", result, tt.expected)
				}
			}
		})
	}
}
