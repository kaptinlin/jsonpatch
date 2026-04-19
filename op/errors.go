package op

import "errors"

var (
	// ErrPathNotFound reports that a JSON Pointer does not resolve.
	ErrPathNotFound = errors.New("path not found")
	// ErrInvalidPath reports that a JSON Pointer is malformed.
	ErrInvalidPath = errors.New("invalid path")
	// ErrPathEmpty reports that path is required but empty.
	ErrPathEmpty = errors.New("path cannot be empty")
	// ErrFromPathEmpty reports that from is required but empty.
	ErrFromPathEmpty = errors.New("from path cannot be empty")
	// ErrPathsIdentical reports that a move target is nested under its source.
	ErrPathsIdentical = errors.New("cannot move into own children")
	// ErrIndexOutOfRange reports that an array index is outside the valid bounds.
	ErrIndexOutOfRange = errors.New("index out of range")
	// ErrNotAnArray reports that an operation expected an array value.
	ErrNotAnArray = errors.New("not an array")
	// ErrArrayTooSmall reports that an array does not have enough elements.
	ErrArrayTooSmall = errors.New("array must have at least 2 elements")
	// ErrPositionOutOfBounds reports that a position is outside the valid range.
	ErrPositionOutOfBounds = errors.New("position out of bounds")
	// ErrPositionNegative reports that a position is negative.
	ErrPositionNegative = errors.New("position cannot be negative")
	// ErrInvalidPosition reports that a position cannot be interpreted.
	ErrInvalidPosition = errors.New("invalid position")
	// ErrMissingStrOrLen reports that str_del requires either str or len.
	ErrMissingStrOrLen = errors.New("str_del requires either str or len")
	// ErrInvalidTarget reports that the operation target is unsupported.
	ErrInvalidTarget = errors.New("invalid target")
	// ErrNotString reports that a string operation received a non-string value.
	ErrNotString = errors.New("value is not a string")
	// ErrNotNumber reports that a numeric operation received a non-number value.
	ErrNotNumber = errors.New("value must be a number")
	// ErrNotObject reports that an object operation received a non-object value.
	ErrNotObject = errors.New("value is not an object")
	// ErrInvalidType reports that a JSON type name is unsupported.
	ErrInvalidType = errors.New("invalid type")
	// ErrEmptyTypeList reports that a type list is empty.
	ErrEmptyTypeList = errors.New("empty type list")
	// ErrContainsValueMustBeString reports that contains requires a string value.
	ErrContainsValueMustBeString = errors.New("contains operation value must be a string")
	// ErrTestFailed reports that a test operation did not match.
	ErrTestFailed = errors.New("test failed")
	// ErrDefinedTestFailed reports that a defined predicate failed.
	ErrDefinedTestFailed = errors.New("defined test failed")
	// ErrUndefinedTestFailed reports that an undefined predicate failed.
	ErrUndefinedTestFailed = errors.New("undefined test failed")
	// ErrAndTestFailed reports that an and predicate failed.
	ErrAndTestFailed = errors.New("and test failed")
	// ErrOrTestFailed reports that an or predicate failed.
	ErrOrTestFailed = errors.New("or test failed")
	// ErrNotTestFailed reports that a not predicate failed.
	ErrNotTestFailed = errors.New("not test failed")
	// ErrCannotReplace reports that replace targeted a missing path.
	ErrCannotReplace = errors.New("path not found for replace")
	// ErrCannotAddToValue reports that add targeted a scalar value.
	ErrCannotAddToValue = errors.New("cannot add to non-object/non-array value")
	// ErrCannotRemoveFromValue reports that remove targeted a scalar document.
	ErrCannotRemoveFromValue = errors.New("cannot remove from non-object/non-array document")
	// ErrCannotMoveIntoChildren reports that move targeted a descendant of its source.
	ErrCannotMoveIntoChildren = errors.New("cannot move into own children")
	// ErrPropertiesNil reports that extend received nil properties.
	ErrPropertiesNil = errors.New("properties cannot be nil")
	// ErrValuesArrayEmpty reports that in requires an array value.
	ErrValuesArrayEmpty = errors.New("'in' operation 'value' must be an array")
	// ErrInvalidKeyTypeMap reports that a map key has the wrong type.
	ErrInvalidKeyTypeMap = errors.New("invalid key type for map")
	// ErrInvalidKeyTypeSlice reports that a slice index has the wrong type.
	ErrInvalidKeyTypeSlice = errors.New("invalid key type for slice")
	// ErrUnsupportedParentType reports that the parent container type is unsupported.
	ErrUnsupportedParentType = errors.New("unsupported parent type")
	// ErrPositionOutOfStringRange reports that a string position is outside the valid range.
	ErrPositionOutOfStringRange = errors.New("position out of string range")
	// ErrSubstringTooLong reports that a substring length exceeds the source string.
	ErrSubstringTooLong = errors.New("value too long")
	// ErrSubstringMismatch reports that a required substring did not match.
	ErrSubstringMismatch = errors.New("substring does not match")
	// ErrStringLengthMismatch reports that a string length check failed.
	ErrStringLengthMismatch = errors.New("string length mismatch")
	// ErrPatternEmpty reports that a regex pattern is empty.
	ErrPatternEmpty = errors.New("pattern cannot be empty")
	// ErrLengthNegative reports that a length is negative.
	ErrLengthNegative = errors.New("length cannot be negative")
	// ErrInvalidLength reports that a length is not an integer.
	ErrInvalidLength = errors.New("length must be an integer")
	// ErrTypeMismatch reports that two values have incompatible JSON types.
	ErrTypeMismatch = errors.New("type mismatch")
	// ErrContainsMismatch reports that a contains predicate failed.
	ErrContainsMismatch = errors.New("contains check failed")
	// ErrInvalidPredicateInAnd reports that and received a non-predicate operand.
	ErrInvalidPredicateInAnd = errors.New("invalid predicate in and operation")
	// ErrInvalidPredicateInNot reports that not received a non-predicate operand.
	ErrInvalidPredicateInNot = errors.New("invalid predicate in not operation")
	// ErrInvalidPredicateInOr reports that or received a non-predicate operand.
	ErrInvalidPredicateInOr = errors.New("invalid predicate in or operation")
	// ErrNotNoOperands reports that not received no operands.
	ErrNotNoOperands = errors.New("not operation requires operands")
	// ErrCannotModifyRootArray reports that the root array cannot be edited in place.
	ErrCannotModifyRootArray = errors.New("cannot modify root array directly")
	// ErrCannotUpdateParent reports that the parent container could not be updated.
	ErrCannotUpdateParent = errors.New("cannot update parent")
	// ErrCannotUpdateGrandparent reports that the grandparent container could not be updated.
	ErrCannotUpdateGrandparent = errors.New("cannot update grandparent")
	// ErrCannotConvertNilToString reports that nil cannot be converted to a string.
	ErrCannotConvertNilToString = errors.New("cannot convert nil to string")
	// ErrCannotConvertToString reports that a value cannot be converted to a string.
	ErrCannotConvertToString = errors.New("cannot convert value to string")
	// ErrTestOperationNumberStringMismatch reports that a number was compared to a string.
	ErrTestOperationNumberStringMismatch = errors.New("number is not equal to string")
	// ErrTestOperationStringNotEquivalent reports that two strings are not equivalent.
	ErrTestOperationStringNotEquivalent = errors.New("string not equivalent")
	// ErrComparisonFailed is the base error for wrapped comparison failures.
	ErrComparisonFailed = errors.New("comparison failed")
	// ErrStringMismatch is the base error for wrapped string comparison failures.
	ErrStringMismatch = errors.New("string mismatch")
	// ErrTestOperationFailed is the base error for wrapped test failures.
	ErrTestOperationFailed = errors.New("test operation failed")
	// ErrInvalidIndex is the base error for wrapped index failures.
	ErrInvalidIndex = errors.New("invalid index")
	// ErrRegexPattern is the base error for wrapped regex failures.
	ErrRegexPattern = errors.New("regex pattern error")
	// ErrOperationFailed is the base error for wrapped operation failures.
	ErrOperationFailed = errors.New("operation failed")
)
