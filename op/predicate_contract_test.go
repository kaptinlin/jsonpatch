package op

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/internal"
)

func TestOperationSerializationContracts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		op          internal.Op
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
			wantJSON:    internal.Operation{Op: "remove", Path: "/name", OldValue: "Ada"},
			wantCompact: internal.CompactOperation{internal.OpRemoveCode, []string{"name"}},
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
			wantJSON:    internal.Operation{Op: "str_del", Path: "/name", Pos: 1, Str: "da"},
			wantCompact: internal.CompactOperation{internal.OpStrDelCode, []string{"name"}, 1, "da"},
		},
		{
			name:        "split with props",
			op:          NewSplit([]string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}),
			wantType:    internal.OpSplitType,
			wantCode:    internal.OpSplitCode,
			wantJSON:    internal.Operation{Op: "split", Path: "/nodes/0", Pos: 1, Props: map[string]any{"kind": "paragraph"}},
			wantCompact: internal.CompactOperation{internal.OpSplitCode, []string{"nodes", "0"}, 1.0, map[string]any{"kind": "paragraph"}},
		},
		{
			name:        "merge with props",
			op:          NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}),
			wantType:    internal.OpMergeType,
			wantCode:    internal.OpMergeCode,
			wantJSON:    internal.Operation{Op: "merge", Path: "/nodes/1", Pos: 1, Props: map[string]any{"merged": true}},
			wantCompact: internal.CompactOperation{internal.OpMergeCode, []string{"nodes", "1"}, map[string]any{"merged": true}},
		},
		{
			name:        "extend with delete null",
			op:          NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, true),
			wantType:    internal.OpExtendType,
			wantCode:    internal.OpExtendCode,
			wantJSON:    internal.Operation{Op: "extend", Path: "/profile", Props: map[string]any{"name": "Ada"}, DeleteNull: true},
			wantCompact: internal.CompactOperation{internal.OpExtendCode, []string{"profile"}, map[string]any{"name": "Ada"}},
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
		{name: "copy empty from", op: NewCopy([]string{"name"}, nil), wantErr: ErrFromPathEmpty},
		{name: "copy identical paths", op: NewCopy([]string{"name"}, []string{"name"}), wantErr: ErrPathsIdentical},
		{name: "extend nil properties", op: NewExtend([]string{"profile"}, nil, false), wantErr: ErrPropertiesNil},
		{name: "merge negative position", op: NewMerge([]string{"nodes"}, -1, nil), wantErr: ErrPositionNegative},
		{name: "str_del missing length", op: NewStrDel([]string{"name"}, 1, 0), wantErr: ErrMissingStrOrLen},
		{name: "test_string empty path", op: NewTestString(nil, "Ada", 0, false, false), wantErr: ErrPathEmpty},
		{name: "test_string_len empty path", op: NewTestStringLen(nil, 3), wantErr: ErrPathEmpty},
		{name: "test_string_len negative length", op: NewTestStringLen([]string{"name"}, -1), wantErr: ErrLengthNegative},
		{name: "test_string_len fractional length", op: NewTestStringLen([]string{"name"}, 1.5), wantErr: ErrInvalidLength},
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

	t.Run("merge text nodes preserves later marks and applies props", func(t *testing.T) {
		t.Parallel()

		doc := map[string]any{"nodes": []any{
			map[string]any{"text": "A", "bold": true},
			map[string]any{"text": "da", "italic": true},
		}}
		result, err := NewMerge([]string{"nodes"}, 1, map[string]any{"merged": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"nodes": []any{
			map[string]any{"text": "Ada", "bold": true, "italic": true, "merged": true},
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
		result, err := NewMerge([]string{"nodes"}, 1, map[string]any{"merged": true}).Apply(doc)
		require.NoError(t, err)

		wantDoc := map[string]any{"nodes": []any{
			map[string]any{"type": "paragraph", "align": "center", "merged": true, "children": []any{map[string]any{"text": "A"}, map[string]any{"text": "B"}}},
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
		op          internal.PredicateOp
		wantJSON    internal.Operation
		wantCompact internal.CompactOperation
	}{
		{
			name:        "contains ignore case",
			op:          NewContainsWithIgnoreCase([]string{"name"}, "lovelace", true),
			wantJSON:    internal.Operation{Op: "contains", Path: "/name", Value: "lovelace", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpContainsCode, []string{"name"}, "lovelace"},
		},
		{
			name:        "starts ignore case",
			op:          NewStartsWithIgnoreCase([]string{"name"}, "ada", true),
			wantJSON:    internal.Operation{Op: "starts", Path: "/name", Value: "ada", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpStartsCode, []string{"name"}, "ada"},
		},
		{
			name:        "ends ignore case",
			op:          NewEndsWithIgnoreCase([]string{"name"}, "LOVELACE", true),
			wantJSON:    internal.Operation{Op: "ends", Path: "/name", Value: "LOVELACE", IgnoreCase: true},
			wantCompact: internal.CompactOperation{internal.OpEndsCode, []string{"name"}, "LOVELACE"},
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
			name:        "test type multiple",
			op:          NewTestTypeMultiple([]string{"name"}, []string{"string", "null"}),
			wantJSON:    internal.Operation{Op: "test_type", Path: "/name", Value: []string{"string", "null"}},
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

	doc := map[string]any{
		"name":  "Ada",
		"role":  "viewer",
		"score": 5.0,
	}

	tests := []struct {
		name    string
		op      internal.PredicateOp
		wantErr error
	}{
		{name: "contains mismatch", op: NewContains([]string{"name"}, "Grace"), wantErr: ErrStringMismatch},
		{name: "starts mismatch", op: NewStarts([]string{"name"}, "Grace"), wantErr: ErrStringMismatch},
		{name: "ends mismatch", op: NewEnds([]string{"name"}, "Grace"), wantErr: ErrStringMismatch},
		{name: "in mismatch", op: NewIn([]string{"role"}, []any{"admin"}), wantErr: ErrOperationFailed},
		{name: "less mismatch", op: NewLess([]string{"score"}, 5), wantErr: ErrComparisonFailed},
		{name: "more mismatch", op: NewMore([]string{"score"}, 5), wantErr: ErrComparisonFailed},
		{name: "type mismatch", op: NewType([]string{"name"}, "number"), wantErr: ErrTypeMismatch},
		{name: "test_type mismatch", op: NewTestType([]string{"name"}, "number"), wantErr: ErrTypeMismatch},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			matched, err := tc.op.Test(doc)
			require.NoError(t, err)
			assert.False(t, matched)

			result, err := tc.op.Apply(doc)
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.wantErr)
			if diff := cmp.Diff(internal.OpResult[any]{}, result); diff != "" {
				t.Errorf("Apply() result mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPredicateOperationValidateRejectsInvalidOperands(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      internal.PredicateOp
		wantErr error
	}{
		{name: "contains empty path", op: NewContains(nil, "a"), wantErr: ErrPathEmpty},
		{name: "starts empty path", op: NewStarts(nil, "a"), wantErr: ErrPathEmpty},
		{name: "ends empty path", op: NewEnds(nil, "a"), wantErr: ErrPathEmpty},
		{name: "in empty path", op: NewIn(nil, []any{"a"}), wantErr: ErrPathEmpty},
		{name: "in empty values", op: NewIn([]string{"role"}, nil), wantErr: ErrValuesArrayEmpty},
		{name: "less empty path", op: NewLess(nil, 1), wantErr: ErrPathEmpty},
		{name: "more empty path", op: NewMore(nil, 1), wantErr: ErrPathEmpty},
		{name: "type empty path", op: NewType(nil, "string"), wantErr: ErrPathEmpty},
		{name: "type invalid type", op: NewType([]string{"name"}, "invalid"), wantErr: ErrInvalidType},
		{name: "test_type empty path", op: NewTestType(nil, "string"), wantErr: ErrPathEmpty},
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
