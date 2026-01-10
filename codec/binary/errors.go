package binary

import "errors"

var (
	// ErrUnsupportedOp indicates an unknown or unsupported operation code.
	ErrUnsupportedOp = errors.New("binary: unsupported operation code")
	// ErrInvalidTestTypeFormat indicates invalid types array for test_type predicate.
	ErrInvalidTestTypeFormat = errors.New("binary: invalid test_type types format")
	// ErrInvalidValueType indicates the decoded value has an unexpected type.
	ErrInvalidValueType = errors.New("binary: invalid value type")
)
