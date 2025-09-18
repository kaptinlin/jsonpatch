// Package testutils provides shared utilities for JSON Patch testing.
package testutils

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/stretchr/testify/require"
)

// ApplyOperation applies a single operation to a document
func ApplyOperation(t *testing.T, doc interface{}, operation jsonpatch.Operation) interface{} {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	options := jsonpatch.WithMutate(true)
	result, err := jsonpatch.ApplyPatch(doc, patch, options)
	require.NoError(t, err)
	return result.Doc
}

// ApplyOperationWithError applies an operation expecting it to fail
func ApplyOperationWithError(t *testing.T, doc interface{}, operation jsonpatch.Operation) error {
	t.Helper()
	patch := []jsonpatch.Operation{operation}
	options := jsonpatch.WithMutate(true)
	_, err := jsonpatch.ApplyPatch(doc, patch, options)
	require.Error(t, err)
	return err
}

// ApplyOperations applies multiple operations to a document
func ApplyOperations(t *testing.T, doc interface{}, operations []jsonpatch.Operation) interface{} {
	t.Helper()
	options := jsonpatch.WithMutate(true)
	result, err := jsonpatch.ApplyPatch(doc, operations, options)
	require.NoError(t, err)
	return result.Doc
}

// ApplyOperationsWithError applies multiple operations expecting them to fail
func ApplyOperationsWithError(t *testing.T, doc interface{}, operations []jsonpatch.Operation) error {
	t.Helper()
	options := jsonpatch.WithMutate(true)
	_, err := jsonpatch.ApplyPatch(doc, operations, options)
	require.Error(t, err)
	return err
}

// TestCase represents a single operation test case
type TestCase struct {
	Name       string              // Test case name
	Doc        interface{}         // Input document
	Operation  jsonpatch.Operation // Operation to apply
	Expected   interface{}         // Expected result
	ShouldFail bool                // Whether operation should fail
	Comment    string              // Comment or source information
}

// MultiOperationTestCase represents a multi-operation test case
type MultiOperationTestCase struct {
	Name       string                // Test case name
	Doc        interface{}           // Input document
	Operations []jsonpatch.Operation // Operations to apply
	Expected   interface{}           // Expected result
	ShouldFail bool                  // Whether operations should fail
	Comment    string                // Comment or source information
}

// RunTestCase executes a single test case
func RunTestCase(t *testing.T, tc TestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		if tc.ShouldFail {
			_ = ApplyOperationWithError(t, tc.Doc, tc.Operation)
		} else {
			result := ApplyOperation(t, tc.Doc, tc.Operation)
			require.Equal(t, tc.Expected, result)
		}
	})
}

// RunMultiOperationTestCase executes a multi-operation test case
func RunMultiOperationTestCase(t *testing.T, tc MultiOperationTestCase) {
	t.Helper()
	t.Run(tc.Name, func(t *testing.T) {
		if tc.ShouldFail {
			_ = ApplyOperationsWithError(t, tc.Doc, tc.Operations)
		} else {
			result := ApplyOperations(t, tc.Doc, tc.Operations)
			require.Equal(t, tc.Expected, result)
		}
	})
}

// RunTestCases executes multiple test cases
func RunTestCases(t *testing.T, testCases []TestCase) {
	t.Helper()
	for _, tc := range testCases {
		RunTestCase(t, tc)
	}
}

// RunMultiOperationTestCases executes multiple multi-operation test cases
func RunMultiOperationTestCases(t *testing.T, testCases []MultiOperationTestCase) {
	t.Helper()
	for _, tc := range testCases {
		RunMultiOperationTestCase(t, tc)
	}
}
