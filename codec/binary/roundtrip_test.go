package binary

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCodecRoundTripPreservesOperationJSON(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		op   internal.Op
	}{
		{name: "add nested object", op: op.NewAdd([]string{"profile"}, map[string]any{"name": "Ada", "tags": []any{"go"}})},
		{name: "remove without old value", op: op.NewRemove([]string{"profile", "name"})},
		{name: "remove with old value", op: op.NewRemoveWithOldValue([]string{"profile", "name"}, "Ada")},
		{name: "replace scalar", op: op.NewReplace([]string{"profile", "name"}, "Grace")},
		{name: "move from source to target", op: op.NewMove([]string{"profile", "displayName"}, []string{"profile", "name"})},
		{name: "copy from source to target", op: op.NewCopy([]string{"profile", "alias"}, []string{"profile", "name"})},
		{name: "test value", op: op.NewTest([]string{"profile", "name"}, "Ada")},
		{name: "defined predicate", op: op.NewDefined([]string{"profile", "name"})},
		{name: "undefined predicate", op: op.NewUndefined([]string{"profile", "deleted"})},
		{name: "test type multiple", op: op.NewTestTypeMultiple([]string{"profile", "name"}, []string{"string", "null"})},
		{name: "less predicate", op: op.NewLess([]string{"score"}, 10)},
		{name: "more predicate", op: op.NewMore([]string{"score"}, 5)},
		{name: "contains predicate", op: op.NewContains([]string{"profile", "name"}, "Ad")},
		{name: "starts predicate", op: op.NewStarts([]string{"profile", "name"}, "A")},
		{name: "ends predicate", op: op.NewEnds([]string{"profile", "name"}, "a")},
		{name: "in predicate", op: op.NewIn([]string{"role"}, []any{"admin", "editor"})},
		{name: "matches predicate", op: op.NewMatches([]string{"profile", "name"}, "^ad", true, nil)},
		{name: "test string predicate", op: op.NewTestString([]string{"profile", "name"}, "da", 1, false, false)},
		{name: "test string length predicate", op: op.NewTestStringLenWithNot([]string{"profile", "name"}, 3, true)},
		{name: "type predicate", op: op.NewType([]string{"profile", "name"}, "string")},
		{name: "flip operation", op: op.NewFlip([]string{"enabled"})},
		{name: "inc operation", op: op.NewInc([]string{"count"}, 2)},
		{name: "str ins operation", op: op.NewStrIns([]string{"profile", "name"}, 1, "d")},
		{name: "str del operation", op: op.NewStrDel([]string{"profile", "name"}, 1, 2)},
		{name: "split without props", op: op.NewSplit([]string{"nodes", "0"}, 1, nil)},
		{name: "split with props", op: op.NewSplit([]string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"})},
		{name: "extend operation", op: op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada", "nested": map[string]any{"ok": true}}, true)},
		{name: "merge without props", op: op.NewMerge([]string{"nodes", "1"}, 1, nil)},
		{name: "merge with props", op: op.NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true})},
	}

	codec := New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := codec.Encode([]internal.Op{tc.op})
			require.NoError(t, err)
			require.NotEmpty(t, encoded)

			decoded, err := codec.Decode(encoded)
			require.NoError(t, err)
			require.Len(t, decoded, 1)

			want, err := tc.op.ToJSON()
			require.NoError(t, err)
			got, err := decoded[0].ToJSON()
			require.NoError(t, err)
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("round-tripped operation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCodecDecodeRejectsMalformedOperations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		ops     []internal.Op
		mutate  func([]byte) []byte
		wantErr error
	}{
		{
			name: "unsupported operation code",
			ops:  []internal.Op{op.NewAdd([]string{"name"}, "Ada")},
			mutate: func(data []byte) []byte {
				var buf bytes.Buffer
				writer := msgp.NewWriter(&buf)
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(2))
				require.NoError(t, writer.WriteUint8(255))
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteString("name"))
				require.NoError(t, writer.Flush())
				return buf.Bytes()
			},
			wantErr: ErrUnsupportedOp,
		},
		{
			name: "test_type value is not array",
			ops:  []internal.Op{op.NewTestTypeMultiple([]string{"name"}, []string{"string"})},
			mutate: func(data []byte) []byte {
				encoded, err := New().Encode([]internal.Op{op.NewTest([]string{"name"}, "string")})
				if err != nil {
					panic(err)
				}
				for i, b := range encoded {
					if b == internal.OpTestCode {
						encoded[i] = internal.OpTestTypeCode
						break
					}
				}
				return encoded
			},
			wantErr: ErrInvalidTestTypeFormat,
		},
		{
			name: "in value is not array",
			ops:  []internal.Op{op.NewIn([]string{"role"}, []any{"admin"})},
			mutate: func(data []byte) []byte {
				encoded, err := New().Encode([]internal.Op{op.NewTest([]string{"role"}, "admin")})
				if err != nil {
					panic(err)
				}
				for i, b := range encoded {
					if b == internal.OpTestCode {
						encoded[i] = internal.OpInCode
						break
					}
				}
				return encoded
			},
			wantErr: ErrInvalidValueType,
		},
		{
			name: "extend properties are not object",
			ops:  []internal.Op{op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, false)},
			mutate: func(data []byte) []byte {
				encoded, err := New().Encode([]internal.Op{op.NewTest([]string{"profile"}, "Ada")})
				if err != nil {
					panic(err)
				}
				for i, b := range encoded {
					if b == internal.OpTestCode {
						encoded[i] = internal.OpExtendCode
						break
					}
				}
				return encoded
			},
			wantErr: ErrInvalidValueType,
		},
		{
			name: "merge properties are not object",
			ops:  []internal.Op{op.NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true})},
			mutate: func(data []byte) []byte {
				encoded, err := New().Encode([]internal.Op{op.NewSplit([]string{"nodes", "1"}, 1, "not-object")})
				if err != nil {
					panic(err)
				}
				for i, b := range encoded {
					if b == internal.OpSplitCode {
						encoded[i] = internal.OpMergeCode
						break
					}
				}
				return encoded
			},
			wantErr: ErrInvalidValueType,
		},
	}

	codec := New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := codec.Encode(tc.ops)
			require.NoError(t, err)

			decoded, err := codec.Decode(tc.mutate(data))
			require.Error(t, err)
			assert.Nil(t, decoded)
			assert.ErrorIs(t, err, tc.wantErr)
		})
	}
}

