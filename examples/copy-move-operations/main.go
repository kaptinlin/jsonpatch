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

	doc := map[string]any{
		"user": map[string]any{
			"name":  "John Doe",
			"email": "john@example.com",
		},
		"settings": map[string]any{
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
			Op:   "copy",
			From: "/user/email",
			Path: "/backup_email",
		},

		// Move theme to user preferences
		{
			Op:   "move",
			From: "/settings/theme",
			Path: "/user/theme",
		},

		// Remove temporary data
		{
			Op:   "remove",
			Path: "/temp",
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
