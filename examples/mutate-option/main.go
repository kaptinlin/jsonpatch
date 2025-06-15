package main

import (
	"encoding/json"
	"fmt"
	"reflect"

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

	result1, err := jsonpatch.ApplyPatch(doc1, patch, jsonpatch.ApplyPatchOptions{
		Mutate: false, // Default behavior
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc1))
	fmt.Printf("Result:   %s\n", toJSON(result1.Doc))
	fmt.Printf("Same object: %v\n", isSameObject(doc1, result1.Doc))

	// Example 2: Performance mode - modifies original
	fmt.Println("\n⚡ Mutate: true (Performance Mode)")
	doc2 := copyDocument(document)

	result2, err := jsonpatch.ApplyPatch(doc2, patch, jsonpatch.ApplyPatchOptions{
		Mutate: true, // High performance
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Original: %s\n", toJSON(doc2))
	fmt.Printf("Result:   %s\n", toJSON(result2.Doc))
	fmt.Printf("Same object: %v\n", isSameObject(doc2, result2.Doc))

	// Performance benefits
	fmt.Println("\n📊 Performance Benefits (Mutate: true)")
	fmt.Println("   • 12-27% faster execution")
	fmt.Println("   • 27% less memory usage")
	fmt.Println("   • Fewer memory allocations")

	// Use cases
	fmt.Println("\n🎯 When to Use Each Mode")
	fmt.Println("Mutate: false (Safe)")
	fmt.Println("   • Need to preserve original document")
	fmt.Println("   • Working with shared data")
	fmt.Println("   • Safety over performance")

	fmt.Println("\nMutate: true (Fast)")
	fmt.Println("   • Large documents")
	fmt.Println("   • High-frequency operations")
	fmt.Println("   • Memory-constrained environments")
	fmt.Println("   • Don't need original document")

	// Go language limitation example
	fmt.Println("\n⚠️  Go Language Limitation")
	primitiveDoc := "hello"
	primitiveResult, _ := jsonpatch.ApplyPatch(primitiveDoc, []jsonpatch.Operation{
		{"op": "replace", "path": "", "value": "world"},
	}, jsonpatch.ApplyPatchOptions{Mutate: true})

	fmt.Printf("Primitive original: %q (unchanged)\n", primitiveDoc)
	fmt.Printf("Primitive result:   %q\n", primitiveResult.Doc)
	fmt.Println("Note: Primitive types cannot be mutated in-place in Go")
}

// Helper functions
func copyDocument(doc map[string]interface{}) map[string]interface{} {
	copy := make(map[string]interface{})
	for k, v := range doc {
		copy[k] = v
	}
	return copy
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
