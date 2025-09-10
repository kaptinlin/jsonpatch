// Package main demonstrates JSON string patching operations.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
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
	patch := []internal.Operation{
		{"op": "replace", "path": "/name", "value": "Charles Brown"},
		{"op": "add", "path": "/company", "value": "TechStart Inc"},
		{"op": "add", "path": "/tags/-", "value": "entrepreneur"},
		{"op": "replace", "path": "/age", "value": 41},
	}

	// Apply patch - returns string
	result, err := jsonpatch.ApplyPatch(jsonStr, patch)
	if err != nil {
		log.Fatal("Failed to patch JSON string:", err)
	}

	fmt.Printf("After:  %s\n", prettyJSONString(result.Doc))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Res))
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
	patch := []internal.Operation{
		{"op": "replace", "path": "/data/user/email", "value": "john.doe@newcompany.com"},
		{"op": "replace", "path": "/data/user/preferences/theme", "value": "light"},
		{"op": "add", "path": "/data/user/preferences/notifications", "value": true},
		{"op": "add", "path": "/data/permissions/-", "value": "admin"},
		{"op": "replace", "path": "/timestamp", "value": "2024-01-15T11:00:00Z"},
	}

	result, err := jsonpatch.ApplyPatch(apiResponse, patch)
	if err != nil {
		log.Fatal("Failed to patch API response:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSONString(result.Doc))
	fmt.Printf("✅ API response updated: %d operations\n", len(result.Res))
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
	patch := []internal.Operation{
		{"op": "replace", "path": "/server/host", "value": "0.0.0.0"},
		{"op": "replace", "path": "/server/ssl", "value": true},
		{"op": "add", "path": "/server/ssl_cert", "value": "/etc/ssl/cert.pem"},
		{"op": "replace", "path": "/features/metrics", "value": true},
		{"op": "add", "path": "/features/monitoring", "value": true},
		{"op": "add", "path": "/redis", "value": map[string]interface{}{
			"host": "redis.example.com",
			"port": 6379,
		}},
	}

	result, err := jsonpatch.ApplyPatch(configStr, patch)
	if err != nil {
		log.Fatal("Failed to patch configuration:", err)
	}

	fmt.Printf("After:\n%s\n", prettyJSONString(result.Doc))
	fmt.Printf("✅ Configuration updated: %d changes\n", len(result.Res))
}

// prettyJSONString formats a JSON string for better readability
func prettyJSONString(jsonStr string) string {
	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		return jsonStr // Return original if parsing fails
	}

	pretty, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return jsonStr // Return original if formatting fails
	}

	return string(pretty)
}
