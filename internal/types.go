package internal

// Document defines the supported document types for JSON Patch operations.
type Document interface {
	~[]byte | ~string | map[string]any | any
}

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

// RegexMatcher tests whether a string value matches a pattern.
type RegexMatcher func(value string) bool

// CreateRegexMatcher creates a RegexMatcher from a pattern
// and case-sensitivity flag.
type CreateRegexMatcher func(pattern string, ignoreCase bool) RegexMatcher

// Options holds configuration for patch operations.
type Options struct {
	Mutate        bool               // Whether to modify the original document.
	CreateMatcher CreateRegexMatcher // Optional regex matcher factory.
}

// Option is a functional option for configuring patch operations.
type Option func(*Options)
