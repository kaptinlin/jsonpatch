package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Decoder decodes compact format operations into operation instances
type Decoder struct {
	options DecoderOptions
}

// NewDecoder creates a new compact decoder with the given options
func NewDecoder(opts ...DecoderOption) *Decoder {
	options := DefaultDecoderOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &Decoder{
		options: options,
	}
}

// Decode decodes a single compact operation into an operation instance
func (d *Decoder) Decode(compactOp Op) (internal.Op, error) {
	return compactToOp(compactOp, d.options)
}

// DecodeSlice decodes multiple compact operations into operation instances
func (d *Decoder) DecodeSlice(compactOps []Op) ([]internal.Op, error) {
	result := make([]internal.Op, 0, len(compactOps))
	for _, compactOp := range compactOps {
		op, err := d.Decode(compactOp)
		if err != nil {
			return nil, err
		}
		result = append(result, op)
	}
	return result, nil
}
