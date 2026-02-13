// Package main demonstrates error handling with JSON Patch operations.
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
	doc := map[string]any{
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
			Op:    "test",
			Path:  "/status",
			Value: "active",
		},
		{
			Op:    "test",
			Path:  "/balance",
			Value: 1000.0,
		},
		// Make changes
		{
			Op:    "replace",
			Path:  "/balance",
			Value: 800.0,
		},
		{
			Op:   "inc",
			Path: "/version",
			Inc:  1,
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
			Op:    "test",
			Path:  "/balance",
			Value: 2000.0, // Wrong value
		},
		{
			Op:    "replace",
			Path:  "/balance",
			Value: 0.0,
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
			Op:    "replace",
			Path:  "/nonexistent",
			Value: "value",
		},
	}

	_, err = jsonpatch.ApplyPatch(doc, invalidPatch)
	if err != nil {
		fmt.Printf("Path error: %v\n", err)
	}
}
