package ops_test

import (
	"fmt"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
)

func TestMoreOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when value is higher than requested", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "more",
				Path:  "",
				Value: 99,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(123, patch, internal.WithMutate(true))
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyPatch() error = %v, want nil", err))
			}
		})

		t.Run("fails when value is not higher than requested", func(t *testing.T) {
			t.Parallel()
			op1 := internal.Operation{
				Op:    "more",
				Path:  "",
				Value: 123,
			}
			patch1 := []internal.Operation{op1}
			_, err1 := jsonpatch.ApplyPatch(123, patch1, internal.WithMutate(true))
			if err1 == nil {
				require.FailNow(t, "ApplyPatch() error = nil, want error for equal value")
			}

			op2 := internal.Operation{
				Op:    "more",
				Path:  "",
				Value: 124,
			}
			patch2 := []internal.Operation{op2}
			_, err2 := jsonpatch.ApplyPatch(123, patch2, internal.WithMutate(true))
			if err2 == nil {
				require.FailNow(t, "ApplyPatch() error = nil, want error for higher value")
			}
		})
	})
}
