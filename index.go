package jsonpatch

import "github.com/kaptinlin/jsonpatch/internal"

// Re-export all types from the internal package for convenience.

// Operation represents a JSON patch operation.
type Operation = internal.Operation

// OpType represents the type of a JSON patch operation.
type OpType = internal.OpType

// Op represents an executable operation.
type Op = internal.Op

// Document defines the supported document types for JSON Patch operations.
type Document = internal.Document

// Option represents functional options for configuring patch operations.
type Option = internal.Option

// Options holds configuration parameters for patch operations.
type Options = internal.Options

// WithMutate configures whether the patch operation should modify the original document.
var WithMutate = internal.WithMutate

// WithMatcher configures a custom regex matcher for pattern operations.
var WithMatcher = internal.WithMatcher

// DefaultOptions returns the default configuration for patch operations.
func DefaultOptions() *internal.Options {
	return &internal.Options{
		Mutate:        false, // Immutable by default for safety
		CreateMatcher: nil,   // Use default regex implementation
	}
}

// OpResult represents the result of a single operation with generic type support.
type OpResult[T Document] = internal.OpResult[T]

// PatchResult represents the result of applying a JSON Patch with generic type support.
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

// Types represents the valid JSON types for type operations.
type Types = internal.JSONPatchTypes

// RegexMatcher is a function type that tests if a value matches a pattern.
type RegexMatcher = internal.RegexMatcher

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
