// Package main demonstrates basic JSON Patch operations.
package main

import (
	"fmt"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/kaptinlin/jsonpatch"
)

var patch = []jsoncodec.Operation{
	{Op: "add", Path: "/email", Value: "john@example.com"},
	{Op: "replace", Path: "/name", Value: "Jane"},
}

func main() {
	if err := run(patch); err != nil {
		return
	}
}

func run(patch []jsoncodec.Operation) error {
	doc := map[string]any{"name": "John", "age": 30}

	compiled, err := jsonpatch.CompileOperations(patch)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}
	result, err := jsonpatch.Apply(compiled, doc)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err
	}

	fmt.Printf("After operations: %+v\n", result.Doc)
	return nil
}
