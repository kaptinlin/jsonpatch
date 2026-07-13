package compact

import (
	"errors"
	"fmt"
	"slices"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
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
			return nil, fmt.Errorf("compact operation %d: %w", i, err)
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

	path, err := parsePathValue(raw[1], ErrPathNotString)
	if err != nil {
		return "", nil, err
	}

	opType, err := resolveOpType(raw[0])
	if err != nil {
		return "", nil, err
	}

	return opType, path, nil
}

func isNumericConstraintError(err error) bool {
	return errors.Is(err, ErrNumberNotFinite) || errors.Is(err, ErrNumberNotInteger)
}

func predicateMaxArity(opType internal.OpType) int {
	switch opType {
	case internal.OpDefinedType, internal.OpUndefinedType:
		return 2
	case internal.OpContainsType, internal.OpStartsType, internal.OpEndsType,
		internal.OpMatchesType, internal.OpTestStringLenType:
		return 4
	case internal.OpTestStringType:
		return 5
	default:
		return 3
	}
}

// parseOp converts a compact operation to an operation instance.
func parseOp(raw Op) (internal.Op, error) {
	opType, path, err := parseHeader(raw)
	if err != nil {
		return nil, err
	}

	switch opType {
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
	switch opType {
	case internal.OpAddType:
		if len(raw) < 3 {
			return nil, ErrAddMissingValue
		}
		if len(raw) != 3 {
			return nil, ErrInvalidArity
		}
		return op.NewAdd(path, raw[2]), nil

	case internal.OpRemoveType:
		if len(raw) > 3 {
			return nil, ErrInvalidArity
		}
		if len(raw) >= 3 {
			return op.NewRemoveWithOldValue(path, raw[2]), nil
		}
		return op.NewRemove(path), nil

	case internal.OpReplaceType:
		if len(raw) < 3 {
			return nil, ErrReplaceMissingValue
		}
		if len(raw) > 4 {
			return nil, ErrInvalidArity
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
		if len(raw) != 3 {
			return nil, ErrInvalidArity
		}
		return op.NewMove(path, from), nil

	case internal.OpCopyType:
		from, err := requireFromPath(raw, ErrCopyMissingFrom, ErrCopyFromNotString)
		if err != nil {
			return nil, err
		}
		if len(raw) != 3 {
			return nil, ErrInvalidArity
		}
		return op.NewCopy(path, from), nil

	case internal.OpTestType:
		if len(raw) < 3 {
			return nil, ErrTestMissingValue
		}
		if len(raw) > 4 {
			return nil, ErrInvalidArity
		}
		not, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewTestWithNot(path, raw[2], not), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parseExtendedOp decodes extended operations (flip, inc, str_ins, str_del, split, merge, extend).
func parseExtendedOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	switch opType {
	case internal.OpFlipType:
		if len(raw) != 2 {
			return nil, ErrInvalidArity
		}
		return op.NewFlip(path), nil

	case internal.OpIncType:
		if len(raw) < 3 {
			return nil, ErrIncMissingDelta
		}
		if len(raw) != 3 {
			return nil, ErrInvalidArity
		}
		delta, err := finiteNumber(raw[2])
		if err != nil {
			if errors.Is(err, ErrNumberNotFinite) {
				return nil, err
			}
			return nil, ErrIncDeltaNotNumber
		}
		return op.NewInc(path, delta), nil

	case internal.OpStrInsType:
		if len(raw) < 4 {
			return nil, ErrStrInsMissingFields
		}
		if len(raw) != 4 {
			return nil, ErrInvalidArity
		}
		pos, err := integerNumber(raw[2])
		if err != nil {
			if isNumericConstraintError(err) {
				return nil, err
			}
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
		if len(raw) > 5 {
			return nil, ErrInvalidArity
		}
		pos, err := integerNumber(raw[2])
		if err != nil {
			if isNumericConstraintError(err) {
				return nil, err
			}
			return nil, ErrStrDelPosNotNumber
		}
		if str, ok := raw[3].(string); ok {
			if len(raw) != 4 {
				return nil, ErrInvalidArity
			}
			return op.NewStrDelWithStr(path, pos, str), nil
		}
		if len(raw) != 5 {
			return nil, ErrInvalidVariant
		}
		if !isZeroMarker(raw[3]) {
			return nil, ErrInvalidVariant
		}
		length, err := integerNumber(raw[4])
		if err != nil {
			if isNumericConstraintError(err) {
				return nil, err
			}
			return nil, ErrStrDelLenNotNumber
		}
		return op.NewStrDel(path, pos, length), nil

	case internal.OpSplitType:
		if len(raw) < 3 {
			return nil, ErrSplitMissingPos
		}
		if len(raw) > 4 {
			return nil, ErrInvalidArity
		}
		pos, err := finiteNumber(raw[2])
		if err != nil {
			if errors.Is(err, ErrNumberNotFinite) {
				return nil, err
			}
			return nil, ErrSplitPosNotNumber
		}
		var props any
		if len(raw) >= 4 {
			var ok bool
			props, ok = raw[3].(map[string]any)
			if !ok {
				return nil, ErrSplitPropsNotObject
			}
		}
		return op.NewSplit(path, pos, props), nil

	case internal.OpMergeType:
		if len(raw) < 3 {
			return nil, ErrMergeMissingPos
		}
		if len(raw) > 4 {
			return nil, ErrInvalidArity
		}
		pos, err := finiteNumber(raw[2])
		if err != nil {
			if errors.Is(err, ErrNumberNotFinite) {
				return nil, err
			}
			return nil, ErrMergePosNotNumber
		}
		var props map[string]any
		if len(raw) >= 4 {
			var ok bool
			props, ok = raw[3].(map[string]any)
			if !ok {
				return nil, ErrMergePropsNotObject
			}
		}
		return op.NewMerge(path, pos, props), nil

	case internal.OpExtendType:
		if len(raw) < 3 {
			return nil, ErrExtendMissingProps
		}
		if len(raw) > 4 {
			return nil, ErrInvalidArity
		}
		props, ok := raw[2].(map[string]any)
		if !ok {
			return nil, ErrExtendPropsNotObject
		}
		deleteNull, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewExtend(path, props, deleteNull), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parsePredicateOp decodes JSON Predicate operations.
func parsePredicateOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	maxArity := predicateMaxArity(opType)
	if len(raw) > maxArity {
		return nil, ErrInvalidArity
	}

	switch opType {
	case internal.OpDefinedType:
		return op.NewDefined(path), nil

	case internal.OpUndefinedType:
		return op.NewUndefined(path), nil

	case internal.OpContainsType:
		value, err := requireString(raw, 2, ErrContainsMissingValue, ErrContainsValueNotString)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewContainsWithIgnoreCase(path, value, ignoreCase), nil

	case internal.OpStartsType:
		value, err := requireString(raw, 2, ErrStartsMissingValue, ErrStartsValueNotString)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewStartsWithIgnoreCase(path, value, ignoreCase), nil

	case internal.OpEndsType:
		value, err := requireString(raw, 2, ErrEndsMissingValue, ErrEndsValueNotString)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewEndsWithIgnoreCase(path, value, ignoreCase), nil

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
			return nil, ErrTestStringMissingPos
		}
		pos, err := integerNumber(raw[2])
		if err != nil {
			if isNumericConstraintError(err) {
				return nil, err
			}
			return nil, ErrTestStringPosNotNumber
		}
		if len(raw) < 4 {
			return nil, ErrTestStringMissingStr
		}
		str, ok := raw[3].(string)
		if !ok {
			return nil, ErrTestStringNotString
		}
		not, err := trueMarkerAt(raw, 4)
		if err != nil {
			return nil, err
		}
		return op.NewTestString(path, str, pos, not, false), nil

	case internal.OpTestStringLenType:
		if len(raw) < 3 {
			return nil, ErrTestStringLenMissingLen
		}
		length, err := integerNumber(raw[2])
		if err != nil {
			if isNumericConstraintError(err) {
				return nil, err
			}
			return nil, ErrTestStringLenNotNumber
		}
		not, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
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
		ignoreCase, err := trueMarkerAt(raw, 3)
		if err != nil {
			return nil, err
		}
		return op.NewMatches(path, pattern, ignoreCase, nil), nil

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedOp, opType)
	}
}

// parseCompositeOp decodes second-order predicate operations (and, or, not).
func parseCompositeOp(opType internal.OpType, path []string, raw Op) (internal.Op, error) {
	if len(raw) < 3 {
		switch opType {
		case internal.OpAndType:
			return nil, ErrAndMissingOps
		case internal.OpOrType:
			return nil, ErrOrMissingOps
		default:
			return nil, ErrNotMissingOps
		}
	}
	if len(raw) != 3 {
		return nil, ErrInvalidArity
	}

	subOps, err := parsePredicateOps(raw[2], path)
	if err != nil {
		return nil, err
	}

	switch opType {
	case internal.OpAndType:
		return op.NewAnd(path, subOps), nil
	case internal.OpOrType:
		return op.NewOr(path, subOps), nil
	default:
		if len(subOps) != 1 {
			return nil, ErrNotSinglePredicate
		}
		return op.NewNotMultiple(path, subOps), nil
	}
}

// requireFromPath extracts and validates the "from" path at index 2.
func requireFromPath(raw Op, errMissing, errNotString error) ([]string, error) {
	if len(raw) < 3 {
		return nil, errMissing
	}
	return parsePathValue(raw[2], errNotString)
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
func parsePredicateOps(value any, base []string) ([]any, error) {
	arr, ok := value.([]any)
	if !ok {
		return nil, ErrPredicateNotArray
	}

	result := make([]any, len(arr))
	for i, item := range arr {
		raw, ok := item.([]any)
		if !ok {
			return nil, ErrPredicateOpInvalid
		}
		merged, err := withMergedPredicatePath(Op(raw), base)
		if err != nil {
			return nil, err
		}
		decoded, err := parseOp(merged)
		if err != nil {
			return nil, err
		}
		if _, ok := decoded.(internal.PredicateOp); !ok {
			return nil, ErrNotPredicate
		}
		result[i] = decoded
	}
	return result, nil
}

func withMergedPredicatePath(raw Op, base []string) (Op, error) {
	if len(raw) < 2 {
		return nil, ErrMinLength
	}

	child, err := parsePathValue(raw[1], ErrPathNotString)
	if err != nil {
		return nil, err
	}

	clone := slices.Clone(raw)
	clone[1] = mergePaths(base, child)
	return clone, nil
}

func parsePathValue(value any, errInvalid error) ([]string, error) {
	switch path := value.(type) {
	case []string:
		return slices.Clone(path), nil
	case []any:
		segments := make([]string, len(path))
		for i, segment := range path {
			value, ok := segment.(string)
			if !ok {
				return nil, errInvalid
			}
			segments[i] = value
		}
		return segments, nil
	default:
		return nil, errInvalid
	}
}

// mergePaths combines parent and child predicate paths.
func mergePaths(base, child []string) []string {
	return slices.Concat(base, child)
}
