package jsonpatch_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

// TestEmptyDocumentHandling tests edge cases with empty/undefined documents
// Ported from TypeScript: patch.scenarious.spec.ts
func TestEmptyDocumentHandling(t *testing.T) {
	t.Run("cannot add key to empty document", func(t *testing.T) {
		op := jsonpatch.Operation{"op": "add", "path": "/foo", "value": 123}
		var doc interface{}
		_ = testutils.ApplyOperationWithError(t, doc, op)
	})

	t.Run("can overwrite empty document", func(t *testing.T) {
		op := jsonpatch.Operation{"op": "add", "path": "", "value": map[string]interface{}{"foo": 123}}
		var doc interface{}
		result := testutils.ApplyOperation(t, doc, op)
		expected := map[string]interface{}{"foo": float64(123)} // JSON unmarshaling converts numbers to float64
		assert.Equal(t, expected, result)
	})

	t.Run("cannot add value to nonexisting path", func(t *testing.T) {
		doc := map[string]interface{}{"foo": 123}
		op := jsonpatch.Operation{"op": "add", "path": "/foo/bar/baz", "value": "test"}
		_ = testutils.ApplyOperationWithError(t, doc, op)
	})
}

// TestNumberTypeCoercion tests number type handling edge cases
func TestNumberTypeCoercion(t *testing.T) {
	t.Run("inc operation with boolean values", func(t *testing.T) {
		doc := map[string]interface{}{
			"trueVal":  true,
			"falseVal": false,
		}
		ops := []jsonpatch.Operation{
			{"op": "inc", "path": "/trueVal", "inc": 1},
			{"op": "inc", "path": "/falseVal", "inc": 1},
		}
		result := testutils.ApplyOperations(t, doc, ops)

		expected := map[string]interface{}{
			"trueVal":  float64(2), // true converts to 1, then +1 = 2
			"falseVal": float64(1), // false converts to 0, then +1 = 1
		}
		assert.Equal(t, expected, result)
	})

	t.Run("inc operation with string numbers", func(t *testing.T) {
		doc := map[string]interface{}{"numStr": "42"}
		op := jsonpatch.Operation{"op": "inc", "path": "/numStr", "inc": 8}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"numStr": float64(50)}
		assert.Equal(t, expected, result)
	})

	t.Run("inc operation with floating point precision", func(t *testing.T) {
		doc := map[string]interface{}{"val": 0.1}
		op := jsonpatch.Operation{"op": "inc", "path": "/val", "inc": 0.2}
		result := testutils.ApplyOperation(t, doc, op)

		// Note: Floating point arithmetic precision
		resultMap := result.(map[string]interface{})
		resultVal := resultMap["val"].(float64)
		assert.InDelta(t, 0.3, resultVal, 0.0001)
	})
}

// TestArrayBoundaryConditions tests array operations at boundaries
func TestArrayBoundaryConditions(t *testing.T) {
	t.Run("add to array at exact length", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{"op": "add", "path": "/3", "value": 4} // Adding at index 3 (length)
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{1, 2, 3, 4}
		assert.Equal(t, expected, result)
	})

	t.Run("remove from array first element", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{"op": "remove", "path": "/0"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{2, 3}
		assert.Equal(t, expected, result)
	})

	t.Run("remove from array last element", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{"op": "remove", "path": "/2"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{1, 2}
		assert.Equal(t, expected, result)
	})
}

// TestStringOperationEdgeCases tests string manipulation edge cases
func TestStringOperationEdgeCases(t *testing.T) {
	t.Run("str_ins at string beginning", func(t *testing.T) {
		doc := map[string]interface{}{"text": "world"}
		op := jsonpatch.Operation{"op": "str_ins", "path": "/text", "pos": 0, "str": "hello "}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "hello world"}
		assert.Equal(t, expected, result)
	})

	t.Run("str_ins at string end", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello"}
		op := jsonpatch.Operation{"op": "str_ins", "path": "/text", "pos": 5, "str": " world"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "hello world"}
		assert.Equal(t, expected, result)
	})

	t.Run("str_del from string beginning", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello world"}
		op := jsonpatch.Operation{"op": "str_del", "path": "/text", "pos": 0, "len": 6}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "world"}
		assert.Equal(t, expected, result)
	})

	t.Run("str_del entire string", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello"}
		op := jsonpatch.Operation{"op": "str_del", "path": "/text", "pos": 0, "len": 5}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": ""}
		assert.Equal(t, expected, result)
	})
}
