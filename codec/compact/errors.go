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
)

// Base errors for dynamic data
var (
	ErrUnsupportedOperationType = errors.New("unsupported operation type")
	ErrUnknownStringOpcode      = errors.New("unknown string opcode")
	ErrInvalidOpcodeType        = errors.New("invalid opcode type")
	ErrUnknownNumericOpcode     = errors.New("unknown numeric opcode")
)
