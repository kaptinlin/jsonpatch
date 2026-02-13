package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestBasicOperationsNumericCodes(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			// Test encoding
			encoder := compact.NewEncoder()
			encoded, err := encoder.Encode(tt.op)
			require.NoError(t, err)

			// Check encoded result
			if got, want := len(encoded), len(tt.expected); got != want {
				assert.Equal(t, want, got, "encoded length")
			}

			// Check opcode
			if got, want := encoded[0], tt.expected[0]; got != want {
				assert.Equal(t, want, got, "opcode")
			}

			// Check path
			if got, want := encoded[1], tt.expected[1]; got != want {
				assert.Equal(t, want, got, "path")
			}

			// Test round-trip decoding
			decoder := compact.NewDecoder()
			decoded, err := decoder.Decode(encoded)
			require.NoError(t, err)

			// Check that decoded operation has the same type and path
			if got, want := decoded.Op(), tt.op.Op(); got != want {
				assert.Equal(t, want, got, "operation type")
			}

			// Check path equality
			originalPath := tt.op.Path()
			decodedPath := decoded.Path()
			if got, want := len(decodedPath), len(originalPath); got != want {
				assert.Equal(t, want, got, "path length")
			}
			for i, segment := range originalPath {
				if got := decodedPath[i]; got != segment {
					assert.Equal(t, got, segment, i, "path segment %d")
				}
			}
		})
	}
}

func TestStringOpcodes(t *testing.T) {
	t.Parallel()
	addOp := op.NewAdd([]string{"foo"}, "bar")

	// Test with string opcodes
	encoder := compact.NewEncoder(compact.WithStringOpcode(true))
	encoded, err := encoder.Encode(addOp)
	require.NoError(t, err)

	// Check that opcode is a string
	opcode, ok := encoded[0].(string)
	if !ok {
		t.Fatal("opcode should be a string")
	}
	assert.Equal(t, "add", opcode, "opcode")

	// Test decoding
	decoder := compact.NewDecoder()
	decoded, err := decoder.Decode(encoded)
	require.NoError(t, err)

	if got, want := decoded.Op(), internal.OpAddType; got != want {
		assert.Equal(t, want, got, "decoded operation type")
	}
}

func TestSliceOperations(t *testing.T) {
	t.Parallel()
	ops := []internal.Op{
		op.NewAdd([]string{"foo"}, "bar"),
		op.NewRemove([]string{"baz"}),
		op.NewReplace([]string{"qux"}, 42),
	}

	// Test encoding slice
	encoded, err := compact.Encode(ops)
	require.NoError(t, err)

	if got, want := len(encoded), len(ops); got != want {
		assert.Equal(t, want, got, "encoded slice length")
	}

	// Test decoding slice
	decoded, err := compact.Decode(encoded)
	require.NoError(t, err)

	if got, want := len(decoded), len(ops); got != want {
		assert.Equal(t, want, got, "decoded slice length")
	}

	for i, decodedOp := range decoded {
		if got, want := decodedOp.Op(), ops[i].Op(); got != want {
			assert.Equal(t, got, want, i, "operation %d type")
		}
	}
}

func TestJSONMarshaling(t *testing.T) {
	t.Parallel()
	ops := []internal.Op{
		op.NewAdd([]string{"foo"}, "bar"),
		op.NewRemove([]string{"baz"}),
	}

	// Test encoding to JSON
	jsonData, err := compact.EncodeJSON(ops)
	require.NoError(t, err)

	// Test decoding from JSON
	decoded, err := compact.DecodeJSON(jsonData)
	require.NoError(t, err)

	if got, want := len(decoded), len(ops); got != want {
		assert.Equal(t, want, got, "JSON round-trip length")
	}

	for i, decodedOp := range decoded {
		if got, want := decodedOp.Op(), ops[i].Op(); got != want {
			assert.Equal(t, got, want, i, "JSON round-trip operation %d type")
		}
	}
}
