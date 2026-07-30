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

type binaryOperationCase struct {
	name string
	op   internal.Op
}

func binaryOperationCases() []binaryOperationCase {
	return []binaryOperationCase{
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
		{name: "test string predicate", op: op.NewTestString([]string{"profile", "name"}, "da", 1, false)},
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
		{name: "and predicate", op: op.NewAnd([]string{"profile"},
			op.NewDefined([]string{"profile", "name"}),
			op.NewContains([]string{"profile", "role"}, "admin"),
		)},
		{name: "or predicate", op: op.NewOr([]string{"profile"},
			op.NewUndefined([]string{"profile", "deleted"}),
			op.NewContainsWithIgnoreCase([]string{"profile", "role"}, "ADMIN", true),
		)},
		{name: "not predicate", op: op.NewNot(
			[]string{"profile"},
			op.NewUndefined([]string{"profile", "deleted"}),
		)},
	}
}

func TestCodecRoundTripPreservesOperationJSON(t *testing.T) {
	t.Parallel()

	codec := New()
	for _, tc := range binaryOperationCases() {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := codec.Encode([]internal.Op{tc.op})
			require.NoError(t, err)
			require.NotEmpty(t, encoded)

			decoded, err := codec.Decode(encoded)
			require.NoError(t, err)
			require.Len(t, decoded, 1)

			want := operationToJSON(t, tc.op)
			got := operationToJSON(t, decoded[0])
			if diff := cmp.Diff(want, got); diff != "" {
				t.Errorf("round-tripped operation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCodecEncodeUsesCanonicalIntegerTrueMarkers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		present     internal.Op
		absent      internal.Op
		presentSize int
	}{
		{name: "test not", present: op.NewTestWithNot([]string{"status"}, "ready", true), absent: op.NewTest([]string{"status"}, "ready"), presentSize: 4},
		{name: "contains ignore case", present: op.NewContainsWithIgnoreCase([]string{"name"}, "ada", true), absent: op.NewContains([]string{"name"}, "Ada"), presentSize: 4},
		{name: "starts ignore case", present: op.NewStartsWithIgnoreCase([]string{"name"}, "a", true), absent: op.NewStarts([]string{"name"}, "A"), presentSize: 4},
		{name: "ends ignore case", present: op.NewEndsWithIgnoreCase([]string{"name"}, "a", true), absent: op.NewEnds([]string{"name"}, "a"), presentSize: 4},
		{name: "matches ignore case", present: op.NewMatches([]string{"name"}, "^a", true, nil), absent: op.NewMatches([]string{"name"}, "^A", false, nil), presentSize: 4},
		{name: "test string not", present: op.NewTestString([]string{"name"}, "ad", 0, true), absent: op.NewTestString([]string{"name"}, "Ad", 0, false), presentSize: 5},
		{name: "test string length not", present: op.NewTestStringLenWithNot([]string{"name"}, 3, true), absent: op.NewTestStringLen([]string{"name"}, 3), presentSize: 4},
		{name: "extend delete null", present: op.NewExtend([]string{"profile"}, map[string]any{"name": nil}, true), absent: op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, false), presentSize: 4},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			present := decodeSingleMessagePackOperation(t, tc.present)
			require.Len(t, present, tc.presentSize)
			assertMessagePackIntegerOne(t, present[len(present)-1])

			absent := decodeSingleMessagePackOperation(t, tc.absent)
			require.Len(t, absent, tc.presentSize-1)
		})
	}
}

func decodeSingleMessagePackOperation(t *testing.T, operation internal.Op) []any {
	t.Helper()

	encoded, err := New().Encode([]internal.Op{operation})
	require.NoError(t, err)
	decoded, err := msgp.NewReader(bytes.NewReader(encoded)).ReadIntf()
	require.NoError(t, err)
	patch, ok := decoded.([]any)
	require.True(t, ok, "decoded patch has type %T", decoded)
	require.Len(t, patch, 1)
	frame, ok := patch[0].([]any)
	require.True(t, ok, "decoded operation has type %T", patch[0])
	return frame
}

func assertMessagePackIntegerOne(t *testing.T, value any) {
	t.Helper()

	switch marker := value.(type) {
	case int64:
		assert.Equal(t, int64(1), marker)
	case uint64:
		assert.Equal(t, uint64(1), marker)
	default:
		t.Errorf("true marker has type %T, want MessagePack integer", value)
	}
}

