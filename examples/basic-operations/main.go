package main

import (
	"encoding/json"
	"fmt"
	"log"

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
	original, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
		// Test: verify current value
		{
			"op":    "test",
			"path":  "/name",
			"value": "John",
		},

		// Replace: update existing field
		{
			"op":    "replace",
			"path":  "/email",
			"value": "new@example.com",
		},

		// Add: create new field
		{
			"op":    "add",
			"path":  "/city",
			"value": "New York",
		},

		// Remove: delete field
		{
			"op":   "remove",
			"path": "/temp",
		},
	}

	options := jsonpatch.ApplyPatchOptions{Mutate: false}
	result, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter operations:")
	updated, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Println(string(updated))
}
