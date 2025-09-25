package op

import "errors"

// Sentinel errors for path and validation related operations
var (
	// Core path errors - aligned with json-joy patterns
	ErrPathNotFound     = errors.New("NOT_FOUND")
	ErrPathDoesNotExist = errors.New("NOT_FOUND")
	ErrInvalidPath      = errors.New("OP_PATH_INVALID")
	ErrPathEmpty        = errors.New("OP_PATH_INVALID")
	ErrFromPathEmpty    = errors.New("OP_FROM_INVALID")
	ErrPathsIdentical   = errors.New("cannot move into own children")

	// Array operation errors - aligned with json-joy patterns
	ErrArrayIndexOutOfBounds = errors.New("INVALID_INDEX")
	ErrIndexOutOfRange       = errors.New("INVALID_INDEX")
	ErrNotAnArray            = errors.New("Not a array")
	ErrArrayTooSmall         = errors.New("array must have at least 2 elements")
	ErrPositionOutOfBounds   = errors.New("INVALID_INDEX")
	ErrPositionNegative      = errors.New("INVALID_INDEX")

	// Type validation errors - aligned with json-joy patterns
	ErrNotString     = errors.New("value is not a string")
	ErrNotNumber     = errors.New("value must be a number")
	ErrNotObject     = errors.New("value is not an object")
	ErrInvalidType   = errors.New("invalid type")
	ErrEmptyTypeList = errors.New("empty type list")

	// Operation execution errors - aligned with json-joy patterns
	ErrTestFailed          = errors.New("test failed")
	ErrDefinedTestFailed   = errors.New("defined test failed")
	ErrUndefinedTestFailed = errors.New("undefined test failed")
	ErrAndTestFailed       = errors.New("and test failed")
	ErrOrTestFailed        = errors.New("or test failed")
	ErrNotTestFailed       = errors.New("not test failed")

	// Value operation errors - aligned with json-joy patterns
	ErrCannotReplace          = errors.New("NOT_FOUND")
	ErrCannotAddToValue       = errors.New("cannot add to non-object/non-array value")
	ErrCannotRemoveFromValue  = errors.New("cannot remove from non-object/non-array document")
	ErrPathMissingRecursive   = errors.New("NOT_FOUND")
	ErrCannotMoveIntoChildren = errors.New("cannot move into own children")
	ErrPropertiesNil          = errors.New("properties cannot be nil")
	ErrValuesArrayEmpty       = errors.New("'in' operation 'value' must be an array")

	// Key type errors
	ErrInvalidKeyTypeMap     = errors.New("invalid key type for map")
	ErrInvalidKeyTypeSlice   = errors.New("invalid key type for slice")
	ErrUnsupportedParentType = errors.New("unsupported parent type")

	// String operation errors - aligned with json-joy patterns
	ErrPositionOutOfStringRange = errors.New("INVALID_INDEX")
	ErrSubstringTooLong         = errors.New("value too long")
	ErrSubstringMismatch        = errors.New("substring does not match")
	ErrStringLengthMismatch     = errors.New("string length mismatch")
	ErrPatternEmpty             = errors.New("pattern cannot be empty")
	ErrLengthNegative           = errors.New("INVALID_INDEX")

	// Type comparison errors
	ErrTypeMismatch = errors.New("type mismatch")

	// Predicate operation errors - aligned with json-joy patterns
	ErrInvalidPredicateInAnd = errors.New("OP_INVALID")
	ErrInvalidPredicateInNot = errors.New("OP_INVALID")
	ErrInvalidPredicateInOr  = errors.New("OP_INVALID")
	ErrAndNoOperands         = errors.New("empty operation patch")
	ErrNotNoOperands         = errors.New("empty operation patch")
	ErrOrNoOperands          = errors.New("empty operation patch")

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
