package jsonpatch_test

import (
	"fmt"
	"testing"

	"github.com/go-json-experiment/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
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

				switch {
				case test.Expected != nil:
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						patchJSON, err := json.Marshal(test.RawPatch)
						require.NoError(t, err)
						patch, err := jsonpatch.CompileJSON(patchJSON)
						if err != nil {
							require.FailNow(t, fmt.Sprintf("CompileJSON() unexpected error: %v", err))
						}
						result, err := jsonpatch.Apply(patch, test.Doc)
						if err != nil {
							require.FailNow(t, fmt.Sprintf("Apply() unexpected error: %v", err))
						}

						assertJSONEqual(t, test.Expected, result.Doc)
					})
				case test.Error != "":
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						if _, err := jsoncodec.Decode(test.RawPatch, jsoncodec.PatchOptions{}); err != nil {
							return
						}
						patchJSON, err := json.Marshal(test.RawPatch)
						require.NoError(t, err)
						patch, err := jsonpatch.CompileJSON(patchJSON)
						if err != nil {
							return
						}
						if _, err := jsonpatch.Apply(patch, test.Doc); err == nil {
							assert.Fail(t, "Apply() expected error, got nil")
						}
					})
				default:
					t.Run(testName, func(t *testing.T) {
						t.Parallel()
						require.FailNow(t, fmt.Sprintf("Invalid test case: %+v", test))
					})
				}
			}
		})
	}
}

func assertJSONEqual(t *testing.T, expected, actual any) {
	t.Helper()
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err)
	actualJSON, err := json.Marshal(actual)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedJSON), string(actualJSON))
}

type AutomatedTestSuite struct {
	Name  string
	Tests []AutomatedTestCase
}

type AutomatedTestCase struct {
	Comment  string
	Doc      any
	RawPatch []map[string]any
	Patch    []jsoncodec.Operation
	Expected any
	Error    string
	Disabled bool
}

func convertSpecTestCases(specCases []data.TestCase) []AutomatedTestCase {
	result := make([]AutomatedTestCase, len(specCases))
	for i, tc := range specCases {
		result[i] = AutomatedTestCase{
			Comment:  tc.Comment,
			Doc:      tc.Doc,
			RawPatch: tc.Patch,
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
			RawPatch: tc.Patch,
			Patch:    convertPatch(tc.Patch),
			Expected: tc.Expected,
			Error:    tc.Error,
			Disabled: tc.Disabled,
		}
	}
	return result
}

func convertPatch(patch []map[string]any) []jsoncodec.Operation {
	result := make([]jsoncodec.Operation, len(patch))
	for i, op := range patch {
		var operation jsoncodec.Operation
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
		if v, ok := op["apply"].([]jsoncodec.Operation); ok {
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
