// Package json implements a codec for JSON Patch operations
// with full RFC 6902, JSON Predicate, and extended operation support.
package json

import (
	"slices"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
	"github.com/kaptinlin/jsonpointer"
)

// Decode converts JSON operation maps to Op instances.
func Decode(operations []map[string]any, opts internal.JSONPatchOptions) ([]internal.Op, error) {
	ops := make([]internal.Op, 0, len(operations))
	for _, m := range operations {
		o, err := decodeOp(m, opts)
		if err != nil {
			return nil, err
		}
		ops = append(ops, o)
	}
	return ops, nil
}

// DecodeOperations converts Operation structs to Op instances.
func DecodeOperations(operations []internal.Operation, opts internal.JSONPatchOptions) ([]internal.Op, error) {
	maps := make([]map[string]any, len(operations))
	for i, o := range operations {
		maps[i] = operationToMap(o, false)
	}
	return Decode(maps, opts)
}

// DecodeJSON converts JSON bytes to Op instances.
func DecodeJSON(data []byte, opts internal.JSONPatchOptions) ([]internal.Op, error) {
	var operations []map[string]any
	if err := json.Unmarshal(data, &operations); err != nil {
		return nil, err
	}
	return Decode(operations, opts)
}

// parseOpHeader extracts and validates the "op" and "path" fields from m.
func parseOpHeader(m map[string]any) (string, []string, error) {
	opType, ok := m["op"].(string)
	if !ok {
		return "", nil, ErrOpMissingOpField
	}

	pathStr, ok := m["path"].(string)
	if !ok {
		return "", nil, ErrOpMissingPathField
	}

	if err := jsonpointer.Validate(pathStr); err != nil {
		return "", nil, ErrInvalidPointer
	}

	return opType, jsonpointer.Parse(pathStr), nil
}

// decodeOp converts a JSON operation map to an Op instance.
func decodeOp(m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	opType, path, err := parseOpHeader(m)
	if err != nil {
		return nil, err
	}

	switch opType {
	case "add", "remove", "replace", "move", "copy":
		return decodeCoreOp(opType, path, m)
	case "flip", "inc", "str_ins", "str_del", "split", "merge", "extend":
		return decodeExtendedOp(opType, path, m)
	case "not":
		return decodeNotOp(path, m, opts)
	default:
		return decodePredicateOp(opType, path, m, opts)
	}
}

// decodePredicateOnly converts a JSON operation map to a PredicateOp.
func decodePredicateOnly(m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	opType, path, err := parseOpHeader(m)
	if err != nil {
		return nil, err
	}
	return decodePredicateOp(opType, path, m, opts)
}

// decodeCoreOp decodes standard JSON Patch (RFC 6902) operations.
func decodeCoreOp(opType string, path []string, m map[string]any) (internal.Op, error) {
	switch opType {
	case "add":
		if _, ok := m["value"]; !ok {
			return nil, ErrAddOpMissingValue
		}
		return op.NewAdd(path, m["value"]), nil

	case "remove":
		if oldVal, ok := m["oldValue"]; ok {
			return op.NewRemoveWithOldValue(path, oldVal), nil
		}
		return op.NewRemove(path), nil

	case "replace":
		if _, ok := m["value"]; !ok {
			return nil, ErrReplaceOpMissingValue
		}
		if oldVal, ok := m["oldValue"]; ok {
			return op.NewReplaceWithOldValue(path, m["value"], oldVal), nil
		}
		return op.NewReplace(path, m["value"]), nil

	case "move":
		from, ok := m["from"].(string)
		if !ok {
			return nil, ErrMoveOpMissingFrom
		}
		return op.NewMove(path, jsonpointer.Parse(from)), nil

	case "copy":
		from, ok := m["from"].(string)
		if !ok {
			return nil, ErrCopyOpMissingFrom
		}
		return op.NewCopy(path, jsonpointer.Parse(from)), nil

	default:
		return nil, ErrCodecOpUnknown
	}
}

