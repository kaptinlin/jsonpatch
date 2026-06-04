package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func applyPatch[T jsonpatch.Document](t *testing.T, doc T, operations []internal.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	return testutils.ApplyInternalOperationsResult(t, doc, operations)
}
