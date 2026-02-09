package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpMore_Basic(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		value         float64
		expectError   bool
		expectedError error
	}{
		{
			name:        "test_greater_than_success",
			doc:         map[string]any{"score": 85.5},
			path:        []string{"score"},
			value:       80.0,
			expectError: false,
		},
		{
			name:        "test_greater_than_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       30.0,
			expectError: true,
		},
		{
			name:        "test_equal_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       25.0,
			expectError: true,
		},
		{
			name:        "test_integer_comparison",
			doc:         map[string]any{"age": 30},
			path:        []string{"age"},
			value:       25.0,
			expectError: false,
		},
		{
			name:          "test_non_numeric_value",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			value:         10.0,
			expectError:   true,
			expectedError: ErrNotNumber,
		},
		{
			name:          "test_missing_path",
			doc:           map[string]any{"score": 85.5},
			path:          []string{"missing"},
			value:         80.0,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "test_nested_path",
			doc: map[string]any{
				"user": map[string]any{
					"stats": map[string]any{
						"score": 95.0,
					},
				},
			},
			path:        []string{"user", "stats", "score"},
			value:       90.0,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewMore(tt.path, tt.value)
			result, err := op.Apply(tt.doc)

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

func TestOpMore_Constructor(t *testing.T) {
	path := []string{"age"}
	value := 18.0

	op := NewMore(path, value)
	assert.Equal(t, path, op.Path())
	assert.Equal(t, value, op.Value)
	assert.Equal(t, internal.OpMoreType, op.Op())
	assert.Equal(t, internal.OpMoreCode, op.Code())
}

func TestOpMore_ToJSON(t *testing.T) {
	op := NewMore([]string{"age"}, 18.0)
	json, err := op.ToJSON()

	require.NoError(t, err)
	assert.Equal(t, string(internal.OpMoreType), json.Op)
	assert.Equal(t, "/age", json.Path)
	assert.Equal(t, 18, json.Value) // Expect int, not float64
}

func TestOpMore_ToCompact(t *testing.T) {
	op := NewMore([]string{"age"}, 18.0)
	compact, err := op.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []any{internal.OpMoreCode, []string{"age"}, 18.0}, compact)
}
