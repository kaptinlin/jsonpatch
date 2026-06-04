// Package main demonstrates JSON bytes patching operations.
package main

import (
	"fmt"
	"log"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/go-json-experiment/json"
	"github.com/go-json-experiment/json/jsontext"

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
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Robert Johnson"},
		{Op: "add", Path: "/email", Value: "bob.johnson@company.com"},
		{Op: "replace", Path: "/age", Value: 36},
		{Op: "add", Path: "/department", Value: "Engineering"},
	}

	result, err := applyCompiled(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch JSON:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Steps))
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
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/user/name", Value: "Alice Johnson"},
		{Op: "add", Path: "/user/profile/skills/-", Value: "Rust"},
		{Op: "replace", Path: "/user/profile/bio", Value: "Senior Software Engineer"},
		{Op: "add", Path: "/user/email", Value: "alice.johnson@tech.com"},
		{Op: "replace", Path: "/settings/theme", Value: "light"},
	}

	result, err := applyCompiled(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch complex JSON:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Complex operations completed: %d\n", len(result.Steps))
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
	patch := []jsoncodec.Operation{
		{Op: "add", Path: "/members/-", Value: map[string]string{"name": "Alice", "role": "DevOps"}},
		{Op: "replace", Path: "/members/0/role", Value: "Tech Lead"},
		{Op: "add", Path: "/technologies/-", Value: "Docker"},
		{Op: "add", Path: "/technologies/0", Value: "Kubernetes"},
		{Op: "remove", Path: "/members/2"},
	}

	result, err := applyCompiled(jsonData, patch)
	if err != nil {
		log.Fatal("Failed to patch arrays:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSON(result.Doc))
	fmt.Printf("✅ Array operations completed: %d\n", len(result.Steps))
}

func applyCompiled[T jsonpatch.Document](doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	patch, err := jsonpatch.CompileOperations(operations)
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
}

// prettyJSON formats JSON bytes for better readability
func prettyJSON(data []byte) string {
	var obj any
	if err := json.Unmarshal(data, &obj); err != nil {
		return string(data) // Return original if parsing fails
	}

	pretty, err := json.Marshal(obj, jsontext.Multiline(true))
	if err != nil {
		return string(data) // Return original if formatting fails
	}

	return string(pretty)
}
