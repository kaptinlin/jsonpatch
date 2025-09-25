// Package main demonstrates map patching operations using JSON Patch.
package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kaptinlin/jsonpatch"
)

func main() {
	fmt.Println("🗺️  JSON Patch for Map Data")
	fmt.Println("===========================")

	// Example 1: Basic map patching
	fmt.Println("\n📝 Example 1: Basic Map Operations")
	demoBasicMapPatch()

	// Example 2: Nested map operations
	fmt.Println("\n🔧 Example 2: Nested Map Operations")
	demoNestedMapPatch()

	// Example 3: Dynamic data manipulation
	fmt.Println("\n⚡ Example 3: Dynamic Data Manipulation")
	demoDynamicMapPatch()
}

// demoBasicMapPatch demonstrates basic map patching
func demoBasicMapPatch() {
	// Document to modify
	doc := map[string]any{
		"name":   "Alice Cooper",
		"age":    25,
		"tags":   []any{"designer", "ui"},
		"active": true,
	}

	fmt.Printf("Before: %s\n", prettyMap(doc))

	// Define patch operations
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Alice Johnson"},
		{Op: "add", Path: "/email", Value: "alice@design.co"},
		{Op: "add", Path: "/tags/-", Value: "ux"},
		{Op: "replace", Path: "/age", Value: 26},
		{Op: "add", Path: "/department", Value: "Product Design"},
	}

	// Apply patch - returns map[string]any
	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatal("Failed to patch document:", err)
	}

	fmt.Printf("After:  %s\n", prettyMap(result.Doc))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Res))
}

// demoNestedMapPatch demonstrates nested map operations
func demoNestedMapPatch() {
	// Complex nested map structure
	doc := map[string]any{
		"user": map[string]any{
			"name": "Bob Smith",
			"profile": map[string]any{
				"bio":    "Full-stack developer",
				"skills": []any{"Go", "React", "PostgreSQL"},
			},
			"contact": map[string]any{
				"email": "bob@example.com",
				"phone": "+1-555-0123",
			},
		},
		"metadata": map[string]any{
			"created": "2024-01-01",
			"version": 1,
		},
	}

	fmt.Printf("Before:\n%s\n", prettyMap(doc))

	// Nested operations
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/user/name", Value: "Robert Smith"},
		{Op: "add", Path: "/user/profile/skills/-", Value: "Docker"},
		{Op: "replace", Path: "/user/profile/bio", Value: "Senior Full-stack Developer"},
		{Op: "add", Path: "/user/contact/linkedin", Value: "linkedin.com/in/robertsmith"},
		{Op: "replace", Path: "/metadata/version", Value: 2},
		{Op: "add", Path: "/metadata/updated", Value: "2024-01-15"},
	}

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		log.Fatal("Failed to patch nested document:", err)
	}

	fmt.Printf("After:\n%s\n", prettyMap(result.Doc))
	fmt.Printf("✅ Nested operations completed: %d\n", len(result.Res))
}

// demoDynamicMapPatch demonstrates dynamic data manipulation
func demoDynamicMapPatch() {
	// Dynamic data structure (like from API or database)
	doc := map[string]any{
		"products": []any{
			map[string]any{"id": 1, "name": "Laptop", "price": 999.99, "stock": 10},
			map[string]any{"id": 2, "name": "Mouse", "price": 29.99, "stock": 50},
			map[string]any{"id": 3, "name": "Keyboard", "price": 79.99, "stock": 25},
		},
		"store": map[string]any{
			"name":     "Tech Store",
			"location": "Downtown",
			"open":     true,
		},
		"stats": map[string]any{
			"total_products": 3,
			"total_value":    1109.97,
		},
	}

	fmt.Printf("Before:\n%s\n", prettyMap(doc))

	// Dynamic updates (price changes, stock updates, new products)
	patch := []jsonpatch.Operation{
		// Update product prices
		{Op: "replace", Path: "/products/0/price", Value: 899.99},
		{Op: "replace", Path: "/products/2/stock", Value: 30},

		// Add new product
		{Op: "add", Path: "/products/-", Value: map[string]any{
			"id": 4, "name": "Monitor", "price": 299.99, "stock": 15,
		}},

		// Update store info
		{Op: "add", Path: "/store/phone", Value: "+1-555-TECH"},
		{Op: "add", Path: "/store/website", Value: "techstore.com"},

		// Update stats
		{Op: "replace", Path: "/stats/total_products", Value: 4},
		{Op: "replace", Path: "/stats/total_value", Value: 1309.96},
		{Op: "add", Path: "/stats/last_updated", Value: "2024-01-15T10:30:00Z"},
	}

	result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(false))
	if err != nil {
		log.Fatal("Failed to patch dynamic data:", err)
	}

	fmt.Printf("After:\n%s\n", prettyMap(result.Doc))
	fmt.Printf("✅ Dynamic updates completed: %d operations\n", len(result.Res))

	// Show that original is unchanged
	fmt.Printf("✅ Original unchanged: total_products = %v\n",
		doc["stats"].(map[string]any)["total_products"])
}

// prettyMap formats a map for better readability
func prettyMap(m map[string]any) string {
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Sprintf("%+v", m) // Fallback to Go's default formatting
	}
	return string(data)
}
