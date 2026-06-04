package jsonpatch

import (
	"errors"
	"fmt"

	"github.com/kaptinlin/jsonpointer"

	"github.com/kaptinlin/jsonpatch/internal"
)

var (
	// ErrPayloadInvalid reports an invalid patch payload or compiled operation.
	ErrPayloadInvalid = errors.New("payload invalid")
	// ErrUnsupportedCapability reports an operation outside the enabled vocabulary.
	ErrUnsupportedCapability = errors.New("unsupported capability")
	// ErrRuntimeConflict reports a valid operation that cannot apply to the document state.
	ErrRuntimeConflict = errors.New("runtime conflict")
	// ErrTestFailed reports a failed predicate or test operation.
	ErrTestFailed = errors.New("test failed")
	// ErrTypeMismatch reports a runtime type mismatch.
	ErrTypeMismatch = errors.New("type mismatch")
	// ErrConversionFailed reports that a patched result could not be converted back.
	ErrConversionFailed = errors.New("failed to convert result back to original type")
)

// Error carries stable patch failure context for programmatic inspection.
type Error struct {
	kind  error
	index int
	op    string
	path  string
	from  string
	codec string
	cause error
}

func newError(kind error, index int, operation internal.Op, codec string, cause error) *Error {
	err := &Error{
		kind:  kind,
		index: index,
		codec: codec,
		cause: cause,
	}
	if operation == nil {
		return err
	}

	err.op = string(operation.Op())
	err.path = jsonpointer.Format(operation.Path()...)
	if from, ok := operation.(interface{ From() []string }); ok {
		err.from = jsonpointer.Format(from.From()...)
	}
	return err
}

func newPayloadError(codec string, cause error) *Error {
	return &Error{
		kind:  ErrPayloadInvalid,
		index: -1,
		codec: codec,
		cause: cause,
	}
}

func newFieldError(kind error, index int, opName, path, from, codec string, cause error) *Error {
	return &Error{
		kind:  kind,
		index: index,
		op:    opName,
		path:  path,
		from:  from,
		codec: codec,
		cause: cause,
	}
}

// Error returns a human-readable failure description.
func (e *Error) Error() string {
	context := "patch"
	if e.index >= 0 {
		context = fmt.Sprintf("operation %d", e.index)
	}
	if e.op != "" {
		context = fmt.Sprintf("%s (%s", context, e.op)
		if e.path != "" {
			context = fmt.Sprintf("%s %q", context, e.path)
		}
		context += ")"
	}
	if e.codec != "" {
		context = fmt.Sprintf("%s [%s codec]", context, e.codec)
	}
	if e.cause == nil {
		return fmt.Sprintf("%s: %v", context, e.kind)
	}
	return fmt.Sprintf("%s: %v: %v", context, e.kind, e.cause)
}

// Unwrap exposes both the stable kind and original cause to errors.Is/As.
func (e *Error) Unwrap() error {
	if e.cause == nil {
		return e.kind
	}
	return errors.Join(e.kind, e.cause)
}

// Kind returns the stable failure class.
func (e *Error) Kind() error {
	return e.kind
}

// Index returns the failing operation index, or -1 when no operation was decoded.
func (e *Error) Index() int {
	return e.index
}

// Op returns the failing operation name.
func (e *Error) Op() string {
	return e.op
}

// Path returns the failing operation target path.
func (e *Error) Path() string {
	return e.path
}

// From returns the failing operation source path when present.
func (e *Error) From() string {
	return e.from
}

// Codec returns the codec boundary associated with the failure.
func (e *Error) Codec() string {
	return e.codec
}

// Cause returns the original wrapped error.
func (e *Error) Cause() error {
	return e.cause
}
