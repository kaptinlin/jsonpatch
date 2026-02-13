package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInOp(t *testing.T) {
	t.Parallel()
	t.Run("positive", func(t *testing.T) {
		t.Parallel()
		t.Run("should test against root (on a json document of type object) - and return true", func(t *testing.T) {
			t.Parallel()
			obj := map[string]any{
				"hello": "world",
			}
			op := internal.Operation{
				Op:   "in",
				Path: "",
				Value: []any{
					1,
					map[string]any{
						"hello": "world",
					},
				},
			}
			result := testutils.ApplyInternalOp(t, obj, op)
			assert.Equal(t, obj, result)
		})

		t.Run("should test against root (on a json document of type object) - and return false", func(t *testing.T) {
			t.Parallel()
			obj := map[string]any{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "in",
				Path:  "",
				Value: []any{1},
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})

		t.Run("should test against root (on a json document of type array) - and return false", func(t *testing.T) {
			t.Parallel()
			obj := []any{
				map[string]any{
					"hello": "world",
				},
			}
			op := internal.Operation{
				Op:    "in",
				Path:  "",
				Value: []any{1},
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})
}
