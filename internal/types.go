package internal

// Document describes the input and output shapes supported by the generic API.
type Document interface {
	~[]byte | ~string | map[string]any | any
}

// Operation describes one JSON Patch, predicate, or extended operation.
type Operation struct {
	// Op names the operation to apply.
	Op string `json:"op"`
	// Path is the target JSON Pointer.
	Path string `json:"path"`
	// Value holds the operation payload when the operation uses one.
	Value any `json:"value,omitempty"`
	// From is the source JSON Pointer for move and copy.
	From string `json:"from,omitempty"`

	// Inc is the numeric delta for inc. It is never omitted because 0 is valid.
	Inc float64 `json:"inc"`
	// Pos is the string or array position for operations that use one. It is never omitted because 0 is valid.
	Pos int `json:"pos"`
	// Str is the string operand for string-editing operations. It is never omitted because the empty string is valid.
	Str string `json:"str"`
	// Len is the requested length for operations that use one. It is never omitted because 0 is valid.
	Len int `json:"len"`

	// Not inverts supported predicate operations.
	Not bool `json:"not,omitempty"`
	// Type holds one JSON type name or a list of JSON type names.
	Type any `json:"type,omitempty"`
	// IgnoreCase makes string matching case-insensitive when supported.
	IgnoreCase bool `json:"ignore_case,omitempty"`
	// Apply contains nested predicate operations for and, or, and not.
	Apply []Operation `json:"apply,omitempty"`

	// Props supplies object properties for extend.
	Props map[string]any `json:"props,omitempty"`
	// DeleteNull removes properties whose incoming value is null during extend.
	DeleteNull bool `json:"deleteNull,omitempty"`
	// OldValue carries the expected previous value for operations that use it.
	OldValue any `json:"oldValue,omitempty"`
}

// CompactOperation is the array representation used by the compact codec.
type CompactOperation = []any

// OpResult holds the result of applying one operation.
type OpResult[T Document] struct {
	// Doc is the document after the operation completes.
	Doc T `json:"doc"`
	// Old is the previous value returned by operations that expose one.
	Old any `json:"old,omitempty"`
}

// PatchResult holds the result of applying a sequence of operations.
type PatchResult[T Document] struct {
	// Doc is the final patched document.
	Doc T
	// Res contains the per-operation results in application order.
	Res []OpResult[T]
}

// RegexMatcher reports whether value matches a compiled pattern.
type RegexMatcher func(value string) bool

// CreateRegexMatcher builds a RegexMatcher from pattern and ignoreCase.
type CreateRegexMatcher func(pattern string, ignoreCase bool) RegexMatcher

// Options configures patch application.
type Options struct {
	// Mutate enables in-place updates instead of cloning the input first.
	Mutate bool
	// CreateMatcher overrides regex compilation for pattern operations.
	CreateMatcher CreateRegexMatcher
}

// Option applies one patch option.
type Option func(*Options)
