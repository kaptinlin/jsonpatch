// Package testutils provides shared utilities for JSON Patch testing.
package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

// ApplyOperation applies a single operation to a document.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyOperation(t *testing.T, doc any, operation jsonpatch.Operation) any {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	result, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyOperationWithError applies an operation expecting it to fail.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyOperationWithError(t *testing.T, doc any, operation jsonpatch.Operation) error {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	_, err := jsonpatch.ApplyPatch(doc, patch, jsonpatch.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
	return err
}

// ApplyOperations applies multiple operations to a document.
func ApplyOperations(t *testing.T, doc any, operations []jsonpatch.Operation) any {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, operations, jsonpatch.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyOperationsWithError applies multiple operations expecting them to fail.
func ApplyOperationsWithError(t *testing.T, doc any, operations []jsonpatch.Operation) error {
	t.Helper()
	_, err := jsonpatch.ApplyPatch(doc, operations, jsonpatch.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
	return err
}

// ApplyInternalOp applies a single internal.Operation to a document.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyInternalOp(t *testing.T, doc any, op internal.Operation) any {
	t.Helper()
	patch := []internal.Operation{op}
	result, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyInternalOpWithError applies an internal.Operation expecting it to fail.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyInternalOpWithError(t *testing.T, doc any, op internal.Operation) {
	t.Helper()
	patch := []internal.Operation{op}
	_, err := jsonpatch.ApplyPatch(doc, patch, internal.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
}

// ApplyInternalOps applies multiple internal.Operations to a document.
func ApplyInternalOps(t *testing.T, doc any, ops []internal.Operation) any {
	t.Helper()
	result, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v, want nil", err)
	}
	return result.Doc
}

// ApplyInternalOpsWithError applies multiple internal.Operations expecting them to fail.
func ApplyInternalOpsWithError(t *testing.T, doc any, ops []internal.Operation) {
	t.Helper()
	_, err := jsonpatch.ApplyPatch(doc, ops, internal.WithMutate(true))
	if err == nil {
		t.Fatalf("ApplyPatch() error = nil, want error")
	}
}

// TestCase represents a single operation test case.
type TestCase struct {
	Name      string
	Doc       any
	Operation jsonpatch.Operation
	Expected  any
	WantErr   bool
	Comment   string
}

// MultiOperationTestCase represents a multi-operation test case.
type MultiOperationTestCase struct {
	Name       string
	Doc        any
	Operations []jsonpatch.Operation
	Expected   any
	WantErr    bool
	Comment    string
}

// RunTestCase executes a single test case.
//
//nolint:gocritic // Table-driven tests are assembled by value for readability.
func RunTestCase(t *testing.T, tc TestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		t.Parallel()
		if tc.WantErr {
			_ = ApplyOperationWithError(t, tc.Doc, tc.Operation)
		} else {
			result := ApplyOperation(t, tc.Doc, tc.Operation)
			assert.Equal(t, tc.Expected, result, "ApplyPatch() result mismatch")
		}
	})
}

// RunMultiOperationTestCase executes a multi-operation test case.
//
//nolint:gocritic // Table-driven tests are assembled by value for readability.
func RunMultiOperationTestCase(t *testing.T, tc MultiOperationTestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		t.Parallel()
		if tc.WantErr {
			_ = ApplyOperationsWithError(t, tc.Doc, tc.Operations)
		} else {
			result := ApplyOperations(t, tc.Doc, tc.Operations)
			assert.Equal(t, tc.Expected, result, "ApplyPatch() result mismatch")
		}
	})
}

// RunTestCases executes multiple test cases.
func RunTestCases(t *testing.T, testCases []TestCase) {
	t.Helper()
	for i := range testCases {
		RunTestCase(t, testCases[i])
	}
}

// RunMultiOperationTestCases executes multiple multi-operation test cases.
func RunMultiOperationTestCases(t *testing.T, testCases []MultiOperationTestCase) {
	t.Helper()
	for i := range testCases {
		RunMultiOperationTestCase(t, testCases[i])
	}
}
