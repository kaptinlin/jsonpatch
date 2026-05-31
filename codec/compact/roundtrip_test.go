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
		{name: "add", raw: Op{CodeAdd, []string{"profile", "name"}, "Ada"}, want: internal.Operation{Op: "add", Path: "/profile/name", Value: "Ada"}},
		{name: "remove", raw: Op{CodeRemove, []string{"profile", "name"}}, want: internal.Operation{Op: "remove", Path: "/profile/name"}},
		{name: "remove with old value", raw: Op{CodeRemove, []string{"profile", "name"}, "Ada"}, want: internal.Operation{Op: "remove", Path: "/profile/name", OldValue: "Ada"}},
		{name: "replace", raw: Op{CodeReplace, []string{"profile", "name"}, "Grace"}, want: internal.Operation{Op: "replace", Path: "/profile/name", Value: "Grace"}},
		{name: "replace with old value", raw: Op{CodeReplace, []string{"profile", "name"}, "Grace", "Ada"}, want: internal.Operation{Op: "replace", Path: "/profile/name", Value: "Grace", OldValue: "Ada"}},
		{name: "move", raw: Op{CodeMove, []string{"profile", "displayName"}, []string{"profile", "name"}}, want: internal.Operation{Op: "move", Path: "/profile/displayName", From: "/profile/name"}},
		{name: "copy", raw: Op{CodeCopy, []string{"profile", "alias"}, []string{"profile", "name"}}, want: internal.Operation{Op: "copy", Path: "/profile/alias", From: "/profile/name"}},
		{name: "test", raw: Op{CodeTest, []string{"profile", "name"}, "Ada", true}, want: internal.Operation{Op: "test", Path: "/profile/name", Value: "Ada", Not: true}},
		{name: "flip", raw: Op{CodeFlip, []string{"enabled"}}, want: internal.Operation{Op: "flip", Path: "/enabled"}},
		{name: "inc", raw: Op{CodeInc, []string{"count"}, 2}, want: internal.Operation{Op: "inc", Path: "/count", Inc: 2}},
		{name: "str_ins", raw: Op{CodeStrIns, []string{"profile", "name"}, 1, "d"}, want: internal.Operation{Op: "str_ins", Path: "/profile/name", Pos: 1, Str: "d"}},
		{name: "str_del string mode", raw: Op{CodeStrDel, []string{"profile", "name"}, 1, "da"}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Str: "da"}},
		{name: "str_del length mode", raw: Op{CodeStrDel, []string{"profile", "name"}, 1, 0, 2}, want: internal.Operation{Op: "str_del", Path: "/profile/name", Pos: 1, Len: 2}},
		{name: "split", raw: Op{CodeSplit, []string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}}, want: internal.Operation{Op: "split", Path: "/nodes/0", Pos: 1, Props: map[string]any{"kind": "paragraph"}}},
		{name: "merge", raw: Op{CodeMerge, []string{"nodes", "1"}, 1, map[string]any{"merged": true}}, want: internal.Operation{Op: "merge", Path: "/nodes/1", Pos: 1, Props: map[string]any{"merged": true}}},
		{name: "extend", raw: Op{CodeExtend, []string{"profile"}, map[string]any{"name": "Ada"}, true}, want: internal.Operation{Op: "extend", Path: "/profile", Props: map[string]any{"name": "Ada"}, DeleteNull: true}},
		{name: "defined", raw: Op{CodeDefined, []string{"profile", "name"}}, want: internal.Operation{Op: "defined", Path: "/profile/name"}},
		{name: "undefined", raw: Op{CodeUndefined, []string{"profile", "deleted"}}, want: internal.Operation{Op: "undefined", Path: "/profile/deleted"}},
		{name: "contains", raw: Op{CodeContains, []string{"profile", "name"}, "Ad", true}, want: internal.Operation{Op: "contains", Path: "/profile/name", Value: "Ad", IgnoreCase: true}},
		{name: "starts", raw: Op{CodeStarts, []string{"profile", "name"}, "A", true}, want: internal.Operation{Op: "starts", Path: "/profile/name", Value: "A", IgnoreCase: true}},
		{name: "ends", raw: Op{CodeEnds, []string{"profile", "name"}, "a", true}, want: internal.Operation{Op: "ends", Path: "/profile/name", Value: "a", IgnoreCase: true}},
		{name: "type", raw: Op{CodeType, []string{"profile", "name"}, "string"}, want: internal.Operation{Op: "type", Path: "/profile/name", Value: "string"}},
		{name: "test_type", raw: Op{CodeTestType, []string{"profile", "name"}, []any{"string", "null"}}, want: internal.Operation{Op: "test_type", Path: "/profile/name", Type: []string{"string", "null"}}},
		{name: "test_string", raw: Op{CodeTestString, []string{"profile", "name"}, 1, "da", true}, want: internal.Operation{Op: "test_string", Path: "/profile/name", Str: "da", Pos: 1, Not: true}},
		{name: "test_string_len", raw: Op{CodeTestStringLen, []string{"profile", "name"}, 3, true}, want: internal.Operation{Op: "test_string_len", Path: "/profile/name", Len: 3, Not: true}},
		{name: "in", raw: Op{CodeIn, []string{"role"}, []any{"admin", "editor"}}, want: internal.Operation{Op: "in", Path: "/role", Value: []any{"admin", "editor"}}},
		{name: "less", raw: Op{CodeLess, []string{"score"}, 10}, want: internal.Operation{Op: "less", Path: "/score", Value: 10}},
		{name: "more", raw: Op{CodeMore, []string{"score"}, 5}, want: internal.Operation{Op: "more", Path: "/score", Value: 5}},
		{name: "matches", raw: Op{CodeMatches, []string{"profile", "name"}, "^ad", true}, want: internal.Operation{Op: "matches", Path: "/profile/name", Value: "^ad", IgnoreCase: true}},
		{
			name: "and",
			raw:  Op{CodeAnd, []string{"profile"}, []any{[]any{CodeDefined, []string{"name"}}, []any{CodeContains, []string{"role"}, "admin"}}},
			want: internal.Operation{Op: "and", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile/name"},
				{Op: "contains", Path: "/profile/role", Value: "admin"},
			}},
		},
		{
			name: "or",
			raw:  Op{CodeOr, []string{"profile"}, []any{[]any{CodeDefined, []string{"name"}}, []any{CodeDefined, []string{"email"}}}},
			want: internal.Operation{Op: "or", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile/name"},
				{Op: "defined", Path: "/profile/email"},
			}},
		},
		{
			name: "not",
			raw:  Op{CodeNot, []string{"profile"}, []any{[]any{CodeUndefined, []string{"deleted"}}}},
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

	want := Op{"move", []string{"with/slash", "~tilde"}, []string{""}}
	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("encoded operation mismatch (-want +got):\n%s", diff)
	}
}

