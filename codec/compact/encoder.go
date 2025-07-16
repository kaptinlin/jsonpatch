package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Encoder encodes operations into compact format
type Encoder struct {
	options EncoderOptions
}

// NewEncoder creates a new compact encoder with the given options
func NewEncoder(opts ...EncoderOption) *Encoder {
	options := DefaultEncoderOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &Encoder{
		options: options,
	}
}

// Encode encodes a single operation into compact format
func (e *Encoder) Encode(op internal.Op) (CompactOp, error) {
	return opToCompact(op, e.options)
}

// EncodeSlice encodes multiple operations into compact format
func (e *Encoder) EncodeSlice(ops []internal.Op) ([]CompactOp, error) {
	result := make([]CompactOp, 0, len(ops))
	for _, op := range ops {
		compactOp, err := e.Encode(op)
		if err != nil {
			return nil, err
		}
		result = append(result, compactOp)
	}
	return result, nil
}
