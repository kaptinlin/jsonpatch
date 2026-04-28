package op

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	if diff := cmp.Diff([]any{"a", "c", nil, nil}, originalValues); diff != "" {
		t.Errorf("source array mutated (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]any{"a", "a", "c", nil, nil}, resultValues); diff != "" {
		t.Errorf("result array mismatch (-want +got):\n%s", diff)
	}
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
	if diff := cmp.Diff([]any{"a", "b", nil, nil}, originalValues); diff != "" {
		t.Errorf("source array mutated (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff([]any{"a", "b", nil, nil, "a"}, resultValues); diff != "" {
		t.Errorf("result array mismatch (-want +got):\n%s", diff)
	}
	assert.NotEqual(t, originalPointer, fmt.Sprintf("%p", resultValues))
	assert.Equal(t, nil, result.Old)
}

func TestCopy_RootReplacementDeepClonesSource(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"source": map[string]any{"name": "Ada"},
		"other":  "keep",
	}
	result, err := NewCopy(nil, []string{"source"}).Apply(doc)
	require.NoError(t, err)

	gotDoc := result.Doc.(map[string]any)
	gotDoc["name"] = "Grace"
	wantOriginal := map[string]any{
		"source": map[string]any{"name": "Ada"},
		"other":  "keep",
	}
	if diff := cmp.Diff(wantOriginal, doc); diff != "" {
		t.Errorf("source document mutated (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff(wantOriginal, result.Old); diff != "" {
		t.Errorf("old document mismatch (-want +got):\n%s", diff)
	}
}

func TestCopy_ReturnsSourcePathError(t *testing.T) {
	t.Parallel()

	doc := map[string]any{"name": "Ada"}
	_, err := NewCopy([]string{"copy"}, []string{"missing"}).Apply(doc)
	require.ErrorIs(t, err, ErrPathNotFound)
}