func TestEncodeCanonicalCompactForm(t *testing.T) {
	t.Parallel()

	encoded, err := Encode([]internal.Op{
		op.NewTest([]string{"profile", "name"}, "Ada"),
		op.NewTestWithNot([]string{"profile", "name"}, "Grace", true),
		op.NewContainsWithIgnoreCase([]string{"profile", "name"}, "ada", true),
		op.NewMatches([]string{"profile", "name"}, "^A", false, nil),
		op.NewTestString([]string{"profile", "name"}, "da", 1, true, false),
		op.NewAnd([]string{"profile"}, []any{
			op.NewDefined([]string{"profile", "name"}),
			op.NewNotMultiple([]string{"profile", "meta"}, []any{
				op.NewUndefined([]string{"profile", "meta", "deleted"}),
			}),
		}),
	})
	require.NoError(t, err)

	want := []Op{
		{internal.OpTestCode, []string{"profile", "name"}, "Ada"},
		{internal.OpTestCode, []string{"profile", "name"}, "Grace", 1},
		{internal.OpContainsCode, []string{"profile", "name"}, "ada", true},
		{internal.OpMatchesCode, []string{"profile", "name"}, "^A"},
		{internal.OpTestStringCode, []string{"profile", "name"}, 1, "da", true},
		{internal.OpAndCode, []string{"profile"}, []any{
			[]any{internal.OpDefinedCode, []string{"name"}},
			[]any{internal.OpNotCode, []string{"meta"}, []any{
				[]any{internal.OpUndefinedCode, []string{"deleted"}},
			}},
		}},
	}
	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("encoded compact form mismatch (-want +got):\n%s", diff)
	}
}

