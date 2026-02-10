package tests

import (
	"testing"

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
			if err != nil {
				t.Fatalf("encoding should not error: %v", err)
			}

			// Check encoded result
			if got, want := len(encoded), len(tt.expected); got != want {
				t.Errorf("encoded length = %d, want %d", got, want)
			}

			// Check opcode
			if got, want := encoded[0], tt.expected[0]; got != want {
				t.Errorf("opcode = %v, want %v", got, want)
			}

			// Check path
			if got, want := encoded[1], tt.expected[1]; got != want {
				t.Errorf("path = %v, want %v", got, want)
			}

			// Test round-trip decoding
			decoder := compact.NewDecoder()
			decoded, err := decoder.Decode(encoded)
			if err != nil {
				t.Fatalf("decoding should not error: %v", err)
			}

			// Check that decoded operation has the same type and path
			if got, want := decoded.Op(), tt.op.Op(); got != want {
				t.Errorf("operation type = %v, want %v", got, want)
			}

			// Check path equality
			originalPath := tt.op.Path()
			decodedPath := decoded.Path()
			if got, want := len(decodedPath), len(originalPath); got != want {
				t.Errorf("path length = %d, want %d", got, want)
			}
			for i, segment := range originalPath {
				if got := decodedPath[i]; got != segment {
					t.Errorf("path segment %d = %v, want %v", i, got, segment)
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
	if err != nil {
		t.Fatalf("encoding with string opcodes should not error: %v", err)
	}

	// Check that opcode is a string
	opcode, ok := encoded[0].(string)
	if !ok {
		t.Fatal("opcode should be a string")
	}
	if opcode != "add" {
		t.Errorf("opcode = %v, want %v", opcode, "add")
	}

	// Test decoding
	decoder := compact.NewDecoder()
	decoded, err := decoder.Decode(encoded)
	if err != nil {
		t.Fatalf("decoding string opcode operation should not error: %v", err)
	}

	if got, want := decoded.Op(), internal.OpAddType; got != want {
		t.Errorf("decoded operation type = %v, want %v", got, want)
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
	if err != nil {
		t.Fatalf("encoding operations slice should not error: %v", err)
	}

	if got, want := len(encoded), len(ops); got != want {
		t.Errorf("encoded slice length = %d, want %d", got, want)
	}

	// Test decoding slice
	decoded, err := compact.Decode(encoded)
	if err != nil {
		t.Fatalf("decoding operations slice should not error: %v", err)
	}

	if got, want := len(decoded), len(ops); got != want {
		t.Errorf("decoded slice length = %d, want %d", got, want)
	}

	for i, decodedOp := range decoded {
		if got, want := decodedOp.Op(), ops[i].Op(); got != want {
			t.Errorf("operation %d type = %v, want %v", i, got, want)
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
	if err != nil {
		t.Fatalf("encoding to JSON should not error: %v", err)
	}

	// Test decoding from JSON
	decoded, err := compact.DecodeJSON(jsonData)
	if err != nil {
		t.Fatalf("decoding from JSON should not error: %v", err)
	}

	if got, want := len(decoded), len(ops); got != want {
		t.Errorf("JSON round-trip length = %d, want %d", got, want)
	}

	for i, decodedOp := range decoded {
		if got, want := decodedOp.Op(), ops[i].Op(); got != want {
			t.Errorf("JSON round-trip operation %d type = %v, want %v", i, got, want)
		}
	}
}
