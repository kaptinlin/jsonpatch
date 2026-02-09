package compact

import (
	"fmt"

	"github.com/go-json-experiment/json"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Pre-built mapping tables for better performance and maintainability.
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

// Decoder decodes compact format operations into operation instances.
type Decoder struct{}

// NewDecoder creates a new compact decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode decodes a single compact operation into an operation instance.
func (d *Decoder) Decode(compactOp Op) (internal.Op, error) {
	return parseOp(compactOp)
}

// DecodeSlice decodes multiple compact operations.
func (d *Decoder) DecodeSlice(compactOps []Op) ([]internal.Op, error) {
	result := make([]internal.Op, len(compactOps))
	for i, compactOp := range compactOps {
		parsed, err := parseOp(compactOp)
		if err != nil {
			return nil, err
		}
		result[i] = parsed
	}
	return result, nil
}

// Decode decodes compact format operations.
func Decode(compactOps []Op) ([]internal.Op, error) {
	return NewDecoder().DecodeSlice(compactOps)
}

// DecodeJSON decodes compact format JSON bytes into operations.
func DecodeJSON(data []byte) ([]internal.Op, error) {
	var compactOps []Op
	if err := json.Unmarshal(data, &compactOps); err != nil {
		return nil, fmt.Errorf("unmarshal compact ops: %w", err)
	}
	return Decode(compactOps)
}

// --- Header parsing ---

// parseHeader extracts and validates the opcode and path from a compact operation.
func parseHeader(compactOp Op) (internal.OpType, []string, error) {
	if len(compactOp) < 2 {
		return "", nil, ErrOpMinLength
	}

	pathStr, ok := compactOp[1].(string)
	if !ok {
		return "", nil, ErrOpPathNotString
	}

	opType, err := resolveOpType(compactOp[0])
	if err != nil {
		return "", nil, err
	}

	return opType, parsePath(pathStr), nil
}

// --- Operation dispatching ---

// parseOp converts a compact operation to an operation instance.
func parseOp(compactOp Op) (internal.Op, error) {
	opType, path, err := parseHeader(compactOp)
	if err != nil {
		return nil, err
	}

	switch opType { //nolint:exhaustive // Intentional grouping; all cases covered across sub-functions.
	case internal.OpAddType, internal.OpRemoveType, internal.OpReplaceType,
		internal.OpMoveType, internal.OpCopyType, internal.OpTestType:
		return parseCoreOp(opType, path, compactOp)
	case internal.OpFlipType, internal.OpIncType,
		internal.OpStrInsType, internal.OpStrDelType,
		internal.OpSplitType, internal.OpMergeType, internal.OpExtendType:
		return parseExtendedOp(opType, path, compactOp)
	case internal.OpAndType, internal.OpOrType, internal.OpNotType:
		return parseCompositeOp(opType, path, compactOp)
	default:
		return parsePredicateOp(opType, path, compactOp)
	}
}

// --- Core operations (RFC 6902) ---

// parseCoreOp decodes standard JSON Patch (RFC 6902) operations.
func parseCoreOp(opType internal.OpType, path []string, compactOp Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles core RFC 6902 operations.
	case internal.OpAddType:
		if len(compactOp) < 3 {
			return nil, ErrAddOpMissingValue
		}
		return op.NewAdd(path, compactOp[2]), nil

	case internal.OpRemoveType:
		if len(compactOp) >= 3 {
			return op.NewRemoveWithOldValue(path, compactOp[2]), nil
		}
		return op.NewRemove(path), nil

	case internal.OpReplaceType:
		if len(compactOp) < 3 {
			return nil, ErrReplaceOpMissingValue
		}
		if len(compactOp) >= 4 {
			return op.NewReplaceWithOldValue(path, compactOp[2], compactOp[3]), nil
		}
		return op.NewReplace(path, compactOp[2]), nil

	case internal.OpMoveType:
		from, err := requireFromPath(compactOp, ErrMoveOpMissingFrom, ErrMoveOpFromNotString)
		if err != nil {
			return nil, err
		}
		return op.NewMove(path, from), nil

	case internal.OpCopyType:
		from, err := requireFromPath(compactOp, ErrCopyOpMissingFrom, ErrCopyOpFromNotString)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(path, from), nil

	case internal.OpTestType:
		if len(compactOp) < 3 {
			return nil, ErrTestOpMissingValue
		}
		not := boolAt(compactOp, 3)
		return op.NewTestWithNot(path, compactOp[2], not), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// --- Extended operations ---

// parseExtendedOp decodes extended operations (flip, inc, str_ins, str_del, split, merge, extend).
func parseExtendedOp(opType internal.OpType, path []string, compactOp Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles extended operations.
	case internal.OpFlipType:
		return op.NewFlip(path), nil

	case internal.OpIncType:
		if len(compactOp) < 3 {
			return nil, ErrIncOpMissingDelta
		}
		delta, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrIncOpDeltaNotNumber
		}
		return op.NewInc(path, delta), nil

	case internal.OpStrInsType:
		if len(compactOp) < 4 {
			return nil, ErrStrInsOpMissingFields
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrStrInsOpPosNotNumber
		}
		str, ok := compactOp[3].(string)
		if !ok {
			return nil, ErrStrInsOpStrNotString
		}
		return op.NewStrIns(path, pos, str), nil

	case internal.OpStrDelType:
		if len(compactOp) < 4 {
			return nil, ErrStrDelOpMissingFields
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrStrDelOpPosNotNumber
		}
		length, err := toFloat64(compactOp[3])
		if err != nil {
			return nil, ErrStrDelOpLenNotNumber
		}
		return op.NewStrDel(path, pos, length), nil

	case internal.OpSplitType:
		if len(compactOp) < 3 {
			return nil, ErrSplitOpMissingPos
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrSplitOpPosNotNumber
		}
		var props any
		if len(compactOp) >= 4 {
			props = compactOp[3]
		}
		return op.NewSplit(path, pos, props), nil

	case internal.OpMergeType:
		if len(compactOp) < 3 {
			return nil, ErrMergeOpMissingPos
		}
		pos, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrMergeOpPosNotNumber
		}
		var props map[string]any
		if len(compactOp) >= 4 {
			if p, ok := compactOp[3].(map[string]any); ok {
				props = p
			}
		}
		return op.NewMerge(path, pos, props), nil

	case internal.OpExtendType:
		if len(compactOp) < 3 {
			return nil, ErrExtendOpMissingProps
		}
		props, ok := compactOp[2].(map[string]any)
		if !ok {
			return nil, ErrExtendOpPropsNotObject
		}
		deleteNull := boolAt(compactOp, 3)
		return op.NewExtend(path, props, deleteNull), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// --- Predicate operations ---

// parsePredicateOp decodes JSON Predicate operations.
func parsePredicateOp(opType internal.OpType, path []string, compactOp Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles predicate operations.
	case internal.OpDefinedType:
		return op.NewDefined(path), nil

	case internal.OpUndefinedType:
		return op.NewUndefined(path), nil

	case internal.OpContainsType:
		value, err := requireString(compactOp, 2, ErrContainsOpMissingValue, ErrContainsOpValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewContainsWithIgnoreCase(path, value, boolAt(compactOp, 3)), nil

	case internal.OpStartsType:
		value, err := requireString(compactOp, 2, ErrStartsOpMissingValue, ErrStartsOpValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewStartsWithIgnoreCase(path, value, boolAt(compactOp, 3)), nil

	case internal.OpEndsType:
		value, err := requireString(compactOp, 2, ErrEndsOpMissingValue, ErrEndsOpValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewEndsWithIgnoreCase(path, value, boolAt(compactOp, 3)), nil

	case internal.OpTypeType:
		if len(compactOp) < 3 {
			return nil, ErrTypeOpMissingType
		}
		expectedType, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrTypeOpTypeNotString
		}
		return op.NewType(path, expectedType), nil

	case internal.OpTestTypeType:
		if len(compactOp) < 3 {
			return nil, ErrTestTypeOpMissingTypes
		}
		types, err := toStringSlice(compactOp[2])
		if err != nil {
			return nil, ErrTestTypeOpTypesNotArray
		}
		return op.NewTestTypeMultiple(path, types), nil

	case internal.OpTestStringType:
		if len(compactOp) < 3 {
			return nil, ErrTestStringOpMissingStr
		}
		str, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrTestStringOpStrNotString
		}
		pos, _ := float64At(compactOp, 3)
		not := boolAt(compactOp, 4)
		return op.NewTestStringFull(path, str, pos, not), nil

	case internal.OpTestStringLenType:
		if len(compactOp) < 3 {
			return nil, ErrTestStringLenOpMissingLen
		}
		length, err := toFloat64(compactOp[2])
		if err != nil {
			return nil, ErrTestStringLenOpLenNotNumber
		}
		not := boolAt(compactOp, 3)
		return op.NewTestStringLenWithNot(path, length, not), nil

	case internal.OpInType:
		if len(compactOp) < 3 {
			return nil, ErrInOpMissingValues
		}
		values, ok := compactOp[2].([]any)
		if !ok {
			return nil, ErrInOpValuesNotArray
		}
		return op.NewIn(path, values), nil

	case internal.OpLessType:
		value, err := requireFloat64(compactOp, 2, ErrLessOpMissingValue, ErrLessOpValueNotNumber)
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, value), nil

	case internal.OpMoreType:
		value, err := requireFloat64(compactOp, 2, ErrMoreOpMissingValue, ErrMoreOpValueNotNumber)
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, value), nil

	case internal.OpMatchesType:
		if len(compactOp) < 3 {
			return nil, ErrMatchesOpMissingPattern
		}
		pattern, ok := compactOp[2].(string)
		if !ok {
			return nil, ErrMatchesOpPatternNotString
		}
		ignoreCase := boolAt(compactOp, 3)
		return op.NewMatches(path, pattern, ignoreCase, nil), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// --- Composite operations (and, or, not) ---

// parseCompositeOp decodes second-order predicate operations (and, or, not).
func parseCompositeOp(opType internal.OpType, path []string, compactOp Op) (internal.Op, error) {
	if len(compactOp) < 3 {
		switch opType { //nolint:exhaustive // Only handles composite operations.
		case internal.OpAndType:
			return nil, ErrAndOpMissingOps
		case internal.OpOrType:
			return nil, ErrOrOpMissingOps
		default:
			return nil, ErrNotOpMissingOps
		}
	}

	subOps, err := parsePredicateOps(compactOp[2])
	if err != nil {
		return nil, err
	}

	switch opType { //nolint:exhaustive // Only handles composite operations.
	case internal.OpAndType:
		return op.NewAnd(path, subOps), nil
	case internal.OpOrType:
		return op.NewOr(path, subOps), nil
	default:
		return op.NewNotMultiple(path, subOps), nil
	}
}

// --- Helper functions ---

// requireFromPath extracts and validates the "from" path at index 2.
func requireFromPath(compactOp Op, errMissing, errNotString error) ([]string, error) {
	if len(compactOp) < 3 {
		return nil, errMissing
	}
	fromStr, ok := compactOp[2].(string)
	if !ok {
		return nil, errNotString
	}
	return parsePath(fromStr), nil
}

// requireString extracts a required string value at the given index.
func requireString(compactOp Op, index int, errMissing, errNotString error) (string, error) {
	if len(compactOp) <= index {
		return "", errMissing
	}
	value, ok := compactOp[index].(string)
	if !ok {
		return "", errNotString
	}
	return value, nil
}

// requireFloat64 extracts a required float64 value at the given index.
func requireFloat64(compactOp Op, index int, errMissing, errNotNumber error) (float64, error) {
	if len(compactOp) <= index {
		return 0, errMissing
	}
	value, err := toFloat64(compactOp[index])
	if err != nil {
		return 0, errNotNumber
	}
	return value, nil
}

// parsePredicateOps decodes an array of compact operations into []any for and/or/not operations.
func parsePredicateOps(value any) ([]any, error) {
	arr, ok := value.([]any)
	if !ok {
		return nil, ErrPredicateOpsNotArray
	}

	result := make([]any, 0, len(arr))
	for _, item := range arr {
		compactOp, ok := item.([]any)
		if !ok {
			return nil, ErrPredicateOpNotArray
		}
		decoded, err := parseOp(compactOp)
		if err != nil {
			return nil, err
		}
		if _, ok := decoded.(internal.PredicateOp); !ok {
			return nil, ErrDecodedOpNotPredicate
		}
		result = append(result, decoded)
	}
	return result, nil
}

// --- Opcode resolution ---

// resolveOpType determines the operation type from the opcode using lookup tables.
func resolveOpType(opcode any) (internal.OpType, error) {
	if codeStr, ok := opcode.(string); ok {
		if opType, exists := stringToOpType[codeStr]; exists {
			return opType, nil
		}
		return "", fmt.Errorf("%w: %s", ErrUnknownStringOpcode, codeStr)
	}

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

// --- Path utilities ---

// parsePath converts a JSON pointer string to a path slice.
func parsePath(pathStr string) []string {
	if pathStr == "" {
		return []string{}
	}
	return []string(jsonpointer.Parse(pathStr))
}

// --- Type conversion utilities ---

// boolAt safely extracts a bool value at the given index, returning false if absent.
func boolAt(compactOp Op, index int) bool {
	if len(compactOp) <= index {
		return false
	}
	return toBool(compactOp[index])
}

// float64At safely extracts a float64 value at the given index.
func float64At(compactOp Op, index int) (float64, error) {
	if len(compactOp) <= index {
		return 0, nil
	}
	return toFloat64(compactOp[index])
}

// toBool converts a value to bool.
func toBool(v any) bool {
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

// toFloat64 converts a value to float64.
func toFloat64(v any) (float64, error) {
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

// toStringSlice converts a value to []string.
func toStringSlice(v any) ([]string, error) {
	arr, ok := v.([]any)
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
