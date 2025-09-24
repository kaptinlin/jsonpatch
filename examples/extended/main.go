// Package main demonstrates extended JSON Patch operations.
package main

import (
	"fmt"
	"log"

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
	patch := []jsonpatch.Operation{
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

	result, err := jsonpatch.ApplyPatch(doc, patch)
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
	incPatch := []jsonpatch.Operation{
		{
			Op:  "inc",
			Path: "",  // Root path for primitive value
			Inc: 8,
		},
	}

	incResult, err := jsonpatch.ApplyPatch(numberDoc, incPatch)
	if err != nil {
		log.Fatalf("Inc operation failed: %v", err)
	}
	fmt.Printf("Inc operation: 42 + 8 = %v\n", incResult.Doc)

	// Test flip operation on root boolean
	boolDoc := true
	flipPatch := []jsonpatch.Operation{
		{
			Op:   "flip",
			Path: "",  // Root path for primitive value
		},
	}

	flipResult, err := jsonpatch.ApplyPatch(boolDoc, flipPatch)
	if err != nil {
		log.Fatalf("Flip operation failed: %v", err)
	}
	fmt.Printf("Flip operation: true -> %v\n", flipResult.Doc)
}