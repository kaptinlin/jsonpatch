package json

import "errors"

// Errors for base operation decoding.
var (
	ErrOpMissingOpField   = errors.New("operation missing 'op' field")
	ErrOpMissingPathField = errors.New("operation missing 'path' field")
	ErrInvalidPointer     = errors.New("invalid pointer")
	ErrCodecOpUnknown     = errors.New("unknown operation")
)

// Errors for core operation (RFC 6902) decoding.
var (
	ErrAddOpMissingValue     = errors.New("add operation missing 'value' field")
	ErrReplaceOpMissingValue = errors.New("replace operation missing 'value' field")
	ErrMissingValueField     = errors.New("missing value field")
	ErrMoveOpMissingFrom     = errors.New("move operation missing 'from' field")
	ErrCopyOpMissingFrom     = errors.New("copy operation missing 'from' field")
)

// Errors for extended operation decoding.
var (
	ErrIncOpMissingInc       = errors.New("inc operation missing 'inc' field")
	ErrIncOpInvalidType      = errors.New("inc operation 'inc' field must be a number")
	ErrStrInsOpMissingPos    = errors.New("str_ins operation missing 'pos' field")
	ErrStrInsOpMissingStr    = errors.New("str_ins operation missing 'str' field")
	ErrStrDelOpMissingPos    = errors.New("str_del operation missing 'pos' field")
	ErrStrDelOpMissingFields = errors.New("str_del operation missing 'str' or 'len' field")
	ErrSplitOpMissingPos     = errors.New("split operation missing 'pos' field")
	ErrValueNotObject        = errors.New("value is not an object")
)

// Errors for predicate operation decoding.
var (
	ErrTypeOpMissingValue        = errors.New("type operation missing string 'value' field")
	ErrTestTypeOpMissingType     = errors.New("test_type operation missing 'type' field")
	ErrTestStringOpMissingStr    = errors.New("test_string operation missing 'str' field")
	ErrTestStringLenOpMissingLen = errors.New("test_string_len operation missing 'len' field")
	ErrContainsOpMissingValue    = errors.New("contains operation missing 'value' field")
	ErrEndsOpMissingValue        = errors.New("ends operation missing 'value' field")
	ErrStartsOpMissingValue      = errors.New("starts operation missing 'value' field")
	ErrMatchesOpMissingValue     = errors.New("matches operation missing 'value' field")
	ErrLessOpMissingValue        = errors.New("less operation missing 'value' field")
	ErrMoreOpMissingValue        = errors.New("more operation missing 'value' field")
	ErrInvalidType               = errors.New("invalid type")
	ErrEmptyTypeList             = errors.New("empty type list")
)

// Errors for composite operation decoding.
var (
	ErrAndOpMissingApply         = errors.New("and operation missing 'apply' field")
	ErrOrOpMissingApply          = errors.New("or operation missing 'apply' field")
	ErrNotOpMissingApply         = errors.New("not operation missing 'apply' field")
	ErrNotOpRequiresOperand      = errors.New("not operation requires at least one operand")
	ErrNotOpRequiresValidOperand = errors.New("not operation requires a valid predicate operand")
)
