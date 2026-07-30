package jsonpatch_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch"
	"github.com/kaptinlin/jsonpatch/codec/binary"
	"github.com/kaptinlin/jsonpatch/codec/compact"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestStringCoordinatesRoundTripAcrossCodecs(t *testing.T) {
	t.Parallel()

	operations := []internal.Op{
		op.NewStrIns([]string{"text"}, 2, "🙂"),
		op.NewStrDel([]string{"text"}, 2, 3),
		op.NewSplit([]string{"text"}, 2, nil),
		op.NewTestString([]string{"text"}, "e\u0301", 3, false),
		op.NewTestStringLen([]string{"text"}, 6),
	}

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		wire, err := jsoncodec.EncodeJSON(operations)
		require.NoError(t, err)
		decoded, err := jsoncodec.DecodeJSON(wire, jsoncodec.PatchOptions{})
		require.NoError(t, err)
		assertStringCoordinateFacts(t, decoded)
	})

	for name, options := range map[string][]compact.Option{
		"numeric compact": nil,
		"string compact":  {compact.WithStringOpcode(true)},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			wire, err := compact.Encode(operations, options...)
			require.NoError(t, err)
			decoded, err := compact.Decode(wire)
			require.NoError(t, err)
			assertStringCoordinateFacts(t, decoded)
		})
	}

	t.Run("binary", func(t *testing.T) {
		t.Parallel()
		codec := binary.New()
		wire, err := codec.Encode(operations)
		require.NoError(t, err)
		decoded, err := codec.Decode(wire)
		require.NoError(t, err)
		assertStringCoordinateFacts(t, decoded)
	})
}

func TestLargeCoordinatesRoundTripAcrossCodecs(t *testing.T) {
	t.Parallel()
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit int")
	}

	large64 := int64(1<<53 + 1)
	large := int(large64)
	operations := []internal.Op{
		op.NewStrIns([]string{"text"}, large, "x"),
		op.NewStrDel([]string{"text"}, large, 1),
		op.NewTestString([]string{"text"}, "x", large, false),
		op.NewTestStringLen([]string{"text"}, large),
		op.NewSplit([]string{"text"}, large, nil),
		op.NewMerge([]string{"nodes", "1"}, large, nil),
		op.NewAnd([]string{"profile"},
			op.NewTestString([]string{"profile", "name"}, "A", large, false),
		),
		op.NewStrDel([]string{"text"}, 1, int(^uint(0)>>1)),
	}
	assertCoordinates := func(t *testing.T, decoded []internal.Op) {
		t.Helper()
		require.Len(t, decoded, 8)
		assert.Equal(t, large, decoded[0].(*op.StrInsOperation).Pos)
		strDel := decoded[1].(*op.StrDelOperation)
		assert.Equal(t, large, strDel.Pos)
		assert.Equal(t, 1, strDel.Len)
		assert.Equal(t, large, decoded[2].(*op.TestStringOperation).Pos)
		assert.Equal(t, large, decoded[3].(*op.TestStringLenOperation).Length)
		assert.Equal(t, large, decoded[4].(*op.SplitOperation).Pos)
		assert.Equal(t, large, decoded[5].(*op.MergeOperation).Pos)
		and := decoded[6].(*op.AndOperation)
		require.Len(t, and.Operations, 1)
		assert.Equal(t, large, and.Operations[0].(*op.TestStringOperation).Pos)
		maxLengthDelete := decoded[7].(*op.StrDelOperation)
		assert.Equal(t, int(^uint(0)>>1), maxLengthDelete.Len)
		result, err := maxLengthDelete.Apply(map[string]any{"text": "abc"})
		require.NoError(t, err)
		assert.Equal(t, map[string]any{"text": "a"}, result.Doc)
	}

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		wire, err := jsoncodec.EncodeJSON(operations)
		require.NoError(t, err)
		decoded, err := jsoncodec.DecodeJSON(wire, jsoncodec.PatchOptions{})
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})

	t.Run("compact", func(t *testing.T) {
		t.Parallel()
		wire, err := compact.Encode(operations)
		require.NoError(t, err)
		decoded, err := compact.Decode(wire)
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})

	t.Run("compact JSON", func(t *testing.T) {
		t.Parallel()
		wire, err := compact.EncodeJSON(operations)
		require.NoError(t, err)
		decoded, err := compact.DecodeJSON(wire)
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})

	t.Run("binary", func(t *testing.T) {
		t.Parallel()
		codec := binary.New()
		wire, err := codec.Encode(operations)
		require.NoError(t, err)
		decoded, err := codec.Decode(wire)
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})
}

