package binary

import (
	"bytes"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/tinylib/msgp/msgp"
)

// Codec encodes and decodes JSON Patch operations in MessagePack binary format.
type Codec struct{}

// NewCodec creates a new binary Codec.
func NewCodec() *Codec {
	return &Codec{}
}

// Encode serializes operations into MessagePack binary format.
func (c *Codec) Encode(ops []internal.Op) ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(len(ops) * 32) // pre-allocate based on typical operation size
	writer := msgp.NewWriter(&buf)
	if err := encodeOps(writer, ops); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Decode deserializes operations from MessagePack binary format.
func (c *Codec) Decode(data []byte) ([]internal.Op, error) {
	reader := msgp.NewReader(bytes.NewReader(data))
	return decodeOps(reader)
}
