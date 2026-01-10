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
	size64, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := uint32(size64)
	ops := make([]internal.Op, size)
	for i := uint32(0); i < size; i++ {
		op, err := decodeOp(reader)
		if err != nil {
			return nil, err
		}
		ops[i] = op
	}
	return ops, nil
}

func decodeOp(reader *msgp.Reader) (internal.Op, error) {
	// Operation size is not used for now, but we must read it.
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
	case internal.OpTestTypeCode:
		typesVal, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		types, ok := typesVal.([]interface{})
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
		return op.NewOpTestTypeOperationMultiple(path, strTypes), nil
	case internal.OpDefinedCode:
		return op.NewOpDefinedOperation(path), nil
	case internal.OpUndefinedCode:
		return op.NewOpUndefinedOperation(path), nil
	case internal.OpLessCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		floatVal, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("%w: less value must be a number, got %T", ErrInvalidValueType, value)
		}
		return op.NewOpLessOperation(path, floatVal), nil
	case internal.OpMoreCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		floatVal, ok := value.(float64)
		if !ok {
			return nil, fmt.Errorf("%w: more value must be a number, got %T", ErrInvalidValueType, value)
		}
		return op.NewOpMoreOperation(path, floatVal), nil
	case internal.OpContainsCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: contains value must be a string, got %T", ErrInvalidValueType, value)
		}
		return op.NewOpContainsOperation(path, strVal), nil
	case internal.OpInCode:
		values, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		arrVal, ok := values.([]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: in values must be an array, got %T", ErrInvalidValueType, values)
		}
		return op.NewOpInOperation(path, arrVal), nil
	case internal.OpStartsCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: starts value must be a string, got %T", ErrInvalidValueType, value)
		}
		return op.NewOpStartsOperation(path, strVal), nil
	case internal.OpEndsCode:
		value, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		strVal, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("%w: ends value must be a string, got %T", ErrInvalidValueType, value)
		}
		return op.NewOpEndsOperation(path, strVal), nil
	case internal.OpMatchesCode:
		pattern, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := reader.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpMatchesOperation(path, pattern, ignoreCase, nil), nil
	case internal.OpTestStringCode:
		str, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpTestStringOperationWithPos(path, str, pos), nil
	case internal.OpTestStringLenCode:
		length, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		not, err := reader.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpTestStringLenOperationWithNot(path, length, not), nil
	case internal.OpTypeCode:
		expectedType, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewOpTypeOperation(path, expectedType), nil
	case internal.OpFlipCode:
		return op.NewOpFlipOperation(path), nil
	case internal.OpIncCode:
		inc, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpIncOperation(path, inc), nil
	case internal.OpStrInsCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		str, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewOpStrInsOperation(path, pos, str), nil
	case internal.OpStrDelCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		length, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpStrDelOperation(path, pos, length), nil
	case internal.OpSplitCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		props, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		return op.NewOpSplitOperation(path, pos, props), nil
	case internal.OpExtendCode:
		properties, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		propsMap, ok := properties.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: extend properties must be an object, got %T", ErrInvalidValueType, properties)
		}
		deleteNull, err := reader.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpExtendOperation(path, propsMap, deleteNull), nil
	case internal.OpMergeCode:
		pos, err := reader.ReadFloat64()
		if err != nil {
			return nil, err
		}
		props, err := decodeValue(reader)
		if err != nil {
			return nil, err
		}
		propsMap, ok := props.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%w: merge properties must be an object, got %T", ErrInvalidValueType, props)
		}
		return op.NewOpMergeOperation(path, pos, propsMap), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

func decodePath(reader *msgp.Reader) ([]string, error) {
	size64, err := reader.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := uint32(size64)
	path := make([]string, size)
	for i := uint32(0); i < size; i++ {
		segment, err := reader.ReadString()
		if err != nil {
			return nil, err
		}
		path[i] = segment
	}
	return path, nil
}

func decodeValue(reader *msgp.Reader) (interface{}, error) {
	v, err := reader.ReadIntf()
	if err != nil {
		return nil, err
	}
	return convertMap(v), nil
}

// convertMap recursively converts map[interface{}]interface{} to map[string]interface{}
func convertMap(v interface{}) interface{} {
	if m, ok := v.(map[interface{}]interface{}); ok {
		res := make(map[string]interface{})
		for key, val := range m {
			if k, ok := key.(string); ok {
				res[k] = convertMap(val)
			}
		}
		return res
	}
	if s, ok := v.([]interface{}); ok {
		for i, val := range s {
			s[i] = convertMap(val)
		}
	}
	return v
}