// decodeExtendedOp decodes extended operations.
func decodeExtendedOp(opType string, path []string, m map[string]any) (internal.Op, error) {
	switch opType {
	case "flip":
		return op.NewFlip(path), nil

	case "inc":
		raw, ok := m["inc"]
		if !ok {
			return nil, ErrIncOpMissingInc
		}
		val, ok := op.ToFloat64(raw)
		if !ok {
			return nil, ErrIncOpInvalidType
		}
		return op.NewInc(path, val), nil

	case "str_ins":
		raw, hasPos := m["pos"]
		if !hasPos {
			return nil, ErrStrInsOpMissingPos
		}
		pos, ok := op.ToFloat64(raw)
		if !ok {
			return nil, ErrStrInsOpMissingPos
		}
		str, ok := m["str"].(string)
		if !ok {
			return nil, ErrStrInsOpMissingStr
		}
		return op.NewStrIns(path, pos, str), nil

	case "str_del":
		raw, hasPos := m["pos"]
		if !hasPos {
			return nil, ErrStrDelOpMissingPos
		}
		pos, ok := op.ToFloat64(raw)
		if !ok {
			return nil, ErrStrDelOpMissingPos
		}
		if str, ok := m["str"].(string); ok && str != "" {
			return op.NewStrDelWithStr(path, pos, str), nil
		}
		if lenVal, ok := op.ToFloat64(m["len"]); ok {
			return op.NewStrDel(path, pos, lenVal), nil
		}
		return nil, ErrStrDelOpMissingFields

	case "split":
		raw, hasPos := m["pos"]
		if !hasPos {
			return nil, ErrSplitOpMissingPos
		}
		pos, ok := op.ToFloat64(raw)
		if !ok {
			return nil, ErrSplitOpMissingPos
		}
		return op.NewSplit(path, pos, m["props"]), nil

	case "merge":
		var props map[string]any
		if p, ok := m["props"].(map[string]any); ok {
			props = p
		}
		var pos float64
		if v, ok := op.ToFloat64(m["pos"]); ok {
			pos = v
		}
		return op.NewMerge(path, pos, props), nil

	case "extend":
		props, ok := m["props"].(map[string]any)
		if !ok {
			return nil, ErrValueNotObject
		}
		deleteNull, _ := m["deleteNull"].(bool)
		return op.NewExtend(path, props, deleteNull), nil

	default:
		return nil, ErrCodecOpUnknown
	}
}

// decodePredicateOp decodes JSON Predicate operations including test,
// type checks, string tests, comparisons, and composite operations.
func decodePredicateOp(opType string, path []string, m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	switch opType {
	case "test":
		val, ok := m["value"]
		if !ok {
			return nil, ErrMissingValueField
		}
		t := op.NewTest(path, val)
		if notVal, ok := m["not"].(bool); ok && notVal {
			t.NotFlag = true
		}
		return t, nil

	case "defined":
		return op.NewDefined(path), nil

	case "undefined":
		return op.NewUndefined(path), nil

	case "type":
		val, ok := m["value"].(string)
		if !ok {
			return nil, ErrTypeOpMissingValue
		}
		return op.NewType(path, val), nil

	case "test_type":
		return decodeTestType(path, m)

	case "test_string":
		return decodeTestString(path, m)

	case "test_string_len":
		return decodeTestStringLen(path, m)

	case "contains":
		val, ok := m["value"]
		if !ok {
			return nil, ErrContainsOpMissingValue
		}
		s, ok := val.(string)
		if !ok {
			return nil, op.ErrContainsValueMustBeString
		}
		return op.NewContainsWithIgnoreCase(path, s, boolField(m, "ignore_case")), nil

	case "ends":
		val, ok := m["value"].(string)
		if !ok {
			return nil, ErrEndsOpMissingValue
		}
		return op.NewEndsWithIgnoreCase(path, val, boolField(m, "ignore_case")), nil

	case "starts":
		val, ok := m["value"].(string)
		if !ok {
			return nil, ErrStartsOpMissingValue
		}
		return op.NewStartsWithIgnoreCase(path, val, boolField(m, "ignore_case")), nil

	case "matches":
		val, ok := m["value"].(string)
		if !ok {
			return nil, ErrMatchesOpMissingValue
		}
		return op.NewMatches(path, val, boolField(m, "ignore_case"), opts.CreateMatcher), nil

	case "in":
		val := m["value"]
		if values, ok := val.([]any); ok {
			return op.NewIn(path, values), nil
		}
		return op.NewIn(path, []any{val}), nil

	case "less":
		val, ok := op.ToFloat64(m["value"])
		if !ok {
			return nil, ErrLessOpMissingValue
		}
		return op.NewLess(path, val), nil

	case "more":
		val, ok := op.ToFloat64(m["value"])
		if !ok {
			return nil, ErrMoreOpMissingValue
		}
		return op.NewMore(path, val), nil

	case "and":
		preds, err := decodeApplyField(m, path, ErrAndOpMissingApply, opts)
		if err != nil {
			return nil, err
		}
		return op.NewAnd(path, preds), nil

	case "or":
		preds, err := decodeApplyField(m, path, ErrOrOpMissingApply, opts)
		if err != nil {
			return nil, err
		}
		return op.NewOr(path, preds), nil

	case "not":
		return decodePredicateNot(path, m, opts)

	default:
		return nil, ErrCodecOpUnknown
	}
}

