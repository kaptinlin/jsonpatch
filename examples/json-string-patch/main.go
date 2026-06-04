// Package main demonstrates JSON string patching operations.
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
	fmt.Println("📝 JSON Patch for String Data")
	fmt.Println("=============================")

	// Example 1: Basic JSON string patching
	fmt.Println("\n📝 Example 1: Basic JSON String Operations")
	demoBasicJSONStringPatch()

	// Example 2: API response patching
	fmt.Println("\n🌐 Example 2: API Response Modification")
	demoAPIResponsePatch()

	// Example 3: Configuration updates
	fmt.Println("\n⚙️  Example 3: Configuration Updates")
	demoConfigurationPatch()
}

// demoBasicJSONStringPatch demonstrates basic JSON string patching
func demoBasicJSONStringPatch() {
	// Original JSON data as string
	jsonStr := `{"name":"Charlie Brown","age":40,"tags":["cto","founder"],"active":true}`

	fmt.Printf("Before: %s\n", prettyJSONString(jsonStr))

	// Define patch operations
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/name", Value: "Charles Brown"},
		{Op: "add", Path: "/company", Value: "TechStart Inc"},
		{Op: "add", Path: "/tags/-", Value: "entrepreneur"},
		{Op: "replace", Path: "/age", Value: 41},
	}

	result, err := applyCompiled(jsonpatch.JSONText(jsonStr), patch)
	if err != nil {
		log.Fatal("Failed to patch JSON string:", err)
	}

	fmt.Printf("After:  %s\n", prettyJSONString(string(result.Doc)))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Steps))
}

// demoAPIResponsePatch demonstrates patching API response data
func demoAPIResponsePatch() {
	// Simulated API response
	apiResponse := `{
		"status": "success",
		"data": {
			"user": {
				"id": 123,
				"username": "johndoe",
				"email": "john@example.com",
				"preferences": {
					"theme": "dark",
					"language": "en"
				}
			},
			"permissions": ["read", "write"]
		},
		"timestamp": "2024-01-15T10:30:00Z"
	}`

	fmt.Printf("Before:\n%s\n", prettyJSONString(apiResponse))

	// Update user preferences and permissions
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/data/user/email", Value: "john.doe@newcompany.com"},
		{Op: "replace", Path: "/data/user/preferences/theme", Value: "light"},
		{Op: "add", Path: "/data/user/preferences/notifications", Value: true},
		{Op: "add", Path: "/data/permissions/-", Value: "admin"},
		{Op: "replace", Path: "/timestamp", Value: "2024-01-15T11:00:00Z"},
	}

	result, err := applyCompiled(jsonpatch.JSONText(apiResponse), patch)
	if err != nil {
		log.Fatal("Failed to patch API response:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSONString(string(result.Doc)))
	fmt.Printf("✅ API response updated: %d operations\n", len(result.Steps))
}

// demoConfigurationPatch demonstrates configuration file updates
func demoConfigurationPatch() {
	// Configuration JSON string
	configStr := `{
		"server": {
			"host": "localhost",
			"port": 8080,
			"ssl": false
		},
		"database": {
			"host": "db.example.com",
			"port": 5432,
			"name": "myapp"
		},
		"features": {
			"logging": true,
			"metrics": false,
			"cache": true
		}
	}`

	fmt.Printf("Before:\n%s\n", prettyJSONString(configStr))

	// Configuration updates
	patch := []jsoncodec.Operation{
		{Op: "replace", Path: "/server/host", Value: "0.0.0.0"},
		{Op: "replace", Path: "/server/ssl", Value: true},
		{Op: "add", Path: "/server/ssl_cert", Value: "/etc/ssl/cert.pem"},
		{Op: "replace", Path: "/features/metrics", Value: true},
		{Op: "add", Path: "/features/monitoring", Value: true},
		{Op: "add", Path: "/redis", Value: map[string]any{
			"host": "redis.example.com",
			"port": 6379,
		}},
	}

	result, err := applyCompiled(jsonpatch.JSONText(configStr), patch)
	if err != nil {
		log.Fatal("Failed to patch configuration:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSONString(string(result.Doc)))
	fmt.Printf("✅ Configuration updated: %d changes\n", len(result.Steps))
}

func applyCompiled[T jsonpatch.Document](doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	patch, err := jsonpatch.CompileOperations(operations)
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
}

// prettyJSONString formats a JSON string for better readability
func prettyJSONString(jsonStr string) string {
	var obj any
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return jsonStr // Return original if parsing fails
	}

	pretty, err := json.Marshal(obj, jsontext.Multiline(true))
	if err != nil {
		return jsonStr // Return original if formatting fails
	}

	return string(pretty)
}
