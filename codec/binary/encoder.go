package binary

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/tinylib/msgp/msgp"
)

// encodeOps writes the operation count followed by each encoded operation.
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

// encodeOp dispatches encoding to the appropriate helper by operation type.
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
		return encodeMatches(writer, o)
	case *op.TestStringOperation:
		return encodeTestString(writer, o)
	case *op.TestStringLenOperation:
		return encodeTestStringLen(writer, o)
	case *op.TypeOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.TypeValue)

	// Extended operations
	case *op.FlipOperation:
		return encodePathOnly(writer, o.Code(), o.Path())
	case *op.IncOperation:
		return encodePathValue(writer, o.Code(), o.Path(), o.Inc)
	case *op.StrInsOperation:
		return encodeStrIns(writer, o)
	case *op.StrDelOperation:
		return encodeStrDel(writer, o)
	case *op.SplitOperation:
		return encodePathFloat64Value(writer, o.Code(), o.Path(), o.Pos, o.Props)
	case *op.ExtendOperation:
		return encodeExtend(writer, o)
	case *op.MergeOperation:
		return encodePathFloat64Value(writer, o.Code(), o.Path(), o.Pos, o.Props)

	default:
		return fmt.Errorf("unsupported op type %T: %w", i, ErrUnsupportedOp)
	}
}

// --- Encoding primitives ---

// writeHeader writes the array header and operation code.
func writeHeader(writer *msgp.Writer, size uint32, code int) error {
	if err := writer.WriteArrayHeader(size); err != nil {
		return err
	}
	return writer.WriteUint8(uint8(code)) //nolint:gosec // Operation codes are bounded constants within uint8 range.
}

// encodePathOnly encodes operations with format: [code, path].
func encodePathOnly(writer *msgp.Writer, code int, path []string) error {
	if err := writeHeader(writer, 2, code); err != nil {
		return err
	}
	return encodePath(writer, path)
}

// encodePathValue encodes operations with format: [code, path, value].
func encodePathValue(writer *msgp.Writer, code int, path []string, value any) error {
	if err := writeHeader(writer, 3, code); err != nil {
		return err
	}
	if err := encodePath(writer, path); err != nil {
		return err
	}
	return encodeValue(writer, value)
}

// encodePathPaths encodes operations with format: [code, from, path].
func encodePathPaths(writer *msgp.Writer, code int, from, path []string) error {
	if err := writeHeader(writer, 3, code); err != nil {
		return err
	}
	if err := encodePath(writer, from); err != nil {
		return err
	}
	return encodePath(writer, path)
}

// encodePathFloat64Value encodes operations with format: [code, path, float64, value].
func encodePathFloat64Value(writer *msgp.Writer, code int, path []string, f float64, value any) error {
	if err := writeHeader(writer, 4, code); err != nil {
		return err
	}
	if err := encodePath(writer, path); err != nil {
		return err
	}
	if err := writer.WriteFloat64(f); err != nil {
		return err
	}
	return encodeValue(writer, value)
}

// --- Operation-specific encoders ---

// encodeMatches encodes a matches predicate: [code, path, pattern, ignoreCase].
func encodeMatches(writer *msgp.Writer, o *op.MatchesOperation) error {
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
}

// encodeTestString encodes a test_string operation: [code, path, str, pos].
func encodeTestString(writer *msgp.Writer, o *op.TestStringOperation) error {
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
}

// encodeTestStringLen encodes a test_string_len operation: [code, path, length, not].
func encodeTestStringLen(writer *msgp.Writer, o *op.TestStringLenOperation) error {
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
}

// encodeStrIns encodes a str_ins operation: [code, path, pos, str].
func encodeStrIns(writer *msgp.Writer, o *op.StrInsOperation) error {
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
}

// encodeStrDel encodes a str_del operation: [code, path, pos, len].
func encodeStrDel(writer *msgp.Writer, o *op.StrDelOperation) error {
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
}

// encodeExtend encodes an extend operation: [code, path, properties, deleteNull].
func encodeExtend(writer *msgp.Writer, o *op.ExtendOperation) error {
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
}

// --- Low-level encoders ---

// encodePath writes a path as a float64 count followed by string segments.
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

// encodeValue writes an arbitrary value using msgp interface encoding.
func encodeValue(writer *msgp.Writer, value any) error {
	return writer.WriteIntf(value)
}
