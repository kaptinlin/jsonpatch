//nolint:gosec // Integer size conversions are safe due to operation code and slice length bounds.
package binary

import (
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	msgpack "github.com/wapc/tinygo-msgpack"
)

// encode converts a slice of operations to a byte slice in binary format.
func (c *Codec) encode(ops []internal.Op) ([]byte, error) {
	var sizer msgpack.Sizer
	if err := encodeOps(&sizer, ops); err != nil {
		return nil, err
	}

	buf := make([]byte, sizer.Len())
	encoder := msgpack.NewEncoder(buf)
	if err := encodeOps(&encoder, ops); err != nil {
		return nil, err
	}

	return buf, nil
}

func encodeOps(encoder msgpack.Writer, ops []internal.Op) error {
	encoder.WriteArraySize(uint32(len(ops)))
	for _, op := range ops {
		if err := encodeOp(encoder, op); err != nil {
			return err
		}
	}
	return encoder.Err()
}

func encodeOp(encoder msgpack.Writer, i internal.Op) error {
	switch o := i.(type) {
	case *op.OpAddOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpRemoveOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.OpReplaceOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpMoveOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.From())
		encodePath(encoder, o.Path())
	case *op.OpCopyOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.From())
		encodePath(encoder, o.Path())
	case *op.OpTestOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	// Predicate operations
	case *op.OpTestTypeOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Types)
	case *op.OpDefinedOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.OpUndefinedOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteBool(o.Not())
		return encoder.Err()
	case *op.OpLessOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpMoreOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpContainsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpInOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Values)
	case *op.OpStartsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpEndsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.OpMatchesOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.Pattern)
		encoder.WriteBool(o.IgnoreCase)
		return encoder.Err()
	case *op.OpTestStringOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.Str)
		encoder.WriteInt64(int64(o.Pos))
		return encoder.Err()
	case *op.OpTestStringLenOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteInt64(int64(o.Length))
		encoder.WriteBool(o.Not)
		return encoder.Err()
	// Type predicate operation
	case *op.OpTypeOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.TypeValue)
		return encoder.Err()
	// JSON Patch Extended
	case *op.OpFlipOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.OpIncOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Inc)
		return encoder.Err()
	case *op.OpStrInsOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteInt64(int64(o.Pos))
		encoder.WriteString(o.Str)
		return encoder.Err()
	case *op.OpStrDelOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteInt64(int64(o.Pos))
		encoder.WriteInt64(int64(o.Len))
		return encoder.Err()
	case *op.OpSplitOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteInt64(int64(o.Pos))
		return encodeValue(encoder, o.Props)
	case *op.OpExtendOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		if err := encodeValue(encoder, o.Properties); err != nil {
			return err
		}
		encoder.WriteBool(o.DeleteNull)
		return encoder.Err()
	case *op.OpMergeOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteInt64(int64(o.Pos))
		return encodeValue(encoder, o.Props)
	}
	return nil
}

func encodePath(encoder msgpack.Writer, path []string) {
	encoder.WriteArraySize(uint32(len(path)))
	for _, segment := range path {
		encoder.WriteString(segment)
	}
}

func encodeValue(encoder msgpack.Writer, value interface{}) error {
	encoder.WriteAny(value)
	return encoder.Err()
}
