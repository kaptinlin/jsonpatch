// Package internal provides shared types, interfaces, and constants for JSON Patch operations.
package internal

// OpType represents a JSON Patch operation name such as "add" or "remove".
type OpType string

// JSON Patch (RFC 6902) operation types.
const (
	OpAddType     OpType = "add"
	OpRemoveType  OpType = "remove"
	OpReplaceType OpType = "replace"
	OpMoveType    OpType = "move"
	OpCopyType    OpType = "copy"
	OpTestType    OpType = "test"
)

// JSON Predicate operation types.
const (
	OpContainsType      OpType = "contains"
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
)

// Composite (second-order) predicate operation types.
const (
	OpAndType OpType = "and"
	OpOrType  OpType = "or"
	OpNotType OpType = "not"
)

// Extended operation types.
const (
	OpFlipType   OpType = "flip"
	OpIncType    OpType = "inc"
	OpStrInsType OpType = "str_ins"
	OpStrDelType OpType = "str_del"
	OpSplitType  OpType = "split"
	OpMergeType  OpType = "merge"
	OpExtendType OpType = "extend"
)

// Operation codes for binary serialization.
const (
	// JSON Patch (RFC 6902).
	OpAddCode     = 0
	OpRemoveCode  = 1
	OpReplaceCode = 2
	OpCopyCode    = 3
	OpMoveCode    = 4
	OpTestCode    = 5

	// String editing.
	OpStrInsCode = 6
	OpStrDelCode = 7

	// Extra operations.
	OpFlipCode = 8
	OpIncCode  = 9

	// Slate.js operations.
	OpSplitCode  = 10
	OpMergeCode  = 11
	OpExtendCode = 12

	// JSON Predicate operations.
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
