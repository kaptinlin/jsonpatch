package compatibility

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

// TypeScriptTestCase represents a test case from TypeScript implementation
type TypeScriptTestCase struct {
	Name       string                `json:"name"`
	Doc        interface{}           `json:"doc"`
	Patch      []jsonpatch.Operation `json:"patch"`
	Expected   interface{}           `json:"expected,omitempty"`
	WantErr    bool                  `json:"shouldFail,omitempty"`
	Comment    string                `json:"comment,omitempty"`
	Source     string                `json:"source,omitempty"` // Which TypeScript file this came from
}

// TestTypeScriptParity verifies that our Go implementation behaves consistently with TypeScript
func TestTypeScriptParity(t *testing.T) {
	// Load test cases from known working operations
	testCases := getKnownWorkingTestCases()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.WantErr {
				_ = testutils.ApplyOperationsWithError(t, tc.Doc, tc.Patch)
			} else {
				result := testutils.ApplyOperations(t, tc.Doc, tc.Patch)
				if diff := cmp.Diff(tc.Expected, result); diff != "" {
					t.Errorf("ApplyOperations() mismatch (-want +got):\n%s\n%s", diff, tc.Comment)
				}
			}
		})
	}
}

func TestBasicOperationParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "add_to_object",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "add", Path: "/baz", Value: "qux"},
			},
			Expected: map[string]interface{}{"foo": "bar", "baz": "qux"},
			Comment:  "Should add new property to object",
			Source:   "typescript:basic.spec.ts",
		},
		{
			Name: "replace_value",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "replace", Path: "/foo", Value: "baz"},
			},
			Expected: map[string]interface{}{"foo": "baz"},
			Comment:  "Should replace existing value",
			Source:   "typescript:basic.spec.ts",
		},
		{
			Name: "remove_property",
			Doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			Patch: []jsonpatch.Operation{
				{Op: "remove", Path: "/baz"},
			},
			Expected: map[string]interface{}{"foo": "bar"},
			Comment:  "Should remove property from object",
			Source:   "typescript:basic.spec.ts",
		},
		{
			Name: "test_successful",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: "bar"},
			},
			Expected: map[string]interface{}{"foo": "bar"},
			Comment:  "Test should pass and leave document unchanged",
			Source:   "typescript:test.spec.ts",
		},
		{
			Name: "test_failure",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "test", Path: "/foo", Value: "baz"},
			},
			WantErr: true,
			Comment: "Test should fail when values don't match",
			Source:     "typescript:test.spec.ts",
		},
		{
			Name: "copy_operation",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "copy", From: "/foo", Path: "/baz"},
			},
			Expected: map[string]interface{}{"foo": "bar", "baz": "bar"},
			Comment:  "Should copy value to new location",
			Source:   "typescript:copy.spec.ts",
		},
		{
			Name: "move_operation",
			Doc:  map[string]interface{}{"foo": "bar", "baz": "qux"},
			Patch: []jsonpatch.Operation{
				{Op: "move", From: "/baz", Path: "/moved"},
			},
			Expected: map[string]interface{}{"foo": "bar", "moved": "qux"},
			Comment:  "Should move value to new location",
			Source:   "typescript:move.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestArrayOperationParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "add_to_array_end",
			Doc:  []interface{}{"foo", "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "add", Path: "/-", Value: "baz"},
			},
			Expected: []interface{}{"foo", "bar", "baz"},
			Comment:  "Should add to end of array",
			Source:   "typescript:array.spec.ts",
		},
		{
			Name: "add_to_array_middle",
			Doc:  []interface{}{"foo", "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "add", Path: "/1", Value: "baz"},
			},
			Expected: []interface{}{"foo", "baz", "bar"},
			Comment:  "Should insert into array at index",
			Source:   "typescript:array.spec.ts",
		},
		{
			Name: "remove_from_array",
			Doc:  []interface{}{"foo", "bar", "baz"},
			Patch: []jsonpatch.Operation{
				{Op: "remove", Path: "/1"},
			},
			Expected: []interface{}{"foo", "baz"},
			Comment:  "Should remove from array at index",
			Source:   "typescript:array.spec.ts",
		},
		{
			Name: "replace_array_element",
			Doc:  []interface{}{"foo", "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "replace", Path: "/0", Value: "baz"},
			},
			Expected: []interface{}{"baz", "bar"},
			Comment:  "Should replace array element",
			Source:   "typescript:array.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestPredicateOperationParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "contains_successful",
			Doc:  map[string]interface{}{"text": "hello world"},
			Patch: []jsonpatch.Operation{
				{Op: "contains", Path: "/text", Value: "world"},
			},
			Expected: map[string]interface{}{"text": "hello world"},
			Comment:  "Contains should pass when substring exists",
			Source:   "typescript:predicates.spec.ts",
		},
		{
			Name: "contains_failure",
			Doc:  map[string]interface{}{"text": "hello world"},
			Patch: []jsonpatch.Operation{
				{Op: "contains", Path: "/text", Value: "xyz"},
			},
			WantErr: true,
			Comment:    "Contains should fail when substring doesn't exist",
			Source:     "typescript:predicates.spec.ts",
		},
		{
			Name: "type_check_number",
			Doc:  map[string]interface{}{"value": 42},
			Patch: []jsonpatch.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
			Expected: map[string]interface{}{"value": 42},
			Comment:  "Type check should pass for correct type",
			Source:   "typescript:type.spec.ts",
		},
		{
			Name: "type_check_failure",
			Doc:  map[string]interface{}{"value": "string"},
			Patch: []jsonpatch.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
			WantErr: true,
			Comment:    "Type check should fail for incorrect type",
			Source:     "typescript:type.spec.ts",
		},
		{
			Name: "less_than_check",
			Doc:  map[string]interface{}{"value": 5},
			Patch: []jsonpatch.Operation{
				{Op: "less", Path: "/value", Value: 10},
			},
			Expected: map[string]interface{}{"value": 5},
			Comment:  "Less than check should pass",
			Source:   "typescript:comparison.spec.ts",
		},
		{
			Name: "more_than_check",
			Doc:  map[string]interface{}{"value": 15},
			Patch: []jsonpatch.Operation{
				{Op: "more", Path: "/value", Value: 10},
			},
			Expected: map[string]interface{}{"value": 15},
			Comment:  "More than check should pass",
			Source:   "typescript:comparison.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestExtendedOperationParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "inc_operation",
			Doc:  map[string]interface{}{"counter": 5},
			Patch: []jsonpatch.Operation{
				{Op: "inc", Path: "/counter", Inc: 3},
			},
			Expected: map[string]interface{}{"counter": float64(8)}, // JSON unmarshaling converts to float64
			Comment:  "Increment should add to numeric value",
			Source:   "typescript:inc.spec.ts",
		},
		{
			Name: "flip_operation",
			Doc:  map[string]interface{}{"enabled": true},
			Patch: []jsonpatch.Operation{
				{Op: "flip", Path: "/enabled"},
			},
			Expected: map[string]interface{}{"enabled": false},
			Comment:  "Flip should toggle boolean value",
			Source:   "typescript:flip.spec.ts",
		},
		{
			Name: "str_ins_operation",
			Doc:  map[string]interface{}{"text": "hello world"},
			Patch: []jsonpatch.Operation{
				{Op: "str_ins", Path: "/text", Pos: 5, Str: " beautiful"},
			},
			Expected: map[string]interface{}{"text": "hello beautiful world"},
			Comment:  "String insert should insert text at position",
			Source:   "typescript:strins.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestErrorHandlingParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "path_not_found",
			Doc:  map[string]interface{}{"foo": "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "remove", Path: "/nonexistent"},
			},
			WantErr: true,
			Comment:    "Should fail when path doesn't exist",
			Source:     "typescript:errors.spec.ts",
		},
		{
			Name: "invalid_array_index",
			Doc:  []interface{}{"foo", "bar"},
			Patch: []jsonpatch.Operation{
				{Op: "remove", Path: "/5"},
			},
			WantErr: true,
			Comment:    "Should fail for out-of-bounds array index",
			Source:     "typescript:errors.spec.ts",
		},
		{
			Name: "type_mismatch_operation",
			Doc:  map[string]interface{}{"value": "string"},
			Patch: []jsonpatch.Operation{
				{Op: "inc", Path: "/value", Inc: 1},
			},
			WantErr: true,
			Comment:    "Should fail when operation type doesn't match value type",
			Source:     "typescript:errors.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestSecondOrderPredicateParity(t *testing.T) {
	testCases := []TypeScriptTestCase{
		{
			Name: "not_predicate_success",
			Doc:  map[string]interface{}{"foo": 1, "bar": 2},
			Patch: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 2},
					},
				},
			},
			Expected: map[string]interface{}{"foo": 1, "bar": 2},
			Comment:  "NOT should succeed when inner predicate fails",
			Source:   "typescript:second-order-predicates.spec.ts",
		},
		{
			Name: "not_predicate_failure",
			Doc:  map[string]interface{}{"foo": 1, "bar": 2},
			Patch: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
					},
				},
			},
			WantErr: true,
			Comment:    "NOT should fail when inner predicate succeeds",
			Source:     "typescript:second-order-predicates.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

