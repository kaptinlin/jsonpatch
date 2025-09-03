// Package json implements JSON codec for JSON Patch operations.
// Provides encoding and decoding functionality for JSON Patch operations with full RFC 6902 support.
package json

import (
	"errors"
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Decode operation errors - define clearly and concisely
var (
	ErrInvalidPointer            = errors.New("invalid pointer")
	ErrAddOpMissingValue         = errors.New("add operation missing 'value' field")
	ErrReplaceOpMissingValue     = errors.New("replace operation missing 'value' field")
	ErrNotOpRequiresOperand      = errors.New("not operation requires at least one operand")
	ErrMissingValueField         = errors.New("missing value field")
	ErrEmptyTypeList             = errors.New("empty type list")
	ErrInvalidType               = errors.New("invalid type")
	ErrNotOpRequiresValidOperand = errors.New("not operation requires a valid predicate operand")
)

// toPath converts string path to jsonpointer.Path.
func toPath(pathStr string) jsonpointer.Path {
	return jsonpointer.Parse(pathStr)
}

// pathToStringSlice converts jsonpointer.Path to []string for op constructors.
func pathToStringSlice(path jsonpointer.Path) []string {
	result := make([]string, len(path))
	for i, token := range path {
		result[i] = fmt.Sprintf("%v", token)
	}
	return result
}

// OperationToOp converts JSON operation to Op instance.
func OperationToOp(operation map[string]interface{}, options internal.JsonPatchOptions) (internal.Op, error) {
	opType, ok := operation["op"].(string)
	if !ok {
		return nil, ErrOpMissingOpField
	}

	pathStr, ok := operation["path"].(string)
	if !ok {
		return nil, ErrOpMissingPathField
	}

	// Validate JSON pointer format
	if err := jsonpointer.Validate(pathStr); err != nil {
		return nil, ErrInvalidPointer
	}

	path := pathToStringSlice(toPath(pathStr))

	switch opType {
	case "add":
		_, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrAddOpMissingValue
		}
		return op.NewOpAddOperation(path, operation["value"]), nil
	case "remove":
		// Check for oldValue field
		if oldValue, hasOldValue := operation["oldValue"]; hasOldValue {
			return op.NewOpRemoveOperationWithOldValue(path, oldValue), nil
		}
		return op.NewOpRemoveOperation(path), nil
	case "replace":
		_, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrReplaceOpMissingValue
		}

		// Check for oldValue field
		if oldValue, hasOldValue := operation["oldValue"]; hasOldValue {
			return op.NewOpReplaceOperationWithOldValue(path, operation["value"], oldValue), nil
		}
		return op.NewOpReplaceOperation(path, operation["value"]), nil
	case "move":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrMoveOpMissingFrom
		}
		return op.NewOpMoveOperation(path, pathToStringSlice(toPath(fromStr))), nil
	case "copy":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrCopyOpMissingFrom
		}
		return op.NewOpCopyOperation(path, pathToStringSlice(toPath(fromStr))), nil
	case "flip":
		return op.NewOpFlipOperation(path), nil
	case "inc":
		incVal, ok := toFloat64(operation["inc"])
		if !ok {
			return nil, ErrIncOpMissingInc
		}
		return op.NewOpIncOperation(path, incVal), nil
	case "str_ins":
		pos, ok := toFloat64(operation["pos"])
		if !ok {
			return nil, ErrStrInsOpMissingPos
		}
		str, ok := operation["str"].(string)
		if !ok {
			return nil, ErrStrInsOpMissingStr
		}
		return op.NewOpStrInsOperation(path, pos, str), nil
	case "str_del":
		pos, ok := toFloat64(operation["pos"])
		if !ok {
			return nil, ErrStrDelOpMissingPos
		}
		// str_del can have either str or len parameter
		if str, ok := operation["str"].(string); ok {
			return op.NewOpStrDelOperationWithStr(path, pos, str), nil
		}
		if lenVal, ok := toFloat64(operation["len"]); ok {
			return op.NewOpStrDelOperation(path, pos, lenVal), nil
		}
		return nil, ErrStrDelOpMissingFields
	case "split":
		pos, ok := toFloat64(operation["pos"])
		if !ok {
			return nil, ErrSplitOpMissingPos
		}
		props := operation["props"]
		return op.NewOpSplitOperation(path, pos, props), nil
	case "merge":
		var props map[string]interface{}
		if p, ok := operation["props"].(map[string]interface{}); ok {
			props = p
		} else {
			props = make(map[string]interface{}) // Default to empty map
		}
		pos := float64(0) // Default position
		if posVal, ok := toFloat64(operation["pos"]); ok {
			pos = posVal
		}
		return op.NewOpMergeOperation(path, pos, props), nil
	case "extend":
		props, ok := operation["props"].(map[string]interface{})
		if !ok {
			return nil, ErrValueNotObject
		}
		deleteNull := false
		if dn, ok := operation["deleteNull"].(bool); ok {
			deleteNull = dn
		}
		return op.NewOpExtendOperation(path, props, deleteNull), nil
	case "and":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrAndOpMissingApply
		}
		// Convert each sub-operation to a proper PredicateOp
		predicateOps := make([]interface{}, 0, len(apply))
		for _, subOp := range apply {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				// Merge paths if needed
				subPath := ""
				if sp, ok := subOpMap["path"].(string); ok {
					subPath = sp
				}
				mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
				subOpMap["path"] = formatPath(mergedPath)

				predicateOp, err := OperationToPredicateOp(subOpMap, options)
				if err != nil {
					return nil, err
				}
				predicateOps = append(predicateOps, predicateOp)
			}
		}
		return op.NewOpAndOperation(path, predicateOps), nil
	case "or":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrOrOpMissingApply
		}
		// Convert each sub-operation to a proper PredicateOp
		predicateOps := make([]interface{}, 0, len(apply))
		for _, subOp := range apply {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				// Merge paths if needed
				subPath := ""
				if sp, ok := subOpMap["path"].(string); ok {
					subPath = sp
				}
				mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
				subOpMap["path"] = formatPath(mergedPath)

				predicateOp, err := OperationToPredicateOp(subOpMap, options)
				if err != nil {
					return nil, err
				}
				predicateOps = append(predicateOps, predicateOp)
			}
		}
		return op.NewOpOrOperation(path, predicateOps), nil
	case "not":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrNotOpMissingApply
		}
		if len(apply) == 0 {
			return nil, ErrNotOpRequiresOperand
		}
		// Convert each sub-operation to a proper PredicateOp
		predicateOps := make([]interface{}, 0, len(apply))
		for _, subOp := range apply {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				// Merge paths if needed
				subPath := ""
				if sp, ok := subOpMap["path"].(string); ok {
					subPath = sp
				}
				mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
				subOpMap["path"] = formatPath(mergedPath)

				predicateOp, err := OperationToPredicateOp(subOpMap, options)
				if err != nil {
					return nil, err
				}
				predicateOps = append(predicateOps, predicateOp)
			}
		}
		return op.NewOpNotOperationMultiple(path, predicateOps), nil
	default:
		return OperationToPredicateOp(operation, options)
	}
}

