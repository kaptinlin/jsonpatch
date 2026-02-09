package binary

import "errors"

var (
	// ErrUnsupportedOp indicates an unknown or unsupported operation code.
	ErrUnsupportedOp = errors.New("unsupported operation code")
	// ErrInvalidTestTypeFormat indicates invalid types array for test_type predicate.
	ErrInvalidTestTypeFormat = errors.New("invalid test_type types format")
	// ErrInvalidValueType indicates the decoded value has an unexpected type.
	ErrInvalidValueType = errors.New("invalid value type")
)
