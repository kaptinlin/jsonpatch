package compact

import (
	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpointer"
)

// Encode encodes operations into compact format using default options.
func Encode(ops []internal.Op, opts ...EncoderOption) ([]Op, error) {
	encoder := NewEncoder(opts...)
	return encoder.EncodeSlice(ops)
}

// EncodeJSON encodes operations into compact format JSON bytes.
func EncodeJSON(ops []internal.Op, opts ...EncoderOption) ([]byte, error) {
	compactOps, err := Encode(ops, opts...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(compactOps)
}

// opToCompact converts a single operation to compact format.
func opToCompact(op internal.Op, options EncoderOptions) (Op, error) {
	// Get the standard compact format from the operation
	compactOp, err := op.ToCompact()
	if err != nil {
		return nil, err
	}

	// Convert to our Op type and adjust opcode if needed
	result := make(Op, len(compactOp))
	copy(result, compactOp)

	// Convert opcode to string if requested
	if options.StringOpcode {
		result[0] = string(op.Op())
	}

	// Convert all path-like fields from []string to string format for JSON compatibility
	if len(result) > 1 {
		if pathSlice, ok := result[1].([]string); ok {
			result[1] = pathToString(pathSlice)
		}
	}

	// Handle 'from' path for move/copy operations (third element)
	if len(result) > 2 {
		if fromSlice, ok := result[2].([]string); ok {
			result[2] = pathToString(fromSlice)
		}
	}

	return result, nil
}

// pathToString converts a JSON pointer path to string representation.
func pathToString(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return jsonpointer.Format(path...)
}
