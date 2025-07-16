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
func Decode(compactOps []CompactOp, opts ...DecoderOption) ([]internal.Op, error) {
	decoder := NewDecoder(opts...)
	return decoder.DecodeSlice(compactOps)
}

// DecodeJSON decodes compact format JSON bytes into operations
func DecodeJSON(data []byte, opts ...DecoderOption) ([]internal.Op, error) {
	var compactOps []CompactOp
	if err := json.Unmarshal(data, &compactOps); err != nil {
		return nil, fmt.Errorf("failed to unmarshal compact operations: %w", err)
	}
	return Decode(compactOps, opts...)
}

// compactToOp converts a compact operation to an operation instance
func compactToOp(compactOp CompactOp, options DecoderOptions) (internal.Op, error) {
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
		return op.NewOpAddOperation(path, compactOp[2]), nil

	case internal.OpRemoveType:
		// Optional oldValue parameter
		if len(compactOp) >= 3 {
			return op.NewOpRemoveOperationWithOldValue(path, compactOp[2]), nil
		}
		return op.NewOpRemoveOperation(path), nil

	case internal.OpReplaceType:
		if len(compactOp) < 3 {
			return nil, ErrReplaceOperationRequiresValue
		}
		// Optional oldValue parameter
		if len(compactOp) >= 4 {
			return op.NewOpReplaceOperationWithOldValue(path, compactOp[2], compactOp[3]), nil
		}
		return op.NewOpReplaceOperation(path, compactOp[2]), nil

	case internal.OpMoveType:
		if len(compactOp) < 3 {
			return nil, ErrMoveOperationRequiresFrom
		}
		fromStr, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrMoveOperationFromNotString
		}
		from := stringToPath(fromStr)
		return op.NewOpMoveOperation(path, from), nil

	case internal.OpCopyType:
		if len(compactOp) < 3 {
			return nil, ErrCopyOperationRequiresFrom
		}
		fromStr, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrCopyOperationFromNotString
		}
		from := stringToPath(fromStr)
		return op.NewOpCopyOperation(path, from), nil

	case internal.OpTestType:
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		// Currently test operation doesn't have a "not" variant in constructors
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpFlipType:
		return op.NewOpFlipOperation(path), nil

	case internal.OpIncType:
		if len(compactOp) < 3 {
			return nil, ErrIncOperationRequiresDelta
		}
		delta, ok := compactOp[2].(float64)
		if !ok {
			return nil, ErrIncOperationDeltaNotNumber
		}
		return op.NewOpIncOperation(path, delta), nil

	case internal.OpDefinedType:
		return op.NewOpDefinedOperation(path), nil

	case internal.OpUndefinedType:
		return op.NewOpUndefinedOperation(path, false), nil

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
			if ignoreCaseVal, ok := compactOp[3].(float64); ok && ignoreCaseVal == 1 {
				ignoreCase = true
			}
		}
		if ignoreCase {
			return op.NewOpContainsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpContainsOperation(path, value), nil

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
			if ignoreCaseVal, ok := compactOp[3].(float64); ok && ignoreCaseVal == 1 {
				ignoreCase = true
			}
		}
		if ignoreCase {
			return op.NewOpStartsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpStartsOperation(path, value), nil

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
			if ignoreCaseVal, ok := compactOp[3].(float64); ok && ignoreCaseVal == 1 {
				ignoreCase = true
			}
		}
		if ignoreCase {
			return op.NewOpEndsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpEndsOperation(path, value), nil

	case internal.OpTypeType:
		// Type test operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpTestTypeType:
		// Test type operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpTestStringType:
		// Test string operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpTestStringLenType:
		// Test string length operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpInType:
		// In operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpLessType:
		// Less operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpMoreType:
		// More operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpMatchesType:
		// Matches operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpAndType:
		// And operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpOrType:
		// Or operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpNotType:
		// Not operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpStrInsType:
		// String insert operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpStrDelType:
		// String delete operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpSplitType:
		// Split operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpMergeType:
		// Merge operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	case internal.OpExtendType:
		// Extend operation - currently fallback to test operation
		if len(compactOp) < 3 {
			return nil, ErrTestOperationRequiresValue
		}
		return op.NewOpTestOperation(path, compactOp[2]), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOperationType, opType)
	}
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
