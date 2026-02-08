// Package json implements JSON Patch decoder functionality.
package json

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Decoder provides JSON patch decoding functionality with configurable options.
type Decoder struct {
	options internal.JSONPatchOptions
}

// NewDecoder creates a new Decoder with the given options.
func NewDecoder(options internal.JSONPatchOptions) *Decoder {
	return &Decoder{
		options: options,
	}
}

// Decode decodes a JSON patch array to operations using the decoder's options.
func (d *Decoder) Decode(patch []map[string]interface{}) ([]internal.Op, error) {
	return Decode(patch, d.options)
}
