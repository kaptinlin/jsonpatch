package json

import "errors"

var (
	// ErrOpMissingOpField indicates operation is missing 'op' field
	ErrOpMissingOpField = errors.New("operation missing 'op' field")
	// ErrOpMissingPathField indicates operation is missing 'path' field
	ErrOpMissingPathField = errors.New("operation missing 'path' field")

	// ErrMoveOpMissingFrom indicates move operation is missing 'from' field
	ErrMoveOpMissingFrom = errors.New("move operation missing 'from' field")
	// ErrCopyOpMissingFrom indicates copy operation is missing 'from' field
	ErrCopyOpMissingFrom = errors.New("copy operation missing 'from' field")

	// ErrIncOpMissingInc indicates inc operation is missing 'inc' field
	ErrIncOpMissingInc = errors.New("inc operation missing 'inc' field")
	// ErrStrInsOpMissingPos indicates str_ins operation is missing 'pos' field
	ErrStrInsOpMissingPos = errors.New("str_ins operation missing 'pos' field")
	// ErrStrInsOpMissingStr indicates str_ins operation is missing 'str' field
	ErrStrInsOpMissingStr = errors.New("str_ins operation missing 'str' field")
	// ErrStrDelOpMissingPos indicates str_del operation is missing 'pos' field
	ErrStrDelOpMissingPos = errors.New("str_del operation missing 'pos' field")
	// ErrStrDelOpMissingFields indicates str_del operation is missing 'str' or 'len' field
	ErrStrDelOpMissingFields = errors.New("str_del operation missing 'str' or 'len' field")
	// ErrSplitOpMissingPos indicates split operation is missing 'pos' field
	ErrSplitOpMissingPos = errors.New("split operation missing 'pos' field")
	// ErrMergeOpMissingPos indicates merge operation is missing 'pos' field
	ErrMergeOpMissingPos = errors.New("merge operation missing 'pos' field")

	// ErrTypeOpMissingValue indicates type operation is missing string 'value' field
	ErrTypeOpMissingValue = errors.New("type operation missing string 'value' field")
	// ErrTestTypeOpMissingType indicates test_type operation is missing 'type' field
	ErrTestTypeOpMissingType = errors.New("test_type operation missing 'type' field")
	// ErrTestStringOpMissingStr indicates test_string operation is missing 'str' field
	ErrTestStringOpMissingStr = errors.New("test_string operation missing 'str' field")
	// ErrTestStringLenOpMissingLen indicates test_string_len operation is missing 'len' field
	ErrTestStringLenOpMissingLen = errors.New("test_string_len operation missing 'len' field")
	// ErrContainsOpMissingValue indicates contains operation is missing 'value' field
	ErrContainsOpMissingValue = errors.New("contains operation missing 'value' field")
	// ErrEndsOpMissingValue indicates ends operation is missing 'value' field
	ErrEndsOpMissingValue = errors.New("ends operation missing 'value' field")
	// ErrStartsOpMissingValue indicates starts operation is missing 'value' field
	ErrStartsOpMissingValue = errors.New("starts operation missing 'value' field")
	// ErrMatchesOpMissingValue indicates matches operation is missing 'value' field
	ErrMatchesOpMissingValue = errors.New("matches operation missing 'value' field")
	// ErrLessOpMissingValue indicates less operation is missing 'value' field
	ErrLessOpMissingValue = errors.New("less operation missing 'value' field")
	// ErrMoreOpMissingValue indicates more operation is missing 'value' field
	ErrMoreOpMissingValue = errors.New("more operation missing 'value' field")

	// ErrAndOpMissingApply indicates and operation is missing 'apply' field
	ErrAndOpMissingApply = errors.New("and operation missing 'apply' field")
	// ErrOrOpMissingApply indicates or operation is missing 'apply' field
	ErrOrOpMissingApply = errors.New("or operation missing 'apply' field")
	// ErrNotOpMissingApply indicates not operation is missing 'apply' field
	ErrNotOpMissingApply = errors.New("not operation missing 'apply' field")

	// ErrValueNotObject indicates value is not an object
	ErrValueNotObject = errors.New("value is not an object")

	// ErrCodecOpUnknown indicates unknown operation
	ErrCodecOpUnknown = errors.New("unknown operation")
)
