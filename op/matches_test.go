package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMatches_Basic(t *testing.T) {
	tests := []struct {
		name          string
		doc           any
		path          []string
		pattern       string
		ignoreCase    bool
		expectError   bool
		expectedError error
	}{
		{
			name:        "simple pattern match",
			doc:         map[string]any{"text": "hello123"},
			path:        []string{"text"},
			pattern:     "hello\\d+",
			ignoreCase:  false,
			expectError: false,
		},
		{
			name:        "pattern no match",
			doc:         map[string]any{"text": "hello"},
			path:        []string{"text"},
			pattern:     "\\d+",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "case sensitive match",
			doc:         map[string]any{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "case insensitive match",
			doc:         map[string]any{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  true,
			expectError: false,
		},
		{
			name:          "non string value",
			doc:           map[string]any{"number": 123},
			path:          []string{"number"},
			pattern:       "\\d+",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "missing path",
			doc:           map[string]any{"text": "hello"},
			path:          []string{"missing"},
			pattern:       "hello",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matchesOp := NewMatches(tt.path, tt.pattern, tt.ignoreCase, nil)

			result, err := matchesOp.Apply(tt.doc)

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

func TestMatches_Constructor(t *testing.T) {
	path := []string{"email"}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	ignoreCase := false

	matchesOp := NewMatches(path, pattern, ignoreCase, nil)

	assert.Equal(t, path, matchesOp.Path())
	assert.Equal(t, pattern, matchesOp.Pattern)
	assert.Equal(t, ignoreCase, matchesOp.IgnoreCase)
	assert.Equal(t, internal.OpMatchesType, matchesOp.Op())
	assert.Equal(t, internal.OpMatchesCode, matchesOp.Code())
}

func TestMatches_InvalidPattern(t *testing.T) {
	path := []string{"email"}
	invalidPattern := `[invalid-regex`

	// Invalid patterns create a matcher that always returns false
	// (aligned with json-joy's behavior)
	matchesOp := NewMatches(path, invalidPattern, false, nil)
	assert.NotNil(t, matchesOp)

	result, _ := matchesOp.Test(map[string]any{"email": "test@example.com"})
	assert.False(t, result)
}

func TestMatches_ToJSON(t *testing.T) {
	matchesOp := NewMatches([]string{"email"}, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, true, nil)

	got, err := matchesOp.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, string(internal.OpMatchesType), got.Op)
	assert.Equal(t, "/email", got.Path)
	assert.Equal(t, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, got.Value)
	assert.Equal(t, true, got.IgnoreCase)
}

func TestMatches_ToCompact(t *testing.T) {
	matchesOp := NewMatches([]string{"name"}, "john", true, nil)
	compact, err := matchesOp.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []any{internal.OpMatchesCode, []string{"name"}, "john", true}, compact)
}

func TestMatches_ToCompact_WithoutIgnoreCase(t *testing.T) {
	matchesOp := NewMatches([]string{"name"}, "john", false, nil)
	compact, err := matchesOp.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []any{internal.OpMatchesCode, []string{"name"}, "john", false}, compact)
}
