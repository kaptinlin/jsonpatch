// Package compact implements a compact array-based codec for JSON Patch operations.
// This codec uses arrays instead of objects to represent operations, significantly
// reducing the physical space required for encoding while maintaining readability.
package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// OpCode represents operation codes for compact format.
type OpCode int

// Operation codes for compact format - derived from internal constants.
const (
	// JSON Patch (RFC 6902) operations
	OpCodeAdd     OpCode = OpCode(internal.OpAddCode)
	OpCodeRemove  OpCode = OpCode(internal.OpRemoveCode)
	OpCodeReplace OpCode = OpCode(internal.OpReplaceCode)
	OpCodeCopy    OpCode = OpCode(internal.OpCopyCode)
	OpCodeMove    OpCode = OpCode(internal.OpMoveCode)
	OpCodeTest    OpCode = OpCode(internal.OpTestCode)

	// String editing
	OpCodeStrIns OpCode = OpCode(internal.OpStrInsCode)
	OpCodeStrDel OpCode = OpCode(internal.OpStrDelCode)

	// Extra
	OpCodeFlip OpCode = OpCode(internal.OpFlipCode)
	OpCodeInc  OpCode = OpCode(internal.OpIncCode)

	// Slate.js
	OpCodeSplit  OpCode = OpCode(internal.OpSplitCode)
	OpCodeMerge  OpCode = OpCode(internal.OpMergeCode)
	OpCodeExtend OpCode = OpCode(internal.OpExtendCode)

	// JSON Predicate operations
	OpCodeContains      OpCode = OpCode(internal.OpContainsCode)
	OpCodeDefined       OpCode = OpCode(internal.OpDefinedCode)
	OpCodeEnds          OpCode = OpCode(internal.OpEndsCode)
	OpCodeIn            OpCode = OpCode(internal.OpInCode)
	OpCodeLess          OpCode = OpCode(internal.OpLessCode)
	OpCodeMatches       OpCode = OpCode(internal.OpMatchesCode)
	OpCodeMore          OpCode = OpCode(internal.OpMoreCode)
	OpCodeStarts        OpCode = OpCode(internal.OpStartsCode)
	OpCodeUndefined     OpCode = OpCode(internal.OpUndefinedCode)
	OpCodeTestType      OpCode = OpCode(internal.OpTestTypeCode)
	OpCodeTestString    OpCode = OpCode(internal.OpTestStringCode)
	OpCodeTestStringLen OpCode = OpCode(internal.OpTestStringLenCode)
	OpCodeType          OpCode = OpCode(internal.OpTypeCode)
	OpCodeAnd           OpCode = OpCode(internal.OpAndCode)
	OpCodeNot           OpCode = OpCode(internal.OpNotCode)
	OpCodeOr            OpCode = OpCode(internal.OpOrCode)
)

// Op represents a compact format operation as an array.
type Op []any

// EncoderOptions configures the compact encoder behavior.
type EncoderOptions struct {
	// StringOpcode uses string opcodes instead of numeric ones.
	StringOpcode bool
}

// EncoderOption is a functional option for configuring the encoder.
type EncoderOption func(*EncoderOptions)

// WithStringOpcode configures the encoder to use string opcodes.
func WithStringOpcode(useString bool) EncoderOption {
	return func(opts *EncoderOptions) {
		opts.StringOpcode = useString
	}
}

// Operation represents a compact format operation.
type Operation = internal.CompactOperation
