package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestDefinedOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when value is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("hello", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when value is undefined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when property is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "/foo",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when property is not defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when index is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "/0",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when index is not defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "defined",
				Path: "/5",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})
}
