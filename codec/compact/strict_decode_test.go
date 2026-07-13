package compact

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeRejectsMalformedCoreOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     Op
		wantErr error
	}{
		{name: "fractional opcode", raw: Op{float64(CodeAdd) + 0.5, []string{"x"}, 1}, wantErr: ErrInvalidCodeType},
		{name: "NaN opcode", raw: Op{math.NaN(), []string{"x"}, 1}, wantErr: ErrInvalidCodeType},
		{name: "infinite opcode", raw: Op{math.Inf(1), []string{"x"}, 1}, wantErr: ErrInvalidCodeType},
		{name: "unknown opcode", raw: Op{999, []string{"x"}}, wantErr: ErrUnknownNumericCode},
		{name: "add missing value", raw: Op{CodeAdd, []string{"x"}}, wantErr: ErrAddMissingValue},
		{name: "add extra field", raw: Op{CodeAdd, []string{"x"}, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "remove extra field", raw: Op{CodeRemove, []string{"x"}, nil, "extra"}, wantErr: ErrInvalidArity},
		{name: "replace missing value", raw: Op{CodeReplace, []string{"x"}}, wantErr: ErrReplaceMissingValue},
		{name: "replace extra field", raw: Op{CodeReplace, []string{"x"}, 1, nil, "extra"}, wantErr: ErrInvalidArity},
		{name: "move missing from", raw: Op{CodeMove, []string{"x"}}, wantErr: ErrMoveMissingFrom},
		{name: "move wrong from type", raw: Op{CodeMove, []string{"x"}, 1}, wantErr: ErrMoveFromNotString},
		{name: "move extra field", raw: Op{CodeMove, []string{"x"}, []string{"y"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "copy missing from", raw: Op{CodeCopy, []string{"x"}}, wantErr: ErrCopyMissingFrom},
		{name: "copy wrong from type", raw: Op{CodeCopy, []string{"x"}, 1}, wantErr: ErrCopyFromNotString},
		{name: "copy extra field", raw: Op{CodeCopy, []string{"x"}, []string{"y"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "test missing value", raw: Op{CodeTest, []string{"x"}}, wantErr: ErrTestMissingValue},
		{name: "test extra field", raw: Op{CodeTest, []string{"x"}, 1, true, "extra"}, wantErr: ErrInvalidArity},
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

func TestDecodeRejectsMalformedPredicateOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     Op
		wantErr error
	}{
		{name: "test boolean marker", raw: Op{CodeTest, []string{"x"}, 1, true}, wantErr: ErrInvalidTrueMarker},
		{name: "test zero marker", raw: Op{CodeTest, []string{"x"}, 1, 0}, wantErr: ErrInvalidTrueMarker},
		{name: "test non-one marker", raw: Op{CodeTest, []string{"x"}, 1, 2}, wantErr: ErrInvalidTrueMarker},
		{name: "test string marker", raw: Op{CodeTest, []string{"x"}, 1, "1"}, wantErr: ErrInvalidTrueMarker},
		{name: "defined extra field", raw: Op{CodeDefined, []string{"x"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "undefined extra field", raw: Op{CodeUndefined, []string{"x"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "contains extra field", raw: Op{CodeContains, []string{"x"}, "a", 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "contains boolean marker", raw: Op{CodeContains, []string{"x"}, "a", true}, wantErr: ErrInvalidTrueMarker},
		{name: "starts extra field", raw: Op{CodeStarts, []string{"x"}, "a", 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "ends extra field", raw: Op{CodeEnds, []string{"x"}, "a", 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "type extra field", raw: Op{CodeType, []string{"x"}, "string", "extra"}, wantErr: ErrInvalidArity},
		{name: "test_type extra field", raw: Op{CodeTestType, []string{"x"}, []any{"string"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "test_string extra field", raw: Op{CodeTestString, []string{"x"}, 0, "a", 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "test_string boolean marker", raw: Op{CodeTestString, []string{"x"}, 0, "a", true}, wantErr: ErrInvalidTrueMarker},
		{name: "test_string fractional pos", raw: Op{CodeTestString, []string{"x"}, 0.5, "a"}, wantErr: ErrNumberNotInteger},
		{name: "test_string non-finite pos", raw: Op{CodeTestString, []string{"x"}, math.NaN(), "a"}, wantErr: ErrNumberNotFinite},
		{name: "test_string_len extra field", raw: Op{CodeTestStringLen, []string{"x"}, 1, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "test_string_len fractional length", raw: Op{CodeTestStringLen, []string{"x"}, 1.5}, wantErr: ErrNumberNotInteger},
		{name: "test_string_len non-finite length", raw: Op{CodeTestStringLen, []string{"x"}, math.Inf(1)}, wantErr: ErrNumberNotFinite},
		{name: "in extra field", raw: Op{CodeIn, []string{"x"}, []any{1}, "extra"}, wantErr: ErrInvalidArity},
		{name: "less extra field", raw: Op{CodeLess, []string{"x"}, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "more extra field", raw: Op{CodeMore, []string{"x"}, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "matches extra field", raw: Op{CodeMatches, []string{"x"}, "a", 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "and extra field", raw: Op{CodeAnd, []string{"x"}, []any{[]any{CodeDefined, []string{"a"}}}, "extra"}, wantErr: ErrInvalidArity},
		{name: "or extra field", raw: Op{CodeOr, []string{"x"}, []any{[]any{CodeDefined, []string{"a"}}}, "extra"}, wantErr: ErrInvalidArity},
		{name: "not extra field", raw: Op{CodeNot, []string{"x"}, []any{[]any{CodeDefined, []string{"a"}}}, "extra"}, wantErr: ErrInvalidArity},
		{name: "nested malformed arity", raw: Op{CodeAnd, []string{"x"}, []any{[]any{CodeDefined, []string{"a"}, "extra"}}}, wantErr: ErrInvalidArity},
		{name: "nested mutation", raw: Op{CodeAnd, []string{"x"}, []any{[]any{CodeAdd, []string{"a"}, 1}}}, wantErr: ErrNotPredicate},
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

	for _, marker := range []any{1, int64(1), float64(1)} {
		decoded, err := decoder.Decode(Op{CodeContains, []string{"x"}, "a", marker})
		require.NoError(t, err)
		assert.NotNil(t, decoded)
	}
}

func TestDecodeRejectsMalformedExtendedOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     Op
		wantErr error
	}{
		{name: "flip extra field", raw: Op{CodeFlip, []string{"x"}, "extra"}, wantErr: ErrInvalidArity},
		{name: "inc extra field", raw: Op{CodeInc, []string{"x"}, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "inc NaN", raw: Op{CodeInc, []string{"x"}, math.NaN()}, wantErr: ErrNumberNotFinite},
		{name: "inc infinity", raw: Op{CodeInc, []string{"x"}, math.Inf(1)}, wantErr: ErrNumberNotFinite},
		{name: "str_ins extra field", raw: Op{CodeStrIns, []string{"x"}, 1, "a", "extra"}, wantErr: ErrInvalidArity},
		{name: "str_ins fractional pos", raw: Op{CodeStrIns, []string{"x"}, 1.5, "a"}, wantErr: ErrNumberNotInteger},
		{name: "str_ins infinite pos", raw: Op{CodeStrIns, []string{"x"}, math.Inf(1), "a"}, wantErr: ErrNumberNotFinite},
		{name: "str_del string extra field", raw: Op{CodeStrDel, []string{"x"}, 1, "a", "extra"}, wantErr: ErrInvalidArity},
		{name: "str_del short numeric variant", raw: Op{CodeStrDel, []string{"x"}, 1, 2}, wantErr: ErrInvalidVariant},
		{name: "str_del wrong length marker", raw: Op{CodeStrDel, []string{"x"}, 1, 1, 2}, wantErr: ErrInvalidVariant},
		{name: "str_del length extra field", raw: Op{CodeStrDel, []string{"x"}, 1, 0, 2, "extra"}, wantErr: ErrInvalidArity},
		{name: "str_del fractional pos", raw: Op{CodeStrDel, []string{"x"}, 1.5, "a"}, wantErr: ErrNumberNotInteger},
		{name: "str_del fractional length", raw: Op{CodeStrDel, []string{"x"}, 1, 0, 1.5}, wantErr: ErrNumberNotInteger},
		{name: "split extra field", raw: Op{CodeSplit, []string{"x"}, 1, map[string]any{}, "extra"}, wantErr: ErrInvalidArity},
		{name: "split non-finite pos", raw: Op{CodeSplit, []string{"x"}, math.NaN()}, wantErr: ErrNumberNotFinite},
		{name: "split props wrong type", raw: Op{CodeSplit, []string{"x"}, 1, "props"}, wantErr: ErrSplitPropsNotObject},
		{name: "merge extra field", raw: Op{CodeMerge, []string{"x"}, 1, map[string]any{}, "extra"}, wantErr: ErrInvalidArity},
		{name: "merge non-finite pos", raw: Op{CodeMerge, []string{"x"}, math.Inf(-1)}, wantErr: ErrNumberNotFinite},
		{name: "merge props wrong type", raw: Op{CodeMerge, []string{"x"}, 1, "props"}, wantErr: ErrMergePropsNotObject},
		{name: "extend extra field", raw: Op{CodeExtend, []string{"x"}, map[string]any{}, 1, "extra"}, wantErr: ErrInvalidArity},
		{name: "extend boolean marker", raw: Op{CodeExtend, []string{"x"}, map[string]any{}, true}, wantErr: ErrInvalidTrueMarker},
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

	for _, raw := range []Op{
		{CodeSplit, []string{"x"}, 1.5},
		{CodeMerge, []string{"x"}, -1.5},
	} {
		decoded, err := decoder.Decode(raw)
		require.NoError(t, err)
		assert.NotNil(t, decoded)
	}
}

func FuzzDecodeMalformedCompactDoesNotPanic(f *testing.F) {
	seeds := [][]byte{
		[]byte(`[[0.5,[],1]]`),
		[]byte(`[[30,[],"a",2]]`),
		[]byte(`[[43,[],[[31,[],"extra"]]]]`),
		[]byte(`[[0,[],1,2,3,4,5,6,7,8,9]]`),
		[]byte(`[[7,[],1,0,2]]`),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(_ *testing.T, data []byte) {
		_, _ = DecodeJSON(data)
	})
}
