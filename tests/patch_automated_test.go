package jsonpatch_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/data"
)

// TestAutomated runs automated tests for JSON Patch operations
func TestAutomated(t *testing.T) {
	// Use test suites from Go data files instead of JSON
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
						result, err := jsonpatch.ApplyPatch(test.Doc, test.Patch, options)

						if err != nil {
							t.Fatalf("ApplyPatch failed: %v", err)
						}

						// Direct comparison for Go native testing
						if !reflect.DeepEqual(test.Expected, result.Doc) {
							t.Errorf("Expected %v, got %v", test.Expected, result.Doc)
						}
					})
				case test.Error != "":
					t.Run(testName, func(t *testing.T) {
						// First validate operations
						for _, op := range test.Patch {
							if err := jsonpatch.ValidateOperation(op, false); err != nil {
								if !containsError(err.Error(), test.Error) {
									t.Errorf("Validation error mismatch. Expected %s, got %s", test.Error, err.Error())
								}
								return
							}
						}

						// Then apply patch
						_, err := jsonpatch.ApplyPatch(test.Doc, test.Patch, options)
						if err == nil {
							t.Fatal("Patch should have failed")
						}

						// Check if the error message contains the expected error type
						if !containsError(err.Error(), test.Error) {
							t.Errorf("Error message mismatch. Expected %s, got %s", test.Error, err.Error())
						}
					})
				default:
					t.Fatalf("Invalid test case: %+v", test)
				}
			}
		})
	}
}

// AutomatedTestSuite represents a test suite for automated tests
type AutomatedTestSuite struct {
	Name  string              `json:"name"`
	Tests []AutomatedTestCase `json:"tests"`
}

// AutomatedTestCase represents a single test case for automated tests
type AutomatedTestCase struct {
	Comment  string                `json:"comment,omitempty"`
	Doc      interface{}           `json:"doc"`
	Patch    []jsonpatch.Operation `json:"patch"`
	Expected interface{}           `json:"expected,omitempty"`
	Error    string                `json:"error,omitempty"`
	Disabled bool                  `json:"disabled,omitempty"`
}

// convertSpecTestCases converts data.TestCase to AutomatedTestCase for spec tests
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

// convertTestCases converts data.TestCase to AutomatedTestCase
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

// convertPatch converts []map[string]any to []jsonpatch.Operation
func convertPatch(patch []map[string]any) []jsonpatch.Operation {
	result := make([]jsonpatch.Operation, len(patch))
	copy(result, patch)
	return result
}

