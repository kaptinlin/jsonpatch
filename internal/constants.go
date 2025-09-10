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
	OpDefinedType OpType = "defined"
	// OpUndefinedType represents the "undefined" operation type for JSON Predicate operations
	OpUndefinedType OpType = "undefined"
	// OpTypeType represents the "type" operation type for JSON Predicate operations
	OpTypeType OpType = "type"
	// OpTestTypeType represents the "test_type" operation type for JSON Predicate operations
	OpTestTypeType OpType = "test_type"
	// OpTestStringType represents the "test_string" operation type for JSON Predicate operations
	OpTestStringType OpType = "test_string"
	// OpTestStringLenType represents the "test_string_len" operation type for JSON Predicate operations
	OpTestStringLenType OpType = "test_string_len"
	// OpEndsType represents the "ends" operation type for JSON Predicate operations
	OpEndsType OpType = "ends"
	// OpStartsType represents the "starts" operation type for JSON Predicate operations
	OpStartsType OpType = "starts"
	// OpInType represents the "in" operation type for JSON Predicate operations
	OpInType OpType = "in"
	// OpLessType represents the "less" operation type for JSON Predicate operations
	OpLessType OpType = "less"
	// OpMoreType represents the "more" operation type for JSON Predicate operations
	OpMoreType OpType = "more"
	// OpMatchesType represents the "matches" operation type for JSON Predicate operations
	OpMatchesType OpType = "matches"

	// OpAndType represents the "and" operation type for composite operations
	OpAndType OpType = "and"
	// OpOrType represents the "or" operation type for composite operations
	OpOrType OpType = "or"
	// OpNotType represents the "not" operation type for composite operations
	OpNotType OpType = "not"

	// OpFlipType represents the "flip" operation type for extended operations
	OpFlipType OpType = "flip"
	// OpIncType represents the "inc" operation type for extended operations
	OpIncType OpType = "inc"
	// OpStrInsType represents the "str_ins" operation type for extended operations
	OpStrInsType OpType = "str_ins"
	// OpStrDelType represents the "str_del" operation type for extended operations
	OpStrDelType OpType = "str_del"
	// OpSplitType represents the "split" operation type for extended operations
	OpSplitType OpType = "split"
	// OpMergeType represents the "merge" operation type for extended operations
	OpMergeType OpType = "merge"
	// OpExtendType represents the "extend" operation type for extended operations
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
