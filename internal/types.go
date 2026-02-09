package internal

import (
	"math"
	"reflect"
)

// Operation represents a JSON Patch operation object.
type Operation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
	From  string `json:"from,omitempty"`

	// Extended operation fields.
	Inc float64 `json:"inc"` // No omitempty — 0 is a valid increment.
	Pos int     `json:"pos"` // No omitempty — 0 is a valid position.
	Str string  `json:"str"`
	Len int     `json:"len"` // No omitempty — 0 is a valid length.

	// Predicate fields.
	Not        bool        `json:"not,omitempty"`
	Type       any         `json:"type,omitempty"`
	IgnoreCase bool        `json:"ignore_case,omitempty"`
	Apply      []Operation `json:"apply,omitempty"`

	// Special fields.
	Props      map[string]any `json:"props,omitempty"`
	DeleteNull bool           `json:"deleteNull,omitempty"`
	OldValue   any            `json:"oldValue,omitempty"`
}

// CompactOperation represents a compact-format operation as an array.
type CompactOperation = []any

// OpResult holds the result of a single operation application.
type OpResult[T Document] struct {
	Doc T   `json:"doc"`
	Old any `json:"old,omitempty"`
}

// PatchResult holds the result of applying an entire JSON Patch.
type PatchResult[T Document] struct {
	Doc T             // The patched document.
	Res []OpResult[T] // Results of individual operations.
}

// Document defines the supported document types for JSON Patch operations.
type Document interface {
	~[]byte | ~string | map[string]any | any
}

// Options holds configuration for patch operations.
type Options struct {
	Mutate        bool               // Whether to modify the original document.
	CreateMatcher CreateRegexMatcher // Optional regex matcher factory.
}

// Option is a functional option for configuring patch operations.
type Option func(*Options)

// WithMutate sets whether the patch should modify the original document.
func WithMutate(mutate bool) Option {
	return func(o *Options) {
		o.Mutate = mutate
	}
}

// WithMatcher sets a custom regex matcher factory for pattern operations.
func WithMatcher(createMatcher CreateRegexMatcher) Option {
	return func(o *Options) {
		o.CreateMatcher = createMatcher
	}
}

// JSONPatchOptions contains options for JSON Patch decoding.
type JSONPatchOptions struct {
	CreateMatcher CreateRegexMatcher
}

// JSONPatchType represents valid JSON types for type-checking operations.
type JSONPatchType string

// Valid JSON Patch types.
const (
	JSONPatchTypeString  JSONPatchType = "string"
	JSONPatchTypeNumber  JSONPatchType = "number"
	JSONPatchTypeBoolean JSONPatchType = "boolean"
	JSONPatchTypeObject  JSONPatchType = "object"
	JSONPatchTypeInteger JSONPatchType = "integer"
	JSONPatchTypeArray   JSONPatchType = "array"
	JSONPatchTypeNull    JSONPatchType = "null"
)

// RegexMatcher tests whether a string value matches a pattern.
type RegexMatcher func(value string) bool

// CreateRegexMatcher creates a RegexMatcher from a pattern and case-sensitivity flag.
type CreateRegexMatcher func(pattern string, ignoreCase bool) RegexMatcher

// IsValidJSONPatchType reports whether typeStr is a valid JSON Patch type name.
func IsValidJSONPatchType(typeStr string) bool {
	switch JSONPatchType(typeStr) {
	case JSONPatchTypeString, JSONPatchTypeNumber, JSONPatchTypeBoolean,
		JSONPatchTypeObject, JSONPatchTypeInteger, JSONPatchTypeArray,
		JSONPatchTypeNull:
		return true
	default:
		return false
	}
}

// GetJSONPatchType returns the JSON Patch type for a Go value.
// Whole-number floats return "integer"; unknown types return "null".
func GetJSONPatchType(value any) JSONPatchType {
	if value == nil {
		return JSONPatchTypeNull
	}

	switch v := value.(type) {
	case string:
		return JSONPatchTypeString

	case bool:
		return JSONPatchTypeBoolean

	case map[string]any:
		return JSONPatchTypeObject

	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return JSONPatchTypeInteger

	case float32:
		if !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) && math.Trunc(float64(v)) == float64(v) {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber

	case float64:
		if !math.IsNaN(v) && !math.IsInf(v, 0) && math.Trunc(v) == v {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber

	default:
		if isSliceKind(value) {
			return JSONPatchTypeArray
		}
		return JSONPatchTypeNull
	}
}

// isSliceKind reports whether value is a slice or array using reflection.
func isSliceKind(value any) bool {
	kind := reflect.TypeOf(value).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

// IsJSONPatchOperation reports whether op is a core JSON Patch (RFC 6902) operation.
func IsJSONPatchOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only RFC 6902 operations
	case OpAddType, OpRemoveType, OpReplaceType,
		OpMoveType, OpCopyType, OpTestType:
		return true
	default:
		return false
	}
}

// IsPredicateOperation reports whether op is any predicate operation.
func IsPredicateOperation(op string) bool {
	return IsFirstOrderPredicateOperation(op) || IsSecondOrderPredicateOperation(op)
}

// IsFirstOrderPredicateOperation reports whether op is a first-order predicate.
func IsFirstOrderPredicateOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only first-order predicates
	case OpTestType, OpDefinedType, OpUndefinedType,
		OpTestTypeType, OpTestStringType, OpTestStringLenType,
		OpContainsType, OpEndsType, OpStartsType,
		OpInType, OpLessType, OpMoreType, OpMatchesType:
		return true
	default:
		return false
	}
}

// IsSecondOrderPredicateOperation reports whether op is a second-order (composite) predicate.
func IsSecondOrderPredicateOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only second-order predicates
	case OpAndType, OpOrType, OpNotType:
		return true
	default:
		return false
	}
}

// IsJSONPatchExtendedOperation reports whether op is an extended operation.
func IsJSONPatchExtendedOperation(op string) bool {
	switch OpType(op) { //nolint:exhaustive // intentionally matches only extended operations
	case OpStrInsType, OpStrDelType, OpFlipType,
		OpIncType, OpSplitType, OpMergeType, OpExtendType:
		return true
	default:
		return false
	}
}
