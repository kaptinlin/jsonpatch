package op

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpCopy_Basic(t *testing.T) {
	doc := map[string]interface{}{
		"foo": "bar",
		"baz": 123,
	}

	copyOp := NewCopy([]string{"baz_copy"}, []string{"baz"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err, "Copy should succeed for existing field")
	modifiedDoc := result.Doc.(map[string]interface{})
	assert.Equal(t, 123, modifiedDoc["baz_copy"], "Copied value should match source")
	assert.Equal(t, 123, modifiedDoc["baz"], "Source value should remain unchanged")
}

func TestOpCopy_Nested(t *testing.T) {
	doc := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "baz",
			"qux": 123,
		},
	}

	copyOp := NewCopy([]string{"foo", "bar_copy"}, []string{"foo", "bar"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err, "Copy should succeed for existing nested field")
	foo := result.Doc.(map[string]interface{})["foo"].(map[string]interface{})
	assert.Equal(t, "baz", foo["bar_copy"], "Copied nested value should match source")
	assert.Equal(t, "baz", foo["bar"], "Source nested value should remain unchanged")
}

func TestOpCopy_Array(t *testing.T) {
	doc := map[string]interface{}{
		"arr": []interface{}{1, 2, 3},
	}

	copyOp := NewCopy([]string{"arr", "3"}, []string{"arr", "1"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err, "Copy should succeed for array element")
	arr := result.Doc.(map[string]interface{})["arr"].([]interface{})
	assert.Equal(t, 2, arr[3], "Copied array value should match source")
	assert.Equal(t, 2, arr[1], "Source array value should remain unchanged")
	assert.Equal(t, 4, len(arr), "Array should have one more element")
}

func TestOpCopy_DeepClone(t *testing.T) {
	doc := map[string]interface{}{
		"obj": map[string]interface{}{"a": 1},
	}
	copyOp := NewCopy([]string{"obj_copy"}, []string{"obj"})
	result, err := copyOp.Apply(doc)
	require.NoError(t, err, "Copy should succeed for object value")
	obj := doc["obj"].(map[string]interface{})
	objCopy := result.Doc.(map[string]interface{})["obj_copy"].(map[string]interface{})
	assert.Equal(t, obj, objCopy, "Copied object should be equal to source")
	obj["a"] = 2
	assert.NotEqual(t, obj["a"], objCopy["a"], "Copied object should be deep cloned")
}

func TestOpCopy_FromNonExistent(t *testing.T) {
	doc := map[string]interface{}{"foo": "bar"}
	copyOp := NewCopy([]string{"baz"}, []string{"qux"})
	_, err := copyOp.Apply(doc)
	assert.Error(t, err, "Copy should fail for non-existent from path")
	assert.Contains(t, err.Error(), "path not found", "Error message should be descriptive")
}

func TestOpCopy_SamePath(t *testing.T) {
	copyOp := NewCopy([]string{"foo"}, []string{"foo"})
	err := copyOp.Validate()
	assert.Error(t, err, "Copy should fail validation for same path and from")
	assert.Contains(t, err.Error(), "path and from cannot be the same", "Error message should mention same paths")
}

func TestOpCopy_EmptyPath(t *testing.T) {
	// Empty path is valid for copy - it copies to root
	copyOp := NewCopy([]string{}, []string{"foo"})
	err := copyOp.Validate()
	assert.NoError(t, err, "Copy to root (empty path) should be valid")

	// Test actual copy to root
	doc := map[string]interface{}{"foo": "bar", "other": "value"}
	result, err := copyOp.Apply(doc)
	assert.NoError(t, err, "Copy to root should succeed")
	assert.Equal(t, "bar", result.Doc, "Root should be replaced with copied value")
}

func TestOpCopy_EmptyFrom(t *testing.T) {
	copyOp := NewCopy([]string{"foo"}, []string{})
	err := copyOp.Validate()
	assert.Error(t, err, "Copy should fail validation for empty from path")
	assert.Contains(t, err.Error(), "from path cannot be empty", "Error message should mention empty from path")
}

func TestOpCopy_InterfaceMethods(t *testing.T) {
	copyOp := NewCopy([]string{"target"}, []string{"source"})
	assert.Equal(t, internal.OpCopyType, copyOp.Op(), "Op() should return correct operation type")
	assert.Equal(t, internal.OpCopyCode, copyOp.Code(), "Code() should return correct operation code")
	assert.Equal(t, []string{"target"}, copyOp.Path(), "Path() should return correct path")
	assert.Equal(t, []string{"source"}, copyOp.From(), "From() should return correct from path")
	assert.True(t, copyOp.HasFrom(), "HasFrom() should return true when from path exists")
}

func TestOpCopy_ToJSON(t *testing.T) {
	copyOp := NewCopy([]string{"target"}, []string{"source"})
	json, err := copyOp.ToJSON()
	require.NoError(t, err, "ToJSON should not fail for valid operation")
	assert.Equal(t, "copy", json["op"], "JSON should contain correct op type")
	assert.Equal(t, "/target", json["path"], "JSON should contain correct formatted path")
	assert.Equal(t, "/source", json["from"], "JSON should contain correct formatted from path")
}

func TestOpCopy_ToCompact(t *testing.T) {
	copyOp := NewCopy([]string{"target"}, []string{"source"})
	compact, err := copyOp.ToCompact()
	require.NoError(t, err, "ToCompact should not fail for valid operation")
	require.Len(t, compact, 3, "Compact format should have 3 elements")
	assert.Equal(t, internal.OpCopyCode, compact[0], "First element should be operation code")
	assert.Equal(t, []string{"target"}, compact[1], "Second element should be path")
	assert.Equal(t, []string{"source"}, compact[2], "Third element should be from path")
}

func TestOpCopy_Validate(t *testing.T) {
	copyOp := NewCopy([]string{"target"}, []string{"source"})
	err := copyOp.Validate()
	assert.NoError(t, err, "Valid operation should not fail validation")

	// Empty path is valid for copy (copies to root)
	copyOp = NewCopy([]string{}, []string{"source"})
	err = copyOp.Validate()
	assert.NoError(t, err, "Copy to root (empty path) should be valid")

	copyOp = NewCopy([]string{"target"}, []string{})
	err = copyOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "from path cannot be empty", "Error message should mention empty from path")

	copyOp = NewCopy([]string{"same"}, []string{"same"})
	err = copyOp.Validate()
	assert.Error(t, err, "Invalid operation should fail validation")
	assert.Contains(t, err.Error(), "path and from cannot be the same", "Error message should mention same paths")
}

func TestOpCopy_RFC6902_CopyToRoot(t *testing.T) {
	// RFC 6902 compliance: copy to root with different types
	tests := []struct {
		name     string
		doc      any
		from     []string
		expected any
	}{
		{
			name:     "copy null to root",
			doc:      map[string]interface{}{"foo": nil, "bar": "value"},
			from:     []string{"foo"},
			expected: nil,
		},
		{
			name:     "copy string to root",
			doc:      map[string]interface{}{"foo": "hello", "bar": "world"},
			from:     []string{"foo"},
			expected: "hello",
		},
		{
			name:     "copy object to root",
			doc:      map[string]interface{}{"foo": map[string]interface{}{"nested": true}, "bar": "value"},
			from:     []string{"foo"},
			expected: map[string]interface{}{"nested": true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			copyOp := NewCopy([]string{}, tt.from)
			result, err := copyOp.Apply(tt.doc)
			require.NoError(t, err, "Copy to root should work")
			assert.Equal(t, tt.expected, result.Doc, "Root should be replaced with copied value")
		})
	}
}
