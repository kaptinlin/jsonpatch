// Package main demonstrates copy and move operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Copy and Move Operations ===")

	doc := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
		},
		"settings": map[string]interface{}{
			"theme": "dark",
		},
		"temp": "temporary_data",
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
		// Copy email to create backup
		{
			"op":   "copy",
			"from": "/user/email",
			"path": "/backup_email",
		},

		// Move theme to user preferences
		{
			"op":   "move",
			"from": "/settings/theme",
			"path": "/user/theme",
		},

		// Remove temporary data
		{
			"op":   "remove",
			"path": "/temp",
		},
	}

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter copy/move:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
