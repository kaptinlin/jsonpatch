package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
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
			_, err := jsonpatch.ApplyPatch(123, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
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
			_, err1 := jsonpatch.ApplyPatch(123, patch1, internal.WithMutate(true))
			if err1 == nil {
				t.Fatal("ApplyPatch() error = nil, want error for equal value")
			}

			op2 := internal.Operation{
				Op:    "less",
				Path:  "",
				Value: 1,
			}
			patch2 := []internal.Operation{op2}
			_, err2 := jsonpatch.ApplyPatch(123, patch2, internal.WithMutate(true))
			if err2 == nil {
				t.Fatal("ApplyPatch() error = nil, want error for smaller value")
			}
		})
	})
}
