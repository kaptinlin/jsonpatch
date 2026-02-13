package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestNot(t *testing.T) {
	t.Parallel()
	testCases := []testutils.MultiOperationTestCase{
		{
			Name: "succeeds_when_all_predicates_fail",
			Doc:  map[string]any{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 3},
					},
				},
			},
			Expected: map[string]any{"foo": 2, "bar": 2},
			Comment:  "NOT should succeed when all predicates fail",
		},
		{
			Name: "fails_when_one_predicate_passes",
			Doc:  map[string]any{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			WantErr: true,
			Comment: "NOT should fail when any predicate passes",
		},
		{
			Name: "fails_when_all_predicates_pass",
			Doc:  map[string]any{"foo": 1, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "not",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			WantErr: true,
			Comment: "NOT should fail when all predicates pass",
		},
	}

	testutils.RunMultiOperationTestCases(t, testCases)
}
