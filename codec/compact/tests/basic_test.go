package tests

import (
	"testing"

	"github.com/kaptinlin/jsonpatch/codec/compact"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

func TestBasicOperationsNumericCodes(t *testing.T) {
	tests := []struct {
		name     string
		op       internal.Op
		expected compact.CompactOp
	}{
		{
			name:     "add operation",
			op:       op.NewOpAddOperation([]string{"foo"}, "bar"),
			expected: compact.CompactOp{0, "/foo", "bar"},
		},
		{
			name:     "remove operation",
			op:       op.NewOpRemoveOperation([]string{"foo"}),
			expected: compact.CompactOp{1, "/foo"},
		},
		{
			name:     "replace operation",
			op:       op.NewOpReplaceOperation([]string{"foo"}, "new_value"),
			expected: compact.CompactOp{2, "/foo", "new_value"},
		},
		{
			name:     "move operation",
			op:       op.NewOpMoveOperation([]string{"foo"}, []string{"bar"}),
			expected: compact.CompactOp{4, "/foo", "/bar"},
		},
		{
			name:     "copy operation",
			op:       op.NewOpCopyOperation([]string{"foo"}, []string{"bar"}),
			expected: compact.CompactOp{3, "/foo", "/bar"},
		},
		{
			name:     "test operation",
			op:       op.NewOpTestOperation([]string{"foo"}, "expected"),
			expected: compact.CompactOp{5, "/foo", "expected"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test encoding
			encoder := compact.NewEncoder()
			encoded, err := encoder.Encode(tt.op)
			if err != nil {
				t.Fatalf("failed to encode operation: %v", err)
			}

			// Check encoded result
			if len(encoded) != len(tt.expected) {
				t.Errorf("encoded length mismatch: got %d, want %d", len(encoded), len(tt.expected))
			}

			// Check opcode
			if encoded[0] != tt.expected[0] {
				t.Errorf("opcode mismatch: got %v, want %v", encoded[0], tt.expected[0])
			}

			// Check path
			if encoded[1] != tt.expected[1] {
				t.Errorf("path mismatch: got %v, want %v", encoded[1], tt.expected[1])
			}

			// Test round-trip decoding
			decoder := compact.NewDecoder()
			decoded, err := decoder.Decode(encoded)
			if err != nil {
				t.Fatalf("failed to decode operation: %v", err)
			}

			// Check that decoded operation has the same type and path
			if decoded.Op() != tt.op.Op() {
				t.Errorf("operation type mismatch: got %v, want %v", decoded.Op(), tt.op.Op())
			}

			// Check path equality
			originalPath := tt.op.Path()
			decodedPath := decoded.Path()
			if len(originalPath) != len(decodedPath) {
				t.Errorf("path length mismatch: got %d, want %d", len(decodedPath), len(originalPath))
			} else {
				for i, segment := range originalPath {
					if decodedPath[i] != segment {
						t.Errorf("path segment %d mismatch: got %v, want %v", i, decodedPath[i], segment)
					}
				}
			}
		})
	}
}

func TestStringOpcodes(t *testing.T) {
	op := op.NewOpAddOperation([]string{"foo"}, "bar")

	// Test with string opcodes
	encoder := compact.NewEncoder(compact.WithStringOpcode(true))
	encoded, err := encoder.Encode(op)
	if err != nil {
		t.Fatalf("failed to encode with string opcodes: %v", err)
	}

	// Check that opcode is a string
	if opcode, ok := encoded[0].(string); !ok || opcode != "add" {
		t.Errorf("expected string opcode 'add', got %v", encoded[0])
	}

	// Test decoding
	decoder := compact.NewDecoder()
	decoded, err := decoder.Decode(encoded)
	if err != nil {
		t.Fatalf("failed to decode string opcode operation: %v", err)
	}

	if decoded.Op() != internal.OpAddType {
		t.Errorf("decoded operation type mismatch: got %v, want %v", decoded.Op(), internal.OpAddType)
	}
}

func TestSliceOperations(t *testing.T) {
	ops := []internal.Op{
		op.NewOpAddOperation([]string{"foo"}, "bar"),
		op.NewOpRemoveOperation([]string{"baz"}),
		op.NewOpReplaceOperation([]string{"qux"}, 42),
	}

	// Test encoding slice
	encoded, err := compact.Encode(ops)
	if err != nil {
		t.Fatalf("failed to encode operations slice: %v", err)
	}

	if len(encoded) != len(ops) {
		t.Errorf("encoded slice length mismatch: got %d, want %d", len(encoded), len(ops))
	}

	// Test decoding slice
	decoded, err := compact.Decode(encoded)
	if err != nil {
		t.Fatalf("failed to decode operations slice: %v", err)
	}

	if len(decoded) != len(ops) {
		t.Errorf("decoded slice length mismatch: got %d, want %d", len(decoded), len(ops))
	}

	for i, decodedOp := range decoded {
		if decodedOp.Op() != ops[i].Op() {
			t.Errorf("operation %d type mismatch: got %v, want %v", i, decodedOp.Op(), ops[i].Op())
		}
	}
}

func TestJSONMarshaling(t *testing.T) {
	ops := []internal.Op{
		op.NewOpAddOperation([]string{"foo"}, "bar"),
		op.NewOpRemoveOperation([]string{"baz"}),
	}

	// Test encoding to JSON
	jsonData, err := compact.EncodeJSON(ops)
	if err != nil {
		t.Fatalf("failed to encode to JSON: %v", err)
	}

	// Test decoding from JSON
	decoded, err := compact.DecodeJSON(jsonData)
	if err != nil {
		t.Fatalf("failed to decode from JSON: %v", err)
	}

	if len(decoded) != len(ops) {
		t.Errorf("JSON round-trip length mismatch: got %d, want %d", len(decoded), len(ops))
	}

	for i, decodedOp := range decoded {
		if decodedOp.Op() != ops[i].Op() {
			t.Errorf("JSON round-trip operation %d type mismatch: got %v, want %v", i, decodedOp.Op(), ops[i].Op())
		}
	}
}
