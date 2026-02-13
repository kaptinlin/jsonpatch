package binarytests

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

var (
	// Patches is a collection of test cases for the binary codec.
	Patches = []struct {
		name  string
		patch []internal.Op
	}{
		{
			name:  "AddOperation",
			patch: []internal.Op{op.NewAdd([]string{"a", "b", "c"}, "foo")},
		},
		{
			name:  "RemoveOperation",
			patch: []internal.Op{op.NewRemove([]string{"a", "b", "c"})},
		},
		{
			name:  "ReplaceOperation",
			patch: []internal.Op{op.NewReplace([]string{"a", "b", "c"}, "bar")},
		},
		{
			name:  "MoveOperation",
			patch: []internal.Op{op.NewMove([]string{"a", "b", "c"}, []string{"a", "b", "d"})},
		},
		{
			name:  "CopyOperation",
			patch: []internal.Op{op.NewCopy([]string{"a", "b", "d"}, []string{"a", "b", "e"})},
		},
		{
			name:  "TestOperation",
			patch: []internal.Op{op.NewTest([]string{"a", "b", "e"}, "bar")},
		},
		{
			name:  "TestTypeOperationMultiple",
			patch: []internal.Op{op.NewTestTypeMultiple([]string{"a", "b", "e"}, []string{"string", "number"})},
		},
		{
			name:  "DefinedOperation",
			patch: []internal.Op{op.NewDefined([]string{"a", "b", "e"})},
		},
		{
			name:  "UndefinedOperation",
			patch: []internal.Op{op.NewUndefined([]string{"a", "b", "f"})},
		},
		{
			name:  "LessOperation",
			patch: []internal.Op{op.NewLess([]string{"a", "num"}, 100)},
		},
		{
			name:  "MoreOperation",
			patch: []internal.Op{op.NewMore([]string{"a", "num"}, 5)},
		},
		{
			name:  "ContainsOperation",
			patch: []internal.Op{op.NewContains([]string{"a", "str"}, "world")},
		},
		{
			name:  "InOperation",
			patch: []internal.Op{op.NewIn([]string{"a", "b", "e"}, []any{"foo", "bar", "baz"})},
		},
		{
			name:  "StartsOperation",
			patch: []internal.Op{op.NewStarts([]string{"a", "str"}, "hello")},
		},
		{
			name:  "EndsOperation",
			patch: []internal.Op{op.NewEnds([]string{"a", "str"}, "world")},
		},
		{
			name:  "MatchesOperation",
			patch: []internal.Op{mustNewMatchesOperation([]string{"a", "str"}, "^hello.*d$", true)},
		},
		{
			name:  "TestStringOperationWithPos",
			patch: []internal.Op{op.NewTestString([]string{"a", "str"}, "lo", 3, false, false)},
		},
		{
			name:  "TestStringLenOperationWithNot",
			patch: []internal.Op{op.NewTestStringLenWithNot([]string{"a", "str"}, 12, false)},
		},
		{
			name:  "TypeOperation",
			patch: []internal.Op{op.NewType([]string{"a", "val"}, "number")},
		},
		{
			name:  "FlipOperation",
			patch: []internal.Op{op.NewFlip([]string{"a", "bool"})},
		},
		{
			name:  "IncOperation",
			patch: []internal.Op{op.NewInc([]string{"a", "num"}, 10)},
		},
		{
			name:  "StrInsOperation",
			patch: []internal.Op{op.NewStrIns([]string{"a", "str"}, 6, " beautiful")},
		},
		{
			name:  "StrDelOperation",
			patch: []internal.Op{op.NewStrDel([]string{"a", "str"}, 0, 5)},
		},
		{
			name:  "SplitOperation",
			patch: []internal.Op{op.NewSplit([]string{"a", "str"}, 5, nil)},
		},
		{
			name:  "ExtendOperation",
			patch: []internal.Op{op.NewExtend([]string{"a", "obj"}, map[string]any{"q": 4, "r": 5}, false)},
		},
		{
			name:  "MergeOperation",
			patch: []internal.Op{op.NewMerge([]string{"a", "arr"}, 1, map[string]any{"merged": true})},
		},
	}
)

