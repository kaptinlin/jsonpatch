package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("=== Batch Update ===")

	// System with multiple servers
	system := map[string]interface{}{
		"servers": []interface{}{
			map[string]interface{}{"id": "srv1", "status": "old", "version": "1.0"},
			map[string]interface{}{"id": "srv2", "status": "old", "version": "1.0"},
			map[string]interface{}{"id": "srv3", "status": "old", "version": "1.0"},
		},
		"updated": 0,
	}

	fmt.Println("\nOriginal:")
	original, _ := json.MarshalIndent(system, "", "  ")
	fmt.Println(string(original))

	// Build batch update patch
	var patch []jsonpatch.Operation

	// Update all servers in batch
	for i := 0; i < 3; i++ {
		serverPath := fmt.Sprintf("/servers/%d", i)

		// Update status
		patch = append(patch, jsonpatch.Operation{
			"op":    "replace",
			"path":  serverPath + "/status",
			"value": "updated",
		})

		// Update version
		patch = append(patch, jsonpatch.Operation{
			"op":    "replace",
			"path":  serverPath + "/version",
			"value": "2.0",
		})
	}

	// Update counter
	patch = append(patch, jsonpatch.Operation{
		"op":    "replace",
		"path":  "/updated",
		"value": 3,
	})

	fmt.Printf("\nApplying %d operations in batch...\n", len(patch))

	options := jsonpatch.ApplyPatchOptions{Mutate: false}
	result, err := jsonpatch.ApplyPatch(system, patch, options)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Println("\nAfter batch update:")
	updated, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Println(string(updated))
}
