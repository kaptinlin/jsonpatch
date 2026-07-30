package binary

import (
	"bytes"
	"io"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestCodecUsesCanonicalCompactOperationTree(t *testing.T) {
	t.Parallel()

	operations := binaryRepresentationOperations()
	compactTree, err := compact.Encode(operations)
	require.NoError(t, err)

	t.Run("encode", func(t *testing.T) {
		t.Parallel()

		encoded, err := New().Encode(operations)
		require.NoError(t, err)
		binaryTree := decodeMessagePackTree(t, encoded)
		if diff := cmp.Diff(normalizeRepresentation(compactTree), normalizeRepresentation(binaryTree)); diff != "" {
			t.Errorf("binary tree differs from compact tree (-compact +binary):\n%s", diff)
		}
	})

	t.Run("decode", func(t *testing.T) {
		t.Parallel()

		encoded := encodeMessagePackTree(t, messagePackValue(compactTree))
		decoded, err := New().Decode(encoded)
		require.NoError(t, err)
		restored, err := compact.Encode(decoded)
		require.NoError(t, err)
		if diff := cmp.Diff(normalizeRepresentation(compactTree), normalizeRepresentation(restored)); diff != "" {
			t.Errorf("restored compact tree mismatch (-want +got):\n%s", diff)
		}
	})
}

func binaryRepresentationOperations() []internal.Op {
	cases := binaryOperationCases()
	operations := make([]internal.Op, 0, len(cases)+10)
	for _, tc := range cases {
		operations = append(operations, tc.op)
	}
	operations = append(operations,
		op.NewRemoveWithOldValue([]string{"nullable"}, nil),
		op.NewReplaceWithOldValue([]string{"nullable"}, "new", nil),
		op.NewTestWithNot([]string{"status"}, "ready", true),
		op.NewStartsWithIgnoreCase([]string{"name"}, "a", true),
		op.NewEndsWithIgnoreCase([]string{"name"}, "a", true),
		op.NewTestString([]string{"name"}, "ad", 0, true),
		op.NewStrDelWithStr([]string{"text"}, 0, ""),
		op.NewAnd([]string{"profile"},
			op.NewNot([]string{"profile", "flags"},
				op.NewContainsWithIgnoreCase([]string{"profile", "flags", "role"}, "ADMIN", true),
			),
		),
	)
	if strconv.IntSize == 64 {
		large := int64(1<<53 + 1)
		operations = append(operations,
			op.NewStrIns([]string{"text"}, int(large), "x"),
			op.NewMerge([]string{"nodes", "1"}, -7, nil),
		)
	}
	return operations
}

func decodeMessagePackTree(t *testing.T, data []byte) any {
	t.Helper()

	reader := msgp.NewReader(bytes.NewReader(data))
	decoded, err := reader.ReadIntf()
	require.NoError(t, err)
	_, err = reader.NextType()
	require.ErrorIs(t, err, io.EOF)
	return decoded
}

func encodeMessagePackTree(t *testing.T, tree any) []byte {
	t.Helper()

	var buffer bytes.Buffer
	writer := msgp.NewWriter(&buffer)
	require.NoError(t, writer.WriteIntf(tree))
	require.NoError(t, writer.Flush())
	return buffer.Bytes()
}

func messagePackValue(value any) any {
	switch typed := value.(type) {
	case []compact.Op:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = messagePackValue(child)
		}
		return result
	case compact.Op:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = messagePackValue(child)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = messagePackValue(child)
		}
		return result
	case []string:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = child
		}
		return result
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, child := range typed {
			result[key] = messagePackValue(child)
		}
		return result
	default:
		return value
	}
}

type normalizedInteger struct {
	Negative  bool
	Magnitude uint64
}

func normalizeRepresentation(value any) any {
	switch typed := value.(type) {
	case []compact.Op:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = normalizeRepresentation(child)
		}
		return result
	case compact.Op:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = normalizeRepresentation(child)
		}
		return result
	case []any:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = normalizeRepresentation(child)
		}
		return result
	case []string:
		result := make([]any, len(typed))
		for i, child := range typed {
			result[i] = child
		}
		return result
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, child := range typed {
			result[key] = normalizeRepresentation(child)
		}
		return result
	case int:
		return normalizeSignedInteger(int64(typed))
	case int8:
		return normalizeSignedInteger(int64(typed))
	case int16:
		return normalizeSignedInteger(int64(typed))
	case int32:
		return normalizeSignedInteger(int64(typed))
	case int64:
		return normalizeSignedInteger(typed)
	case uint:
		return normalizedInteger{Magnitude: uint64(typed)}
	case uint8:
		return normalizedInteger{Magnitude: uint64(typed)}
	case uint16:
		return normalizedInteger{Magnitude: uint64(typed)}
	case uint32:
		return normalizedInteger{Magnitude: uint64(typed)}
	case uint64:
		return normalizedInteger{Magnitude: typed}
	default:
		return value
	}
}

func normalizeSignedInteger(value int64) normalizedInteger {
	if value >= 0 {
		return normalizedInteger{Magnitude: uint64(value)}
	}
	return normalizedInteger{Negative: true, Magnitude: uint64(-(value + 1)) + 1}
}
