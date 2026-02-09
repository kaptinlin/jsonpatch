// Package compact implements a compact array-based codec for JSON Patch operations.
// It uses arrays instead of objects to represent operations, reducing the physical
// space required for encoding while maintaining readability.
package compact

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// Code represents a numeric operation code in compact format.
type Code int

// Numeric operation codes derived from internal constants.
const (
	// JSON Patch (RFC 6902)
	CodeAdd     Code = Code(internal.OpAddCode)
	CodeRemove  Code = Code(internal.OpRemoveCode)
	CodeReplace Code = Code(internal.OpReplaceCode)
	CodeCopy    Code = Code(internal.OpCopyCode)
	CodeMove    Code = Code(internal.OpMoveCode)
	CodeTest    Code = Code(internal.OpTestCode)

	// String editing
	CodeStrIns Code = Code(internal.OpStrInsCode)
	CodeStrDel Code = Code(internal.OpStrDelCode)

	// Extra
	CodeFlip Code = Code(internal.OpFlipCode)
	CodeInc  Code = Code(internal.OpIncCode)

	// Slate.js
	CodeSplit  Code = Code(internal.OpSplitCode)
	CodeMerge  Code = Code(internal.OpMergeCode)
	CodeExtend Code = Code(internal.OpExtendCode)

	// JSON Predicate
	CodeContains      Code = Code(internal.OpContainsCode)
	CodeDefined       Code = Code(internal.OpDefinedCode)
	CodeEnds          Code = Code(internal.OpEndsCode)
	CodeIn            Code = Code(internal.OpInCode)
	CodeLess          Code = Code(internal.OpLessCode)
	CodeMatches       Code = Code(internal.OpMatchesCode)
	CodeMore          Code = Code(internal.OpMoreCode)
	CodeStarts        Code = Code(internal.OpStartsCode)
	CodeUndefined     Code = Code(internal.OpUndefinedCode)
	CodeTestType      Code = Code(internal.OpTestTypeCode)
	CodeTestString    Code = Code(internal.OpTestStringCode)
	CodeTestStringLen Code = Code(internal.OpTestStringLenCode)
	CodeType          Code = Code(internal.OpTypeCode)
	CodeAnd           Code = Code(internal.OpAndCode)
	CodeNot           Code = Code(internal.OpNotCode)
	CodeOr            Code = Code(internal.OpOrCode)
)

// Op represents a compact format operation as an array.
type Op []any

// Options configures the compact encoder behavior.
type Options struct {
	// StringOpcode uses string opcodes instead of numeric ones.
	StringOpcode bool
}

// Option is a functional option for configuring the encoder.
type Option func(*Options)

// WithStringOpcode configures the encoder to use string opcodes.
func WithStringOpcode(useString bool) Option {
	return func(o *Options) {
		o.StringOpcode = useString
	}
}

// Operation represents a compact format operation.
type Operation = internal.CompactOperation
