// Package main demonstrates conditional operations using JSON Patch.
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
			Op:    "test",
			Path:  "/status",
			Value: "active",
		},
		// Test: verify sufficient balance
		{
			Op:    "test",
			Path:  "/balance",
			Value: 1000.0,
		},
		// Perform withdrawal
		{
			Op:    "replace",
			Path:  "/balance",
			Value: 700.0,
		},
		// Increment version
		{
			Op:   "inc",
			Path: "/version",
			Inc:  1,
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
			Op:    "test",
			Path:  "/balance",
			Value: 2000.0, // Wrong amount
		},
		{
			Op:    "replace",
			Path:  "/balance",
			Value: 0.0,
		},
	}

	_, err = jsonpatch.ApplyPatch(account, failPatch)
	if err != nil {
		fmt.Printf("Expected failure: %v\n", err)
	}
}
