package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestUndefinedOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("throws when value is defined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("hello", patch, internal.ApplyPatchOptions{Mutate: true})
			require.Error(t, err)
		})

		t.Run("succeeds when value is undefined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{}, patch, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("throws when property is defined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "/foo",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.ApplyPatchOptions{Mutate: true})
			require.Error(t, err)
		})

		t.Run("succeeds when property is not defined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("throws when index is defined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "/0",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.ApplyPatchOptions{Mutate: true})
			require.Error(t, err)
		})

		t.Run("succeeds when index is not defined", func(t *testing.T) {
			op := internal.Operation{
				"op":   "undefined",
				"path": "/5",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
		})
	})
}
