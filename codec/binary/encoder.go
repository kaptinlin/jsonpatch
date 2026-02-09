// Package binary implements a MessagePack-based binary codec for JSON Patch operations.
//
// Limitations:
//   - Second-order predicates (and, or, not) are NOT supported in binary codec.
//     These operations are skipped during encoding with a warning.
//     Use the JSON or compact codec if you need second-order predicate support.
//
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
	for _, o := range ops {
		if err := encodeOp(writer, o); err != nil {
			return err
		}
	}
	return nil
}

// writeHeader writes the array header and operation code.
func writeHeader(writer *msgp.Writer, size uint32, code int) error {
	if err := writer.WriteArrayHeader(size); err != nil {
		return err
	}
	return writer.WriteUint8(uint8(code))
}

// encodePathOnly encodes operations with format: [code, path] (e.g. remove, defined, undefined, flip).
func encodePathOnly(writer *msgp.Writer, code int, path []string) error {
	if err := writeHeader(writer, 2, code); err != nil {
		return err
	}
	return encodePath(writer, path)
}

// encodePathValue encodes operations with format: [code, path, value] (e.g. add, replace, test).
func encodePathValue(writer *msgp.Writer, code int, path []string, value any) error {
	if err := writeHeader(writer, 3, code); err != nil {
		return err
	}
	if err := encodePath(writer, path); err != nil {
		return err
	}
	return encodeValue(writer, value)
}

// encodePathPaths encodes operations with format: [code, from, path] (e.g. move, copy).
func encodePathPaths(writer *msgp.Writer, code int, from, path []string) error {
	if err := writeHeader(writer, 3, code); err != nil {
		return err
	}
	if err := encodePath(writer, from); err != nil {
		return err
	}
	return encodePath(writer, path)
}

func encodeOp(writer *msgp.Writer, i internal.Op) error {
	switch o := i.(type) {
	// Standard RFC 6902
	case *op.AddOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.RemoveOperation:
		return encodePathOnly(writer, o.Code(), o.Path())
	case *op.ReplaceOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.MoveOperation:
		return encodePathPaths(writer, o.Code(), o.From(), o.Path())
	case *op.CopyOperation:
		return encodePathPaths(writer, o.Code(), o.From(), o.Path())
	case *op.TestOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)

	// Predicate operations with value
	case *op.TestTypeOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Types)
	case *op.LessOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.MoreOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.ContainsOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.InOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.StartsOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)
	case *op.EndsOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Value)

	// Predicate operations with path only
	case *op.DefinedOperation:
		return encodePathOnly(writer, o.Code(), o.Path())
	case *op.UndefinedOperation:
		return encodePathOnly(writer, o.Code(), o.Path())

	// Predicate operations with custom fields
	case *op.MatchesOperation:
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
			return err
		}
		if err := encodePath(writer, o.Path()); err != nil {
			return err
		}
		if err := writer.WriteFloat64(o.Length); err != nil {
			return err
		}
		return writer.WriteBool(o.Not())

	case *op.TypeOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.TypeValue)

	// Extended operations
	case *op.FlipOperation:
		return encodePathOnly(writer, o.Code(), o.Path())

	case *op.IncOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Inc)

	case *op.StrInsOperation:
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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
		if err := writeHeader(writer, 4, o.Code()); err != nil {
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

func encodeValue(writer *msgp.Writer, value any) error {
	return writer.WriteIntf(value)
}
