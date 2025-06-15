package json

import "errors"

// Core error definitions - simple and clear
var (
	// Core field validation errors
	ErrOpMissingOpField   = errors.New("operation missing 'op' field")
	ErrOpMissingPathField = errors.New("operation missing 'path' field")

	// Move/Copy operation errors
	ErrMoveOpMissingFrom = errors.New("move operation missing 'from' field")
	ErrCopyOpMissingFrom = errors.New("copy operation missing 'from' field")

	// Extended operation errors
	ErrIncOpMissingInc       = errors.New("inc operation missing 'inc' field")
	ErrStrInsOpMissingPos    = errors.New("str_ins operation missing 'pos' field")
	ErrStrInsOpMissingStr    = errors.New("str_ins operation missing 'str' field")
	ErrStrDelOpMissingPos    = errors.New("str_del operation missing 'pos' field")
	ErrStrDelOpMissingFields = errors.New("str_del operation missing 'str' or 'len' field")
	ErrSplitOpMissingPos     = errors.New("split operation missing 'pos' field")
	ErrMergeOpMissingPos     = errors.New("merge operation missing 'pos' field")

	// Predicate operation errors
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

	// Composite operation errors
	ErrAndOpMissingApply = errors.New("and operation missing 'apply' field")
	ErrOrOpMissingApply  = errors.New("or operation missing 'apply' field")
	ErrNotOpMissingApply = errors.New("not operation missing 'apply' field")

	// Type errors
	ErrValueNotObject = errors.New("value is not an object")

	// Unknown operation error
	ErrCodecOpUnknown = errors.New("unknown operation")
)
