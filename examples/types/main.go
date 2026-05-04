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

var userPatch = []jsonpatch.Operation{
	{Op: "add", Path: "/email", Value: "john@example.com"},
}

var jsonPatch = []jsonpatch.Operation{
	{Op: "replace", Path: "/age", Value: 26},
}

func main() {
	if err := run(userPatch, jsonPatch); err != nil {
		return
	}
}

func run(userPatch, jsonPatch []jsonpatch.Operation) error {
	user := User{Name: "John", Age: 30}

	result, err := jsonpatch.ApplyPatch(user, userPatch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("Updated user: %+v\n", result.Doc)

	jsonStr := `{"name":"Alice","age":25}`

	jsonResult, err := jsonpatch.ApplyPatch(jsonStr, jsonPatch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("Updated JSON: %s\n", jsonResult.Doc)
	return nil
}
