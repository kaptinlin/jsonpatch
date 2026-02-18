package op

import (
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpointer"
)

// extractString extracts a string from a value, handling string and []byte types.
// Returns the string and true if successful, or empty string and false otherwise.
func extractString(val any) (string, bool) {
	switch v := val.(type) {
	case string:
		return v, true
	case []byte:
		return string(v), true
	default:
		return "", false
	}
}

// extractPredicateOps converts a slice of any to []internal.PredicateOp.
func extractPredicateOps(operations []any) []internal.PredicateOp {
	ops := make([]internal.PredicateOp, 0, len(operations))
	for _, op := range operations {
		if predicateOp, ok := op.(internal.PredicateOp); ok {
			ops = append(ops, predicateOp)
		}
	}
	return ops
}

// pathEquals checks if two paths are equal.
func pathEquals(p1, p2 []string) bool {
	return slices.Equal(p1, p2)
}

// formatPath formats a path slice into a JSON Pointer string.
// Optimized for common cases with pre-allocated buffer and single-pass construction.
func formatPath(path []string) string {
	if len(path) == 0 {
		return ""
	}

	// Fast path for single segment (very common)
	if len(path) == 1 {
		return "/" + path[0]
	}

	var builder strings.Builder
	// More accurate size estimation based on actual segment lengths
	totalLen := len(path) // One '/' per segment
	for _, segment := range path {
		totalLen += len(segment)
	}
	builder.Grow(totalLen)

	for _, segment := range path {
		builder.WriteByte('/')
		builder.WriteString(segment)
	}

	return builder.String()
}

// value retrieves a value from a document using a path.
// Returns the value at the path or an error if the path is not found.
func value(doc any, path []string) (any, error) {
	if len(path) == 0 {
		return doc, nil
	}

	val, err := jsonpointer.Get(doc, path...)
	if err != nil {
		return nil, ErrPathNotFound
	}
	return val, nil
}

// numericValue retrieves a numeric value from the document at the given path.
// Returns the original value, the float64 conversion, and any error.
func numericValue(doc any, path []string) (any, float64, error) {
	val, err := value(doc, path)
	if err != nil {
		return nil, 0, ErrPathNotFound
	}
	numValue, ok := ToFloat64(val)
	if !ok {
		return nil, 0, ErrNotNumber
	}
	return val, numValue, nil
}

// deepEqual performs a deep equality check between two values.
// Optimized to avoid expensive reflect.DeepEqual for common types.
func deepEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	// Fast path: strings (very common)
	if aStr, aIsStr := a.(string); aIsStr {
		bStr, bIsStr := b.(string)
		return bIsStr && aStr == bStr
	}

	// Fast path: booleans
	if aBool, aIsBool := a.(bool); aIsBool {
		bBool, bIsBool := b.(bool)
		return bIsBool && aBool == bBool
	}

	// Fast path: numeric types
	aFloat, aIsNum := toNumericValue(a)
	bFloat, bIsNum := toNumericValue(b)
	if aIsNum != bIsNum {
		return false
	}
	if aIsNum {
		return aFloat == bFloat
	}

	// Fast path: try direct comparison for comparable types
	if reflect.TypeOf(a).Comparable() && reflect.TypeOf(b).Comparable() {
		return a == b
	}

	// Slow path: complex types (maps, slices, structs)
	return reflect.DeepEqual(a, b)
}

// toNumericValue converts a value to float64 if it's an actual numeric type
// (not a string that could be parsed as a number).
// Returns the float64 value and true if successful, 0 and false otherwise.
func toNumericValue(val any) (float64, bool) {
	switch v := val.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	default:
		return 0, false
	}
}

// DeepClone performs a deep clone of a value.
func DeepClone(value any) (any, error) {
	cloned := deepclone.Clone(value)
	return cloned, nil
}

// parseArrayIndex parses a string token as an array index.
func parseArrayIndex(token string) (int, error) {
	index, err := strconv.Atoi(token)
	if err != nil {
		return 0, ErrPathNotFound
	}
	return index, nil
}

