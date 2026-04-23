package testutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestHelperFunctions(t *testing.T) {
	t.Parallel()

	t.Run("apply operation helpers", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"name": "Ada"}
		operation := jsonpatch.Operation{Op: "replace", Path: "/name", Value: "Grace"}
		result := ApplyOperation(t, doc, operation)
		assert.Equal(t, map[string]any{"name": "Grace"}, result)

		err := ApplyOperationWithError(t, map[string]any{"name": "Ada"}, jsonpatch.Operation{Op: "remove", Path: "/missing"})
		require.Error(t, err)
	})

	t.Run("apply operations helpers", func(t *testing.T) {
		t.Parallel()

		result := ApplyOperations(t, map[string]any{"name": "Ada"}, []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Grace"},
			{Op: "add", Path: "/role", Value: "admin"},
		})
		assert.Equal(t, map[string]any{"name": "Grace", "role": "admin"}, result)

		err := ApplyOperationsWithError(t, map[string]any{"name": "Ada"}, []jsonpatch.Operation{{Op: "remove", Path: "/missing"}})
		require.Error(t, err)
	})

	t.Run("internal operation helpers", func(t *testing.T) {
		t.Parallel()

		result := ApplyInternalOp(t, map[string]any{"name": "Ada"}, internal.Operation{Op: "replace", Path: "/name", Value: "Grace"})
		assert.Equal(t, map[string]any{"name": "Grace"}, result)

		ApplyInternalOpWithError(t, map[string]any{"name": "Ada"}, internal.Operation{Op: "remove", Path: "/missing"})

		result = ApplyInternalOps(t, map[string]any{"name": "Ada"}, []internal.Operation{{Op: "add", Path: "/role", Value: "admin"}})
		assert.Equal(t, map[string]any{"name": "Ada", "role": "admin"}, result)

		ApplyInternalOpsWithError(t, map[string]any{"name": "Ada"}, []internal.Operation{{Op: "remove", Path: "/missing"}})
	})
}

func TestRunCaseHelpers(t *testing.T) {
	t.Parallel()

	RunTestCase(t, TestCase{
		Name:      "single success",
		Doc:       map[string]any{"name": "Ada"},
		Operation: jsonpatch.Operation{Op: "replace", Path: "/name", Value: "Grace"},
		Expected:  map[string]any{"name": "Grace"},
	})

	RunMultiOperationTestCase(t, MultiOperationTestCase{
		Name: "multi success",
		Doc:  map[string]any{"name": "Ada"},
		Operations: []jsonpatch.Operation{
			{Op: "replace", Path: "/name", Value: "Grace"},
			{Op: "add", Path: "/role", Value: "admin"},
		},
		Expected: map[string]any{"name": "Grace", "role": "admin"},
	})

	RunTestCases(t, []TestCase{{
		Name:      "run test cases",
		Doc:       map[string]any{"name": "Ada"},
		Operation: jsonpatch.Operation{Op: "replace", Path: "/name", Value: "Grace"},
		Expected:  map[string]any{"name": "Grace"},
	}})

	RunMultiOperationTestCases(t, []MultiOperationTestCase{{
		Name:       "run multi test cases",
		Doc:        map[string]any{"name": "Ada"},
		Operations: []jsonpatch.Operation{{Op: "add", Path: "/role", Value: "admin"}},
		Expected:   map[string]any{"name": "Ada", "role": "admin"},
	}})
}