func TestCompileJSONAcceptsExactIntBoundaryCoordinates(t *testing.T) {
	t.Parallel()
	if strconv.IntSize < 64 {
		t.Skip("requires a 64-bit int")
	}

	patch, err := jsonpatch.CompileJSON(
		[]byte(`[
            {"op":"str_ins","path":"/text","pos":9223372036854775807.0,"str":"x"},
            {"op":"and","path":"/profile","apply":[{"op":"test_string","path":"/name","pos":9.223372036854775807e18,"str":"A"}]}
        ]`),
		jsonpatch.WithCapabilities(jsonpatch.Extended, jsonpatch.Predicate),
	)
	require.NoError(t, err)
	assert.Equal(t, 2, patch.Len())
}

func TestMergeMetadataCoordinatesRoundTripAcrossCodecs(t *testing.T) {
	t.Parallel()

	operations := []internal.Op{
		op.NewMerge([]string{"nodes", "1"}, -7, nil),
		op.NewMerge([]string{"nodes", "1"}, 0, nil),
		op.NewMerge([]string{"nodes", "1"}, 3, nil),
	}
	assertCoordinates := func(t *testing.T, decoded []internal.Op) {
		t.Helper()
		require.Len(t, decoded, 3)
		for i, want := range []int{-7, 0, 3} {
			merge := decoded[i].(*op.MergeOperation)
			assert.Equal(t, want, merge.Pos)
			result, err := merge.Apply(map[string]any{"nodes": []any{"a", "b", "c"}})
			require.NoError(t, err)
			assert.Equal(t, map[string]any{"nodes": []any{"ab", "c"}}, result.Doc)
		}
	}

	t.Run("JSON", func(t *testing.T) {
		t.Parallel()
		wire, err := jsoncodec.EncodeJSON(operations)
		require.NoError(t, err)
		decoded, err := jsoncodec.DecodeJSON(wire, jsoncodec.PatchOptions{})
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})

	t.Run("compact", func(t *testing.T) {
		t.Parallel()
		wire, err := compact.Encode(operations)
		require.NoError(t, err)
		decoded, err := compact.Decode(wire)
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})

	t.Run("binary", func(t *testing.T) {
		t.Parallel()
		codec := binary.New()
		wire, err := codec.Encode(operations)
		require.NoError(t, err)
		decoded, err := codec.Decode(wire)
		require.NoError(t, err)
		assertCoordinates(t, decoded)
	})
}

func assertStringCoordinateFacts(t *testing.T, operations []internal.Op) {
	t.Helper()
	require.Len(t, operations, 5)

	strIns, ok := operations[0].(*op.StrInsOperation)
	require.True(t, ok)
	assert.Equal(t, 2, strIns.Pos)

	strDel, ok := operations[1].(*op.StrDelOperation)
	require.True(t, ok)
	assert.Equal(t, 2, strDel.Pos)
	assert.Equal(t, 3, strDel.Len)

	split, ok := operations[2].(*op.SplitOperation)
	require.True(t, ok)
	assert.Equal(t, 2, split.Pos)

	testString, ok := operations[3].(*op.TestStringOperation)
	require.True(t, ok)
	assert.Equal(t, 3, testString.Pos)
	assert.Equal(t, "e\u0301", testString.Str)
	matched, err := testString.Test(map[string]any{"text": "ab🙂e\u0301"})
	require.NoError(t, err)
	assert.True(t, matched)

	testStringLen, ok := operations[4].(*op.TestStringLenOperation)
	require.True(t, ok)
	assert.Equal(t, 6, testStringLen.Length)
}
