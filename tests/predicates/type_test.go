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
				value interface{}
				typ   string
			}{
				{1, "number"},
				{"hello", "string"},
				{true, "boolean"},
				{nil, "null"},
				{[]interface{}{}, "array"},
				{map[string]interface{}{}, "object"},
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
				value interface{}
				typ   string
			}{
				{1, "string"},
				{"hello", "number"},
				{true, "null"},
				{nil, "boolean"},
				{[]interface{}{}, "object"},
				{map[string]interface{}{}, "array"},
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
				obj map[string]interface{}
				typ string
			}{
				{map[string]interface{}{"foo": 1}, "number"},
				{map[string]interface{}{"foo": "hello"}, "string"},
				{map[string]interface{}{"foo": true}, "boolean"},
				{map[string]interface{}{"foo": nil}, "null"},
				{map[string]interface{}{"foo": []interface{}{}}, "array"},
				{map[string]interface{}{"foo": map[string]interface{}{}}, "object"},
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
				obj map[string]interface{}
				typ string
			}{
				{map[string]interface{}{"foo": 1}, "string"},
				{map[string]interface{}{"foo": "hello"}, "number"},
				{map[string]interface{}{"foo": true}, "null"},
				{map[string]interface{}{"foo": nil}, "boolean"},
				{map[string]interface{}{"foo": []interface{}{}}, "object"},
				{map[string]interface{}{"foo": map[string]interface{}{}}, "array"},
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
				arr []interface{}
				typ string
			}{
				{[]interface{}{1}, "number"},
				{[]interface{}{"hello"}, "string"},
				{[]interface{}{true}, "boolean"},
				{[]interface{}{nil}, "null"},
				{[]interface{}{[]interface{}{}}, "array"},
				{[]interface{}{map[string]interface{}{}}, "object"},
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
				arr []interface{}
				typ string
			}{
				{[]interface{}{1}, "string"},
				{[]interface{}{"hello"}, "number"},
				{[]interface{}{true}, "null"},
				{[]interface{}{nil}, "boolean"},
				{[]interface{}{[]interface{}{}}, "object"},
				{[]interface{}{map[string]interface{}{}}, "array"},
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
