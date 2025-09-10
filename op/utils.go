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
func formatPath(path []string) string {
	if len(path) == 0 {
		return ""
	}

	var builder strings.Builder
	estimatedSize := len(path) * 9
	builder.Grow(estimatedSize)

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

// deepEqual performs a deep equality check between two values.
func deepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
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
		return 0, ErrInvalidPath
	}
	return index, nil
}

// navigateToParent navigates to the parent of the target path and returns the parent, key, and any error
func navigateToParent(doc interface{}, path []string) (interface{}, interface{}, error) {
	if len(path) == 0 {
		return nil, nil, ErrInvalidPath
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
		return nil, nil, ErrInvalidPath
	default:
		return nil, nil, ErrInvalidPath
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
		return ErrInvalidPath
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

// deleteFromParent removes a value from the parent container
func deleteFromParent(parent interface{}, key interface{}) error {
	switch p := parent.(type) {
	case map[string]interface{}:
		if k, ok := key.(string); ok {
			if _, exists := p[k]; !exists {
				return ErrKeyDoesNotExist
			}
			delete(p, k)
			return nil
		}
		return ErrInvalidKeyTypeMap
	case []interface{}:
		if k, ok := key.(int); ok {
			if k < 0 || k >= len(p) {
				return ErrIndexOutOfRange
			}
			// For slices, we need to modify the original slice reference
			// This is a limitation - we can't easily modify slice length in place
			// Return an error that indicates the caller needs to handle this differently
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

// toFloat64 converts a value to float64, handling various numeric types, booleans, and strings.
func toFloat64(val interface{}) (float64, bool) {
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
	case bool:
		// Convert boolean to number: true -> 1, false -> 0
		if v {
			return 1, true
		}
		return 0, true
	case string:
		// Try to parse string as number
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
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
