package binary

import (
	"fmt"
	"slices"

	"github.com/tinylib/msgp/msgp"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

// decodeOps reads the operation count and decodes each operation.
func decodeOps(r *msgp.Reader) ([]internal.Op, error) {
	arrSize, err := r.ReadArrayHeader()
	if err != nil {
		return nil, err
	}
	size := int(arrSize)
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
	return decodeOpWithParent(r, nil)
}

func decodeOpWithParent(r *msgp.Reader, parent []string) (internal.Op, error) {
	arrSize, err := r.ReadArrayHeader()
	if err != nil {
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
	if parent != nil {
		path = mergePaths(parent, path)
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
		if arrSize >= 3 {
			oldValue, err := decodeValue(r)
			if err != nil {
				return nil, err
			}
			return op.NewRemoveWithOldValue(path, oldValue), nil
		}
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
		return op.NewMove(path, from), nil
	case internal.OpCopyCode:
		from, err := decodePath(r)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(path, from), nil
	case internal.OpTestCode:
		value, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		not, err := decodeOptionalBool(r, arrSize, 4)
		if err != nil {
			return nil, err
		}
		return op.NewTestWithNot(path, value, not), nil

	// Predicate operations
	case internal.OpDefinedCode:
		return op.NewDefined(path), nil
	case internal.OpUndefinedCode:
		return op.NewUndefined(path), nil
	case internal.OpTestTypeCode:
		return decodeTestType(r, path)
	case internal.OpLessCode:
		v, err := r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, v), nil
	case internal.OpMoreCode:
		v, err := r.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, v), nil
	case internal.OpContainsCode:
		v, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := decodeOptionalBool(r, arrSize, 4)
		if err != nil {
			return nil, err
		}
		return op.NewContainsWithIgnoreCase(path, v, ignoreCase), nil
	case internal.OpStartsCode:
		v, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := decodeOptionalBool(r, arrSize, 4)
		if err != nil {
			return nil, err
		}
		return op.NewStartsWithIgnoreCase(path, v, ignoreCase), nil
	case internal.OpEndsCode:
		v, err := r.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := decodeOptionalBool(r, arrSize, 4)
		if err != nil {
			return nil, err
		}
		return op.NewEndsWithIgnoreCase(path, v, ignoreCase), nil
	case internal.OpInCode:
		return decodeIn(r, path)
	case internal.OpMatchesCode:
		return decodeMatches(r, path, arrSize)
	case internal.OpTestStringCode:
		return decodeTestString(r, path, arrSize)
	case internal.OpTestStringLenCode:
		return decodeTestStringLen(r, path, arrSize)
	case internal.OpTypeCode:
		return decodeType(r, path)
	case internal.OpAndCode:
		return decodeComposite(r, path, internal.OpAndType)
	case internal.OpOrCode:
		return decodeComposite(r, path, internal.OpOrType)
	case internal.OpNotCode:
		return decodeComposite(r, path, internal.OpNotType)

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
		return decodeSplit(r, path, arrSize)
	case internal.OpExtendCode:
		return decodeExtend(r, path, arrSize)
	case internal.OpMergeCode:
		return decodeMerge(r, path, arrSize)

	default:
		return nil, fmt.Errorf("unsupported op code %d: %w",
			code, ErrUnsupportedOp)
	}
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
func decodeMatches(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	pattern, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	ignoreCase, err := decodeOptionalBool(r, arrSize, 4)
	if err != nil {
		return nil, err
	}
	return op.NewMatches(path, pattern, ignoreCase, nil), nil
}

// decodeTestString decodes a test_string operation.
func decodeTestString(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	str, err := r.ReadString()
	if err != nil {
		return nil, err
	}
	not, err := decodeOptionalBool(r, arrSize, 5)
	if err != nil {
		return nil, err
	}
	return op.NewTestString(path, str, pos, not, false), nil
}

// decodeTestStringLen decodes a test_string_len operation.
func decodeTestStringLen(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	length, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	not, err := decodeOptionalBool(r, arrSize, 4)
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
func decodeSplit(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	var props any
	if arrSize >= 4 {
		props, err = decodeValue(r)
		if err != nil {
			return nil, err
		}
	}
	return op.NewSplit(path, pos, props), nil
}

// decodeExtend decodes an extend operation.
func decodeExtend(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	raw, err := decodeValue(r)
	if err != nil {
		return nil, err
	}
	props, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("extend properties must be an object, got %T: %w", raw, ErrInvalidValueType)
	}
	deleteNull, err := decodeOptionalBool(r, arrSize, 4)
	if err != nil {
		return nil, err
	}
	return op.NewExtend(path, props, deleteNull), nil
}

func decodeComposite(r *msgp.Reader, path []string, opType internal.OpType) (internal.Op, error) {
	ops, err := decodePredicateOps(r, path)
	if err != nil {
		return nil, err
	}

	switch opType {
	case internal.OpAndType:
		return op.NewAnd(path, ops), nil
	case internal.OpOrType:
		return op.NewOr(path, ops), nil
	default:
		if len(ops) != 1 {
			return nil, ErrNotSinglePredicate
		}
		return op.NewNotMultiple(path, ops), nil
	}
}

func decodePredicateOps(r *msgp.Reader, parent []string) ([]any, error) {
	size, err := r.ReadArrayHeader()
	if err != nil {
		return nil, err
	}
	ops := make([]any, int(size))
	for i := range ops {
		decoded, err := decodeOpWithParent(r, parent)
		if err != nil {
			return nil, err
		}
		predicate, ok := decoded.(internal.PredicateOp)
		if !ok {
			return nil, ErrInvalidPredicate
		}
		ops[i] = predicate
	}
	return ops, nil
}

func decodeOptionalBool(r *msgp.Reader, arrSize, presentAt uint32) (bool, error) {
	if arrSize < presentAt {
		return false, nil
	}
	return r.ReadBool()
}

// decodeMerge decodes a merge operation.
func decodeMerge(r *msgp.Reader, path []string, arrSize uint32) (internal.Op, error) {
	pos, err := r.ReadFloat64()
	if err != nil {
		return nil, err
	}
	var props map[string]any
	if arrSize >= 4 {
		raw, err := decodeValue(r)
		if err != nil {
			return nil, err
		}
		var ok bool
		props, ok = raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("merge properties must be an object, got %T: %w", raw, ErrInvalidValueType)
		}
	}
	return op.NewMerge(path, pos, props), nil
}

// decodePath reads a path as a msgpack native array of string segments.
func decodePath(r *msgp.Reader) ([]string, error) {
	size, err := r.ReadArrayHeader()
	if err != nil {
		return nil, err
	}
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

func mergePaths(base, child []string) []string {
	if len(child) == 0 {
		return slices.Clone(base)
	}
	if slices.Equal(base, child) {
		return slices.Clone(child)
	}
	return slices.Concat(base, child)
}

// decodeValue reads an arbitrary msgp value and normalizes map types.
func decodeValue(r *msgp.Reader) (any, error) {
	v, err := r.ReadIntf()
	if err != nil {
		return nil, err
	}
	return normalizeMap(v), nil
}