func TestEncodeCanonicalOptionalFieldGolden(t *testing.T) {
	t.Parallel()

	encoded, err := Encode([]internal.Op{
		op.NewTest([]string{"flag"}, true),
		op.NewTestWithNot([]string{"flag"}, false, true),
		op.NewContains([]string{"name"}, "Ada"),
		op.NewContainsWithIgnoreCase([]string{"name"}, "ada", true),
		op.NewStarts([]string{"name"}, "A"),
		op.NewStartsWithIgnoreCase([]string{"name"}, "a", true),
		op.NewEnds([]string{"name"}, "a"),
		op.NewEndsWithIgnoreCase([]string{"name"}, "A", true),
		op.NewMatches([]string{"name"}, "^A", false, nil),
		op.NewMatches([]string{"name"}, "^a", true, nil),
		op.NewTestString([]string{"name"}, "Ad", 0, false, false),
		op.NewTestString([]string{"name"}, "ad", 0, true, false),
		op.NewTestStringLen([]string{"name"}, 3),
		op.NewTestStringLenWithNot([]string{"name"}, 4, true),
		op.NewSplit([]string{"nodes", "0"}, 1, nil),
		op.NewSplit([]string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}),
		op.NewMerge([]string{"nodes", "1"}, 1, nil),
		op.NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}),
		op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, false),
		op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, true),
	})
	require.NoError(t, err)

	want := []Op{
		{internal.OpTestCode, []string{"flag"}, true},
		{internal.OpTestCode, []string{"flag"}, false, 1},
		{internal.OpContainsCode, []string{"name"}, "Ada"},
		{internal.OpContainsCode, []string{"name"}, "ada", true},
		{internal.OpStartsCode, []string{"name"}, "A"},
		{internal.OpStartsCode, []string{"name"}, "a", true},
		{internal.OpEndsCode, []string{"name"}, "a"},
		{internal.OpEndsCode, []string{"name"}, "A", true},
		{internal.OpMatchesCode, []string{"name"}, "^A"},
		{internal.OpMatchesCode, []string{"name"}, "^a", true},
		{internal.OpTestStringCode, []string{"name"}, 0, "Ad"},
		{internal.OpTestStringCode, []string{"name"}, 0, "ad", true},
		{internal.OpTestStringLenCode, []string{"name"}, 3.0},
		{internal.OpTestStringLenCode, []string{"name"}, 4.0, true},
		{internal.OpSplitCode, []string{"nodes", "0"}, 1.0},
		{internal.OpSplitCode, []string{"nodes", "0"}, 1.0, map[string]any{"kind": "paragraph"}},
		{internal.OpMergeCode, []string{"nodes", "1"}, 1.0},
		{internal.OpMergeCode, []string{"nodes", "1"}, 1.0, map[string]any{"merged": true}},
		{internal.OpExtendCode, []string{"profile"}, map[string]any{"name": "Ada"}},
		{internal.OpExtendCode, []string{"profile"}, map[string]any{"name": "Ada"}, true},
	}
	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("canonical compact optional fields mismatch (-want +got):\n%s", diff)
	}
}

