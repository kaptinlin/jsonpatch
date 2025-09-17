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
	encoder.WriteFloat64(float64(len(ops)))
	for _, op := range ops {
		if err := encodeOp(encoder, op); err != nil {
			return err
		}
	}
	return encoder.Err()
}

func encodeOp(encoder msgpack.Writer, i internal.Op) error {
	switch o := i.(type) {
	case *op.AddOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.RemoveOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.ReplaceOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.MoveOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.From())
		encodePath(encoder, o.Path())
	case *op.CopyOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.From())
		encodePath(encoder, o.Path())
	case *op.TestOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	// Predicate operations
	case *op.TestTypeOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Types)
	case *op.DefinedOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.UndefinedOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteBool(o.Not())
		return encoder.Err()
	case *op.LessOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.MoreOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.ContainsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.InOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.StartsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.EndsOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		return encodeValue(encoder, o.Value)
	case *op.MatchesOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.Pattern)
		encoder.WriteBool(o.IgnoreCase)
		return encoder.Err()
	case *op.TestStringOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.Str)
		encoder.WriteFloat64(o.Pos)
		return encoder.Err()
	case *op.TestStringLenOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Length)
		encoder.WriteBool(o.Not)
		return encoder.Err()
	// Type predicate operation
	case *op.TypeOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteString(o.TypeValue)
		return encoder.Err()
	// JSON Patch Extended
	case *op.FlipOperation:
		encoder.WriteArraySize(2)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
	case *op.IncOperation:
		encoder.WriteArraySize(3)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Inc)
		return encoder.Err()
	case *op.StrInsOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Pos)
		encoder.WriteString(o.Str)
		return encoder.Err()
	case *op.StrDelOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Pos)
		encoder.WriteFloat64(o.Len)
		return encoder.Err()
	case *op.SplitOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Pos)
		return encodeValue(encoder, o.Props)
	case *op.ExtendOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		if err := encodeValue(encoder, o.Properties); err != nil {
			return err
		}
		encoder.WriteBool(o.DeleteNull)
		return encoder.Err()
	case *op.MergeOperation:
		encoder.WriteArraySize(4)
		encoder.WriteUint8(uint8(o.Code()))
		encodePath(encoder, o.Path())
		encoder.WriteFloat64(o.Pos)
		return encodeValue(encoder, o.Props)
	}
	return nil
}

func encodePath(encoder msgpack.Writer, path []string) {
	encoder.WriteFloat64(float64(len(path)))
	for _, segment := range path {
		encoder.WriteString(segment)
	}
}

func encodeValue(encoder msgpack.Writer, value interface{}) error {
	encoder.WriteAny(value)
	return encoder.Err()
}
