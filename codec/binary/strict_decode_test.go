package binary

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCodecDecodeRejectsInvalidFixedFrameArity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		operation internal.Op
		arity     uint32
	}{
		{name: "add", operation: op.NewAdd(nil, "value"), arity: 3},
		{name: "move", operation: op.NewMove([]string{"to"}, []string{"from"}), arity: 3},
		{name: "copy", operation: op.NewCopy([]string{"to"}, []string{"from"}), arity: 3},
		{name: "test type", operation: op.NewTestType(nil, "string"), arity: 3},
		{name: "type", operation: op.NewType(nil, "string"), arity: 3},
		{name: "in", operation: op.NewIn(nil, []any{"value"}), arity: 3},
		{name: "less", operation: op.NewLess(nil, 1), arity: 3},
		{name: "more", operation: op.NewMore(nil, 1), arity: 3},
		{name: "inc", operation: op.NewInc(nil, 1), arity: 3},
		{name: "defined", operation: op.NewDefined(nil), arity: 2},
		{name: "undefined", operation: op.NewUndefined(nil), arity: 2},
		{name: "flip", operation: op.NewFlip(nil), arity: 2},
		{name: "str ins", operation: op.NewStrIns(nil, 0, "value"), arity: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, declared := range []uint32{tc.arity - 1, tc.arity + 1} {
				t.Run(fmt.Sprintf("declared_%d", declared), func(t *testing.T) {
					t.Parallel()

					decoded, err := New().Decode(binaryWithFrameArity(t, tc.operation, declared))
					require.Error(t, err)
					assert.Nil(t, decoded)
					assert.ErrorIs(t, err, ErrInvalidFrameArity)
				})
			}
		})
	}
}

func TestCodecDecodeDoesNotConsumeExtraFramePayloadAsNextOperation(t *testing.T) {
	t.Parallel()

	data := binaryFixture(t, func(writer *msgp.Writer) {
		require.NoError(t, writer.WriteArrayHeader(2))
		writeBinaryHeader(t, writer, 4, internal.OpAddCode, "first")
		require.NoError(t, writer.WriteString("value"))
		writeBinaryHeader(t, writer, 2, internal.OpDefinedCode, "second")
	})

	decoded, err := New().Decode(data)
	require.Error(t, err)
	assert.Nil(t, decoded)
	assert.ErrorIs(t, err, ErrInvalidFrameArity)
}

func TestCodecDecodeAcceptsOnlyOptionalFrameArities(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		operations []internal.Op
		minArity   uint32
		maxArity   uint32
	}{
		{name: "remove", operations: []internal.Op{op.NewRemove(nil), op.NewRemoveWithOldValue(nil, nil)}, minArity: 2, maxArity: 3},
		{name: "replace", operations: []internal.Op{op.NewReplace(nil, "value"), op.NewReplaceWithOldValue(nil, "value", nil)}, minArity: 3, maxArity: 4},
		{name: "test", operations: []internal.Op{op.NewTest(nil, "value"), op.NewTestWithNot(nil, "value", true)}, minArity: 3, maxArity: 4},
		{name: "contains", operations: []internal.Op{op.NewContains(nil, "value"), op.NewContainsWithIgnoreCase(nil, "value", true)}, minArity: 3, maxArity: 4},
		{name: "starts", operations: []internal.Op{op.NewStarts(nil, "value"), op.NewStartsWithIgnoreCase(nil, "value", true)}, minArity: 3, maxArity: 4},
		{name: "ends", operations: []internal.Op{op.NewEnds(nil, "value"), op.NewEndsWithIgnoreCase(nil, "value", true)}, minArity: 3, maxArity: 4},
		{name: "matches", operations: []internal.Op{op.NewMatches(nil, "value", false, nil), op.NewMatches(nil, "value", true, nil)}, minArity: 3, maxArity: 4},
		{name: "test string", operations: []internal.Op{op.NewTestString(nil, "value", 0, false, false), op.NewTestString(nil, "value", 0, true, false)}, minArity: 4, maxArity: 5},
		{name: "test string len", operations: []internal.Op{op.NewTestStringLen(nil, 1), op.NewTestStringLenWithNot(nil, 1, true)}, minArity: 3, maxArity: 4},
		{name: "str del", operations: []internal.Op{op.NewStrDelWithStr(nil, 0, "value"), op.NewStrDel(nil, 0, 1)}, minArity: 4, maxArity: 5},
		{name: "split", operations: []internal.Op{op.NewSplit(nil, 1, nil), op.NewSplit(nil, 1, map[string]any{"kind": "value"})}, minArity: 3, maxArity: 4},
		{name: "merge", operations: []internal.Op{op.NewMerge(nil, 1, nil), op.NewMerge(nil, 1, map[string]any{"kind": "value"})}, minArity: 3, maxArity: 4},
		{name: "extend", operations: []internal.Op{op.NewExtend(nil, map[string]any{"key": "value"}, false), op.NewExtend(nil, map[string]any{"key": "value"}, true)}, minArity: 3, maxArity: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, operation := range tc.operations {
				data, err := New().Encode([]internal.Op{operation})
				require.NoError(t, err)
				decoded, err := New().Decode(data)
				require.NoError(t, err)
				require.Len(t, decoded, 1)
			}

			for _, declared := range []uint32{tc.minArity - 1, tc.maxArity + 1} {
				t.Run(fmt.Sprintf("declared_%d", declared), func(t *testing.T) {
					t.Parallel()

					decoded, err := New().Decode(binaryWithFrameArity(t, tc.operations[0], declared))
					require.Error(t, err)
					assert.Nil(t, decoded)
					assert.ErrorIs(t, err, ErrInvalidFrameArity)
				})
			}
		})
	}
}

