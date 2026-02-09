package compact

import (
	"errors"
)

// --- Base operation errors ---

var (
	// ErrOpMinLength indicates compact operation must have at least opcode and path.
	ErrOpMinLength = errors.New("compact operation must have at least opcode and path")
	// ErrOpPathNotString indicates compact operation path must be a string.
	ErrOpPathNotString = errors.New("compact operation path must be a string")
)

// --- Core operation errors ---

var (
	// ErrAddOpMissingValue indicates add operation is missing value.
	ErrAddOpMissingValue = errors.New("add operation requires value")
	// ErrReplaceOpMissingValue indicates replace operation is missing value.
	ErrReplaceOpMissingValue = errors.New("replace operation requires value")
	// ErrMoveOpMissingFrom indicates move operation is missing from path.
	ErrMoveOpMissingFrom = errors.New("move operation requires from path")
	// ErrMoveOpFromNotString indicates move operation from must be a string.
	ErrMoveOpFromNotString = errors.New("move operation from must be a string")
	// ErrCopyOpMissingFrom indicates copy operation is missing from path.
	ErrCopyOpMissingFrom = errors.New("copy operation requires from path")
	// ErrCopyOpFromNotString indicates copy operation from must be a string.
	ErrCopyOpFromNotString = errors.New("copy operation from must be a string")
	// ErrTestOpMissingValue indicates test operation is missing value.
	ErrTestOpMissingValue = errors.New("test operation requires value")
)

// --- Extended operation errors ---

var (
	// ErrIncOpMissingDelta indicates inc operation is missing delta.
	ErrIncOpMissingDelta = errors.New("inc operation requires delta")
	// ErrIncOpDeltaNotNumber indicates inc operation delta must be a number.
	ErrIncOpDeltaNotNumber = errors.New("inc operation delta must be a number")
	// ErrStrInsOpMissingFields indicates str_ins operation is missing pos and str.
	ErrStrInsOpMissingFields = errors.New("str_ins operation requires pos and str")
	// ErrStrInsOpPosNotNumber indicates str_ins operation pos must be a number.
	ErrStrInsOpPosNotNumber = errors.New("str_ins operation pos must be a number")
	// ErrStrInsOpStrNotString indicates str_ins operation str must be a string.
	ErrStrInsOpStrNotString = errors.New("str_ins operation str must be a string")
	// ErrStrDelOpMissingFields indicates str_del operation is missing pos and len.
	ErrStrDelOpMissingFields = errors.New("str_del operation requires pos and len")
	// ErrStrDelOpPosNotNumber indicates str_del operation pos must be a number.
	ErrStrDelOpPosNotNumber = errors.New("str_del operation pos must be a number")
	// ErrStrDelOpLenNotNumber indicates str_del operation len must be a number.
	ErrStrDelOpLenNotNumber = errors.New("str_del operation len must be a number")
	// ErrSplitOpMissingPos indicates split operation is missing pos.
	ErrSplitOpMissingPos = errors.New("split operation requires pos")
	// ErrSplitOpPosNotNumber indicates split operation pos must be a number.
	ErrSplitOpPosNotNumber = errors.New("split operation pos must be a number")
	// ErrMergeOpMissingPos indicates merge operation is missing pos.
	ErrMergeOpMissingPos = errors.New("merge operation requires pos")
	// ErrMergeOpPosNotNumber indicates merge operation pos must be a number.
	ErrMergeOpPosNotNumber = errors.New("merge operation pos must be a number")
	// ErrExtendOpMissingProps indicates extend operation is missing props.
	ErrExtendOpMissingProps = errors.New("extend operation requires props")
	// ErrExtendOpPropsNotObject indicates extend operation props must be an object.
	ErrExtendOpPropsNotObject = errors.New("extend operation props must be an object")
)

// --- Predicate operation errors ---

