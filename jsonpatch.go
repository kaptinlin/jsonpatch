// Package jsonpatch provides JSON Patch (RFC 6902) and JSON Predicate operations.
//
// Implements JSON mutation operations including:
//   - JSON Patch (RFC 6902): Standard operations (add, remove, replace, move, copy, test)
//     https://tools.ietf.org/html/rfc6902
//   - JSON Predicate: Test operations (contains, defined, type, less, more, etc.)
//     https://tools.ietf.org/id/draft-snell-json-test-01.html
//   - Extended operations: Additional operations (flip, inc, str_ins, str_del, split, merge)
//
// Core API Functions:
//   - ApplyOp: Apply a single operation
//   - ApplyOps: Apply multiple operations
//   - ApplyPatch: Apply a JSON Patch to a document
//   - ValidateOperations: Validate an array of operations
//   - ValidateOperation: Validate a single operation
//
// Basic usage:
//
//	doc := map[string]interface{}{"name": "John", "age": 30}
//	patch := []Operation{
//		{"op": "replace", "path": "/name", "value": "Jane"},
//		{"op": "add", "path": "/email", "value": "jane@example.com"},
//	}
//	result, err := ApplyPatch(doc, patch, DefaultApplyPatchOptions())
//
// The library provides optimized performance for map[string]interface{} documents
// and supports both generic and type-specific APIs.
package jsonpatch

import (
	"errors"
	"fmt"

	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
)

// Operation application errors
var (
	ErrNoOperationDecoded = errors.New("no operation decoded")
)

// Error message templates
const (
	errOperationFailed       = "operation %d failed: %w"
	errOperationDecodeFailed = "operation %d decode failed: %w"
)

// ApplyOp applies a single operation to a document.
func ApplyOp(doc interface{}, operation internal.Op, mutate bool) (*internal.OpResult, error) {
	workingDoc := doc
	if !mutate {
		workingDoc = deepclone.Clone(doc)
	}

	opResult, err := operation.Apply(workingDoc)
	if err != nil {
		return nil, err
	}

	return &opResult, nil
}

// ApplyOps applies multiple operations to a document.
func ApplyOps(doc interface{}, operations []internal.Op, mutate bool) (*internal.PatchResult, error) {
	workingDoc := doc
	if !mutate {
		workingDoc = deepclone.Clone(doc)
	}

	results := make([]internal.OpResult, 0, len(operations))
	for i, operation := range operations {
		opResult, err := operation.Apply(workingDoc)
		if err != nil {
			return nil, fmt.Errorf(errOperationFailed, i, err)
		}
		workingDoc = opResult.Doc
		results = append(results, opResult)
	}

	return &internal.PatchResult{Doc: workingDoc, Res: results}, nil
}

// ApplyPatch applies a JSON Patch to a document.
func ApplyPatch(doc interface{}, patch []internal.Operation, options internal.ApplyPatchOptions) (*internal.PatchResult, error) {
	workingDoc := doc
	if !options.Mutate {
		workingDoc = deepclone.Clone(doc)
	}

	results := make([]internal.OpResult, 0, len(patch))

	// Use codec/json decoder to convert operations to Op instances
	decoder := json.NewDecoder(internal.JsonPatchOptions{
		CreateMatcher: options.JsonPatchOptions.CreateMatcher,
	})

	for i, operation := range patch {
		// Convert operation to Op instance using operationToOp equivalent
		opInstance, err := decoder.Decode([]map[string]interface{}{operation})
		if err != nil {
			return nil, fmt.Errorf(errOperationDecodeFailed, i, err)
		}
		if len(opInstance) == 0 {
			return nil, fmt.Errorf(errOperationFailed, i, ErrNoOperationDecoded)
		}

		// Apply operation
		opResult, err := opInstance[0].Apply(workingDoc)
		if err != nil {
			return nil, fmt.Errorf(errOperationFailed, i, err)
		}
		workingDoc = opResult.Doc
		results = append(results, opResult)
	}

	return &internal.PatchResult{Doc: workingDoc, Res: results}, nil
}
