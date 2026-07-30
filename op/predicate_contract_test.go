package op

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

type codecOperation interface {
	internal.JSONOp
	internal.CompactOp
}

type predicateCodecOperation interface {
	internal.PredicateOp
	internal.JSONOp
	internal.CompactOp
}

func TestOperationSerializationContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		op          codecOperation
		wantType    internal.OpType
		wantCode    int
		wantJSON    internal.Operation
		wantCompact internal.CompactOperation
	}{
		{
			name:        "remove with old value",
			op:          NewRemoveWithOldValue([]string{"name"}, "Ada"),
			wantType:    internal.OpRemoveType,
			wantCode:    internal.OpRemoveCode,
			wantJSON:    internal.Operation{Op: "remove", Path: "/name", OldValue: "Ada", HasOldValue: true},
			wantCompact: internal.CompactOperation{internal.OpRemoveCode, []string{"name"}, "Ada"},
		},
		{
			name:        "copy",
			op:          NewCopy([]string{"profile", "alias"}, []string{"profile", "name"}),
			wantType:    internal.OpCopyType,
			wantCode:    internal.OpCopyCode,
			wantJSON:    internal.Operation{Op: "copy", Path: "/profile/alias", From: "/profile/name"},
			wantCompact: internal.CompactOperation{internal.OpCopyCode, []string{"profile", "alias"}, []string{"profile", "name"}},
		},
		{
			name:        "flip",
			op:          NewFlip([]string{"enabled"}),
			wantType:    internal.OpFlipType,
			wantCode:    internal.OpFlipCode,
			wantJSON:    internal.Operation{Op: "flip", Path: "/enabled"},
			wantCompact: internal.CompactOperation{internal.OpFlipCode, []string{"enabled"}},
		},
		{
			name:        "inc",
			op:          NewInc([]string{"count"}, 2),
			wantType:    internal.OpIncType,
			wantCode:    internal.OpIncCode,
			wantJSON:    internal.Operation{Op: "inc", Path: "/count", Inc: 2},
			wantCompact: internal.CompactOperation{internal.OpIncCode, []string{"count"}, 2.0},
		},
		{
			name:        "str_ins",
			op:          NewStrIns([]string{"name"}, 1, "d"),
			wantType:    internal.OpStrInsType,
			wantCode:    internal.OpStrInsCode,
			wantJSON:    internal.Operation{Op: "str_ins", Path: "/name", Pos: 1, Str: "d"},
			wantCompact: internal.CompactOperation{internal.OpStrInsCode, []string{"name"}, 1, "d"},
		},
		{
			name:        "str_del length mode",
			op:          NewStrDel([]string{"name"}, 1, 2),
			wantType:    internal.OpStrDelType,
			wantCode:    internal.OpStrDelCode,
			wantJSON:    internal.Operation{Op: "str_del", Path: "/name", Pos: 1, Len: 2},
			wantCompact: internal.CompactOperation{internal.OpStrDelCode, []string{"name"}, 1, 0, 2},
		},
		{
			name:        "str_del string mode",
			op:          NewStrDelWithStr([]string{"name"}, 1, "da"),
			wantType:    internal.OpStrDelType,
			wantCode:    internal.OpStrDelCode,
			wantJSON:    internal.Operation{Op: "str_del", Path: "/name", Pos: 1, Str: "da", HasStr: true},
			wantCompact: internal.CompactOperation{internal.OpStrDelCode, []string{"name"}, 1, "da"},
		},
		{
			name:        "split with props",
			op:          NewSplit([]string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}),
			wantType:    internal.OpSplitType,
			wantCode:    internal.OpSplitCode,
			wantJSON:    internal.Operation{Op: "split", Path: "/nodes/0", Pos: 1, Props: map[string]any{"kind": "paragraph"}},
			wantCompact: internal.CompactOperation{internal.OpSplitCode, []string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}},
		},
		{
			name:        "merge with props",
			op:          NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}),
			wantType:    internal.OpMergeType,
			wantCode:    internal.OpMergeCode,
			wantJSON:    internal.Operation{Op: "merge", Path: "/nodes/1", Pos: 1, Props: map[string]any{"merged": true}},
			wantCompact: internal.CompactOperation{internal.OpMergeCode, []string{"nodes", "1"}, 1, map[string]any{"merged": true}},
		},
		{
			name:        "extend with delete null",
			op:          NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, true),
			wantType:    internal.OpExtendType,
			wantCode:    internal.OpExtendCode,
			wantJSON:    internal.Operation{Op: "extend", Path: "/profile", Props: map[string]any{"name": "Ada"}, DeleteNull: true},
			wantCompact: internal.CompactOperation{internal.OpExtendCode, []string{"profile"}, map[string]any{"name": "Ada"}, 1},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantType, tc.op.Op())
			assert.Equal(t, tc.wantCode, tc.op.Code())

			gotJSON, err := tc.op.ToJSON()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.wantJSON, gotJSON); diff != "" {
				t.Errorf("ToJSON() mismatch (-want +got):\n%s", diff)
			}

			gotCompact, err := tc.op.ToCompact()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.wantCompact, gotCompact); diff != "" {
				t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
			}

			require.NoError(t, tc.op.Validate())
		})
	}
}

