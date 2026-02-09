// Package json implements JSON codec for JSON Patch operations.
// It provides encoding and decoding for JSON Patch operations with full RFC 6902 support.
package json

import (
	"slices"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Decode converts JSON operations to Op instances.
func Decode(operations []map[string]any, options internal.JSONPatchOptions) ([]internal.Op, error) {
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

// DecodeOperations converts Operation structs to Op instances.
func DecodeOperations(operations []internal.Operation, options internal.JSONPatchOptions) ([]internal.Op, error) {
	operationMaps := make([]map[string]any, len(operations))
	for i, o := range operations {
		operationMaps[i] = toMap(o)
	}
	return Decode(operationMaps, options)
}

// DecodeJSON converts JSON bytes to Op instances.
func DecodeJSON(data []byte, options internal.JSONPatchOptions) ([]internal.Op, error) {
	var operations []map[string]any
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, err
	}
	return Decode(operations, options)
}

// parseOpHeader extracts and validates the common "op" and "path" fields
// from an operation map.
func parseOpHeader(operation map[string]any) (opType, pathStr string, path []string, err error) {
	opType, ok := operation["op"].(string)
	if !ok {
		return "", "", nil, ErrOpMissingOpField
	}

	pathStr, ok = operation["path"].(string)
	if !ok {
		return "", "", nil, ErrOpMissingPathField
	}

	if err := jsonpointer.Validate(pathStr); err != nil {
		return "", "", nil, ErrInvalidPointer
	}

	return opType, pathStr, jsonpointer.Parse(pathStr), nil
}

// OperationToOp converts a JSON operation map to an Op instance.
func OperationToOp(operation map[string]any, options internal.JSONPatchOptions) (internal.Op, error) {
	opType, pathStr, path, err := parseOpHeader(operation)
	if err != nil {
		return nil, err
	}

	switch opType {
	case "add", "remove", "replace", "move", "copy":
		return parseCoreOp(opType, path, operation)
	case "flip", "inc", "str_ins", "str_del", "split", "merge", "extend":
		return parseExtendedOp(opType, path, operation)
	case "not":
		return parseNotOp(path, pathStr, operation, options)
	default:
		return parsePredicateOp(opType, path, pathStr, operation, options)
	}
}

// OperationToPredicateOp converts a JSON operation map to a PredicateOp instance.
func OperationToPredicateOp(operation map[string]any, options internal.JSONPatchOptions) (internal.Op, error) {
	opType, pathStr, path, err := parseOpHeader(operation)
	if err != nil {
		return nil, err
	}
	return parsePredicateOp(opType, path, pathStr, operation, options)
}

// parseCoreOp decodes standard JSON Patch (RFC 6902) operations.
func parseCoreOp(opType string, path []string, operation map[string]any) (internal.Op, error) {
	switch opType {
	case "add":
		if _, ok := operation["value"]; !ok {
			return nil, ErrAddOpMissingValue
		}
		return op.NewAdd(path, operation["value"]), nil

	case "remove":
		if oldValue, ok := operation["oldValue"]; ok {
			return op.NewRemoveWithOldValue(path, oldValue), nil
		}
		return op.NewRemove(path), nil

	case "replace":
		if _, ok := operation["value"]; !ok {
			return nil, ErrReplaceOpMissingValue
		}
		if oldValue, ok := operation["oldValue"]; ok {
			return op.NewReplaceWithOldValue(path, operation["value"], oldValue), nil
		}
		return op.NewReplace(path, operation["value"]), nil

	case "move":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrMoveOpMissingFrom
		}
		return op.NewMove(path, jsonpointer.Parse(fromStr)), nil

	case "copy":
		fromStr, ok := operation["from"].(string)
		if !ok {
			return nil, ErrCopyOpMissingFrom
		}
		return op.NewCopy(path, jsonpointer.Parse(fromStr)), nil

	default:
		return nil, ErrCodecOpUnknown
	}
}

