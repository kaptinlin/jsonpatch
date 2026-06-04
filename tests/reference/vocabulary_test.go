package reference

import (
	"fmt"
	"slices"
	"testing"

	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

type ReferenceCase struct {
	Name     string
	Doc      any
	Patch    []jsoncodec.Operation
	Expected any
	WantErr  bool
	Comment  string
	Evidence string
}

func TestReferenceVocabulary(t *testing.T) {
	t.Parallel()
	testCases := getKnownWorkingTestCases()

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			if tc.WantErr {
				_ = testutils.ApplyOperationsWithError(t, tc.Doc, tc.Patch)
			} else {
				result := testutils.ApplyOperations(t, tc.Doc, tc.Patch)
				assert.Empty(t, cmp.Diff(tc.Expected, result), "ApplyOperations() mismatch: %s", tc.Comment)
			}
		})
	}
}

func TestBasicOperationReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "add_to_object",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "add", Path: "/baz", Value: "qux"},
			},
			Expected: map[string]any{"foo": "bar", "baz": "qux"},
			Comment:  "Should add new property to object",
			Evidence: "reference:basic.spec.ts",
		},
		{
			Name: "replace_value",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "replace", Path: "/foo", Value: "baz"},
			},
			Expected: map[string]any{"foo": "baz"},
			Comment:  "Should replace existing value",
			Evidence: "reference:basic.spec.ts",
		},
		{
			Name: "remove_property",
			Doc:  map[string]any{"foo": "bar", "baz": "qux"},
			Patch: []jsoncodec.Operation{
				{Op: "remove", Path: "/baz"},
			},
			Expected: map[string]any{"foo": "bar"},
			Comment:  "Should remove property from object",
			Evidence: "reference:basic.spec.ts",
		},
		{
			Name: "test_successful",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "test", Path: "/foo", Value: "bar"},
			},
			Expected: map[string]any{"foo": "bar"},
			Comment:  "Test should pass and leave document unchanged",
			Evidence: "reference:test.spec.ts",
		},
		{
			Name: "test_failure",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "test", Path: "/foo", Value: "baz"},
			},
			WantErr:  true,
			Comment:  "Test should fail when values don't match",
			Evidence: "reference:test.spec.ts",
		},
		{
			Name: "copy_operation",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "copy", From: "/foo", Path: "/baz"},
			},
			Expected: map[string]any{"foo": "bar", "baz": "bar"},
			Comment:  "Should copy value to new location",
			Evidence: "reference:copy.spec.ts",
		},
		{
			Name: "move_operation",
			Doc:  map[string]any{"foo": "bar", "baz": "qux"},
			Patch: []jsoncodec.Operation{
				{Op: "move", From: "/baz", Path: "/moved"},
			},
			Expected: map[string]any{"foo": "bar", "moved": "qux"},
			Comment:  "Should move value to new location",
			Evidence: "reference:move.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestArrayOperationReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "add_to_array_end",
			Doc:  []any{"foo", "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "add", Path: "/-", Value: "baz"},
			},
			Expected: []any{"foo", "bar", "baz"},
			Comment:  "Should add to end of array",
			Evidence: "reference:array.spec.ts",
		},
		{
			Name: "add_to_array_middle",
			Doc:  []any{"foo", "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "add", Path: "/1", Value: "baz"},
			},
			Expected: []any{"foo", "baz", "bar"},
			Comment:  "Should insert into array at index",
			Evidence: "reference:array.spec.ts",
		},
		{
			Name: "remove_from_array",
			Doc:  []any{"foo", "bar", "baz"},
			Patch: []jsoncodec.Operation{
				{Op: "remove", Path: "/1"},
			},
			Expected: []any{"foo", "baz"},
			Comment:  "Should remove from array at index",
			Evidence: "reference:array.spec.ts",
		},
		{
			Name: "replace_array_element",
			Doc:  []any{"foo", "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "replace", Path: "/0", Value: "baz"},
			},
			Expected: []any{"baz", "bar"},
			Comment:  "Should replace array element",
			Evidence: "reference:array.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestPredicateOperationReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "contains_successful",
			Doc:  map[string]any{"text": "hello world"},
			Patch: []jsoncodec.Operation{
				{Op: "contains", Path: "/text", Value: "world"},
			},
			Expected: map[string]any{"text": "hello world"},
			Comment:  "Contains should pass when substring exists",
			Evidence: "reference:predicates.spec.ts",
		},
		{
			Name: "contains_failure",
			Doc:  map[string]any{"text": "hello world"},
			Patch: []jsoncodec.Operation{
				{Op: "contains", Path: "/text", Value: "xyz"},
			},
			WantErr:  true,
			Comment:  "Contains should fail when substring doesn't exist",
			Evidence: "reference:predicates.spec.ts",
		},
		{
			Name: "type_check_number",
			Doc:  map[string]any{"value": 42},
			Patch: []jsoncodec.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
			Expected: map[string]any{"value": 42},
			Comment:  "Type check should pass for correct type",
			Evidence: "reference:type.spec.ts",
		},
		{
			Name: "type_check_failure",
			Doc:  map[string]any{"value": "string"},
			Patch: []jsoncodec.Operation{
				{Op: "type", Path: "/value", Value: "number"},
			},
			WantErr:  true,
			Comment:  "Type check should fail for incorrect type",
			Evidence: "reference:type.spec.ts",
		},
		{
			Name: "less_than_check",
			Doc:  map[string]any{"value": 5},
			Patch: []jsoncodec.Operation{
				{Op: "less", Path: "/value", Value: 10},
			},
			Expected: map[string]any{"value": 5},
			Comment:  "Less than check should pass",
			Evidence: "reference:comparison.spec.ts",
		},
		{
			Name: "more_than_check",
			Doc:  map[string]any{"value": 15},
			Patch: []jsoncodec.Operation{
				{Op: "more", Path: "/value", Value: 10},
			},
			Expected: map[string]any{"value": 15},
			Comment:  "More than check should pass",
			Evidence: "reference:comparison.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestExtendedOperationReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "inc_operation",
			Doc:  map[string]any{"counter": 5},
			Patch: []jsoncodec.Operation{
				{Op: "inc", Path: "/counter", Inc: 3},
			},
			Expected: map[string]any{"counter": float64(8)}, // JSON unmarshaling converts to float64
			Comment:  "Increment should add to numeric value",
			Evidence: "reference:inc.spec.ts",
		},
		{
			Name: "flip_operation",
			Doc:  map[string]any{"enabled": true},
			Patch: []jsoncodec.Operation{
				{Op: "flip", Path: "/enabled"},
			},
			Expected: map[string]any{"enabled": false},
			Comment:  "Flip should toggle boolean value",
			Evidence: "reference:flip.spec.ts",
		},
		{
			Name: "str_ins_operation",
			Doc:  map[string]any{"text": "hello world"},
			Patch: []jsoncodec.Operation{
				{Op: "str_ins", Path: "/text", Pos: 5, Str: " beautiful"},
			},
			Expected: map[string]any{"text": "hello beautiful world"},
			Comment:  "String insert should insert text at position",
			Evidence: "reference:strins.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestErrorHandlingReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "path_not_found",
			Doc:  map[string]any{"foo": "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "remove", Path: "/nonexistent"},
			},
			WantErr:  true,
			Comment:  "Should fail when path doesn't exist",
			Evidence: "reference:errors.spec.ts",
		},
		{
			Name: "invalid_array_index",
			Doc:  []any{"foo", "bar"},
			Patch: []jsoncodec.Operation{
				{Op: "remove", Path: "/5"},
			},
			WantErr:  true,
			Comment:  "Should fail for out-of-bounds array index",
			Evidence: "reference:errors.spec.ts",
		},
		{
			Name: "type_mismatch_operation",
			Doc:  map[string]any{"value": "string"},
			Patch: []jsoncodec.Operation{
				{Op: "inc", Path: "/value", Inc: 1},
			},
			WantErr:  true,
			Comment:  "Should fail when operation type doesn't match value type",
			Evidence: "reference:errors.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func TestSecondOrderPredicateReference(t *testing.T) {
	t.Parallel()
	testCases := []ReferenceCase{
		{
			Name: "not_predicate_success",
			Doc:  map[string]any{"foo": 1, "bar": 2},
			Patch: []jsoncodec.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsoncodec.Operation{
						{Op: "test", Path: "/foo", Value: 2},
					},
				},
			},
			Expected: map[string]any{"foo": 1, "bar": 2},
			Comment:  "NOT should succeed when inner predicate fails",
			Evidence: "reference:second-order-predicates.spec.ts",
		},
		{
			Name: "not_predicate_failure",
			Doc:  map[string]any{"foo": 1, "bar": 2},
			Patch: []jsoncodec.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsoncodec.Operation{
						{Op: "test", Path: "/foo", Value: 1},
					},
				},
			},
			WantErr:  true,
			Comment:  "NOT should fail when inner predicate succeeds",
			Evidence: "reference:second-order-predicates.spec.ts",
		},
	}

	testutils.RunMultiOperationTestCases(t, convertToMultiOpTestCases(testCases))
}

