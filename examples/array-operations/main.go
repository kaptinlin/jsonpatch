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
			"op":    "add",
			"path":  "/users/-",
			"value": map[string]interface{}{"id": 4, "name": "David"},
		},

		// Insert at beginning
		{
			"op":    "add",
			"path":  "/tags/0",
			"value": "patch",
		},

		// Replace array element
		{
			"op":    "replace",
			"path":  "/users/1/name",
			"value": "Bobby",
		},

		// Remove array element
		{
			"op":   "remove",
			"path": "/users/2",
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