// decodeTestType decodes a test_type operation supporting both single and multiple types.
func decodeTestType(path []string, m map[string]any) (internal.Op, error) {
	// Check "type" field first (standard), then fall back to "value" (compatibility).
	typeField := m["type"]
	if typeField == nil {
		typeField = m["value"]
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
		types := make([]string, len(v))
		for i, t := range v {
			s, ok := t.(string)
			if !ok || !internal.IsValidJSONPatchType(s) {
				return nil, ErrInvalidType
			}
			types[i] = s
		}
		return op.NewTestTypeMultiple(path, types), nil
	case []string:
		if len(v) == 0 {
			return nil, ErrEmptyTypeList
		}
		for _, s := range v {
			if !internal.IsValidJSONPatchType(s) {
				return nil, ErrInvalidType
			}
		}
		return op.NewTestTypeMultiple(path, v), nil
	default:
		return nil, ErrTestTypeOpMissingType
	}
}

// decodeTestString decodes a test_string operation.
func decodeTestString(path []string, m map[string]any) (internal.Op, error) {
	str, ok := m["str"].(string)
	if !ok {
		return nil, ErrTestStringOpMissingStr
	}
	var pos float64
	if v, ok := op.ToFloat64(m["pos"]); ok {
		pos = v
	}
	notFlag, _ := m["not"].(bool)
	ignoreCase, _ := m["ignore_case"].(bool)

	return op.NewTestString(path, str, pos, notFlag, ignoreCase), nil
}

// decodeTestStringLen decodes a test_string_len operation.
func decodeTestStringLen(path []string, m map[string]any) (internal.Op, error) {
	lenVal, ok := op.ToFloat64(m["len"])
	if !ok {
		return nil, ErrTestStringLenOpMissingLen
	}
	not, _ := m["not"].(bool)
	return op.NewTestStringLenWithNot(path, lenVal, not), nil
}

// extractApplyField extracts and validates the "apply" array from m.
func extractApplyField(m map[string]any) ([]any, error) {
	apply, ok := m["apply"].([]any)
	if !ok {
		return nil, ErrNotOpMissingApply
	}
	if len(apply) == 0 {
		return nil, ErrNotOpRequiresOperand
	}
	return apply, nil
}

// decodeApplyField extracts the "apply" array and decodes its sub-predicates.
func decodeApplyField(m map[string]any, path []string, missingErr error, opts internal.JSONPatchOptions) ([]any, error) {
	apply, ok := m["apply"].([]any)
	if !ok {
		return nil, missingErr
	}
	return decodeSubPredicates(apply, path, opts)
}

// decodeNotOp decodes a top-level not operation with multiple predicates.
func decodeNotOp(path []string, m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	apply, err := extractApplyField(m)
	if err != nil {
		return nil, err
	}
	preds, err := decodeSubPredicates(apply, path, opts)
	if err != nil {
		return nil, err
	}
	return op.NewNotMultiple(path, preds), nil
}

// decodePredicateNot decodes a not predicate with a single operand.
func decodePredicateNot(path []string, m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	apply, err := extractApplyField(m)
	if err != nil {
		return nil, err
	}
	sub, ok := apply[0].(map[string]any)
	if !ok {
		return nil, ErrNotOpRequiresValidOperand
	}

	sp, _ := sub["path"].(string)
	merged := mergePaths(path, jsonpointer.Parse(sp))
	sub["path"] = jsonpointer.Format(merged...)

	operand, err := decodePredicateOnly(sub, opts)
	if err != nil {
		return nil, err
	}
	pred, ok := operand.(internal.PredicateOp)
	if !ok {
		return nil, ErrNotOpRequiresValidOperand
	}
	return op.NewNot(pred), nil
}

