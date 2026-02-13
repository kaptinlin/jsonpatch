package ops_test

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/tests/testutils"
)

func TestTypeOp(t *testing.T) {
	t.Parallel()
	t.Run("root", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when type matches", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				value any
				typ   string
			}{
				{1, "number"},
				{"hello", "string"},
				{true, "boolean"},
				{nil, "null"},
				{[]any{}, "array"},
				{map[string]any{}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "",
					Value: test.typ,
				}
				testutils.ApplyInternalOp(t, test.value, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				value any
				typ   string
			}{
				{1, "string"},
				{"hello", "number"},
				{true, "null"},
				{nil, "boolean"},
				{[]any{}, "object"},
				{map[string]any{}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "",
					Value: test.typ,
				}
				testutils.ApplyInternalOpWithError(t, test.value, op)
			}
		})
	})

	t.Run("object", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when type matches", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				obj map[string]any
				typ string
			}{
				{map[string]any{"foo": 1}, "number"},
				{map[string]any{"foo": "hello"}, "string"},
				{map[string]any{"foo": true}, "boolean"},
				{map[string]any{"foo": nil}, "null"},
				{map[string]any{"foo": []any{}}, "array"},
				{map[string]any{"foo": map[string]any{}}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "/foo",
					Value: test.typ,
				}
				testutils.ApplyInternalOp(t, test.obj, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				obj map[string]any
				typ string
			}{
				{map[string]any{"foo": 1}, "string"},
				{map[string]any{"foo": "hello"}, "number"},
				{map[string]any{"foo": true}, "null"},
				{map[string]any{"foo": nil}, "boolean"},
				{map[string]any{"foo": []any{}}, "object"},
				{map[string]any{"foo": map[string]any{}}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "/foo",
					Value: test.typ,
				}
				testutils.ApplyInternalOpWithError(t, test.obj, op)
			}
		})
	})

	t.Run("array", func(t *testing.T) {
		t.Parallel()
		t.Run("succeeds when type matches", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				arr []any
				typ string
			}{
				{[]any{1}, "number"},
				{[]any{"hello"}, "string"},
				{[]any{true}, "boolean"},
				{[]any{nil}, "null"},
				{[]any{[]any{}}, "array"},
				{[]any{map[string]any{}}, "object"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "/0",
					Value: test.typ,
				}
				testutils.ApplyInternalOp(t, test.arr, op)
			}
		})

		t.Run("throws when type does not match", func(t *testing.T) {
			t.Parallel()
			tests := []struct {
				arr []any
				typ string
			}{
				{[]any{1}, "string"},
				{[]any{"hello"}, "number"},
				{[]any{true}, "null"},
				{[]any{nil}, "boolean"},
				{[]any{[]any{}}, "object"},
				{[]any{map[string]any{}}, "array"},
			}

			for _, test := range tests {
				op := internal.Operation{
					Op:    "type",
					Path:  "/0",
					Value: test.typ,
				}
				testutils.ApplyInternalOpWithError(t, test.arr, op)
			}
		})
	})
}
