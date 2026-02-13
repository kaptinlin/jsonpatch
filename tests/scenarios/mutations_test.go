package jsonpatch_test

import (
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/assert"
)

func TestMutateOptionFunctionality(t *testing.T) {
	t.Parallel()
	t.Run("Mutate False - Document Preservation", func(t *testing.T) {
		t.Parallel()
		original := map[string]any{
			"name": "John",
			"age":  30,
			"city": "Boston",
		}
		originalSnapshot := copyMap(original)

		patch := []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
			{Op: "add", Path: "/email", Value: "jane@example.com"},
			{Op: "remove", Path: "/city"},
		}

		options := jsonpatch.WithMutate(false)
		result, err := jsonpatch.ApplyPatch(original, patch, options)
		if err != nil {
			t.Fatalf("ApplyPatch() error = %v, want nil", err)
		}

		assert.Equal(t, originalSnapshot, original)

		resultDoc := result.Doc
		if got := resultDoc["name"]; got != "Jane" {
			t.Errorf("result name = %v, want %q", got, "Jane")
		}
		if got := resultDoc["email"]; got != "jane@example.com" {
			t.Errorf("result email = %v, want %q", got, "jane@example.com")
		}
		if _, ok := resultDoc["city"]; ok {
			assert.Fail(t, "result should not have city field")
		}

		if isSameMapObject(original, resultDoc) {
			assert.Fail(t, "result should be a different object from original")
		}
	})

	t.Run("Mutate True - In-Place Modification", func(t *testing.T) {
		t.Parallel()
		original := map[string]any{
			"name": "John",
			"age":  30,
			"city": "Boston",
		}

		patch := []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Jane"},
			{Op: "add", Path: "/email", Value: "jane@example.com"},
			{Op: "remove", Path: "/city"},
		}

		options := jsonpatch.WithMutate(true)
		result, err := jsonpatch.ApplyPatch(original, patch, options)
		if err != nil {
			t.Fatalf("ApplyPatch() error = %v, want nil", err)
		}

		if got := original["name"]; got != "Jane" {
			t.Errorf("original name = %v, want %q", got, "Jane")
		}
		if got := original["email"]; got != "jane@example.com" {
			t.Errorf("original email = %v, want %q", got, "jane@example.com")
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

	t.Run("Array Operations with Mutate", func(t *testing.T) {
		t.Parallel()
		t.Run("Mutate False - Array Preservation", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}
			originalSnapshot := make([]any, len(original))
			copy(originalSnapshot, original)

			patch := []jsonpatch.Operation{
				{Op: "replace", Path: "/1", Value: "blueberry"},
				{Op: "add", Path: "/-", Value: "date"},
			}

			options := jsonpatch.WithMutate(false)
			result, err := jsonpatch.ApplyPatch(original, patch, options)
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}

			assert.Equal(t, originalSnapshot, original)

			resultArray := result.Doc
			if got := resultArray[1]; got != "blueberry" {
				t.Errorf("result[1] = %v, want %q", got, "blueberry")
			}
			if got := resultArray[3]; got != "date" {
				t.Errorf("result[3] = %v, want %q", got, "date")
			}
		})

		t.Run("Mutate True - Array Modification (Go Limitation)", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}

			patch := []jsonpatch.Operation{
				{Op: "replace", Path: "/1", Value: "blueberry"},
			}

			options := jsonpatch.WithMutate(true)
			result, err := jsonpatch.ApplyPatch(original, patch, options)
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}

			if got := original[1]; got != "blueberry" {
				t.Errorf("original[1] = %v, want %q", got, "blueberry")
			}

			resultArray := result.Doc
			if !isSameSliceObject(original, resultArray) {
				assert.Fail(t, "result should be the same slice as original")
			}
		})

		t.Run("Array Addition - New Slice Required", func(t *testing.T) {
			t.Parallel()
			original := []any{"apple", "banana", "cherry"}

			patch := []jsonpatch.Operation{
				{Op: "add", Path: "/-", Value: "date"},
			}

			options := jsonpatch.WithMutate(true)
			result, err := jsonpatch.ApplyPatch(original, patch, options)
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}

			if got := len(original); got != 3 {
				t.Errorf("original slice length = %d, want 3 (Go limitation)", got)
			}

			resultArray := result.Doc
			if got := len(resultArray); got != 4 {
				t.Errorf("result slice length = %d, want 4", got)
			}
			if got := resultArray[3]; got != "date" {
				t.Errorf("result[3] = %v, want %q", got, "date")
			}
		})
	})

	t.Run("Primitive Types - Go Language Behavior", func(t *testing.T) {
		t.Parallel()
		testCases := []struct {
			name     string
			value    any
			patch    []jsonpatch.Operation
			expected any
		}{
			{
				name:     "String Replacement",
				value:    "hello",
				patch:    []jsonpatch.Operation{{Op: "replace", Path: "", Value: "world"}},
				expected: "world",
			},
			{
				name:     "Integer Replacement",
				value:    42,
				patch:    []jsonpatch.Operation{{Op: "replace", Path: "", Value: 99}},
				expected: 99,
			},
			{
				name:     "Boolean Replacement",
				value:    true,
				patch:    []jsonpatch.Operation{{Op: "replace", Path: "", Value: false}},
				expected: false,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				original := tc.value

				for _, mutate := range []bool{false, true} {
					t.Run(fmt.Sprintf("Mutate_%v", mutate), func(t *testing.T) {
						options := jsonpatch.WithMutate(mutate)
						result, err := jsonpatch.ApplyPatch(original, tc.patch, options)
						if err != nil {
							t.Fatalf("ApplyPatch() error = %v, want nil", err)
						}

						if original != tc.value {
							t.Errorf("original = %v, want %v (primitives are immutable in Go)", original, tc.value)
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

		patch := []jsonpatch.Operation{
			{Op: "replace", Path: "/user/name", Value: "Jane"},
			{Op: "add", Path: "/user/email", Value: "jane@example.com"},
			{Op: "replace", Path: "/items/0/name", Value: "updated_item1"},
		}

		t.Run("Mutate True - Deep Modification", func(t *testing.T) {
			t.Parallel()
			testDoc := deepCopyMap(original)

			options := jsonpatch.WithMutate(true)
			result, err := jsonpatch.ApplyPatch(testDoc, patch, options)
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}

			if got := testDoc["user"].(map[string]any)["name"]; got != "Jane" {
				t.Errorf("testDoc user name = %v, want %q", got, "Jane")
			}
			if got := testDoc["user"].(map[string]any)["email"]; got != "jane@example.com" {
				t.Errorf("testDoc user email = %v, want %q", got, "jane@example.com")
			}
			if got := testDoc["items"].([]any)[0].(map[string]any)["name"]; got != "updated_item1" {
				t.Errorf("testDoc items[0] name = %v, want %q", got, "updated_item1")
			}

			if !isSameMapObject(testDoc, result.Doc) {
				assert.Fail(t, "result should be the same object as testDoc")
			}
		})
	})
}

func TestMutatePerformanceCharacteristics(t *testing.T) {
	t.Parallel()
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	largeDoc := createLargeDocument(1000)
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/field_500", Value: "modified"},
		{Op: "add", Path: "/new_field", Value: "new_value"},
	}

	t.Run("Performance Validation", func(t *testing.T) {
		t.Parallel()
		optionsFalse := jsonpatch.WithMutate(false)
		resultFalse, err := jsonpatch.ApplyPatch(deepCopyMap(largeDoc), patch, optionsFalse)
		if err != nil {
			t.Fatalf("ApplyPatch(mutate=false) error = %v, want nil", err)
		}

		optionsTrue := jsonpatch.WithMutate(true)
		resultTrue, err := jsonpatch.ApplyPatch(deepCopyMap(largeDoc), patch, optionsTrue)
		if err != nil {
			t.Fatalf("ApplyPatch(mutate=true) error = %v, want nil", err)
		}

		assert.Equal(t, resultFalse.Doc, resultTrue.Doc)
	})
}

func copyMap(original map[string]any) map[string]any {
	mapCopy := make(map[string]any)
	for k, v := range original {
		mapCopy[k] = v
	}
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

func BenchmarkMutateVsClone(b *testing.B) {
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/field_500", Value: "modified"},
	}

	b.Run("Mutate=false", func(b *testing.B) {
		options := jsonpatch.WithMutate(false)
		b.ResetTimer()
		for b.Loop() {
			doc := createLargeDocument(1000)
			_, err := jsonpatch.ApplyPatch(doc, patch, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Mutate=true", func(b *testing.B) {
		options := jsonpatch.WithMutate(true)
		b.ResetTimer()
		for b.Loop() {
			doc := createLargeDocument(1000)
			_, err := jsonpatch.ApplyPatch(doc, patch, options)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
