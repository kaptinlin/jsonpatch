package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kaptinlin/jsonpointer"
)

// Base validation errors - define clearly and concisely
var (
	ErrNotArray             = errors.New("not an array")
	ErrEmptyPatch           = errors.New("empty operation patch")
	ErrInvalidOperation     = errors.New("invalid operation")
	ErrMissingPath          = errors.New("missing required field 'path'")
	ErrMissingOp            = errors.New("missing required field 'op'")
	ErrMissingValue         = errors.New("missing required field 'value'")
	ErrMissingFrom          = errors.New("missing required field 'from'")
	ErrInvalidPath          = errors.New("field 'path' must be a string")
	ErrInvalidOp            = errors.New("field 'op' must be a string")
	ErrInvalidFrom          = errors.New("field 'from' must be a string")
	ErrInvalidJSONPointer   = errors.New("invalid JSON pointer")
	ErrInvalidOldValue      = errors.New("invalid oldValue")
	ErrCannotMoveToChildren = errors.New("cannot move into own children")
	ErrInvalidIncValue      = errors.New("invalid inc value")
	ErrExpectedStringField  = errors.New("expected string field")
	ErrExpectedBooleanField = errors.New("expected field to be boolean")
	ErrExpectedIntegerField = errors.New("not an integer")
	ErrNegativeNumber       = errors.New("number is negative")
	ErrInvalidProps         = errors.New("invalid props field")
	ErrInvalidTypeField     = errors.New("invalid type field")
	ErrEmptyTypeList        = errors.New("empty type list")
	ErrInvalidType          = errors.New("invalid type")
	ErrValueMustBeString    = errors.New("value must be a string")
	ErrValueMustBeNumber    = errors.New("value must be a number")
	ErrValueMustBeArray     = errors.New("value must be an array")
	ErrValueTooLong         = errors.New("value too long")
	ErrInvalidNotModifier   = errors.New("invalid not modifier")
	ErrMatchesNotAllowed    = errors.New("matches operation not allowed")
	ErrMustBeArray          = errors.New("must be an array")
	ErrEmptyPredicateList   = errors.New("predicate list is empty")
	ErrEitherStrOrLen       = errors.New("either str or len must be set")
	ErrPosGreaterThanZero   = errors.New("expected pos field to be greater than 0")

	// Additional static errors for err113 compliance
	ErrInOperationValueMustBeArray = errors.New("in operation value must be an array")
	ErrExpectedValueToBeString     = errors.New("expected value to be string")
	ErrExpectedIgnoreCaseBoolean   = errors.New("expected ignore_case to be boolean")
	ErrExpectedFieldString         = errors.New("expected field to be string")
)

// ValidateOperations validates an array of JSON Patch operations.
func ValidateOperations(ops []Operation, allowMatchesOp bool) error {
	if ops == nil {
		return ErrNotArray
	}
	if len(ops) == 0 {
		return ErrEmptyPatch
	}

	for i, operation := range ops {
		if err := ValidateOperation(operation, allowMatchesOp); err != nil {
			return fmt.Errorf("error in operation [index = %d] (%w)", i, err)
		}
	}
	return nil
}

// ValidateOperation validates a single JSON Patch operation.
func ValidateOperation(operation Operation, allowMatchesOp bool) error {
	if operation == nil {
		return ErrInvalidOperation
	}

	opMap := operation

	// Validate path field
	path, pathExists := opMap["path"]
	if !pathExists {
		return ErrMissingPath
	}
	pathStr, isString := path.(string)
	if !isString {
		return ErrInvalidPath
	}
	if err := validateJSONPointer(pathStr); err != nil {
		return ErrInvalidJSONPointer
	}

	// Validate op field
	op, opExists := opMap["op"]
	if !opExists {
		return ErrMissingOp
	}
	opStr, isString := op.(string)
	if !isString {
		return ErrInvalidOp
	}

	// Validate operation by type
	switch opStr {
	case "add":
		return validateOperationAdd(opMap)
	case "remove":
		return validateOperationRemove(opMap)
	case "replace":
		return validateOperationReplace(opMap)
	case "copy":
		return validateOperationCopy(opMap)
	case "move":
		return validateOperationMove(opMap)
	case "flip":
		return nil
	case "inc":
		return validateOperationInc(opMap)
	case "str_ins":
		return validateOperationStrIns(opMap)
	case "str_del":
		return validateOperationStrDel(opMap)
	case "extend":
		return validateOperationExtend(opMap)
	case "merge":
		return validateOperationMerge(opMap)
	case "split":
		return validateOperationSplit(opMap)
	default:
		return validatePredicateOperation(opMap, opStr, allowMatchesOp)
	}
}

