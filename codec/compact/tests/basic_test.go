package tests

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBasicOperationsNumericCodes(t *testing.T) {
	tests := []struct {
		name     string
		op       internal.Op
		expected compact.Op
	}{
		{
			name:     "add operation",
			op:       op.NewAdd([]string{"foo"}, "bar"),
			expected: compact.Op{0, "/foo", "bar"},
		},
		{
			name:     "remove operation",
			op:       op.NewRemove([]string{"foo"}),
			expected: compact.Op{1, "/foo"},
		},
		{
			name:     "replace operation",
			op:       op.NewReplace([]string{"foo"}, "new_value"),
			expected: compact.Op{2, "/foo", "new_value"},
		},
		{
			name:     "move operation",
			op:       op.NewMove([]string{"foo"}, []string{"bar"}),
			expected: compact.Op{4, "/foo", "/bar"},
		},
		{
			name:     "copy operation",
			op:       op.NewCopy([]string{"foo"}, []string{"bar"}),
			expected: compact.Op{3, "/foo", "/bar"},
		},
		{
			name:     "test operation",
			op:       op.NewTest([]string{"foo"}, "expected"),
			expected: compact.Op{5, "/foo", "expected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoder := compact.NewEncoder()
			encoded, err := encoder.Encode(tt.op)
			require.NoError(t, err, "encoding should not error")

			// Check encoded result
			assert.Equal(t, len(tt.expected), len(encoded), "encoded length should match expected")

			// Check opcode
			assert.Equal(t, tt.expected[0], encoded[0], "opcode should match")

			// Check path
			assert.Equal(t, tt.expected[1], encoded[1], "path should match")

			// Test round-trip decoding
			decoder := compact.NewDecoder()
			decoded, err := decoder.Decode(encoded)
			require.NoError(t, err, "decoding should not error")

			// Check that decoded operation has the same type and path
			assert.Equal(t, tt.op.Op(), decoded.Op(), "operation type should match")

			// Check path equality
			originalPath := tt.op.Path()
			decodedPath := decoded.Path()
			assert.Equal(t, len(originalPath), len(decodedPath), "path length should match")
			for i, segment := range originalPath {
				assert.Equal(t, segment, decodedPath[i], "path segment %d should match", i)
			}
		})
	}
}

func TestStringOpcodes(t *testing.T) {
	op := op.NewAdd([]string{"foo"}, "bar")

	// Test with string opcodes
	encoder := compact.NewEncoder(compact.WithStringOpcode(true))
	encoded, err := encoder.Encode(op)
	require.NoError(t, err, "encoding with string opcodes should not error")

	// Check that opcode is a string
	opcode, ok := encoded[0].(string)
	assert.True(t, ok, "opcode should be a string")
	assert.Equal(t, "add", opcode, "opcode should be 'add'")

	// Test decoding
	decoder := compact.NewDecoder()
	decoded, err := decoder.Decode(encoded)
	require.NoError(t, err, "decoding string opcode operation should not error")

	assert.Equal(t, internal.OpAddType, decoded.Op(), "decoded operation type should match")
}

func TestSliceOperations(t *testing.T) {
	ops := []internal.Op{
		op.NewAdd([]string{"foo"}, "bar"),
		op.NewRemove([]string{"baz"}),
		op.NewReplace([]string{"qux"}, 42),
	}

	// Test encoding slice
	encoded, err := compact.Encode(ops)
	require.NoError(t, err, "encoding operations slice should not error")

	assert.Equal(t, len(ops), len(encoded), "encoded slice length should match")

	// Test decoding slice
	decoded, err := compact.Decode(encoded)
	require.NoError(t, err, "decoding operations slice should not error")

	assert.Equal(t, len(ops), len(decoded), "decoded slice length should match")

	for i, decodedOp := range decoded {
		assert.Equal(t, ops[i].Op(), decodedOp.Op(), "operation %d type should match", i)
	}
}

func TestJSONMarshaling(t *testing.T) {
	ops := []internal.Op{
		op.NewAdd([]string{"foo"}, "bar"),
		op.NewRemove([]string{"baz"}),
	}

	// Test encoding to JSON
	jsonData, err := compact.EncodeJSON(ops)
	require.NoError(t, err, "encoding to JSON should not error")

	// Test decoding from JSON
	decoded, err := compact.DecodeJSON(jsonData)
	require.NoError(t, err, "decoding from JSON should not error")

	assert.Equal(t, len(ops), len(decoded), "JSON round-trip length should match")

	for i, decodedOp := range decoded {
		assert.Equal(t, ops[i].Op(), decodedOp.Op(), "JSON round-trip operation %d type should match", i)
	}
}
