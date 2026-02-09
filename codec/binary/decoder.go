package binary

import (
	"bytes"
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/tinylib/msgp/msgp"
)

// decode decodes a slice of operations from binary format.
func (c *Codec) decode(data []byte) ([]internal.Op, error) {
	reader := msgp.NewReader(bytes.NewReader(data))
	return decodeOps(reader)
}

// decodeOps reads the operation count and decodes each operation sequentially.
func decodeOps(reader *msgp.Reader) ([]internal.Op, error) {
	count, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := int(count)
	ops := make([]internal.Op, size)
	for i := range size {
		decoded, err := decodeOp(reader)
		if err != nil {
			return nil, err
		}
		ops[i] = decoded
	}
	return ops, nil
}

// decodeOp reads the array header, operation code, and path, then dispatches
// to the appropriate sub-function by operation category.
func decodeOp(reader *msgp.Reader) (internal.Op, error) {
	// Read and discard array header (size is implicit from operation type).
	if _, err := reader.ReadArrayHeader(); err != nil {
		return nil, err
	}

	code, err := reader.ReadUint8()
	if err != nil {
		return nil, err
	}

	path, err := decodePath(reader)
	if err != nil {
		return nil, err
	}

	switch {
	case isCoreOp(code):
		return decodeCoreOp(reader, code, path)
	case isPredicateOp(code):
		return decodePredicateOp(reader, code, path)
	case isExtendedOp(code):
		return decodeExtendedOp(reader, code, path)
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

// --- Operation category checks ---

func isCoreOp(code uint8) bool {
	switch code {
	case internal.OpAddCode, internal.OpRemoveCode, internal.OpReplaceCode,
		internal.OpMoveCode, internal.OpCopyCode, internal.OpTestCode:
		return true
	default:
		return false
	}
}

func isPredicateOp(code uint8) bool {
	switch code {
	case internal.OpTestTypeCode, internal.OpDefinedCode, internal.OpUndefinedCode,
		internal.OpLessCode, internal.OpMoreCode, internal.OpContainsCode,
		internal.OpInCode, internal.OpStartsCode, internal.OpEndsCode,
		internal.OpMatchesCode, internal.OpTestStringCode, internal.OpTestStringLenCode,
		internal.OpTypeCode:
		return true
	default:
		return false
	}
}

func isExtendedOp(code uint8) bool {
	switch code {
	case internal.OpFlipCode, internal.OpIncCode,
		internal.OpStrInsCode, internal.OpStrDelCode,
		internal.OpSplitCode, internal.OpExtendCode, internal.OpMergeCode:
		return true
	default:
		return false
	}
}

// --- Core RFC 6902 operations ---

// decodeCoreOp decodes standard JSON Patch (RFC 6902) operations.
func decodeCoreOp(reader *msgp.Reader, code uint8, path []string) (internal.Op, error) {
	switch code {
	case internal.OpAddCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewAdd(path, value), nil

	case internal.OpRemoveCode:
		return op.NewRemove(path), nil

	case internal.OpReplaceCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewReplace(path, value), nil

	case internal.OpMoveCode:
		from, err := decodePath(reader)
		if err != nil {
			return nil, err
		}
		return op.NewMove(from, path), nil

	case internal.OpCopyCode:
		from, err := decodePath(reader)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(from, path), nil

	case internal.OpTestCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewTest(path, value), nil

	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

// --- Predicate operations ---

// decodePredicateOp decodes JSON Predicate operations.
func decodePredicateOp(reader *msgp.Reader, code uint8, path []string) (internal.Op, error) {
	switch code {
	case internal.OpDefinedCode:
		return op.NewDefined(path), nil

	case internal.OpUndefinedCode:
		return op.NewUndefined(path), nil

	case internal.OpTestTypeCode:
		return decodeTestType(reader, path)

	case internal.OpLessCode:
		v, err := readFloat64Value(reader)
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, v), nil

	case internal.OpMoreCode:
		v, err := readFloat64Value(reader)
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, v), nil

	case internal.OpContainsCode:
		v, err := readStringValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewContains(path, v), nil

	case internal.OpStartsCode:
		v, err := readStringValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewStarts(path, v), nil

	case internal.OpEndsCode:
		v, err := readStringValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewEnds(path, v), nil

	case internal.OpInCode:
		return decodeIn(reader, path)

	case internal.OpMatchesCode:
		return decodeMatches(reader, path)

	case internal.OpTestStringCode:
		return decodeTestString(reader, path)

	case internal.OpTestStringLenCode:
		return decodeTestStringLen(reader, path)

	case internal.OpTypeCode:
		return decodeType(reader, path)

	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

// --- Extended operations ---

// decodeExtendedOp decodes extended operations (flip, inc, str_ins, str_del, split, extend, merge).
func decodeExtendedOp(reader *msgp.Reader, code uint8, path []string) (internal.Op, error) {
	switch code {
	case internal.OpFlipCode:
		return op.NewFlip(path), nil

	case internal.OpIncCode:
		inc, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewInc(path, inc), nil

	case internal.OpStrInsCode:
		return decodeStrIns(reader, path)

	case internal.OpStrDelCode:
		return decodeStrDel(reader, path)

	case internal.OpSplitCode:
		return decodeSplit(reader, path)

	case internal.OpExtendCode:
		return decodeExtend(reader, path)

	case internal.OpMergeCode:
		return decodeMerge(reader, path)

	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

// --- Typed value readers ---

// readFloat64Value reads an interface value and asserts it is float64.
func readFloat64Value(reader *msgp.Reader) (float64, error) {
	value, err := decodeValue(reader)
	if err != nil {
		return 0, err
	}
	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("%w: expected number, got %T", ErrInvalidValueType, value)
	}
	return f, nil
}

// readStringValue reads an interface value and asserts it is string.
func readStringValue(reader *msgp.Reader) (string, error) {
	value, err := decodeValue(reader)
	if err != nil {
		return "", err
	}
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%w: expected string, got %T", ErrInvalidValueType, value)
	}
	return s, nil
}

// --- Operation-specific decoders ---

// decodeTestType decodes a test_type operation.
func decodeTestType(reader *msgp.Reader, path []string) (internal.Op, error) {
	typesVal, err := decodeValue(reader)
	if err != nil {
		return nil, err
	}
	types, ok := typesVal.([]any)
	if !ok {
		return nil, ErrInvalidTestTypeFormat
	}
	strTypes := make([]string, len(types))
	for i, v := range types {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("%w: expected string at index %d, got %T", ErrInvalidTestTypeFormat, i, v)
		}
		strTypes[i] = str
	}
	return op.NewTestTypeMultiple(path, strTypes), nil
}

// decodeIn decodes an in predicate operation.
func decodeIn(reader *msgp.Reader, path []string) (internal.Op, error) {
	values, err := decodeValue(reader)
	if err != nil {
		return nil, err
	}
	arr, ok := values.([]any)
	if !ok {
		return nil, fmt.Errorf("%w: in values must be an array, got %T", ErrInvalidValueType, values)
	}
	return op.NewIn(path, arr), nil
}

// decodeMatches decodes a matches predicate operation.
func decodeMatches(reader *msgp.Reader, path []string) (internal.Op, error) {
	pattern, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	ignoreCase, err := reader.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewMatches(path, pattern, ignoreCase, nil), nil
}

// decodeTestString decodes a test_string operation.
func decodeTestString(reader *msgp.Reader, path []string) (internal.Op, error) {
	str, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	pos, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return op.NewTestStringWithPos(path, str, pos), nil
}

// decodeTestStringLen decodes a test_string_len operation.
func decodeTestStringLen(reader *msgp.Reader, path []string) (internal.Op, error) {
	length, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	not, err := reader.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewTestStringLenWithNot(path, length, not), nil
}

// decodeType decodes a type predicate operation.
func decodeType(reader *msgp.Reader, path []string) (internal.Op, error) {
	expectedType, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	return op.NewType(path, expectedType), nil
}

// decodeStrIns decodes a str_ins operation.
func decodeStrIns(reader *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	str, err := reader.ReadString()
	if err != nil {
		return nil, err
	}
	return op.NewStrIns(path, pos, str), nil
}

// decodeStrDel decodes a str_del operation.
func decodeStrDel(reader *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	length, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	return op.NewStrDel(path, pos, length), nil
}

// decodeSplit decodes a split operation.
func decodeSplit(reader *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	props, err := decodeValue(reader)
	if err != nil {
		return nil, err
	}
	return op.NewSplit(path, pos, props), nil
}

// decodeExtend decodes an extend operation.
func decodeExtend(reader *msgp.Reader, path []string) (internal.Op, error) {
	properties, err := decodeValue(reader)
	if err != nil {
		return nil, err
	}
	propsMap, ok := properties.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: extend properties must be an object, got %T", ErrInvalidValueType, properties)
	}
	deleteNull, err := reader.ReadBool()
	if err != nil {
		return nil, err
	}
	return op.NewExtend(path, propsMap, deleteNull), nil
}

// decodeMerge decodes a merge operation.
func decodeMerge(reader *msgp.Reader, path []string) (internal.Op, error) {
	pos, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	props, err := decodeValue(reader)
	if err != nil {
		return nil, err
	}
	propsMap, ok := props.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: merge properties must be an object, got %T", ErrInvalidValueType, props)
	}
	return op.NewMerge(path, pos, propsMap), nil
}

// --- Low-level decoders ---

// decodePath reads a path as a float64 count followed by string segments.
func decodePath(reader *msgp.Reader) ([]string, error) {
	count, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := int(count)
	path := make([]string, size)
	for i := range size {
		segment, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		path[i] = segment
	}
	return path, nil
}

// decodeValue reads an arbitrary msgp value and normalizes map types.
func decodeValue(reader *msgp.Reader) (any, error) {
	v, err := reader.ReadIntf()
	if err != nil {
		return nil, err
	}
	return normalizeMap(v), nil
}

// normalizeMap recursively converts map[any]any to map[string]any.
func normalizeMap(v any) any {
	switch m := v.(type) {
	case map[any]any:
		res := make(map[string]any, len(m))
		for key, val := range m {
			if k, ok := key.(string); ok {
				res[k] = normalizeMap(val)
			}
		}
		return res
	case []any:
		for i, val := range m {
			m[i] = normalizeMap(val)
		}
		return m
	default:
		return v
	}
}
