package compact

import (
	"fmt"

	"github.com/kaptinlin/jsonpatch/internal"
)

// Lookup tables for opcode resolution.
var (
	// numericToOpType maps numeric opcodes to operation types.
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

	// stringToOpType maps string opcodes to operation types.
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

// resolveOpType determines the operation type from the opcode.
func resolveOpType(opcode any) (internal.OpType, error) {
	if s, ok := opcode.(string); ok {
		if opType, exists := stringToOpType[s]; exists {
			return opType, nil
		}
		return "", fmt.Errorf("%w: %s", ErrUnknownStringCode, s)
	}

	var code int
	switch v := opcode.(type) {
	case int:
		code = v
	case float64:
		code = int(v)
	case Code:
		code = int(v)
	default:
		return "", fmt.Errorf("%w: %T", ErrInvalidCodeType, opcode)
	}

	if opType, exists := numericToOpType[code]; exists {
		return opType, nil
	}
	return "", fmt.Errorf("%w: %d", ErrUnknownNumericCode, code)
}
