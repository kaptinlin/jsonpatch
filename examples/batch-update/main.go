// Package main demonstrates batch update operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Batch Update ===")

	// System with multiple servers
	system := map[string]any{
		"servers": []any{
			map[string]any{"id": "srv1", "status": "inactive", "version": "1.0"},
			map[string]any{"id": "srv2", "status": "inactive", "version": "1.0"},
			map[string]any{"id": "srv3", "status": "inactive", "version": "1.0"},
		},
		"updated": 0,
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(system)
	fmt.Println(string(original))

	// Build batch update patch
	// 3 servers * 2 operations (status + version) + 1 counter update = 7 operations
	patch := make([]jsonpatch.Operation, 0, 7)

	// Update all servers in batch
	for i := range 3 {
		serverPath := fmt.Sprintf("/servers/%d", i)

		// Update status
		patch = append(patch, jsonpatch.Operation{
			Op:    "replace",
			Path:  serverPath + "/status",
			Value: "updated",
		})

		// Update version
		patch = append(patch, jsonpatch.Operation{
			Op:    "replace",
			Path:  serverPath + "/version",
			Value: "2.0",
		})
	}

	// Update counter
	patch = append(patch, jsonpatch.Operation{
		Op:    "replace",
		Path:  "/updated",
		Value: 3,
	})

	fmt.Printf("\nApplying %d operations in batch...\n", len(patch))

	result, err := jsonpatch.ApplyPatch(system, patch)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter batch update:")
	updated, _ := json.Marshal(result.Doc)
	fmt.Println(string(updated))
}