func TestCodecEncodeCanonicalOptionalFields(t *testing.T) {
	t.Parallel()

	encoded, err := New().Encode([]internal.Op{
		op.NewTest([]string{"profile", "name"}, "Ada"),
		op.NewTestWithNot([]string{"profile", "name"}, "Grace", true),
		op.NewMatches([]string{"profile", "name"}, "^A", false, nil),
		op.NewContainsWithIgnoreCase([]string{"profile", "name"}, "ada", true),
		op.NewTestString([]string{"profile", "name"}, "da", 1, true),
		op.NewTestStringLenWithNot([]string{"profile", "name"}, 3, false),
		op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, true),
	})
	require.NoError(t, err)

	reader := msgp.NewReader(bytes.NewReader(encoded))
	count, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(7), count)

	readBinaryOpHeader(t, reader, 3, internal.OpTestCode, []string{"profile", "name"})
	value, err := reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "Ada", value)

	readBinaryOpHeader(t, reader, 4, internal.OpTestCode, []string{"profile", "name"})
	value, err = reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "Grace", value)
	not, err := reader.ReadUint8()
	require.NoError(t, err)
	assert.Equal(t, uint8(1), not)

	readBinaryOpHeader(t, reader, 3, internal.OpMatchesCode, []string{"profile", "name"})
	pattern, err := reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "^A", pattern)

	readBinaryOpHeader(t, reader, 4, internal.OpContainsCode, []string{"profile", "name"})
	needle, err := reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "ada", needle)
	ignoreCase, err := reader.ReadUint8()
	require.NoError(t, err)
	assert.Equal(t, uint8(1), ignoreCase)

	readBinaryOpHeader(t, reader, 5, internal.OpTestStringCode, []string{"profile", "name"})
	pos, err := reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(1), pos)
	str, err := reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "da", str)
	not, err = reader.ReadUint8()
	require.NoError(t, err)
	assert.Equal(t, uint8(1), not)

	readBinaryOpHeader(t, reader, 3, internal.OpTestStringLenCode, []string{"profile", "name"})
	length, err := reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(3), length)

	readBinaryOpHeader(t, reader, 4, internal.OpExtendCode, []string{"profile"})
	props, err := reader.ReadIntf()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"name": "Ada"}, props)
	deleteNull, err := reader.ReadUint8()
	require.NoError(t, err)
	assert.Equal(t, uint8(1), deleteNull)
}

