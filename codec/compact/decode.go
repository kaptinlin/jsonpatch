package compact

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Pre-built mapping tables for better performance and maintainability
var (
	numericToOpType = map[int]internal.OpType{
		0:  internal.OpAddType,
		1:  internal.OpRemoveType,
		2:  internal.OpReplaceType,
		3:  internal.OpCopyType,
		4:  internal.OpMoveType,
		5:  internal.OpTestType,
		6:  internal.OpStrInsType,
		7:  internal.OpStrDelType,
		8:  internal.OpFlipType,
		9:  internal.OpIncType,
		10: internal.OpSplitType,
		11: internal.OpMergeType,
		12: internal.OpExtendType,
		30: internal.OpContainsType,
		31: internal.OpDefinedType,
		32: internal.OpEndsType,
		33: internal.OpInType,
		34: internal.OpLessType,
		35: internal.OpMatchesType,
		36: internal.OpMoreType,
		37: internal.OpStartsType,
		38: internal.OpUndefinedType,
		39: internal.OpTestTypeType,
		40: internal.OpTestStringType,
		41: internal.OpTestStringLenType,
		42: internal.OpTypeType,
		43: internal.OpAndType,
		44: internal.OpNotType,
		45: internal.OpOrType,
	}

	stringToOpType = map[string]internal.OpType{
		"add":             internal.OpAddType,
		"remove":          internal.OpRemoveType,
		"replace":         internal.OpReplaceType,
		"copy":            internal.OpCopyType,
		"move":            internal.OpMoveType,
		"test":            internal.OpTestType,
		"str_ins":         internal.OpStrInsType,
		"str_del":         internal.OpStrDelType,
		"flip":            internal.OpFlipType,
		"inc":             internal.OpIncType,
		"split":           internal.OpSplitType,
		"merge":           internal.OpMergeType,
		"extend":          internal.OpExtendType,
		"contains":        internal.OpContainsType,
		"defined":         internal.OpDefinedType,
		"ends":            internal.OpEndsType,
		"in":              internal.OpInType,
		"less":            internal.OpLessType,
		"matches":         internal.OpMatchesType,
		"more":            internal.OpMoreType,
		"starts":          internal.OpStartsType,
		"undefined":       internal.OpUndefinedType,
		"test_type":       internal.OpTestTypeType,
		"test_string":     internal.OpTestStringType,
		"test_string_len": internal.OpTestStringLenType,
		"type":            internal.OpTypeType,
		"and":             internal.OpAndType,
		"not":             internal.OpNotType,
		"or":              internal.OpOrType,
	}
)

// Decode decodes compact format operations using default options
func Decode(compactOps []Op, opts ...DecoderOption) ([]internal.Op, error) {
	decoder := NewDecoder(opts...)
	return decoder.DecodeSlice(compactOps)
}

// DecodeJSON decodes compact format JSON bytes into operations
func DecodeJSON(data []byte, opts ...DecoderOption) ([]internal.Op, error) {
	var compactOps []Op
	if err := json.Unmarshal(data, &compactOps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal compact operations: %w", err)
	}
	return Decode(compactOps, opts...)
}

