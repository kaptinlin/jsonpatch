// Package main demonstrates struct patching operations using JSON Patch.
package main

import (
	"fmt"
	"log"

	"encoding/json"

	"github.com/kaptinlin/jsonpatch"
)

// User represents a user profile with JSON tags
type User struct {
	Name    string   `json:"name"`
	Email   string   `json:"email,omitempty"`
	Age     int      `json:"age"`
	Tags    []string `json:"tags"`
	Active  bool     `json:"active"`
	Profile *Profile `json:"profile,omitempty"`
}

// Profile represents nested profile information
type Profile struct {
	Bio     string `json:"bio,omitempty"`
	Website string `json:"website,omitempty"`
}

func main() {
	fmt.Println("🏗️  JSON Patch for Go Structs")
	fmt.Println("=============================")

	// Example 1: Basic struct patching
	fmt.Println("\n📝 Example 1: Basic User Updates")
	demoBasicStructPatch()

	// Example 2: Complex nested struct operations
	fmt.Println("\n🔧 Example 2: Complex Nested Operations")
	demoComplexStructPatch()

	// Example 3: Array operations on struct fields
	fmt.Println("\n📋 Example 3: Array Field Operations")
	demoArrayFieldOperations()
}

// demoBasicStructPatch demonstrates basic struct patching
func demoBasicStructPatch() {
	// Original user
	user := User{
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
		Tags:   []string{"developer", "golang"},
		Active: true,
	}

	fmt.Printf("Before: %+v\n", user)

	// Define patch operations
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Jane Smith"},
		{Op: "replace", Path: "/age", Value: 28},
		{Op: "replace", Path: "/email", Value: "jane.smith@example.com"},
	}

	// Apply patch - preserves struct type
	result, err := jsonpatch.ApplyPatch(user, patch)
	if err != nil {
		log.Fatal("Failed to patch user:", err)
	}

	fmt.Printf("After:  %+v\n", result.Doc)
	fmt.Printf("✅ Original unchanged: %+v\n", user)
}

// demoComplexStructPatch demonstrates complex nested struct operations
func demoComplexStructPatch() {
	// Complex user with nested profile
	user := User{
		Name:   "Diana Prince",
		Email:  "diana@themyscira.com",
		Age:    500,
		Tags:   []string{"warrior", "princess"},
		Active: true,
		Profile: &Profile{
			Bio:     "Amazon warrior princess",
			Website: "https://wonderwoman.dc.com",
		},
	}

	// Print original (pretty formatted)
	originalJSON, _ := json.MarshalIndent(user, "", "  ")
	fmt.Printf("Before:\n%s\n", string(originalJSON))

	// Complex patch operations including nested paths
	patch := []jsonpatch.Operation{
		{Op: "replace", Path: "/name", Value: "Wonder Woman"},
		{Op: "replace", Path: "/profile/bio", Value: "DC Comics superhero and Amazon warrior"},
		{Op: "add", Path: "/headquarters", Value: "Hall of Justice"},
	}

	// Apply patch with immutable option (default)
	result, err := jsonpatch.ApplyPatch(user, patch, jsonpatch.WithMutate(false))
	if err != nil {
		log.Fatal("Failed to patch complex user:", err)
	}

	// Print result (pretty formatted)
	patchedJSON, _ := json.MarshalIndent(result.Doc, "", "  ")
	fmt.Printf("After:\n%s\n", string(patchedJSON))
	fmt.Printf("✅ Operations applied: %d\n", len(result.Res))
}

// demoArrayFieldOperations demonstrates operations on array fields
func demoArrayFieldOperations() {
	user := User{
		Name: "Alice Cooper",
		Age:  25,
		Tags: []string{"designer", "ui"},
	}

	fmt.Printf("Before: Tags = %v\n", user.Tags)

	// Array operations on Tags field
	patch := []jsonpatch.Operation{
		{Op: "add", Path: "/tags/-", Value: "ux"},          // Append
		{Op: "add", Path: "/tags/0", Value: "senior"},      // Insert at beginning
		{Op: "replace", Path: "/tags/2", Value: "product"}, // Replace middle
	}

	result, err := jsonpatch.ApplyPatch(user, patch)
	if err != nil {
		log.Fatal("Failed to patch tags:", err)
	}

	fmt.Printf("After:  Tags = %v\n", result.Doc.Tags)
	fmt.Printf("✅ Tag operations completed successfully\n")
}
