package binary

import (
	"errors"
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCodecDecodeAcceptsCanonicalIntegerTrueMarker(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		writeMarker func(*msgp.Writer) error
	}{
		{name: "signed integer", writeMarker: func(writer *msgp.Writer) error { return writer.WriteInt8(1) }},
		{name: "unsigned integer", writeMarker: func(writer *msgp.Writer) error { return writer.WriteUint64(1) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := New().Decode(binaryTestWithMarker(t, tc.writeMarker))
			require.NoError(t, err)
			require.Len(t, decoded, 1)
			testOperation, ok := decoded[0].(*op.TestOperation)
			require.True(t, ok)
			assert.True(t, testOperation.Not())
		})
	}
}

func TestCodecDecodeRejectsInvalidTrueMarkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		writeMarker func(*msgp.Writer) error
	}{
		{name: "boolean", writeMarker: func(writer *msgp.Writer) error { return writer.WriteBool(true) }},
		{name: "zero", writeMarker: func(writer *msgp.Writer) error { return writer.WriteInt64(0) }},
		{name: "integer other than one", writeMarker: func(writer *msgp.Writer) error { return writer.WriteUint8(2) }},
		{name: "float one", writeMarker: func(writer *msgp.Writer) error { return writer.WriteFloat64(1) }},
		{name: "string one", writeMarker: func(writer *msgp.Writer) error { return writer.WriteString("1") }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := New().Decode(binaryTestWithMarker(t, tc.writeMarker))
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, ErrInvalidTrueMarker)
		})
	}
}

func binaryTestWithMarker(t *testing.T, writeMarker func(*msgp.Writer) error) []byte {
	t.Helper()

	return binaryFixture(t, func(writer *msgp.Writer) {
		require.NoError(t, writer.WriteArrayHeader(1))
		writeBinaryHeader(t, writer, 4, internal.OpTestCode, "status")
		require.NoError(t, writer.WriteString("ready"))
		require.NoError(t, writeMarker(writer))
	})
}

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

