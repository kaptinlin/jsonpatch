package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Conditional Operations ===")

	// Account data
	account := map[string]interface{}{
		"balance": 1000.0,
		"status":  "active",
		"version": 1,
	}

	fmt.Println("\nOriginal account:")
	original, _ := json.Marshal(account)
	fmt.Println(string(original))

	// Safe withdrawal: only if conditions are met
	patch := []jsonpatch.Operation{
		// Test: verify account is active
		{
			"op":    "test",
			"path":  "/status",
			"value": "active",
		},
		// Test: verify sufficient balance
		{
			"op":    "test",
			"path":  "/balance",
			"value": 1000.0,
		},
		// Perform withdrawal
		{
			"op":    "replace",
			"path":  "/balance",
			"value": 700.0,
		},
		// Increment version
		{
			"op":   "inc",
			"path": "/version",
			"inc":  1,
		},
	}

	fmt.Println("\nAttempting safe withdrawal...")
	result, err := jsonpatch.ApplyPatch(account, patch)
	if err != nil {
		log.Printf("Transaction failed: %v", err)
		return
	}

	fmt.Println("\nAfter successful transaction:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))

	// Demonstrate failed condition
	fmt.Println("\n--- Failed Condition Example ---")

	failPatch := []jsonpatch.Operation{
		// This test will fail - wrong balance
		{
			"op":    "test",
			"path":  "/balance",
			"value": 2000.0, // Wrong amount
		},
		{
			"op":    "replace",
			"path":  "/balance",
			"value": 0.0,
		},
	}

	_, err = jsonpatch.ApplyPatch(account, failPatch)
	if err != nil {
		fmt.Printf("Expected failure: %v\n", err)
	}
}
