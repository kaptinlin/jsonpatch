// Package json implements JSON codec for JSON Patch operations.
// Provides encoding and decoding functionality for JSON Patch operations with full RFC 6902 support.
package json

import (
	"fmt"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Note: All errors are now defined in errors.go for consistency

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
func OperationToOp(operation map[string]interface{}, options internal.JSONPatchOptions) (internal.Op, error) {
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
		return op.NewAdd(path, operation["value"]), nil
	case "remove":
		// Check for oldValue field
		if oldValue, hasOldValue := operation["oldValue"]; hasOldValue {
			return op.NewRemoveWithOldValue(path, oldValue), nil
		}
		return op.NewRemove(path), nil
	case "replace":
		_, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrReplaceOpMissingValue
		}

		// Check for oldValue field
		if oldValue, hasOldValue := operation["oldValue"]; hasOldValue {
			return op.NewReplaceWithOldValue(path, operation["value"], oldValue), nil
		}
		return op.NewReplace(path, operation["value"]), nil
	case "move":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrMoveOpMissingFrom
		}
		return op.NewMove(path, pathToStringSlice(toPath(fromStr))), nil
	case "copy":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrCopyOpMissingFrom
		}
		return op.NewCopy(path, pathToStringSlice(toPath(fromStr))), nil
	case "flip":
		return op.NewOpFlipOperation(path), nil
	case "inc":
		incField, hasInc := operation["inc"]
		if !hasInc {
			return nil, ErrIncOpMissingInc
		}
		incVal, ok := op.ToFloat64(incField)
		if !ok {
			return nil, ErrIncOpInvalidType
		}
		return op.NewOpIncOperation(path, incVal), nil
	case "str_ins":
		posVal, hasPosField := operation["pos"]
		if !hasPosField {
			return nil, ErrStrInsOpMissingPos
		}
		pos, ok := op.ToFloat64(posVal)
		if !ok {
			return nil, ErrStrInsOpMissingPos
		}
		str, ok := operation["str"].(string)
		if !ok {
			return nil, ErrStrInsOpMissingStr
		}
		return op.NewOpStrInsOperation(path, pos, str), nil
	case "str_del":
		posVal, hasPosField := operation["pos"]
		if !hasPosField {
			return nil, ErrStrDelOpMissingPos
		}
		pos, ok := op.ToFloat64(posVal)
		if !ok {
			return nil, ErrStrDelOpMissingPos
		}

		// str_del can have either str or len parameter
		if str, ok := operation["str"].(string); ok && str != "" {
			return op.NewOpStrDelOperationWithStr(path, pos, str), nil
		}
		if lenVal, ok := op.ToFloat64(operation["len"]); ok {
			return op.NewOpStrDelOperation(path, pos, lenVal), nil
		}
		return nil, ErrStrDelOpMissingFields
	case "split":
		posVal, hasPosField := operation["pos"]
		if !hasPosField {
			return nil, ErrSplitOpMissingPos
		}
		pos, ok := op.ToFloat64(posVal)
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
		if posVal, ok := op.ToFloat64(operation["pos"]); ok {
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
	case "not":
		// Note: "not" case uses NewOpNotOperationMultiple for multiple operands support
		// The "and" and "or" cases are handled by OperationToPredicateOp to avoid duplication
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrNotOpMissingApply
		}
		if len(apply) == 0 {
			return nil, ErrNotOpRequiresOperand
		}
		predicateOps, err := decodeCompositePredicates(apply, toPath(pathStr), options)
		if err != nil {
			return nil, err
		}
		return op.NewOpNotOperationMultiple(path, predicateOps), nil
	default:
		// Handle "and", "or", and all predicate operations
		return OperationToPredicateOp(operation, options)
	}
}

