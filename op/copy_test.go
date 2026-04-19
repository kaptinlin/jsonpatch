package op

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopy_ArrayInsertCreatesNewSlice(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"values": []any{"a", "c", nil, nil},
	}
	originalValues := doc["values"].([]any)
	originalPointer := fmt.Sprintf("%p", originalValues)

	copyOp := NewCopy([]string{"values", "1"}, []string{"values", "0"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err)

	resultDoc := result.Doc.(map[string]any)
	resultValues := resultDoc["values"].([]any)
	assert.Equal(t, []any{"a", "c", nil, nil}, originalValues)
	assert.Equal(t, []any{"a", "a", "c", nil, nil}, resultValues)
	assert.NotEqual(t, originalPointer, fmt.Sprintf("%p", resultValues))
	assert.Equal(t, "c", result.Old)
}

func TestCopy_ArrayAppendCreatesNewSlice(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"values": []any{"a", "b", nil, nil},
	}
	originalValues := doc["values"].([]any)
	originalPointer := fmt.Sprintf("%p", originalValues)

	copyOp := NewCopy([]string{"values", "4"}, []string{"values", "0"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err)

	resultDoc := result.Doc.(map[string]any)
	resultValues := resultDoc["values"].([]any)
	assert.Equal(t, []any{"a", "b", nil, nil}, originalValues)
	assert.Equal(t, []any{"a", "b", nil, nil, "a"}, resultValues)
	assert.NotEqual(t, originalPointer, fmt.Sprintf("%p", resultValues))
	assert.Equal(t, nil, result.Old)
}
