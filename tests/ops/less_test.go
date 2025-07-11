package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestLessOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when value is lower than requested", func(t *testing.T) {
			op := internal.Operation{
				"op":    "less",
				"path":  "",
				"value": 124,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(123, patch, internal.WithMutate(true))
			require.NoError(t, err)
		})

		t.Run("fails when value is not lower than requested", func(t *testing.T) {
			op1 := internal.Operation{
				"op":    "less",
				"path":  "",
				"value": 123,
			}
			patch1 := []internal.Operation{op1}
			_, err1 := jsonpatch.ApplyPatch(123, patch1, internal.WithMutate(true))
			require.Error(t, err1)

			op2 := internal.Operation{
				"op":    "less",
				"path":  "",
				"value": 1,
			}
			patch2 := []internal.Operation{op2}
			_, err2 := jsonpatch.ApplyPatch(123, patch2, internal.WithMutate(true))
			require.Error(t, err2)
		})
	})
}
