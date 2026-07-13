package binary

import "errors"

var (
	// ErrUnsupportedOp indicates an unknown or unsupported operation code.
	ErrUnsupportedOp = errors.New("unsupported operation code")
	// ErrInvalidPredicate indicates a composite predicate contains a non-predicate operation.
	ErrInvalidPredicate = errors.New("invalid predicate operation")
	// ErrNotSinglePredicate indicates a not operation has anything other than one child predicate.
	ErrNotSinglePredicate = errors.New("not operation requires exactly one predicate")
	// ErrInvalidTestTypeFormat indicates invalid types array for test_type predicate.
	ErrInvalidTestTypeFormat = errors.New("invalid test_type types format")
	// ErrInvalidValueType indicates the decoded value has an unexpected type.
	ErrInvalidValueType = errors.New("invalid value type")
	// ErrInvalidFrameArity indicates an operation frame has the wrong number of fields.
	ErrInvalidFrameArity = errors.New("invalid operation frame arity")
	// ErrTrailingData indicates bytes remain after a complete patch value.
	ErrTrailingData = errors.New("trailing data after patch")
)
