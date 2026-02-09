package json

import "errors"

// Base operation errors.
var (
	// ErrOpMissingOpField indicates operation is missing 'op' field.
	ErrOpMissingOpField = errors.New("operation missing 'op' field")
	// ErrOpMissingPathField indicates operation is missing 'path' field.
	ErrOpMissingPathField = errors.New("operation missing 'path' field")
	// ErrInvalidPointer indicates an invalid JSON pointer.
	ErrInvalidPointer = errors.New("invalid pointer")
	// ErrCodecOpUnknown indicates unknown operation.
	ErrCodecOpUnknown = errors.New("unknown operation")
)

// Core operation errors.
var (
	// ErrAddOpMissingValue indicates add operation is missing 'value' field.
	ErrAddOpMissingValue = errors.New("add operation missing 'value' field")
	// ErrReplaceOpMissingValue indicates replace operation is missing 'value' field.
	ErrReplaceOpMissingValue = errors.New("replace operation missing 'value' field")
	// ErrMissingValueField indicates a missing value field.
	ErrMissingValueField = errors.New("missing value field")
	// ErrMoveOpMissingFrom indicates move operation is missing 'from' field.
	ErrMoveOpMissingFrom = errors.New("move operation missing 'from' field")
	// ErrCopyOpMissingFrom indicates copy operation is missing 'from' field.
	ErrCopyOpMissingFrom = errors.New("copy operation missing 'from' field")
)

// Extended operation errors.
var (
	// ErrIncOpMissingInc indicates inc operation is missing 'inc' field.
	ErrIncOpMissingInc = errors.New("inc operation missing 'inc' field")
	// ErrIncOpInvalidType indicates inc operation 'inc' field must be a number.
	ErrIncOpInvalidType = errors.New("inc operation 'inc' field must be a number")
	// ErrStrInsOpMissingPos indicates str_ins operation is missing 'pos' field.
	ErrStrInsOpMissingPos = errors.New("str_ins operation missing 'pos' field")
	// ErrStrInsOpMissingStr indicates str_ins operation is missing 'str' field.
	ErrStrInsOpMissingStr = errors.New("str_ins operation missing 'str' field")
	// ErrStrDelOpMissingPos indicates str_del operation is missing 'pos' field.
	ErrStrDelOpMissingPos = errors.New("str_del operation missing 'pos' field")
	// ErrStrDelOpMissingFields indicates str_del operation is missing 'str' or 'len' field.
	ErrStrDelOpMissingFields = errors.New("str_del operation missing 'str' or 'len' field")
	// ErrSplitOpMissingPos indicates split operation is missing 'pos' field.
	ErrSplitOpMissingPos = errors.New("split operation missing 'pos' field")
	// ErrValueNotObject indicates value is not an object.
	ErrValueNotObject = errors.New("value is not an object")
)

// Predicate operation errors.
var (
	// ErrTypeOpMissingValue indicates type operation is missing string 'value' field.
	ErrTypeOpMissingValue = errors.New("type operation missing string 'value' field")
	// ErrTestTypeOpMissingType indicates test_type operation is missing 'type' field.
	ErrTestTypeOpMissingType = errors.New("test_type operation missing 'type' field")
	// ErrTestStringOpMissingStr indicates test_string operation is missing 'str' field.
	ErrTestStringOpMissingStr = errors.New("test_string operation missing 'str' field")
	// ErrTestStringLenOpMissingLen indicates test_string_len operation is missing 'len' field.
	ErrTestStringLenOpMissingLen = errors.New("test_string_len operation missing 'len' field")
	// ErrContainsOpMissingValue indicates contains operation is missing 'value' field.
	ErrContainsOpMissingValue = errors.New("contains operation missing 'value' field")
	// ErrEndsOpMissingValue indicates ends operation is missing 'value' field.
	ErrEndsOpMissingValue = errors.New("ends operation missing 'value' field")
	// ErrStartsOpMissingValue indicates starts operation is missing 'value' field.
	ErrStartsOpMissingValue = errors.New("starts operation missing 'value' field")
	// ErrMatchesOpMissingValue indicates matches operation is missing 'value' field.
	ErrMatchesOpMissingValue = errors.New("matches operation missing 'value' field")
	// ErrLessOpMissingValue indicates less operation is missing 'value' field.
	ErrLessOpMissingValue = errors.New("less operation missing 'value' field")
	// ErrMoreOpMissingValue indicates more operation is missing 'value' field.
	ErrMoreOpMissingValue = errors.New("more operation missing 'value' field")
	// ErrInvalidType indicates an invalid type.
	ErrInvalidType = errors.New("invalid type")
	// ErrEmptyTypeList indicates an empty type list.
	ErrEmptyTypeList = errors.New("empty type list")
)

// Composite operation errors.
var (
	// ErrAndOpMissingApply indicates and operation is missing 'apply' field.
	ErrAndOpMissingApply = errors.New("and operation missing 'apply' field")
	// ErrOrOpMissingApply indicates or operation is missing 'apply' field.
	ErrOrOpMissingApply = errors.New("or operation missing 'apply' field")
	// ErrNotOpMissingApply indicates not operation is missing 'apply' field.
	ErrNotOpMissingApply = errors.New("not operation missing 'apply' field")
	// ErrNotOpRequiresOperand indicates not operation requires at least one operand.
	ErrNotOpRequiresOperand = errors.New("not operation requires at least one operand")
	// ErrNotOpRequiresValidOperand indicates not operation requires a valid predicate operand.
	ErrNotOpRequiresValidOperand = errors.New("not operation requires a valid predicate operand")
)
