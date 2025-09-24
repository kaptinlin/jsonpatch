// Package main demonstrates JSON bytes patching operations.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("📄 JSON Patch for Byte Data")
	fmt.Println("===========================")

	// Example 1: Basic JSON bytes patching
	fmt.Println("\n📝 Example 1: Basic JSON Bytes Operations")
	demoBasicJSONBytesPatch()

	// Example 2: Complex JSON document patching
	fmt.Println("\n🔧 Example 2: Complex JSON Document")
	demoComplexJSONBytesPatch()

	// Example 3: Array operations in JSON bytes
	fmt.Println("\n📋 Example 3: Array Operations")
	demoJSONBytesArrayOperations()
}

// demoBasicJSONBytesPatch demonstrates basic JSON bytes patching
func demoBasicJSONBytesPatch() {
	// Original JSON data as bytes
	jsonData := []byte(`{
		"name": "Bob Johnson",
		"age": 35,
		"active": true
	}`)

	fmt.Printf("Before:\n%s\n", prettyJSON(jsonData))

	// Define patch operations
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Robert Johnson"},
		{Op: "add", Path: "/email", Value: "bob.johnson@company.com"},
		{Op: "replace", Path: "/age", Value: 36},
		{Op: "add", Path: "/department", Value: "Engineering"},
	}

	// Apply patch - returns []byte
	result, err := jsonpatch.ApplyPatch(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch JSON:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Res))
}

// demoComplexJSONBytesPatch demonstrates complex JSON document patching
func demoComplexJSONBytesPatch() {
	// Complex JSON document
	jsonData := []byte(`{
		"user": {
			"name": "Alice Smith",
			"profile": {
				"bio": "Software Engineer",
				"skills": ["Go", "JavaScript", "Python"]
			}
		},
		"settings": {
			"theme": "dark",
			"notifications": true
		}
	}`)

	fmt.Printf("Before:\n%s\n", prettyJSON(jsonData))

	// Complex nested operations
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/user/name", Value: "Alice Johnson"},
		{Op: "add", Path: "/user/profile/skills/-", Value: "Rust"},
		{Op: "replace", Path: "/user/profile/bio", Value: "Senior Software Engineer"},
		{Op: "add", Path: "/user/email", Value: "alice.johnson@tech.com"},
		{Op: "replace", Path: "/settings/theme", Value: "light"},
	}

	result, err := jsonpatch.ApplyPatch(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch complex JSON:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Complex operations completed: %d\n", len(result.Res))
}

// demoJSONBytesArrayOperations demonstrates array operations
func demoJSONBytesArrayOperations() {
	// JSON with arrays
	jsonData := []byte(`{
		"team": "Development",
		"members": [
			{"name": "John", "role": "Lead"},
			{"name": "Jane", "role": "Developer"},
			{"name": "Bob", "role": "Designer"}
		],
		"technologies": ["Go", "React", "PostgreSQL"]
	}`)

	fmt.Printf("Before:\n%s\n", prettyJSON(jsonData))

	// Array-focused operations
	patch := []jsonpatch.Operation{
		{Op: "add", Path: "/members/-", Value: map[string]string{"name": "Alice", "role": "DevOps"}},
		{Op: "replace", Path: "/members/0/role", Value: "Tech Lead"},
		{Op: "add", Path: "/technologies/-", Value: "Docker"},
		{Op: "add", Path: "/technologies/0", Value: "Kubernetes"},
		{Op: "remove", Path: "/members/2"},
	}

	result, err := jsonpatch.ApplyPatch(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch arrays:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Array operations completed: %d\n", len(result.Res))
}

// prettyJSON formats JSON bytes for better readability
func prettyJSON(data []byte) string {
	var obj interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data) // Return original if parsing fails
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return string(data) // Return original if formatting fails
	}

	return string(pretty)
}
