package ops_test

import (
	"fmt"
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/stretchr/testify/require"
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
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyPatch() error = %v, want nil", err))
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
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err == nil {
				require.FailNow(t, "ApplyPatch() error = nil, want error")
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
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyPatch() error = %v, want nil", err))
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
			_, err := jsonpatch.ApplyPatch(map[string]any{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyPatch() error = %v, want nil", err))
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
			_, err := jsonpatch.ApplyPatch(map[string]any{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				require.FailNow(t, "ApplyPatch() error = nil, want error")
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
			_, err := jsonpatch.ApplyPatch([]any{"Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				require.FailNow(t, fmt.Sprintf("ApplyPatch() error = %v, want nil", err))
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
			_, err := jsonpatch.ApplyPatch([]any{"Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				require.FailNow(t, "ApplyPatch() error = nil, want error")
			}
		})
	})
}
