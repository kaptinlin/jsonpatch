package op

import "errors"

// Sentinel errors for path and validation related operations
var (
	// Core path errors
	ErrPathNotFound     = errors.New("path not found")
	ErrPathDoesNotExist = errors.New("path does not exist")
	ErrInvalidPath      = errors.New("invalid path")
	ErrPathEmpty        = errors.New("path cannot be empty")
	ErrFromPathEmpty    = errors.New("from path cannot be empty")
	ErrPathsIdentical   = errors.New("path and from cannot be the same")

	// Array operation errors
	ErrArrayIndexOutOfBounds = errors.New("array index out of bounds")
	ErrIndexOutOfRange       = errors.New("index out of range")
	ErrNotAnArray            = errors.New("not an array")
	ErrArrayTooSmall         = errors.New("array must have at least 2 elements")
	ErrPositionOutOfBounds   = errors.New("position out of bounds")
	ErrPositionNegative      = errors.New("position cannot be negative")

	// Type validation errors
	ErrNotString     = errors.New("value is not a string")
	ErrNotNumber     = errors.New("value is not a number")
	ErrNotObject     = errors.New("value is not an object")
	ErrInvalidType   = errors.New("invalid type")
	ErrEmptyTypeList = errors.New("types cannot be empty")

	// Operation execution errors
	ErrTestFailed          = errors.New("test failed")
	ErrDefinedTestFailed   = errors.New("defined test failed")
	ErrUndefinedTestFailed = errors.New("undefined test failed")
	ErrAndTestFailed       = errors.New("and test failed")
	ErrOrTestFailed        = errors.New("or test failed")
	ErrNotTestFailed       = errors.New("not test failed")

	// Value operation errors
	ErrCannotReplace          = errors.New("cannot replace key in non-object")
	ErrCannotAddToValue       = errors.New("cannot add to non-object/non-array value")
	ErrCannotRemoveFromValue  = errors.New("cannot remove from non-object/non-array document")
	ErrPathMissingRecursive   = errors.New("path does not exist -- missing objects are not created recursively")
	ErrCannotMoveIntoChildren = errors.New("cannot move into own children")
	ErrPropertiesNil          = errors.New("properties cannot be nil")
	ErrValuesArrayEmpty       = errors.New("values array cannot be empty")

	// Key type errors
	ErrInvalidKeyTypeMap     = errors.New("invalid key type for map")
	ErrInvalidKeyTypeSlice   = errors.New("invalid key type for slice")
	ErrUnsupportedParentType = errors.New("unsupported parent type")

	// String operation errors
	ErrPositionOutOfStringRange = errors.New("position out of range")
	ErrSubstringTooLong         = errors.New("substring extends beyond string length")
	ErrSubstringMismatch        = errors.New("substring does not match")
	ErrStringLengthMismatch     = errors.New("string length mismatch")
	ErrPatternEmpty             = errors.New("pattern cannot be empty")
	ErrLengthNegative           = errors.New("length cannot be negative")

	// Type comparison errors
	ErrTypeMismatch = errors.New("type mismatch")

	// Predicate operation errors
	ErrInvalidPredicateInAnd = errors.New("invalid predicate operation in AND")
	ErrInvalidPredicateInNot = errors.New("invalid predicate operation in NOT")
	ErrInvalidPredicateInOr  = errors.New("invalid predicate operation in OR")
	ErrAndNoOperands         = errors.New("and operation must have at least one operand")
	ErrNotNoOperands         = errors.New("not operation must have at least one operand")
	ErrOrNoOperands          = errors.New("or operation must have at least one operand")

	// Operation modification errors
	ErrCannotModifyRootArray     = errors.New("cannot modify root array directly")
	ErrCannotUpdateParent        = errors.New("cannot update parent")
	ErrCannotUpdateGrandparent   = errors.New("cannot update grandparent")
	ErrCannotAppendInPlace       = errors.New("cannot append to slice in place - caller must handle")
	ErrSliceDeletionNotSupported = errors.New("slice deletion not supported in place - caller must handle")
	ErrKeyDoesNotExist           = errors.New("key does not exist")

	// Value conversion errors
	ErrCannotConvertNilToString = errors.New("cannot convert nil to string")

	// Test operation errors
	ErrTestOperationNumberStringMismatch = errors.New("test operation failed: number is not equal to string")
	ErrTestOperationStringNotEquivalent  = errors.New("test operation failed: string not equivalent")

	// Base errors for dynamic wrapping with fmt.Errorf
	ErrComparisonFailed    = errors.New("comparison failed")
	ErrStringMismatch      = errors.New("string mismatch")
	ErrTestOperationFailed = errors.New("test operation failed")
	ErrInvalidIndex        = errors.New("invalid index")
	ErrRegexPattern        = errors.New("regex pattern error")
	ErrOperationFailed     = errors.New("operation failed")
)
