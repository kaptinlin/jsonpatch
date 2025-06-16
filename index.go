package jsonpatch

import "github.com/kaptinlin/jsonpatch/internal"

// Re-export all types from the internal package for convenience.

// Core types
type (
	Operation = internal.Operation
	OpType    = internal.OpType
	Op        = internal.Op
)

// Generic types for type-safe JSON Patch operations
type (
	Document = internal.Document
	Option   = internal.Option
	Options  = internal.Options
)

var WithMutate = internal.WithMutate
var WithMatcher = internal.WithMatcher

// DefaultOptions returns the default configuration for patch operations.
func DefaultOptions() *internal.Options {
	return &internal.Options{
		Mutate:        false, // Immutable by default for safety
		CreateMatcher: nil,   // Use default regex implementation
	}
}

// Generic result types (requires Go 1.18+)
type OpResult[T Document] = internal.OpResult[T]
type PatchResult[T Document] = internal.PatchResult[T]

// Operation type constants (string constants)
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

// JSON Patch types
type (
	JsonPatchTypes   = internal.JsonPatchTypes
	RegexMatcher     = internal.RegexMatcher
	JsonPatchOptions = internal.JsonPatchOptions
)

const (
	JsonPatchTypeString  = internal.JsonPatchTypeString
	JsonPatchTypeNumber  = internal.JsonPatchTypeNumber
	JsonPatchTypeBoolean = internal.JsonPatchTypeBoolean
	JsonPatchTypeObject  = internal.JsonPatchTypeObject
	JsonPatchTypeInteger = internal.JsonPatchTypeInteger
	JsonPatchTypeArray   = internal.JsonPatchTypeArray
	JsonPatchTypeNull    = internal.JsonPatchTypeNull
)

// Re-export functions
var (
	IsValidJsonPatchType = internal.IsValidJsonPatchType
	GetJsonPatchType     = internal.GetJsonPatchType

	// Operation type checking functions
	IsJsonPatchOperation            = internal.IsJsonPatchOperation
	IsPredicateOperation            = internal.IsPredicateOperation
	IsFirstOrderPredicateOperation  = internal.IsFirstOrderPredicateOperation
	IsSecondOrderPredicateOperation = internal.IsSecondOrderPredicateOperation
	IsJsonPatchExtendedOperation    = internal.IsJsonPatchExtendedOperation
)
