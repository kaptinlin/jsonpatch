package op

import (
	"errors"
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpMatches_Basic(t *testing.T) {
	tests := []struct {
		name          string
		doc           interface{}
		path          []string
		pattern       string
		ignoreCase    bool
		expectError   bool
		expectedError error
	}{
		{
			name:        "test simple pattern match",
			doc:         map[string]interface{}{"text": "hello123"},
			path:        []string{"text"},
			pattern:     "hello\\d+",
			ignoreCase:  false,
			expectError: false,
		},
		{
			name:        "test pattern no match",
			doc:         map[string]interface{}{"text": "hello"},
			path:        []string{"text"},
			pattern:     "\\d+",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "test case sensitive match",
			doc:         map[string]interface{}{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  false,
			expectError: true,
		},
		{
			name:        "test case insensitive match",
			doc:         map[string]interface{}{"text": "Hello"},
			path:        []string{"text"},
			pattern:     "hello",
			ignoreCase:  true,
			expectError: false,
		},
		{
			name:          "test non string value",
			doc:           map[string]interface{}{"number": 123},
			path:          []string{"number"},
			pattern:       "\\d+",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrNotString,
		},
		{
			name:          "test missing path",
			doc:           map[string]interface{}{"text": "hello"},
			path:          []string{"missing"},
			pattern:       "hello",
			ignoreCase:    false,
			expectError:   true,
			expectedError: ErrPathNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := NewMatches(tt.path, tt.pattern, tt.ignoreCase)
			assert.NoError(t, err)

			result, err := op.Apply(tt.doc)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.True(t, errors.Is(err, tt.expectedError))
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

func TestOpMatches_Constructor(t *testing.T) {
	path := []string{"email"}
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	ignoreCase := false

	op, err := NewMatches(path, pattern, ignoreCase)
	require.NoError(t, err)

	assert.Equal(t, path, op.Path())
	assert.Equal(t, pattern, op.Pattern)
	assert.Equal(t, ignoreCase, op.IgnoreCase)
	assert.Equal(t, internal.OpMatchesType, op.Op())
	assert.Equal(t, internal.OpMatchesCode, op.Code())
}

func TestOpMatches_InvalidPattern(t *testing.T) {
	path := []string{"email"}
	invalidPattern := `[invalid-regex`

	_, err := NewMatches(path, invalidPattern, false)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "regex pattern error")
}

func TestOpMatches_ToJSON(t *testing.T) {
	op, err := NewMatches([]string{"email"}, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, true)
	require.NoError(t, err)

	json, err := op.ToJSON()
	require.NoError(t, err)

	assert.Equal(t, string(internal.OpMatchesType), json["op"])
	assert.Equal(t, "/email", json["path"])
	assert.Equal(t, `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, json["value"])
	assert.Equal(t, true, json["ignore_case"])
}

func TestOpMatches_ToCompact(t *testing.T) {
	op, err := NewMatches([]string{"name"}, "john", true)
	require.NoError(t, err)
	compact, err := op.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{internal.OpMatchesCode, []string{"name"}, "john", true}, compact)
}

func TestOpMatches_ToCompact_WithoutIgnoreCase(t *testing.T) {
	op, err := NewMatches([]string{"name"}, "john", false)
	require.NoError(t, err)
	compact, err := op.ToCompact()
	assert.NoError(t, err)
	assert.Equal(t, []interface{}{internal.OpMatchesCode, []string{"name"}, "john", false}, compact)
}
