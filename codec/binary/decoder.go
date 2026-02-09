package binary

import (
	"fmt"

	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

// decodeOps reads the operation count and decodes each operation.
func decodeOps(r *msgp.Reader) ([]internal.Op, error) {
	count, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := int(count)
	ops := make([]internal.Op, size)
	for i := range size {
		decoded, err := decodeOp(r)
		if err != nil {
			return nil, err
		}
		ops[i] = decoded
	}
	return ops, nil
}

// decodeOp reads the array header, operation code, and path,
// then dispatches to the appropriate decoder.
func decodeOp(r *msgp.Reader) (internal.Op, error) {
	if _, err := r.ReadArrayHeader(); err != nil {
		return nil, err
	}
	code, err := r.ReadUint8()
	if err != nil {
		return nil, err
	}
	path, err := decodePath(r)
	if err != nil {
		return nil, err
	}

	switch code {
	// Standard RFC 6902
	case internal.OpAddCode:
		value, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		return op.NewAdd(path, value), nil
	case internal.OpRemoveCode:
		return op.NewRemove(path), nil
	case internal.OpReplaceCode:
		value, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		return op.NewReplace(path, value), nil
	case internal.OpMoveCode:
		from, err := decodePath(r)
		if err != nil {
			return nil, err
		}
		return op.NewMove(from, path), nil
	case internal.OpCopyCode:
		from, err := decodePath(r)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(from, path), nil
	case internal.OpTestCode:
		value, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		return op.NewTest(path, value), nil

	// Predicate operations
	case internal.OpDefinedCode:
		return op.NewDefined(path), nil
	case internal.OpUndefinedCode:
		return op.NewUndefined(path), nil
	case internal.OpTestTypeCode:
		return decodeTestType(r, path)
	case internal.OpLessCode:
		v, err := readFloat64(r)
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, v), nil
	case internal.OpMoreCode:
		v, err := readFloat64(r)
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, v), nil
	case internal.OpContainsCode:
		v, err := readString(r)
		if err != nil {
			return nil, err
		}
		return op.NewContains(path, v), nil
	case internal.OpStartsCode:
		v, err := readString(r)
		if err != nil {
			return nil, err
		}
		return op.NewStarts(path, v), nil
	case internal.OpEndsCode:
		v, err := readString(r)
		if err != nil {
			return nil, err
		}
		return op.NewEnds(path, v), nil
	case internal.OpInCode:
		return decodeIn(r, path)
	case internal.OpMatchesCode:
		return decodeMatches(r, path)
	case internal.OpTestStringCode:
		return decodeTestString(r, path)
	case internal.OpTestStringLenCode:
		return decodeTestStringLen(r, path)
	case internal.OpTypeCode:
		return decodeType(r, path)

	// Extended operations
	case internal.OpFlipCode:
		return op.NewFlip(path), nil
	case internal.OpIncCode:
		inc, err := r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewInc(path, inc), nil
	case internal.OpStrInsCode:
		return decodeStrIns(r, path)
	case internal.OpStrDelCode:
		return decodeStrDel(r, path)
	case internal.OpSplitCode:
		return decodeSplit(r, path)
	case internal.OpExtendCode:
		return decodeExtend(r, path)
	case internal.OpMergeCode:
		return decodeMerge(r, path)

	default:
		return nil, fmt.Errorf("unsupported op code %d: %w",
			code, ErrUnsupportedOp)
	}
}

// readFloat64 reads an interface value and asserts it is float64.
func readFloat64(r *msgp.Reader) (float64, error) {
	value, err := decodeValue(r)
	if err != nil {
		return 0, err
	}
	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("expected number, got %T: %w", value, ErrInvalidValueType)
	}
	return f, nil
}

// readString reads an interface value and asserts it is string.
func readString(r *msgp.Reader) (string, error) {
	value, err := decodeValue(r)
	if err != nil {
		return "", err
	}
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("expected string, got %T: %w", value, ErrInvalidValueType)
	}
	return s, nil
}

// decodeTestType decodes a test_type operation.
func decodeTestType(r *msgp.Reader, path []string) (internal.Op, error) {
	raw, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	types, ok := raw.([]any)
	if !ok {
		return nil, ErrInvalidTestTypeFormat
	}
	strs := make([]string, len(types))
	for i, v := range types {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("expected string at index %d, got %T: %w", i, v, ErrInvalidTestTypeFormat)
		}
		strs[i] = str
	}
	return op.NewTestTypeMultiple(path, strs), nil
}

// decodeIn decodes an in predicate operation.
func decodeIn(r *msgp.Reader, path []string) (internal.Op, error) {
	raw, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, fmt.Errorf("in values must be an array, got %T: %w", raw, ErrInvalidValueType)
	}
	return op.NewIn(path, arr), nil
}

// decodeMatches decodes a matches predicate operation.
func decodeMatches(r *msgp.Reader, path []string) (internal.Op, error) {
	pattern, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	ignoreCase, err := r.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewMatches(path, pattern, ignoreCase, nil), nil
}

// decodeTestString decodes a test_string operation.
func decodeTestString(r *msgp.Reader, path []string) (internal.Op, error) {
	str, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return op.NewTestString(path, str, pos, false, false), nil
}

// decodeTestStringLen decodes a test_string_len operation.
func decodeTestStringLen(r *msgp.Reader, path []string) (internal.Op, error) {
	length, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	not, err := r.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewTestStringLenWithNot(path, length, not), nil
}

// decodeType decodes a type predicate operation.
func decodeType(r *msgp.Reader, path []string) (internal.Op, error) {
	expected, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	return op.NewType(path, expected), nil
}

// decodeStrIns decodes a str_ins operation.
func decodeStrIns(r *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	str, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	return op.NewStrIns(path, pos, str), nil
}

// decodeStrDel decodes a str_del operation.
func decodeStrDel(r *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	length, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return op.NewStrDel(path, pos, length), nil
}

// decodeSplit decodes a split operation.
func decodeSplit(r *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	props, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	return op.NewSplit(path, pos, props), nil
}

// decodeExtend decodes an extend operation.
func decodeExtend(r *msgp.Reader, path []string) (internal.Op, error) {
	raw, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	props, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("extend properties must be an object, got %T: %w", raw, ErrInvalidValueType)
	}
	deleteNull, err := r.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewExtend(path, props, deleteNull), nil
}

// decodeMerge decodes a merge operation.
func decodeMerge(r *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	raw, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	props, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("merge properties must be an object, got %T: %w", raw, ErrInvalidValueType)
	}
	return op.NewMerge(path, pos, props), nil
}

// decodePath reads a path as a float64 count followed by string segments.
func decodePath(r *msgp.Reader) ([]string, error) {
	count, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := int(count)
	path := make([]string, size)
	for i := range size {
		seg, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		path[i] = seg
	}
	return path, nil
}

// decodeValue reads an arbitrary msgp value and normalizes map types.
func decodeValue(r *msgp.Reader) (any, error) {
	v, err := r.ReadIntf()
	if err != nil {
		return nil, err
	}
	return normalizeMap(v), nil
}

