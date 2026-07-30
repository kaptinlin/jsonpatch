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
	// ErrMapKeyNotString indicates a MessagePack object contains a non-string key.
	ErrMapKeyNotString = errors.New("map key must be a string")
	// ErrInvalidFrameArity indicates an operation frame has the wrong number of fields.
	ErrInvalidFrameArity = errors.New("invalid operation frame arity")
	// ErrInvalidTrueMarker indicates an optional true flag is not encoded as integer 1.
	ErrInvalidTrueMarker = errors.New("binary optional true marker must be integer 1")
	// ErrTrailingData indicates bytes remain after a complete patch value.
	ErrTrailingData = errors.New("trailing data after patch")
	// ErrNumberNotFinite indicates a coordinate is NaN or infinite.
	ErrNumberNotFinite = errors.New("numeric field must be finite")
	// ErrNumberNotInteger indicates a coordinate has a fractional component.
	ErrNumberNotInteger = errors.New("coordinate field must be an integer")
	// ErrNumberOutOfRange indicates a coordinate cannot be represented as int.
	ErrNumberOutOfRange = errors.New("coordinate field is outside the int range")
)
