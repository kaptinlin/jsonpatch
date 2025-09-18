package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

// TestNegation tests predicate negation using the "not" flag
func TestNegation(t *testing.T) {
	testCases := []testutils.TestCase{
		{
			Name:      "not_flag_succeeds_when_values_differ",
			Doc:       map[string]interface{}{"value": 42},
			Operation: jsonpatch.Operation{"op": "test", "path": "/value", "value": 100, "not": true},
			Expected:  map[string]interface{}{"value": 42},
			Comment:   "Test with not=true should succeed when values differ",
		},
		{
			Name:       "not_flag_fails_when_values_match",
			Doc:        map[string]interface{}{"value": 42},
			Operation:  jsonpatch.Operation{"op": "test", "path": "/value", "value": 42, "not": true},
			ShouldFail: true,
			Comment:    "Test with not=true should fail when values match",
		},
	}

	testutils.RunTestCases(t, testCases)
}
