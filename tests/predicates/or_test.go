package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestOr(t *testing.T) {
	t.Parallel()
	testCases := []testutils.MultiOperationTestCase{
		{
			Name: "succeeds_when_one_predicate_passes",
			Doc:  map[string]any{"foo": 2, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "or",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			Expected: map[string]any{"foo": 2, "bar": 2},
			Comment:  "OR should succeed when at least one predicate passes",
		},
		{
			Name: "succeeds_when_both_predicates_pass",
			Doc:  map[string]any{"foo": 1, "bar": 2},
			Operations: []jsonpatch.Operation{
				{
					Op:   "or",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			Expected: map[string]any{"foo": 1, "bar": 2},
			Comment:  "OR should succeed when all predicates pass",
		},
		{
			Name: "fails_when_all_predicates_fail",
			Doc:  map[string]any{"foo": 3, "bar": 4},
			Operations: []jsonpatch.Operation{
				{
					Op:   "or",
					Path: "",
					Apply: []jsonpatch.Operation{
						{Op: "test", Path: "/foo", Value: 1},
						{Op: "test", Path: "/bar", Value: 2},
					},
				},
			},
			WantErr: true,
			Comment: "OR should fail when all predicates fail",
		},
	}

	testutils.RunMultiOperationTestCases(t, testCases)
}
