//nolint:gosec // Integer size conversions are safe due to operation code and slice length bounds.
package binary

import (
	"bytes"
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/tinylib/msgp/msgp"
)

// encode converts a slice of operations to a byte slice in binary format.
func (c *Codec) encode(ops []internal.Op) ([]byte, error) {
	var buf bytes.Buffer
	writer := msgp.NewWriter(&buf)
	if err := encodeOps(writer, ops); err != nil {
		return nil, err
	}
	if err := writer.Flush(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encodeOps(writer *msgp.Writer, ops []internal.Op) error {
	if err := writer.WriteFloat64(float64(len(ops))); err != nil {
		return err
	}
	for _, op := range ops {
		if err := encodeOp(writer, op); err != nil {
			return err
		}
	}
	return nil
}

func encodeOp(writer *msgp.Writer, i internal.Op) error {
	switch o := i.(type) {
	case *op.AddOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.RemoveOperation:
		if err := writer.WriteArrayHeader(2); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.ReplaceOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.MoveOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.From()); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.CopyOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.From()); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.TestOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	// Predicate operations
	case *op.TestTypeOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Types)
	case *op.DefinedOperation:
		if err := writer.WriteArrayHeader(2); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.UndefinedOperation:
		if err := writer.WriteArrayHeader(2); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.LessOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.MoreOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.ContainsOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.InOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.StartsOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.EndsOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return encodeValue(writer, o.Value)
	case *op.MatchesOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteString(o.Pattern); err != nil {
			return err
		}
		return writer.WriteBool(o.IgnoreCase)
	case *op.TestStringOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteString(o.Str); err != nil {
			return err
		}
		return writer.WriteFloat64(float64(o.Pos))
	case *op.TestStringLenOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Length); err != nil {
			return err
		}
		return writer.WriteBool(o.Not())
	// Type predicate operation
	case *op.TypeOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return writer.WriteString(o.TypeValue)
	// JSON Patch Extended
	case *op.FlipOperation:
		if err := writer.WriteArrayHeader(2); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		return encodePath(writer, o.Path())
	case *op.IncOperation:
		if err := writer.WriteArrayHeader(3); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		return writer.WriteFloat64(o.Inc)
	case *op.StrInsOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Pos); err != nil {
			return err
		}
		return writer.WriteString(o.Str)
	case *op.StrDelOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Pos); err != nil {
			return err
		}
		return writer.WriteFloat64(o.Len)
	case *op.SplitOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Pos); err != nil {
			return err
		}
		return encodeValue(writer, o.Props)
	case *op.ExtendOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := encodeValue(writer, o.Properties); err != nil {
			return err
		}
		return writer.WriteBool(o.DeleteNull)
	case *op.MergeOperation:
		if err := writer.WriteArrayHeader(4); err != nil {
			return err
		}
		if err := writer.WriteUint8(uint8(o.Code())); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Pos); err != nil {
			return err
		}
		return encodeValue(writer, o.Props)
	default:
		return fmt.Errorf("%w: %T", ErrUnsupportedOp, i)
	}
}

func encodePath(writer *msgp.Writer, path []string) error {
	if err := writer.WriteFloat64(float64(len(path))); err != nil {
		return err
	}
	for _, segment := range path {
		if err := writer.WriteString(segment); err != nil {
			return err
		}
	}
	return nil
}

func encodeValue(writer *msgp.Writer, value interface{}) error {
	return writer.WriteIntf(value)
}
