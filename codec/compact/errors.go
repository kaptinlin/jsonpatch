package compact

import "errors"

// Base errors for compact operation validation.
var (
	ErrMinLength     = errors.New("compact operation must have at least opcode and path")
	ErrPathNotString = errors.New("compact operation path must be a string")
)

// Core operation (RFC 6902) errors.
var (
	ErrAddMissingValue     = errors.New("add operation requires value")
	ErrReplaceMissingValue = errors.New("replace operation requires value")
	ErrMoveMissingFrom     = errors.New("move operation requires from path")
	ErrMoveFromNotString   = errors.New("move operation from must be a string")
	ErrCopyMissingFrom     = errors.New("copy operation requires from path")
	ErrCopyFromNotString   = errors.New("copy operation from must be a string")
	ErrTestMissingValue    = errors.New("test operation requires value")
)

// Extended operation errors.
var (
	ErrIncMissingDelta      = errors.New("inc operation requires delta")
	ErrIncDeltaNotNumber    = errors.New("inc operation delta must be a number")
	ErrStrInsMissingFields  = errors.New("str_ins operation requires pos and str")
	ErrStrInsPosNotNumber   = errors.New("str_ins operation pos must be a number")
	ErrStrInsStrNotString   = errors.New("str_ins operation str must be a string")
	ErrStrDelMissingFields  = errors.New("str_del operation requires pos and len")
	ErrStrDelPosNotNumber   = errors.New("str_del operation pos must be a number")
	ErrStrDelLenNotNumber   = errors.New("str_del operation len must be a number")
	ErrSplitMissingPos      = errors.New("split operation requires pos")
	ErrSplitPosNotNumber    = errors.New("split operation pos must be a number")
	ErrMergeMissingPos      = errors.New("merge operation requires pos")
	ErrMergePosNotNumber    = errors.New("merge operation pos must be a number")
	ErrExtendMissingProps   = errors.New("extend operation requires props")
	ErrExtendPropsNotObject = errors.New("extend operation props must be an object")
)

// Predicate operation errors.
var (
	ErrContainsMissingValue    = errors.New("contains operation requires value")
	ErrContainsValueNotString  = errors.New("contains operation value must be a string")
	ErrStartsMissingValue      = errors.New("starts operation requires value")
	ErrStartsValueNotString    = errors.New("starts operation value must be a string")
	ErrEndsMissingValue        = errors.New("ends operation requires value")
	ErrEndsValueNotString      = errors.New("ends operation value must be a string")
	ErrTypeMissingType         = errors.New("type operation requires type")
	ErrTypeNotString           = errors.New("type operation type must be a string")
	ErrTestTypeMissingTypes    = errors.New("test_type operation requires types")
	ErrTestTypeTypesNotArray   = errors.New("test_type operation types must be an array")
	ErrTestStringMissingStr    = errors.New("test_string operation requires str")
	ErrTestStringNotString     = errors.New("test_string operation str must be a string")
	ErrTestStringLenMissingLen = errors.New("test_string_len operation requires len")
	ErrTestStringLenNotNumber  = errors.New("test_string_len operation len must be a number")
	ErrInMissingValues         = errors.New("in operation requires values")
	ErrInValuesNotArray        = errors.New("in operation values must be an array")
	ErrLessMissingValue        = errors.New("less operation requires value")
	ErrLessValueNotNumber      = errors.New("less operation value must be a number")
	ErrMoreMissingValue        = errors.New("more operation requires value")
	ErrMoreValueNotNumber      = errors.New("more operation value must be a number")
	ErrMatchesMissingPattern   = errors.New("matches operation requires pattern")
	ErrMatchesPatternNotString = errors.New("matches operation pattern must be a string")
)

// Composite operation errors.
var (
	ErrAndMissingOps      = errors.New("and operation requires ops")
	ErrOrMissingOps       = errors.New("or operation requires ops")
	ErrNotMissingOps      = errors.New("not operation requires ops")
	ErrPredicateNotArray  = errors.New("predicate ops must be an array")
	ErrPredicateOpInvalid = errors.New("predicate op must be an array")
	ErrNotPredicate       = errors.New("decoded operation is not a predicate")
)

// Resolution and conversion errors.
var (
	ErrUnsupportedOp      = errors.New("unsupported operation type")
	ErrUnknownStringCode  = errors.New("unknown string opcode")
	ErrInvalidCodeType    = errors.New("invalid opcode type")
	ErrUnknownNumericCode = errors.New("unknown numeric opcode")
	ErrNotFloat64         = errors.New("cannot convert to float64")
	ErrExpectedArray      = errors.New("expected array")
	ErrExpectedString     = errors.New("expected string in array")
)
