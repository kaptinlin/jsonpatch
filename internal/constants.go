package internal

// OpType represents the string type for JSON Patch operation names,
// such as "add", "remove", "replace", etc.
// Used for type safety and constant references only.

type OpType string

const (
	// JSON Patch (RFC 6902) operations
	OpAddType     OpType = "add"
	OpRemoveType  OpType = "remove"
	OpReplaceType OpType = "replace"
	OpMoveType    OpType = "move"
	OpCopyType    OpType = "copy"
	OpTestType    OpType = "test"

	// JSON Predicate operations
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

	// Composite operations
	OpAndType OpType = "and"
	OpOrType  OpType = "or"
	OpNotType OpType = "not"

	// Extended operations
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
