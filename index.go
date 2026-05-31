package jsonpatch

import "github.com/kaptinlin/jsonpatch/internal"

// Operation aliases internal.Operation.
type Operation = internal.Operation

// OpType names a JSON Patch operation.
type OpType = internal.OpType

// Op applies itself to a document.
type Op = internal.Op

// Document is the set of document types supported by the generic API.
type Document = internal.Document

// Option applies one patch option.
type Option = internal.Option

// Options configures patch application.
type Options = internal.Options

// WithMutate enables or disables in-place mutation.
var WithMutate = internal.WithMutate

// WithMatcher sets the regex matcher factory used by pattern operations.
var WithMatcher = internal.WithMatcher

// OpResult is the result of applying one operation.
type OpResult[T Document] = internal.OpResult[T]

// PatchResult is the result of applying a sequence of operations.
type PatchResult[T Document] = internal.PatchResult[T]

// These constants name the supported operations.
const (
	// JSON Patch (RFC 6902) operations
	OpAddType     = internal.OpAddType
	OpRemoveType  = internal.OpRemoveType
	OpReplaceType = internal.OpReplaceType
	OpMoveType    = internal.OpMoveType
	OpCopyType    = internal.OpCopyType
	OpTestType    = internal.OpTestType

	// JSON Predicate operations
	OpContainsType      = internal.OpContainsType
	OpDefinedType       = internal.OpDefinedType
	OpUndefinedType     = internal.OpUndefinedType
	OpTypeType          = internal.OpTypeType
	OpTestTypeType      = internal.OpTestTypeType
	OpTestStringType    = internal.OpTestStringType
	OpTestStringLenType = internal.OpTestStringLenType
	OpEndsType          = internal.OpEndsType
	OpStartsType        = internal.OpStartsType
	OpInType            = internal.OpInType
	OpLessType          = internal.OpLessType
	OpMoreType          = internal.OpMoreType
	OpMatchesType       = internal.OpMatchesType

	// Composite operations
	OpAndType = internal.OpAndType
	OpOrType  = internal.OpOrType
	OpNotType = internal.OpNotType

	// Extended operations
	OpFlipType   = internal.OpFlipType
	OpIncType    = internal.OpIncType
	OpStrInsType = internal.OpStrInsType
	OpStrDelType = internal.OpStrDelType
	OpSplitType  = internal.OpSplitType
	OpMergeType  = internal.OpMergeType
	OpExtendType = internal.OpExtendType
)

// RegexMatcher tests if a value matches a pattern.
type RegexMatcher = internal.RegexMatcher

// CreateRegexMatcher creates a RegexMatcher from a pattern.
type CreateRegexMatcher = internal.CreateRegexMatcher