// parseExtendedOp decodes extended operations (flip, inc, str_ins, str_del, split, merge, extend).
func parseExtendedOp(opType string, path []string, operation map[string]any) (internal.Op, error) {
	switch opType {
	case "flip":
		return op.NewFlip(path), nil

	case "inc":
		incField, hasInc := operation["inc"]
		if !hasInc {
			return nil, ErrIncOpMissingInc
		}
		incVal, ok := op.ToFloat64(incField)
		if !ok {
			return nil, ErrIncOpInvalidType
		}
		return op.NewInc(path, incVal), nil

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
		return op.NewStrIns(path, pos, str), nil

	case "str_del":
		posVal, hasPosField := operation["pos"]
		if !hasPosField {
			return nil, ErrStrDelOpMissingPos
		}
		pos, ok := op.ToFloat64(posVal)
		if !ok {
			return nil, ErrStrDelOpMissingPos
		}
		if str, ok := operation["str"].(string); ok && str != "" {
			return op.NewStrDelWithStr(path, pos, str), nil
		}
		if lenVal, ok := op.ToFloat64(operation["len"]); ok {
			return op.NewStrDel(path, pos, lenVal), nil
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
		return op.NewSplit(path, pos, operation["props"]), nil

	case "merge":
		var props map[string]any
		if p, ok := operation["props"].(map[string]any); ok {
			props = p
		}
		var pos float64
		if posVal, ok := op.ToFloat64(operation["pos"]); ok {
			pos = posVal
		}
		return op.NewMerge(path, pos, props), nil

	case "extend":
		props, ok := operation["props"].(map[string]any)
		if !ok {
			return nil, ErrValueNotObject
		}
		deleteNull, _ := operation["deleteNull"].(bool)
		return op.NewExtend(path, props, deleteNull), nil

	default:
		return nil, ErrCodecOpUnknown
	}
}

// parsePredicateOp decodes JSON Predicate operations including test, type checks,
// string tests, comparisons, and composite operations (and, or, not).
func parsePredicateOp(opType string, path []string, pathStr string, operation map[string]any, options internal.JSONPatchOptions) (internal.Op, error) {
	switch opType {
	case "test":
		value, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrMissingValueField
		}
		testOp := op.NewTest(path, value)
		if notVal, ok := operation["not"].(bool); ok && notVal {
			testOp.NotFlag = true
		}
		return testOp, nil

	case "defined":
		return op.NewDefined(path), nil

	case "undefined":
		return op.NewUndefined(path), nil

	case "type":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrTypeOpMissingValue
		}
		return op.NewType(path, value), nil

	case "test_type":
		return parseTestType(path, operation)

	case "test_string":
		return parseTestString(path, operation)

	case "test_string_len":
		return parseTestStringLen(path, operation)

	case "contains":
		value, hasValue := operation["value"]
		if !hasValue {
			return nil, ErrContainsOpMissingValue
		}
		stringValue, ok := value.(string)
		if !ok {
			return nil, op.ErrContainsValueMustBeString
		}
		return op.NewContainsWithIgnoreCase(path, stringValue, boolField(operation, "ignore_case")), nil

	case "ends":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrEndsOpMissingValue
		}
		return op.NewEndsWithIgnoreCase(path, value, boolField(operation, "ignore_case")), nil

	case "starts":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrStartsOpMissingValue
		}
		return op.NewStartsWithIgnoreCase(path, value, boolField(operation, "ignore_case")), nil

	case "matches":
		value, ok := operation["value"].(string)
		if !ok {
			return nil, ErrMatchesOpMissingValue
		}
		return op.NewMatches(path, value, boolField(operation, "ignore_case"), options.CreateMatcher), nil

	case "in":
		value := operation["value"]
		if values, ok := value.([]any); ok {
			return op.NewIn(path, values), nil
		}
		return op.NewIn(path, []any{value}), nil

	case "less":
		value, ok := op.ToFloat64(operation["value"])
		if !ok {
			return nil, ErrLessOpMissingValue
		}
		return op.NewLess(path, value), nil

	case "more":
		value, ok := op.ToFloat64(operation["value"])
		if !ok {
			return nil, ErrMoreOpMissingValue
		}
		return op.NewMore(path, value), nil

	case "and":
		apply, ok := operation["apply"].([]any)
		if !ok {
			return nil, ErrAndOpMissingApply
		}
		predicateOps, err := decodeCompositePredicates(apply, jsonpointer.Parse(pathStr), options)
		if err != nil {
			return nil, err
		}
		return op.NewAnd(path, predicateOps), nil

	case "or":
		apply, ok := operation["apply"].([]any)
		if !ok {
			return nil, ErrOrOpMissingApply
		}
		predicateOps, err := decodeCompositePredicates(apply, jsonpointer.Parse(pathStr), options)
		if err != nil {
			return nil, err
		}
		return op.NewOr(path, predicateOps), nil

	case "not":
		return parsePredicateNot(pathStr, operation, options)

	default:
		return nil, ErrCodecOpUnknown
	}
}

