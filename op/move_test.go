package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMove_Basic(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]any{
			"nested": "value",
		},
	}

	moveOp := NewMove([]string{"qux", "moved"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move should succeed for existing field")

	modifiedDoc := result.Doc.(map[string]any)
	assert.Nil(t, result.Old, "Old value should be nil when moving to new location")
	assert.NotContains(t, modifiedDoc, "foo", "Source field should be removed")
	assert.Equal(t, "bar", modifiedDoc["qux"].(map[string]any)["moved"], "Field should be moved to target path")
	assert.Equal(t, 123, modifiedDoc["baz"], "Other fields should remain unchanged")
}

func TestMove_Array(t *testing.T) {
	doc := map[string]any{
		"items": []any{
			"first",
			"second",
			"third",
		},
		"target": map[string]any{},
	}

	moveOp := NewMove([]string{"target", "moved"}, []string{"items", "1"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move should succeed for existing array element")

	modifiedDoc := result.Doc.(map[string]any)
	items := modifiedDoc["items"].([]any)
	target := modifiedDoc["target"].(map[string]any)

	assert.Nil(t, result.Old, "Old value should be nil when moving to new location")
	assert.Len(t, items, 2, "Array should have one less element")
	assert.Equal(t, "first", items[0], "First element should remain")
	assert.Equal(t, "third", items[1], "Third element should become second")
	assert.Equal(t, "second", target["moved"], "Element should be moved to target path")
}

func TestMove_FromNonExistent(t *testing.T) {
	doc := map[string]any{
		"foo": "bar",
	}

	moveOp := NewMove([]string{"target"}, []string{"qux"})
	_, err := moveOp.Apply(doc)
	assert.Error(t, err, "Move should fail for non-existent from path")
	assert.ErrorIs(t, err, ErrPathNotFound)
}

func TestMove_SamePath(t *testing.T) {
	doc := map[string]any{"foo": 1}
	moveOp := NewMove([]string{"foo"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move to same location should have no effect")
	assert.Equal(t, doc, result.Doc, "Document should remain unchanged")
	assert.Nil(t, result.Old, "Old value should be nil for no-op")
}

func TestMove_RootArray(t *testing.T) {
	doc := []any{"first", "second", "third"}
	moveOp := NewMove([]string{"0"}, []string{"2"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move within root array should succeed")

	resultArray := result.Doc.([]any)
	assert.Equal(t, []any{"third", "first", "second"}, resultArray, "Root array should be properly reordered")
	assert.Equal(t, "first", result.Old, "Old value should be the displaced element")
}

func TestMove_EmptyPath(t *testing.T) {
	moveOp := NewMove([]string{}, []string{"foo"})
	err := moveOp.Validate()
	assert.Error(t, err, "Move should fail validation for empty path")
	assert.ErrorIs(t, err, ErrPathEmpty)
}

func TestMove_EmptyFrom(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{})
	err := moveOp.Validate()
	assert.Error(t, err, "Move should fail validation for empty from path")
	assert.ErrorIs(t, err, ErrFromPathEmpty)
}

func TestMove_InterfaceMethods(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	assert.Equal(t, internal.OpMoveType, moveOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpMoveCode, moveOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"target"}, moveOp.Path(), "Path() should return correct path")
	assert.Equal(t, []string{"source"}, moveOp.From(), "From() should return correct from path")
	assert.True(t, moveOp.HasFrom(), "HasFrom() should return true when from path exists")
}

func TestMove_ToJSON(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	got, err := moveOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "move", got.Op, "JSON should contain correct op type")
	assert.Equal(t, "/target", got.Path, "JSON should contain correct formatted path")
	assert.Equal(t, "/source", got.From, "JSON should contain correct formatted from path")
}

func TestMove_ToCompact(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	compact, err := moveOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpMoveCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"target"}, compact[1], "Second element should be path")
	assert.Equal(t, []string{"source"}, compact[2], "Third element should be from path")
}

func TestMove_Validate(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})
	err := moveOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	moveOp = NewMove([]string{}, []string{"source"})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrPathEmpty)

	moveOp = NewMove([]string{"target"}, []string{})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrFromPathEmpty)

	moveOp = NewMove([]string{"same"}, []string{"same"})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.ErrorIs(t, err, ErrPathsIdentical)
}

func TestMove_RFC6902_RemoveAddPattern(t *testing.T) {
	tests := []struct {
		name     string
		doc      map[string]any
		from     []string
		path     []string
		expected map[string]any
	}{
		{
			name: "object property to array element",
			doc: map[string]any{
				"baz": []any{map[string]any{"qux": "hello"}},
				"bar": 1,
			},
			from: []string{"baz", "0", "qux"},
			path: []string{"baz", "1"},
			expected: map[string]any{
				"baz": []any{map[string]any{}, "hello"},
				"bar": 1,
			},
		},
		{
			name: "array element to front",
			doc: map[string]any{
				"users": []any{
					map[string]any{"name": "Alice"},
					map[string]any{"name": "Bob"},
				},
			},
			from: []string{"users", "1"},
			path: []string{"users", "0"},
			expected: map[string]any{
				"users": []any{
					map[string]any{"name": "Bob"},
					map[string]any{"name": "Alice"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			moveOp := NewMove(tt.path, tt.from)
			result, err := moveOp.Apply(tt.doc)
			require.NoError(t, err, "Move operation should work")
			assert.Equal(t, tt.expected, result.Doc, "Move should follow remove->add pattern per RFC 6902")
		})
	}
}