func TestCodecDecodeRejectsInvalidCompositeFrameArity(t *testing.T) {
	t.Parallel()

	operations := []struct {
		name      string
		operation internal.Op
	}{
		{name: "and", operation: op.NewAnd(nil, []any{op.NewDefined([]string{"child"})})},
		{name: "or", operation: op.NewOr(nil, []any{op.NewDefined([]string{"child"})})},
		{name: "not", operation: op.NewNotMultiple(nil, []any{op.NewDefined([]string{"child"})})},
	}
	for _, tc := range operations {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			for _, declared := range []uint32{2, 4} {
				decoded, err := New().Decode(binaryWithFrameArity(t, tc.operation, declared))
				require.Error(t, err)
				assert.Nil(t, decoded)
				assert.ErrorIs(t, err, ErrInvalidFrameArity)
			}
		})
	}

	nested := []struct {
		name  string
		write func(*testing.T, *msgp.Writer)
	}{
		{
			name: "child too small",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 1, internal.OpDefinedCode, "child")
			},
		},
		{
			name: "child too large",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpDefinedCode, "child")
			},
		},
		{
			name: "malformed child before sibling",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(2))
				writeBinaryHeader(t, writer, 3, internal.OpDefinedCode, "first")
				writeBinaryHeader(t, writer, 2, internal.OpUndefinedCode, "second")
			},
		},
	}
	for _, tc := range nested {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAndCode, "root")
				tc.write(t, writer)
			})
			decoded, err := New().Decode(data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, ErrInvalidFrameArity)
		})
	}
}

func TestCodecDecodeRejectsTrailingData(t *testing.T) {
	t.Parallel()

	codec := New()
	valid, err := codec.Encode([]internal.Op{op.NewAdd([]string{"name"}, "Ada")})
	require.NoError(t, err)
	decoded, err := codec.Decode(valid)
	require.NoError(t, err)
	require.Len(t, decoded, 1)

	secondPatch, err := codec.Encode([]internal.Op{op.NewRemove([]string{"name"})})
	require.NoError(t, err)
	tests := []struct {
		name   string
		suffix []byte
	}{
		{name: "reserved byte", suffix: []byte{0xc1}},
		{name: "scalar", suffix: []byte{0xc0}},
		{name: "empty patch", suffix: []byte{0x90}},
		{name: "second patch", suffix: secondPatch},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := append(append([]byte(nil), valid...), tc.suffix...)
			decoded, err := codec.Decode(data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, ErrTrailingData)
		})
	}

	t.Run("truncated payload", func(t *testing.T) {
		t.Parallel()

		decoded, err := codec.Decode(valid[:len(valid)-1])
		require.Error(t, err)
		assert.Nil(t, decoded)
		assert.NotErrorIs(t, err, ErrTrailingData)
	})

	for _, tc := range []struct {
		name string
		data []byte
	}{
		{name: "empty input", data: nil},
		{name: "non-array top level", data: []byte{0xc0}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := codec.Decode(tc.data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.NotErrorIs(t, err, ErrTrailingData)
		})
	}
}

func binaryWithFrameArity(t *testing.T, operation internal.Op, arity uint32) []byte {
	t.Helper()

	data, err := New().Encode([]internal.Op{operation})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(data), 2)
	require.Equal(t, byte(0x90), data[1]&0xf0, "operation must use a MessagePack fixarray header")
	require.LessOrEqual(t, arity, uint32(15))
	data[1] = 0x90 | byte(arity)
	return data
}
