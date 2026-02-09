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

// decodeFloat64Value reads an interface value and asserts it is float64.
func decodeFloat64Value(reader *msgp.Reader, field string) (float64, error) {
	value, err := decodeValue(reader)
	if err != nil {
		return 0, err
	}
	f, ok := value.(float64)
	if !ok {
		return 0, fmt.Errorf("%w: %s must be a number, got %T", ErrInvalidValueType, field, value)
	}
	return f, nil
}

// decodeStringValue reads an interface value and asserts it is string.
func decodeStringValue(reader *msgp.Reader, field string) (string, error) {
	value, err := decodeValue(reader)
	if err != nil {
		return "", err
	}
	s, ok := value.(string)
	if !ok {
		return "", fmt.Errorf("%w: %s must be a string, got %T", ErrInvalidValueType, field, value)
	}
	return s, nil
}

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

	switch code {
	// Standard RFC 6902
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

	// Predicate operations
	case internal.OpTestTypeCode:
		return decodeTestType(reader, path)

	case internal.OpDefinedCode:
		return op.NewDefined(path), nil

	case internal.OpUndefinedCode:
		return op.NewUndefined(path), nil

	case internal.OpLessCode:
		v, err := decodeFloat64Value(reader, "less value")
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, v), nil

	case internal.OpMoreCode:
		v, err := decodeFloat64Value(reader, "more value")
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, v), nil

	case internal.OpContainsCode:
		v, err := decodeStringValue(reader, "contains value")
		if err != nil {
			return nil, err
		}
		return op.NewContains(path, v), nil

	case internal.OpInCode:
		values, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		arr, ok := values.([]any)
		if !ok {
			return nil, fmt.Errorf("%w: in values must be an array, got %T", ErrInvalidValueType, values)
		}
		return op.NewIn(path, arr), nil

	case internal.OpStartsCode:
		v, err := decodeStringValue(reader, "starts value")
		if err != nil {
			return nil, err
		}
		return op.NewStarts(path, v), nil

	case internal.OpEndsCode:
		v, err := decodeStringValue(reader, "ends value")
		if err != nil {
			return nil, err
		}
		return op.NewEnds(path, v), nil

	case internal.OpMatchesCode:
		pattern, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := reader.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewMatches(path, pattern, ignoreCase, nil), nil

	case internal.OpTestStringCode:
		str, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewTestStringWithPos(path, str, pos), nil

	case internal.OpTestStringLenCode:
		length, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		not, err := reader.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewTestStringLenWithNot(path, length, not), nil

	case internal.OpTypeCode:
		expectedType, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewType(path, expectedType), nil

	// Extended operations
	case internal.OpFlipCode:
		return op.NewFlip(path), nil

	case internal.OpIncCode:
		inc, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewInc(path, inc), nil

	case internal.OpStrInsCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		str, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewStrIns(path, pos, str), nil

	case internal.OpStrDelCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		length, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewStrDel(path, pos, length), nil

	case internal.OpSplitCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		props, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewSplit(path, pos, props), nil

	case internal.OpExtendCode:
		return decodeExtend(reader, path)

	case internal.OpMergeCode:
		return decodeMerge(reader, path)

	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

// decodeTestType decodes a test_type operation from the reader.
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

// decodeExtend decodes an extend operation from the reader.
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

// decodeMerge decodes a merge operation from the reader.
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

func decodeValue(reader *msgp.Reader) (any, error) {
	v, err := reader.ReadIntf()
	if err != nil {
		return nil, err
	}
	return convertMap(v), nil
}

// convertMap recursively converts map[any]any to map[string]any.
func convertMap(v any) any {
	switch m := v.(type) {
	case map[any]any:
		res := make(map[string]any, len(m))
		for key, val := range m {
			if k, ok := key.(string); ok {
				res[k] = convertMap(val)
			}
		}
		return res
	case []any:
		for i, val := range m {
			m[i] = convertMap(val)
		}
		return m
	default:
		return v
	}
}
