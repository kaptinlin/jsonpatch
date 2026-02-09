package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestStartsOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("Hello world", patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("can ignore case", func(t *testing.T) {
			op := internal.Operation{
				Op:         "starts",
				Path:       "",
				Value:      "hello",
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
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "/msg",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "/msg",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"msg": "Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("succeeds when string starts with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "/0",
				Value: "Hello",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})

		t.Run("throws when string does not start with prefix", func(t *testing.T) {
			op := internal.Operation{
				Op:    "starts",
				Path:  "/0",
				Value: "World",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"Hello world"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})
	})
}
