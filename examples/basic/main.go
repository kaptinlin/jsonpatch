// Package main demonstrates basic JSON Patch operations.
package main

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch"
)

var patch = []jsonpatch.Operation{
	{Op: "add", Path: "/email", Value: "john@example.com"},
	{Op: "replace", Path: "/name", Value: "Jane"},
}

func main() {
	if err := run(patch); err != nil {
		return
	}
}

func run(patch []jsonpatch.Operation) error {
	doc := map[string]any{"name": "John", "age": 30}

	result, err := jsonpatch.ApplyPatch(doc, patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("After operations: %+v\n", result.Doc)
	return nil
}
