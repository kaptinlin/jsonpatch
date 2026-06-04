// Package main demonstrates array operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Array Operations ===")

	// Document with arrays
	doc := map[string]any{
		"users": []any{
			map[string]any{"id": 1, "name": "Alice"},
			map[string]any{"id": 2, "name": "Bob"},
			map[string]any{"id": 3, "name": "Charlie"},
		},
		"tags": []any{"go", "json"},
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsoncodec.Operation{
		// Add to end of array using "-"
		{
			Op:    "add",
			Path:  "/users/-",
			Value: map[string]any{"id": 4, "name": "David"},
		},

		// Insert at beginning
		{
			Op:    "add",
			Path:  "/tags/0",
			Value: "patch",
		},

		// Replace array element
		{
			Op:    "replace",
			Path:  "/users/1/name",
			Value: "Bobby",
		},

		// Remove array element
		{
			Op:   "remove",
			Path: "/users/2",
		},
	}

	compiled, err := jsonpatch.CompileOperations(patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}
	result, err := jsonpatch.Apply(compiled, doc)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter array operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
