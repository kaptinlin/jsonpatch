// Package main demonstrates array operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Array Operations ===")

	// Document with arrays
	doc := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"id": 1, "name": "Alice"},
			map[string]interface{}{"id": 2, "name": "Bob"},
			map[string]interface{}{"id": 3, "name": "Charlie"},
		},
		"tags": []interface{}{"go", "json"},
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
		// Add to end of array using "-"
		{
			Op:    "add",
			Path:  "/users/-",
			Value: map[string]interface{}{"id": 4, "name": "David"},
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

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter array operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
