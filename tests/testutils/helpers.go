// Package testutils provides shared utilities for JSON Patch testing.
package testutils

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"

	"github.com/stretchr/testify/assert"
)

// ApplyOperation applies a single operation to a document.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyOperation(t *testing.T, doc any, operation jsoncodec.Operation) any {
	t.Helper()
	result, err := compileAndApplyInPlace(t, doc, []jsoncodec.Operation{operation})
	if err != nil {
		t.Fatalf("ApplyInPlace() error = %v, want nil", err)
	}
	return result
}

// ApplyOperationWithError applies an operation expecting it to fail.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyOperationWithError(t *testing.T, doc any, operation jsoncodec.Operation) error {
	t.Helper()
	_, err := compileAndApplyInPlace(t, doc, []jsoncodec.Operation{operation})
	if err == nil {
		t.Fatalf("Compile/ApplyInPlace error = nil, want error")
	}
	return err
}

// ApplyOperations applies multiple operations to a document.
func ApplyOperations(t *testing.T, doc any, operations []jsoncodec.Operation) any {
	t.Helper()
	result, err := compileAndApplyInPlace(t, doc, operations)
	if err != nil {
		t.Fatalf("ApplyInPlace() error = %v, want nil", err)
	}
	return result
}

// ApplyOperationsWithError applies multiple operations expecting them to fail.
func ApplyOperationsWithError(t *testing.T, doc any, operations []jsoncodec.Operation) error {
	t.Helper()
	_, err := compileAndApplyInPlace(t, doc, operations)
	if err == nil {
		t.Fatalf("Compile/ApplyInPlace error = nil, want error")
	}
	return err
}

// ApplyOperationsResult compiles and applies JSON-shaped operations.
func ApplyOperationsResult[T jsonpatch.Document](t testing.TB, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	return jsonpatch.Apply(patch, doc)
}

// ApplyInternalOp applies a single internal.Operation to a document.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyInternalOp(t *testing.T, doc any, op internal.Operation) any {
	t.Helper()
	result, err := compileAndApplyInPlace(t, doc, []jsoncodec.Operation{jsoncodec.Operation(op)})
	if err != nil {
		t.Fatalf("ApplyInPlace() error = %v, want nil", err)
	}
	return result
}

// ApplyInternalOpWithError applies an internal.Operation expecting it to fail.
//
//nolint:gocritic // Tests pass operations by value for readability at call sites.
func ApplyInternalOpWithError(t *testing.T, doc any, op internal.Operation) {
	t.Helper()
	_, err := compileAndApplyInPlace(t, doc, []jsoncodec.Operation{jsoncodec.Operation(op)})
	if err == nil {
		t.Fatalf("Compile/ApplyInPlace error = nil, want error")
	}
}

// ApplyInternalOps applies multiple internal.Operations to a document.
func ApplyInternalOps(t *testing.T, doc any, ops []internal.Operation) any {
	t.Helper()
	operations := make([]jsoncodec.Operation, len(ops))
	for i := range ops {
		operations[i] = jsoncodec.Operation(ops[i])
	}
	result, err := compileAndApplyInPlace(t, doc, operations)
	if err != nil {
		t.Fatalf("ApplyInPlace() error = %v, want nil", err)
	}
	return result
}

// ApplyInternalOpsWithError applies multiple internal.Operations expecting them to fail.
func ApplyInternalOpsWithError(t *testing.T, doc any, ops []internal.Operation) {
	t.Helper()
	operations := make([]jsoncodec.Operation, len(ops))
	for i := range ops {
		operations[i] = jsoncodec.Operation(ops[i])
	}
	_, err := compileAndApplyInPlace(t, doc, operations)
	if err == nil {
		t.Fatalf("Compile/ApplyInPlace error = nil, want error")
	}
}

// ApplyInternalOperationsResult compiles and applies internal.Operation test values.
func ApplyInternalOperationsResult[T jsonpatch.Document](t testing.TB, doc T, ops []internal.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	operations := make([]jsoncodec.Operation, len(ops))
	for i := range ops {
		operations[i] = jsoncodec.Operation(ops[i])
	}
	return ApplyOperationsResult(t, doc, operations)
}

func compileAndApplyInPlace(t *testing.T, doc any, operations []jsoncodec.Operation) (any, error) {
	t.Helper()
	patch, err := jsonpatch.CompileOperations(operations, jsonpatch.WithCapabilities(jsonpatch.AllCapabilities))
	if err != nil {
		return nil, err
	}
	if err := jsonpatch.ApplyInPlace(patch, &doc); err != nil {
		return nil, err
	}
	return doc, nil
}

// TestCase represents a single operation test case.
type TestCase struct {
	Name      string
	Doc       any
	Operation jsoncodec.Operation
	Expected  any
	WantErr   bool
	Comment   string
}

// MultiOperationTestCase represents a multi-operation test case.
type MultiOperationTestCase struct {
	Name       string
	Doc        any
	Operations []jsoncodec.Operation
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
			assert.Equal(t, tc.Expected, result, "ApplyInPlace() result mismatch")
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
			assert.Equal(t, tc.Expected, result, "ApplyInPlace() result mismatch")
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
