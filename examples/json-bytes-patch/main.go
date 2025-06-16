package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
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
	patch := []internal.Operation{
		{"op": "replace", "path": "/name", "value": "Robert Johnson"},
		{"op": "add", "path": "/email", "value": "bob.johnson@company.com"},
		{"op": "replace", "path": "/age", "value": 36},
		{"op": "add", "path": "/department", "value": "Engineering"},
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
	patch := []internal.Operation{
		{"op": "replace", "path": "/user/name", "value": "Alice Johnson"},
		{"op": "add", "path": "/user/profile/skills/-", "value": "Rust"},
		{"op": "replace", "path": "/user/profile/bio", "value": "Senior Software Engineer"},
		{"op": "add", "path": "/user/email", "value": "alice.johnson@tech.com"},
		{"op": "replace", "path": "/settings/theme", "value": "light"},
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
	patch := []internal.Operation{
		{"op": "add", "path": "/members/-", "value": map[string]string{"name": "Alice", "role": "DevOps"}},
		{"op": "replace", "path": "/members/0/role", "value": "Tech Lead"},
		{"op": "add", "path": "/technologies/-", "value": "Docker"},
		{"op": "add", "path": "/technologies/0", "value": "Kubernetes"},
		{"op": "remove", "path": "/members/2"},
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
