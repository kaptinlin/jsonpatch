package jsonpatch_test

import (
	"math"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

// TestEmptyDocumentHandling tests edge cases with empty/undefined documents
// Ported from TypeScript: patch.scenarious.spec.ts
func TestEmptyDocumentHandling(t *testing.T) {
	t.Run("cannot add key to empty document", func(t *testing.T) {
		op := jsonpatch.Operation{Op: "add", Path: "/foo", Value: 123}
		var doc interface{}
		_ = testutils.ApplyOperationWithError(t, doc, op)
	})

	t.Run("can overwrite empty document", func(t *testing.T) {
		op := jsonpatch.Operation{Op: "add", Path: "", Value: map[string]interface{}{"foo": 123}}
		var doc interface{}
		result := testutils.ApplyOperation(t, doc, op)
		expected := map[string]interface{}{"foo": float64(123)} // JSON unmarshaling converts numbers to float64
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("cannot add value to nonexisting path", func(t *testing.T) {
		doc := map[string]interface{}{"foo": 123}
		op := jsonpatch.Operation{Op: "add", Path: "/foo/bar/baz", Value: "test"}
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
			{Op: "inc", Path: "/trueVal", Inc: 1},
			{Op: "inc", Path: "/falseVal", Inc: 1},
		}
		result := testutils.ApplyOperations(t, doc, ops)

		expected := map[string]interface{}{
			"trueVal":  float64(2), // true converts to 1, then +1 = 2
			"falseVal": float64(1), // false converts to 0, then +1 = 1
		}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperations() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("inc operation with string numbers", func(t *testing.T) {
		doc := map[string]interface{}{"numStr": "42"}
		op := jsonpatch.Operation{Op: "inc", Path: "/numStr", Inc: 8}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"numStr": float64(50)}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("inc operation with floating point precision", func(t *testing.T) {
		doc := map[string]interface{}{"val": 0.1}
		op := jsonpatch.Operation{Op: "inc", Path: "/val", Inc: 0.2}
		result := testutils.ApplyOperation(t, doc, op)

		// Note: Floating point arithmetic precision
		resultMap := result.(map[string]interface{})
		resultVal := resultMap["val"].(float64)
		if math.Abs(resultVal-0.3) > 0.0001 {
			t.Errorf("result val = %v, want ~0.3 (within 0.0001)", resultVal)
		}
	})
}

// TestArrayBoundaryConditions tests array operations at boundaries
func TestArrayBoundaryConditions(t *testing.T) {
	t.Run("add to array at exact length", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{Op: "add", Path: "/3", Value: 4} // Adding at index 3 (length)
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{1, 2, 3, 4}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("remove from array first element", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{Op: "remove", Path: "/0"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{2, 3}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("remove from array last element", func(t *testing.T) {
		doc := []interface{}{1, 2, 3}
		op := jsonpatch.Operation{Op: "remove", Path: "/2"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := []interface{}{1, 2}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})
}

// TestStringOperationEdgeCases tests string manipulation edge cases
func TestStringOperationEdgeCases(t *testing.T) {
	t.Run("str_ins at string beginning", func(t *testing.T) {
		doc := map[string]interface{}{"text": "world"}
		op := jsonpatch.Operation{Op: "str_ins", Path: "/text", Pos: 0, Str: "hello "}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "hello world"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("str_ins at string end", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello"}
		op := jsonpatch.Operation{Op: "str_ins", Path: "/text", Pos: 5, Str: " world"}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "hello world"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("str_del from string beginning", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello world"}
		op := jsonpatch.Operation{Op: "str_del", Path: "/text", Pos: 0, Len: 6}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": "world"}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("str_del entire string", func(t *testing.T) {
		doc := map[string]interface{}{"text": "hello"}
		op := jsonpatch.Operation{Op: "str_del", Path: "/text", Pos: 0, Len: 5}
		result := testutils.ApplyOperation(t, doc, op)

		expected := map[string]interface{}{"text": ""}
		if diff := cmp.Diff(expected, result); diff != "" {
			t.Errorf("ApplyOperation() mismatch (-want +got):\n%s", diff)
		}
	})
}
