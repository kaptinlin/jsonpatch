// Package main demonstrates mutate option usage with JSON Patch.
package main

import (
	"fmt"
	"reflect"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("🚀 JSON Patch Mutate Option Usage")
	fmt.Println("==================================")

	// Sample document and patch
	document := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "john@example.com",
	}

	patch := []jsonpatch.Operation{
		{"op": "replace", "path": "/name", "value": "Jane"},
		{"op": "add", "path": "/city", "value": "New York"},
	}

	// Example 1: Safe mode (default) - preserves original
	fmt.Println("\n📋 Mutate: false (Safe Mode)")
	doc1 := copyDocument(document)

	result1, err := jsonpatch.ApplyPatch(doc1, patch, jsonpatch.WithMutate(false))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc1))
	fmt.Printf("Result:   %s\n", toJSON(result1.Doc))
	fmt.Printf("Same object: %v\n", isSameObject(doc1, result1.Doc))

	// Example 2: Performance mode - modifies original
	fmt.Println("\n⚡ Mutate: true (Performance Mode)")
	doc2 := copyDocument(document)

	result2, err := jsonpatch.ApplyPatch(doc2, patch, jsonpatch.WithMutate(true))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc2))
	fmt.Printf("Result:   %s\n", toJSON(result2.Doc))
	fmt.Printf("Same object: %v\n", isSameObject(doc2, result2.Doc))

	// Go language limitation example
	fmt.Println("\n⚠️  Go Language Limitation")
	primitiveDoc := "hello"
	primitiveResult, _ := jsonpatch.ApplyPatch(primitiveDoc, []jsonpatch.Operation{
		{"op": "replace", "path": "", "value": "world"},
	}, jsonpatch.WithMutate(true))

	fmt.Printf("Primitive original: %q (unchanged)\n", primitiveDoc)
	fmt.Printf("Primitive result:   %q\n", primitiveResult.Doc)
	fmt.Println("Note: Primitive types cannot be mutated in-place in Go")
}

// Helper functions
func copyDocument(doc map[string]interface{}) map[string]interface{} {
	docCopy := make(map[string]interface{})
	for k, v := range doc {
		docCopy[k] = v
	}
	return docCopy
}

func toJSON(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func isSameObject(a, b interface{}) bool {
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
