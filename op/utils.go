package op

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/kaptinlin/deepclone"
	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpointer"
)

// pathEquals checks if two paths are equal.
func pathEquals(p1, p2 []string) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i, v := range p1 {
		if v != p2[i] {
			return false
		}
	}
	return true
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

// IsPredicateOp checks if an operation type is a predicate operation.
func IsPredicateOp(opType internal.OpType) bool {
	switch opType {
	case internal.OpTestType, internal.OpDefinedType, internal.OpUndefinedType, internal.OpTestTypeType,
		internal.OpTestStringType, internal.OpTestStringLenType, internal.OpContainsType, internal.OpEndsType,
		internal.OpStartsType, internal.OpInType, internal.OpLessType, internal.OpMoreType, internal.OpMatchesType,
		internal.OpAndType, internal.OpOrType, internal.OpNotType:
		return true
	case internal.OpAddType, internal.OpRemoveType, internal.OpReplaceType, internal.OpMoveType, internal.OpCopyType,
		internal.OpTypeType, internal.OpFlipType, internal.OpIncType, internal.OpStrInsType, internal.OpStrDelType,
		internal.OpSplitType, internal.OpMergeType, internal.OpExtendType:
		return false
	default:
		return false
	}
}

// IsMutationOp checks if an operation type is a mutation operation.
func IsMutationOp(opType internal.OpType) bool {
	switch opType {
	case internal.OpAddType, internal.OpRemoveType, internal.OpReplaceType, internal.OpMoveType, internal.OpCopyType,
		internal.OpIncType, internal.OpFlipType, internal.OpStrInsType, internal.OpStrDelType, internal.OpSplitType,
		internal.OpMergeType, internal.OpExtendType:
		return true
	case internal.OpTestType, internal.OpContainsType, internal.OpDefinedType, internal.OpUndefinedType,
		internal.OpTypeType, internal.OpTestTypeType, internal.OpTestStringType, internal.OpTestStringLenType,
		internal.OpEndsType, internal.OpStartsType, internal.OpInType, internal.OpLessType, internal.OpMoreType,
		internal.OpMatchesType, internal.OpAndType, internal.OpOrType, internal.OpNotType:
		return false
	default:
		return false
	}
}

// IsSecondOrderPredicateOp checks if an operation type is a second-order predicate operation.
func IsSecondOrderPredicateOp(opType internal.OpType) bool {
	switch opType {
	case internal.OpAndType, internal.OpOrType, internal.OpNotType:
		return true
	case internal.OpAddType, internal.OpRemoveType, internal.OpReplaceType, internal.OpMoveType, internal.OpCopyType,
		internal.OpTestType, internal.OpContainsType, internal.OpDefinedType, internal.OpUndefinedType,
		internal.OpTypeType, internal.OpTestTypeType, internal.OpTestStringType, internal.OpTestStringLenType,
		internal.OpEndsType, internal.OpStartsType, internal.OpInType, internal.OpLessType, internal.OpMoreType,
		internal.OpMatchesType, internal.OpFlipType, internal.OpIncType, internal.OpStrInsType, internal.OpStrDelType,
		internal.OpSplitType, internal.OpMergeType, internal.OpExtendType:
		return false
	default:
		return false
	}
}

// getValue retrieves a value from a document using a path.
func getValue(doc interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return doc, nil
	}

	value, err := jsonpointer.Get(doc, path...)
	if err != nil {
		return nil, ErrPathNotFound
	}
	return value, nil
}

// getNumericValue retrieves a numeric value from the document at the given path.
// Returns the original value, the float64 conversion, and any error.
func getNumericValue(doc any, path []string) (any, float64, error) {
	value, err := getValue(doc, path)
	if err != nil {
		return nil, 0, ErrPathNotFound
	}
	numValue, ok := ToFloat64(value)
	if !ok {
		return nil, 0, ErrNotNumber
	}
	return value, numValue, nil
}