// OperationToPredicateOp converts JSON operation to PredicateOp instance.
func OperationToPredicateOp(operation map[string]interface{}, options internal.JsonPatchOptions) (internal.Op, error) {
	opType, ok := operation["op"].(string)
	if !ok {
		return nil, ErrOpMissingOpField
	}

	pathStr, ok := operation["path"].(string)
	if !ok {
		return nil, ErrOpMissingPathField
	}

	// Validate JSON pointer format
	if err := jsonpointer.Validate(pathStr); err != nil {
		return nil, ErrInvalidPointer
	}

	path := pathToStringSlice(toPath(pathStr))

	switch opType {
	case "test":
		// Check if value field exists (it's required for test operations)
		value, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrMissingValueField
		}
		testOp := op.NewOpTestOperation(path, value)
		// Check for "not" field
		if notVal, ok := operation["not"].(bool); ok && notVal {
			testOp.NotFlag = true
		}
		return testOp, nil
	case "defined":
		return op.NewOpDefinedOperation(path), nil
	case "undefined":
		return op.NewOpUndefinedOperation(path, false), nil
	case "type":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrTypeOpMissingValue
		}
		return op.NewOpTypeOperation(path, value), nil
	case "test_type":
		// Handle both single type string and array of types
		typeField := operation["type"]
		if typeStr, ok := typeField.(string); ok {
			// Validate single type
			if err := validateSingleTestType(typeStr); err != nil {
				return nil, err
			}
			return op.NewOpTestTypeOperation(path, typeStr), nil
		} else if typeSlice, ok := typeField.([]interface{}); ok {
			if len(typeSlice) == 0 {
				return nil, ErrEmptyTypeList
			}
			// Convert to []string and validate all types
			typeStrings := make([]string, len(typeSlice))
			for i, t := range typeSlice {
				typeStr, isString := t.(string)
				if !isString {
					return nil, ErrInvalidType
				}
				if err := validateSingleTestType(typeStr); err != nil {
					return nil, err
				}
				typeStrings[i] = typeStr
			}
			return op.NewOpTestTypeOperationMultiple(path, typeStrings), nil
		} else if typeStringSlice, ok := typeField.([]string); ok {
			if len(typeStringSlice) == 0 {
				return nil, ErrEmptyTypeList
			}
			// Validate all types in the array
			for _, typeStr := range typeStringSlice {
				if err := validateSingleTestType(typeStr); err != nil {
					return nil, err
				}
			}
			return op.NewOpTestTypeOperationMultiple(path, typeStringSlice), nil
		}
		return nil, ErrTestTypeOpMissingType
	case "test_string":
		str, hasStr := operation["str"].(string)
		if !hasStr {
			return nil, ErrTestStringOpMissingStr
		}
		pos := float64(0)
		if posVal, ok := toFloat64(operation["pos"]); ok {
			pos = posVal
		}
		if pos != 0 {
			return op.NewOpTestStringOperationWithPos(path, str, pos), nil
		}
		return op.NewOpTestStringOperation(path, str), nil
	case "test_string_len":
		lenVal, ok := toFloat64(operation["len"])
		if !ok {
			return nil, ErrTestStringLenOpMissingLen
		}

		// Check for not flag
		not := false
		if n, ok := operation["not"].(bool); ok {
			not = n
		}

		if not {
			return op.NewOpTestStringLenOperationWithNot(path, lenVal, not), nil
		}
		return op.NewOpTestStringLenOperation(path, lenVal), nil
	case "contains":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrContainsOpMissingValue
		}

		// Check for ignore_case flag
		ignoreCase := false
		if ic, ok := operation["ignore_case"].(bool); ok {
			ignoreCase = ic
		}

		if ignoreCase {
			return op.NewOpContainsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpContainsOperation(path, value), nil
	case "ends":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrEndsOpMissingValue
		}

		// Check for ignore_case flag
		ignoreCase := false
		if ic, ok := operation["ignore_case"].(bool); ok {
			ignoreCase = ic
		}

		if ignoreCase {
			return op.NewOpEndsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpEndsOperation(path, value), nil
	case "starts":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrStartsOpMissingValue
		}

		// Check for ignore_case flag
		ignoreCase := false
		if ic, ok := operation["ignore_case"].(bool); ok {
			ignoreCase = ic
		}

		if ignoreCase {
			return op.NewOpStartsOperationWithIgnoreCase(path, value, ignoreCase), nil
		}
		return op.NewOpStartsOperation(path, value), nil
	case "matches":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrMatchesOpMissingValue
		}
		ignoreCase := false
		if ic, ok := operation["ignore_case"].(bool); ok {
			ignoreCase = ic
		}
		matchesOp, err := op.NewOpMatchesOperation(path, value, ignoreCase)
		if err != nil {
			return nil, err
		}
		return matchesOp, nil
	case "in":
		value := operation["value"]
		if values, ok := value.([]interface{}); ok {
			return op.NewOpInOperation(path, values), nil
		}
		return op.NewOpInOperation(path, []interface{}{value}), nil
	case "less":
		value, ok := toFloat64(operation["value"])
		if !ok {
			return nil, ErrLessOpMissingValue
		}
		return op.NewOpLessOperation(path, value), nil
	case "more":
		value, ok := toFloat64(operation["value"])
		if !ok {
			return nil, ErrMoreOpMissingValue
		}
		return op.NewOpMoreOperation(path, value), nil
	case "and":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrAndOpMissingApply
		}
		// Convert each sub-operation to a proper PredicateOp
		predicateOps := make([]interface{}, 0, len(apply))
		for _, subOp := range apply {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				// Merge paths if needed
				subPath := ""
				if sp, ok := subOpMap["path"].(string); ok {
					subPath = sp
				}
				mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
				subOpMap["path"] = formatPath(mergedPath)

				predicateOp, err := OperationToPredicateOp(subOpMap, options)
				if err != nil {
					return nil, err
				}
				predicateOps = append(predicateOps, predicateOp)
			}
		}
		return op.NewOpAndOperation(path, predicateOps), nil
	case "or":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrOrOpMissingApply
		}
		// Convert each sub-operation to a proper PredicateOp
		predicateOps := make([]interface{}, 0, len(apply))
		for _, subOp := range apply {
			if subOpMap, ok := subOp.(map[string]interface{}); ok {
				// Merge paths if needed
				subPath := ""
				if sp, ok := subOpMap["path"].(string); ok {
					subPath = sp
				}
				mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
				subOpMap["path"] = formatPath(mergedPath)

				predicateOp, err := OperationToPredicateOp(subOpMap, options)
				if err != nil {
					return nil, err
				}
				predicateOps = append(predicateOps, predicateOp)
			}
		}
		return op.NewOpOrOperation(path, predicateOps), nil
	case "not":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrNotOpMissingApply
		}
		if len(apply) == 0 {
			return nil, ErrNotOpRequiresOperand
		}
		// For not operation, we need to create a single predicate op
		if applyMap, ok := apply[0].(map[string]interface{}); ok {
			// Merge paths
			subPath := ""
			if sp, ok := applyMap["path"].(string); ok {
				subPath = sp
			}
			mergedPath := mergePaths(toPath(pathStr), toPath(subPath))
			applyMap["path"] = formatPath(mergedPath)

			operand, err := OperationToPredicateOp(applyMap, options)
			if err != nil {
				return nil, err
			}
			if predicateOp, ok := operand.(internal.PredicateOp); ok {
				return op.NewOpNotOperation(predicateOp), nil
			}
		}
		return nil, ErrNotOpRequiresValidOperand
	default:
		return nil, ErrCodecOpUnknown
	}
}