func TestCodecEncodeCanonicalWireBytesGolden(t *testing.T) {
	t.Parallel()

	encoded, err := New().Encode([]internal.Op{
		op.NewTest([]string{"flag"}, true),
		op.NewTestWithNot([]string{"flag"}, false, true),
		op.NewContains([]string{"name"}, "Ada"),
		op.NewContainsWithIgnoreCase([]string{"name"}, "ada", true),
		op.NewStarts([]string{"name"}, "A"),
		op.NewStartsWithIgnoreCase([]string{"name"}, "a", true),
		op.NewEnds([]string{"name"}, "a"),
		op.NewEndsWithIgnoreCase([]string{"name"}, "A", true),
		op.NewMatches([]string{"name"}, "^A", false, nil),
		op.NewMatches([]string{"name"}, "^a", true, nil),
		op.NewTestString([]string{"name"}, "Ad", 0, false),
		op.NewTestString([]string{"name"}, "ad", 0, true),
		op.NewTestStringLen([]string{"name"}, 3),
		op.NewTestStringLenWithNot([]string{"name"}, 4, true),
		op.NewSplit([]string{"nodes", "0"}, 1, nil),
		op.NewSplit([]string{"nodes", "0"}, 1, map[string]any{"kind": "paragraph"}),
		op.NewMerge([]string{"nodes", "1"}, 1, nil),
		op.NewMerge([]string{"nodes", "1"}, 1, map[string]any{"merged": true}),
		op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, false),
		op.NewExtend([]string{"profile"}, map[string]any{"name": "Ada"}, true),
		op.NewAnd([]string{"profile"},
			op.NewDefined([]string{"profile", "name"}),
			op.NewOr([]string{"profile", "flags"},
				op.NewUndefined([]string{"profile", "flags", "deleted"}),
				op.NewContains([]string{"profile", "flags", "role"}, "admin"),
			),
		),
	})
	require.NoError(t, err)

	want := binaryFixture(t, func(writer *msgp.Writer) {
		require.NoError(t, writer.WriteArrayHeader(21))

		writeBinaryHeader(t, writer, 3, internal.OpTestCode, "flag")
		require.NoError(t, writer.WriteBool(true))
		writeBinaryHeader(t, writer, 4, internal.OpTestCode, "flag")
		require.NoError(t, writer.WriteBool(false))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpContainsCode, "name")
		require.NoError(t, writer.WriteString("Ada"))
		writeBinaryHeader(t, writer, 4, internal.OpContainsCode, "name")
		require.NoError(t, writer.WriteString("ada"))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpStartsCode, "name")
		require.NoError(t, writer.WriteString("A"))
		writeBinaryHeader(t, writer, 4, internal.OpStartsCode, "name")
		require.NoError(t, writer.WriteString("a"))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpEndsCode, "name")
		require.NoError(t, writer.WriteString("a"))
		writeBinaryHeader(t, writer, 4, internal.OpEndsCode, "name")
		require.NoError(t, writer.WriteString("A"))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpMatchesCode, "name")
		require.NoError(t, writer.WriteString("^A"))
		writeBinaryHeader(t, writer, 4, internal.OpMatchesCode, "name")
		require.NoError(t, writer.WriteString("^a"))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 4, internal.OpTestStringCode, "name")
		require.NoError(t, writer.WriteInt64(0))
		require.NoError(t, writer.WriteString("Ad"))
		writeBinaryHeader(t, writer, 5, internal.OpTestStringCode, "name")
		require.NoError(t, writer.WriteInt64(0))
		require.NoError(t, writer.WriteString("ad"))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpTestStringLenCode, "name")
		require.NoError(t, writer.WriteInt64(3))
		writeBinaryHeader(t, writer, 4, internal.OpTestStringLenCode, "name")
		require.NoError(t, writer.WriteInt64(4))
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpSplitCode, "nodes", "0")
		require.NoError(t, writer.WriteInt64(1))
		writeBinaryHeader(t, writer, 4, internal.OpSplitCode, "nodes", "0")
		require.NoError(t, writer.WriteInt64(1))
		writeBinaryStringMap(t, writer, "kind", "paragraph")

		writeBinaryHeader(t, writer, 3, internal.OpMergeCode, "nodes", "1")
		require.NoError(t, writer.WriteInt64(1))
		writeBinaryHeader(t, writer, 4, internal.OpMergeCode, "nodes", "1")
		require.NoError(t, writer.WriteInt64(1))
		require.NoError(t, writer.WriteMapHeader(1))
		require.NoError(t, writer.WriteString("merged"))
		require.NoError(t, writer.WriteBool(true))

		writeBinaryHeader(t, writer, 3, internal.OpExtendCode, "profile")
		writeBinaryStringMap(t, writer, "name", "Ada")
		writeBinaryHeader(t, writer, 4, internal.OpExtendCode, "profile")
		writeBinaryStringMap(t, writer, "name", "Ada")
		require.NoError(t, writer.WriteUint8(1))

		writeBinaryHeader(t, writer, 3, internal.OpAndCode, "profile")
		require.NoError(t, writer.WriteArrayHeader(2))
		writeBinaryHeader(t, writer, 2, internal.OpDefinedCode, "name")
		writeBinaryHeader(t, writer, 3, internal.OpOrCode, "flags")
		require.NoError(t, writer.WriteArrayHeader(2))
		writeBinaryHeader(t, writer, 2, internal.OpUndefinedCode, "deleted")
		writeBinaryHeader(t, writer, 3, internal.OpContainsCode, "role")
		require.NoError(t, writer.WriteString("admin"))
	})

	if diff := cmp.Diff(want, encoded); diff != "" {
		t.Errorf("canonical binary bytes mismatch (-want +got):\n%s", diff)
	}
}

func TestCodecEncodeCompositeUsesParentRelativePaths(t *testing.T) {
	t.Parallel()

	encoded, err := New().Encode([]internal.Op{
		op.NewAnd([]string{"profile"},
			op.NewDefined([]string{"profile", "name"}),
			op.NewNot(
				[]string{"profile", "meta"},
				op.NewUndefined([]string{"profile", "meta", "deleted"}),
			),
		),
	})
	require.NoError(t, err)

	reader := msgp.NewReader(bytes.NewReader(encoded))
	count, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(1), count)

	readBinaryOpHeader(t, reader, 3, internal.OpAndCode, []string{"profile"})
	childCount, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(2), childCount)

	readBinaryOpHeader(t, reader, 2, internal.OpDefinedCode, []string{"name"})

	readBinaryOpHeader(t, reader, 3, internal.OpNotCode, []string{"meta"})
	grandchildCount, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(1), grandchildCount)
	readBinaryOpHeader(t, reader, 2, internal.OpUndefinedCode, []string{"deleted"})
}

func TestCodecDecodeCompositePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		data []byte
		want internal.Operation
	}{
		{
			name: "and merges relative child paths",
			data: binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAndCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(2))
				writeBinaryHeader(t, writer, 2, internal.OpDefinedCode, "name")
				writeBinaryHeader(t, writer, 3, internal.OpContainsCode, "role")
				require.NoError(t, writer.WriteString("admin"))
			}),
			want: internal.Operation{Op: "and", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/name"},
				{Op: "contains", Path: "/role", Value: "admin"},
			}},
		},
		{
			name: "or treats repeated parent segment as relative",
			data: binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpOrCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 2, internal.OpDefinedCode, "profile")
			}),
			want: internal.Operation{Op: "or", Path: "/profile", Apply: []internal.Operation{
				{Op: "defined", Path: "/profile"},
			}},
		},
		{
			name: "not uses parent path for empty child path",
			data: binaryFixture(t, func(writer *msgp.Writer) {
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpNotCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 2, internal.OpUndefinedCode)
			}),
			want: internal.Operation{Op: "not", Path: "/profile", Apply: []internal.Operation{
				{Op: "undefined", Path: ""},
			}},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			decoded, err := New().Decode(tc.data)
			require.NoError(t, err)
			require.Len(t, decoded, 1)

			got := operationToJSON(t, decoded[0])
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("decoded operation mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCodecDecodeRejectsInvalidCompositePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		write   func(t *testing.T, writer *msgp.Writer)
		wantErr error
	}{
		{
			name: "predicate list is not an array",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAndCode, "profile")
				require.NoError(t, writer.WriteString("not-array"))
			},
		},
		{
			name: "child operation is malformed",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAndCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(3))
				require.NoError(t, writer.WriteUint8(internal.OpDefinedCode))
				require.NoError(t, writer.WriteString("name"))
			},
		},
		{
			name: "child operation is not a predicate",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAndCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpAddCode, "name")
				require.NoError(t, writer.WriteString("Ada"))
			},
			wantErr: ErrInvalidPredicate,
		},
		{
			name: "not contains no predicate",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpNotCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(0))
			},
			wantErr: ErrNotSinglePredicate,
		},
		{
			name: "not contains more than one predicate",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				writeBinaryHeader(t, writer, 3, internal.OpNotCode, "profile")
				require.NoError(t, writer.WriteArrayHeader(2))
				writeBinaryHeader(t, writer, 2, internal.OpDefinedCode, "name")
				writeBinaryHeader(t, writer, 2, internal.OpUndefinedCode, "deleted")
			},
			wantErr: ErrNotSinglePredicate,
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

func TestCodecEncodeRejectsInvalidCompositePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		op      internal.Op
		wantErr error
	}{
		{
			name:    "or rejects child outside parent path",
			op:      op.NewOr([]string{"profile"}, op.NewDefined([]string{"settings", "enabled"})),
			wantErr: op.ErrPredicatePathOutsideParent,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			data, err := New().Encode([]internal.Op{tc.op})
			require.Error(t, err)
			assert.Nil(t, data)
			assert.ErrorIs(t, err, tc.wantErr)
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

	got := operationToJSON(t, decoded[0])
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

func operationToJSON(t *testing.T, operation internal.Op) internal.Operation {
	t.Helper()

	jsonOp, ok := operation.(internal.JSONOp)
	require.True(t, ok, "operation %T should encode to JSON", operation)
	result, err := jsonOp.ToJSON()
	require.NoError(t, err)
	return result
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
			name: "matches ignore case marker is string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpMatchesCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteString("^A"))
				require.NoError(t, writer.WriteString("true"))
			},
			wantErr: ErrInvalidTrueMarker,
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
			name: "test string len marker is string",
			write: func(t *testing.T, writer *msgp.Writer) {
				t.Helper()
				require.NoError(t, writer.WriteArrayHeader(1))
				require.NoError(t, writer.WriteArrayHeader(4))
				require.NoError(t, writer.WriteUint8(internal.OpTestStringLenCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteFloat64(3))
				require.NoError(t, writer.WriteString("false"))
			},
			wantErr: ErrInvalidTrueMarker,
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
				require.NoError(t, writer.WriteArrayHeader(5))
				require.NoError(t, writer.WriteUint8(internal.OpStrDelCode))
				writeBinaryPath(t, writer, "name")
				require.NoError(t, writer.WriteFloat64(1))
				require.NoError(t, writer.WriteUint8(0))
				require.NoError(t, writer.WriteString("two"))
			},
		},
		{
			name: "extend delete null marker is string",
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
			wantErr: ErrInvalidTrueMarker,
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

func TestCodecPreservesStrDelVariants(t *testing.T) {
	t.Parallel()

	operations := []internal.Op{
		op.NewStrDelWithStr([]string{"text"}, 1, "da"),
		op.NewStrDelWithStr([]string{"text"}, 2, ""),
		op.NewStrDel([]string{"text"}, 3, 4),
	}
	codec := New()
	encoded, err := codec.Encode(operations)
	require.NoError(t, err)

	reader := msgp.NewReader(bytes.NewReader(encoded))
	count, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, uint32(3), count)

	readBinaryOpHeader(t, reader, 4, internal.OpStrDelCode, []string{"text"})
	pos, err := reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(1), pos)
	str, err := reader.ReadString()
	require.NoError(t, err)
	assert.Equal(t, "da", str)

	readBinaryOpHeader(t, reader, 4, internal.OpStrDelCode, []string{"text"})
	pos, err = reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(2), pos)
	str, err = reader.ReadString()
	require.NoError(t, err)
	assert.Empty(t, str)

	readBinaryOpHeader(t, reader, 5, internal.OpStrDelCode, []string{"text"})
	pos, err = reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(3), pos)
	discriminator, err := reader.ReadUint8()
	require.NoError(t, err)
	assert.Zero(t, discriminator)
	length, err := reader.ReadInt64()
	require.NoError(t, err)
	assert.Equal(t, int64(4), length)

	decoded, err := codec.Decode(encoded)
	require.NoError(t, err)
	require.Len(t, decoded, 3)
	assertStrDelVariant(t, decoded[0], true, "da", 2)
	assertStrDelVariant(t, decoded[1], true, "", 0)
	assertStrDelVariant(t, decoded[2], false, "", 4)
}