var (
	// ErrContainsOpMissingValue indicates contains operation is missing value.
	ErrContainsOpMissingValue = errors.New("contains operation requires value")
	// ErrContainsOpValueNotString indicates contains operation value must be a string.
	ErrContainsOpValueNotString = errors.New("contains operation value must be a string")
	// ErrStartsOpMissingValue indicates starts operation is missing value.
	ErrStartsOpMissingValue = errors.New("starts operation requires value")
	// ErrStartsOpValueNotString indicates starts operation value must be a string.
	ErrStartsOpValueNotString = errors.New("starts operation value must be a string")
	// ErrEndsOpMissingValue indicates ends operation is missing value.
	ErrEndsOpMissingValue = errors.New("ends operation requires value")
	// ErrEndsOpValueNotString indicates ends operation value must be a string.
	ErrEndsOpValueNotString = errors.New("ends operation value must be a string")
	// ErrTypeOpMissingType indicates type operation is missing type.
	ErrTypeOpMissingType = errors.New("type operation requires type")
	// ErrTypeOpTypeNotString indicates type operation type must be a string.
	ErrTypeOpTypeNotString = errors.New("type operation type must be a string")
	// ErrTestTypeOpMissingTypes indicates test_type operation is missing types.
	ErrTestTypeOpMissingTypes = errors.New("test_type operation requires types")
	// ErrTestTypeOpTypesNotArray indicates test_type operation types must be an array.
	ErrTestTypeOpTypesNotArray = errors.New("test_type operation types must be an array")
	// ErrTestStringOpMissingStr indicates test_string operation is missing str.
	ErrTestStringOpMissingStr = errors.New("test_string operation requires str")
	// ErrTestStringOpStrNotString indicates test_string operation str must be a string.
	ErrTestStringOpStrNotString = errors.New("test_string operation str must be a string")
	// ErrTestStringLenOpMissingLen indicates test_string_len operation is missing len.
	ErrTestStringLenOpMissingLen = errors.New("test_string_len operation requires len")
	// ErrTestStringLenOpLenNotNumber indicates test_string_len operation len must be a number.
	ErrTestStringLenOpLenNotNumber = errors.New("test_string_len operation len must be a number")
	// ErrInOpMissingValues indicates in operation is missing values.
	ErrInOpMissingValues = errors.New("in operation requires values")
	// ErrInOpValuesNotArray indicates in operation values must be an array.
	ErrInOpValuesNotArray = errors.New("in operation values must be an array")
	// ErrLessOpMissingValue indicates less operation is missing value.
	ErrLessOpMissingValue = errors.New("less operation requires value")
	// ErrLessOpValueNotNumber indicates less operation value must be a number.
	ErrLessOpValueNotNumber = errors.New("less operation value must be a number")
	// ErrMoreOpMissingValue indicates more operation is missing value.
	ErrMoreOpMissingValue = errors.New("more operation requires value")
	// ErrMoreOpValueNotNumber indicates more operation value must be a number.
	ErrMoreOpValueNotNumber = errors.New("more operation value must be a number")
	// ErrMatchesOpMissingPattern indicates matches operation is missing pattern.
	ErrMatchesOpMissingPattern = errors.New("matches operation requires pattern")
	// ErrMatchesOpPatternNotString indicates matches operation pattern must be a string.
	ErrMatchesOpPatternNotString = errors.New("matches operation pattern must be a string")
)

// --- Composite operation errors ---

var (
	// ErrAndOpMissingOps indicates and operation is missing ops.
	ErrAndOpMissingOps = errors.New("and operation requires ops")
	// ErrOrOpMissingOps indicates or operation is missing ops.
	ErrOrOpMissingOps = errors.New("or operation requires ops")
	// ErrNotOpMissingOps indicates not operation is missing ops.
	ErrNotOpMissingOps = errors.New("not operation requires ops")
	// ErrPredicateOpsNotArray indicates predicate ops must be an array.
	ErrPredicateOpsNotArray = errors.New("predicate ops must be an array")
	// ErrPredicateOpNotArray indicates predicate op must be an array.
	ErrPredicateOpNotArray = errors.New("predicate op must be an array")
	// ErrDecodedOpNotPredicate indicates decoded operation is not a predicate.
	ErrDecodedOpNotPredicate = errors.New("decoded operation is not a predicate")
)

// --- Dynamic errors ---

var (
	// ErrUnsupportedOp indicates unsupported operation type.
	ErrUnsupportedOp = errors.New("unsupported operation type")
	// ErrUnknownStringOpcode indicates unknown string opcode.
	ErrUnknownStringOpcode = errors.New("unknown string opcode")
	// ErrInvalidOpcodeType indicates invalid opcode type.
	ErrInvalidOpcodeType = errors.New("invalid opcode type")
	// ErrUnknownNumericOpcode indicates unknown numeric opcode.
	ErrUnknownNumericOpcode = errors.New("unknown numeric opcode")
)

// --- Type conversion errors ---

var (
	// ErrCannotConvertToFloat64 indicates value cannot be converted to float64.
	ErrCannotConvertToFloat64 = errors.New("cannot convert to float64")
	// ErrExpectedArray indicates expected array value.
	ErrExpectedArray = errors.New("expected array")
	// ErrExpectedStringInArray indicates expected string in array.
	ErrExpectedStringInArray = errors.New("expected string in array")
)