// navigateToParent navigates to the parent of the target path and returns the parent, key, and any error
func navigateToParent(doc any, path []string) (any, any, error) {
	if len(path) == 0 {
		return nil, nil, ErrPathNotFound
	}

	parentPath := path[:len(path)-1]
	key := path[len(path)-1]

	parent := doc
	if len(parentPath) > 0 {
		var err error
		parent, err = value(doc, parentPath)
		if err != nil {
			return nil, nil, ErrPathNotFound
		}
	}

	// Convert key to appropriate type based on parent
	switch parent.(type) {
	case map[string]any:
		return parent, key, nil
	case []any:
		index, err := parseArrayIndex(key)
		if err != nil {
			return nil, nil, ErrPathNotFound
		}
		return parent, index, nil
	default:
		return nil, nil, ErrPathNotFound
	}
}

// valueFromParent retrieves a value from a parent container using a key.
// Supports both map[string]any (with string keys) and []any (with int keys).
func valueFromParent(parent any, key any) any {
	switch p := parent.(type) {
	case map[string]any:
		if k, ok := key.(string); ok {
			return p[k]
		}
		return nil
	case []any:
		if k, ok := key.(int); ok {
			if k >= 0 && k < len(p) {
				return p[k]
			}
		}
		return nil
	default:
		return nil
	}
}

// setValueAtPath sets a value at a specific path in the document
func setValueAtPath(doc any, path []string, value any) error {
	return setValueAtPathWithMode(doc, path, value, false)
}

// insertValueAtPath inserts a value at a specific path in the document (for arrays, it inserts rather than replaces)
func insertValueAtPath(doc any, path []string, value any) error {
	return setValueAtPathWithMode(doc, path, value, true)
}

// setValueAtPathWithMode sets or inserts a value at a specific path in the document
func setValueAtPathWithMode(doc any, path []string, value any, insertMode bool) error {
	if len(path) == 0 {
		// Root level set - this should be handled by the caller
		return ErrPathNotFound
	}

	parent, key, err := navigateToParent(doc, path)
	if err != nil {
		return err
	}

	// Handle array operations specially
	if slice, ok := parent.([]any); ok {
		if index, ok := key.(int); ok && index >= 0 && index <= len(slice) {
			if insertMode {
				// Insert mode: always insert/append
				var newSlice []any
				if index == len(slice) {
					// Append to end
					newSlice = append(newSlice, slice...)
					newSlice = append(newSlice, value)
				} else {
					// Insert at index - move elements to make room
					newSlice = make([]any, len(slice)+1)
					copy(newSlice[:index], slice[:index])
					newSlice[index] = value
					copy(newSlice[index+1:], slice[index:])
				}

				// We need to replace the array in its parent context
				if len(path) == 1 {
					// This is modifying the root array, but we can't change doc directly
					return ErrCannotModifyRootArray
				}
				return updateGrandparent(doc, path, newSlice)
			}
			// Set mode: replace if within bounds, append if at end
			if index < len(slice) {
				// This is a replacement, use normal updateParent
				return updateParent(parent, key, value)
			}
			// This is append at end
			newSlice := make([]any, len(slice)+1)
			copy(newSlice, slice)
			newSlice[len(slice)] = value

			// We need to replace the array in its parent context
			if len(path) == 1 {
				// This is modifying the root array, but we can't change doc directly
				return ErrCannotModifyRootArray
			}
			return updateGrandparent(doc, path, newSlice)
		}
	}

	return updateParent(parent, key, value)
}

// updateGrandparent updates a grandparent container with a new value for the given key.
// It handles both root-level parents and nested grandparents.
func updateGrandparent(doc any, path []string, newSlice []any) error {
	grandParentPath := path[:len(path)-2]
	grandParentKey := path[len(path)-2]

	if len(grandParentPath) == 0 {
		docMap, ok := doc.(map[string]any)
		if !ok {
			return ErrCannotUpdateParent
		}
		docMap[grandParentKey] = newSlice
		return nil
	}

	grandParent, err := value(doc, grandParentPath)
	if err != nil {
		return err
	}
	grandParentMap, ok := grandParent.(map[string]any)
	if !ok {
		return ErrCannotUpdateGrandparent
	}
	grandParentMap[grandParentKey] = newSlice
	return nil
}

// updateParent updates the parent container with a new value
func updateParent(parent any, key any, value any) error {
	switch p := parent.(type) {
	case map[string]any:
		k, ok := key.(string)
		if !ok {
			return ErrInvalidKeyTypeMap
		}
		p[k] = value
		return nil
	case []any:
		k, ok := key.(int)
		if !ok {
			return ErrInvalidKeyTypeSlice
		}
		if k >= 0 && k < len(p) {
			p[k] = value
			return nil
		}
		return ErrIndexOutOfRange
	default:
		return ErrUnsupportedParentType
	}
}