// validatePredicateOperation validates predicate operations
func validatePredicateOperation(opMap map[string]interface{}, opStr string, allowMatchesOp bool) error {
	switch opStr {
	case "test":
		return validateOperationTest(opMap)
	case "test_type":
		return validateOperationTestType(opMap)
	case "test_string":
		return validateOperationTestString(opMap)
	case "test_string_len":
		return validateOperationTestStringLen(opMap)
	case "matches":
		if !allowMatchesOp {
			return ErrMatchesNotAllowed
		}
		return validateOperationMatches(opMap)
	case "contains":
		return validateOperationContains(opMap)
	case "ends":
		return validateOperationEnds(opMap)
	case "starts":
		return validateOperationStarts(opMap)
	case "in":
		return validateOperationIn(opMap)
	case "more":
		return validateOperationMore(opMap)
	case "less":
		return validateOperationLess(opMap)
	case "type":
		return validateOperationType(opMap)
	case "defined":
		return nil
	case "undefined":
		return nil
	case "and", "or", "not":
		return validateCompositeOperation(opMap, opStr, allowMatchesOp)
	default:
		return fmt.Errorf("%w: unknown operation '%s'", ErrInvalidOperation, opStr)
	}
}

// Core operation validators
func validateOperationAdd(opMap map[string]interface{}) error {
	if _, exists := opMap["value"]; !exists {
		return ErrMissingValue
	}
	return nil
}

func validateOperationRemove(opMap map[string]interface{}) error {
	if oldValue, exists := opMap["oldValue"]; exists && oldValue == nil {
		return ErrInvalidOldValue
	}
	return nil
}

func validateOperationReplace(opMap map[string]interface{}) error {
	if _, exists := opMap["value"]; !exists {
		return ErrMissingValue
	}
	if oldValue, exists := opMap["oldValue"]; exists && oldValue == nil {
		return ErrInvalidOldValue
	}
	return nil
}

func validateOperationCopy(opMap map[string]interface{}) error {
	from, exists := opMap["from"]
	if !exists {
		return ErrMissingFrom
	}
	fromStr, isString := from.(string)
	if !isString {
		return ErrInvalidFrom
	}
	return validateJSONPointer(fromStr)
}

func validateOperationMove(opMap map[string]interface{}) error {
	from, exists := opMap["from"]
	if !exists {
		return ErrMissingFrom
	}
	fromStr, isString := from.(string)
	if !isString {
		return ErrInvalidFrom
	}
	if err := validateJSONPointer(fromStr); err != nil {
		return err
	}

	pathStr, _ := opMap["path"].(string)
	if strings.HasPrefix(pathStr, fromStr+"/") {
		return ErrCannotMoveToChildren
	}
	return nil
}

func validateOperationTest(opMap map[string]interface{}) error {
	if _, exists := opMap["value"]; !exists {
		return ErrMissingValue
	}
	return validateNot(opMap)
}

// Extended operation validators
func validateOperationInc(opMap map[string]interface{}) error {
	inc, exists := opMap["inc"]
	if !exists {
		return ErrInvalidIncValue
	}
	if !isNumber(inc) {
		return ErrInvalidIncValue
	}
	return nil
}

func validateOperationStrIns(opMap map[string]interface{}) error {
	if err := validateNonNegativeInteger(opMap, "pos"); err != nil {
		return err
	}
	str, exists := opMap["str"]
	if !exists {
		return ErrExpectedStringField
	}
	if _, isString := str.(string); !isString {
		return ErrExpectedStringField
	}
	return nil
}

