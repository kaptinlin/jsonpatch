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

// DefaultOptions returns the default patch options.
func DefaultOptions() *internal.Options {
	return &internal.Options{}
}

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

// Operation code constants (numeric constants)
const (
	// JSON Patch (RFC 6902) operations
	OpAddCode     = internal.OpAddCode
	OpRemoveCode  = internal.OpRemoveCode
	OpReplaceCode = internal.OpReplaceCode
	OpCopyCode    = internal.OpCopyCode
	OpMoveCode    = internal.OpMoveCode
	OpTestCode    = internal.OpTestCode

	// String editing
	OpStrInsCode = internal.OpStrInsCode
	OpStrDelCode = internal.OpStrDelCode

	// Extra
	OpFlipCode = internal.OpFlipCode
	OpIncCode  = internal.OpIncCode

	// Slate.js
	OpSplitCode  = internal.OpSplitCode
	OpMergeCode  = internal.OpMergeCode
	OpExtendCode = internal.OpExtendCode

	// JSON Predicate
	OpContainsCode      = internal.OpContainsCode
	OpDefinedCode       = internal.OpDefinedCode
	OpEndsCode          = internal.OpEndsCode
	OpInCode            = internal.OpInCode
	OpLessCode          = internal.OpLessCode
	OpMatchesCode       = internal.OpMatchesCode
	OpMoreCode          = internal.OpMoreCode
	OpStartsCode        = internal.OpStartsCode
	OpUndefinedCode     = internal.OpUndefinedCode
	OpTestTypeCode      = internal.OpTestTypeCode
	OpTestStringCode    = internal.OpTestStringCode
	OpTestStringLenCode = internal.OpTestStringLenCode
	OpTypeCode          = internal.OpTypeCode
	OpAndCode           = internal.OpAndCode
	OpNotCode           = internal.OpNotCode
	OpOrCode            = internal.OpOrCode
)

// Types represents the valid JSON types for type operations.
type Types = internal.JSONPatchType

// RegexMatcher is a function type that tests if a value matches a pattern.
// This aligns with json-joy's RegexMatcher type.
type RegexMatcher = internal.RegexMatcher

// CreateRegexMatcher is a function type that creates a RegexMatcher from a pattern.
// This aligns with json-joy's CreateRegexMatcher type.
type CreateRegexMatcher = internal.CreateRegexMatcher

// JSONPatchTypeString represents the JSON string type.
const JSONPatchTypeString = internal.JSONPatchTypeString

// JSONPatchTypeNumber represents the JSON number type.
const JSONPatchTypeNumber = internal.JSONPatchTypeNumber

// JSONPatchTypeBoolean represents the JSON boolean type.
const JSONPatchTypeBoolean = internal.JSONPatchTypeBoolean

// JSONPatchTypeObject represents the JSON object type.
const JSONPatchTypeObject = internal.JSONPatchTypeObject

// JSONPatchTypeInteger represents the JSON integer type.
const JSONPatchTypeInteger = internal.JSONPatchTypeInteger

// JSONPatchTypeArray represents the JSON array type.
const JSONPatchTypeArray = internal.JSONPatchTypeArray

// JSONPatchTypeNull represents the JSON null type.
const JSONPatchTypeNull = internal.JSONPatchTypeNull

// Re-export functions
var (
	IsValidJSONPatchType = internal.IsValidJSONPatchType
	GetJSONPatchType     = internal.GetJSONPatchType

	// Operation type checking functions
	IsJSONPatchOperation            = internal.IsJSONPatchOperation
	IsPredicateOperation            = internal.IsPredicateOperation
	IsFirstOrderPredicateOperation  = internal.IsFirstOrderPredicateOperation
	IsSecondOrderPredicateOperation = internal.IsSecondOrderPredicateOperation
	IsJSONPatchExtendedOperation    = internal.IsJSONPatchExtendedOperation
)
