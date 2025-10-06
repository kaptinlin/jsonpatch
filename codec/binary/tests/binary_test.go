package binarytests

import (
	"reflect"
	"testing"

	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			patch: []internal.Op{op.NewOpTestTypeOperationMultiple([]string{"a", "b", "e"}, []string{"string", "number"})},
		},
		{
			name:  "DefinedOperation",
			patch: []internal.Op{op.NewOpDefinedOperation([]string{"a", "b", "e"})},
		},
		{
			name:  "UndefinedOperation",
			patch: []internal.Op{op.NewOpUndefinedOperation([]string{"a", "b", "f"})},
		},
		{
			name:  "LessOperation",
			patch: []internal.Op{op.NewOpLessOperation([]string{"a", "num"}, 100)},
		},
		{
			name:  "MoreOperation",
			patch: []internal.Op{op.NewOpMoreOperation([]string{"a", "num"}, 5)},
		},
		{
			name:  "ContainsOperation",
			patch: []internal.Op{op.NewOpContainsOperation([]string{"a", "str"}, "world")},
		},
		{
			name:  "InOperation",
			patch: []internal.Op{op.NewOpInOperation([]string{"a", "b", "e"}, []interface{}{"foo", "bar", "baz"})},
		},
		{
			name:  "StartsOperation",
			patch: []internal.Op{op.NewOpStartsOperation([]string{"a", "str"}, "hello")},
		},
		{
			name:  "EndsOperation",
			patch: []internal.Op{op.NewOpEndsOperation([]string{"a", "str"}, "world")},
		},
		{
			name:  "MatchesOperation",
			patch: []internal.Op{mustNewMatchesOperation([]string{"a", "str"}, "^hello.*d$", true)},
		},
		{
			name:  "TestStringOperationWithPos",
			patch: []internal.Op{op.NewOpTestStringOperationWithPos([]string{"a", "str"}, "lo", 3)},
		},
		{
			name:  "TestStringLenOperationWithNot",
			patch: []internal.Op{op.NewOpTestStringLenOperationWithNot([]string{"a", "str"}, 12, false)},
		},
		{
			name:  "TypeOperation",
			patch: []internal.Op{op.NewOpTypeOperation([]string{"a", "val"}, "number")},
		},
		{
			name:  "FlipOperation",
			patch: []internal.Op{op.NewOpFlipOperation([]string{"a", "bool"})},
		},
		{
			name:  "IncOperation",
			patch: []internal.Op{op.NewOpIncOperation([]string{"a", "num"}, 10)},
		},
		{
			name:  "StrInsOperation",
			patch: []internal.Op{op.NewOpStrInsOperation([]string{"a", "str"}, 6, " beautiful")},
		},
		{
			name:  "StrDelOperation",
			patch: []internal.Op{op.NewOpStrDelOperation([]string{"a", "str"}, 0, 5)},
		},
		{
			name:  "SplitOperation",
			patch: []internal.Op{op.NewOpSplitOperation([]string{"a", "str"}, 5, nil)},
		},
		{
			name:  "ExtendOperation",
			patch: []internal.Op{op.NewOpExtendOperation([]string{"a", "obj"}, map[string]interface{}{"q": 4, "r": 5}, false)},
		},
		{
			name:  "MergeOperation",
			patch: []internal.Op{op.NewOpMergeOperation([]string{"a", "arr"}, 1, map[string]interface{}{"merged": true})},
		},
	}
)

func mustNewMatchesOperation(path []string, pattern string, ignoreCase bool) *op.MatchesOperation {
	op, err := op.NewOpMatchesOperation(path, pattern, ignoreCase)
	if err != nil {
		panic(err)
	}
	return op
}

func TestRoundtrip(t *testing.T) {
	codec := binary.Codec{}
	for _, tt := range Patches {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := codec.Encode(tt.patch)
			require.NoError(t, err, "Encode should not error")

			decoded, err := codec.Decode(encoded)
			require.NoError(t, err, "Decode should not error")

			assert.True(t, areOpsEqual(tt.patch, decoded), "decoded patch should equal original patch")
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
	// Compare basic fields
	if a.Op != b.Op || a.Path != b.Path {
		return false
	}

	// Compare values with special handling for numbers
	if !areValuesEqual(a.Value, b.Value) {
		return false
	}

	// Compare other fields
	if a.From != b.From || a.Str != b.Str || a.Type != b.Type {
		return false
	}

	// Compare numeric fields with tolerance
	if !areNumericEqual(a.Inc, b.Inc) {
		return false
	}

	// Compare integer fields
	if a.Pos != b.Pos || a.Len != b.Len {
		return false
	}

	// Compare boolean fields
	if a.Not != b.Not || a.IgnoreCase != b.IgnoreCase || a.DeleteNull != b.DeleteNull {
		return false
	}

	// Compare map fields
	if !areMapsEqual(a.Props, b.Props) {
		return false
	}

	// Compare OldValue
	if !areValuesEqual(a.OldValue, b.OldValue) {
		return false
	}

	// Compare Apply operations
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
	// Handle zero values
	if a == 0 && b == 0 {
		return true
	}

	// Simple equality check for non-zero values
	return a == b
}

func areMapsEqual(a, b map[string]interface{}) bool {
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

func areValuesEqual(a, b interface{}) bool {
	if reflect.TypeOf(a) != reflect.TypeOf(b) {
		// Attempt numeric conversion if types are different but numeric
		aFloat, aIsNum := op.ToFloat64(a)
		bFloat, bIsNum := op.ToFloat64(b)
		if aIsNum && bIsNum {
			return aFloat == bFloat
		}
	}

	switch aVal := a.(type) {
	case []interface{}:
		bVal, ok := b.([]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		for i := range aVal {
			if !areValuesEqual(aVal[i], bVal[i]) {
				return false
			}
		}
		return true
	case map[string]interface{}:
		bVal, ok := b.(map[string]interface{})
		if !ok || len(aVal) != len(bVal) {
			return false
		}
		return areMapsEqual(aVal, bVal)
	case map[interface{}]interface{}:
		// Handle comparison with map[string]interface{}
		bVal, ok := b.(map[string]interface{})
		if !ok {
			// If b is also map[interface{}]interface{}, convert both to map[string]interface{}
			if bVal, ok := b.(map[interface{}]interface{}); ok {
				return areMapsEqual(convertMapToS(aVal), convertMapToS(bVal))
			}
			return false
		}
		return areMapsEqual(convertMapToS(aVal), bVal)
	default:
		return reflect.DeepEqual(a, b)
	}
}

func convertMapToS(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		if strKey, ok := k.(string); ok {
			res[strKey] = v
		}
	}
	return res
}