// compactToOp converts a compact operation to an operation instance
func compactToOp(compactOp Op, opts DecoderOptions) (internal.Op, error) {
	if len(compactOp) < 2 {
		return nil, ErrCompactOperationMinLength
	}

	// Parse path
	pathStr, ok := compactOp[1].(string)
	if !ok {
		return nil, ErrCompactOperationPathNotString
	}
	path := stringToPath(pathStr)

	// Determine operation type from opcode
	opType, err := getOpTypeFromOpcode(compactOp[0])
	if err != nil {
		return nil, err
	}

	// Create operation based on type
	switch opType {
	case internal.OpAddType:
		if len(compactOp) < 3 {
			return nil, ErrAddOperationRequiresValue
		}
		return op.NewAdd(path, compactOp[2]), nil

	case internal.OpRemoveType:
		// Optional oldValue parameter
		if len(compactOp) >= 3 {
			return op.NewRemoveWithOldValue(path, compactOp[2]), nil
		}
		return op.NewRemove(path), nil

	case internal.OpReplaceType:
		if len(compactOp) < 3 {
			return nil, ErrReplaceOperationRequiresValue
		}
		// Optional oldValue parameter
		if len(compactOp) >= 4 {
			return op.NewReplaceWithOldValue(path, compactOp[2], compactOp[3]), nil
		}
		return op.NewReplace(path, compactOp[2]), nil

	case internal.OpMoveType:
		if len(compactOp) < 3 {
			return nil, ErrMoveOperationRequiresFrom
		}
		fromStr, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrMoveOperationFromNotString
		}
		from := stringToPath(fromStr)
		return op.NewMove(path, from), nil

	case internal.OpCopyType:
		if len(compactOp) < 3 {
			return nil, ErrCopyOperationRequiresFrom
		}
		fromStr, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrCopyOperationFromNotString
		}
		from := stringToPath(fromStr)
		return op.NewCopy(path, from), nil

	case internal.OpTestType:
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		not := false
		if len(compactOp) >= 4 {
			not = toBool(compactOp[3])
		}
		return op.NewOpTestOperationWithNot(path, compactOp[2], not), nil

	case internal.OpFlipType:
		return op.NewOpFlipOperation(path), nil

	case internal.OpIncType:
		if len(compactOp) < 3 {
			return nil, ErrIncOperationRequiresDelta
		}
		delta, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrIncOperationDeltaNotNumber
		}
		return op.NewOpIncOperation(path, delta), nil

	case internal.OpDefinedType:
		return op.NewOpDefinedOperation(path), nil

	case internal.OpUndefinedType:
		return op.NewOpUndefinedOperation(path), nil

	case internal.OpContainsType:
		if len(compactOp) < 3 {
			return nil, ErrContainsOperationRequiresValue
		}
		value, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrContainsOperationValueNotString
		}
		ignoreCase := false
		if len(compactOp) >= 4 {
			ignoreCase = toBool(compactOp[3])
		}
		return op.NewOpContainsOperationWithIgnoreCase(path, value, ignoreCase), nil

	case internal.OpStartsType:
		if len(compactOp) < 3 {
			return nil, ErrStartsOperationRequiresValue
		}
		value, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrStartsOperationValueNotString
		}
		ignoreCase := false
		if len(compactOp) >= 4 {
			ignoreCase = toBool(compactOp[3])
		}
		return op.NewOpStartsOperationWithIgnoreCase(path, value, ignoreCase), nil

	case internal.OpEndsType:
		if len(compactOp) < 3 {
			return nil, ErrEndsOperationRequiresValue
		}
		value, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrEndsOperationValueNotString
		}
		ignoreCase := false
		if len(compactOp) >= 4 {
			ignoreCase = toBool(compactOp[3])
		}
		return op.NewOpEndsOperationWithIgnoreCase(path, value, ignoreCase), nil

	case internal.OpTypeType:
		if len(compactOp) < 3 {
			return nil, ErrTypeOperationRequiresType
		}
		expectedType, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrTypeOperationTypeNotString
		}
		return op.NewOpTypeOperation(path, expectedType), nil

	case internal.OpTestTypeType:
		if len(compactOp) < 3 {
			return nil, ErrTestTypeOperationRequiresTypes
		}
		types, err := toStringSlice(compactOp[2])
		if err != nil {
			return nil, ErrTestTypeOperationTypesNotArray
		}
		return op.NewOpTestTypeOperationMultiple(path, types), nil

	case internal.OpTestStringType:
		if len(compactOp) < 3 {
			return nil, ErrTestStringOperationRequiresStr
		}
		str, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrTestStringOperationStrNotString
		}
		pos := float64(0)
		if len(compactOp) >= 4 {
			pos, _ = toFloat64(compactOp[3])
		}
		not := false
		if len(compactOp) >= 5 {
			not = toBool(compactOp[4])
		}
		return op.NewOpTestStringOperationFull(path, str, pos, not), nil

	case internal.OpTestStringLenType:
		if len(compactOp) < 3 {
			return nil, ErrTestStringLenOperationRequiresLen
		}
		length, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrTestStringLenOperationLenNotNumber
		}
		not := false
		if len(compactOp) >= 4 {
			not = toBool(compactOp[3])
		}
		return op.NewOpTestStringLenOperationWithNot(path, length, not), nil

	case internal.OpInType:
		if len(compactOp) < 3 {
			return nil, ErrInOperationRequiresValues
		}
		values, ok := compactOp[2].([]interface{})
		if !ok {
			return nil, ErrInOperationValuesNotArray
		}
		return op.NewOpInOperation(path, values), nil

	case internal.OpLessType:
		if len(compactOp) < 3 {
			return nil, ErrLessOperationRequiresValue
		}
		value, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrLessOperationValueNotNumber
		}
		return op.NewOpLessOperation(path, value), nil

	case internal.OpMoreType:
		if len(compactOp) < 3 {
			return nil, ErrMoreOperationRequiresValue
		}
		value, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrMoreOperationValueNotNumber
		}
		return op.NewOpMoreOperation(path, value), nil

	case internal.OpMatchesType:
		if len(compactOp) < 3 {
			return nil, ErrMatchesOperationRequiresPattern
		}
		pattern, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrMatchesOperationPatternNotString
		}
		ignoreCase := false
		if len(compactOp) >= 4 {
			ignoreCase = toBool(compactOp[3])
		}
		return op.NewOpMatchesOperation(path, pattern, ignoreCase, nil), nil

	case internal.OpAndType:
		if len(compactOp) < 3 {
			return nil, ErrAndOperationRequiresOps
		}
		subOps, err := decodePredicateOpsAsInterface(compactOp[2], opts)
		if err != nil {
			return nil, err
		}
		return op.NewOpAndOperation(path, subOps), nil

	case internal.OpOrType:
		if len(compactOp) < 3 {
			return nil, ErrOrOperationRequiresOps
		}
		subOps, err := decodePredicateOpsAsInterface(compactOp[2], opts)
		if err != nil {
			return nil, err
		}
		return op.NewOpOrOperation(path, subOps), nil

	case internal.OpNotType:
		if len(compactOp) < 3 {
			return nil, ErrNotOperationRequiresOps
		}
		subOps, err := decodePredicateOpsAsInterface(compactOp[2], opts)
		if err != nil {
			return nil, err
		}
		return op.NewOpNotOperationMultiple(path, subOps), nil

	case internal.OpStrInsType:
		if len(compactOp) < 4 {
			return nil, ErrStrInsOperationRequiresPosAndStr
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrStrInsOperationPosNotNumber
		}
		str, ok := compactOp[3].(string)
		if !ok {
			return nil, ErrStrInsOperationStrNotString
		}
		return op.NewOpStrInsOperation(path, pos, str), nil

	case internal.OpStrDelType:
		if len(compactOp) < 4 {
			return nil, ErrStrDelOperationRequiresPosAndLen
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrStrDelOperationPosNotNumber
		}
		length, err := toFloat64(compactOp[3])
		if err != nil {
			return nil, ErrStrDelOperationLenNotNumber
		}
		return op.NewOpStrDelOperation(path, pos, length), nil

	case internal.OpSplitType:
		if len(compactOp) < 3 {
			return nil, ErrSplitOperationRequiresPos
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrSplitOperationPosNotNumber
		}
		var props interface{}
		if len(compactOp) >= 4 {
			props = compactOp[3]
		}
		return op.NewOpSplitOperation(path, pos, props), nil

	case internal.OpMergeType:
		if len(compactOp) < 3 {
			return nil, ErrMergeOperationRequiresPos
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrMergeOperationPosNotNumber
		}
		var props map[string]interface{}
		if len(compactOp) >= 4 {
			if p, ok := compactOp[3].(map[string]interface{}); ok {
				props = p
			}
		}
		return op.NewOpMergeOperation(path, pos, props), nil

	case internal.OpExtendType:
		if len(compactOp) < 3 {
			return nil, ErrExtendOperationRequiresProps
		}
		props, ok := compactOp[2].(map[string]interface{})
		if !ok {
			return nil, ErrExtendOperationPropsNotObject
		}
		deleteNull := false
		if len(compactOp) >= 4 {
			deleteNull = toBool(compactOp[3])
		}
		return op.NewOpExtendOperation(path, props, deleteNull), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOperationType, opType)
	}
}

