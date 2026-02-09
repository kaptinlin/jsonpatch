package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestTestStringLen_Apply(t *testing.T) {
	tests := []struct {
		name           string
		doc            any
		path           []string
		expectedLength float64
		expectError    bool
		expectedError  error
	}{
		{
			name:           "string length match success",
			doc:            map[string]any{"name": "John"},
			path:           []string{"name"},
			expectedLength: 4.0,
			expectError:    false,
		},
		{
			name:           "empty string length",
			doc:            map[string]any{"description": ""},
			path:           []string{"description"},
			expectedLength: 0.0,
			expectError:    false,
		},
		{
			name:           "long string length",
			doc:            map[string]any{"text": "Hello, World! 123"},
			path:           []string{"text"},
			expectedLength: 17.0,
			expectError:    false,
		},
		{
			name:           "unicode string length",
			doc:            map[string]any{"text": "你好世界"},
			path:           []string{"text"},
			expectedLength: 12.0, // 4 Chinese characters = 12 bytes in UTF-8
			expectError:    false,
		},
		{
			name:           "string length mismatch",
			doc:            map[string]any{"name": "John"},
			path:           []string{"name"},
			expectedLength: 5.0,
			expectError:    true,
			expectedError:  ErrStringLengthMismatch,
		},
		{
			name:           "non-string value",
			doc:            map[string]any{"age": 25},
			path:           []string{"age"},
			expectedLength: 2.0,
			expectError:    true,
			expectedError:  ErrNotString,
		},
		{
			name:           "null value",
			doc:            map[string]any{"value": nil},
			path:           []string{"value"},
			expectedLength: 0.0,
			expectError:    true,
			expectedError:  ErrNotString,
		},
		{
			name:           "path not found",
			doc:            map[string]any{"name": "John"},
			path:           []string{"nonexistent"},
			expectedLength: 4.0,
			expectError:    true,
			expectedError:  ErrPathNotFound,
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
			path:           []string{"user", "profile", "email"},
			expectedLength: 16.0,
			expectError:    false,
		},
		{
			name: "array index success",
			doc: map[string]any{
				"items": []any{"item1", "item2", "item3"},
			},
			path:           []string{"items", "1"},
			expectedLength: 5.0,
			expectError:    false,
		},
		{
			name:           "byte slice as string",
			doc:            map[string]any{"data": []byte("hello")},
			path:           []string{"data"},
			expectedLength: 5.0,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strLenOp := NewTestStringLen(tt.path, tt.expectedLength)
			result, err := strLenOp.Apply(tt.doc)

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