func TestCodecDecodeNormalizesNestedMessagePackMaps(t *testing.T) {
	t.Parallel()

	data := binaryFixture(t, func(writer *msgp.Writer) {
		require.NoError(t, writer.WriteArrayHeader(1))
		require.NoError(t, writer.WriteArrayHeader(3))
		require.NoError(t, writer.WriteUint8(internal.OpAddCode))
		writeBinaryPath(t, writer, "profile")
		require.NoError(t, writer.WriteMapHeader(2))
		require.NoError(t, writer.WriteString("name"))
		require.NoError(t, writer.WriteString("Ada"))
		require.NoError(t, writer.WriteString("meta"))
		require.NoError(t, writer.WriteMapHeader(2))
		require.NoError(t, writer.WriteString("active"))
		require.NoError(t, writer.WriteBool(true))
		require.NoError(t, writer.WriteString("tags"))
		require.NoError(t, writer.WriteArrayHeader(2))
		require.NoError(t, writer.WriteString("math"))
		require.NoError(t, writer.WriteString("logic"))
	})

	decoded, err := New().Decode(data)
	require.NoError(t, err)
	require.Len(t, decoded, 1)

	got, err := decoded[0].ToJSON()
	require.NoError(t, err)
	want := internal.Operation{
		Op:   "add",
		Path: "/profile",
		Value: map[string]any{
			"name": "Ada",
			"meta": map[string]any{
				"active": true,
				"tags":   []any{"math", "logic"},
			},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("decoded operation mismatch (-want +got):\n%s", diff)
	}
}

func TestCodecDecodeRejectsMalformedValuePayloads(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		write func(t *testing.T, writer *msgp.Writer)
	}{
		{
			name: "add value is incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpAddCode, 3)
			},
		},
		{
			name: "replace value is incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpReplaceCode, 3)
			},
		},
		{
			name: "test value is incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpTestCode, 3)
			},
		},
		{
			name: "test type value is incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpTestTypeCode, 3)
			},
		},
		{
			name: "in value is incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpInCode, 3)
			},
		},
		{
			name: "split properties are incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpSplitCode))
				writeBinaryPath(t, writer, "nodes", "0")
				require.NoError(t, writer.WriteFloat64(1))
				require.NoError(t, writer.WriteArrayHeader(1))
			},
		},
		{
			name: "extend properties are incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryMalformedValueOp(t, writer, internal.OpExtendCode, 4)
			},
		},
		{
			name: "merge properties are incomplete",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpMergeCode))
				writeBinaryPath(t, writer, "nodes", "1")
				require.NoError(t, writer.WriteFloat64(1))
				require.NoError(t, writer.WriteArrayHeader(1))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := New().Decode(binaryFixture(t, func(writer *msgp.Writer) {
				tc.write(t, writer)
			}))
			require.Error(t, err)
			assert.Nil(t, decoded)
		})
	}
}