func validateOperationStrDel(opMap map[string]interface{}) error {
	if err := validateNonNegativeInteger(opMap, "pos"); err != nil {
		return err
	}

	_, hasStr := opMap["str"]
	_, hasLen := opMap["len"]

	if !hasStr && !hasLen {
		return ErrEitherStrOrLen
	}

	if hasStr {
		str := opMap["str"]
		if _, isString := str.(string); !isString {
			return ErrExpectedStringField
		}
	}

	if hasLen {
		return validateNonNegativeInteger(opMap, "len")
	}

	return nil
}

func validateOperationExtend(opMap map[string]interface{}) error {
	props, exists := opMap["props"]
	if !exists || props == nil {
		return ErrInvalidProps
	}
	if _, isMap := props.(map[string]interface{}); !isMap {
		return ErrInvalidProps
	}

	if deleteNull, exists := opMap["deleteNull"]; exists {
		if _, isBool := deleteNull.(bool); !isBool {
			return ErrExpectedBooleanField
		}
	}

	return nil
}

func validateOperationMerge(opMap map[string]interface{}) error {
	pos, exists := opMap["pos"]
	if !exists {
		return ErrPosGreaterThanZero
	}
	if !isInteger(pos) {
		return ErrExpectedIntegerField
	}
	posInt := getIntValue(pos)
	if posInt < 1 {
		return ErrPosGreaterThanZero
	}

	if props, exists := opMap["props"]; exists {
		if _, isMap := props.(map[string]interface{}); !isMap {
			return ErrInvalidProps
		}
	}

	return nil
}

func validateOperationSplit(opMap map[string]interface{}) error {
	pos, exists := opMap["pos"]
	if !exists {
		return ErrExpectedIntegerField
	}
	if !isInteger(pos) {
		return ErrExpectedIntegerField
	}

	if props, exists := opMap["props"]; exists {
		if _, isMap := props.(map[string]interface{}); !isMap {
			return ErrInvalidProps
		}
	}

	return nil
}

// Predicate operation validators
func validateOperationTestType(opMap map[string]interface{}) error {
	typeField, exists := opMap["type"]
	if !exists {
		return ErrInvalidTypeField
	}

	typeSlice, isSlice := typeField.([]interface{})
	if !isSlice {
		return ErrInvalidTypeField
	}

	if len(typeSlice) < 1 {
		return ErrEmptyTypeList
	}

	for _, t := range typeSlice {
		typeStr, isString := t.(string)
		if !isString {
			return ErrInvalidType
		}
		if err := validateTestType(typeStr); err != nil {
			return err
		}
	}
	return nil
}

func validateOperationTestString(opMap map[string]interface{}) error {
	if err := validateNot(opMap); err != nil {
		return err
	}
	if err := validateNonNegativeInteger(opMap, "pos"); err != nil {
		return err
	}
	str, exists := opMap["str"]
	if !exists {
		return ErrValueMustBeString
	}
	if _, isString := str.(string); !isString {
		return ErrValueMustBeString
	}
	return nil
}

func validateOperationTestStringLen(opMap map[string]interface{}) error {
	if err := validateNot(opMap); err != nil {
		return err
	}
	return validateNonNegativeInteger(opMap, "len")
}

func validateOperationMatches(opMap map[string]interface{}) error {
	if err := validateValueString(opMap, "value"); err != nil {
		return err
	}
	return validateIgnoreCase(opMap)
}

func validateOperationContains(opMap map[string]interface{}) error {
	if err := validateValueString(opMap, "value"); err != nil {
		return err
	}
	return validateIgnoreCase(opMap)
}

func validateOperationEnds(opMap map[string]interface{}) error {
	if err := validateValueString(opMap, "value"); err != nil {
		return err
	}
	return validateIgnoreCase(opMap)
}

func validateOperationStarts(opMap map[string]interface{}) error {
	if err := validateValueString(opMap, "value"); err != nil {
		return err
	}
	return validateIgnoreCase(opMap)
}

func validateOperationIn(opMap map[string]interface{}) error {
	value, exists := opMap["value"]
	if !exists {
		return ErrInOperationValueMustBeArray
	}
	if !isArray(value) {
		return ErrInOperationValueMustBeArray
	}
	return nil
}