func mustNewMatchesOperation(path []string, pattern string, ignoreCase bool) *op.MatchesOperation {
	return op.NewMatches(path, pattern, ignoreCase, nil)
}

func TestRoundtrip(t *testing.T) {
	t.Parallel()
	codec := binary.Codec{}
	for _, tt := range Patches {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			encoded, err := codec.Encode(tt.patch)
			require.NoError(t, err)

			decoded, err := codec.Decode(encoded)
			require.NoError(t, err)

			if !areOpsEqual(tt.patch, decoded) {
				assert.Fail(t, "decoded patch should equal original patch")
			}
		})
	}
}

// areOpsEqual is a helper function to compare two slices of operations.
// It is needed because reflect.DeepEqual might fail on comparing regex objects
// and different number types (e.g. int vs int64).
func areOpsEqual(a, b []internal.Op) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !isOpEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func isOpEqual(a, b internal.Op) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		return false
	}
	// Special case for 'matches' operation due to embedded regexp object.
	if ma, ok := a.(*op.MatchesOperation); ok {
		mb, ok := b.(*op.MatchesOperation)
		if !ok {
			return false // Should not happen if types are the same.
		}
		// Compare public fields that define the operation.
		return reflect.DeepEqual(ma.Path(), mb.Path()) && ma.Pattern == mb.Pattern && ma.IgnoreCase == mb.IgnoreCase
	}

	// For other operations, convert to standard JSON map format and compare deeply.
	jsonA, errA := a.ToJSON()
	jsonB, errB := b.ToJSON()

	if errA != nil || errB != nil {
		// If conversion fails, they are not equal.
		return false
	}

	return areOperationsEqual(jsonA, jsonB)
}

func areOperationsEqual(a, b internal.Operation) bool {
	if a.Op != b.Op || a.Path != b.Path {
		return false
	}
	if !areValuesEqual(a.Value, b.Value) {
		return false
	}
	if a.From != b.From || a.Str != b.Str || a.Type != b.Type {
		return false
	}
	if !areNumericEqual(a.Inc, b.Inc) {
		return false
	}
	if a.Pos != b.Pos || a.Len != b.Len {
		return false
	}
	if a.Not != b.Not || a.IgnoreCase != b.IgnoreCase || a.DeleteNull != b.DeleteNull {
		return false
	}
	if !areMapsEqual(a.Props, b.Props) {
		return false
	}
	if !areValuesEqual(a.OldValue, b.OldValue) {
		return false
	}
	if len(a.Apply) != len(b.Apply) {
		return false
	}
	for i := range a.Apply {
		if !areOperationsEqual(a.Apply[i], b.Apply[i]) {
			return false
		}
	}
	return true
}

func areNumericEqual(a, b float64) bool {
	if a == 0 && b == 0 {
		return true
	}
	return a == b
}

func areMapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, vA := range a {
		vB, ok := b[k]
		if !ok {
			return false
		}
		if !areValuesEqual(vA, vB) {
			return false
		}
	}
	return true
}

func areValuesEqual(a, b any) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		aFloat, aIsNum := op.ToFloat64(a)
		bFloat, bIsNum := op.ToFloat64(b)
		if aIsNum && bIsNum {
			return aFloat == bFloat
		}
	}
	switch aVal := a.(type) {
	case []any:
		bVal, ok := b.([]any)
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !areValuesEqual(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	case map[string]any:
		bVal, ok := b.(map[string]any)
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		return areMapsEqual(aVal, bVal)
	case map[any]any:
		bVal, ok := b.(map[string]any)
		if !ok {
			if bVal, ok := b.(map[any]any); ok {
				return areMapsEqual(convertMapToS(aVal), convertMapToS(bVal))
			}
			return false
		}
		return areMapsEqual(convertMapToS(aVal), bVal)
	default:
		return reflect.DeepEqual(a, b)
	}
}

func convertMapToS(m map[any]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		if strKey, ok := k.(string); ok {
			res[strKey] = v
		}
	}
	return res
}
