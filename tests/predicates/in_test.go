package ops_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestInOp(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("should test against root (on a json document of type object) - and return true", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:   "in",
				Path: "",
				Value: []interface{}{
					1,
					map[string]interface{}{
						"hello": "world",
					},
				},
			}
			result := testutils.ApplyInternalOp(t, obj, op)
			if diff := cmp.Diff(obj, result); diff != "" {
				t.Errorf("ApplyInternalOp() mismatch (-want +got):\n%s", diff)
			}
		})

		t.Run("should test against root (on a json document of type object) - and return false", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "in",
				Path:  "",
				Value: []interface{}{1},
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})

		t.Run("should test against root (on a json document of type array) - and return false", func(t *testing.T) {
			obj := []interface{}{
				map[string]interface{}{
					"hello": "world",
				},
			}
			op := internal.Operation{
				Op:    "in",
				Path:  "",
				Value: []interface{}{1},
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})
}
