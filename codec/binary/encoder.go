package binary

import (
	"fmt"
	"slices"

	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

// encodeOps writes the operation count followed by each encoded operation.
func encodeOps(w *msgp.Writer, ops []internal.Op) error {
	if err := w.WriteArrayHeader(uint32(len(ops))); err != nil { //nolint:gosec // ops length is bounded by practical limits.
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
	return encodeOpWithParent(w, v, nil)
}

func encodeOpWithParent(w *msgp.Writer, v internal.Op, parent []string) error {
	path, err := relativePath(parent, v.Path())
	if err != nil {
		return err
	}

	switch o := v.(type) {
	// Standard RFC 6902
	case *op.AddOperation:
		return encodePathValue(w, o.Code(), path, o.Value)
	case *op.RemoveOperation:
		if o.HasOldValue {
			return encodePathValue(w, o.Code(), path, o.OldValue)
		}
		return encodePathOnly(w, o.Code(), path)
	case *op.ReplaceOperation:
		return encodePathValue(w, o.Code(), path, o.Value)
	case *op.MoveOperation:
		return encodePathPaths(w, o.Code(), path, o.From())
	case *op.CopyOperation:
		return encodePathPaths(w, o.Code(), path, o.From())
	case *op.TestOperation:
		return encodeTest(w, o, path)

	// Predicate operations with value
	case *op.TestTypeOperation:
		return encodePathValue(w, o.Code(), path, o.Types)
	case *op.LessOperation:
		return encodePathValue(w, o.Code(), path, o.Value)
	case *op.MoreOperation:
		return encodePathValue(w, o.Code(), path, o.Value)
	case *op.ContainsOperation:
		return encodeStringPredicate(w, o.Code(), path, o.Value, o.IgnoreCase)
	case *op.InOperation:
		return encodePathValue(w, o.Code(), path, o.Value)
	case *op.StartsOperation:
		return encodeStringPredicate(w, o.Code(), path, o.Value, o.IgnoreCase)
	case *op.EndsOperation:
		return encodeStringPredicate(w, o.Code(), path, o.Value, o.IgnoreCase)

	// Predicate operations with path only
	case *op.DefinedOperation:
		return encodePathOnly(w, o.Code(), path)
	case *op.UndefinedOperation:
		return encodePathOnly(w, o.Code(), path)

	// Predicate operations with custom fields
	case *op.MatchesOperation:
		return encodeMatches(w, o, path)
	case *op.TestStringOperation:
		return encodeTestString(w, o, path)
	case *op.TestStringLenOperation:
		return encodeTestStringLen(w, o, path)
	case *op.TypeOperation:
		return encodePathValue(w, o.Code(), path, o.TypeValue)
	case *op.AndOperation:
		return encodeComposite(w, o.Code(), path, o.Path(), o.Operations, op.ErrInvalidPredicateInAnd)
	case *op.OrOperation:
		return encodeComposite(w, o.Code(), path, o.Path(), o.Operations, op.ErrInvalidPredicateInOr)
	case *op.NotOperation:
		if err := o.Validate(); err != nil {
			return err
		}
		return encodeComposite(w, o.Code(), path, o.Path(), o.Operations, op.ErrInvalidPredicateInNot)

	// Extended operations
	case *op.FlipOperation:
		return encodePathOnly(w, o.Code(), path)
	case *op.IncOperation:
		return encodePathValue(w, o.Code(), path, o.Inc)
	case *op.StrInsOperation:
		return encodeStrIns(w, o, path)
	case *op.StrDelOperation:
		return encodeStrDel(w, o, path)
	case *op.SplitOperation:
		return encodeSplitOrMerge(w, o.Code(), path, o.Pos, o.Props)
	case *op.ExtendOperation:
		return encodeExtend(w, o, path)
	case *op.MergeOperation:
		var props any
		if o.Props != nil {
			props = o.Props
		}
		return encodeSplitOrMerge(w, o.Code(), path, o.Pos, props)

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
	return w.WriteIntf(value)
}

// encodePathPaths encodes operations with format: [code, path, from].
func encodePathPaths(w *msgp.Writer, code int, path, from []string) error {
	if err := writeHeader(w, 3, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	return encodePath(w, from)
}

func encodeTest(w *msgp.Writer, o *op.TestOperation, path []string) error {
	size := uint32(3)
	if o.Not() {
		size = 4
	}
	if err := writeHeader(w, size, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteIntf(o.Value); err != nil {
		return err
	}
	if o.Not() {
		return w.WriteBool(true)
	}
	return nil
}

func encodeStringPredicate(w *msgp.Writer, code int, path []string, value string, ignoreCase bool) error {
	size := uint32(3)
	if ignoreCase {
		size = 4
	}
	if err := writeHeader(w, size, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteString(value); err != nil {
		return err
	}
	if ignoreCase {
		return w.WriteBool(true)
	}
	return nil
}

// encodeSplitOrMerge encodes split/merge operations with optional props.
// When props is nil, format is [code, path, float64]; otherwise [code, path, float64, props].
func encodeSplitOrMerge(w *msgp.Writer, code int, path []string, f float64, props any) error {
	if props == nil {
		if err := writeHeader(w, 3, code); err != nil {
			return err
		}
		if err := encodePath(w, path); err != nil {
			return err
		}
		return w.WriteFloat64(f)
	}
	if err := writeHeader(w, 4, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(f); err != nil {
		return err
	}
	return w.WriteIntf(props)
}

// encodeMatches encodes a matches predicate: [code, path, pattern, ignoreCase?].
func encodeMatches(w *msgp.Writer, o *op.MatchesOperation, path []string) error {
	size := uint32(3)
	if o.IgnoreCase {
		size = 4
	}
	if err := writeHeader(w, size, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteString(o.Pattern); err != nil {
		return err
	}
	if o.IgnoreCase {
		return w.WriteBool(true)
	}
	return nil
}

// encodeTestString encodes a test_string operation: [code, path, pos, str, not?].
func encodeTestString(w *msgp.Writer, o *op.TestStringOperation, path []string) error {
	size := uint32(4)
	if o.Not() {
		size = 5
	}
	if err := writeHeader(w, size, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(float64(o.Pos)); err != nil {
		return err
	}
	if err := w.WriteString(o.Str); err != nil {
		return err
	}
	if o.Not() {
		return w.WriteBool(true)
	}
	return nil
}

// encodeTestStringLen encodes a test_string_len operation: [code, path, length, not?].
func encodeTestStringLen(w *msgp.Writer, o *op.TestStringLenOperation, path []string) error {
	size := uint32(3)
	if o.Not() {
		size = 4
	}
	if err := writeHeader(w, size, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(o.Length); err != nil {
		return err
	}
	if o.Not() {
		return w.WriteBool(true)
	}
	return nil
}

// encodeStrIns encodes a str_ins operation: [code, path, pos, str].
func encodeStrIns(w *msgp.Writer, o *op.StrInsOperation, path []string) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(float64(o.Pos)); err != nil {
		return err
	}
	return w.WriteString(o.Str)
}

// encodeStrDel encodes a str_del operation: [code, path, pos, len].
func encodeStrDel(w *msgp.Writer, o *op.StrDelOperation, path []string) error {
	if err := writeHeader(w, 4, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteFloat64(float64(o.Pos)); err != nil {
		return err
	}
	return w.WriteFloat64(float64(o.Len))
}

// encodeExtend encodes an extend operation: [code, path, properties, deleteNull?].
func encodeExtend(w *msgp.Writer, o *op.ExtendOperation, path []string) error {
	size := uint32(3)
	if o.DeleteNull {
		size = 4
	}
	if err := writeHeader(w, size, o.Code()); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteIntf(o.Properties); err != nil {
		return err
	}
	if o.DeleteNull {
		return w.WriteBool(true)
	}
	return nil
}

func encodeComposite(w *msgp.Writer, code int, path, fullPath []string, operations []any, errInvalid error) error {
	if err := writeHeader(w, 3, code); err != nil {
		return err
	}
	if err := encodePath(w, path); err != nil {
		return err
	}
	if err := w.WriteArrayHeader(uint32(len(operations))); err != nil { //nolint:gosec // operation count is bounded by practical limits.
		return err
	}
	for _, candidate := range operations {
		predicate, ok := candidate.(internal.PredicateOp)
		if !ok {
			return errInvalid
		}
		if err := encodeOpWithParent(w, predicate, fullPath); err != nil {
			return err
		}
	}
	return nil
}

func relativePath(parent, path []string) ([]string, error) {
	if len(parent) == 0 {
		return slices.Clone(path), nil
	}
	if len(path) < len(parent) || !slices.Equal(path[:len(parent)], parent) {
		return nil, op.ErrPredicatePathOutsideParent
	}
	return slices.Clone(path[len(parent):]), nil
}

// encodePath writes a path as a msgpack native array of string segments.
func encodePath(w *msgp.Writer, path []string) error {
	if err := w.WriteArrayHeader(uint32(len(path))); err != nil { //nolint:gosec // path length is bounded by practical limits.
		return err
	}
	for _, seg := range path {
		if err := w.WriteString(seg); err != nil {
			return err
		}
	}
	return nil
}