func TestCodecDecodeRejectsMalformedFieldTypes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		write   func(t *testing.T, writer *msgp.Writer)
		wantErr error
	}{
		{
			name: "operations value is not array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteString("ops"))
			},
		},
		{
			name: "operation value is not array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteString("op"))
			},
		},
		{
			name: "operation code is not uint8",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(2))
				require.NoError(t, writer.WriteString("add"))
				writeBinaryPath(t, writer, "name")
			},
		},
		{
			name: "path is not array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpAddCode))
				require.NoError(t, writer.WriteString("name"))
				require.NoError(t, writer.WriteString("Ada"))
			},
		},
		{
			name: "path segment is not string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpAddCode))
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteInt64(1))
				require.NoError(t, writer.WriteString("Ada"))
			},
		},
		{
			name: "move from path is not array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpMoveCode))
				writeBinaryPath(t, writer, "target")
				require.NoError(t, writer.WriteString("source"))
			},
		},
		{
			name: "copy from path is not array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpCopyCode))
				writeBinaryPath(t, writer, "target")
				require.NoError(t, writer.WriteString("source"))
			},
		},
		{
			name: "test type array member is not string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpTestTypeCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteInt64(1))
			},
			wantErr: ErrInvalidTestTypeFormat,
		},
		{
			name: "less value is not float",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryNumericOp(t, writer, internal.OpLessCode, "ten")
			},
		},
		{
			name: "more value is not float",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				writeBinaryNumericOp(t, writer, internal.OpMoreCode, "five")
			},
		},
		{
			name: "contains value is not string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpContainsCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteInt64(1))
			},
		},
		{
			name: "type expected value is not string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpTypeCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteInt64(1))
			},
		},
		{
			name: "matches ignore case is not bool",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpMatchesCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteString("^A"))
				require.NoError(t, writer.WriteString("true"))
			},
		},
		{
			name: "test string position is not float",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpTestStringCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteString("Ad"))
				require.NoError(t, writer.WriteString("zero"))
			},
		},
		{
			name: "test string len not flag is not bool",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpTestStringLenCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteFloat64(3))
				require.NoError(t, writer.WriteString("false"))
			},
		},
		{
			name: "str ins string is not string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpStrInsCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteFloat64(1))
				require.NoError(t, writer.WriteInt64(1))
			},
		},
		{
			name: "str del length is not float",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpStrDelCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteFloat64(1))
				require.NoError(t, writer.WriteString("two"))
			},
		},
		{
			name: "extend delete null is not bool",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpExtendCode))
				writeBinaryPath(t, writer, "profile")
				require.NoError(t, writer.WriteMapHeader(1))
				require.NoError(t, writer.WriteString("name"))
				require.NoError(t, writer.WriteString("Ada"))
				require.NoError(t, writer.WriteString("false"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := New().Decode(binaryFixture(t, func(writer *msgp.Writer) {
				tc.write(t, writer)
			}))
			require.Error(t, err)
			assert.Nil(t, decoded)
			if tc.wantErr != nil {
				assert.ErrorIs(t, err, tc.wantErr)
			}
		})
	}
}

func binaryFixture(t *testing.T, write func(*msgp.Writer)) []byte {
	t.Helper()

	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	write(writer)
	require.NoError(t, writer.Flush())
	return buf.Bytes()
}

func writeBinaryPath(t *testing.T, writer *msgp.Writer, path ...string) {
	t.Helper()

	require.NoError(t, writer.WriteArrayHeader(uint32(len(path))))
	for _, segment := range path {
		require.NoError(t, writer.WriteString(segment))
	}
}

func writeBinaryNumericOp(t *testing.T, writer *msgp.Writer, code uint8, value string) {
	t.Helper()

	require.NoError(t, writer.WriteArrayHeader(1))
	require.NoError(t, writer.WriteArrayHeader(3))
	require.NoError(t, writer.WriteUint8(code))
	writeBinaryPath(t, writer, "score")
	require.NoError(t, writer.WriteString(value))
}

func writeBinaryMalformedValueOp(t *testing.T, writer *msgp.Writer, code uint8, opSize uint32) {
	t.Helper()

	require.NoError(t, writer.WriteArrayHeader(1))
	require.NoError(t, writer.WriteArrayHeader(opSize))
	require.NoError(t, writer.WriteUint8(code))
	writeBinaryPath(t, writer, "name")
	require.NoError(t, writer.WriteArrayHeader(1))
}

type unsupportedOp struct{}

func (unsupportedOp) Op() internal.OpType { return "unsupported" }
func (unsupportedOp) Code() int           { return 254 }
func (unsupportedOp) Path() []string      { return nil }
func (unsupportedOp) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (unsupportedOp) ToJSON() (internal.Operation, error)           { return internal.Operation{}, nil }
func (unsupportedOp) ToCompact() (internal.CompactOperation, error) { return nil, nil }
func (unsupportedOp) Validate() error                               { return nil }

func TestCodecEncodeRejectsUnsupportedOperation(t *testing.T) {
	t.Parallel()

	data, err := New().Encode([]internal.Op{unsupportedOp{}})
	require.Error(t, err)
	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrUnsupportedOp)
}
