// Package main demonstrates extended JSON Patch operations.
package main

import (
	"fmt"
	"log"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Extended Operations ===")

	// Original document
	doc := map[string]any{
		"counter": 10,
		"enabled": false,
		"name":    "example",
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	// Extended operations patch
	patch := []jsoncodec.Operation{
		// Inc: increment counter by 15
		{
			Op:   "inc",
			Path: "/counter",
			Inc:  15,
		},
		// Flip: toggle boolean value
		{
			Op:   "flip",
			Path: "/enabled",
		},
		// Add new field
		{
			Op:    "add",
			Path:  "/status",
			Value: "active",
		},
	}

	compiled, err := jsonpatch.CompileOperations(patch, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	result, err := jsonpatch.Apply(compiled, doc)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter extended operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))

	// Demonstrate individual operations
	fmt.Println("\n--- Individual Operations ---")

	// Test inc operation on root value
	numberDoc := 42.0
	incPatch := []jsoncodec.Operation{
		{
			Op:   "inc",
			Path: "", // Root path for primitive value
			Inc:  8,
		},
	}

	compiled, err = jsonpatch.CompileOperations(incPatch, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		log.Fatalf("Inc operation failed: %v", err)
	}
	incResult, err := jsonpatch.Apply(compiled, numberDoc)
	if err != nil {
		log.Fatalf("Inc operation failed: %v", err)
	}
	fmt.Printf("Inc operation: 42 + 8 = %v\n", incResult.Doc)

	// Test flip operation on root boolean
	boolDoc := true
	flipPatch := []jsoncodec.Operation{
		{
			Op:   "flip",
			Path: "", // Root path for primitive value
		},
	}

	compiled, err = jsonpatch.CompileOperations(flipPatch, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		log.Fatalf("Flip operation failed: %v", err)
	}
	flipResult, err := jsonpatch.Apply(compiled, boolDoc)
	if err != nil {
		log.Fatalf("Flip operation failed: %v", err)
	}
	fmt.Printf("Flip operation: true -> %v\n", flipResult.Doc)
}
