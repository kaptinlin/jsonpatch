package internal

// Operation represents a JSON Patch operation object
// compatible with map format
type Operation = map[string]interface{}

// CompactOperation represents a compact format operation
// actually []interface{}, but with clearer semantics
type CompactOperation = []interface{}

// OpResult represents the result of a single operation
type OpResult struct {
	Doc interface{} `json:"doc"`
	Old interface{} `json:"old,omitempty"`
}

// PatchResult represents the result of a patch operation
type PatchResult struct {
	Doc interface{} `json:"doc"`
	Res []OpResult  `json:"res"`
}

// JsonPatchTypes represents the valid JSON types for type operations.
type JsonPatchTypes string

const (
	JsonPatchTypeString  JsonPatchTypes = "string"
	JsonPatchTypeNumber  JsonPatchTypes = "number"
	JsonPatchTypeBoolean JsonPatchTypes = "boolean"
	JsonPatchTypeObject  JsonPatchTypes = "object"
	JsonPatchTypeInteger JsonPatchTypes = "integer"
	JsonPatchTypeArray   JsonPatchTypes = "array"
	JsonPatchTypeNull    JsonPatchTypes = "null"
)

// CreateRegexMatcher is a function type that creates a regular expression matcher.
type CreateRegexMatcher func(pattern string, ignoreCase bool) RegexMatcher

// RegexMatcher is a function type that tests if a value matches a pattern.
type RegexMatcher func(value string) bool

// JsonPatchOptions contains options for JSON Patch operations.
type JsonPatchOptions struct {
	CreateMatcher CreateRegexMatcher
}

// ApplyPatchOptions contains options for applying patches.
type ApplyPatchOptions struct {
	Mutate           bool
	JsonPatchOptions JsonPatchOptions
}

// IsValidJsonPatchType checks if a type string is a valid JSON Patch type
func IsValidJsonPatchType(typeStr string) bool {
	switch JsonPatchTypes(typeStr) {
	case JsonPatchTypeString, JsonPatchTypeNumber, JsonPatchTypeBoolean,
		JsonPatchTypeObject, JsonPatchTypeInteger, JsonPatchTypeArray,
		JsonPatchTypeNull:
		return true
	default:
		return false
	}
}

// GetJsonPatchType returns the JSON Patch type for a given value
func GetJsonPatchType(value interface{}) JsonPatchTypes {
	if value == nil {
		return JsonPatchTypeNull
	}

	switch v := value.(type) {
	case string:
		return JsonPatchTypeString
	case bool:
		return JsonPatchTypeBoolean
	case []interface{}, []string, []int, []float64:
		return JsonPatchTypeArray
	case map[string]interface{}:
		return JsonPatchTypeObject
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return JsonPatchTypeInteger
	case float32, float64:
		// Check if it's actually an integer
		switch f := v.(type) {
		case float32:
			if f == float32(int32(f)) {
				return JsonPatchTypeInteger
			}
		case float64:
			if f == float64(int64(f)) {
				return JsonPatchTypeInteger
			}
		}
		return JsonPatchTypeNumber
	default:
		// For other types, use simple checks
		if isArrayType(value) {
			return JsonPatchTypeArray
		}
		if isObjectType(value) {
			return JsonPatchTypeObject
		}
		return JsonPatchTypeNull
	}
}

func isArrayType(value interface{}) bool {
	switch value.(type) {
	case []interface{}, []string, []int, []float64:
		return true
	default:
		return false
	}
}

func isObjectType(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

// Operation type checking functions provide efficient type detection

// IsJsonPatchOperation checks if operation is a core JSON Patch operation
func IsJsonPatchOperation(op string) bool {
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

// IsJsonPatchExtendedOperation checks if operation is an extended operation.
func IsJsonPatchExtendedOperation(op string) bool {
	switch op {
	case string(OpStrInsType), string(OpStrDelType), string(OpFlipType),
		string(OpIncType), string(OpSplitType), string(OpMergeType),
		string(OpExtendType):
		return true
	default:
		return false
	}
}
