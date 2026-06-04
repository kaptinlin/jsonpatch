// Package main demonstrates type operations with JSON Patch.
package main

import (
	"fmt"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/kaptinlin/jsonpatch"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age"`
}

var userPatch = []jsoncodec.Operation{
	{Op: "add", Path: "/email", Value: "john@example.com"},
}

var jsonPatch = []jsoncodec.Operation{
	{Op: "replace", Path: "/age", Value: 26},
}

func main() {
	if err := run(userPatch, jsonPatch); err != nil {
		return
	}
}

func run(userPatch, jsonPatch []jsoncodec.Operation) error {
	user := User{Name: "John", Age: 30}

	compiled, err := jsonpatch.CompileOperations(userPatch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	result, err := jsonpatch.Apply(compiled, user)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("Updated user: %+v\n", result.Doc)

	jsonStr := `{"name":"Alice","age":25}`

	compiled, err = jsonpatch.CompileOperations(jsonPatch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	jsonResult, err := jsonpatch.Apply(compiled, jsonpatch.JSONText(jsonStr))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("Updated JSON: %s\n", jsonResult.Doc)
	return nil
}
