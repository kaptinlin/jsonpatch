// Package json implements a codec for JSON Patch operations
// with full RFC 6902, JSON Predicate, and extended operation support.
package json

import (
	"fmt"
	"slices"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/jsonpointer"

	"github.com/kaptinlin/jsonpatch/internal"
	"github.com/kaptinlin/jsonpatch/op"
)

// Decode converts JSON operation maps to Op instances.
func Decode(operations []map[string]any, opts internal.JSONPatchOptions) ([]internal.Op, error) {
	ops := make([]internal.Op, len(operations))
	for i, m := range operations {
		o, err := decodeOp(m, opts)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", decodeOperationContext(i, m), err)
		}
		ops[i] = o
	}
	return ops, nil
}

// DecodeOperations converts Operation structs to Op instances.
func DecodeOperations(operations []internal.Operation, opts internal.JSONPatchOptions) ([]internal.Op, error) {
	maps := make([]map[string]any, len(operations))
	for i := range operations {
		maps[i] = operationToMap(&operations[i], false)
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

func decodeOperationContext(index int, m map[string]any) string {
	opType, hasOp := m["op"].(string)
	path, hasPath := m["path"].(string)

	switch {
	case hasOp && hasPath:
		return fmt.Sprintf("operation %d (%s %q)", index, opType, path)
	case hasOp:
		return fmt.Sprintf("operation %d (%s)", index, opType)
	case hasPath:
		return fmt.Sprintf("operation %d (%q)", index, path)
	default:
		return fmt.Sprintf("operation %d", index)
	}
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
		from, err := requiredPointer(m, "from", ErrMoveOpMissingFrom)
		if err != nil {
			return nil, err
		}
		return op.NewMove(path, from), nil

	case "copy":
		from, err := requiredPointer(m, "from", ErrCopyOpMissingFrom)
		if err != nil {
			return nil, err
		}
		return op.NewCopy(path, from), nil

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
		pos, ok := requiredFloat64(m, "pos")
		if !ok {
			return nil, ErrStrInsOpMissingPos
		}
		str, ok := m["str"].(string)
		if !ok {
			return nil, ErrStrInsOpMissingStr
		}
		return op.NewStrIns(path, pos, str), nil

	case "str_del":
		pos, ok := requiredFloat64(m, "pos")
		if !ok {
			return nil, ErrStrDelOpMissingPos
		}
		if str, ok := m["str"].(string); ok {
			return op.NewStrDelWithStr(path, pos, str), nil
		}
		if lenVal, ok := op.ToFloat64(m["len"]); ok {
			return op.NewStrDel(path, pos, lenVal), nil
		}
		return nil, ErrStrDelOpMissingFields

	case "split":
		pos, ok := requiredFloat64(m, "pos")
		if !ok {
			return nil, ErrSplitOpMissingPos
		}
		return op.NewSplit(path, pos, m["props"]), nil

	case "merge":
		pos, ok := requiredFloat64(m, "pos")
		if !ok {
			return nil, ErrMergeOpMissingPos
		}
		var props map[string]any
		if raw, ok := m["props"]; ok {
			p, ok := raw.(map[string]any)
			if !ok {
				return nil, ErrValueNotObject
			}
			props = p
		}
		return op.NewMerge(path, pos, props), nil

	case "extend":
		props, ok := m["props"].(map[string]any)
		if !ok {
			return nil, ErrValueNotObject
		}
		deleteNull, err := optionalBool(m, "deleteNull")
		if err != nil {
			return nil, err
		}
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
		notVal, err := optionalBool(m, "not")
		if err != nil {
			return nil, err
		}
		if notVal {
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
		val, err := requiredString(m, "value", ErrContainsOpMissingValue, op.ErrContainsValueMustBeString)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := optionalBool(m, "ignore_case")
		if err != nil {
			return nil, err
		}
		return op.NewContainsWithIgnoreCase(path, val, ignoreCase), nil

	case "ends":
		val, err := requiredString(m, "value", ErrEndsOpMissingValue, ErrEndsOpMissingValue)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := optionalBool(m, "ignore_case")
		if err != nil {
			return nil, err
		}
		return op.NewEndsWithIgnoreCase(path, val, ignoreCase), nil

	case "starts":
		val, err := requiredString(m, "value", ErrStartsOpMissingValue, ErrStartsOpMissingValue)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := optionalBool(m, "ignore_case")
		if err != nil {
			return nil, err
		}
		return op.NewStartsWithIgnoreCase(path, val, ignoreCase), nil

	case "matches":
		val, err := requiredString(m, "value", ErrMatchesOpMissingValue, ErrMatchesOpMissingValue)
		if err != nil {
			return nil, err
		}
		ignoreCase, err := optionalBool(m, "ignore_case")
		if err != nil {
			return nil, err
		}
		return op.NewMatches(path, val, ignoreCase, opts.CreateMatcher), nil

	case "in":
		values, ok := m["value"].([]any)
		if !ok {
			return nil, ErrInOpValueMustBeArray
		}
		return op.NewIn(path, values), nil

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

// decodeTestType decodes a test_type operation supporting one or more type names.
func decodeTestType(path []string, m map[string]any) (internal.Op, error) {
	typeField := m["type"]

	switch v := typeField.(type) {
	case string:
		if !internal.IsValidJSONPatchType(v) {
			return nil, ErrInvalidType
		}
		return op.NewTestType(path, v), nil
	case []any:
		return decodeTestTypeArray(path, v)
	case []string:
		return newTestTypeMultiple(path, v)
	default:
		return nil, ErrTestTypeOpMissingType
	}
}

// decodeTestTypeArray decodes a test_type operation with []any type list.
func decodeTestTypeArray(path []string, v []any) (internal.Op, error) {
	types := make([]string, len(v))
	for i, t := range v {
		s, ok := t.(string)
		if !ok {
			return nil, ErrInvalidType
		}
		types[i] = s
	}
	return newTestTypeMultiple(path, types)
}

func newTestTypeMultiple(path []string, types []string) (internal.Op, error) {
	if len(types) == 0 {
		return nil, ErrEmptyTypeList
	}
	if slices.ContainsFunc(types, func(s string) bool {
		return !internal.IsValidJSONPatchType(s)
	}) {
		return nil, ErrInvalidType
	}
	return op.NewTestTypeMultiple(path, types), nil
}

// decodeTestString decodes a test_string operation.
func decodeTestString(path []string, m map[string]any) (internal.Op, error) {
	str, ok := m["str"].(string)
	if !ok {
		return nil, ErrTestStringOpMissingStr
	}
	pos, ok := requiredFloat64(m, "pos")
	if !ok {
		return nil, ErrTestStringOpMissingPos
	}
	notFlag, err := optionalBool(m, "not")
	if err != nil {
		return nil, err
	}
	ignoreCase, err := optionalBool(m, "ignore_case")
	if err != nil {
		return nil, err
	}

	return op.NewTestString(path, str, pos, notFlag, ignoreCase), nil
}

// decodeTestStringLen decodes a test_string_len operation.
func decodeTestStringLen(path []string, m map[string]any) (internal.Op, error) {
	lenVal, ok := op.ToFloat64(m["len"])
	if !ok {
		return nil, ErrTestStringLenOpMissingLen
	}
	not, err := optionalBool(m, "not")
	if err != nil {
		return nil, err
	}
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
	if len(apply) != 1 {
		return nil, ErrNotOpRequiresSingleOperand
	}
	pred, err := decodeSubPredicate(apply[0], path, opts)
	if err != nil {
		return nil, err
	}
	return op.NewNotMultiple(path, []any{pred}), nil
}

// decodePredicateNot decodes a not predicate with a single operand.
func decodePredicateNot(path []string, m map[string]any, opts internal.JSONPatchOptions) (internal.Op, error) {
	apply, err := extractApplyField(m)
	if err != nil {
		return nil, err
	}
	if len(apply) != 1 {
		return nil, ErrNotOpRequiresSingleOperand
	}
	pred, err := decodeSubPredicate(apply[0], path, opts)
	if err != nil {
		return nil, err
	}
	return op.NewNotMultiple(path, []any{pred}), nil
}

// decodeSubPredicates decodes an array of sub-operations for and/or/not operations,
// handling path merging and recursive predicate decoding.
func decodeSubPredicates(apply []any, base jsonpointer.Path, opts internal.JSONPatchOptions) ([]any, error) {
	preds := make([]any, 0, len(apply))
	for _, raw := range apply {
		pred, err := decodeSubPredicate(raw, base, opts)
		if err != nil {
			return nil, err
		}
		preds = append(preds, pred)
	}
	return preds, nil
}

func decodeSubPredicate(raw any, base jsonpointer.Path, opts internal.JSONPatchOptions) (internal.PredicateOp, error) {
	sub, ok := raw.(map[string]any)
	if !ok {
		return nil, ErrInvalidPredicateOperand
	}

	child := cloneOperationMap(sub)
	pathStr, ok := child["path"].(string)
	if !ok {
		return nil, ErrOpMissingPathField
	}
	if err := jsonpointer.Validate(pathStr); err != nil {
		return nil, ErrInvalidPointer
	}
	merged := mergePaths(base, jsonpointer.Parse(pathStr))
	child["path"] = jsonpointer.Format(merged...)

	decoded, err := decodePredicateOnly(child, opts)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidPredicateOperand, err)
	}
	pred, ok := decoded.(internal.PredicateOp)
	if !ok {
		return nil, ErrInvalidPredicateOperand
	}
	return pred, nil
}

func cloneOperationMap(m map[string]any) map[string]any {
	clone := make(map[string]any, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func optionalBool(m map[string]any, field string) (bool, error) {
	raw, ok := m[field]
	if !ok {
		return false, nil
	}
	value, ok := raw.(bool)
	if !ok {
		return false, ErrInvalidBooleanField
	}
	return value, nil
}

func requiredFloat64(m map[string]any, field string) (float64, bool) {
	raw, ok := m[field]
	if !ok {
		return 0, false
	}
	return op.ToFloat64(raw)
}

func requiredString(m map[string]any, field string, missingErr, typeErr error) (string, error) {
	raw, ok := m[field]
	if !ok {
		return "", missingErr
	}
	value, ok := raw.(string)
	if !ok {
		return "", typeErr
	}
	return value, nil
}

func requiredPointer(m map[string]any, field string, missingErr error) ([]string, error) {
	raw, ok := m[field]
	if !ok {
		return nil, missingErr
	}
	value, ok := raw.(string)
	if !ok {
		return nil, missingErr
	}
	if err := jsonpointer.Validate(value); err != nil {
		return nil, ErrInvalidPointer
	}
	return jsonpointer.Parse(value), nil
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
	return slices.Concat(base, sub)
}

// operationToMap converts an Operation struct to a map for decoding.
// When nested is true, zero-value fields are omitted and top-level-only
// fields (Apply, Props, DeleteNull, OldValue) are excluded.
func operationToMap(o *internal.Operation, nested bool) map[string]any {
	m := make(map[string]any, 8)
	m["op"] = o.Op
	m["path"] = o.Path

	setValueField(m, o, nested)
	if o.Op == "move" || o.Op == "copy" || o.From != "" {
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
	setApplyField(m, o)
	if !nested {
		setTopLevelFields(m, o)
	}
	return m
}

// setValueField sets the "value" field based on operation type and nesting.
func setValueField(m map[string]any, o *internal.Operation, nested bool) {
	if nested && o.Value == nil {
		return
	}
	if nested || o.Value != nil || o.Op == "add" || o.Op == "replace" || o.Op == "test" {
		m["value"] = o.Value
	}
}

// setNumericFields sets inc, pos, str, and len fields for extended operations.
// In nested mode, zero-value fields are omitted.
func setNumericFields(m map[string]any, o *internal.Operation, nested bool) {
	switch o.Op {
	case "inc":
		m["inc"] = o.Inc
	case "str_ins", "test_string":
		m["pos"] = float64(o.Pos)
		m["str"] = o.Str
	case "str_del":
		m["pos"] = float64(o.Pos)
		if o.Str != "" {
			m["str"] = o.Str
		} else {
			m["len"] = float64(o.Len)
		}
	case "split", "merge":
		m["pos"] = float64(o.Pos)
	case "test_string_len":
		m["len"] = float64(o.Len)
	}
}

// setTypeField sets the "type" field based on operation type and nesting.
func setTypeField(m map[string]any, o *internal.Operation, nested bool) {
	if o.Type == nil {
		return
	}
	if nested || o.Op == "test_type" {
		m["type"] = o.Type
		return
	}
	if s, ok := o.Type.(string); ok && s != "" {
		m["type"] = s
	}
}

func setApplyField(m map[string]any, o *internal.Operation) {
	if len(o.Apply) > 0 {
		subs := make([]any, len(o.Apply))
		for i := range o.Apply {
			subs[i] = operationToMap(&o.Apply[i], true)
		}
		m["apply"] = subs
	}
}

// setTopLevelFields sets fields only present at the top level (not nested).
func setTopLevelFields(m map[string]any, o *internal.Operation) {
	if o.Props != nil {
		m["props"] = o.Props
	}
	if o.DeleteNull {
		m["deleteNull"] = true
	}
	if o.OldValue != nil {
		m["oldValue"] = o.OldValue
	}
}