func TestOperationValidateRejectsInvalidOperands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      internal.Op
		wantErr error
	}{
		{name: "extend nil properties", op: NewExtend([]string{"profile"}, nil, false), wantErr: ErrPropertiesNil},
		{name: "str_del missing length", op: NewStrDel([]string{"name"}, 1, 0), wantErr: ErrMissingStrOrLen},
		{name: "test_string_len negative length", op: NewTestStringLen([]string{"name"}, -1), wantErr: ErrLengthNegative},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.op.Validate()
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestSlateNodeSplitAndMergeContracts(t *testing.T) {
	t.Parallel()

	t.Run("split text node preserves marks and applies props", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"node": map[string]any{"text": "Ada", "bold": true}}
		result, err := NewSplit([]string{"node"}, 1, map[string]any{"italic": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"node": []any{
			map[string]any{"text": "A", "bold": true, "italic": true},
			map[string]any{"text": "da", "bold": true, "italic": true},
		}}
		if diff := cmp.Diff(wantDoc, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
		assert.Equal(t, map[string]any{"text": "Ada", "bold": true}, result.Old)
	})

	t.Run("split element node splits children", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"node": map[string]any{
			"type": "paragraph",
			"children": []any{
				map[string]any{"text": "A"},
				map[string]any{"text": "B"},
			},
		}}
		result, err := NewSplit([]string{"node"}, 1, map[string]any{"split": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"node": []any{
			map[string]any{"type": "paragraph", "split": true, "children": []any{map[string]any{"text": "A"}}},
			map[string]any{"type": "paragraph", "split": true, "children": []any{map[string]any{"text": "B"}}},
		}}
		if diff := cmp.Diff(wantDoc, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("merge text nodes preserves later marks and ignores metadata props", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"nodes": []any{
			map[string]any{"text": "A", "bold": true},
			map[string]any{"text": "da", "italic": true},
		}}
		result, err := NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"nodes": []any{
			map[string]any{"text": "Ada", "bold": true, "italic": true},
		}}
		if diff := cmp.Diff(wantDoc, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("merge element nodes concatenates children", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"nodes": []any{
			map[string]any{"type": "paragraph", "children": []any{map[string]any{"text": "A"}}},
			map[string]any{"align": "center", "children": []any{map[string]any{"text": "B"}}},
		}}
		result, err := NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"nodes": []any{
			map[string]any{"type": "paragraph", "align": "center", "children": []any{map[string]any{"text": "A"}, map[string]any{"text": "B"}}},
		}}
		if diff := cmp.Diff(wantDoc, result.Doc); diff != "" {
			t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestPredicateOperationContracts(t *testing.T) {
	t.Parallel()

	doc := map[string]any{
		"name":    "Ada Lovelace",
		"role":    "admin",
		"score":   9.5,
		"age":     36.0,
		"active":  true,
		"tags":    []any{"math", "logic"},
		"profile": map[string]any{"name": "Ada"},
	}

	tests := []struct {
		name        string
		op          predicateCodecOperation
		wantJSON    internal.Operation
		wantCompact internal.CompactOperation
	}{
		{
			name:        "contains ignore case",
			op:          NewContainsWithIgnoreCase([]string{"name"}, "lovelace", true),
			wantJSON:    internal.Operation{Op: "contains", Path: "/name", Value: "lovelace", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpContainsCode, []string{"name"}, "lovelace", 1},
		},
		{
			name:        "starts ignore case",
			op:          NewStartsWithIgnoreCase([]string{"name"}, "ada", true),
			wantJSON:    internal.Operation{Op: "starts", Path: "/name", Value: "ada", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpStartsCode, []string{"name"}, "ada", 1},
		},
		{
			name:        "ends ignore case",
			op:          NewEndsWithIgnoreCase([]string{"name"}, "LOVELACE", true),
			wantJSON:    internal.Operation{Op: "ends", Path: "/name", Value: "LOVELACE", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpEndsCode, []string{"name"}, "LOVELACE", 1},
		},
		{
			name:        "in",
			op:          NewIn([]string{"role"}, []any{"admin", "editor"}),
			wantJSON:    internal.Operation{Op: "in", Path: "/role", Value: []any{"admin", "editor"}},
			wantCompact: internal.CompactOperation{internal.OpInCode, []string{"role"}, []any{"admin", "editor"}},
		},
		{
			name:        "less",
			op:          NewLess([]string{"score"}, 10),
			wantJSON:    internal.Operation{Op: "less", Path: "/score", Value: 10},
			wantCompact: internal.CompactOperation{internal.OpLessCode, []string{"score"}, 10.0},
		},
		{
			name:        "more",
			op:          NewMore([]string{"score"}, 5),
			wantJSON:    internal.Operation{Op: "more", Path: "/score", Value: 5},
			wantCompact: internal.CompactOperation{internal.OpMoreCode, []string{"score"}, 5.0},
		},
		{
			name:        "type integer",
			op:          NewType([]string{"age"}, "integer"),
			wantJSON:    internal.Operation{Op: "type", Path: "/age", Value: "integer"},
			wantCompact: internal.CompactOperation{internal.OpTypeCode, []string{"age"}, "integer"},
		},
		{
			name:        "matches ignore case",
			op:          NewMatches([]string{"name"}, `^ada\s+lovelace$`, true, nil),
			wantJSON:    internal.Operation{Op: "matches", Path: "/name", Value: `^ada\s+lovelace$`, IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpMatchesCode, []string{"name"}, `^ada\s+lovelace$`, 1},
		},
		{
			name:        "test type multiple",
			op:          NewTestTypeMultiple([]string{"name"}, []string{"string", "null"}),
			wantJSON:    internal.Operation{Op: "test_type", Path: "/name", Type: []string{"string", "null"}},
			wantCompact: internal.CompactOperation{internal.OpTestTypeCode, []string{"name"}, []string{"string", "null"}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := tc.op.Test(doc)
			require.NoError(t, err)
			assert.True(t, matched)

			result, err := tc.op.Apply(doc)
			require.NoError(t, err)
			if diff := cmp.Diff(doc, result.Doc); diff != "" {
				t.Errorf("Apply() document mismatch (-want +got):\n%s", diff)
			}

			gotJSON, err := tc.op.ToJSON()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.wantJSON, gotJSON); diff != "" {
				t.Errorf("ToJSON() mismatch (-want +got):\n%s", diff)
			}

			gotCompact, err := tc.op.ToCompact()
			require.NoError(t, err)
			if diff := cmp.Diff(tc.wantCompact, gotCompact); diff != "" {
				t.Errorf("ToCompact() mismatch (-want +got):\n%s", diff)
			}

			require.NoError(t, tc.op.Validate())
		})
	}
}

func TestPredicateOperationsReportFailures(t *testing.T) {
	t.Parallel()

	defaultDoc := map[string]any{
		"name":  "Ada",
		"role":  "viewer",
		"score": 5.0,
		"count": 5.0,
	}

	tests := []struct {
		name    string
		op      internal.PredicateOp
		doc     any
		wantErr error
	}{
		{name: "test mismatch", op: NewTest([]string{"name"}, "Grace"), doc: defaultDoc, wantErr: ErrTestOperationStringNotEquivalent},
		{name: "contains mismatch", op: NewContains([]string{"name"}, "Grace"), doc: defaultDoc, wantErr: ErrStringMismatch},
		{name: "defined missing", op: NewDefined([]string{"missing"}), doc: defaultDoc, wantErr: ErrDefinedTestFailed},
		{name: "undefined existing", op: NewUndefined([]string{"name"}), doc: defaultDoc, wantErr: ErrUndefinedTestFailed},
		{name: "starts mismatch", op: NewStarts([]string{"name"}, "Grace"), doc: defaultDoc, wantErr: ErrStringMismatch},
		{name: "ends mismatch", op: NewEnds([]string{"name"}, "Grace"), doc: defaultDoc, wantErr: ErrStringMismatch},
		{name: "in mismatch", op: NewIn([]string{"role"}, []any{"admin"}), doc: defaultDoc, wantErr: ErrOperationFailed},
		{name: "less mismatch", op: NewLess([]string{"score"}, 5), doc: defaultDoc, wantErr: ErrComparisonFailed},
		{name: "more mismatch", op: NewMore([]string{"score"}, 5), doc: defaultDoc, wantErr: ErrComparisonFailed},
		{name: "matches mismatch", op: NewMatches([]string{"name"}, `^Grace$`, false, nil), doc: defaultDoc, wantErr: ErrStringMismatch},
		{name: "type mismatch", op: NewType([]string{"name"}, "number"), doc: defaultDoc, wantErr: ErrTypeMismatch},
		{name: "test_type mismatch", op: NewTestType([]string{"name"}, "number"), doc: defaultDoc, wantErr: ErrTypeMismatch},
		{name: "test_string mismatch", op: NewTestString([]string{"name"}, "Gra", 0, false), doc: defaultDoc, wantErr: ErrSubstringMismatch},
		{name: "test_string_len mismatch", op: NewTestStringLen([]string{"name"}, 4), doc: defaultDoc, wantErr: ErrStringLengthMismatch},
		{name: "and mismatch", op: NewAnd(nil, NewDefined([]string{"missing"})), doc: defaultDoc, wantErr: ErrAndTestFailed},
		{name: "or mismatch", op: NewOr(nil, NewDefined([]string{"missing"})), doc: defaultDoc, wantErr: ErrOrTestFailed},
		{name: "not mismatch", op: NewNot(nil, NewDefined([]string{"name"})), doc: defaultDoc, wantErr: ErrNotTestFailed},
		{name: "missing path", op: NewContains([]string{"missing"}, "Ada"), doc: defaultDoc, wantErr: ErrPathNotFound},
		{name: "wrong target type", op: NewContains([]string{"count"}, "5"), doc: defaultDoc, wantErr: ErrNotString},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := tc.op.Test(tc.doc)
			require.NoError(t, err)
			assert.False(t, matched)

			result, err := tc.op.Apply(tc.doc)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrTestFailed)
			assert.ErrorIs(t, err, tc.wantErr)
			if diff := cmp.Diff(internal.OpResult[any]{}, result); diff != "" {
				t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
			}
		})
	}

	t.Run("matcher error is not a predicate mismatch", func(t *testing.T) {
		t.Parallel()

		cause := errors.New("matcher unavailable")
		operation := NewMatches(nil, "pattern", false, func(string, bool) (internal.RegexMatcher, error) {
			return nil, cause
		})
		matched, err := operation.Test("value")
		assert.False(t, matched)
		require.ErrorIs(t, err, cause)
		assert.NotErrorIs(t, err, ErrTestFailed)

		result, err := operation.Apply("value")
		assert.Equal(t, internal.OpResult[any]{}, result)
		require.ErrorIs(t, err, cause)
		assert.NotErrorIs(t, err, ErrTestFailed)
	})
}

