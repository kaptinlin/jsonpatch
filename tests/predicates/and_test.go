package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

// TestAnd tests the AND second-order predicate operation
func TestAnd(t *testing.T) {
	testCases := []testutils.MultiOperationTestCase{
		{
			Name: "succeeds_when_both_predicates_pass",
			Doc:  map[string]interface{}{"foo": 1, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "and",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			Expected: map[string]interface{}{"foo": 1, "bar": 2},
			Comment:  "AND should succeed when all predicates pass",
		},
		{
			Name: "fails_when_one_predicate_fails",
			Doc:  map[string]interface{}{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "and",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			WantErr: true,
			Comment:    "AND should fail when any predicate fails",
		},
	}

	testutils.RunMultiOperationTestCases(t, testCases)
}
