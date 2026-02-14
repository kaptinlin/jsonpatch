package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
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

	for i, op := range ops {
		if err := ValidateOperation(op, allowMatchesOp); err != nil {
			return fmt.Errorf("error in operation [index = %d] (%w)", i, err)
		}
	}
	return nil
}

// ValidateOperation validates a single JSON Patch operation.
func ValidateOperation(op Operation, allowMatchesOp bool) error {
	// Validate op field first
	if op.Op == "" {
		return ErrMissingOp
	}

	// Validate path field
	if op.Path == "" {
		return ErrMissingPath
	}
	if err := validateJSONPointer(op.Path); err != nil {
		return ErrInvalidJSONPointer
	}

	// Validate operation by type
	switch op.Op {
	case "add":
		return validateOperationAdd(op)
	case "remove":
		return validateOperationRemove(op)
	case "replace":
		return validateOperationReplace(op)
	case "copy":
		return validateOperationCopy(op)
	case "move":
		return validateOperationMove(op)
	case "flip":
		return nil
	case "inc":
		return validateOperationInc(op)
	case "str_ins":
		return validateOperationStrIns(op)
	case "str_del":
		return validateOperationStrDel(op)
	case "extend":
		return validateOperationExtend(op)
	case "merge":
		return validateOperationMerge(op)
	case "split":
		return validateOperationSplit(op)
	default:
		return validatePredicateOperation(op, op.Op, allowMatchesOp)
	}
}

// validatePredicateOperation validates predicate operations
func validatePredicateOperation(op Operation, opStr string, allowMatchesOp bool) error {
	switch opStr {
	case "test":
		return validateOperationTest(op)
	case "test_type":
		return validateOperationTestType(op)
	case "test_string":
		return validateOperationTestString(op)
	case "test_string_len":
		return validateOperationTestStringLen(op)
	case "matches":
		if !allowMatchesOp {
			return ErrMatchesNotAllowed
		}
		return validateOperationMatches(op)
	case "contains":
		return validateOperationContains(op)
	case "ends":
		return validateOperationEnds(op)
	case "starts":
		return validateOperationStarts(op)
	case "in":
		return validateOperationIn(op)
	case "more":
		return validateOperationMore(op)
	case "less":
		return validateOperationLess(op)
	case "type":
		return validateOperationType(op)
	case "defined":
		return nil
	case "undefined":
		return nil
	case "and", "or", "not":
		return validateCompositeOperation(op, allowMatchesOp)
	default:
		return fmt.Errorf("%w: unknown operation '%s'", ErrInvalidOperation, opStr)
	}
}

// Core operation validators
func validateOperationAdd(op Operation) error {
	if op.Value == nil {
		return ErrMissingValue
	}
	return nil
}

func validateOperationRemove(_ Operation) error {
	// OldValue is optional, no validation needed for struct-based approach
	return nil
}

func validateOperationReplace(op Operation) error {
	if op.Value == nil {
		return ErrMissingValue
	}
	// OldValue is optional, no validation needed for struct-based approach
	return nil
}

func validateOperationCopy(op Operation) error {
	if op.From == "" {
		return ErrMissingFrom
	}
	return validateJSONPointer(op.From)
}

func validateOperationMove(op Operation) error {
	if op.From == "" {
		return ErrMissingFrom
	}
	if err := validateJSONPointer(op.From); err != nil {
		return err
	}

	if strings.HasPrefix(op.Path, op.From+"/") {
		return ErrCannotMoveToChildren
	}
	return nil
}

func validateOperationTest(op Operation) error {
	if op.Value == nil {
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

func validateOperationStrIns(op Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	// Str field can be empty (for inserting empty string)
	return nil
}

func validateOperationStrDel(op Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}

	// Either Str or Len should be provided (but not required to error if both missing, as len=0 is valid)
	if op.Len < 0 {
		return ErrNegativeNumber
	}

	return nil
}

func validateOperationExtend(_ Operation) error {
	// Props can be nil (treated as empty object)
	return nil
}

func validateOperationMerge(op Operation) error {
	if op.Pos < 1 {
		return ErrPosGreaterThanZero
	}
	return nil
}

func validateOperationSplit(_ Operation) error {
	// Pos can be any integer for split operation
	return nil
}

// Predicate operation validators
func validateOperationTestType(op Operation) error {
	if op.Type == nil {
		return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
	}

	// Handle single type string
	if typeStr, ok := op.Type.(string); ok {
		if typeStr == "" {
			return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
		}
		if !slices.Contains(validTypes, typeStr) {
			return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
		}
		return nil
	}

	// Handle array of types
	if typeSlice, ok := op.Type.([]any); ok {
		if len(typeSlice) == 0 {
			return fmt.Errorf("%w: type array cannot be empty", ErrInvalidTypeField)
		}
		for _, t := range typeSlice {
			typeStr, isString := t.(string)
			if !isString {
				return fmt.Errorf("%w: all types must be strings", ErrInvalidType)
			}
			if !slices.Contains(validTypes, typeStr) {
				return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
			}
		}
		return nil
	}

	return fmt.Errorf("%w: type field must be string or array of strings", ErrInvalidType)
}

func validateOperationTestString(op Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	// Str can be empty (to test for empty string at position)
	return nil
}

func validateOperationTestStringLen(op Operation) error {
	if op.Len < 0 {
		return ErrNegativeNumber
	}
	return nil
}

// requireStringValue validates that op.Value is a non-nil string.
// This eliminates code duplication across string predicate validators.
func requireStringValue(op Operation) error {
	if op.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := op.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

// requireNumberValue validates that op.Value is a non-nil number.
// This eliminates code duplication across numeric predicate validators.
func requireNumberValue(op Operation, err error) error {
	if op.Value == nil {
		return err
	}
	if !isNumber(op.Value) {
		return err
	}
	return nil
}

func validateOperationMatches(op Operation) error {
	return requireStringValue(op)
}

func validateOperationContains(op Operation) error {
	return requireStringValue(op)
}

func validateOperationEnds(op Operation) error {
	return requireStringValue(op)
}

func validateOperationStarts(op Operation) error {
	return requireStringValue(op)
}

func validateOperationIn(op Operation) error {
	if op.Value == nil {
		return ErrInOperationValueMustBeArray
	}
	if !isArray(op.Value) {
		return ErrInOperationValueMustBeArray
	}
	return nil
}

func validateOperationMore(op Operation) error {
	return requireNumberValue(op, ErrValueMustBeNumber)
}

func validateOperationLess(op Operation) error {
	return requireNumberValue(op, ErrValueMustBeNumber)
}

func validateOperationType(op Operation) error {
	if op.Value == nil {
		return ErrExpectedValueToBeString
	}
	valueStr, isString := op.Value.(string)
	if !isString {
		return ErrExpectedValueToBeString
	}
	return validateTestType(valueStr)
}

func validateCompositeOperation(op internal.Operation, allowMatchesOp bool) error {
	if len(op.Apply) == 0 {
		return ErrEmptyPredicateList
	}

	for _, predicate := range op.Apply {
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

var validTypes = []string{
	"string", "number", "boolean", "object", "integer", "array", "null",
}

func validateTestType(typeStr string) error {
	if !slices.Contains(validTypes, typeStr) {
		return ErrInvalidType
	}
	return nil
}

// Type checking helper functions
func isNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func isArray(value any) bool {
	if value == nil {
		return false
	}
	switch value.(type) {
	case []any, []string, []int, []float64:
		return true
	default:
		rv := reflect.ValueOf(value)
		return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
	}
}
