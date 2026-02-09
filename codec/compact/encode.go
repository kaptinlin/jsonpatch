package compact

import (
	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpointer"
)

// Encoder encodes operations into compact format.
type Encoder struct {
	options EncoderOptions
}

// NewEncoder creates a new compact encoder with the given options.
func NewEncoder(opts ...EncoderOption) *Encoder {
	var options EncoderOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &Encoder{options: options}
}

// Encode encodes a single operation into compact format.
func (e *Encoder) Encode(o internal.Op) (Op, error) {
	return encodeOp(o, e.options)
}

// EncodeSlice encodes multiple operations into compact format.
func (e *Encoder) EncodeSlice(ops []internal.Op) ([]Op, error) {
	result := make([]Op, len(ops))
	for i, o := range ops {
		encoded, err := encodeOp(o, e.options)
		if err != nil {
			return nil, err
		}
		result[i] = encoded
	}
	return result, nil
}

// Encode encodes operations into compact format.
func Encode(ops []internal.Op, opts ...EncoderOption) ([]Op, error) {
	return NewEncoder(opts...).EncodeSlice(ops)
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

	// Convert []string path fields to JSON Pointer strings.
	for i := 1; i <= 2 && i < len(result); i++ {
		if pathSlice, ok := result[i].([]string); ok {
			result[i] = formatPath(pathSlice)
		}
	}

	return result, nil
}

// formatPath converts a path slice to JSON Pointer string.
func formatPath(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return jsonpointer.Format(path...)
}
