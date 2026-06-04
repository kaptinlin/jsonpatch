package reference

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func applyPatch[T jsonpatch.Document](t testing.TB, doc T, operations []jsoncodec.Operation) (*jsonpatch.Result[T], error) {
	t.Helper()
	return testutils.ApplyOperationsResult(t, doc, operations)
}
