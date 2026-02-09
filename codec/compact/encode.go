package compact

import (
	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpointer"
)

// Encoder encodes operations into compact format.
type Encoder struct {
	opts Options
}

// NewEncoder creates a new compact encoder with the given options.
func NewEncoder(opts ...Option) *Encoder {
	var o Options
	for _, opt := range opts {
		opt(&o)
	}
	return &Encoder{opts: o}
}

// Encode encodes a single operation into compact format.
func (e *Encoder) Encode(o internal.Op) (Op, error) {
	return encodeOp(o, e.opts)
}

// EncodeSlice encodes multiple operations into compact format.
func (e *Encoder) EncodeSlice(ops []internal.Op) ([]Op, error) {
	result := make([]Op, len(ops))
	for i, o := range ops {
		encoded, err := encodeOp(o, e.opts)
		if err != nil {
			return nil, err
		}
		result[i] = encoded
	}
	return result, nil
}

// Encode encodes operations into compact format.
func Encode(ops []internal.Op, opts ...Option) ([]Op, error) {
	return NewEncoder(opts...).EncodeSlice(ops)
}

// EncodeJSON encodes operations into compact JSON bytes.
func EncodeJSON(ops []internal.Op, opts ...Option) ([]byte, error) {
	compact, err := Encode(ops, opts...)
	if err != nil {
		return nil, err
	}
	return json.Marshal(compact)
}

// encodeOp converts a single operation to compact format.
func encodeOp(o internal.Op, opts Options) (Op, error) {
	raw, err := o.ToCompact()
	if err != nil {
		return nil, err
	}

	result := make(Op, len(raw))
	copy(result, raw)

	if opts.StringOpcode {
		result[0] = string(o.Op())
	}

	// Convert []string path fields to JSON Pointer strings.
	for i := 1; i <= 2 && i < len(result); i++ {
		if segments, ok := result[i].([]string); ok {
			result[i] = formatPath(segments)
		}
	}

	return result, nil
}

// formatPath converts a path slice to a JSON Pointer string.
func formatPath(segments []string) string {
	if len(segments) == 0 {
		return ""
	}
	return jsonpointer.Format(segments...)
}
