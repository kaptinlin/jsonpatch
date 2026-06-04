// Package main demonstrates error handling in JSON Patch operations.
package main

import (
	"fmt"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/op"
)

func main() {
	doc := map[string]any{"name": "John", "age": 30}

	fmt.Println("=== Error Handling Example ===")
	fmt.Printf("Original document: %+v\n\n", doc)

	// Example 1: Compile Go-built operations from the op package.
	fmt.Println("--- Compiling Go-built operations ---")

	// Test that passes
	fmt.Println("\n1. Test that should pass:")
	testOp := op.NewTest([]string{"name"}, "John")
	compiled, err := jsonpatch.Compile(testOp)
	if err == nil {
		result, err := jsonpatch.Apply(compiled, doc)
		if err == nil {
			fmt.Printf("Test passed: %+v\n", result.Doc)
		}
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Test that fails
	fmt.Println("\n2. Test that should fail:")
	failOp := op.NewTest([]string{"name"}, "Jane")
	compiled, err = jsonpatch.Compile(failOp)
	if err == nil {
		_, err = jsonpatch.Apply(compiled, doc)
	}
	if err != nil {
		fmt.Printf("Test failed as expected: %v\n", err)
	} else {
		fmt.Println("Test unexpectedly passed!")
	}

	// Example 2: Compile JSON-shaped operations.
	fmt.Println("\n--- Compiling JSON-shaped operations ---")

	fmt.Println("\n3. Test with invalid path:")
	invalidPatch := []jsoncodec.Operation{
		{Op: "test", Path: "/nonexistent", Value: "value"},
	}
	compiled, err = jsonpatch.CompileOperations(invalidPatch)
	if err == nil {
		_, err = jsonpatch.Apply(compiled, doc)
	}
	if err != nil {
		fmt.Printf("Invalid path error: %v\n", err)
	}

	fmt.Println("\n4. Multiple operations with one failing:")
	mixedPatch := []jsoncodec.Operation{
		{Op: "test", Path: "/name", Value: "John"},    // Should pass
		{Op: "test", Path: "/age", Value: "wrong"},    // Should fail
		{Op: "replace", Path: "/name", Value: "Jane"}, // Won't be reached
	}
	compiled, err = jsonpatch.CompileOperations(mixedPatch)
	if err == nil {
		_, err = jsonpatch.Apply(compiled, doc)
	}
	if err != nil {
		fmt.Printf("Patch failed (as expected): %v\n", err)
	}

	fmt.Println("\n5. Successful multi-operation patch:")
	successPatch := []jsoncodec.Operation{
		{Op: "test", Path: "/name", Value: "John"},
		{Op: "replace", Path: "/name", Value: "Jane"},
		{Op: "add", Path: "/city", Value: "New York"},
	}
	compiled, err = jsonpatch.CompileOperations(successPatch)
	if err == nil {
		patchResult, err := jsonpatch.Apply(compiled, doc)
		if err == nil {
			fmt.Printf("Patch succeeded: %+v\n", patchResult.Doc)
			fmt.Printf("Number of operations applied: %d\n", len(patchResult.Steps))
		}
	}
	if err != nil {
		fmt.Printf("Unexpected error: %v\n", err)
	}
}