// parseTestType decodes a test_type operation supporting both single and multiple types.
func parseTestType(path []string, operation map[string]any) (internal.Op, error) {
	// Check "type" field first (standard), then fall back to "value" field (compatibility).
	typeField := operation["type"]
	if typeField == nil {
		typeField = operation["value"]
	}

	switch v := typeField.(type) {
	case string:
		if !internal.IsValidJSONPatchType(v) {
			return nil, ErrInvalidType
		}
		return op.NewTestType(path, v), nil
	case []any:
		if len(v) == 0 {
			return nil, ErrEmptyTypeList
		}
		typeStrings := make([]string, len(v))
		for i, t := range v {
			typeStr, ok := t.(string)
			if !ok {
				return nil, ErrInvalidType
			}
			if !internal.IsValidJSONPatchType(typeStr) {
				return nil, ErrInvalidType
			}
			typeStrings[i] = typeStr
		}
		return op.NewTestTypeMultiple(path, typeStrings), nil
	case []string:
		if len(v) == 0 {
			return nil, ErrEmptyTypeList
		}
		for _, typeStr := range v {
			if !internal.IsValidJSONPatchType(typeStr) {
				return nil, ErrInvalidType
			}
		}
		return op.NewTestTypeMultiple(path, v), nil
	default:
		return nil, ErrTestTypeOpMissingType
	}
}

// parseTestString decodes a test_string operation.
func parseTestString(path []string, operation map[string]any) (internal.Op, error) {
	str, ok := operation["str"].(string)
	if !ok {
		return nil, ErrTestStringOpMissingStr
	}
	var pos float64
	if posVal, ok := op.ToFloat64(operation["pos"]); ok {
		pos = posVal
	}
	notFlag, _ := operation["not"].(bool)
	ignoreCase, _ := operation["ignore_case"].(bool)

	if pos != 0 || notFlag || ignoreCase {
		return op.NewTestStringWithIgnoreCase(path, str, pos, notFlag, ignoreCase), nil
	}
	return op.NewTestString(path, str), nil
}

// parseTestStringLen decodes a test_string_len operation.
func parseTestStringLen(path []string, operation map[string]any) (internal.Op, error) {
	lenVal, ok := op.ToFloat64(operation["len"])
	if !ok {
		return nil, ErrTestStringLenOpMissingLen
	}
	not, _ := operation["not"].(bool)
	if not {
		return op.NewTestStringLenWithNot(path, lenVal, not), nil
	}
	return op.NewTestStringLen(path, lenVal), nil
}

// parseNotOp decodes a top-level not operation with multiple predicates.
func parseNotOp(path []string, pathStr string, operation map[string]any, options internal.JSONPatchOptions) (internal.Op, error) {
	apply, ok := operation["apply"].([]any)
	if !ok {
		return nil, ErrNotOpMissingApply
	}
	if len(apply) == 0 {
		return nil, ErrNotOpRequiresOperand
	}
	predicateOps, err := decodeCompositePredicates(apply, jsonpointer.Parse(pathStr), options)
	if err != nil {
		return nil, err
	}
	return op.NewNotMultiple(path, predicateOps), nil
}

// parsePredicateNot decodes a not predicate with a single operand.
func parsePredicateNot(pathStr string, operation map[string]any, options internal.JSONPatchOptions) (internal.Op, error) {
	apply, ok := operation["apply"].([]any)
	if !ok {
		return nil, ErrNotOpMissingApply
	}
	if len(apply) == 0 {
		return nil, ErrNotOpRequiresOperand
	}
	applyMap, ok := apply[0].(map[string]any)
	if !ok {
		return nil, ErrNotOpRequiresValidOperand
	}

	subPath := ""
	if sp, ok := applyMap["path"].(string); ok {
		subPath = sp
	}
	mergedPath := mergePaths(jsonpointer.Parse(pathStr), jsonpointer.Parse(subPath))
	applyMap["path"] = jsonpointer.Format(mergedPath...)

	operand, err := OperationToPredicateOp(applyMap, options)
	if err != nil {
		return nil, err
	}
	predicateOp, ok := operand.(internal.PredicateOp)
	if !ok {
		return nil, ErrNotOpRequiresValidOperand
	}
	return op.NewNot(predicateOp), nil
}