// getKnownWorkingTestCases returns a curated set of test cases that are known to work
func getKnownWorkingTestCases() []TypeScriptTestCase {
	var allTestCases []TypeScriptTestCase

	// Combine all the working test cases
	allTestCases = append(allTestCases, getBasicOperationTestCases()...)
	allTestCases = append(allTestCases, getArrayOperationTestCases()...)
	allTestCases = append(allTestCases, getPredicateOperationTestCases()...)
	allTestCases = append(allTestCases, getExtendedOperationTestCases()...)
	allTestCases = append(allTestCases, getErrorHandlingTestCases()...)
	allTestCases = append(allTestCases, getSecondOrderPredicateTestCases()...)

	return allTestCases
}

// Helper functions to get specific test case categories
func getBasicOperationTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name:     "basic_add",
			Doc:      map[string]interface{}{"a": 1},
			Patch:    []jsonpatch.Operation{{Op: "add", Path: "/b", Value: 2}},
			Expected: map[string]interface{}{"a": 1, "b": 2},
			Comment:  "Basic add operation",
			Source:   "typescript:add.spec.ts",
		},
		{
			Name:     "basic_remove",
			Doc:      map[string]interface{}{"a": 1, "b": 2},
			Patch:    []jsonpatch.Operation{{Op: "remove", Path: "/b"}},
			Expected: map[string]interface{}{"a": 1},
			Comment:  "Basic remove operation",
			Source:   "typescript:remove.spec.ts",
		},
	}
}

func getArrayOperationTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name:     "array_append",
			Doc:      []interface{}{1, 2},
			Patch:    []jsonpatch.Operation{{Op: "add", Path: "/-", Value: 3}},
			Expected: []interface{}{1, 2, 3},
			Comment:  "Array append operation",
			Source:   "typescript:array.spec.ts",
		},
	}
}

func getPredicateOperationTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name:     "predicate_contains",
			Doc:      map[string]interface{}{"text": "hello"},
			Patch:    []jsonpatch.Operation{{Op: "contains", Path: "/text", Value: "ell"}},
			Expected: map[string]interface{}{"text": "hello"},
			Comment:  "Predicate contains operation",
			Source:   "typescript:contains.spec.ts",
		},
	}
}

func getExtendedOperationTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name:     "extended_inc",
			Doc:      map[string]interface{}{"num": 5},
			Patch:    []jsonpatch.Operation{{Op: "inc", Path: "/num", Inc: 2}},
			Expected: map[string]interface{}{"num": float64(7)},
			Comment:  "Extended inc operation",
			Source:   "typescript:inc.spec.ts",
		},
	}
}

func getErrorHandlingTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name:       "error_path_not_found",
			Doc:        map[string]interface{}{"a": 1},
			Patch:      []jsonpatch.Operation{{Op: "test", Path: "/b", Value: 1}},
			WantErr: true,
			Comment:    "Path not found error",
			Source:     "typescript:errors.spec.ts",
		},
	}
}

func getSecondOrderPredicateTestCases() []TypeScriptTestCase {
	return []TypeScriptTestCase{
		{
			Name: "second_order_not",
			Doc:  map[string]interface{}{"val": 1},
			Patch: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/val", Value: 2},
					},
				},
			},
			Expected: map[string]interface{}{"val": 1},
			Comment:  "Second-order NOT predicate",
			Source:   "typescript:second-order-predicates.spec.ts",
		},
	}
}

// convertToMultiOpTestCases converts TypeScript test cases to multi-operation test cases
func convertToMultiOpTestCases(tsCases []TypeScriptTestCase) []testutils.MultiOperationTestCase {
	multiOpCases := make([]testutils.MultiOperationTestCase, 0, len(tsCases))
	for _, tc := range tsCases {
		multiOpCases = append(multiOpCases, testutils.MultiOperationTestCase{
			Name:       tc.Name,
			Doc:        tc.Doc,
			Operations: tc.Patch,
			Expected:   tc.Expected,
			WantErr: tc.WantErr,
			Comment:    fmt.Sprintf("%s (Source: %s)", tc.Comment, tc.Source),
		})
	}
	return multiOpCases
}

func BenchmarkTypeScriptParity(b *testing.B) {
	testCases := getKnownWorkingTestCases()

	b.ResetTimer()
	for b.Loop() {
		for _, tc := range testCases {
			if !tc.WantErr {
				_, err := jsonpatch.ApplyPatch(tc.Doc, tc.Patch, jsonpatch.WithMutate(true))
				if err != nil {
					b.Errorf("Operation %s failed: %v", tc.Name, err)
				}
			}
		}
	}
}
