package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func TestTestOp(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		t.Run("should test against root on json document of type object and return true", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: map[string]interface{}{"hello": "world"},
			}
			result := testutils.ApplyInternalOp(t, obj, op)
			assert.Equal(t, obj, result)
		})

		t.Run("should test against root on json document of type object and return false", func(t *testing.T) {
			obj := map[string]interface{}{
				"hello": "world",
			}
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 1,
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})

		t.Run("should test against root on json document of type array and return false", func(t *testing.T) {
			obj := []interface{}{
				map[string]interface{}{
					"hello": "world",
				},
			}
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 1,
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})

		t.Run("should throw against root", func(t *testing.T) {
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 2,
				Not:   false,
			}
			testutils.ApplyInternalOpWithError(t, 1, op)
		})

		t.Run("should throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 2,
				Not:   false,
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})

		t.Run("should not throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 1,
				Not:   false,
			}
			testutils.ApplyInternalOp(t, obj, op)
		})
	})

	t.Run("negative", func(t *testing.T) {
		t.Run("should test against root", func(t *testing.T) {
			op := internal.Operation{
				Op:    "test",
				Path:  "",
				Value: 2,
				Not:   true,
			}
			testutils.ApplyInternalOp(t, 1, op)
		})

		t.Run("should not throw when object key is different", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 2,
				Not:   true,
			}
			testutils.ApplyInternalOp(t, obj, op)
		})

		t.Run("should throw when object key is the same", func(t *testing.T) {
			obj := map[string]interface{}{"foo": 1}
			op := internal.Operation{
				Op:    "test",
				Path:  "/foo",
				Value: 1,
				Not:   true,
			}
			testutils.ApplyInternalOpWithError(t, obj, op)
		})
	})
}