func getKnownWorkingTestCases() []ReferenceCase {
	return slices.Concat(
		getBasicOperationTestCases(),
		getArrayOperationTestCases(),
		getPredicateOperationTestCases(),
		getExtendedOperationTestCases(),
		getErrorHandlingTestCases(),
		getSecondOrderPredicateTestCases(),
	)
}

func getBasicOperationTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name:     "basic_add",
			Doc:      map[string]any{"a": 1},
			Patch:    []jsoncodec.Operation{{Op: "add", Path: "/b", Value: 2}},
			Expected: map[string]any{"a": 1, "b": 2},
			Comment:  "Basic add operation",
			Evidence: "reference:add.spec.ts",
		},
		{
			Name:     "basic_remove",
			Doc:      map[string]any{"a": 1, "b": 2},
			Patch:    []jsoncodec.Operation{{Op: "remove", Path: "/b"}},
			Expected: map[string]any{"a": 1},
			Comment:  "Basic remove operation",
			Evidence: "reference:remove.spec.ts",
		},
	}
}

func getArrayOperationTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name:     "array_append",
			Doc:      []any{1, 2},
			Patch:    []jsoncodec.Operation{{Op: "add", Path: "/-", Value: 3}},
			Expected: []any{1, 2, 3},
			Comment:  "Array append operation",
			Evidence: "reference:array.spec.ts",
		},
	}
}

func getPredicateOperationTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name:     "predicate_contains",
			Doc:      map[string]any{"text": "hello"},
			Patch:    []jsoncodec.Operation{{Op: "contains", Path: "/text", Value: "ell"}},
			Expected: map[string]any{"text": "hello"},
			Comment:  "Predicate contains operation",
			Evidence: "reference:contains.spec.ts",
		},
	}
}

func getExtendedOperationTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name:     "extended_inc",
			Doc:      map[string]any{"num": 5},
			Patch:    []jsoncodec.Operation{{Op: "inc", Path: "/num", Inc: 2}},
			Expected: map[string]any{"num": float64(7)},
			Comment:  "Extended inc operation",
			Evidence: "reference:inc.spec.ts",
		},
	}
}

func getErrorHandlingTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name:     "error_path_not_found",
			Doc:      map[string]any{"a": 1},
			Patch:    []jsoncodec.Operation{{Op: "test", Path: "/b", Value: 1}},
			WantErr:  true,
			Comment:  "Path not found error",
			Evidence: "reference:errors.spec.ts",
		},
	}
}

func getSecondOrderPredicateTestCases() []ReferenceCase {
	return []ReferenceCase{
		{
			Name: "second_order_not",
			Doc:  map[string]any{"val": 1},
			Patch: []jsoncodec.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsoncodec.Operation{
						{Op: "test", Path: "/val", Value: 2},
					},
				},
			},
			Expected: map[string]any{"val": 1},
			Comment:  "Second-order NOT predicate",
			Evidence: "reference:second-order-predicates.spec.ts",
		},
	}
}

func convertToMultiOpTestCases(referenceCases []ReferenceCase) []testutils.MultiOperationTestCase {
	multiOpCases := make([]testutils.MultiOperationTestCase, 0, len(referenceCases))
	for _, tc := range referenceCases {
		multiOpCases = append(multiOpCases, testutils.MultiOperationTestCase{
			Name:       tc.Name,
			Doc:        tc.Doc,
			Operations: tc.Patch,
			Expected:   tc.Expected,
			WantErr:    tc.WantErr,
			Comment:    fmt.Sprintf("%s (Evidence: %s)", tc.Comment, tc.Evidence),
		})
	}
	return multiOpCases
}

func BenchmarkReferenceVocabulary(b *testing.B) {
	testCases := getKnownWorkingTestCases()

	b.ResetTimer()
	for b.Loop() {
		for _, tc := range testCases {
			if !tc.WantErr {
				_, err := applyPatch(b, tc.Doc, tc.Patch)
				if err != nil {
					b.Errorf("Operation %s failed: %v", tc.Name, err)
				}
			}
		}
	}
}
