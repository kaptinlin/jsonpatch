package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"strings"

	"github.com/kaptinlin/jsonpointer"

	"github.com/kaptinlin/jsonpatch/internal"
)

var (
	// ErrNotArray reports that a patch value is not an operation array.
	ErrNotArray = errors.New("not an array")
	// ErrEmptyPatch reports that no operations were provided.
	ErrEmptyPatch = errors.New("empty operation patch")
	// ErrInvalidOperation reports that an operation name is unknown.
	ErrInvalidOperation = errors.New("invalid operation")
	// ErrMissingPath reports that an operation is missing its path field.
	ErrMissingPath = errors.New("missing required field 'path'")
	// ErrMissingOp reports that an operation is missing its op field.
	ErrMissingOp = errors.New("missing required field 'op'")
	// ErrMissingValue reports that an operation is missing its value field.
	ErrMissingValue = errors.New("missing required field 'value'")
	// ErrMissingFrom reports that an operation is missing its from field.
	ErrMissingFrom = errors.New("missing required field 'from'")
	// ErrInvalidPath reports that path has the wrong JSON type.
	ErrInvalidPath = errors.New("field 'path' must be a string")
	// ErrInvalidOp reports that op has the wrong JSON type.
	ErrInvalidOp = errors.New("field 'op' must be a string")
	// ErrInvalidFrom reports that from has the wrong JSON type.
	ErrInvalidFrom = errors.New("field 'from' must be a string")
	// ErrInvalidJSONPointer reports that a path or from value is not a valid JSON Pointer.
	ErrInvalidJSONPointer = errors.New("invalid JSON pointer")
	// ErrInvalidOldValue reports that oldValue has an invalid shape.
	ErrInvalidOldValue = errors.New("invalid oldValue")
	// ErrCannotMoveToChildren reports that a move target is nested under its source.
	ErrCannotMoveToChildren = errors.New("cannot move into own children")
	// ErrInvalidIncValue reports that an inc operand has an invalid type.
	ErrInvalidIncValue = errors.New("invalid inc value")
	// ErrExpectedStringField reports that a field must contain a string.
	ErrExpectedStringField = errors.New("expected string field")
	// ErrExpectedBooleanField reports that a field must contain a boolean.
	ErrExpectedBooleanField = errors.New("expected field to be boolean")
	// ErrExpectedIntegerField reports that a field must contain an integer.
	ErrExpectedIntegerField = errors.New("not an integer")
	// ErrNegativeNumber reports that a numeric field is negative.
	ErrNegativeNumber = errors.New("number is negative")
	// ErrInvalidProps reports that props has an invalid shape.
	ErrInvalidProps = errors.New("invalid props field")
	// ErrInvalidTypeField reports that type has an invalid shape.
	ErrInvalidTypeField = errors.New("invalid type field")
	// ErrEmptyTypeList reports that a type list is empty.
	ErrEmptyTypeList = errors.New("empty type list")
	// ErrInvalidType reports that a JSON type name is unsupported.
	ErrInvalidType = errors.New("invalid type")
	// ErrValueMustBeString reports that value must be a string.
	ErrValueMustBeString = errors.New("value must be a string")
	// ErrValueMustBeNumber reports that value must be numeric.
	ErrValueMustBeNumber = errors.New("value must be a number")
	// ErrValueMustBeArray reports that value must be an array.
	ErrValueMustBeArray = errors.New("value must be an array")
	// ErrValueTooLong reports that a provided value exceeds the allowed length.
	ErrValueTooLong = errors.New("value too long")
	// ErrInvalidNotModifier reports that not was supplied where it is unsupported.
	ErrInvalidNotModifier = errors.New("invalid not modifier")
	// ErrMatchesNotAllowed reports that matches is disabled for this validation pass.
	ErrMatchesNotAllowed = errors.New("matches operation not allowed")
	// ErrMustBeArray reports that a value must decode to an array.
	ErrMustBeArray = errors.New("must be an array")
	// ErrEmptyPredicateList reports that a composite predicate has no operands.
	ErrEmptyPredicateList = errors.New("predicate list is empty")
	// ErrPosGreaterThanZero reports that pos must be greater than zero.
	ErrPosGreaterThanZero = errors.New("expected pos field to be greater than 0")
	// ErrInOperationValueMustBeArray reports that in requires an array value.
	ErrInOperationValueMustBeArray = errors.New("in operation value must be an array")
	// ErrExpectedValueToBeString reports that value must decode as a string.
	ErrExpectedValueToBeString = errors.New("expected value to be string")
	// ErrExpectedIgnoreCaseBoolean reports that ignore_case must decode as a boolean.
	ErrExpectedIgnoreCaseBoolean = errors.New("expected ignore_case to be boolean")
	// ErrExpectedFieldString reports that a field must decode as a string.
	ErrExpectedFieldString = errors.New("expected field to be string")
)

// ValidateOperations validates an array of JSON Patch operations.
func ValidateOperations(ops []Operation, allowMatchesOp bool) error {
	if ops == nil {
		return ErrNotArray
	}
	if len(ops) == 0 {
		return ErrEmptyPatch
	}

	for i := range ops {
		if err := validateOperation(&ops[i], allowMatchesOp); err != nil {
			return fmt.Errorf("error in operation [index = %d] (%w)", i, err)
		}
	}
	return nil
}

