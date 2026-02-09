package json

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Decoder provides JSON patch decoding with configurable options.
type Decoder struct {
	options internal.JSONPatchOptions
}

// NewDecoder creates a new Decoder with the given options.
func NewDecoder(options internal.JSONPatchOptions) *Decoder {
	return &Decoder{options: options}
}

// Decode decodes a JSON patch array to operations.
func (d *Decoder) Decode(patch []map[string]any) ([]internal.Op, error) {
	return Decode(patch, d.options)
}
