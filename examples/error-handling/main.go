package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Error Handling ===")

	// Test document
	doc := map[string]interface{}{
		"balance": 1000.0,
		"status":  "active",
		"version": 1,
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(doc)
	fmt.Println(string(original))

	// Example 1: Successful validation
	fmt.Println("\n--- Successful Validation ---")

	successPatch := []jsonpatch.Operation{
		// Test current values
		{
			"op":    "test",
			"path":  "/status",
			"value": "active",
		},
		{
			"op":    "test",
			"path":  "/balance",
			"value": 1000.0,
		},
		// Make changes
		{
			"op":    "replace",
			"path":  "/balance",
			"value": 800.0,
		},
		{
			"op":   "inc",
			"path": "/version",
			"inc":  1,
		},
	}

	result, err := jsonpatch.ApplyPatch(doc, successPatch)
	if err != nil {
		log.Printf("Patch failed: %v", err)
	} else {
		fmt.Println("Success:")
		updated, _ := json.Marshal(result.Doc)
		fmt.Println(string(updated))
	}

	// Example 2: Failed test condition
	fmt.Println("\n--- Failed Test Condition ---")

	failPatch := []jsonpatch.Operation{
		{
			"op":    "test",
			"path":  "/balance",
			"value": 2000.0, // Wrong value
		},
		{
			"op":    "replace",
			"path":  "/balance",
			"value": 0.0,
		},
	}

	_, err = jsonpatch.ApplyPatch(doc, failPatch)
	if err != nil {
		fmt.Printf("Expected failure: %v\n", err)
	}

	// Example 3: Invalid path
	fmt.Println("\n--- Invalid Path ---")

	invalidPatch := []jsonpatch.Operation{
		{
			"op":    "replace",
			"path":  "/nonexistent",
			"value": "value",
		},
	}

	_, err = jsonpatch.ApplyPatch(doc, invalidPatch)
	if err != nil {
		fmt.Printf("Path error: %v\n", err)
	}
}
