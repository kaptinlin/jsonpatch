package ops_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestLessOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when value is lower than requested", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "less",
				Path:  "",
				Value: 124,
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, 123, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})

		t.Run("fails when value is not lower than requested", func(t *testing.T) {
			t.Parallel()
			op1 := internal.Operation{
				Op:    "less",
				Path:  "",
				Value: 123,
			}
			patch1 := []internal.Operation{op1}
			_, err1 := applyPatch(t, 123, patch1)
			if err1 == nil {
				require.FailNow(t, "Apply() error = nil, want error for equal value")
			}

			op2 := internal.Operation{
				Op:    "less",
				Path:  "",
				Value: 1,
			}
			patch2 := []internal.Operation{op2}
			_, err2 := applyPatch(t, 123, patch2)
			if err2 == nil {
				require.FailNow(t, "Apply() error = nil, want error for smaller value")
			}
		})
	})
}
