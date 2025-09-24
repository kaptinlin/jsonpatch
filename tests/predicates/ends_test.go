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
				Op: "ends",
				Path: "",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.Error(t, err)
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "",
				Value: "WORLD",
				IgnoreCase: true,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			require.NoError(t, err)
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "/msg",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "/msg",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "/0",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			op := internal.Operation{
				Op: "ends",
				Path: "/0",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			require.Error(t, err)
		})
	})
}
