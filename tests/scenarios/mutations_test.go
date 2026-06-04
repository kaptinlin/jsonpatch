package jsonpatch_test

import (
	"fmt"
	"maps"
	"reflect"
	"testing"
	"unsafe"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
)

func TestApplyInPlaceFunctionality(t *testing.T) {
	t.Parallel()
	t.Run("Apply - Document Preservation", func(t *testing.T) {
		t.Parallel()
		original := map[string]any{
			"name": "John",
			"age":  30,
			"city": "Boston",
		}
		originalSnapshot := copyMap(original)

		patch := []jsoncodec.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
			{Op: "add", Path: "/email", Value: "jane@example.com"},
			{Op: "remove", Path: "/city"},
		}

		result, err := apply(t, original, patch)
		if err != nil {
			require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
		}

		assert.Equal(t, originalSnapshot, original)

		resultDoc := result.Doc
		if got := resultDoc["name"]; got != "Jane" {
			assert.Fail(t, fmt.Sprintf("result name = %v, want %q", got, "Jane"))
		}
		if got := resultDoc["email"]; got != "jane@example.com" {
			assert.Fail(t, fmt.Sprintf("result email = %v, want %q", got, "jane@example.com"))
		}
		if _, ok := resultDoc["city"]; ok {
			assert.Fail(t, "result should not have city field")
		}

		if isSameMapObject(original, resultDoc) {
			assert.Fail(t, "result should be a different object from original")
		}
	})

	t.Run("ApplyInPlace - Modification", func(t *testing.T) {
		t.Parallel()
		original := map[string]any{
			"name": "John",
			"age":  30,
			"city": "Boston",
		}

		patch := []jsoncodec.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
			{Op: "add", Path: "/email", Value: "jane@example.com"},
			{Op: "remove", Path: "/city"},
		}

		result, err := applyInPlace(t, original, patch)
		if err != nil {
			require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v, want nil", err))
		}

		if got := original["name"]; got != "Jane" {
			assert.Fail(t, fmt.Sprintf("original name = %v, want %q", got, "Jane"))
		}
		if got := original["email"]; got != "jane@example.com" {
			assert.Fail(t, fmt.Sprintf("original email = %v, want %q", got, "jane@example.com"))
		}
		if _, ok := original["city"]; ok {
			assert.Fail(t, "original should not have city field")
		}

		resultDoc := result.Doc
		if !isSameMapObject(original, resultDoc) {
			assert.Fail(t, "result should be the same object as original")
		}

		assert.Equal(t, original, resultDoc)
	})

	t.Run("Array Operations with ApplyInPlace", func(t *testing.T) {
		t.Parallel()
		t.Run("Apply - Array Preservation", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}
			originalSnapshot := make([]any, len(original))
			copy(originalSnapshot, original)

			patch := []jsoncodec.Operation{
				{Op: "replace", Path: "/1", Value: "blueberry"},
				{Op: "add", Path: "/-", Value: "date"},
			}

			result, err := apply(t, original, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}

			assert.Equal(t, originalSnapshot, original)

			resultArray := result.Doc
			if got := resultArray[1]; got != "blueberry" {
				assert.Fail(t, fmt.Sprintf("result[1] = %v, want %q", got, "blueberry"))
			}
			if got := resultArray[3]; got != "date" {
				assert.Fail(t, fmt.Sprintf("result[3] = %v, want %q", got, "date"))
			}
		})

		t.Run("ApplyInPlace - Array Modification", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}

			patch := []jsoncodec.Operation{
				{Op: "replace", Path: "/1", Value: "blueberry"},
			}

			result, err := applyInPlace(t, original, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v, want nil", err))
			}

			if got := original[1]; got != "blueberry" {
				assert.Fail(t, fmt.Sprintf("original[1] = %v, want %q", got, "blueberry"))
			}

			resultArray := result.Doc
			if !isSameSliceObject(original, resultArray) {
				assert.Fail(t, "result should be the same slice as original")
			}
		})

		t.Run("Array Addition - New Slice Required", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}

			patch := []jsoncodec.Operation{
				{Op: "add", Path: "/-", Value: "date"},
			}

			result, err := applyInPlace(t, original, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v, want nil", err))
			}

			if got := len(original); got != 3 {
				assert.Fail(t, fmt.Sprintf("original slice length = %d, want 3 (Go limitation)", got))
			}

			resultArray := result.Doc
			if got := len(resultArray); got != 4 {
				assert.Fail(t, fmt.Sprintf("result slice length = %d, want 4", got))
			}
			if got := resultArray[3]; got != "date" {
				assert.Fail(t, fmt.Sprintf("result[3] = %v, want %q", got, "date"))
			}
		})
	})

	t.Run("Primitive Types - Go Language Behavior", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name     string
			value    any
			patch    []jsoncodec.Operation
			expected any
		}{
			{
				name:     "String Replacement",
				value:    "hello",
				patch:    []jsoncodec.Operation{{Op: "replace", Path: "", Value: "world"}},
				expected: "world",
			},
			{
				name:     "Integer Replacement",
				value:    42,
				patch:    []jsoncodec.Operation{{Op: "replace", Path: "", Value: 99}},
				expected: 99,
			},
			{
				name:     "Boolean Replacement",
				value:    true,
				patch:    []jsoncodec.Operation{{Op: "replace", Path: "", Value: false}},
				expected: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				original := tc.value

				for _, inPlace := range []bool{false, true} {
					t.Run(fmt.Sprintf("InPlace_%v", inPlace), func(t *testing.T) {
						var (
							result *jsonpatch.Result[any]
							err    error
						)
						if inPlace {
							result, err = applyInPlace(t, original, tc.patch)
						} else {
							result, err = apply(t, original, tc.patch)
						}
						if err != nil {
							require.FailNow(t, fmt.Sprintf("Apply error = %v, want nil", err))
						}

						if original != tc.value {
							assert.Fail(t, fmt.Sprintf("original = %v, want %v (primitives are immutable in Go)", original, tc.value))
						}

						assert.Equal(t, tc.expected, result.Doc, "result.Doc")
					})
				}
			})
		}
	})

	t.Run("Complex Nested Structures", func(t *testing.T) {
		t.Parallel()
		original := map[string]any{
			"user": map[string]any{
				"name": "John",
				"preferences": map[string]any{
					"theme":         "dark",
					"notifications": true,
				},
			},
			"items": []any{
				map[string]any{"id": 1, "name": "item1"},
				map[string]any{"id": 2, "name": "item2"},
			},
		}

		patch := []jsoncodec.Operation{
			{Op: "replace", Path: "/user/name", Value: "Jane"},
			{Op: "add", Path: "/user/email", Value: "jane@example.com"},
			{Op: "replace", Path: "/items/0/name", Value: "updated_item1"},
		}

		t.Run("ApplyInPlace - Deep Modification", func(t *testing.T) {
			t.Parallel()
			testDoc := deepCopyMap(original)

			result, err := applyInPlace(t, testDoc, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v, want nil", err))
			}

			if got := testDoc["user"].(map[string]any)["name"]; got != "Jane" {
				assert.Fail(t, fmt.Sprintf("testDoc user name = %v, want %q", got, "Jane"))
			}
			if got := testDoc["user"].(map[string]any)["email"]; got != "jane@example.com" {
				assert.Fail(t, fmt.Sprintf("testDoc user email = %v, want %q", got, "jane@example.com"))
			}
			if got := testDoc["items"].([]any)[0].(map[string]any)["name"]; got != "updated_item1" {
				assert.Fail(t, fmt.Sprintf("testDoc items[0] name = %v, want %q", got, "updated_item1"))
			}

			if !isSameMapObject(testDoc, result.Doc) {
				assert.Fail(t, "result should be the same object as testDoc")
			}
		})
	})
}