// decodeSubPredicates decodes an array of sub-operations for and/or/not operations,
// handling path merging and recursive predicate decoding.
func decodeSubPredicates(apply []any, base jsonpointer.Path, opts internal.JSONPatchOptions) ([]any, error) {
	preds := make([]any, 0, len(apply))
	for _, raw := range apply {
		sub, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		sp, _ := sub["path"].(string)
		merged := mergePaths(base, jsonpointer.Parse(sp))
		sub["path"] = jsonpointer.Format(merged...)

		pred, err := decodePredicateOnly(sub, opts)
		if err != nil {
			return nil, err
		}
		preds = append(preds, pred)
	}
	return preds, nil
}

// boolField extracts a boolean field from m, defaulting to false.
func boolField(m map[string]any, field string) bool {
	v, _ := m[field].(bool)
	return v
}

// mergePaths combines base and sub paths for composite operations.
// Returns base if sub is empty, sub if equal, otherwise concatenates them.
func mergePaths(base, sub jsonpointer.Path) jsonpointer.Path {
	if len(sub) == 0 {
		return base
	}
	if slices.Equal(base, sub) {
		return sub
	}
	result := make(jsonpointer.Path, 0, len(base)+len(sub))
	return append(append(result, base...), sub...)
}

// operationToMap converts an Operation struct to a map for decoding.
// When nested is true, zero-value fields are omitted and top-level-only
// fields (Apply, Props, DeleteNull, OldValue) are excluded.
func operationToMap(o internal.Operation, nested bool) map[string]any {
	m := make(map[string]any, 8)
	m["op"] = o.Op
	m["path"] = o.Path

	setValueField(m, o, nested)
	if o.From != "" {
		m["from"] = o.From
	}
	setNumericFields(m, o, nested)
	if o.Not {
		m["not"] = true
	}
	setTypeField(m, o, nested)
	if o.IgnoreCase {
		m["ignore_case"] = true
	}
	if !nested {
		setTopLevelFields(m, o)
	}
	return m
}

// setValueField sets the "value" field based on operation type and nesting.
func setValueField(m map[string]any, o internal.Operation, nested bool) {
	if nested {
		if o.Value != nil {
			m["value"] = o.Value
		}
		return
	}
	if o.Value != nil || o.Op == "add" || o.Op == "replace" || o.Op == "test" {
		m["value"] = o.Value
	}
}

// setNumericFields sets inc, pos, str, and len fields for extended operations.
// In nested mode, zero-value fields are omitted.
func setNumericFields(m map[string]any, o internal.Operation, nested bool) {
	m["inc"] = o.Inc
	if nested {
		if o.Pos != 0 {
			m["pos"] = float64(o.Pos)
		}
		if o.Str != "" {
			m["str"] = o.Str
		}
		if o.Len != 0 {
			m["len"] = float64(o.Len)
		}
		return
	}
	m["pos"] = float64(o.Pos)
	m["str"] = o.Str
	m["len"] = float64(o.Len)
}

// setTypeField sets the "type" field based on operation type and nesting.
func setTypeField(m map[string]any, o internal.Operation, nested bool) {
	if nested {
		if o.Type != nil {
			m["type"] = o.Type
		}
		return
	}
	if o.Type == nil {
		return
	}
	if o.Op == "test_type" {
		m["type"] = o.Type
		return
	}
	if s, ok := o.Type.(string); ok && s != "" {
		m["type"] = s
	}
}

// setTopLevelFields sets fields only present at the top level (not nested).
func setTopLevelFields(m map[string]any, o internal.Operation) {
	if len(o.Apply) > 0 {
		subs := make([]any, len(o.Apply))
		for i, sub := range o.Apply {
			subs[i] = operationToMap(sub, true)
		}
		m["apply"] = subs
	}
	if len(o.Props) > 0 {
		m["props"] = o.Props
	}
	if o.DeleteNull {
		m["deleteNull"] = true
	}
	if o.OldValue != nil {
		m["oldValue"] = o.OldValue
	}
}
