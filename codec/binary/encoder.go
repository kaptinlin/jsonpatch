package binary

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/tinylib/msgp/msgp"
)

// encodeOps writes the operation count followed by each encoded operation.
func encodeOps(w *msgp.Writer, ops []internal.Op) error {
	if err := w.WriteFloat64(float64(len(ops))); err != nil {
		return err
	}
	for _, o := range ops {
		if err := encodeOp(w, o); err != nil {
			return err
		}
	}
	return nil
}

// encodeOp dispatches encoding to the appropriate helper by operation type.
func encodeOp(w *msgp.Writer, v internal.Op) error {
	switch o := v.(type) {
	// Standard RFC 6902
	case *op.AddOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.RemoveOperation:
		return encodePathOnly(w, o.Code(), o.Path())
	case *op.ReplaceOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.MoveOperation:
		return encodePathPaths(w, o.Code(), o.From(), o.Path())
	case *op.CopyOperation:
		return encodePathPaths(w, o.Code(), o.From(), o.Path())
	case *op.TestOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)

	// Predicate operations with value
	case *op.TestTypeOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Types)
	case *op.LessOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.MoreOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.ContainsOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.InOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.StartsOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)
	case *op.EndsOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Value)

	// Predicate operations with path only
	case *op.DefinedOperation:
		return encodePathOnly(w, o.Code(), o.Path())
	case *op.UndefinedOperation:
		return encodePathOnly(w, o.Code(), o.Path())

	// Predicate operations with custom fields
	case *op.MatchesOperation:
		return encodeMatches(w, o)
	case *op.TestStringOperation:
		return encodeTestString(w, o)
	case *op.TestStringLenOperation:
		return encodeTestStringLen(w, o)
	case *op.TypeOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.TypeValue)

	// Extended operations
	case *op.FlipOperation:
		return encodePathOnly(w, o.Code(), o.Path())
	case *op.IncOperation:
		return encodePathValue(w, o.Code(), o.Path(), o.Inc)
	case *op.StrInsOperation:
		return encodeStrIns(w, o)
	case *op.StrDelOperation:
		return encodeStrDel(w, o)
	case *op.SplitOperation:
		return encodePathFloat64Value(w, o.Code(), o.Path(), o.Pos, o.Props)
	case *op.ExtendOperation:
		return encodeExtend(w, o)
	case *op.MergeOperation:
		return encodePathFloat64Value(w, o.Code(), o.Path(), o.Pos, o.Props)

	default:
		return fmt.Errorf("unsupported op type %T: %w", v, ErrUnsupportedOp)
	}
}

// writeHeader writes the array header and operation code.
func writeHeader(w *msgp.Writer, size uint32, code int) error {
	if err := w.WriteArrayHeader(size); err != nil {
		return err
	}
	return w.WriteUint8(uint8(code)) //nolint:gosec // Operation codes are bounded constants within uint8 range.
}

// encodePathOnly encodes operations with format: [code, path].
func encodePathOnly(w *msgp.Writer, code int, path []string) error {
	if err := writeHeader(w, 2, code); err != nil {
		return err
	}
	return encodePath(w, path)
}

// encodePathValue encodes operations with format: [code, path, value].
func encodePathValue(w *msgp.Writer, code int, path []string, value any) error {
	if err := writeHeader(w, 3, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	return encodeValue(w, value)
}

// encodePathPaths encodes operations with format: [code, from, path].
func encodePathPaths(w *msgp.Writer, code int, from, path []string) error {
	if err := writeHeader(w, 3, code); err != nil {
		return err
	}
	if err := encodePath(w, from); err != nil {
		return err
	}
	return encodePath(w, path)
}

// encodePathFloat64Value encodes operations with format: [code, path, float64, value].
func encodePathFloat64Value(w *msgp.Writer, code int, path []string, f float64, value any) error {
	if err := writeHeader(w, 4, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(f); err != nil {
		return err
	}
	return encodeValue(w, value)
}

// encodeMatches encodes a matches predicate: [code, path, pattern, ignoreCase].
func encodeMatches(w *msgp.Writer, o *op.MatchesOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := w.WriteString(o.Pattern); err != nil {
		return err
	}
	return w.WriteBool(o.IgnoreCase)
}

// encodeTestString encodes a test_string operation: [code, path, str, pos].
func encodeTestString(w *msgp.Writer, o *op.TestStringOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := w.WriteString(o.Str); err != nil {
		return err
	}
	return w.WriteFloat64(float64(o.Pos))
}

// encodeTestStringLen encodes a test_string_len operation: [code, path, length, not].
func encodeTestStringLen(w *msgp.Writer, o *op.TestStringLenOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := w.WriteFloat64(o.Length); err != nil {
		return err
	}
	return w.WriteBool(o.Not())
}

// encodeStrIns encodes a str_ins operation: [code, path, pos, str].
func encodeStrIns(w *msgp.Writer, o *op.StrInsOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := w.WriteFloat64(o.Pos); err != nil {
		return err
	}
	return w.WriteString(o.Str)
}

// encodeStrDel encodes a str_del operation: [code, path, pos, len].
func encodeStrDel(w *msgp.Writer, o *op.StrDelOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := w.WriteFloat64(o.Pos); err != nil {
		return err
	}
	return w.WriteFloat64(o.Len)
}

// encodeExtend encodes an extend operation: [code, path, properties, deleteNull].
func encodeExtend(w *msgp.Writer, o *op.ExtendOperation) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, o.Path()); err != nil {
		return err
	}
	if err := encodeValue(w, o.Properties); err != nil {
		return err
	}
	return w.WriteBool(o.DeleteNull)
}

// encodePath writes a path as a float64 count followed by string segments.
func encodePath(w *msgp.Writer, path []string) error {
	if err := w.WriteFloat64(float64(len(path))); err != nil {
		return err
	}
	for _, seg := range path {
		if err := w.WriteString(seg); err != nil {
			return err
		}
	}
	return nil
}

// encodeValue writes an arbitrary value using msgp interface encoding.
func encodeValue(w *msgp.Writer, value any) error {
	return w.WriteIntf(value)
}
