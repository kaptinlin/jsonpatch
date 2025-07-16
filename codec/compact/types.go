// Package compact implements a compact array-based codec for JSON Patch operations.
// This codec uses arrays instead of objects to represent operations, significantly
// reducing the physical space required for encoding while maintaining readability.
package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpCode represents operation codes for compact format
type OpCode int

// Operation codes for compact format
const (
	// JSON Patch (RFC 6902) operations - match internal/constants.go
	OpCodeAdd     OpCode = 0
	OpCodeRemove  OpCode = 1
	OpCodeReplace OpCode = 2
	OpCodeCopy    OpCode = 3
	OpCodeMove    OpCode = 4
	OpCodeTest    OpCode = 5

	// String editing
	OpCodeStrIns OpCode = 6
	OpCodeStrDel OpCode = 7

	// Extra
	OpCodeFlip OpCode = 8
	OpCodeInc  OpCode = 9

	// Slate.js
	OpCodeSplit  OpCode = 10
	OpCodeMerge  OpCode = 11
	OpCodeExtend OpCode = 12

	// JSON Predicate operations
	OpCodeContains      OpCode = 30
	OpCodeDefined       OpCode = 31
	OpCodeEnds          OpCode = 32
	OpCodeIn            OpCode = 33
	OpCodeLess          OpCode = 34
	OpCodeMatches       OpCode = 35
	OpCodeMore          OpCode = 36
	OpCodeStarts        OpCode = 37
	OpCodeUndefined     OpCode = 38
	OpCodeTestType      OpCode = 39
	OpCodeTestString    OpCode = 40
	OpCodeTestStringLen OpCode = 41
	OpCodeType          OpCode = 42
	OpCodeAnd           OpCode = 43
	OpCodeNot           OpCode = 44
	OpCodeOr            OpCode = 45
)

// Note: String operation codes are defined in decode.go lookup table for better performance

// CompactOp represents a compact format operation as an array
type CompactOp []interface{}

// EncoderOptions configures the compact encoder behavior
type EncoderOptions struct {
	// StringOpcode determines whether to use string opcodes instead of numeric ones
	StringOpcode bool
}

// DecoderOptions configures the compact decoder behavior
type DecoderOptions struct {
	// Reserved for future options
}

// EncoderOption is a functional option for configuring the encoder
type EncoderOption func(*EncoderOptions)

// DecoderOption is a functional option for configuring the decoder
type DecoderOption func(*DecoderOptions)

// WithStringOpcode configures the encoder to use string opcodes
func WithStringOpcode(useString bool) EncoderOption {
	return func(opts *EncoderOptions) {
		opts.StringOpcode = useString
	}
}

// Default options
var (
	DefaultEncoderOptions = EncoderOptions{
		StringOpcode: false,
	}
	DefaultDecoderOptions = DecoderOptions{}
)

// Operation type aliases for compatibility
type Operation = internal.Operation
type CompactOperation = internal.CompactOperation
type JsonPatchOptions = internal.Options
