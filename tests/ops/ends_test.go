package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestEndsOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "",
				"value": "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.Error(t, err)
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				"op":          "ends",
				"path":        "",
				"value":       "WORLD",
				"ignore_case": true,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "/msg",
				"value": "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "/msg",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "/0",
				"value": "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				"op":    "ends",
				"path":  "/0",
				"value": "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})
}
