package main

import (
	"encoding/json"
	"fmt"
	"log"

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
	original, _ := json.MarshalIndent(doc, "", "  ")
	fmt.Println(string(original))

	patch := []jsonpatch.Operation{
		// Insert at beginning
		{
			"op":   "str_ins",
			"path": "/content",
			"pos":  0,
			"str":  "Welcome! ",
		},

		// Insert at end
		{
			"op":   "str_ins",
			"path": "/content",
			"pos":  21, // After "Welcome! Hello world!"
			"str":  " How are you?",
		},

		// Insert in middle of title
		{
			"op":   "str_ins",
			"path": "/title",
			"pos":  2, // After "My"
			"str":  " First",
		},

		// Add new tag
		{
			"op":   "str_ins",
			"path": "/tags",
			"pos":  7, // At end of "go,json"
			"str":  ",patch",
		},
	}

	options := jsonpatch.ApplyPatchOptions{Mutate: false}
	result, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter string operations:")
	updated, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Println(string(updated))
}
