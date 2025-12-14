package compact

import (
	"errors"
)

// Predefined errors for compact codec operations
var (
	// Base operation errors
	ErrCompactOperationMinLength     = errors.New("compact operation must have at least opcode and path")
	ErrCompactOperationPathNotString = errors.New("compact operation path must be a string")

	// Operation-specific errors
	ErrAddOperationRequiresValue       = errors.New("add operation requires value")
	ErrReplaceOperationRequiresValue   = errors.New("replace operation requires value")
	ErrMoveOperationRequiresFrom       = errors.New("move operation requires from path")
	ErrMoveOperationFromNotString      = errors.New("move operation from must be a string")
	ErrCopyOperationRequiresFrom       = errors.New("copy operation requires from path")
	ErrCopyOperationFromNotString      = errors.New("copy operation from must be a string")
	ErrTestOperationRequiresValue      = errors.New("test operation requires value")
	ErrIncOperationRequiresDelta       = errors.New("inc operation requires delta")
	ErrIncOperationDeltaNotNumber      = errors.New("inc operation delta must be a number")
	ErrContainsOperationRequiresValue  = errors.New("contains operation requires value")
	ErrContainsOperationValueNotString = errors.New("contains operation value must be a string")
	ErrStartsOperationRequiresValue    = errors.New("starts operation requires value")
	ErrStartsOperationValueNotString   = errors.New("starts operation value must be a string")
	ErrEndsOperationRequiresValue      = errors.New("ends operation requires value")
	ErrEndsOperationValueNotString     = errors.New("ends operation value must be a string")

	// Type operation errors
	ErrTypeOperationRequiresType  = errors.New("type operation requires type")
	ErrTypeOperationTypeNotString = errors.New("type operation type must be a string")

	// Test type operation errors
	ErrTestTypeOperationRequiresTypes = errors.New("test_type operation requires types")
	ErrTestTypeOperationTypesNotArray = errors.New("test_type operation types must be an array")

	// Test string operation errors
	ErrTestStringOperationRequiresStr  = errors.New("test_string operation requires str")
	ErrTestStringOperationStrNotString = errors.New("test_string operation str must be a string")

	// Test string len operation errors
	ErrTestStringLenOperationRequiresLen  = errors.New("test_string_len operation requires len")
	ErrTestStringLenOperationLenNotNumber = errors.New("test_string_len operation len must be a number")

	// In operation errors
	ErrInOperationRequiresValues = errors.New("in operation requires values")
	ErrInOperationValuesNotArray = errors.New("in operation values must be an array")

	// Less operation errors
	ErrLessOperationRequiresValue  = errors.New("less operation requires value")
	ErrLessOperationValueNotNumber = errors.New("less operation value must be a number")

	// More operation errors
	ErrMoreOperationRequiresValue  = errors.New("more operation requires value")
	ErrMoreOperationValueNotNumber = errors.New("more operation value must be a number")

	// Matches operation errors
	ErrMatchesOperationRequiresPattern  = errors.New("matches operation requires pattern")
	ErrMatchesOperationPatternNotString = errors.New("matches operation pattern must be a string")

	// Composite operation errors
	ErrAndOperationRequiresOps = errors.New("and operation requires ops")
	ErrOrOperationRequiresOps  = errors.New("or operation requires ops")
	ErrNotOperationRequiresOps = errors.New("not operation requires ops")
	ErrPredicateOpsNotArray    = errors.New("predicate ops must be an array")
	ErrPredicateOpNotArray     = errors.New("predicate op must be an array")
	ErrDecodedOpNotPredicate   = errors.New("decoded operation is not a predicate")

	// String operation errors
	ErrStrInsOperationRequiresPosAndStr = errors.New("str_ins operation requires pos and str")
	ErrStrInsOperationPosNotNumber      = errors.New("str_ins operation pos must be a number")
	ErrStrInsOperationStrNotString      = errors.New("str_ins operation str must be a string")
	ErrStrDelOperationRequiresPosAndLen = errors.New("str_del operation requires pos and len")
	ErrStrDelOperationPosNotNumber      = errors.New("str_del operation pos must be a number")
	ErrStrDelOperationLenNotNumber      = errors.New("str_del operation len must be a number")

	// Split/Merge/Extend operation errors
	ErrSplitOperationRequiresPos     = errors.New("split operation requires pos")
	ErrSplitOperationPosNotNumber    = errors.New("split operation pos must be a number")
	ErrMergeOperationRequiresPos     = errors.New("merge operation requires pos")
	ErrMergeOperationPosNotNumber    = errors.New("merge operation pos must be a number")
	ErrExtendOperationRequiresProps  = errors.New("extend operation requires props")
	ErrExtendOperationPropsNotObject = errors.New("extend operation props must be an object")
)

// Base errors for dynamic data
var (
	ErrUnsupportedOperationType = errors.New("unsupported operation type")
	ErrUnknownStringOpcode      = errors.New("unknown string opcode")
	ErrInvalidOpcodeType        = errors.New("invalid opcode type")
	ErrUnknownNumericOpcode     = errors.New("unknown numeric opcode")
)

// Type conversion errors
var (
	ErrCannotConvertToFloat64 = errors.New("cannot convert to float64")
	ErrExpectedArray          = errors.New("expected array")
	ErrExpectedStringInArray  = errors.New("expected string in array")
)
