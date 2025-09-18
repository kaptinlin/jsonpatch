package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

// TestNot tests the NOT second-order predicate operation
func TestNot(t *testing.T) {
	testCases := []testutils.MultiOperationTestCase{
		{
			Name: "succeeds_when_all_predicates_fail",
			Doc:  map[string]interface{}{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 1},
						map[string]interface{}{"op": "test", "path": "/bar", "value": 3},
					},
				},
			},
			Expected: map[string]interface{}{"foo": 2, "bar": 2},
			Comment:  "NOT should succeed when all predicates fail",
		},
		{
			Name: "fails_when_one_predicate_passes",
			Doc:  map[string]interface{}{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 1},
						map[string]interface{}{"op": "test", "path": "/bar", "value": 2},
					},
				},
			},
			ShouldFail: true,
			Comment:    "NOT should fail when any predicate passes",
		},
		{
			Name: "fails_when_all_predicates_pass",
			Doc:  map[string]interface{}{"foo": 1, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					"op":   "not",
					"path": "",
					"apply": []interface{}{
						map[string]interface{}{"op": "test", "path": "/foo", "value": 1},
						map[string]interface{}{"op": "test", "path": "/bar", "value": 2},
					},
				},
			},
			ShouldFail: true,
			Comment:    "NOT should fail when all predicates pass",
		},
	}

	testutils.RunMultiOperationTestCases(t, testCases)
}