// deepEqual performs a deep equality check between two values.
// Optimized to avoid expensive reflect.DeepEqual for common types.
func deepEqual(a, b interface{}) bool {
	// Fast path: both nil
	if a == nil && b == nil {
		return true
	}

	// Fast path: one nil, other not
	if a == nil || b == nil {
		return false
	}

	// Fast path: strings (very common)
	if aStr, aIsStr := a.(string); aIsStr {
		if bStr, bIsStr := b.(string); bIsStr {
			return aStr == bStr
		}
		return false
	}

	// Fast path: booleans
	if aBool, aIsBool := a.(bool); aIsBool {
		if bBool, bIsBool := b.(bool); bIsBool {
			return aBool == bBool
		}
		return false
	}

	// Fast path: numeric types (only if both are actual numbers, not string conversions)
	aFloat, aIsNum := toNumericValue(a)
	bFloat, bIsNum := toNumericValue(b)
	if aIsNum && bIsNum {
		return aFloat == bFloat
	}

	// Only one is numeric - not equal
	if aIsNum != bIsNum {
		return false
	}

	// Fast path: try direct comparison for comparable types
	// Use defer+recover to handle uncomparable types gracefully
	equal := false
	canCompare := true
	func() {
		defer func() {
			if recover() != nil {
				canCompare = false
			}
		}()
		equal = (a == b)
	}()

	if canCompare {
		return equal
	}

	// Slow path: complex types (maps, slices, structs)
	return reflect.DeepEqual(a, b)
}