// OperationToPredicateOp converts JSON operation to PredicateOp instance.
func OperationToPredicateOp(operation map[string]interface{}, options internal.JSONPatchOptions) (internal.Op, error) {
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
		testOp := op.NewTest(path, value)
		// Check for "not" field
		if notVal, ok := operation["not"].(bool); ok && notVal {
			testOp.NotFlag = true
		}
		return testOp, nil
	case "defined":
		return op.NewOpDefinedOperation(path), nil
	case "undefined":
		return op.NewOpUndefinedOperation(path), nil
	case "type":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrTypeOpMissingValue
		}
		return op.NewOpTypeOperation(path, value), nil
	case "test_type":
		// Handle both single type string and array of types
		// First check for "type" field (standard), then fall back to "value" field (compatibility)
		typeField := operation["type"]
		if typeField == nil {
			// Check for value field as fallback for compatibility
			typeField = operation["value"]
		}

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
		if posVal, ok := op.ToFloat64(operation["pos"]); ok {
			pos = posVal
		}
		notFlag, _ := operation["not"].(bool)
		ignoreCase, _ := operation["ignore_case"].(bool)

		// Use the most specific constructor based on what fields are set
		if pos != 0 || notFlag || ignoreCase {
			return op.NewOpTestStringOperationWithIgnoreCase(path, str, pos, notFlag, ignoreCase), nil
		}
		return op.NewOpTestStringOperation(path, str), nil
	case "test_string_len":
		lenVal, ok := op.ToFloat64(operation["len"])
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
		// Contains operation can have any value type (string for string contains, any for array contains)
		value, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrContainsOpMissingValue
		}
		// Convert value to string (contains only works with strings)
		stringValue, ok := value.(string)
		if !ok {
			return nil, op.ErrContainsValueMustBeString
		}
		return op.NewOpContainsOperationWithIgnoreCase(path, stringValue, getBoolField(operation, "ignore_case")), nil
	case "ends":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrEndsOpMissingValue
		}
		return op.NewOpEndsOperationWithIgnoreCase(path, value, getBoolField(operation, "ignore_case")), nil
	case "starts":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrStartsOpMissingValue
		}
		return op.NewOpStartsOperationWithIgnoreCase(path, value, getBoolField(operation, "ignore_case")), nil
	case "matches":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrMatchesOpMissingValue
		}
		return op.NewOpMatchesOperation(path, value, getBoolField(operation, "ignore_case"), options.CreateMatcher), nil
	case "in":
		value := operation["value"]
		if values, ok := value.([]interface{}); ok {
			return op.NewOpInOperation(path, values), nil
		}
		return op.NewOpInOperation(path, []interface{}{value}), nil
	case "less":
		value, ok := op.ToFloat64(operation["value"])
		if !ok {
			return nil, ErrLessOpMissingValue
		}
		return op.NewOpLessOperation(path, value), nil
	case "more":
		value, ok := op.ToFloat64(operation["value"])
		if !ok {
			return nil, ErrMoreOpMissingValue
		}
		return op.NewOpMoreOperation(path, value), nil
	case "and":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrAndOpMissingApply
		}
		predicateOps, err := decodeCompositePredicates(apply, toPath(pathStr), options)
		if err != nil {
			return nil, err
		}
		return op.NewOpAndOperation(path, predicateOps), nil
	case "or":
		apply, ok := operation["apply"].([]interface{})
		if !ok {
			return nil, ErrOrOpMissingApply
		}
		predicateOps, err := decodeCompositePredicates(apply, toPath(pathStr), options)
		if err != nil {
			return nil, err
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

// decodeCompositePredicates decodes an array of sub-operations for and/or/not operations.
// It handles path merging and recursive predicate decoding.
func decodeCompositePredicates(
	apply []interface{},
	basePath jsonpointer.Path,
	options internal.JSONPatchOptions,
) ([]interface{}, error) {
	predicateOps := make([]interface{}, 0, len(apply))
	for _, subOp := range apply {
		subOpMap, ok := subOp.(map[string]interface{})
		if !ok {
			continue
		}
		// Merge paths if needed
		subPath := ""
		if sp, ok := subOpMap["path"].(string); ok {
			subPath = sp
		}
		mergedPath := mergePaths(basePath, toPath(subPath))
		subOpMap["path"] = formatPath(mergedPath)

		predicateOp, err := OperationToPredicateOp(subOpMap, options)
		if err != nil {
			return nil, err
		}
		predicateOps = append(predicateOps, predicateOp)
	}
	return predicateOps, nil
}

// getBoolField extracts a boolean field from an operation map with a default of false.
func getBoolField(operation map[string]interface{}, field string) bool {
	if v, ok := operation[field].(bool); ok {
		return v
	}
	return false
}

// mergePaths merges two paths for composite operations.
// If subPath is absolute (starts from root), it takes precedence over basePath.
// If subPath is empty, use basePath.
// Otherwise, combine them.
func mergePaths(basePath, subPath jsonpointer.Path) jsonpointer.Path {
	// If subPath is empty, use basePath
	if len(subPath) == 0 {
		return basePath
	}

	// If both paths are the same, don't merge - just use one
	if len(basePath) == len(subPath) {
		same := true
		for i := range basePath {
			if fmt.Sprintf("%v", basePath[i]) != fmt.Sprintf("%v", subPath[i]) {
				same = false
				break
			}
		}
		if same {
			return subPath // Use child path
		}
	}

	// Normal merging for different paths
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
func Decode(operations []map[string]interface{}, options internal.JSONPatchOptions) ([]internal.Op, error) {
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

// DecodeOperations converts Operation structs to Op instances using json/v2
func DecodeOperations(operations []internal.Operation, options internal.JSONPatchOptions) ([]internal.Op, error) {
	// Convert Operation structs to maps manually to handle special float values
	operationMaps := make([]map[string]interface{}, len(operations))

	for i, op := range operations {
		opMap := make(map[string]interface{})

		// Always include op and path
		opMap["op"] = op.Op
		opMap["path"] = op.Path

		// Handle Value field - include for operations that require it
		// For add/replace operations, even nil is a valid value
		if op.Value != nil || op.Op == "add" || op.Op == "replace" || op.Op == "test" {
			opMap["value"] = op.Value
		}
		if op.From != "" {
			opMap["from"] = op.From
		}

		// Handle Inc field specially to support NaN/Inf values
		// Inc has no omitempty tag, so we include it for all operations
		opMap["inc"] = op.Inc

		// Handle Pos field specially - include for all operations since 0 is a valid position
		// This matches the struct tag change where we removed omitempty from Pos
		opMap["pos"] = float64(op.Pos)
		// Always include str field for test_string operations, even if empty
		opMap["str"] = op.Str
		// Handle Len field specially - include for all operations since 0 is a valid length
		// This matches the struct tag change where we removed omitempty from Len
		opMap["len"] = float64(op.Len)
		if op.Not {
			opMap["not"] = op.Not
		}
		// Handle Type field - could be string or interface{} for test_type operations
		if op.Type != nil {
			// For test_type operations, prefer Type field over Value field
			if op.Op == "test_type" {
				opMap["type"] = op.Type
			} else if typeStr, ok := op.Type.(string); ok && typeStr != "" {
				opMap["type"] = typeStr
			}
		}
		if op.IgnoreCase {
			opMap["ignore_case"] = op.IgnoreCase
		}
		if len(op.Apply) > 0 {
			// Convert nested operations recursively
			nestedOps := make([]interface{}, len(op.Apply))
			for j, nestedOp := range op.Apply {
				nestedOpMap := make(map[string]interface{})
				nestedOpMap["op"] = nestedOp.Op
				nestedOpMap["path"] = nestedOp.Path
				if nestedOp.Value != nil {
					nestedOpMap["value"] = nestedOp.Value
				}
				if nestedOp.From != "" {
					nestedOpMap["from"] = nestedOp.From
				}
				nestedOpMap["inc"] = nestedOp.Inc
				if nestedOp.Pos != 0 {
					nestedOpMap["pos"] = float64(nestedOp.Pos)
				}
				if nestedOp.Str != "" {
					nestedOpMap["str"] = nestedOp.Str
				}
				if nestedOp.Len != 0 {
					nestedOpMap["len"] = float64(nestedOp.Len)
				}
				if nestedOp.Not {
					nestedOpMap["not"] = nestedOp.Not
				}
				if nestedOp.Type != "" {
					nestedOpMap["type"] = nestedOp.Type
				}
				if nestedOp.IgnoreCase {
					nestedOpMap["ignore_case"] = nestedOp.IgnoreCase
				}
				nestedOps[j] = nestedOpMap
			}
			opMap["apply"] = nestedOps
		}
		if len(op.Props) > 0 {
			opMap["props"] = op.Props
		}
		if op.DeleteNull {
			opMap["deleteNull"] = op.DeleteNull
		}
		if op.OldValue != nil {
			opMap["oldValue"] = op.OldValue
		}

		operationMaps[i] = opMap
	}

	// Use existing map-based decoder
	return Decode(operationMaps, options)
}

// DecodeJSON converts JSON bytes to Op instances.
func DecodeJSON(data []byte, options internal.JSONPatchOptions) ([]internal.Op, error) {
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