// decodeCompositePredicates decodes an array of sub-operations for and/or/not operations.
// It handles path merging and recursive predicate decoding.
func decodeCompositePredicates(apply []any, basePath jsonpointer.Path, options internal.JSONPatchOptions) ([]any, error) {
	predicateOps := make([]any, 0, len(apply))
	for _, subOp := range apply {
		subOpMap, ok := subOp.(map[string]any)
		if !ok {
			continue
		}
		subPath := ""
		if sp, ok := subOpMap["path"].(string); ok {
			subPath = sp
		}
		mergedPath := mergePaths(basePath, jsonpointer.Parse(subPath))
		subOpMap["path"] = jsonpointer.Format(mergedPath...)

		predicateOp, err := OperationToPredicateOp(subOpMap, options)
		if err != nil {
			return nil, err
		}
		predicateOps = append(predicateOps, predicateOp)
	}
	return predicateOps, nil
}

// boolField extracts a boolean field from an operation map with a default of false.
func boolField(operation map[string]any, field string) bool {
	v, _ := operation[field].(bool)
	return v
}

// mergePaths merges two paths for composite operations.
// If subPath is empty, use basePath. If paths are identical, use subPath.
// Otherwise, concatenate them.
func mergePaths(basePath, subPath jsonpointer.Path) jsonpointer.Path {
	if len(subPath) == 0 {
		return basePath
	}
	if slices.Equal(basePath, subPath) {
		return subPath
	}

	result := make(jsonpointer.Path, 0, len(basePath)+len(subPath))
	result = append(result, basePath...)
	result = append(result, subPath...)
	return result
}

// toMap converts an internal.Operation struct to a map for decoding.
func toMap(o internal.Operation) map[string]any {
	m := make(map[string]any, 8)

	m["op"] = o.Op
	m["path"] = o.Path

	// Value field - include for operations that require it (even nil is valid).
	if o.Value != nil || o.Op == "add" || o.Op == "replace" || o.Op == "test" {
		m["value"] = o.Value
	}
	if o.From != "" {
		m["from"] = o.From
	}

	// Numeric fields without omitempty: 0 is a valid value.
	m["inc"] = o.Inc
	m["pos"] = float64(o.Pos)
	m["str"] = o.Str
	m["len"] = float64(o.Len)

	if o.Not {
		m["not"] = o.Not
	}
	if o.Type != nil {
		if o.Op == "test_type" {
			m["type"] = o.Type
		} else if typeStr, ok := o.Type.(string); ok && typeStr != "" {
			m["type"] = typeStr
		}
	}
	if o.IgnoreCase {
		m["ignore_case"] = o.IgnoreCase
	}
	if len(o.Apply) > 0 {
		nestedOps := make([]any, len(o.Apply))
		for j, nested := range o.Apply {
			nestedOps[j] = nestedToMap(nested)
		}
		m["apply"] = nestedOps
	}
	if len(o.Props) > 0 {
		m["props"] = o.Props
	}
	if o.DeleteNull {
		m["deleteNull"] = o.DeleteNull
	}
	if o.OldValue != nil {
		m["oldValue"] = o.OldValue
	}

	return m
}

// nestedToMap converts a nested Operation to a map.
// Nested operations use conditional inclusion for pos/str/len fields.
func nestedToMap(o internal.Operation) map[string]any {
	m := make(map[string]any, 8)

	m["op"] = o.Op
	m["path"] = o.Path

	if o.Value != nil {
		m["value"] = o.Value
	}
	if o.From != "" {
		m["from"] = o.From
	}
	m["inc"] = o.Inc
	if o.Pos != 0 {
		m["pos"] = float64(o.Pos)
	}
	if o.Str != "" {
		m["str"] = o.Str
	}
	if o.Len != 0 {
		m["len"] = float64(o.Len)
	}
	if o.Not {
		m["not"] = o.Not
	}
	if o.Type != "" {
		m["type"] = o.Type
	}
	if o.IgnoreCase {
		m["ignore_case"] = o.IgnoreCase
	}

	return m
}
