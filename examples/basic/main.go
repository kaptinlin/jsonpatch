// Package main demonstrates basic JSON Patch operations.
package main

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	doc := map[string]any{"name": "John", "age": 30}

	// Create patch operations
	patch := []jsonpatch.Operation{
		{Op: "add", Path: "/email", Value: "john@example.com"},
		{Op: "replace", Path: "/name", Value: "Jane"},
	}

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("After operations: %+v\n", result.Doc)
}
