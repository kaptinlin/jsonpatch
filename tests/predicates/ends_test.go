package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestEndsOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("can ignore case", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:         "ends",
				Path:       "",
				Value:      "WORLD",
				IgnoreCase: true,
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "/msg",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]any{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "/msg",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]any{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when string ends with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "/0",
				Value: "world",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]any{"Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not end with suffix", func(t *testing.T) {
			t.Parallel()
			op := internal.Operation{
				Op:    "ends",
				Path:  "/0",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]any{"Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})
	})
}
