package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Encoder encodes operations into compact format.
type Encoder struct {
	options EncoderOptions
}

// NewEncoder creates a new compact encoder with the given options.
func NewEncoder(opts ...EncoderOption) *Encoder {
	options := defaultEncoderOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &Encoder{options: options}
}

// Encode encodes a single operation into compact format.
func (e *Encoder) Encode(op internal.Op) (Op, error) {
	return encodeOp(op, e.options)
}

// EncodeSlice encodes multiple operations into compact format.
func (e *Encoder) EncodeSlice(ops []internal.Op) ([]Op, error) {
	result := make([]Op, len(ops))
	for i, op := range ops {
		compactOp, err := e.Encode(op)
		if err != nil {
			return nil, err
		}
		result[i] = compactOp
	}
	return result, nil
}
