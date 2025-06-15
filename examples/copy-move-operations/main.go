package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Copy and Move Operations ===")

	// Simple document for demonstration
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
	original, _ := json.MarshalIndent(doc, "", "  ")
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

	options := jsonpatch.ApplyPatchOptions{Mutate: false}
	result, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter copy/move:")
	updated, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Println(string(updated))
}
