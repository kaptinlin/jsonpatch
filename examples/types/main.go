// Package main demonstrates type operations with JSON Patch.
package main

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age"`
}

func main() {
	// Struct operations using new Operation syntax
	user := User{Name: "John", Age: 30}
	patch := []jsonpatch.Operation{
		{Op: "add", Path: "/email", Value: "john@example.com"},
	}
	
	result, err := jsonpatch.ApplyPatch(user, patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Updated user: %+v\n", result.Doc)

	// JSON operations using new Operation syntax
	jsonStr := `{"name":"Alice","age":25}`
	jsonPatch := []jsonpatch.Operation{
		{Op: "replace", Path: "/age", Value: 26},
	}
	
	jsonResult, err := jsonpatch.ApplyPatch(jsonStr, jsonPatch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	
	fmt.Printf("Updated JSON: %s\n", jsonResult.Doc)
}