// pathExists checks if a path exists in the document
func pathExists(doc any, path []string) bool {
	if len(path) == 0 {
		return true
	}

	_, err := jsonpointer.Get(doc, path...)
	return err == nil
}

// ToFloat64 converts a value to float64, handling various numeric types, booleans, and strings.
// This matches JavaScript's Number() behavior for type coercion.
// Optimized for common numeric types with fast paths.
func ToFloat64(val any) (float64, bool) {
	// Handle nil (null in JSON) - JavaScript Number(null) returns 0
	if val == nil {
		return 0, true
	}

	switch v := val.(type) {
	case float64:
		return v, true
	case int:
		return float64(v), true
	case float32:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case bool:
		// Convert boolean to number: true -> 1, false -> 0
		if v {
			return 1, true
		}
		return 0, true
	case uint:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case string:
		return parseStringToFloat(v)
	default:
		return 0, false
	}
}

// parseStringToFloat optimizes string-to-float conversion with fast paths
func parseStringToFloat(v string) (float64, bool) {
	if len(v) == 0 {
		return 0, true
	}

	// Fast path: check for simple numeric strings without trimming
	if len(v) <= 15 && isSimpleNumeric(v) {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}

	// Slower path: trim whitespace and try parsing
	trimmed := strings.TrimSpace(v)
	if trimmed == "" {
		return 0, true
	}

	f, err := strconv.ParseFloat(trimmed, 64)
	return f, err == nil
}

// isSimpleNumeric checks if string contains only numeric characters (fast path)
func isSimpleNumeric(s string) bool {
	if len(s) == 0 {
		return false
	}

	i := 0
	// Allow leading minus sign
	if s[0] == '-' {
		if len(s) == 1 {
			return false
		}
		i = 1
	}

	hasDigit := false
	hasDot := false

	for ; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= '0' && c <= '9':
			hasDigit = true
		case c == '.' && !hasDot:
			hasDot = true
		case c == 'e' || c == 'E':
			// Handle scientific notation - delegate to strconv.ParseFloat
			return hasDigit
		default:
			return false
		}
	}

	return hasDigit
}

// toString converts a value to string, handling string and []byte types.
func toString(value any) (string, error) {
	if value == nil {
		return "", ErrCannotConvertNilToString
	}
	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", ErrCannotConvertToString
	}
}

// validateString retrieves and validates the string value at the path.
// Returns the original value, the string representation, and any error.
func validateString(doc any, path []string) (any, string, error) {
	val, err := value(doc, path)
	if err != nil {
		return nil, "", ErrPathNotFound
	}
	str, err := toString(val)
	if err != nil {
		return nil, "", ErrNotString
	}
	return val, str, nil
}

// floatToJSONValue converts a float64 to int if it's a whole number,
// for cleaner JSON serialization (e.g., 5.0 -> 5).
func floatToJSONValue(f float64) any {
	if f == float64(int(f)) {
		return int(f)
	}
	return f
}

// predicateOpsToJSON converts a slice of predicate operations to JSON format.
// Used by composite predicates (And, Or, Not) to serialize their sub-operations.
func predicateOpsToJSON(operations []any, errInvalid error) ([]internal.Operation, error) {
	result := make([]internal.Operation, 0, len(operations))
	for _, op := range operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, errInvalid
		}
		jsonVal, err := predicateOp.ToJSON()
		if err != nil {
			return nil, err
		}
		result = append(result, jsonVal)
	}
	return result, nil
}

// predicateOpsToCompact converts a slice of predicate operations to compact format.
// Used by composite predicates (And, Or, Not) to serialize their sub-operations.
func predicateOpsToCompact(operations []any, errInvalid error) ([]any, error) {
	result := make([]any, 0, len(operations))
	for _, op := range operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return nil, errInvalid
		}
		compact, err := predicateOp.ToCompact()
		if err != nil {
			return nil, err
		}
		result = append(result, compact)
	}
	return result, nil
}

// validatePredicateOps validates a slice of predicate operations.
// Used by composite predicates (And, Or, Not) to validate their sub-operations.
func validatePredicateOps(operations []any, errInvalid error) error {
	for _, op := range operations {
		predicateOp, ok := op.(internal.PredicateOp)
		if !ok {
			return errInvalid
		}
		if err := predicateOp.Validate(); err != nil {
			return err
		}
	}
	return nil
}
