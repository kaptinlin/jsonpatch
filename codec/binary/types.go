package binary

import (
	"bytes"

	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
)

// Codec encodes and decodes JSON Patch operations in MessagePack binary format.
type Codec struct{}

// New creates a new binary Codec.
func New() *Codec {
	return &Codec{}
}

// Encode serializes operations into MessagePack binary format.
func (c *Codec) Encode(ops []internal.Op) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(len(ops) * 32) // pre-allocate based on typical operation size
	w := msgp.NewWriter(&buf)
	if err := encodeOps(w, ops); err != nil {
		return nil, err
	}
	if err := w.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes operations from MessagePack binary format.
func (c *Codec) Decode(data []byte) ([]internal.Op, error) {
	r := msgp.NewReader(bytes.NewReader(data))
	return decodeOps(r)
}
