// Package testutils provides shared utilities for JSON Patch testing.
package testutils

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

// ApplyOperation applies a single operation to a document.
func ApplyOperation(t *testing.T, doc interface{}, operation jsonpatch.Operation) interface{} {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	options := jsonpatch.WithMutate(true)
	result, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyOperationWithError applies an operation expecting it to fail.
func ApplyOperationWithError(t *testing.T, doc interface{}, operation jsonpatch.Operation) error {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	options := jsonpatch.WithMutate(true)
	_, err := jsonpatch.ApplyPatch(doc, patch, options)
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
	return err
}

// ApplyOperations applies multiple operations to a document.
func ApplyOperations(t *testing.T, doc interface{}, operations []jsonpatch.Operation) interface{} {
	t.Helper()
	options := jsonpatch.WithMutate(true)
	result, err := jsonpatch.ApplyPatch(doc, operations, options)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyOperationsWithError applies multiple operations expecting them to fail.
func ApplyOperationsWithError(t *testing.T, doc interface{}, operations []jsonpatch.Operation) error {
	t.Helper()
	options := jsonpatch.WithMutate(true)
	_, err := jsonpatch.ApplyPatch(doc, operations, options)
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
	return err
}

// ApplyInternalOp applies a single internal.Operation to a document.
func ApplyInternalOp(t *testing.T, doc interface{}, op internal.Operation) interface{} {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyInternalOpWithError applies an internal.Operation expecting it to fail.
func ApplyInternalOpWithError(t *testing.T, doc interface{}, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
}

// ApplyInternalOps applies multiple internal.Operations to a document.
func ApplyInternalOps(t *testing.T, doc interface{}, ops []internal.Operation) interface{} {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyInternalOpsWithError applies multiple internal.Operations expecting them to fail.
func ApplyInternalOpsWithError(t *testing.T, doc interface{}, ops []internal.Operation) {
	t.Helper()
	_, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
}

// TestCase represents a single operation test case.
type TestCase struct {
	Name      string              // Test case name
	Doc       interface{}         // Input document
	Operation jsonpatch.Operation // Operation to apply
	Expected  interface{}         // Expected result
	WantErr   bool                // Whether operation should fail
	Comment   string              // Comment or source information
}

// MultiOperationTestCase represents a multi-operation test case.
type MultiOperationTestCase struct {
	Name       string                // Test case name
	Doc        interface{}           // Input document
	Operations []jsonpatch.Operation // Operations to apply
	Expected   interface{}           // Expected result
	WantErr    bool                  // Whether operations should fail
	Comment    string                // Comment or source information
}

// RunTestCase executes a single test case.
func RunTestCase(t *testing.T, tc TestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		t.Parallel()
		if tc.WantErr {
			_ = ApplyOperationWithError(t, tc.Doc, tc.Operation)
		} else {
			result := ApplyOperation(t, tc.Doc, tc.Operation)
			if diff := cmp.Diff(tc.Expected, result); diff != "" {
				t.Fatalf("ApplyPatch() mismatch (-want +got):\n%s", diff)
			}
		}
	})
}

// RunMultiOperationTestCase executes a multi-operation test case.
func RunMultiOperationTestCase(t *testing.T, tc MultiOperationTestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		t.Parallel()
		if tc.WantErr {
			_ = ApplyOperationsWithError(t, tc.Doc, tc.Operations)
		} else {
			result := ApplyOperations(t, tc.Doc, tc.Operations)
			if diff := cmp.Diff(tc.Expected, result); diff != "" {
				t.Fatalf("ApplyPatch() mismatch (-want +got):\n%s", diff)
			}
		}
	})
}

// RunTestCases executes multiple test cases.
func RunTestCases(t *testing.T, testCases []TestCase) {
	t.Helper()
	for _, tc := range testCases {
		RunTestCase(t, tc)
	}
}

// RunMultiOperationTestCases executes multiple multi-operation test cases.
func RunMultiOperationTestCases(t *testing.T, testCases []MultiOperationTestCase) {
	t.Helper()
	for _, tc := range testCases {
		RunMultiOperationTestCase(t, tc)
	}
}