func TestCodecDecodeRejectsInvalidStringEditCoordinates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		write   func(*testing.T, *msgp.Writer)
		wantErr error
	}{
		{name: "str_ins fractional position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpStrInsCode, "text")
			require.NoError(t, writer.WriteFloat64(1.5))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberNotInteger},
		{name: "str_ins non-finite position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpStrInsCode, "text")
			require.NoError(t, writer.WriteFloat64(math.Inf(1)))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberNotFinite},
		{name: "str_ins overflowing position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpStrInsCode, "text")
			require.NoError(t, writer.WriteFloat64(math.MaxFloat64))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberOutOfRange},
		{name: "str_del fractional position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpStrDelCode, "text")
			require.NoError(t, writer.WriteFloat64(1.5))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberNotInteger},
		{name: "str_del non-finite length", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 5, internal.OpStrDelCode, "text")
			require.NoError(t, writer.WriteFloat64(0))
			require.NoError(t, writer.WriteUint8(0))
			require.NoError(t, writer.WriteFloat64(math.NaN()))
		}, wantErr: ErrNumberNotFinite},
		{name: "str_del overflowing length", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 5, internal.OpStrDelCode, "text")
			require.NoError(t, writer.WriteFloat64(0))
			require.NoError(t, writer.WriteUint8(0))
			require.NoError(t, writer.WriteFloat64(math.MaxFloat64))
		}, wantErr: ErrNumberOutOfRange},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				tc.write(t, writer)
			})
			decoded, err := New().Decode(data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCodecDecodeRejectsInvalidStringPredicateCoordinates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		write   func(*testing.T, *msgp.Writer)
		wantErr error
	}{
		{name: "test_string fractional position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpTestStringCode, "text")
			require.NoError(t, writer.WriteFloat64(1.5))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberNotInteger},
		{name: "test_string non-finite position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpTestStringCode, "text")
			require.NoError(t, writer.WriteFloat64(math.Inf(1)))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberNotFinite},
		{name: "test_string overflowing position", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 4, internal.OpTestStringCode, "text")
			require.NoError(t, writer.WriteFloat64(math.MaxFloat64))
			require.NoError(t, writer.WriteString("x"))
		}, wantErr: ErrNumberOutOfRange},
		{name: "test_string_len fractional length", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 3, internal.OpTestStringLenCode, "text")
			require.NoError(t, writer.WriteFloat64(1.5))
		}, wantErr: ErrNumberNotInteger},
		{name: "test_string_len non-finite length", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 3, internal.OpTestStringLenCode, "text")
			require.NoError(t, writer.WriteFloat64(math.NaN()))
		}, wantErr: ErrNumberNotFinite},
		{name: "test_string_len overflowing length", write: func(t *testing.T, writer *msgp.Writer) {
			writeBinaryHeader(t, writer, 3, internal.OpTestStringLenCode, "text")
			require.NoError(t, writer.WriteFloat64(math.MaxFloat64))
		}, wantErr: ErrNumberOutOfRange},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data := binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				tc.write(t, writer)
			})
			decoded, err := New().Decode(data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCodecDecodeRejectsInvalidSplitMergeCoordinates(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		name    string
		code    uint8
		value   float64
		wantErr error
	}{
		{name: "split fractional position", code: internal.OpSplitCode, value: 1.5, wantErr: ErrNumberNotInteger},
		{name: "split non-finite position", code: internal.OpSplitCode, value: math.Inf(1), wantErr: ErrNumberNotFinite},
		{name: "split overflowing position", code: internal.OpSplitCode, value: math.MaxFloat64, wantErr: ErrNumberOutOfRange},
		{name: "merge fractional position", code: internal.OpMergeCode, value: -1.5, wantErr: ErrNumberNotInteger},
		{name: "merge non-finite position", code: internal.OpMergeCode, value: math.NaN(), wantErr: ErrNumberNotFinite},
		{name: "merge overflowing position", code: internal.OpMergeCode, value: math.MaxFloat64, wantErr: ErrNumberOutOfRange},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data := binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, int(tc.code), "node")
				require.NoError(t, writer.WriteFloat64(tc.value))
			})
			decoded, err := New().Decode(data)
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCodecDecodeRejectsNonStringMapKey(t *testing.T) {
	t.Parallel()

	keys := []struct {
		name  string
		write func(*testing.T, *msgp.Writer)
	}{
		{name: "int", write: func(t *testing.T, writer *msgp.Writer) {
			require.NoError(t, writer.WriteInt64(7))
		}},
		{name: "bool", write: func(t *testing.T, writer *msgp.Writer) {
			require.NoError(t, writer.WriteBool(true))
		}},
		{name: "array", write: func(t *testing.T, writer *msgp.Writer) {
			require.NoError(t, writer.WriteArrayHeader(1))
			require.NoError(t, writer.WriteString("key"))
		}},
		{name: "map", write: func(t *testing.T, writer *msgp.Writer) {
			require.NoError(t, writer.WriteMapHeader(1))
			require.NoError(t, writer.WriteString("key"))
			require.NoError(t, writer.WriteString("value"))
		}},
	}

	for _, nested := range []bool{false, true} {
		location := "top_level"
		if nested {
			location = "nested"
		}
		for _, key := range keys {
			t.Run(location+"_"+key.name, func(t *testing.T) {
				t.Parallel()
				data := binaryFixture(t, func(writer *msgp.Writer) {
					require.NoError(t, writer.WriteArrayHeader(1))
					writeBinaryHeader(t, writer, 3, internal.OpAddCode)
					if nested {
						require.NoError(t, writer.WriteMapHeader(1))
						require.NoError(t, writer.WriteString("values"))
						require.NoError(t, writer.WriteArrayHeader(1))
					}
					require.NoError(t, writer.WriteMapHeader(1))
					key.write(t, writer)
					require.NoError(t, writer.WriteString("discarded"))
				})

				decoded, err := New().Decode(data)
				require.Error(t, err)
				assert.Nil(t, decoded)
				assert.ErrorIs(t, err, ErrMapKeyNotString)
				var typeErr msgp.TypeError
				assert.True(t, errors.As(err, &typeErr))
				assert.Equal(t, msgp.StrType, typeErr.Method)
			})
		}
	}
}

func TestCodecDecodeNonStringMapKeyReturnsNoPartialOperations(t *testing.T) {
	t.Parallel()

	data := binaryFixture(t, func(writer *msgp.Writer) {
		require.NoError(t, writer.WriteArrayHeader(2))
		writeBinaryHeader(t, writer, 3, internal.OpAddCode, "first")
		require.NoError(t, writer.WriteString("ok"))
		writeBinaryHeader(t, writer, 3, internal.OpAddCode, "second")
		require.NoError(t, writer.WriteMapHeader(1))
		require.NoError(t, writer.WriteInt64(7))
		require.NoError(t, writer.WriteString("discarded"))
	})

	decoded, err := New().Decode(data)
	require.ErrorIs(t, err, ErrMapKeyNotString)
	assert.Nil(t, decoded)
}

func TestCodecDecodePreservesNestedStringKeyMaps(t *testing.T) {
	t.Parallel()

	want := map[string]any{
		"values": []any{map[string]any{"nested": "value"}},
	}
	codec := New()
	wire, err := codec.Encode([]internal.Op{op.NewAdd(nil, want)})
	require.NoError(t, err)
	decoded, err := codec.Decode(wire)
	require.NoError(t, err)
	require.Len(t, decoded, 1)
	add, ok := decoded[0].(*op.AddOperation)
	require.True(t, ok)
	assert.Equal(t, want, add.Value)
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
		{name: "test string", operations: []internal.Op{op.NewTestString(nil, "value", 0, false), op.NewTestString(nil, "value", 0, true)}, minArity: 4, maxArity: 5},
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
		{name: "and", operation: op.NewAnd(nil, op.NewDefined([]string{"child"}))},
		{name: "or", operation: op.NewOr(nil, op.NewDefined([]string{"child"}))},
		{name: "not", operation: op.NewNot(nil, op.NewDefined([]string{"child"}))},
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
