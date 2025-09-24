// Package main demonstrates string operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== String Operations ===")

	// Document with text content
	doc := map[string]interface{}{
		"title":   "My Document",
		"content": "Hello world!",
		"tags":    "go,json",
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
		// Insert at beginning
		{
			Op:   "str_ins",
			Path: "/content",
			Pos:  0,
			Str:  "Welcome! ",
		},

		// Insert at end
		{
			Op:   "str_ins",
			Path: "/content",
			Pos:  21, // After "Welcome! Hello world!"
			Str:  " How are you?",
		},

		// Insert in middle of title
		{
			Op:   "str_ins",
			Path: "/title",
			Pos:  2, // After "My"
			Str:  " First",
		},

		// Add new tag
		{
			Op:   "str_ins",
			Path: "/tags",
			Pos:  7, // At end of "go,json"
			Str:  ",patch",
		},
	}

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter string operations:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