func TestEncodeCanonicalCompositeRelativePathGolden(t *testing.T) {
	t.Parallel()

	encoded, err := Encode([]internal.Op{
		op.NewAnd([]string{"profile"}, []any{
			op.NewDefined([]string{"profile", "name"}),
			op.NewOr([]string{"profile", "flags"}, []any{
				op.NewUndefined([]string{"profile", "flags", "deleted"}),
				op.NewContains([]string{"profile", "flags", "role"}, "admin"),
			}),
		}),
	})
	require.NoError(t, err)

	want := []Op{
		{internal.OpAndCode, []string{"profile"}, []any{
			[]any{internal.OpDefinedCode, []string{"name"}},
			[]any{internal.OpOrCode, []string{"flags"}, []any{
				[]any{internal.OpUndefinedCode, []string{"deleted"}},
				[]any{internal.OpContainsCode, []string{"role"}, "admin"},
			}},
		}},
	}
	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("canonical compact composite paths mismatch (-want +got):\n%s", diff)
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
		{name: "unknown string opcode", raw: Op{"unknown", []string{"x"}}, wantErr: ErrUnknownStringCode},
		{name: "unknown numeric opcode", raw: Op{999, []string{"x"}}, wantErr: ErrUnknownNumericCode},
		{name: "unknown float opcode", raw: Op{999.0, []string{"x"}}, wantErr: ErrUnknownNumericCode},
		{name: "invalid opcode type", raw: Op{true, []string{"x"}}, wantErr: ErrInvalidCodeType},
		{name: "add missing value", raw: Op{CodeAdd, []string{"x"}}, wantErr: ErrAddMissingValue},
		{name: "replace missing value", raw: Op{CodeReplace, []string{"x"}}, wantErr: ErrReplaceMissingValue},
		{name: "move missing from", raw: Op{CodeMove, []string{"x"}}, wantErr: ErrMoveMissingFrom},
		{name: "move from not string", raw: Op{CodeMove, []string{"x"}, 1}, wantErr: ErrMoveFromNotString},
		{name: "copy missing from", raw: Op{CodeCopy, []string{"x"}}, wantErr: ErrCopyMissingFrom},
		{name: "copy from not string", raw: Op{CodeCopy, []string{"x"}, 1}, wantErr: ErrCopyFromNotString},
		{name: "test missing value", raw: Op{CodeTest, []string{"x"}}, wantErr: ErrTestMissingValue},
		{name: "inc missing delta", raw: Op{CodeInc, []string{"x"}}, wantErr: ErrIncMissingDelta},
		{name: "inc delta not number", raw: Op{CodeInc, []string{"x"}, "1"}, wantErr: ErrIncDeltaNotNumber},
		{name: "str_ins missing fields", raw: Op{CodeStrIns, []string{"x"}}, wantErr: ErrStrInsMissingFields},
		{name: "str_ins pos not number", raw: Op{CodeStrIns, []string{"x"}, "1", "a"}, wantErr: ErrStrInsPosNotNumber},
		{name: "str_ins str not string", raw: Op{CodeStrIns, []string{"x"}, 1, false}, wantErr: ErrStrInsStrNotString},
		{name: "str_del missing fields", raw: Op{CodeStrDel, []string{"x"}}, wantErr: ErrStrDelMissingFields},
		{name: "str_del pos not number", raw: Op{CodeStrDel, []string{"x"}, "1", 2}, wantErr: ErrStrDelPosNotNumber},
		{name: "str_del len not number", raw: Op{CodeStrDel, []string{"x"}, 1, false}, wantErr: ErrStrDelLenNotNumber},
		{name: "str_del explicit len not number", raw: Op{CodeStrDel, []string{"x"}, 1, 0, false}, wantErr: ErrStrDelLenNotNumber},
		{name: "split missing pos", raw: Op{CodeSplit, []string{"x"}}, wantErr: ErrSplitMissingPos},
		{name: "split pos not number", raw: Op{CodeSplit, []string{"x"}, "1"}, wantErr: ErrSplitPosNotNumber},
		{name: "merge missing pos", raw: Op{CodeMerge, []string{"x"}}, wantErr: ErrMergeMissingPos},
		{name: "merge pos not number", raw: Op{CodeMerge, []string{"x"}, "1"}, wantErr: ErrMergePosNotNumber},
		{name: "extend missing props", raw: Op{CodeExtend, []string{"x"}}, wantErr: ErrExtendMissingProps},
		{name: "extend props not object", raw: Op{CodeExtend, []string{"x"}, "props"}, wantErr: ErrExtendPropsNotObject},
		{name: "contains missing value", raw: Op{CodeContains, []string{"x"}}, wantErr: ErrContainsMissingValue},
		{name: "contains value not string", raw: Op{CodeContains, []string{"x"}, 1}, wantErr: ErrContainsValueNotString},
		{name: "starts missing value", raw: Op{CodeStarts, []string{"x"}}, wantErr: ErrStartsMissingValue},
		{name: "starts value not string", raw: Op{CodeStarts, []string{"x"}, 1}, wantErr: ErrStartsValueNotString},
		{name: "ends missing value", raw: Op{CodeEnds, []string{"x"}}, wantErr: ErrEndsMissingValue},
		{name: "ends value not string", raw: Op{CodeEnds, []string{"x"}, 1}, wantErr: ErrEndsValueNotString},
		{name: "type missing type", raw: Op{CodeType, []string{"x"}}, wantErr: ErrTypeMissingType},
		{name: "type not string", raw: Op{CodeType, []string{"x"}, 1}, wantErr: ErrTypeNotString},
		{name: "test_type missing types", raw: Op{CodeTestType, []string{"x"}}, wantErr: ErrTestTypeMissingTypes},
		{name: "test_type types not array", raw: Op{CodeTestType, []string{"x"}, "string"}, wantErr: ErrTestTypeTypesNotArray},
		{name: "test_type item not string", raw: Op{CodeTestType, []string{"x"}, []any{"string", 1}}, wantErr: ErrTestTypeTypesNotArray},
		{name: "test_string missing pos", raw: Op{CodeTestString, []string{"x"}}, wantErr: ErrTestStringMissingPos},
		{name: "test_string pos not number", raw: Op{CodeTestString, []string{"x"}, "1", "a"}, wantErr: ErrTestStringPosNotNumber},
		{name: "test_string missing str", raw: Op{CodeTestString, []string{"x"}, 1}, wantErr: ErrTestStringMissingStr},
		{name: "test_string not string", raw: Op{CodeTestString, []string{"x"}, 1, 1}, wantErr: ErrTestStringNotString},
		{name: "test_string_len missing len", raw: Op{CodeTestStringLen, []string{"x"}}, wantErr: ErrTestStringLenMissingLen},
		{name: "test_string_len not number", raw: Op{CodeTestStringLen, []string{"x"}, "3"}, wantErr: ErrTestStringLenNotNumber},
		{name: "in missing values", raw: Op{CodeIn, []string{"x"}}, wantErr: ErrInMissingValues},
		{name: "in values not array", raw: Op{CodeIn, []string{"x"}, "admin"}, wantErr: ErrInValuesNotArray},
		{name: "less missing value", raw: Op{CodeLess, []string{"x"}}, wantErr: ErrLessMissingValue},
		{name: "less value not number", raw: Op{CodeLess, []string{"x"}, "1"}, wantErr: ErrLessValueNotNumber},
		{name: "more missing value", raw: Op{CodeMore, []string{"x"}}, wantErr: ErrMoreMissingValue},
		{name: "more value not number", raw: Op{CodeMore, []string{"x"}, "1"}, wantErr: ErrMoreValueNotNumber},
		{name: "matches missing pattern", raw: Op{CodeMatches, []string{"x"}}, wantErr: ErrMatchesMissingPattern},
		{name: "matches pattern not string", raw: Op{CodeMatches, []string{"x"}, 1}, wantErr: ErrMatchesPatternNotString},
		{name: "and missing ops", raw: Op{CodeAnd, []string{"x"}}, wantErr: ErrAndMissingOps},
		{name: "or missing ops", raw: Op{CodeOr, []string{"x"}}, wantErr: ErrOrMissingOps},
		{name: "not missing ops", raw: Op{CodeNot, []string{"x"}}, wantErr: ErrNotMissingOps},
		{name: "predicate ops not array", raw: Op{CodeAnd, []string{"x"}, "ops"}, wantErr: ErrPredicateNotArray},
		{name: "predicate op invalid", raw: Op{CodeAnd, []string{"x"}, []any{"ops"}}, wantErr: ErrPredicateOpInvalid},
		{name: "predicate op must be predicate", raw: Op{CodeAnd, []string{"x"}, []any{[]any{CodeAdd, []string{"name"}, "Ada"}}}, wantErr: ErrNotPredicate},
		{name: "not requires one predicate", raw: Op{CodeNot, []string{"x"}, []any{[]any{CodeDefined, []string{"a"}}, []any{CodeDefined, []string{"b"}}}}, wantErr: ErrNotSinglePredicate},
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