func TestApplyInPlacePerformanceCharacteristics(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	largeDoc := createLargeDocument(1000)
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/field_500", Value: "modified"},
		{Op: "add", Path: "/new_field", Value: "new_value"},
	}

	t.Run("Performance Validation", func(t *testing.T) {
		t.Parallel()
		resultFalse, err := apply(t, deepCopyMap(largeDoc), patch)
		if err != nil {
			require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
		}

		resultTrue, err := applyInPlace(t, deepCopyMap(largeDoc), patch)
		if err != nil {
			require.FailNow(t, fmt.Sprintf("ApplyInPlace() error = %v, want nil", err))
		}

		assert.Equal(t, resultFalse.Doc, resultTrue.Doc)
	})
}

func apply[T jsonpatch.Document](t testing.TB, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
}

func applyInPlace[T jsonpatch.Document](t testing.TB, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	if err := jsonpatch.ApplyInPlace(patch, &doc); err != nil {
		return nil, err
	}
	return &jsonpatch.Result[T]{Doc: doc}, nil
}

func copyMap(original map[string]any) map[string]any {
	mapCopy := make(map[string]any)
	maps.Copy(mapCopy, original)
	return mapCopy
}

func deepCopyMap(original map[string]any) map[string]any {
	mapCopy := make(map[string]any)
	for k, v := range original {
		switch val := v.(type) {
		case map[string]any:
			mapCopy[k] = deepCopyMap(val)
		case []any:
			copySlice := make([]any, len(val))
			for i, item := range val {
				if itemMap, ok := item.(map[string]any); ok {
					copySlice[i] = deepCopyMap(itemMap)
				} else {
					copySlice[i] = item
				}
			}
			mapCopy[k] = copySlice
		default:
			mapCopy[k] = v
		}
	}
	return mapCopy
}

func isSameMapObject(a, b map[string]any) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)
	return aVal.Pointer() == bVal.Pointer()
}

func isSameSliceObject(a, b []any) bool {
	return unsafe.SliceData(a) == unsafe.SliceData(b) // #nosec G103 - intentional pointer comparison for test utility
}

func createLargeDocument(size int) map[string]any {
	doc := make(map[string]any)
	for i := range size {
		doc[fmt.Sprintf("field_%d", i)] = fmt.Sprintf("value_%d", i)
	}
	return doc
}

func BenchmarkApplyInPlaceVsApply(b *testing.B) {
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/field_500", Value: "modified"},
	}

	b.Run("Apply", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			doc := createLargeDocument(1000)
			_, err := apply(b, doc, patch)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("ApplyInPlace", func(b *testing.B) {
		b.ResetTimer()
		for b.Loop() {
			doc := createLargeDocument(1000)
			_, err := applyInPlace(b, doc, patch)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
