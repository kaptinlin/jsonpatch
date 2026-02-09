package compact

import (
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Decoder decodes compact format operations.
type Decoder struct{}

// NewDecoder creates a new compact decoder.
func NewDecoder() *Decoder {
	return &Decoder{}
}

// Decode decodes a single compact operation.
func (d *Decoder) Decode(raw Op) (internal.Op, error) {
	return parseOp(raw)
}

// DecodeSlice decodes multiple compact operations.
func (d *Decoder) DecodeSlice(ops []Op) ([]internal.Op, error) {
	result := make([]internal.Op, len(ops))
	for i, raw := range ops {
		parsed, err := parseOp(raw)
		if err != nil {
			return nil, err
		}
		result[i] = parsed
	}
	return result, nil
}

// Decode decodes compact format operations.
func Decode(ops []Op) ([]internal.Op, error) {
	return NewDecoder().DecodeSlice(ops)
}

// DecodeJSON decodes compact format JSON bytes into operations.
func DecodeJSON(data []byte) ([]internal.Op, error) {
	var ops []Op
	if err := json.Unmarshal(data, &ops); err != nil {
		return nil, fmt.Errorf("unmarshal compact ops: %w", err)
	}
	return Decode(ops)
}

// parseHeader extracts and validates the opcode and path from a compact operation.
func parseHeader(raw Op) (internal.OpType, []string, error) {
	if len(raw) < 2 {
		return "", nil, ErrMinLength
	}

	pathStr, ok := raw[1].(string)
	if !ok {
		return "", nil, ErrPathNotString
	}

	opType, err := resolveOpType(raw[0])
	if err != nil {
		return "", nil, err
	}

	return opType, parsePath(pathStr), nil
}

// parseOp converts a compact operation to an operation instance.
func parseOp(raw Op) (internal.Op, error) {
	opType, path, err := parseHeader(raw)
	if err != nil {
		return nil, err
	}

	switch opType { //nolint:exhaustive // All cases covered across sub-functions.
	case internal.OpAddType, internal.OpRemoveType, internal.OpReplaceType,
		internal.OpMoveType, internal.OpCopyType, internal.OpTestType:
		return parseCoreOp(opType, path, raw)
	case internal.OpFlipType, internal.OpIncType,
		internal.OpStrInsType, internal.OpStrDelType,
		internal.OpSplitType, internal.OpMergeType, internal.OpExtendType:
		return parseExtendedOp(opType, path, raw)
	case internal.OpAndType, internal.OpOrType, internal.OpNotType:
		return parseCompositeOp(opType, path, raw)
	default:
		return parsePredicateOp(opType, path, raw)
	}
}

// parseCoreOp decodes standard JSON Patch (RFC 6902) operations.
func parseCoreOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles core RFC 6902 operations.
	case internal.OpAddType:
		if len(raw) < 3 {
			return nil, ErrAddMissingValue
		}
		return op.NewAdd(path, raw[2]), nil

	case internal.OpRemoveType:
		if len(raw) >= 3 {
			return op.NewRemoveWithOldValue(path, raw[2]), nil
		}
		return op.NewRemove(path), nil

	case internal.OpReplaceType:
		if len(raw) < 3 {
			return nil, ErrReplaceMissingValue
		}
		if len(raw) >= 4 {
			return op.NewReplaceWithOldValue(path, raw[2], raw[3]), nil
		}
		return op.NewReplace(path, raw[2]), nil

	case internal.OpMoveType:
		from, err := requireFromPath(raw, ErrMoveMissingFrom, ErrMoveFromNotString)
		if err != nil {
			return nil, err
		}
		return op.NewMove(path, from), nil

	case internal.OpCopyType:
		from, err := requireFromPath(raw, ErrCopyMissingFrom, ErrCopyFromNotString)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(path, from), nil

	case internal.OpTestType:
		if len(raw) < 3 {
			return nil, ErrTestMissingValue
		}
		not := boolAt(raw, 3)
		return op.NewTestWithNot(path, raw[2], not), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parseExtendedOp decodes extended operations (flip, inc, str_ins, str_del, split, merge, extend).
func parseExtendedOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles extended operations.
	case internal.OpFlipType:
		return op.NewFlip(path), nil

	case internal.OpIncType:
		if len(raw) < 3 {
			return nil, ErrIncMissingDelta
		}
		delta, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrIncDeltaNotNumber
		}
		return op.NewInc(path, delta), nil

	case internal.OpStrInsType:
		if len(raw) < 4 {
			return nil, ErrStrInsMissingFields
		}
		pos, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrStrInsPosNotNumber
		}
		str, ok := raw[3].(string)
		if !ok {
			return nil, ErrStrInsStrNotString
		}
		return op.NewStrIns(path, pos, str), nil

	case internal.OpStrDelType:
		if len(raw) < 4 {
			return nil, ErrStrDelMissingFields
		}
		pos, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrStrDelPosNotNumber
		}
		length, err := toFloat64(raw[3])
		if err != nil {
			return nil, ErrStrDelLenNotNumber
		}
		return op.NewStrDel(path, pos, length), nil

	case internal.OpSplitType:
		if len(raw) < 3 {
			return nil, ErrSplitMissingPos
		}
		pos, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrSplitPosNotNumber
		}
		var props any
		if len(raw) >= 4 {
			props = raw[3]
		}
		return op.NewSplit(path, pos, props), nil

	case internal.OpMergeType:
		if len(raw) < 3 {
			return nil, ErrMergeMissingPos
		}
		pos, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrMergePosNotNumber
		}
		var props map[string]any
		if len(raw) >= 4 {
			if p, ok := raw[3].(map[string]any); ok {
				props = p
			}
		}
		return op.NewMerge(path, pos, props), nil

	case internal.OpExtendType:
		if len(raw) < 3 {
			return nil, ErrExtendMissingProps
		}
		props, ok := raw[2].(map[string]any)
		if !ok {
			return nil, ErrExtendPropsNotObject
		}
		deleteNull := boolAt(raw, 3)
		return op.NewExtend(path, props, deleteNull), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parsePredicateOp decodes JSON Predicate operations.
func parsePredicateOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	switch opType { //nolint:exhaustive // Only handles predicate operations.
	case internal.OpDefinedType:
		return op.NewDefined(path), nil

	case internal.OpUndefinedType:
		return op.NewUndefined(path), nil

	case internal.OpContainsType:
		value, err := requireString(raw, 2, ErrContainsMissingValue, ErrContainsValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewContainsWithIgnoreCase(path, value, boolAt(raw, 3)), nil

	case internal.OpStartsType:
		value, err := requireString(raw, 2, ErrStartsMissingValue, ErrStartsValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewStartsWithIgnoreCase(path, value, boolAt(raw, 3)), nil

	case internal.OpEndsType:
		value, err := requireString(raw, 2, ErrEndsMissingValue, ErrEndsValueNotString)
		if err != nil {
			return nil, err
		}
		return op.NewEndsWithIgnoreCase(path, value, boolAt(raw, 3)), nil

	case internal.OpTypeType:
		if len(raw) < 3 {
			return nil, ErrTypeMissingType
		}
		expected, ok := raw[2].(string)
		if !ok {
			return nil, ErrTypeNotString
		}
		return op.NewType(path, expected), nil

	case internal.OpTestTypeType:
		if len(raw) < 3 {
			return nil, ErrTestTypeMissingTypes
		}
		types, err := toStringSlice(raw[2])
		if err != nil {
			return nil, ErrTestTypeTypesNotArray
		}
		return op.NewTestTypeMultiple(path, types), nil

	case internal.OpTestStringType:
		if len(raw) < 3 {
			return nil, ErrTestStringMissingStr
		}
		str, ok := raw[2].(string)
		if !ok {
			return nil, ErrTestStringNotString
		}
		pos, _ := float64At(raw, 3)
		not := boolAt(raw, 4)
		return op.NewTestString(path, str, pos, not, false), nil

	case internal.OpTestStringLenType:
		if len(raw) < 3 {
			return nil, ErrTestStringLenMissingLen
		}
		length, err := toFloat64(raw[2])
		if err != nil {
			return nil, ErrTestStringLenNotNumber
		}
		not := boolAt(raw, 3)
		return op.NewTestStringLenWithNot(path, length, not), nil

	case internal.OpInType:
		if len(raw) < 3 {
			return nil, ErrInMissingValues
		}
		values, ok := raw[2].([]any)
		if !ok {
			return nil, ErrInValuesNotArray
		}
		return op.NewIn(path, values), nil

	case internal.OpLessType:
		value, err := requireFloat64(raw, 2, ErrLessMissingValue, ErrLessValueNotNumber)
		if err != nil {
			return nil, err
		}
		return op.NewLess(path, value), nil

	case internal.OpMoreType:
		value, err := requireFloat64(raw, 2, ErrMoreMissingValue, ErrMoreValueNotNumber)
		if err != nil {
			return nil, err
		}
		return op.NewMore(path, value), nil

	case internal.OpMatchesType:
		if len(raw) < 3 {
			return nil, ErrMatchesMissingPattern
		}
		pattern, ok := raw[2].(string)
		if !ok {
			return nil, ErrMatchesPatternNotString
		}
		ignoreCase := boolAt(raw, 3)
		return op.NewMatches(path, pattern, ignoreCase, nil), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parseCompositeOp decodes second-order predicate operations (and, or, not).
func parseCompositeOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	if len(raw) < 3 {
		switch opType { //nolint:exhaustive // Only handles composite operations.
		case internal.OpAndType:
			return nil, ErrAndMissingOps
		case internal.OpOrType:
			return nil, ErrOrMissingOps
		default:
			return nil, ErrNotMissingOps
		}
	}

	subOps, err := parsePredicateOps(raw[2])
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

// requireFromPath extracts and validates the "from" path at index 2.
func requireFromPath(raw Op, errMissing, errNotString error) ([]string, error) {
	if len(raw) < 3 {
		return nil, errMissing
	}
	fromStr, ok := raw[2].(string)
	if !ok {
		return nil, errNotString
	}
	return parsePath(fromStr), nil
}

// requireString extracts a required string value at the given index.
func requireString(raw Op, index int, errMissing, errNotString error) (string, error) {
	if len(raw) <= index {
		return "", errMissing
	}
	v, ok := raw[index].(string)
	if !ok {
		return "", errNotString
	}
	return v, nil
}

// requireFloat64 extracts a required float64 value at the given index.
func requireFloat64(raw Op, index int, errMissing, errNotNumber error) (float64, error) {
	if len(raw) <= index {
		return 0, errMissing
	}
	v, err := toFloat64(raw[index])
	if err != nil {
		return 0, errNotNumber
	}
	return v, nil
}

// parsePredicateOps decodes an array of compact operations into predicate ops.
func parsePredicateOps(value any) ([]any, error) {
	arr, ok := value.([]any)
	if !ok {
		return nil, ErrPredicateNotArray
	}

	result := make([]any, 0, len(arr))
	for _, item := range arr {
		raw, ok := item.([]any)
		if !ok {
			return nil, ErrPredicateOpInvalid
		}
		decoded, err := parseOp(raw)
		if err != nil {
			return nil, err
		}
		if _, ok := decoded.(internal.PredicateOp); !ok {
			return nil, ErrNotPredicate
		}
		result = append(result, decoded)
	}
	return result, nil
}

// parsePath converts a JSON pointer string to a path slice.
func parsePath(s string) []string {
	if s == "" {
		return []string{}
	}
	return []string(jsonpointer.Parse(s))
}
