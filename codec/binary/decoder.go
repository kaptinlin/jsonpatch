package binary

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	msgpack "github.com/wapc/tinygo-msgpack"
)

// decode decodes a slice of operations from binary format.
func (c *Codec) decode(data []byte) ([]internal.Op, error) {
	decoder := msgpack.NewDecoder(data)
	return decodeOps(&decoder)
}

func decodeOps(decoder *msgpack.Decoder) ([]internal.Op, error) {
	size64, err := decoder.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := uint32(size64)
	ops := make([]internal.Op, size)
	for i := uint32(0); i < size; i++ {
		op, err := decodeOp(decoder)
		if err != nil {
			return nil, err
		}
		ops[i] = op
	}
	return ops, nil
}

func decodeOp(decoder *msgpack.Decoder) (internal.Op, error) {
	// Operation size is not used for now, but we must read it.
	if _, err := decoder.ReadArraySize(); err != nil {
		return nil, err
	}

	code, err := decoder.ReadUint8()
	if err != nil {
		return nil, err
	}

	path, err := decodePath(decoder)
	if err != nil {
		return nil, err
	}

	switch code {
	case internal.OpAddCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewAdd(path, value), nil
	case internal.OpRemoveCode:
		return op.NewRemove(path), nil
	case internal.OpReplaceCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewReplace(path, value), nil
	case internal.OpMoveCode:
		from, err := decodePath(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewMove(from, path), nil
	case internal.OpCopyCode:
		from, err := decodePath(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(from, path), nil
	case internal.OpTestCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewTest(path, value), nil
	case internal.OpTestTypeCode:
		typesVal, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		types, ok := typesVal.([]interface{})
		if !ok {
			return nil, ErrInvalidTestTypeFormat
		}
		strTypes := make([]string, len(types))
		for i, v := range types {
			strTypes[i] = v.(string)
		}
		return op.NewOpTestTypeOperationMultiple(path, strTypes), nil
	case internal.OpDefinedCode:
		return op.NewOpDefinedOperation(path), nil
	case internal.OpUndefinedCode:
		not, err := decoder.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpUndefinedOperation(path, not), nil
	case internal.OpLessCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpLessOperation(path, value.(float64)), nil
	case internal.OpMoreCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpMoreOperation(path, value.(float64)), nil
	case internal.OpContainsCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpContainsOperation(path, value.(string)), nil
	case internal.OpInCode:
		values, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpInOperation(path, values.([]interface{})), nil
	case internal.OpStartsCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpStartsOperation(path, value.(string)), nil
	case internal.OpEndsCode:
		value, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpEndsOperation(path, value.(string)), nil
	case internal.OpMatchesCode:
		pattern, err := decoder.ReadString()
		if err != nil {
			return nil, err
		}
		ignoreCase, err := decoder.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpMatchesOperation(path, pattern, ignoreCase)
	case internal.OpTestStringCode:
		str, err := decoder.ReadString()
		if err != nil {
			return nil, err
		}
		pos, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpTestStringOperationWithPos(path, str, pos), nil
	case internal.OpTestStringLenCode:
		length, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		not, err := decoder.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpTestStringLenOperationWithNot(path, length, not), nil
	case internal.OpTypeCode:
		expectedType, err := decoder.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewOpTypeOperation(path, expectedType), nil
	case internal.OpFlipCode:
		return op.NewOpFlipOperation(path), nil
	case internal.OpIncCode:
		inc, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpIncOperation(path, inc), nil
	case internal.OpStrInsCode:
		pos, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		str, err := decoder.ReadString()
		if err != nil {
			return nil, err
		}
		return op.NewOpStrInsOperation(path, pos, str), nil
	case internal.OpStrDelCode:
		pos, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		length, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		return op.NewOpStrDelOperation(path, pos, length), nil
	case internal.OpSplitCode:
		pos, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		props, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpSplitOperation(path, pos, props), nil
	case internal.OpExtendCode:
		properties, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		deleteNull, err := decoder.ReadBool()
		if err != nil {
			return nil, err
		}
		return op.NewOpExtendOperation(path, properties.(map[string]interface{}), deleteNull), nil
	case internal.OpMergeCode:
		pos, err := decoder.ReadFloat64()
		if err != nil {
			return nil, err
		}
		props, err := decodeValue(decoder)
		if err != nil {
			return nil, err
		}
		return op.NewOpMergeOperation(path, pos, props.(map[string]interface{})), nil
	default:
		return nil, fmt.Errorf("%w: %d", ErrUnsupportedOp, code)
	}
}

func decodePath(decoder *msgpack.Decoder) ([]string, error) {
	size64, err := decoder.ReadFloat64()
	if err != nil {
		return nil, err
	}
	size := uint32(size64)
	path := make([]string, size)
	for i := uint32(0); i < size; i++ {
		segment, err := decoder.ReadString()
		if err != nil {
			return nil, err
		}
		path[i] = segment
	}
	return path, nil
}

func decodeValue(decoder *msgpack.Decoder) (interface{}, error) {
	v, err := decoder.ReadAny()
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
