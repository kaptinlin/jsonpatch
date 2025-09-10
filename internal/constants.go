// Package internal provides internal types and constants for JSON Patch operations.
package internal

// OpType represents the string type for JSON Patch operation names,
// such as "add", "remove", "replace", etc.
// Used for type safety and constant references only.
type OpType string

const (
	// OpAddType represents the "add" operation type for JSON Patch (RFC 6902)
	OpAddType OpType = "add"
	// OpRemoveType represents the "remove" operation type for JSON Patch (RFC 6902)
	OpRemoveType OpType = "remove"
	// OpReplaceType represents the "replace" operation type for JSON Patch (RFC 6902)
	OpReplaceType OpType = "replace"
	// OpMoveType represents the "move" operation type for JSON Patch (RFC 6902)
	OpMoveType OpType = "move"
	// OpCopyType represents the "copy" operation type for JSON Patch (RFC 6902)
	OpCopyType OpType = "copy"
	// OpTestType represents the "test" operation type for JSON Patch (RFC 6902)
	OpTestType OpType = "test"

	// OpContainsType represents the "contains" operation type for JSON Predicate operations
	OpContainsType OpType = "contains"
	// OpDefinedType represents the "defined" operation type for JSON Predicate operations
	OpDefinedType       OpType = "defined"
	OpUndefinedType     OpType = "undefined"
	OpTypeType          OpType = "type"
	OpTestTypeType      OpType = "test_type"
	OpTestStringType    OpType = "test_string"
	OpTestStringLenType OpType = "test_string_len"
	OpEndsType          OpType = "ends"
	OpStartsType        OpType = "starts"
	OpInType            OpType = "in"
	OpLessType          OpType = "less"
	OpMoreType          OpType = "more"
	OpMatchesType       OpType = "matches"

	// OpAndType represents the "and" operation type for composite operations
	OpAndType OpType = "and"
	OpOrType  OpType = "or"
	OpNotType OpType = "not"

	// OpFlipType represents the "flip" operation type for extended operations
	OpFlipType   OpType = "flip"
	OpIncType    OpType = "inc"
	OpStrInsType OpType = "str_ins"
	OpStrDelType OpType = "str_del"
	OpSplitType  OpType = "split"
	OpMergeType  OpType = "merge"
	OpExtendType OpType = "extend"
)

// Operation code constants, numeric codes
const (
	// JSON Patch (RFC 6902) operations
	OpAddCode     = 0
	OpRemoveCode  = 1
	OpReplaceCode = 2
	OpCopyCode    = 3
	OpMoveCode    = 4
	OpTestCode    = 5

	// String editing
	OpStrInsCode = 6
	OpStrDelCode = 7

	// Extra
	OpFlipCode = 8
	OpIncCode  = 9

	// Slate.js
	OpSplitCode  = 10
	OpMergeCode  = 11
	OpExtendCode = 12

	// JSON Predicate
	OpContainsCode      = 30
	OpDefinedCode       = 31
	OpEndsCode          = 32
	OpInCode            = 33
	OpLessCode          = 34
	OpMatchesCode       = 35
	OpMoreCode          = 36
	OpStartsCode        = 37
	OpUndefinedCode     = 38
	OpTestTypeCode      = 39
	OpTestStringCode    = 40
	OpTestStringLenCode = 41
	OpTypeCode          = 42
	OpAndCode           = 43
	OpNotCode           = 44
	OpOrCode            = 45
)
