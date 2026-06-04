// Package main demonstrates basic JSON Patch operations.
package main

import (
	"fmt"
	"log"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Basic Operations ===")

	// Document to modify
	doc := map[string]any{
		"name":  "John",
		"age":   30,
		"email": "original@example.com",
		"temp":  "remove_me",
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsoncodec.Operation{
		// Test: verify current value
		{
			Op:    "test",
			Path:  "/name",
			Value: "John",
		},

		// Replace: update existing field
		{
			Op:    "replace",
			Path:  "/email",
			Value: "updated@example.com",
		},

		// Add: create new field
		{
			Op:    "add",
			Path:  "/city",
			Value: "New York",
		},

		// Remove: delete field
		{
			Op:   "remove",
			Path: "/temp",
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

	fmt.Println("\nAfter operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
