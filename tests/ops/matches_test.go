package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestMatchesOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when matches correctly a substring", func(t *testing.T) {
			op := internal.Operation{
				"op":    "matches",
				"path":  "",
				"value": "\\d+",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("123", patch, internal.ApplyPatchOptions{Mutate: true})
			require.NoError(t, err)
		})

		t.Run("fails when does not match the string", func(t *testing.T) {
			op := internal.Operation{
				"op":    "matches",
				"path":  "",
				"value": "\\d+",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("asdf", patch, internal.ApplyPatchOptions{Mutate: true})
			require.Error(t, err)
		})
	})
}