// toNumericValue converts a value to float64 if it's an actual numeric type
// (not a string that could be parsed as a number).
// Returns the float64 value and true if successful, 0 and false otherwise.
func toNumericValue(val interface{}) (float64, bool) {
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
func DeepClone(value interface{}) (interface{}, error) {
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
func navigateToParent(doc interface{}, path []string) (interface{}, interface{}, error) {
	if len(path) == 0 {
		return nil, nil, ErrPathNotFound
	}

	parentPath := path[:len(path)-1]
	key := path[len(path)-1]

	// Navigate to parent
	var parent interface{}
	if len(parentPath) == 0 {
		parent = doc
	} else {
		var err error
		parent, err = getValue(doc, parentPath)
		if err != nil {
			return nil, nil, ErrPathNotFound
		}
	}

	// Convert key to appropriate type
	switch parent.(type) {
	case map[string]interface{}:
		return parent, key, nil
	case []interface{}:
		// Try to parse as integer for array index
		if index, err := parseArrayIndex(key); err == nil {
			return parent, index, nil
		}
		return nil, nil, ErrPathNotFound
	default:
		return nil, nil, ErrPathNotFound
	}
}

// getValueFromParent retrieves a value from a parent container using a key
func getValueFromParent(parent interface{}, key interface{}) interface{} {
	switch p := parent.(type) {
	case map[string]interface{}:
		if k, ok := key.(string); ok {
			return p[k]
		}
		return nil
	case []interface{}:
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
func setValueAtPath(doc interface{}, path []string, value interface{}) error {
	return setValueAtPathWithMode(doc, path, value, false)
}

// insertValueAtPath inserts a value at a specific path in the document (for arrays, it inserts rather than replaces)
func insertValueAtPath(doc interface{}, path []string, value interface{}) error {
	return setValueAtPathWithMode(doc, path, value, true)
}

// setValueAtPathWithMode sets or inserts a value at a specific path in the document
func setValueAtPathWithMode(doc interface{}, path []string, value interface{}, insertMode bool) error {
	if len(path) == 0 {
		// Root level set - this should be handled by the caller
		return ErrPathNotFound
	}

	parent, key, err := navigateToParent(doc, path)
	if err != nil {
		return err
	}

	// Handle array operations specially
	if slice, ok := parent.([]interface{}); ok {
		if index, ok := key.(int); ok && index >= 0 && index <= len(slice) {
			if insertMode {
				// Insert mode: always insert/append
				var newSlice []interface{}
				if index == len(slice) {
					// Append to end
					newSlice = append(newSlice, slice...)
					newSlice = append(newSlice, value)
				} else {
					// Insert at index - move elements to make room
					newSlice = make([]interface{}, len(slice)+1)
					copy(newSlice[:index], slice[:index])
					newSlice[index] = value
					copy(newSlice[index+1:], slice[index:])
				}

				// We need to replace the array in its parent context
				if len(path) == 1 {
					// This is modifying the root array, but we can't change doc directly
					return ErrCannotModifyRootArray
				}
				// Get grandparent and update the parent
				grandParentPath := path[:len(path)-2]
				grandParentKey := path[len(path)-2]
				if len(grandParentPath) == 0 {
					// Parent is in root
					docMap, ok := doc.(map[string]interface{})
					if !ok {
						return ErrCannotUpdateParent
					}
					docMap[grandParentKey] = newSlice
					return nil
				}
				grandParent, err := getValue(doc, grandParentPath)
				if err != nil {
					return err
				}
				grandParentMap, ok := grandParent.(map[string]interface{})
				if !ok {
					return ErrCannotUpdateGrandparent
				}
				grandParentMap[grandParentKey] = newSlice
				return nil
			}
			// Set mode: replace if within bounds, append if at end
			if index < len(slice) {
				// This is a replacement, use normal updateParent
				return updateParent(parent, key, value)
			}
			// This is append at end
			newSlice := make([]interface{}, len(slice)+1)
			copy(newSlice, slice)
			newSlice[len(slice)] = value

			// We need to replace the array in its parent context
			if len(path) == 1 {
				// This is modifying the root array, but we can't change doc directly
				return ErrCannotModifyRootArray
			}
			// Get grandparent and update the parent
			grandParentPath := path[:len(path)-2]
			grandParentKey := path[len(path)-2]
			if len(grandParentPath) == 0 {
				// Parent is in root
				docMap, ok := doc.(map[string]interface{})
				if !ok {
					return ErrCannotUpdateParent
				}
				docMap[grandParentKey] = newSlice
				return nil
			}
			grandParent, err := getValue(doc, grandParentPath)
			if err != nil {
				return err
			}
			grandParentMap, ok := grandParent.(map[string]interface{})
			if !ok {
				return ErrCannotUpdateGrandparent
			}
			grandParentMap[grandParentKey] = newSlice
			return nil
		}
	}

	return updateParent(parent, key, value)
}

// updateParent updates the parent container with a new value
func updateParent(parent interface{}, key interface{}, value interface{}) error {
	switch p := parent.(type) {
	case map[string]interface{}:
		if k, ok := key.(string); ok {
			p[k] = value
			return nil
		}
		return ErrInvalidKeyTypeMap
	case []interface{}:
		if k, ok := key.(int); ok {
			if k >= 0 && k < len(p) {
				p[k] = value
				return nil
			} else if k == len(p) {
				// Allow appending to the end of the array
				// Note: We can't modify the slice header here, this needs to be handled by the caller
				return ErrIndexOutOfRange
			}
			return ErrIndexOutOfRange
		}
		return ErrInvalidKeyTypeSlice
	default:
		return ErrUnsupportedParentType
	}
}

// pathExists checks if a path exists in the document
func pathExists(doc interface{}, path []string) bool {
	if len(path) == 0 {
		return true
	}

	_, err := jsonpointer.Get(doc, path...)
	return err == nil
}

// ToFloat64 converts a value to float64, handling various numeric types, booleans, and strings.
// This matches JavaScript's Number() behavior for type coercion.
// Optimized for common numeric types with fast paths.
func ToFloat64(val interface{}) (float64, bool) {
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
	// Fast path: check for empty string first (common case)
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

	if f, err := strconv.ParseFloat(trimmed, 64); err == nil {
		return f, true
	}

	return 0, false
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

// toString converts a value to string, handling various internal.
func toString(value interface{}) (string, error) {
	if value == nil {
		return "", ErrCannotConvertNilToString
	}

	switch v := value.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		// For other types, try to convert using reflection
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.String:
			return reflect.ValueOf(v).String(), nil
		case reflect.Slice:
			if rt.Elem().Kind() == reflect.Uint8 {
				return string(reflect.ValueOf(v).Bytes()), nil
			}
			return "", ErrCannotConvertNilToString
		case reflect.Invalid:
			return "", ErrCannotConvertNilToString
		case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
			return "", ErrCannotConvertNilToString
		case reflect.Array, reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
			reflect.Ptr, reflect.Struct, reflect.UnsafePointer:
			return "", ErrCannotConvertNilToString
		default:
			return "", ErrCannotConvertNilToString
		}
	}
}
