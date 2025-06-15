// Package json implements JSON Patch encoder functionality.
package json

import "github.com/kaptinlin/jsonpatch/internal"

// Encoder provides JSON patch encoding functionality.
type Encoder struct{}

// NewEncoder creates a new Encoder instance.
func NewEncoder() *Encoder {
	return &Encoder{}
}

// Encode encodes operations to JSON format.
func (e *Encoder) Encode(ops []internal.Op) ([]map[string]interface{}, error) {
	return Encode(ops)
}

// EncodeJSON encodes operations to JSON bytes.
func (e *Encoder) EncodeJSON(ops []internal.Op) ([]byte, error) {
	return EncodeJSON(ops)
}
