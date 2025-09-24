package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/kaptinlin/jsonpointer"
	"github.com/kaptinlin/jsonpatch/internal"
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
	// Validate op field first
	if operation.Op == "" {
		return ErrMissingOp
	}

	// Validate path field
	if operation.Path == "" {
		return ErrMissingPath
	}
	if err := validateJSONPointer(operation.Path); err != nil {
		return ErrInvalidJSONPointer
	}

	// Validate operation by type
	switch operation.Op {
	case "add":
		return validateOperationAdd(operation)
	case "remove":
		return validateOperationRemove(operation)
	case "replace":
		return validateOperationReplace(operation)
	case "copy":
		return validateOperationCopy(operation)
	case "move":
		return validateOperationMove(operation)
	case "flip":
		return nil
	case "inc":
		return validateOperationInc(operation)
	case "str_ins":
		return validateOperationStrIns(operation)
	case "str_del":
		return validateOperationStrDel(operation)
	case "extend":
		return validateOperationExtend(operation)
	case "merge":
		return validateOperationMerge(operation)
	case "split":
		return validateOperationSplit(operation)
	default:
		return validatePredicateOperation(operation, operation.Op, allowMatchesOp)
	}
}

// validatePredicateOperation validates predicate operations
func validatePredicateOperation(operation Operation, opStr string, allowMatchesOp bool) error {
	switch opStr {
	case "test":
		return validateOperationTest(operation)
	case "test_type":
		return validateOperationTestType(operation)
	case "test_string":
		return validateOperationTestString(operation)
	case "test_string_len":
		return validateOperationTestStringLen(operation)
	case "matches":
		if !allowMatchesOp {
			return ErrMatchesNotAllowed
		}
		return validateOperationMatches(operation)
	case "contains":
		return validateOperationContains(operation)
	case "ends":
		return validateOperationEnds(operation)
	case "starts":
		return validateOperationStarts(operation)
	case "in":
		return validateOperationIn(operation)
	case "more":
		return validateOperationMore(operation)
	case "less":
		return validateOperationLess(operation)
	case "type":
		return validateOperationType(operation)
	case "defined":
		return nil
	case "undefined":
		return nil
	case "and", "or", "not":
		return validateCompositeOperation(operation, opStr, allowMatchesOp)
	default:
		return fmt.Errorf("%w: unknown operation '%s'", ErrInvalidOperation, opStr)
	}
}

// Core operation validators
func validateOperationAdd(operation Operation) error {
	if operation.Value == nil {
		return ErrMissingValue
	}
	return nil
}

func validateOperationRemove(_ Operation) error {
	// OldValue is optional, no validation needed for struct-based approach
	return nil
}

func validateOperationReplace(operation Operation) error {
	if operation.Value == nil {
		return ErrMissingValue
	}
	// OldValue is optional, no validation needed for struct-based approach
	return nil
}

func validateOperationCopy(operation Operation) error {
	if operation.From == "" {
		return ErrMissingFrom
	}
	return validateJSONPointer(operation.From)
}

func validateOperationMove(operation Operation) error {
	if operation.From == "" {
		return ErrMissingFrom
	}
	if err := validateJSONPointer(operation.From); err != nil {
		return err
	}

	if strings.HasPrefix(operation.Path, operation.From+"/") {
		return ErrCannotMoveToChildren
	}
	return nil
}

func validateOperationTest(operation Operation) error {
	if operation.Value == nil {
		return ErrMissingValue
	}
	return nil
}

// Extended operation validators
func validateOperationInc(_ Operation) error {
	// Inc field can be any number, including 0
	// The field is already defined, no validation needed
	return nil
}

func validateOperationStrIns(operation Operation) error {
	if operation.Pos < 0 {
		return ErrNegativeNumber
	}
	// Str field can be empty (for inserting empty string)
	return nil
}

func validateOperationStrDel(operation Operation) error {
	if operation.Pos < 0 {
		return ErrNegativeNumber
	}

	// Either Str or Len should be provided (but not required to error if both missing, as len=0 is valid)
	if operation.Len < 0 {
		return ErrNegativeNumber
	}

	return nil
}

func validateOperationExtend(_ Operation) error {
	// Props can be nil (treated as empty object)
	return nil
}

func validateOperationMerge(operation Operation) error {
	if operation.Pos < 1 {
		return ErrPosGreaterThanZero
	}
	return nil
}

func validateOperationSplit(_ Operation) error {
	// Pos can be any integer for split operation
	return nil
}

// Predicate operation validators
func validateOperationTestType(operation Operation) error {
	if operation.Type == nil {
		return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
	}
	
	// Handle single type string
	if typeStr, ok := operation.Type.(string); ok {
		if typeStr == "" {
			return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
		}
		if !validTypesMap[typeStr] {
			return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
		}
		return nil
	}
	
	// Handle array of types
	if typeSlice, ok := operation.Type.([]interface{}); ok {
		if len(typeSlice) == 0 {
			return fmt.Errorf("%w: type array cannot be empty", ErrInvalidTypeField)
		}
		for _, t := range typeSlice {
			typeStr, isString := t.(string)
			if !isString {
				return fmt.Errorf("%w: all types must be strings", ErrInvalidType)
			}
			if !validTypesMap[typeStr] {
				return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
			}
		}
		return nil
	}
	
	return fmt.Errorf("%w: type field must be string or array of strings", ErrInvalidType)
}

func validateOperationTestString(operation Operation) error {
	if operation.Pos < 0 {
		return ErrNegativeNumber
	}
	// Str can be empty (to test for empty string at position)
	return nil
}

func validateOperationTestStringLen(operation Operation) error {
	if operation.Len < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateOperationMatches(operation Operation) error {
	if operation.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := operation.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func validateOperationContains(operation Operation) error {
	if operation.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := operation.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func validateOperationEnds(operation Operation) error {
	if operation.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := operation.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func validateOperationStarts(operation Operation) error {
	if operation.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := operation.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func validateOperationIn(operation Operation) error {
	if operation.Value == nil {
		return ErrInOperationValueMustBeArray
	}
	if !isArray(operation.Value) {
		return ErrInOperationValueMustBeArray
	}
	return nil
}

func validateOperationMore(operation Operation) error {
	if operation.Value == nil {
		return ErrValueMustBeNumber
	}
	if !isNumber(operation.Value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationLess(operation Operation) error {
	if operation.Value == nil {
		return ErrValueMustBeNumber
	}
	if !isNumber(operation.Value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationType(operation Operation) error {
	if operation.Value == nil {
		return ErrExpectedValueToBeString
	}
	valueStr, isString := operation.Value.(string)
	if !isString {
		return ErrExpectedValueToBeString
	}
	return validateTestType(valueStr)
}

func validateCompositeOperation(operation internal.Operation, _ string, allowMatchesOp bool) error {
	if len(operation.Apply) == 0 {
		return ErrEmptyPredicateList
	}

	for _, predicate := range operation.Apply {
		if err := ValidateOperation(predicate, allowMatchesOp); err != nil {
			return err
		}
	}
	return nil
}

// Helper validation functions

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
