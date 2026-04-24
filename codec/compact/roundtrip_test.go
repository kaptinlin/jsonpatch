package compact

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestDecodeOperationFamilies(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw  Op
		want internal.Operation
	}{
		{name: "add", raw: Op{CodeAdd, "/profile/name", "Ada"}, want: internal.Operation{Op: "add", Path: "/profile/name", Value: "Ada"}},
		{name: "remove", raw: Op{CodeRemove, "/profile/name"}, want: internal.Operation{Op: "remove", Path: "/profile/name"}},
		{name: "remove with old value", raw: Op{CodeRemove, "/profile/name", "Ada"}, want: internal.Operation{Op: "remove", Path: "/profile/name", OldValue: "Ada"}},
		{name: "replace", raw: Op{CodeReplace, "/profile/name", "Grace"}, want: internal.Operation{Op: "replace", Path: "/profile/name", Value: "Grace"}},
		{name: "replace with old value", raw: Op{CodeReplace, "/profile/name", "Grace", "Ada"}, want: internal.Operation{Op: "replace", Path: "/profile/name", Value: "Grace", OldValue: "Ada"}},
		{name: "move", raw: Op{CodeMove, "/profile/displayName", "/profile/name"}, want: internal.Operation{Op: "move", Path: "/profile/displayName", From: "/profile/name"}},
		{name: "copy", raw: Op{CodeCopy, "/profile/alias", "/profile/name"}, want: internal.Operation{Op: "copy", Path: "/profile/alias", From: "/profile/name"}},
		{name: "test", raw: Op{CodeTest, "/profile/name", "Ada", true}, want: internal.Operation{Op: "test", Path: "/profile/name", Value: "Ada", Not: true}},
		{name: "flip", raw: Op{CodeFlip, "/enabled"}, want: internal.Operation{Op: "flip", Path: "/enabled"}},
		{name: "inc", raw: Op{CodeInc, "/count", 2}, want: internal.Operation{Op: "inc", Path: "/count", Inc: 2}},
		{name: "str_ins", raw: Op{CodeStrIns, "/profile/name", 1, "d"}, want: internal.Operation{Op: "str_ins", Path: "/profile/name", Pos: 1, Str: "d"}},
		{name: "str_del string mode", raw: Op{CodeStrDel, "/profile/name", 1, "da"}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Str: "da"}},
		{name: "str_del length mode", raw: Op{CodeStrDel, "/profile/name", 1, 0, 2}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Len: 2}},
		{name: "split", raw: Op{CodeSplit, "/nodes/0", 1, map[string]any{"kind": "paragraph"}}, want: internal.Operation{Op: "split", Path: "/nodes/0", Pos: 1, Props: map[string]any{"kind": "paragraph"}}},
		{name: "merge", raw: Op{CodeMerge, "/nodes/1", 1, map[string]any{"merged": true}}, want: internal.Operation{Op: "merge", Path: "/nodes/1", Pos: 1, Props: map[string]any{"merged": true}}},
		{name: "extend", raw: Op{CodeExtend, "/profile", map[string]any{"name": "Ada"}, true}, want: internal.Operation{Op: "extend", Path: "/profile", Props: map[string]any{"name": "Ada"}, DeleteNull: true}},
		{name: "defined", raw: Op{CodeDefined, "/profile/name"}, want: internal.Operation{Op: "defined", Path: "/profile/name"}},
		{name: "undefined", raw: Op{CodeUndefined, "/profile/deleted"}, want: internal.Operation{Op: "undefined", Path: "/profile/deleted"}},
		{name: "contains", raw: Op{CodeContains, "/profile/name", "Ad", true}, want: internal.Operation{Op: "contains", Path: "/profile/name", Value: "Ad", IgnoreCase: true}},
		{name: "starts", raw: Op{CodeStarts, "/profile/name", "A", true}, want: internal.Operation{Op: "starts", Path: "/profile/name", Value: "A", IgnoreCase: true}},
		{name: "ends", raw: Op{CodeEnds, "/profile/name", "a", true}, want: internal.Operation{Op: "ends", Path: "/profile/name", Value: "a", IgnoreCase: true}},
		{name: "type", raw: Op{CodeType, "/profile/name", "string"}, want: internal.Operation{Op: "type", Path: "/profile/name", Value: "string"}},
		{name: "test_type", raw: Op{CodeTestType, "/profile/name", []any{"string", "null"}}, want: internal.Operation{Op: "test_type", Path: "/profile/name", Value: []string{"string", "null"}}},
		{name: "test_string", raw: Op{CodeTestString, "/profile/name", "da", 1, true}, want: internal.Operation{Op: "test_string", Path: "/profile/name", Str: "da", Pos: 1, Not: true}},
		{name: "test_string_len", raw: Op{CodeTestStringLen, "/profile/name", 3, true}, want: internal.Operation{Op: "test_string_len", Path: "/profile/name", Len: 3, Not: true}},
		{name: "in", raw: Op{CodeIn, "/role", []any{"admin", "editor"}}, want: internal.Operation{Op: "in", Path: "/role", Value: []any{"admin", "editor"}}},
		{name: "less", raw: Op{CodeLess, "/score", 10}, want: internal.Operation{Op: "less", Path: "/score", Value: 10}},
		{name: "more", raw: Op{CodeMore, "/score", 5}, want: internal.Operation{Op: "more", Path: "/score", Value: 5}},
		{name: "matches", raw: Op{CodeMatches, "/profile/name", "^ad", true}, want: internal.Operation{Op: "matches", Path: "/profile/name", Value: "^ad", IgnoreCase: true}},
		{
			name: "and",
			raw:  Op{CodeAnd, "/profile", []any{[]any{CodeDefined, "/profile/name"}, []any{CodeContains, "/profile/role", "admin"}}},
			want: internal.Operation{Op: "and", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile/name"},
				{Op: "contains", Path: "/profile/role", Value: "admin"},
			}},
		},
		{
			name: "or",
			raw:  Op{CodeOr, "/profile", []any{[]any{CodeDefined, "/profile/name"}, []any{CodeDefined, "/profile/email"}}},
			want: internal.Operation{Op: "or", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile/name"},
				{Op: "defined", Path: "/profile/email"},
			}},
		},
		{
			name: "not",
			raw:  Op{CodeNot, "/profile", []any{[]any{CodeUndefined, "/profile/deleted"}}},
			want: internal.Operation{Op: "not", Path: "/profile", Apply: []internal.Operation{
				{Op: "undefined", Path: "/profile/deleted"},
			}},
		},
	}

	decoder := NewDecoder()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := decoder.Decode(tc.raw)
			require.NoError(t, err)

			got, err := decoded.ToJSON()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("decoded operation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestEncodeFormatsPathsAndOpcodes(t *testing.T) {
	t.Parallel()

	encoded, err := NewEncoder(WithStringOpcode(true)).Encode(op.NewMove([]string{"with/slash", "~tilde"}, []string{""}))
	require.NoError(t, err)

	want := Op{"move", "/with~1slash/~0tilde", "/"}
	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("encoded operation mismatch (-want +got):\n%s", diff)
	}
}

func TestDecodeRejectsInvalidOperationPayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     Op
		wantErr error
	}{
		{name: "too short", raw: Op{CodeAdd}, wantErr: ErrMinLength},
		{name: "path is not string", raw: Op{CodeAdd, 1, "x"}, wantErr: ErrPathNotString},
		{name: "unknown string opcode", raw: Op{"unknown", "/x"}, wantErr: ErrUnknownStringCode},
		{name: "unknown numeric opcode", raw: Op{999, "/x"}, wantErr: ErrUnknownNumericCode},
		{name: "unknown float opcode", raw: Op{999.0, "/x"}, wantErr: ErrUnknownNumericCode},
		{name: "invalid opcode type", raw: Op{true, "/x"}, wantErr: ErrInvalidCodeType},
		{name: "add missing value", raw: Op{CodeAdd, "/x"}, wantErr: ErrAddMissingValue},
		{name: "replace missing value", raw: Op{CodeReplace, "/x"}, wantErr: ErrReplaceMissingValue},
		{name: "move missing from", raw: Op{CodeMove, "/x"}, wantErr: ErrMoveMissingFrom},
		{name: "move from not string", raw: Op{CodeMove, "/x", 1}, wantErr: ErrMoveFromNotString},
		{name: "copy missing from", raw: Op{CodeCopy, "/x"}, wantErr: ErrCopyMissingFrom},
		{name: "copy from not string", raw: Op{CodeCopy, "/x", 1}, wantErr: ErrCopyFromNotString},
		{name: "test missing value", raw: Op{CodeTest, "/x"}, wantErr: ErrTestMissingValue},
		{name: "inc missing delta", raw: Op{CodeInc, "/x"}, wantErr: ErrIncMissingDelta},
		{name: "inc delta not number", raw: Op{CodeInc, "/x", "1"}, wantErr: ErrIncDeltaNotNumber},
		{name: "str_ins missing fields", raw: Op{CodeStrIns, "/x"}, wantErr: ErrStrInsMissingFields},
		{name: "str_ins pos not number", raw: Op{CodeStrIns, "/x", "1", "a"}, wantErr: ErrStrInsPosNotNumber},
		{name: "str_ins str not string", raw: Op{CodeStrIns, "/x", 1, false}, wantErr: ErrStrInsStrNotString},
		{name: "str_del missing fields", raw: Op{CodeStrDel, "/x"}, wantErr: ErrStrDelMissingFields},
		{name: "str_del pos not number", raw: Op{CodeStrDel, "/x", "1", 2}, wantErr: ErrStrDelPosNotNumber},
		{name: "str_del len not number", raw: Op{CodeStrDel, "/x", 1, false}, wantErr: ErrStrDelLenNotNumber},
		{name: "str_del explicit len not number", raw: Op{CodeStrDel, "/x", 1, 0, false}, wantErr: ErrStrDelLenNotNumber},
		{name: "split missing pos", raw: Op{CodeSplit, "/x"}, wantErr: ErrSplitMissingPos},
		{name: "split pos not number", raw: Op{CodeSplit, "/x", "1"}, wantErr: ErrSplitPosNotNumber},
		{name: "merge missing pos", raw: Op{CodeMerge, "/x"}, wantErr: ErrMergeMissingPos},
		{name: "merge pos not number", raw: Op{CodeMerge, "/x", "1"}, wantErr: ErrMergePosNotNumber},
		{name: "extend missing props", raw: Op{CodeExtend, "/x"}, wantErr: ErrExtendMissingProps},
		{name: "extend props not object", raw: Op{CodeExtend, "/x", "props"}, wantErr: ErrExtendPropsNotObject},
		{name: "contains missing value", raw: Op{CodeContains, "/x"}, wantErr: ErrContainsMissingValue},
		{name: "contains value not string", raw: Op{CodeContains, "/x", 1}, wantErr: ErrContainsValueNotString},
		{name: "starts missing value", raw: Op{CodeStarts, "/x"}, wantErr: ErrStartsMissingValue},
		{name: "starts value not string", raw: Op{CodeStarts, "/x", 1}, wantErr: ErrStartsValueNotString},
		{name: "ends missing value", raw: Op{CodeEnds, "/x"}, wantErr: ErrEndsMissingValue},
		{name: "ends value not string", raw: Op{CodeEnds, "/x", 1}, wantErr: ErrEndsValueNotString},
		{name: "type missing type", raw: Op{CodeType, "/x"}, wantErr: ErrTypeMissingType},
		{name: "type not string", raw: Op{CodeType, "/x", 1}, wantErr: ErrTypeNotString},
		{name: "test_type missing types", raw: Op{CodeTestType, "/x"}, wantErr: ErrTestTypeMissingTypes},
		{name: "test_type types not array", raw: Op{CodeTestType, "/x", "string"}, wantErr: ErrTestTypeTypesNotArray},
		{name: "test_type item not string", raw: Op{CodeTestType, "/x", []any{"string", 1}}, wantErr: ErrTestTypeTypesNotArray},
		{name: "test_string missing str", raw: Op{CodeTestString, "/x"}, wantErr: ErrTestStringMissingStr},
		{name: "test_string not string", raw: Op{CodeTestString, "/x", 1}, wantErr: ErrTestStringNotString},
		{name: "test_string_len missing len", raw: Op{CodeTestStringLen, "/x"}, wantErr: ErrTestStringLenMissingLen},
		{name: "test_string_len not number", raw: Op{CodeTestStringLen, "/x", "3"}, wantErr: ErrTestStringLenNotNumber},
		{name: "in missing values", raw: Op{CodeIn, "/x"}, wantErr: ErrInMissingValues},
		{name: "in values not array", raw: Op{CodeIn, "/x", "admin"}, wantErr: ErrInValuesNotArray},
		{name: "less missing value", raw: Op{CodeLess, "/x"}, wantErr: ErrLessMissingValue},
		{name: "less value not number", raw: Op{CodeLess, "/x", "1"}, wantErr: ErrLessValueNotNumber},
		{name: "more missing value", raw: Op{CodeMore, "/x"}, wantErr: ErrMoreMissingValue},
		{name: "more value not number", raw: Op{CodeMore, "/x", "1"}, wantErr: ErrMoreValueNotNumber},
		{name: "matches missing pattern", raw: Op{CodeMatches, "/x"}, wantErr: ErrMatchesMissingPattern},
		{name: "matches pattern not string", raw: Op{CodeMatches, "/x", 1}, wantErr: ErrMatchesPatternNotString},
		{name: "and missing ops", raw: Op{CodeAnd, "/x"}, wantErr: ErrAndMissingOps},
		{name: "or missing ops", raw: Op{CodeOr, "/x"}, wantErr: ErrOrMissingOps},
		{name: "not missing ops", raw: Op{CodeNot, "/x"}, wantErr: ErrNotMissingOps},
		{name: "predicate ops not array", raw: Op{CodeAnd, "/x", "ops"}, wantErr: ErrPredicateNotArray},
		{name: "predicate op invalid", raw: Op{CodeAnd, "/x", []any{"ops"}}, wantErr: ErrPredicateOpInvalid},
		{name: "predicate op must be predicate", raw: Op{CodeAnd, "/x", []any{[]any{CodeAdd, "/name", "Ada"}}}, wantErr: ErrNotPredicate},
	}

	decoder := NewDecoder()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := decoder.Decode(tc.raw)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestDecodeJSONRejectsInvalidPayload(t *testing.T) {
	t.Parallel()

	decoded, err := DecodeJSON([]byte(`{"op":"add"}`))
	require.Error(t, err)
	assert.Nil(t, decoded)
}
