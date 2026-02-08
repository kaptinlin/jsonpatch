package internal

// Operation represents a JSON Patch operation object
// Matches json-joy Operation interface exactly
type Operation struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value any    `json:"value,omitempty"`
	From  string `json:"from,omitempty"`

	// Extended operation fields
	Inc float64 `json:"inc"` // No omitempty - 0 is a valid increment
	Pos int     `json:"pos"` // No omitempty - 0 is a valid position
	Str string  `json:"str"`
	Len int     `json:"len"` // No omitempty - 0 is a valid length

	// Predicate fields
	Not        bool        `json:"not,omitempty"`
	Type       any         `json:"type,omitempty"`
	IgnoreCase bool        `json:"ignore_case,omitempty"`
	Apply      []Operation `json:"apply,omitempty"`

	// Special fields
	Props      map[string]any `json:"props,omitempty"`
	DeleteNull bool           `json:"deleteNull,omitempty"`
	OldValue   any            `json:"oldValue,omitempty"`
}

// CompactOperation represents a compact format operation as an array with clearer semantics.
type CompactOperation = []any

// OpResult represents the result of a single operation with generic type support.
type OpResult[T Document] struct {
	Doc T   `json:"doc"`
	Old any `json:"old,omitempty"`
}

// PatchResult represents the result of applying a JSON Patch with generic type support.
// It contains the patched document and the results of individual operations.
type PatchResult[T Document] struct {
	Doc T             // The patched document of the original type
	Res []OpResult[T] // Results of individual patch operations
}

// Document defines the supported document types for JSON Patch operations.
// Supports: structs, map[string]any, []byte (JSON), and string (JSON).
type Document interface {
	~[]byte | ~string | map[string]any | any
}

// Options holds configuration parameters for patch operations.
// This is the unified configuration struct following Go best practices.
type Options struct {
	Mutate        bool               // Whether to modify the original document
	CreateMatcher CreateRegexMatcher // Optional regex matcher creator
}

// Option represents functional options for configuring patch operations.
// This follows the standard Go functional options pattern.
type Option func(*Options)

// WithMutate configures whether the patch operation should modify the original document.
// When false (default), returns a new copy. When true, modifies the original.
func WithMutate(mutate bool) Option {
	return func(opts *Options) {
		opts.Mutate = mutate
	}
}

// WithMatcher configures a custom regex matcher for pattern operations.
// The createMatcher function should create a RegexMatcher from a pattern and ignoreCase flag.
func WithMatcher(createMatcher CreateRegexMatcher) Option {
	return func(opts *Options) {
		opts.CreateMatcher = createMatcher
	}
}

// JSONPatchTypes represents the valid JSON types for type operations.
type JSONPatchTypes string

const (
	// JSONPatchTypeString represents the string data type for JSON Patch type operations
	JSONPatchTypeString JSONPatchTypes = "string"
	// JSONPatchTypeNumber represents the number data type for JSON Patch type operations
	JSONPatchTypeNumber JSONPatchTypes = "number"
	// JSONPatchTypeBoolean represents the boolean data type for JSON Patch type operations
	JSONPatchTypeBoolean JSONPatchTypes = "boolean"
	// JSONPatchTypeObject represents the object data type for JSON Patch type operations
	JSONPatchTypeObject JSONPatchTypes = "object"
	// JSONPatchTypeInteger represents the integer data type for JSON Patch type operations
	JSONPatchTypeInteger JSONPatchTypes = "integer"
	// JSONPatchTypeArray represents the array data type for JSON Patch type operations
	JSONPatchTypeArray JSONPatchTypes = "array"
	// JSONPatchTypeNull represents the null data type for JSON Patch type operations
	JSONPatchTypeNull JSONPatchTypes = "null"
)

// RegexMatcher is a function type that tests if a value matches a pattern.
// This aligns with json-joy's RegexMatcher type.
type RegexMatcher func(value string) bool

// CreateRegexMatcher is a function type that creates a RegexMatcher from a pattern.
// This aligns with json-joy's CreateRegexMatcher type.
// The function should return a matcher that always returns false if the pattern is invalid.
type CreateRegexMatcher func(pattern string, ignoreCase bool) RegexMatcher

// JSONPatchOptions contains options for JSON Patch decoding operations.
// This aligns with json-joy's JsonPatchOptions interface.
type JSONPatchOptions struct {
	CreateMatcher CreateRegexMatcher
}

// IsValidJSONPatchType checks if a type string is a valid JSON Patch type
func IsValidJSONPatchType(typeStr string) bool {
	switch JSONPatchTypes(typeStr) {
	case JSONPatchTypeString, JSONPatchTypeNumber, JSONPatchTypeBoolean,
		JSONPatchTypeObject, JSONPatchTypeInteger, JSONPatchTypeArray,
		JSONPatchTypeNull:
		return true
	default:
		return false
	}
}

// GetJSONPatchType returns the JSON Patch type for a given value.
func GetJSONPatchType(value any) JSONPatchTypes {
	if value == nil {
		return JSONPatchTypeNull
	}

	switch v := value.(type) {
	case string:
		return JSONPatchTypeString
	case bool:
		return JSONPatchTypeBoolean
	case []any, []string, []int, []float64:
		return JSONPatchTypeArray
	case map[string]any:
		return JSONPatchTypeObject
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64:
		return JSONPatchTypeInteger
	case float32:
		if v == float32(int32(v)) {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber
	case float64:
		if v == float64(int64(v)) {
			return JSONPatchTypeInteger
		}
		return JSONPatchTypeNumber
	default:
		return JSONPatchTypeNull
	}
}

// Operation type checking functions provide efficient type detection

// IsJSONPatchOperation checks if operation is a core JSON Patch operation
func IsJSONPatchOperation(op string) bool {
	switch op {
	case string(OpAddType), string(OpRemoveType), string(OpReplaceType),
		string(OpMoveType), string(OpCopyType), string(OpTestType):
		return true
	default:
		return false
	}
}

// IsPredicateOperation checks if operation is a predicate operation
func IsPredicateOperation(op string) bool {
	return IsFirstOrderPredicateOperation(op) || IsSecondOrderPredicateOperation(op)
}

// IsFirstOrderPredicateOperation checks if operation is a first-order predicate
func IsFirstOrderPredicateOperation(op string) bool {
	switch op {
	case string(OpTestType), string(OpDefinedType), string(OpUndefinedType),
		string(OpTestTypeType), string(OpTestStringType), string(OpTestStringLenType),
		string(OpContainsType), string(OpEndsType), string(OpStartsType),
		string(OpInType), string(OpLessType), string(OpMoreType),
		string(OpMatchesType):
		return true
	default:
		return false
	}
}

// IsSecondOrderPredicateOperation checks if operation is a second-order predicate.
func IsSecondOrderPredicateOperation(op string) bool {
	switch op {
	case string(OpAndType), string(OpOrType), string(OpNotType):
		return true
	default:
		return false
	}
}

// IsJSONPatchExtendedOperation checks if operation is an extended operation.
func IsJSONPatchExtendedOperation(op string) bool {
	switch op {
	case string(OpStrInsType), string(OpStrDelType), string(OpFlipType),
		string(OpIncType), string(OpSplitType), string(OpMergeType),
		string(OpExtendType):
		return true
	default:
		return false
	}
}
