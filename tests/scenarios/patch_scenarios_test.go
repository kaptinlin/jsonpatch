package jsonpatch_test

import (
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/google/go-cmp/cmp"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/data"
)

func TestAutomated(t *testing.T) {
	t.Parallel()
	testSuites := []AutomatedTestSuite{
		{
			Name:  "JSON Patch spec",
			Tests: convertSpecTestCases(data.SpecTestCases),
		},
		{
			Name:  "JSON Patch smoke tests",
			Tests: convertTestCases(data.TestCases),
		},
	}

	for _, suite := range testSuites {
		t.Run(suite.Name, func(t *testing.T) {
			t.Parallel()
			for _, test := range suite.Tests {
				if test.Disabled {
					continue
				}

				testName := test.Comment
				if testName == "" {
					if test.Error != "" {
						testName = test.Error
					} else {
						patchJSON, _ := json.Marshal(test.Patch)
						testName = string(patchJSON)
					}
				}

				options := jsonpatch.WithMutate(false)

				switch {
				case test.Expected != nil:
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						result, err := jsonpatch.ApplyPatch(test.Doc, test.Patch, options)
						if err != nil {
							t.Fatalf("ApplyPatch() unexpected error: %v", err)
						}

						if diff := cmp.Diff(test.Expected, result.Doc); diff != "" {
							t.Errorf("ApplyPatch() mismatch (-want +got):\n%s", diff)
						}
					})
				case test.Error != "":
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						// First validate operations
						validationFailed := false
						for _, op := range test.Patch {
							if err := jsonpatch.ValidateOperation(op, false); err != nil {
								validationFailed = true
								break
							}
						}

						// If validation passed, try applying patch
						if !validationFailed {
							_, err := jsonpatch.ApplyPatch(test.Doc, test.Patch, options)
							if err == nil {
								t.Error("ApplyPatch() expected error, got nil")
							}
						}
					})
				default:
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						t.Fatalf("Invalid test case: %+v", test)
					})
				}
			}
		})
	}
}

type AutomatedTestSuite struct {
	Name  string
	Tests []AutomatedTestCase
}

type AutomatedTestCase struct {
	Comment  string
	Doc      interface{}
	Patch    []jsonpatch.Operation
	Expected interface{}
	Error    string
	Disabled bool
}

func convertSpecTestCases(specCases []data.TestCase) []AutomatedTestCase {
	result := make([]AutomatedTestCase, len(specCases))
	for i, tc := range specCases {
		result[i] = AutomatedTestCase{
			Comment:  tc.Comment,
			Doc:      tc.Doc,
			Patch:    convertPatch(tc.Patch),
			Expected: tc.Expected,
			Error:    tc.Error,
			Disabled: tc.Disabled,
		}
	}
	return result
}

func convertTestCases(testCases []data.TestCase) []AutomatedTestCase {
	result := make([]AutomatedTestCase, len(testCases))
	for i, tc := range testCases {
		result[i] = AutomatedTestCase{
			Comment:  tc.Comment,
			Doc:      tc.Doc,
			Patch:    convertPatch(tc.Patch),
			Expected: tc.Expected,
			Error:    tc.Error,
			Disabled: tc.Disabled,
		}
	}
	return result
}

func convertPatch(patch []map[string]any) []jsonpatch.Operation {
	result := make([]jsonpatch.Operation, len(patch))
	for i, op := range patch {
		var operation jsonpatch.Operation
		if v, ok := op["op"].(string); ok {
			operation.Op = v
		}
		if v, ok := op["path"].(string); ok {
			operation.Path = v
		}
		if v, exists := op["value"]; exists {
			operation.Value = v
		}
		if v, ok := op["from"].(string); ok {
			operation.From = v
		}
		if v, ok := op["inc"].(float64); ok {
			operation.Inc = v
		}
		if v, ok := op["pos"].(int); ok {
			operation.Pos = v
		} else if v, ok := op["pos"].(float64); ok {
			operation.Pos = int(v)
		}
		if v, ok := op["str"].(string); ok {
			operation.Str = v
		}
		if v, ok := op["len"].(int); ok {
			operation.Len = v
		} else if v, ok := op["len"].(float64); ok {
			operation.Len = int(v)
		}
		if v, ok := op["not"].(bool); ok {
			operation.Not = v
		}
		if v, ok := op["type"].(string); ok {
			operation.Type = v
		}
		if v, ok := op["ignore_case"].(bool); ok {
			operation.IgnoreCase = v
		}
		if v, ok := op["apply"].([]jsonpatch.Operation); ok {
			operation.Apply = v
		}
		if v, ok := op["props"].(map[string]any); ok {
			operation.Props = v
		}
		if v, ok := op["deleteNull"].(bool); ok {
			operation.DeleteNull = v
		}
		if v, ok := op["oldValue"]; ok {
			operation.OldValue = v
		}
		result[i] = operation
	}
	return result
}
