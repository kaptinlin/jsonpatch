package op

import "errors"

// Sentinel errors for path and validation related operations
var (
	// Core path errors
	ErrPathNotFound   = errors.New("path not found")
	ErrInvalidPath    = errors.New("invalid path")
	ErrPathEmpty      = errors.New("path cannot be empty")
	ErrFromPathEmpty  = errors.New("from path cannot be empty")
	ErrPathsIdentical = errors.New("cannot move into own children")

	// Array operation errors
	ErrIndexOutOfRange     = errors.New("index out of range")
	ErrNotAnArray          = errors.New("not an array")
	ErrArrayTooSmall       = errors.New("array must have at least 2 elements")
	ErrPositionOutOfBounds = errors.New("position out of bounds")
	ErrPositionNegative    = errors.New("position cannot be negative")
	ErrInvalidTarget       = errors.New("invalid target")

	// Type validation errors
	ErrNotString                 = errors.New("value is not a string")
	ErrNotNumber                 = errors.New("value must be a number")
	ErrNotObject                 = errors.New("value is not an object")
	ErrInvalidType               = errors.New("invalid type")
	ErrEmptyTypeList             = errors.New("empty type list")
	ErrContainsValueMustBeString = errors.New("contains operation value must be a string")

	// Operation execution errors
	ErrTestFailed          = errors.New("test failed")
	ErrDefinedTestFailed   = errors.New("defined test failed")
	ErrUndefinedTestFailed = errors.New("undefined test failed")
	ErrAndTestFailed       = errors.New("and test failed")
	ErrOrTestFailed        = errors.New("or test failed")
	ErrNotTestFailed       = errors.New("not test failed")

	// Value operation errors
	ErrCannotReplace          = errors.New("NOT_FOUND")
	ErrCannotAddToValue       = errors.New("cannot add to non-object/non-array value")
	ErrCannotRemoveFromValue  = errors.New("cannot remove from non-object/non-array document")
	ErrCannotMoveIntoChildren = errors.New("cannot move into own children")
	ErrPropertiesNil          = errors.New("properties cannot be nil")
	ErrValuesArrayEmpty       = errors.New("'in' operation 'value' must be an array")

	// Key type errors
	ErrInvalidKeyTypeMap     = errors.New("invalid key type for map")
	ErrInvalidKeyTypeSlice   = errors.New("invalid key type for slice")
	ErrUnsupportedParentType = errors.New("unsupported parent type")

	// String operation errors
	ErrPositionOutOfStringRange = errors.New("position out of string range")
	ErrSubstringTooLong         = errors.New("value too long")
	ErrSubstringMismatch        = errors.New("substring does not match")
	ErrStringLengthMismatch     = errors.New("string length mismatch")
	ErrPatternEmpty             = errors.New("pattern cannot be empty")
	ErrLengthNegative           = errors.New("length cannot be negative")

	// Type comparison errors
	ErrTypeMismatch     = errors.New("type mismatch")
	ErrContainsMismatch = errors.New("contains check failed")

	// Predicate operation errors
	ErrInvalidPredicateInAnd = errors.New("invalid predicate in and operation")
	ErrInvalidPredicateInNot = errors.New("invalid predicate in not operation")
	ErrInvalidPredicateInOr  = errors.New("invalid predicate in or operation")
	ErrNotNoOperands         = errors.New("not operation requires operands")

	// Operation modification errors
	ErrCannotModifyRootArray   = errors.New("cannot modify root array directly")
	ErrCannotUpdateParent      = errors.New("cannot update parent")
	ErrCannotUpdateGrandparent = errors.New("cannot update grandparent")

	// Value conversion errors
	ErrCannotConvertNilToString = errors.New("cannot convert nil to string")
	ErrCannotConvertToString    = errors.New("cannot convert value to string")

	// Test operation errors
	ErrTestOperationNumberStringMismatch = errors.New("number is not equal to string")
	ErrTestOperationStringNotEquivalent  = errors.New("string not equivalent")

	// Base errors for dynamic wrapping with fmt.Errorf
	ErrComparisonFailed    = errors.New("comparison failed")
	ErrStringMismatch      = errors.New("string mismatch")
	ErrTestOperationFailed = errors.New("test operation failed")
	ErrInvalidIndex        = errors.New("invalid index")
	ErrRegexPattern        = errors.New("regex pattern error")
	ErrOperationFailed     = errors.New("operation failed")
)
