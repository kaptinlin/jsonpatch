package binary

import "github.com/kaptinlin/jsonpatch/internal"

// Codec implements the Codec interface for the binary format.
type Codec struct{}

// NewCodec creates a new Codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Encode encodes a slice of operations into binary format.
func (c *Codec) Encode(ops []internal.Op) ([]byte, error) {
	return c.encode(ops)
}

// Decode decodes a slice of operations from binary format.
func (c *Codec) Decode(data []byte) ([]internal.Op, error) {
	return c.decode(data)
}
