package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMore_Basic(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		value         float64
		expectError   bool
		expectedError error
	}{
		{
			name:        "greater_than_success",
			doc:         map[string]any{"score": 85.5},
			path:        []string{"score"},
			value:       80.0,
			expectError: false,
		},
		{
			name:        "greater_than_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       30.0,
			expectError: true,
		},
		{
			name:        "equal_failure",
			doc:         map[string]any{"score": 25.0},
			path:        []string{"score"},
			value:       25.0,
			expectError: true,
		},
		{
			name:        "integer_comparison",
			doc:         map[string]any{"age": 30},
			path:        []string{"age"},
			value:       25.0,
			expectError: false,
		},
		{
			name:          "non_numeric_value",
			doc:           map[string]any{"name": "John"},
			path:          []string{"name"},
			value:         10.0,
			expectError:   true,
			expectedError: ErrNotNumber,
		},
		{
			name:          "missing_path",
			doc:           map[string]any{"score": 85.5},
			path:          []string{"missing"},
			value:         80.0,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
		{
			name: "nested_path",
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
			moreOp := NewMore(tt.path, tt.value)
			result, err := moreOp.Apply(tt.doc)

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

func TestMore_Constructor(t *testing.T) {
	path := []string{"age"}
	value := 18.0

	moreOp := NewMore(path, value)
	assert.Equal(t, path, moreOp.Path())
	assert.Equal(t, value, moreOp.Value)
	assert.Equal(t, internal.OpMoreType, moreOp.Op())
	assert.Equal(t, internal.OpMoreCode, moreOp.Code())
}

func TestMore_ToJSON(t *testing.T) {
	moreOp := NewMore([]string{"age"}, 18.0)
	got, err := moreOp.ToJSON()

	require.NoError(t, err)
	assert.Equal(t, string(internal.OpMoreType), got.Op)
	assert.Equal(t, "/age", got.Path)
	assert.Equal(t, 18, got.Value) // Expect int, not float64
}

func TestMore_ToCompact(t *testing.T) {
	moreOp := NewMore([]string{"age"}, 18.0)
	compact, err := moreOp.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []any{internal.OpMoreCode, []string{"age"}, 18.0}, compact)
}