// ValidateOperation validates a single JSON Patch operation.
//
//nolint:gocritic // Preserve the exported by-value API.
func ValidateOperation(op Operation, allowMatchesOp bool) error {
	return validateOperation(&op, allowMatchesOp)
}

func validateOperation(op *Operation, allowMatchesOp bool) error {
	if op.Op == "" {
		return ErrMissingOp
	}

	if op.Path == "" {
		return ErrMissingPath
	}
	if err := validateJSONPointer(op.Path); err != nil {
		return ErrInvalidJSONPointer
	}

	switch op.Op {
	case "add":
		return validateOperationAdd(op)
	case "remove", "flip", "inc", "extend", "split":
		return nil
	case "replace":
		return validateOperationReplace(op)
	case "copy":
		return validateOperationCopy(op)
	case "move":
		return validateOperationMove(op)
	case "str_ins":
		return validateOperationStrIns(op)
	case "str_del":
		return validateOperationStrDel(op)
	case "merge":
		return validateOperationMerge(op)
	default:
		return validatePredicateOperation(op, op.Op, allowMatchesOp)
	}
}

func validatePredicateOperation(op *Operation, opStr string, allowMatchesOp bool) error {
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
	case "defined", "undefined":
		return nil
	case "and", "or", "not":
		internalOp := internal.Operation(*op)
		return validateCompositeOperation(&internalOp, allowMatchesOp)
	default:
		return fmt.Errorf("%w: unknown operation '%s'", ErrInvalidOperation, opStr)
	}
}

func validateOperationAdd(op *Operation) error {
	if op.Value == nil {
		return ErrMissingValue
	}
	return nil
}

func validateOperationReplace(op *Operation) error {
	if op.Value == nil {
		return ErrMissingValue
	}
	return nil
}

func validateOperationCopy(op *Operation) error {
	if op.From == "" {
		return ErrMissingFrom
	}
	return validateJSONPointer(op.From)
}

func validateOperationMove(op *Operation) error {
	if op.From == "" {
		return ErrMissingFrom
	}
	if err := validateJSONPointer(op.From); err != nil {
		return err
	}

	if _, found := strings.CutPrefix(op.Path, op.From+"/"); found {
		return ErrCannotMoveToChildren
	}
	return nil
}

func validateOperationTest(op *Operation) error {
	if op.Value == nil {
		return ErrMissingValue
	}
	return nil
}

func validateOperationStrIns(op *Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateOperationStrDel(op *Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	if op.Len < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateOperationMerge(op *Operation) error {
	if op.Pos < 1 {
		return ErrPosGreaterThanZero
	}
	return nil
}

func validateOperationTestType(op *Operation) error {
	if op.Type == nil {
		return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
	}

	if typeStr, ok := op.Type.(string); ok {
		if typeStr == "" {
			return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
		}
		if !slices.Contains(validTypes, typeStr) {
			return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
		}
		return nil
	}

	if typeSlice, ok := op.Type.([]any); ok {
		if len(typeSlice) == 0 {
			return fmt.Errorf("%w: type array cannot be empty", ErrInvalidTypeField)
		}
		for i := range typeSlice {
			typeStr, isString := typeSlice[i].(string)
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

func validateOperationTestString(op *Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateOperationTestStringLen(op *Operation) error {
	if op.Len < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func requireStringValue(op *Operation) error {
	if op.Value == nil {
		return ErrExpectedValueToBeString
	}
	if _, isString := op.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func requireNumberValue(op *Operation) error {
	if op.Value == nil {
		return ErrValueMustBeNumber
	}
	if !isNumber(op.Value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationMatches(op *Operation) error {
	return requireStringValue(op)
}

func validateOperationContains(op *Operation) error {
	return requireStringValue(op)
}

func validateOperationEnds(op *Operation) error {
	return requireStringValue(op)
}

func validateOperationStarts(op *Operation) error {
	return requireStringValue(op)
}

func validateOperationIn(op *Operation) error {
	if op.Value == nil {
		return ErrInOperationValueMustBeArray
	}
	if !isArray(op.Value) {
		return ErrInOperationValueMustBeArray
	}
	return nil
}

func validateOperationMore(op *Operation) error {
	return requireNumberValue(op)
}

func validateOperationLess(op *Operation) error {
	return requireNumberValue(op)
}

func validateOperationType(op *Operation) error {
	if op.Value == nil {
		return ErrExpectedValueToBeString
	}
	valueStr, isString := op.Value.(string)
	if !isString {
		return ErrExpectedValueToBeString
	}
	return validateTestType(valueStr)
}

func validateCompositeOperation(op *internal.Operation, allowMatchesOp bool) error {
	if len(op.Apply) == 0 {
		return ErrEmptyPredicateList
	}

	for i := range op.Apply {
		if err := validateOperation(&op.Apply[i], allowMatchesOp); err != nil {
			return err
		}
	}
	return nil
}

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

// isNumber reports whether value is a numeric type.
func isNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

// isArray reports whether value is an array or slice type.
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
