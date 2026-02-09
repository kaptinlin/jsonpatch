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

// encodeOp converts a single operation to compact format.
func encodeOp(o internal.Op, options EncoderOptions) (Op, error) {
	compactOp, err := o.ToCompact()
	if err != nil {
		return nil, err
	}

	result := make(Op, len(compactOp))
	copy(result, compactOp)

	if options.StringOpcode {
		result[0] = string(o.Op())
	}

	// Convert path fields from []string to JSON Pointer string format.
	if len(result) > 1 {
		if pathSlice, ok := result[1].([]string); ok {
			result[1] = formatPath(pathSlice)
		}
	}
	if len(result) > 2 {
		if fromSlice, ok := result[2].([]string); ok {
			result[2] = formatPath(fromSlice)
		}
	}

	return result, nil
}

// formatPath converts a JSON pointer path to string representation.
func formatPath(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return jsonpointer.Format(path...)
}
