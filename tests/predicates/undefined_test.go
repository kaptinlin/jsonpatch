package ops_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestUndefinedOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("throws when value is defined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, "hello", patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})

		t.Run("succeeds when value is undefined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, map[string]any{}, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("throws when property is defined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "/foo",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, map[string]any{"foo": "bar"}, patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})

		t.Run("succeeds when property is not defined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, map[string]any{"foo": "bar"}, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("throws when index is defined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "/0",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, []any{"hello"}, patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})

		t.Run("succeeds when index is not defined", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:   "undefined",
				Path: "/5",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, []any{"hello"}, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})
	})
}
