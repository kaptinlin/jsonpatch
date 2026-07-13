package jsonpatch_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		op.NewTestString([]string{"text"}, "e\u0301", 3, false, false),
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
	assert.Equal(t, 2.0, split.Pos)

	testString, ok := operations[3].(*op.TestStringOperation)
	require.True(t, ok)
	assert.Equal(t, 3, testString.Pos)
	assert.Equal(t, "e\u0301", testString.Str)

	testStringLen, ok := operations[4].(*op.TestStringLenOperation)
	require.True(t, ok)
	assert.Equal(t, 6.0, testStringLen.Length)
}
