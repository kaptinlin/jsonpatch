// Package main demonstrates basic JSON Patch operations.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Basic Operations ===")

	// Original document
	doc := map[string]interface{}{
		"name":  "John",
		"age":   30,
		"email": "old@example.com",
		"temp":  "remove_me",
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
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
			Value: "new@example.com",
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

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
