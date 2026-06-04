// Package main demonstrates explicit in-place JSON Patch application.
package main

import (
	"fmt"
	"maps"
	"reflect"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("🚀 JSON Patch In-Place Usage")
	fmt.Println("============================")

	// Sample document and patch
	document := map[string]any{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/city", Value: "New York"},
	}
	compiled, err := jsonpatch.CompileOperations(patch)
	if err != nil {
		panic(err)
	}

	// Example 1: Immutable apply preserves original
	fmt.Println("\n📋 Apply (immutable)")
	doc1 := copyDocument(document)

	result1, err := jsonpatch.Apply(compiled, doc1)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc1))
	fmt.Printf("Result:   %s\n", toJSON(result1.Doc))
	fmt.Printf("Same object: %v\n", isSameObject(doc1, result1.Doc))

	// Example 2: ApplyInPlace modifies the provided variable
	fmt.Println("\n⚡ ApplyInPlace")
	doc2 := copyDocument(document)
	beforeInPlace := doc2

	if err := jsonpatch.ApplyInPlace(compiled, &doc2); err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc2))
	fmt.Printf("Same object: %v\n", isSameObject(beforeInPlace, doc2))

	// Primitive values are updated through the pointer.
	fmt.Println("\n⚙️  Primitive Value")
	primitiveDoc := "hello"
	primitivePatch, err := jsonpatch.CompileOperations([]jsoncodec.Operation{
		{Op: "replace", Path: "", Value: "world"},
	})
	if err != nil {
		panic(err)
	}
	if err := jsonpatch.ApplyInPlace(primitivePatch, &primitiveDoc); err != nil {
		panic(err)
	}

	fmt.Printf("Primitive result: %q\n", primitiveDoc)
}

// Helper functions
func copyDocument(doc map[string]any) map[string]any {
	return maps.Clone(doc)
}

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func isSameObject(a, b any) bool {
	aVal := reflect.ValueOf(a)
	bVal := reflect.ValueOf(b)

	if aVal.Kind() != bVal.Kind() {
		return false
	}

	if aVal.Kind() == reflect.Map {
		return aVal.Pointer() == bVal.Pointer()
	}

	return false
}