func TestPredicateOperationValidateRejectsInvalidOperands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      internal.PredicateOp
		wantErr error
	}{
		{name: "in empty values", op: NewIn([]string{"role"}, nil), wantErr: ErrValuesArrayEmpty},
		{name: "matches empty pattern", op: NewMatches([]string{"name"}, "", false, nil), wantErr: ErrPatternEmpty},
		{name: "type invalid type", op: NewType([]string{"name"}, "invalid"), wantErr: ErrInvalidType},
		{name: "test_type empty list", op: NewTestTypeMultiple([]string{"name"}, nil), wantErr: ErrEmptyTypeList},
		{name: "test_type invalid type", op: NewTestType([]string{"name"}, "invalid"), wantErr: ErrInvalidType},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.op.Validate()
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCompositePredicatesRejectNilOperands(t *testing.T) {
	t.Parallel()

	var typedNil *TestOperation
	tests := []struct {
		name string
		op   interface {
			internal.PredicateOp
			internal.JSONOp
			internal.CompactOp
		}
		wantErr error
	}{
		{name: "and nil", op: NewAnd(nil, nil), wantErr: ErrInvalidPredicateInAnd},
		{name: "and typed nil", op: NewAnd(nil, typedNil), wantErr: ErrInvalidPredicateInAnd},
		{name: "or nil", op: NewOr(nil, nil), wantErr: ErrInvalidPredicateInOr},
		{name: "or typed nil", op: NewOr(nil, typedNil), wantErr: ErrInvalidPredicateInOr},
		{name: "not nil", op: NewNot(nil, nil), wantErr: ErrInvalidPredicateInNot},
		{name: "not typed nil", op: NewNot(nil, typedNil), wantErr: ErrInvalidPredicateInNot},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := tc.op.Test(map[string]any{})
			assert.False(t, matched)
			assert.ErrorIs(t, err, tc.wantErr)

			result, err := tc.op.Apply(map[string]any{})
			assert.Equal(t, internal.OpResult[any]{}, result)
			assert.ErrorIs(t, err, tc.wantErr)
			assert.ErrorIs(t, tc.op.Validate(), tc.wantErr)

			_, err = tc.op.Clone()
			assert.ErrorIs(t, err, tc.wantErr)
			_, err = tc.op.ToJSON()
			assert.ErrorIs(t, err, tc.wantErr)
			_, err = tc.op.ToCompact()
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestTestTypeRecognizesReflectiveJSONTypes(t *testing.T) {
	t.Parallel()

	type namedString string
	type namedStruct struct{ Name string }

	name := namedString("Ada")
	age := 36
	tests := []struct {
		name     string
		value    any
		wantType string
	}{
		{name: "named integer", value: int64(1), wantType: "integer"},
		{name: "whole float", value: 1.0, wantType: "integer"},
		{name: "named string", value: namedString("Ada"), wantType: "string"},
		{name: "array", value: [1]string{"Ada"}, wantType: "array"},
		{name: "struct", value: namedStruct{Name: "Ada"}, wantType: "object"},
		{name: "complex number", value: complex64(1), wantType: "number"},
		{name: "pointer to named string", value: &name, wantType: "string"},
		{name: "pointer to number", value: &age, wantType: "number"},
		{name: "function value", value: func() {}, wantType: "object"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			doc := map[string]any{"value": tc.value}
			op := NewTestType([]string{"value"}, tc.wantType)
			matched, err := op.Test(doc)
			require.NoError(t, err)
			assert.True(t, matched)
		})
	}
}
