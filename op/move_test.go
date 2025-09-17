package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpMove_Basic(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
		"qux": map[string]interface{}{
			"nested": "value",
		},
	}

	// Test moving a simple field
	moveOp := NewMove([]string{"qux", "moved"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move should succeed for existing field")

	// Check that the field was moved
	modifiedDoc := result.Doc.(map[string]interface{})
	assert.Nil(t, result.Old, "Old value should be nil when moving to new location")
	assert.NotContains(t, modifiedDoc, "foo", "Source field should be removed")
	assert.Equal(t, "bar", modifiedDoc["qux"].(map[string]interface{})["moved"], "Field should be moved to target path")
	assert.Equal(t, 123, modifiedDoc["baz"], "Other fields should remain unchanged")
}

func TestOpMove_Array(t *testing.T) {
	// Create a test document with array
	doc := map[string]interface{}{
		"items": []interface{}{
			"first",
			"second",
			"third",
		},
		"target": map[string]interface{}{},
	}

	// Test moving an array element
	moveOp := NewMove([]string{"target", "moved"}, []string{"items", "1"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move should succeed for existing array element")

	// Check that the element was moved
	modifiedDoc := result.Doc.(map[string]interface{})
	items := modifiedDoc["items"].([]interface{})
	target := modifiedDoc["target"].(map[string]interface{})

	assert.Nil(t, result.Old, "Old value should be nil when moving to new location")
	assert.Len(t, items, 2, "Array should have one less element")
	assert.Equal(t, "first", items[0], "First element should remain")
	assert.Equal(t, "third", items[1], "Third element should become second")
	assert.Equal(t, "second", target["moved"], "Element should be moved to target path")
}

func TestOpMove_FromNonExistent(t *testing.T) {
	// Create a test document
	doc := map[string]interface{}{
		"foo": "bar",
	}

	// Test moving from non-existent path
	moveOp := NewMove([]string{"target"}, []string{"qux"})
	_, err := moveOp.Apply(doc)
	assert.Error(t, err, "Move should fail for non-existent from path")
	assert.Contains(t, err.Error(), "path not found", "Error message should be descriptive")
}

func TestOpMove_SamePath(t *testing.T) {
	// Test moving to the same path (runtime behavior - should be no-op)
	doc := map[string]interface{}{"foo": 1}
	moveOp := NewMove([]string{"foo"}, []string{"foo"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move to same location should have no effect")
	assert.Equal(t, doc, result.Doc, "Document should remain unchanged")
	assert.Nil(t, result.Old, "Old value should be nil for no-op")
}

func TestOpMove_RootArray(t *testing.T) {
	// Test moving within root array
	doc := []interface{}{"first", "second", "third"}
	moveOp := NewMove([]string{"0"}, []string{"2"})
	result, err := moveOp.Apply(doc)
	require.NoError(t, err, "Move within root array should succeed")

	resultArray := result.Doc.([]interface{})
	// Moving "third" (index 2) to position 0, displacing "first"
	assert.Equal(t, []interface{}{"third", "first", "second"}, resultArray, "Root array should be properly reordered")
	assert.Equal(t, "first", result.Old, "Old value should be the displaced element")
}

func TestOpMove_EmptyPath(t *testing.T) {
	// Test moving with empty path
	moveOp := NewMove([]string{}, []string{"foo"})
	err := moveOp.Validate()
	assert.Error(t, err, "Move should fail validation for empty path")
	assert.Contains(t, err.Error(), "path cannot be empty", "Error message should mention empty path")
}

func TestOpMove_EmptyFrom(t *testing.T) {
	// Test moving with empty from path
	moveOp := NewMove([]string{"target"}, []string{})
	err := moveOp.Validate()
	assert.Error(t, err, "Move should fail validation for empty from path")
	assert.Contains(t, err.Error(), "from path cannot be empty", "Error message should mention empty from path")
}

func TestOpMove_InterfaceMethods(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	// Test Op() method
	assert.Equal(t, internal.OpMoveType, moveOp.Op(), "Op() should return correct operation type")

	// Test Code() method
	assert.Equal(t, internal.OpMoveCode, moveOp.Code(), "Code() should return correct operation code")

	// Test Path() method
	assert.Equal(t, []string{"target"}, moveOp.Path(), "Path() should return correct path")

	// Test From() method
	assert.Equal(t, []string{"source"}, moveOp.From(), "From() should return correct from path")

	// Test HasFrom() method
	assert.True(t, moveOp.HasFrom(), "HasFrom() should return true when from path exists")
}

func TestOpMove_ToJSON(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	json, err := moveOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")

	assert.Equal(t, "move", json["op"], "JSON should contain correct op type")
	assert.Equal(t, "/target", json["path"], "JSON should contain correct formatted path")
	assert.Equal(t, "/source", json["from"], "JSON should contain correct formatted from path")
}

func TestOpMove_ToCompact(t *testing.T) {
	moveOp := NewMove([]string{"target"}, []string{"source"})

	// Test verbose format
	compact, err := moveOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpMoveCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"target"}, compact[1], "Second element should be path")
	assert.Equal(t, []string{"source"}, compact[2], "Third element should be from path")

	// Test non-verbose format
	compact, err = moveOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
}

func TestOpMove_Validate(t *testing.T) {
	// Test valid operation
	moveOp := NewMove([]string{"target"}, []string{"source"})
	err := moveOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Test invalid operation (empty path)
	moveOp = NewMove([]string{}, []string{"source"})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "path cannot be empty", "Error message should mention empty path")

	// Test invalid operation (empty from)
	moveOp = NewMove([]string{"target"}, []string{})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "from path cannot be empty", "Error message should mention empty from path")

	// Test invalid operation (same path and from)
	moveOp = NewMove([]string{"same"}, []string{"same"})
	err = moveOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "path and from cannot be the same", "Error message should mention same paths")
}
