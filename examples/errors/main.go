// Package main demonstrates error handling in JSON Patch operations.
package main

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpatch"
)

func main() {
	doc := map[string]any{"name": "John", "age": 30}

	fmt.Println("=== Error Handling Example ===")
	fmt.Printf("Original document: %+v\n\n", doc)

	// Example 1: Using ApplyOp with individual operations (demonstrates op package usage)
	fmt.Println("--- Using ApplyOp with individual operations ---")
	
	// Test that passes
	fmt.Println("\n1. Test that should pass:")
	testOp := op.NewOpTestOperation([]string{"name"}, "John")
	result, err := jsonpatch.ApplyOp(doc, testOp)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Test passed: %+v\n", result.Doc)
	}

	// Test that fails
	fmt.Println("\n2. Test that should fail:")
	failOp := op.NewOpTestOperation([]string{"name"}, "Jane")
	_, err = jsonpatch.ApplyOp(doc, failOp)
	if err != nil {
		fmt.Printf("Test failed as expected: %v\n", err)
	} else {
		fmt.Println("Test unexpectedly passed!")
	}

	// Example 2: Using ApplyPatch with Operation structs (consistent with other examples)
	fmt.Println("\n--- Using ApplyPatch with Operation structs ---")
	
	fmt.Println("\n3. Test with invalid path:")
	invalidPatch := []jsonpatch.Operation{
		{Op: "test", Path: "/nonexistent", Value: "value"},
	}
	_, err = jsonpatch.ApplyPatch(doc, invalidPatch)
	if err != nil {
		fmt.Printf("Invalid path error: %v\n", err)
	}

	fmt.Println("\n4. Multiple operations with one failing:")
	mixedPatch := []jsonpatch.Operation{
		{Op: "test", Path: "/name", Value: "John"},  // Should pass
		{Op: "test", Path: "/age", Value: "wrong"},  // Should fail
		{Op: "replace", Path: "/name", Value: "Jane"}, // Won't be reached
	}
	_, err = jsonpatch.ApplyPatch(doc, mixedPatch)
	if err != nil {
		fmt.Printf("Patch failed (as expected): %v\n", err)
	}

	fmt.Println("\n5. Successful multi-operation patch:")
	successPatch := []jsonpatch.Operation{
		{Op: "test", Path: "/name", Value: "John"},
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/city", Value: "New York"},
	}
	patchResult, err := jsonpatch.ApplyPatch(doc, successPatch)
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	} else {
		fmt.Printf("Patch succeeded: %+v\n", patchResult.Doc)
		fmt.Printf("Number of operations applied: %d\n", len(patchResult.Res))
	}
}