// decodePredicateOpsAsInterface decodes an array of compact operations into []interface{} for And/Or/Not operations
func decodePredicateOpsAsInterface(value interface{}, opts DecoderOptions) ([]interface{}, error) {
	arr, ok := value.([]interface{})
	if !ok {
		return nil, ErrPredicateOpsNotArray
	}

	result := make([]interface{}, 0, len(arr))
	for _, item := range arr {
		compactOp, ok := item.([]interface{})
		if !ok {
			return nil, ErrPredicateOpNotArray
		}
		decoded, err := compactToOp(compactOp, opts)
		if err != nil {
			return nil, err
		}
		_, ok = decoded.(internal.PredicateOp)
		if !ok {
			return nil, ErrDecodedOpNotPredicate
		}
		result = append(result, decoded)
	}
	return result, nil
}

// getOpTypeFromOpcode determines the operation type from the opcode using lookup tables
func getOpTypeFromOpcode(opcode interface{}) (internal.OpType, error) {
	// Try string opcode first (most common for user input)
	if codeStr, ok := opcode.(string); ok {
		if opType, exists := stringToOpType[codeStr]; exists {
			return opType, nil
		}
		return "", fmt.Errorf("%w: %s", ErrUnknownStringOpcode, codeStr)
	}

	// Try numeric opcodes
	var code int
	switch v := opcode.(type) {
	case int:
		code = v
	case float64:
		code = int(v)
	case OpCode:
		code = int(v)
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidOpcodeType, opcode)
	}

	if opType, exists := numericToOpType[code]; exists {
		return opType, nil
	}
	return "", fmt.Errorf("%w: %d", ErrUnknownNumericOpcode, code)
}

// stringToPath converts a JSON pointer string to path slice
func stringToPath(pathStr string) []string {
	if pathStr == "" {
		return []string{}
	}
	path := jsonpointer.Parse(pathStr)
	result := make([]string, len(path))
	for i, token := range path {
		result[i] = fmt.Sprintf("%v", token)
	}
	return result
}

// toBool converts a value to bool
func toBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case float64:
		return val != 0
	case int:
		return val != 0
	default:
		return false
	}
}

// toFloat64 converts a value to float64
func toFloat64(v interface{}) (float64, error) {
	switch val := v.(type) {
	case float64:
		return val, nil
	case int:
		return float64(val), nil
	case int64:
		return float64(val), nil
	default:
		return 0, ErrCannotConvertToFloat64
	}
}

// toStringSlice converts a value to []string
func toStringSlice(v interface{}) ([]string, error) {
	arr, ok := v.([]interface{})
	if !ok {
		return nil, ErrExpectedArray
	}
	result := make([]string, len(arr))
	for i, item := range arr {
		str, ok := item.(string)
		if !ok {
			return nil, ErrExpectedStringInArray
		}
		result[i] = str
	}
	return result, nil
}
