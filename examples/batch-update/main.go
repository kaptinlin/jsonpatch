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
	system := map[string]interface{}{
		"servers": []interface{}{
			map[string]interface{}{"id": "srv1", "status": "inactive", "version": "1.0"},
			map[string]interface{}{"id": "srv2", "status": "inactive", "version": "1.0"},
			map[string]interface{}{"id": "srv3", "status": "inactive", "version": "1.0"},
		},
		"updated": 0,
	}

	fmt.Println("\nOriginal:")
	original, _ := json.Marshal(system)
	fmt.Println(string(original))

	// Build batch update patch
	var patch []jsonpatch.Operation

	// Update all servers in batch
	for i := 0; i < 3; i++ {
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
