package json

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestDecodeOperationFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  map[string]any
		want internal.Operation
	}{
		{name: "add", raw: map[string]any{"op": "add", "path": "/profile/name", "value": "Ada"}, want: internal.Operation{Op: "add", Path: "/profile/name", Value: "Ada"}},
		{name: "remove", raw: map[string]any{"op": "remove", "path": "/profile/name"}, want: internal.Operation{Op: "remove", Path: "/profile/name"}},
		{name: "remove with old value", raw: map[string]any{"op": "remove", "path": "/profile/name", "oldValue": "Ada"}, want: internal.Operation{Op: "remove", Path: "/profile/name", OldValue: "Ada"}},
		{name: "replace with old value", raw: map[string]any{"op": "replace", "path": "/profile/name", "value": "Grace", "oldValue": "Ada"}, want: internal.Operation{Op: "replace", Path: "/profile/name", Value: "Grace", OldValue: "Ada"}},
		{name: "move", raw: map[string]any{"op": "move", "path": "/profile/displayName", "from": "/profile/name"}, want: internal.Operation{Op: "move", Path: "/profile/displayName", From: "/profile/name"}},
		{name: "copy", raw: map[string]any{"op": "copy", "path": "/profile/alias", "from": "/profile/name"}, want: internal.Operation{Op: "copy", Path: "/profile/alias", From: "/profile/name"}},
		{name: "test with not", raw: map[string]any{"op": "test", "path": "/profile/name", "value": "Ada", "not": true}, want: internal.Operation{Op: "test", Path: "/profile/name", Value: "Ada", Not: true}},
		{name: "flip", raw: map[string]any{"op": "flip", "path": "/enabled"}, want: internal.Operation{Op: "flip", Path: "/enabled"}},
		{name: "inc", raw: map[string]any{"op": "inc", "path": "/count", "inc": 2}, want: internal.Operation{Op: "inc", Path: "/count", Inc: 2}},
		{name: "str_ins", raw: map[string]any{"op": "str_ins", "path": "/profile/name", "pos": 1, "str": "d"}, want: internal.Operation{Op: "str_ins", Path: "/profile/name", Pos: 1, Str: "d"}},
		{name: "str_del string mode", raw: map[string]any{"op": "str_del", "path": "/profile/name", "pos": 1, "str": "da"}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Str: "da"}},
		{name: "str_del length mode", raw: map[string]any{"op": "str_del", "path": "/profile/name", "pos": 1, "len": 2}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Len: 2}},
		{name: "split", raw: map[string]any{"op": "split", "path": "/nodes/0", "pos": 1, "props": map[string]any{"kind": "paragraph"}}, want: internal.Operation{Op: "split", Path: "/nodes/0", Pos: 1, Props: map[string]any{"kind": "paragraph"}}},
		{name: "merge", raw: map[string]any{"op": "merge", "path": "/nodes/1", "pos": 1, "props": map[string]any{"merged": true}}, want: internal.Operation{Op: "merge", Path: "/nodes/1", Pos: 1, Props: map[string]any{"merged": true}}},
		{name: "extend", raw: map[string]any{"op": "extend", "path": "/profile", "props": map[string]any{"name": "Ada"}, "deleteNull": true}, want: internal.Operation{Op: "extend", Path: "/profile", Props: map[string]any{"name": "Ada"}, DeleteNull: true}},
		{name: "defined", raw: map[string]any{"op": "defined", "path": "/profile/name"}, want: internal.Operation{Op: "defined", Path: "/profile/name"}},
		{name: "undefined", raw: map[string]any{"op": "undefined", "path": "/profile/deleted"}, want: internal.Operation{Op: "undefined", Path: "/profile/deleted"}},
		{name: "type", raw: map[string]any{"op": "type", "path": "/profile/name", "value": "string"}, want: internal.Operation{Op: "type", Path: "/profile/name", Value: "string"}},
		{name: "test_type single", raw: map[string]any{"op": "test_type", "path": "/profile/name", "type": "string"}, want: internal.Operation{Op: "test_type", Path: "/profile/name", Type: "string"}},
		{name: "test_type value fallback", raw: map[string]any{"op": "test_type", "path": "/profile/name", "value": []any{"string", "null"}}, want: internal.Operation{Op: "test_type", Path: "/profile/name", Value: []string{"string", "null"}}},
		{name: "test_string", raw: map[string]any{"op": "test_string", "path": "/profile/name", "str": "da", "pos": 1, "not": true, "ignore_case": true}, want: internal.Operation{Op: "test_string", Path: "/profile/name", Str: "da", Pos: 1, Not: true, IgnoreCase: true}},
		{name: "test_string_len", raw: map[string]any{"op": "test_string_len", "path": "/profile/name", "len": 3, "not": true}, want: internal.Operation{Op: "test_string_len", Path: "/profile/name", Len: 3, Not: true}},
		{name: "contains", raw: map[string]any{"op": "contains", "path": "/profile/name", "value": "Ad", "ignore_case": true}, want: internal.Operation{Op: "contains", Path: "/profile/name", Value: "Ad", IgnoreCase: true}},
		{name: "starts", raw: map[string]any{"op": "starts", "path": "/profile/name", "value": "A", "ignore_case": true}, want: internal.Operation{Op: "starts", Path: "/profile/name", Value: "A", IgnoreCase: true}},
		{name: "ends", raw: map[string]any{"op": "ends", "path": "/profile/name", "value": "a", "ignore_case": true}, want: internal.Operation{Op: "ends", Path: "/profile/name", Value: "a", IgnoreCase: true}},
		{name: "matches", raw: map[string]any{"op": "matches", "path": "/profile/name", "value": "^ad", "ignore_case": true}, want: internal.Operation{Op: "matches", Path: "/profile/name", Value: "^ad", IgnoreCase: true}},
		{name: "in wraps scalar", raw: map[string]any{"op": "in", "path": "/role", "value": "admin"}, want: internal.Operation{Op: "in", Path: "/role", Value: []any{"admin"}}},
		{name: "less", raw: map[string]any{"op": "less", "path": "/score", "value": 10}, want: internal.Operation{Op: "less", Path: "/score", Value: 10}},
		{name: "more", raw: map[string]any{"op": "more", "path": "/score", "value": 5}, want: internal.Operation{Op: "more", Path: "/score", Value: 5}},
		{
			name: "and merges paths",
			raw: map[string]any{"op": "and", "path": "/profile", "apply": []any{
				map[string]any{"op": "defined", "path": "/name"},
				map[string]any{"op": "contains", "path": "/role", "value": "admin"},
			}},
			want: internal.Operation{Op: "and", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile/name"},
				{Op: "contains", Path: "/profile/role", Value: "admin"},
			}},
		},
		{
			name: "or equal path remains sub path",
			raw: map[string]any{"op": "or", "path": "/profile", "apply": []any{
				map[string]any{"op": "defined", "path": "/profile"},
			}},
			want: internal.Operation{Op: "or", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile"},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := Decode([]map[string]any{tc.raw}, PatchOptions{})
			require.NoError(t, err)
			require.Len(t, decoded, 1)

			got, err := decoded[0].ToJSON()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("decoded operation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDecodeRejectsInvalidOperationPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     map[string]any
		wantErr error
	}{
		{name: "missing op", raw: map[string]any{"path": "/x"}, wantErr: ErrOpMissingOpField},
		{name: "missing path", raw: map[string]any{"op": "add", "value": "x"}, wantErr: ErrOpMissingPathField},
		{name: "invalid pointer", raw: map[string]any{"op": "add", "path": "x", "value": "x"}, wantErr: ErrInvalidPointer},
		{name: "unknown operation", raw: map[string]any{"op": "unknown", "path": "/x"}, wantErr: ErrCodecOpUnknown},
		{name: "add missing value", raw: map[string]any{"op": "add", "path": "/x"}, wantErr: ErrAddOpMissingValue},
		{name: "replace missing value", raw: map[string]any{"op": "replace", "path": "/x"}, wantErr: ErrReplaceOpMissingValue},
		{name: "move missing from", raw: map[string]any{"op": "move", "path": "/x"}, wantErr: ErrMoveOpMissingFrom},
		{name: "copy missing from", raw: map[string]any{"op": "copy", "path": "/x"}, wantErr: ErrCopyOpMissingFrom},
		{name: "inc missing value", raw: map[string]any{"op": "inc", "path": "/x"}, wantErr: ErrIncOpMissingInc},
		{name: "inc invalid value", raw: map[string]any{"op": "inc", "path": "/x", "inc": "one"}, wantErr: ErrIncOpInvalidType},
		{name: "str_ins missing pos", raw: map[string]any{"op": "str_ins", "path": "/x", "str": "a"}, wantErr: ErrStrInsOpMissingPos},
		{name: "str_ins missing str", raw: map[string]any{"op": "str_ins", "path": "/x", "pos": 1}, wantErr: ErrStrInsOpMissingStr},
		{name: "str_del missing pos", raw: map[string]any{"op": "str_del", "path": "/x", "len": 1}, wantErr: ErrStrDelOpMissingPos},
		{name: "str_del missing len or str", raw: map[string]any{"op": "str_del", "path": "/x", "pos": 1, "len": struct{}{}}, wantErr: ErrStrDelOpMissingFields},
		{name: "split missing pos", raw: map[string]any{"op": "split", "path": "/x"}, wantErr: ErrSplitOpMissingPos},
		{name: "extend props not object", raw: map[string]any{"op": "extend", "path": "/x", "props": "props"}, wantErr: ErrValueNotObject},
		{name: "test missing value", raw: map[string]any{"op": "test", "path": "/x"}, wantErr: ErrMissingValueField},
		{name: "type missing value", raw: map[string]any{"op": "type", "path": "/x"}, wantErr: ErrTypeOpMissingValue},
		{name: "test_type missing type", raw: map[string]any{"op": "test_type", "path": "/x"}, wantErr: ErrTestTypeOpMissingType},
		{name: "test_type invalid scalar", raw: map[string]any{"op": "test_type", "path": "/x", "type": "invalid"}, wantErr: ErrInvalidType},
		{name: "test_type empty array", raw: map[string]any{"op": "test_type", "path": "/x", "type": []any{}}, wantErr: ErrEmptyTypeList},
		{name: "test_type invalid item", raw: map[string]any{"op": "test_type", "path": "/x", "type": []any{"string", 1}}, wantErr: ErrInvalidType},
		{name: "test_type empty string array", raw: map[string]any{"op": "test_type", "path": "/x", "type": []string{}}, wantErr: ErrEmptyTypeList},
		{name: "test_type invalid string array item", raw: map[string]any{"op": "test_type", "path": "/x", "type": []string{"string", "invalid"}}, wantErr: ErrInvalidType},
		{name: "test_string missing str", raw: map[string]any{"op": "test_string", "path": "/x"}, wantErr: ErrTestStringOpMissingStr},
		{name: "test_string_len invalid len", raw: map[string]any{"op": "test_string_len", "path": "/x", "len": struct{}{}}, wantErr: ErrTestStringLenOpMissingLen},
		{name: "contains missing value", raw: map[string]any{"op": "contains", "path": "/x"}, wantErr: ErrContainsOpMissingValue},
		{name: "ends missing value", raw: map[string]any{"op": "ends", "path": "/x"}, wantErr: ErrEndsOpMissingValue},
		{name: "starts missing value", raw: map[string]any{"op": "starts", "path": "/x"}, wantErr: ErrStartsOpMissingValue},
		{name: "matches missing value", raw: map[string]any{"op": "matches", "path": "/x"}, wantErr: ErrMatchesOpMissingValue},
		{name: "less invalid value", raw: map[string]any{"op": "less", "path": "/x", "value": struct{}{}}, wantErr: ErrLessOpMissingValue},
		{name: "more invalid value", raw: map[string]any{"op": "more", "path": "/x", "value": struct{}{}}, wantErr: ErrMoreOpMissingValue},
		{name: "and missing apply", raw: map[string]any{"op": "and", "path": "/x"}, wantErr: ErrAndOpMissingApply},
		{name: "or missing apply", raw: map[string]any{"op": "or", "path": "/x"}, wantErr: ErrOrOpMissingApply},
		{name: "not missing apply", raw: map[string]any{"op": "not", "path": "/x"}, wantErr: ErrNotOpMissingApply},
		{name: "not empty apply", raw: map[string]any{"op": "not", "path": "/x", "apply": []any{}}, wantErr: ErrNotOpRequiresOperand},
		{name: "nested not invalid operand", raw: map[string]any{"op": "and", "path": "/x", "apply": []any{map[string]any{"op": "not", "path": "/y", "apply": []any{"bad"}}}}, wantErr: ErrNotOpRequiresValidOperand},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := Decode([]map[string]any{tc.raw}, PatchOptions{})
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}
