package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/internal"
)

func TestUndefinedOp(t *testing.T) {
	t.Run("root", func(t *testing.T) {
		t.Run("throws when value is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch("hello", patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("succeeds when value is undefined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Run("throws when property is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "/foo",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("succeeds when property is not defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "/missing",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch(map[string]interface{}{"foo": "bar"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Run("throws when index is defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "/0",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.WithMutate(true))
			if err == nil {
				t.Fatal("ApplyPatch() error = nil, want error")
			}
		})

		t.Run("succeeds when index is not defined", func(t *testing.T) {
			op := internal.Operation{
				Op:   "undefined",
				Path: "/5",
			}
			patch := []internal.Operation{op}
			_, err := jsonpatch.ApplyPatch([]interface{}{"hello"}, patch, internal.WithMutate(true))
			if err != nil {
				t.Fatalf("ApplyPatch() error = %v, want nil", err)
			}
		})
	})
}