// mergePaths merges two paths for composite operations.
func mergePaths(basePath, subPath jsonpointer.Path) jsonpointer.Path {
	result := make(jsonpointer.Path, 0, len(basePath)+len(subPath))
	result = append(result, basePath...)
	result = append(result, subPath...)
	return result
}

// formatPath converts path back to string format for JSON operations.
func formatPath(path jsonpointer.Path) string {
	// Convert jsonpointer.Path to []string first
	pathSlice := make([]string, len(path))
	for i, token := range path {
		pathSlice[i] = fmt.Sprintf("%v", token)
	}
	return jsonpointer.Format(pathSlice...)
}

// Decode converts JSON operations to Op instances.
func Decode(operations []map[string]interface{}, options internal.JsonPatchOptions) ([]internal.Op, error) {
	ops := make([]internal.Op, 0, len(operations))
	for _, operation := range operations {
		o, err := OperationToOp(operation, options)
		if err != nil {
			return nil, err
		}
		ops = append(ops, o)
	}
	return ops, nil
}

// DecodeJSON converts JSON bytes to Op instances.
func DecodeJSON(data []byte, options internal.JsonPatchOptions) ([]internal.Op, error) {
	var operations []map[string]interface{}
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, err
	}
	return Decode(operations, options)
}

// validateSingleTestType validates a single type string for test_type operations
func validateSingleTestType(typeStr string) error {
	validTypes := map[string]bool{
		"string":  true,
		"number":  true,
		"boolean": true,
		"object":  true,
		"integer": true,
		"array":   true,
		"null":    true,
	}
	if !validTypes[typeStr] {
		return ErrInvalidType
	}
	return nil
}

// toFloat64 converts various numeric types to float64
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
	default:
		return 0, false
	}
}
