package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Decoder decodes compact format operations into operation instances.
type Decoder struct {
	options DecoderOptions
}

// NewDecoder creates a new compact decoder with the given options.
func NewDecoder(opts ...DecoderOption) *Decoder {
	options := defaultDecoderOptions
	for _, opt := range opts {
		opt(&options)
	}
	return &Decoder{options: options}
}

// Decode decodes a single compact operation into an operation instance.
func (d *Decoder) Decode(compactOp Op) (internal.Op, error) {
	return decodeOp(compactOp, d.options)
}

// DecodeSlice decodes multiple compact operations into operation instances.
func (d *Decoder) DecodeSlice(compactOps []Op) ([]internal.Op, error) {
	result := make([]internal.Op, len(compactOps))
	for i, compactOp := range compactOps {
		op, err := d.Decode(compactOp)
		if err != nil {
			return nil, err
		}
		result[i] = op
	}
	return result, nil
}
