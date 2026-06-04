package ops_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestStartsOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, "Hello world", patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, "Hello world", patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})

		t.Run("can ignore case", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:         "starts",
				Path:       "",
				Value:      "hello",
				IgnoreCase: true,
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, "Hello world", patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "/msg",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, map[string]any{"msg": "Hello world"}, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "/msg",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, map[string]any{"msg": "Hello world"}, patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "/0",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, []any{"Hello world"}, patch)
			if err != nil {
				require.FailNow(t, fmt.Sprintf("Apply() error = %v, want nil", err))
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "starts",
				Path:  "/0",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := applyPatch(t, []any{"Hello world"}, patch)
			if err == nil {
				require.FailNow(t, "Apply() error = nil, want error")
			}
		})
	})
}