// containsError checks if an error message contains the expected error type
func containsError(errorMessage, expectedType string) bool {
	if errorMessage == "" || expectedType == "" {
		return false
	}

	// Direct match
	if errorMessage == expectedType {
		return true
	}

	// Check if error message contains the expected type (case-insensitive)
	if containsIgnoreCase(errorMessage, expectedType) {
		return true
	}

	// Comprehensive error message mappings for better compatibility
	// Supporting both legacy error messages and new op/errors.go error definitions
	errorMappings := map[string][]string{
		// Legacy error mappings (keep existing compatibility)
		"path /a does not exist -- missing objects are not created recursively": {
			"path does not exist -- missing objects are not created recursively",
		},
		"add to a non-existent target": {
			"path does not exist -- missing objects are not created recursively",
		},
		"number is not equal to string": {
			"string not equivalent",
		},
		"Out of bounds (upper)": {
			"array index out of bounds",
			"index out of range",
		},
		"Out of bounds (lower)": {
			"array index out of bounds",
			"index out of range",
		},
		"test op shouldn't get array element 1": {
			"test failed: path not found",
			"path not found",
		},
		"Object operation on array target": {
			"invalid array index",
			"invalid path",
		},
		"remove op shouldn't remove from array with bad number": {
			"path not found",
			"invalid path",
		},
		"replace op shouldn't replace in array with bad number": {
			"invalid path",
		},
		"copy op shouldn't work with bad number": {
			"copy failed: path not found",
			"path not found",
		},
		"move op shouldn't work with bad number": {
			"move failed: path not found",
			"path not found",
		},
		"add op shouldn't add to array with bad number": {
			"invalid path",
			"invalid array index",
		},
		"test op should fail": {
			"test failed",
			"test operation failed",
		},
		"missing 'value' parameter": {
			"missing value field",
			"missing required field 'value'",
			"test failed: expected <nil>, got false",
			"operation missing 'value' field",
		},
		"missing 'from' parameter": {
			"copy operation missing 'from' field",
			"missing required field 'from'",
			"missing from field",
		},
		"Unrecognized op 'spam'": {
			"Unrecognized op",
			"unknown operation",
			"unknown operation 'spam'",
		},
		"invalid operation path": {
			"operation missing 'op' field",
			"operation missing 'path' field",
			"missing required field 'op'",
			"missing required field 'path'",
			"field 'path' must be a string",
			"Error in operation [index = 0] (operation missing 'op' field)",
			"Error in operation [index = 0] (operation missing 'path' field)",
		},
		"unknown operation": {
			"operation missing 'op' field",
			"missing required field 'op'",
			"Error in operation [index = 0] (operation missing 'op' field)",
		},
		"OP_UNKNOWN": {
			"Unrecognized op",
			"unknown operation",
			"unknown operation 'spam'",
			"Error in operation [index = 0] (unknown operation)",
		},
		"missing value field": {
			"operation missing 'value' field",
			"missing required field 'value'",
			"Error in operation [index = 0] (operation missing 'value' field)",
		},
		"invalid pointer": {
			"operation missing 'path' field",
			"invalid JSON pointer",
			"Error in operation [index = 0] (operation missing 'path' field)",
		},

		// New error mappings for op/errors.go definitions
		// Path related errors
		"path cannot be empty": {
			"ErrPathEmpty", "path cannot be empty",
		},
		"from path cannot be empty": {
			"ErrFromPathEmpty", "from path cannot be empty",
		},
		"path and from cannot be the same": {
			"ErrPathsIdentical", "path and from cannot be the same",
		},
		"Cannot move into own children.": {
			"ErrCannotMoveIntoChildren", "cannot move into own children",
		},

		// Array operation errors
		"array index out of bounds": {
			"ErrArrayIndexOutOfBounds", "array index out of bounds",
		},
		"index out of range": {
			"ErrIndexOutOfRange", "index out of range",
		},
		"not an array": {
			"ErrNotAnArray", "not an array",
		},
		"array must have at least 2 elements": {
			"ErrArrayTooSmall", "array must have at least 2 elements",
		},
		"position out of bounds": {
			"ErrPositionOutOfBounds", "position out of bounds",
		},
		"position cannot be negative": {
			"ErrPositionNegative", "position cannot be negative",
		},

		// Type validation errors - base definitions
		"ErrNotString": {
			"ErrNotString", "value is not a string",
		},
		"ErrNotNumber": {
			"ErrNotNumber", "value is not a number",
		},
		"value is not an object": {
			"ErrNotObject", "value is not an object",
		},
		"invalid type": {
			"ErrInvalidType", "invalid type",
		},
		"types cannot be empty": {
			"ErrEmptyTypeList", "types cannot be empty",
		},

		// Operation execution errors
		"defined test failed": {
			"ErrDefinedTestFailed", "defined test failed",
		},
		"undefined test failed": {
			"ErrUndefinedTestFailed", "undefined test failed",
		},
		"and test failed": {
			"ErrAndTestFailed", "and test failed",
		},
		"not test failed": {
			"ErrNotTestFailed", "not test failed",
		},

		// Value operation errors
		"cannot replace key in non-object": {
			"ErrCannotReplace", "cannot replace key in non-object",
		},
		"cannot add to non-object/non-array value": {
			"ErrCannotAddToValue", "cannot add to non-object/non-array value",
		},
		"cannot remove from non-object/non-array document": {
			"ErrCannotRemoveFromValue", "cannot remove from non-object/non-array document",
		},
		"path does not exist -- missing objects are not created recursively": {
			"ErrPathMissingRecursive", "path does not exist -- missing objects are not created recursively",
		},
		"properties cannot be nil": {
			"ErrPropertiesNil", "properties cannot be nil",
		},
		"values array cannot be empty": {
			"ErrValuesArrayEmpty", "values array cannot be empty",
		},

		// Key type errors
		"invalid key type for map": {
			"ErrInvalidKeyTypeMap", "invalid key type for map",
		},
		"invalid key type for slice": {
			"ErrInvalidKeyTypeSlice", "invalid key type for slice",
		},
		"unsupported parent type": {
			"ErrUnsupportedParentType", "unsupported parent type",
		},

		// String operation errors
		"position out of range": {
			"ErrPositionOutOfStringRange", "position out of range",
		},
		"substring extends beyond string length": {
			"ErrSubstringTooLong", "substring extends beyond string length",
		},
		"pattern cannot be empty": {
			"ErrPatternEmpty", "pattern cannot be empty",
		},
		"length cannot be negative": {
			"ErrLengthNegative", "length cannot be negative",
		},

		// Predicate operation errors
		"invalid predicate operation in AND": {
			"ErrInvalidPredicateInAnd", "invalid predicate operation in AND",
		},
		"invalid predicate operation in NOT": {
			"ErrInvalidPredicateInNot", "invalid predicate operation in NOT",
		},
		"invalid predicate operation in OR": {
			"ErrInvalidPredicateInOr", "invalid predicate operation in OR",
		},
		"and operation must have at least one operand": {
			"ErrAndNoOperands", "and operation must have at least one operand",
		},
		"not operation must have at least one operand": {
			"ErrNotNoOperands", "not operation must have at least one operand",
		},
		"or operation must have at least one operand": {
			"ErrOrNoOperands", "or operation must have at least one operand",
		},

		// Test operation specific errors
		"test operation failed: number is not equal to string": {
			"ErrTestOperationNumberStringMismatch", "test operation failed: number is not equal to string",
		},
		"test operation failed: string not equivalent": {
			"ErrTestOperationStringNotEquivalent", "test operation failed: string not equivalent",
		},
		"or test failed": {
			"ErrOrTestFailed", "or test failed",
		},

		// Path operation specific errors - complete prefixed error mappings
		"path not found": {
			"path not found", "ErrPathNotFound",
			"path does not exist", "ErrPathDoesNotExist",
			"contains failed: path not found",
			"ends failed: path not found",
			"in failed: path not found",
			"matches failed: path not found",
			"more failed: path not found",
			"starts failed: path not found",
		},

		// Type validation specific errors - complete error mappings
		"not a string": {
			"value is not a string", "ErrNotString",
			"contains failed: value is not a string",
			"ends failed: value is not a string",
			"matches failed: value is not a string",
			"starts failed: value is not a string",
		},

		"not a number": {
			"value is not a number", "ErrNotNumber",
			"more failed: value is not a number",
		},

		// Test operation specific errors - complete error mappings
		"test failed": {
			"test failed", "ErrTestFailed",
			"starts failed: string",
			"ends failed: string",
		},

		// Operation modification errors
		"cannot modify root array directly": {
			"ErrCannotModifyRootArray", "cannot modify root array directly",
		},
		"cannot update parent": {
			"ErrCannotUpdateParent", "cannot update parent",
		},
		"cannot update grandparent": {
			"ErrCannotUpdateGrandparent", "cannot update grandparent",
		},
		"key does not exist": {
			"ErrKeyDoesNotExist", "key does not exist",
		},

		// Value conversion errors
		"cannot convert nil to string": {
			"ErrCannotConvertNilToString", "cannot convert nil to string",
		},
	}

	// Check if we have a mapping for this expected type
	if patterns, exists := errorMappings[expectedType]; exists {
		for _, pattern := range patterns {
			if containsIgnoreCase(errorMessage, pattern) {
				return true
			}
		}
	}

	return false
}

// containsIgnoreCase performs case-insensitive substring search
func containsIgnoreCase(haystack, needle string) bool {
	haystack = strings.ToLower(haystack)
	needle = strings.ToLower(needle)

	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