func assertStrDelVariant(t *testing.T, operation internal.Op, wantHasStr bool, wantStr string, wantLen int) {
	t.Helper()

	strDel, ok := operation.(*op.StrDelOperation)
	require.True(t, ok)
	assert.Equal(t, wantHasStr, strDel.HasStr)
	assert.Equal(t, wantStr, strDel.Str)
	assert.Equal(t, wantLen, strDel.Len)
}

func writeBinaryPath(t *testing.T, writer *msgp.Writer, path ...string) {
	t.Helper()

	require.NoError(t, writer.WriteArrayHeader(uint32(len(path))))
	for _, segment := range path {
		require.NoError(t, writer.WriteString(segment))
	}
}

func writeBinaryHeader(t *testing.T, writer *msgp.Writer, size uint32, code int, path ...string) {
	t.Helper()

	require.NoError(t, writer.WriteArrayHeader(size))
	require.NoError(t, writer.WriteUint8(uint8(code)))
	writeBinaryPath(t, writer, path...)
}

func writeBinaryStringMap(t *testing.T, writer *msgp.Writer, key, value string) {
	t.Helper()

	require.NoError(t, writer.WriteMapHeader(1))
	require.NoError(t, writer.WriteString(key))
	require.NoError(t, writer.WriteString(value))
}

func readBinaryOpHeader(t *testing.T, reader *msgp.Reader, wantSize uint32, wantCode int, wantPath []string) {
	t.Helper()

	size, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	assert.Equal(t, wantSize, size)

	code, err := reader.ReadUint8()
	require.NoError(t, err)
	assert.Equal(t, wantCode, int(code))

	path := readBinaryPath(t, reader)
	assert.Equal(t, wantPath, path)
}

func readBinaryPath(t *testing.T, reader *msgp.Reader) []string {
	t.Helper()

	size, err := reader.ReadArrayHeader()
	require.NoError(t, err)
	path := make([]string, int(size))
	for i := range path {
		segment, err := reader.ReadString()
		require.NoError(t, err)
		path[i] = segment
	}
	return path
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
func (unsupportedOp) Path() []string      { return nil }
func (unsupportedOp) Apply(doc any) (internal.OpResult[any], error) {
	return internal.OpResult[any]{Doc: doc}, nil
}
func (unsupportedOp) Validate() error             { return nil }
func (unsupportedOp) Clone() (internal.Op, error) { return unsupportedOp{}, nil }

func TestCodecEncodeRejectsUnsupportedOperation(t *testing.T) {
	t.Parallel()

	data, err := New().Encode([]internal.Op{unsupportedOp{}})
	require.Error(t, err)
	assert.Nil(t, data)
	assert.ErrorIs(t, err, ErrUnsupportedOp)
}
