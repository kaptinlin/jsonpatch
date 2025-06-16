package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestStartsOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "",
				"value": "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.Error(t, err)
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				"op":          "starts",
				"path":        "",
				"value":       "hello",
				"ignore_case": true,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "/msg",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "/msg",
				"value": "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "/0",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "starts",
				"path":  "/0",
				"value": "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})
}