func validateOperationMore(opMap map[string]interface{}) error {
	value, exists := opMap["value"]
	if !exists {
		return ErrValueMustBeNumber
	}
	if !isNumber(value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationLess(opMap map[string]interface{}) error {
	value, exists := opMap["value"]
	if !exists {
		return ErrValueMustBeNumber
	}
	if !isNumber(value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationType(opMap map[string]interface{}) error {
	value, exists := opMap["value"]
	if !exists {
		return ErrExpectedValueToBeString
	}
	valueStr, isString := value.(string)
	if !isString {
		return ErrExpectedValueToBeString
	}
	return validateTestType(valueStr)
}

func validateCompositeOperation(opMap map[string]interface{}, opStr string, allowMatchesOp bool) error {
	apply, exists := opMap["apply"]
	if !exists {
		return fmt.Errorf("%w: %s predicate operators must be an array", ErrMustBeArray, opStr)
	}

	applySlice, isSlice := apply.([]interface{})
	if !isSlice {
		return fmt.Errorf("%w: %s predicate operators must be an array", ErrMustBeArray, opStr)
	}

	if len(applySlice) == 0 {
		return ErrEmptyPredicateList
	}

	for _, predicate := range applySlice {
		predicateMap, isMap := predicate.(map[string]interface{})
		if !isMap {
			return ErrInvalidOperation
		}
		if err := ValidateOperation(predicateMap, allowMatchesOp); err != nil {
			return err
		}
	}
	return nil
}

// Helper validation functions
func validateValueString(opMap map[string]interface{}, fieldName string) error {
	value, exists := opMap[fieldName]
	if !exists {
		//nolint:err113 // Dynamic error message needed for field name specificity
		return fmt.Errorf("expected %s to be string", fieldName)
	}
	valueStr, isString := value.(string)
	if !isString {
		//nolint:err113 // Dynamic error message needed for field name specificity
		return fmt.Errorf("expected %s to be string", fieldName)
	}
	if len(valueStr) > 20000 {
		return ErrValueTooLong
	}
	return nil
}

func validateIgnoreCase(opMap map[string]interface{}) error {
	if ignoreCase, exists := opMap["ignore_case"]; exists {
		if _, isBool := ignoreCase.(bool); !isBool {
			return ErrExpectedIgnoreCaseBoolean
		}
	}
	return nil
}

func validateNot(opMap map[string]interface{}) error {
	if not, exists := opMap["not"]; exists {
		if _, isBool := not.(bool); !isBool {
			return ErrInvalidNotModifier
		}
	}
	return nil
}

func validateNonNegativeInteger(opMap map[string]interface{}, fieldName string) error {
	value, exists := opMap[fieldName]
	if !exists {
		return ErrExpectedIntegerField
	}
	if !isInteger(value) {
		return ErrExpectedIntegerField
	}
	intValue := getIntValue(value)
	if intValue < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateJSONPointer(path string) error {
	return jsonpointer.Validate(path)
}

var validTypesMap = map[string]bool{
	"string":  true,
	"number":  true,
	"boolean": true,
	"object":  true,
	"integer": true,
	"array":   true,
	"null":    true,
}

func validateTestType(typeStr string) error {
	if !validTypesMap[typeStr] {
		return ErrInvalidType
	}
	return nil
}

// Type checking helper functions
func isNumber(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func isInteger(value interface{}) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float64:
		return v == float64(int64(v))
	case float32:
		return v == float32(int32(v))
	default:
		return false
	}
}

func getIntValue(value interface{}) int64 {
	switch v := value.(type) {
	case int:
		return int64(v)
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case uint:
		if v > uint(^uint64(0)>>1) {
			return 0 // Return 0 for overflow cases
		}
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		if v > (^uint64(0) >> 1) {
			return 0 // Return 0 for overflow cases
		}
		return int64(v)
	case float64:
		return int64(v)
	case float32:
		return int64(v)
	default:
		return 0
	}
}

func isArray(value interface{}) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case []interface{}, []string, []int, []float64:
		return true
	default:
		rv := reflect.ValueOf(value)
		return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
	}
}
