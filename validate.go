package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"
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
	// ErrMissingOp reports that an operation is missing its op field.
	ErrMissingOp = errors.New("missing required field 'op'")
	// ErrInvalidJSONPointer reports that a path or from value is not a valid JSON Pointer.
	ErrInvalidJSONPointer = errors.New("invalid JSON pointer")
	// ErrCannotMoveToChildren reports that a move target is nested under its source.
	ErrCannotMoveToChildren = errors.New("cannot move into own children")
	// ErrNegativeNumber reports that a numeric field is negative.
	ErrNegativeNumber = errors.New("number is negative")
	// ErrInvalidTypeField reports that type has an invalid shape.
	ErrInvalidTypeField = errors.New("invalid type field")
	// ErrInvalidType reports that a JSON type name is unsupported.
	ErrInvalidType = errors.New("invalid type")
	// ErrValueMustBeNumber reports that value must be numeric.
	ErrValueMustBeNumber = errors.New("value must be a number")
	// ErrMatchesNotAllowed reports that matches is disabled for this validation pass.
	ErrMatchesNotAllowed = errors.New("matches operation not allowed")
	// ErrEmptyPredicateList reports that a composite predicate has no operands.
	ErrEmptyPredicateList = errors.New("predicate list is empty")
	// ErrNotRequiresSinglePredicate reports that not received anything other than one operand.
	ErrNotRequiresSinglePredicate = errors.New("not operation requires exactly one predicate")
	// ErrPosGreaterThanZero reports that pos must be greater than zero.
	ErrPosGreaterThanZero = errors.New("expected pos field to be greater than 0")
	// ErrInOperationValueMustBeArray reports that in requires an array value.
	ErrInOperationValueMustBeArray = errors.New("in operation value must be an array")
	// ErrExpectedValueToBeString reports that value must decode as a string.
	ErrExpectedValueToBeString = errors.New("expected value to be string")
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

func invalidJSONPointerError(err error) error {
	if err == nil {
		return nil
	}
	return errors.Join(ErrInvalidJSONPointer, err)
}

func validateOperation(op *Operation, allowMatchesOp bool) error {
	if op.Op == "" {
		return ErrMissingOp
	}

	if err := invalidJSONPointerError(jsonpointer.Validate(op.Path)); err != nil {
		return err
	}

	switch op.Op {
	case "add", "replace":
		return validateValueRequired(op)
	case "remove", "flip", "inc", "extend", "split":
		return nil
	case "copy":
		return validateOperationCopy(op)
	case "move":
		return validateOperationMove(op)
	case "str_ins":
		return validateOperationPosition(op)
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
		return validateValueRequired(op)
	case "test_type":
		return validateOperationTestType(op)
	case "test_string":
		return validateOperationPosition(op)
	case "test_string_len":
		if op.Len < 0 {
			return ErrNegativeNumber
		}
		return nil
	case "matches":
		if !allowMatchesOp {
			return ErrMatchesNotAllowed
		}
		return requireStringValue(op)
	case "contains", "ends", "starts":
		return requireStringValue(op)
	case "in":
		return validateOperationIn(op)
	case "more", "less":
		return requireNumberValue(op)
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

func validateValueRequired(op *Operation) error {
	// Operation is a Go builder shape, not a raw JSON object. It cannot
	// distinguish an omitted value field from a deliberate JSON null, so
	// wire-level value presence is validated by the JSON codec.
	return nil
}

func validateOperationCopy(op *Operation) error {
	return invalidJSONPointerError(jsonpointer.Validate(op.From))
}

func validateOperationMove(op *Operation) error {
	if err := invalidJSONPointerError(jsonpointer.Validate(op.From)); err != nil {
		return err
	}

	if _, found := strings.CutPrefix(op.Path, op.From+"/"); found {
		return ErrCannotMoveToChildren
	}
	return nil
}

func validateOperationPosition(op *Operation) error {
	if op.Pos < 0 {
		return ErrNegativeNumber
	}
	return nil
}

func validateOperationStrDel(op *Operation) error {
	if err := validateOperationPosition(op); err != nil {
		return err
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

	switch typeValue := op.Type.(type) {
	case string:
		if typeValue == "" {
			return fmt.Errorf("%w: missing required field 'type'", ErrInvalidTypeField)
		}
		return validateTypeName(typeValue)
	case []any:
		if len(typeValue) == 0 {
			return fmt.Errorf("%w: type array cannot be empty", ErrInvalidTypeField)
		}
		for i := range typeValue {
			typeStr, isString := typeValue[i].(string)
			if !isString {
				return fmt.Errorf("%w: all types must be strings", ErrInvalidType)
			}
			if err := validateTypeName(typeStr); err != nil {
				return err
			}
		}
		return nil
	case []string:
		if len(typeValue) == 0 {
			return fmt.Errorf("%w: type array cannot be empty", ErrInvalidTypeField)
		}
		for i := range typeValue {
			if err := validateTypeName(typeValue[i]); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("%w: type field must be string or array of strings", ErrInvalidType)
	}
}

func requireStringValue(op *Operation) error {
	if _, isString := op.Value.(string); !isString {
		return ErrExpectedValueToBeString
	}
	return nil
}

func requireNumberValue(op *Operation) error {
	if !isNumber(op.Value) {
		return ErrValueMustBeNumber
	}
	return nil
}

func validateOperationIn(op *Operation) error {
	if !isArray(op.Value) {
		return ErrInOperationValueMustBeArray
	}
	return nil
}

func validateOperationType(op *Operation) error {
	if op.Value == nil {
		return ErrExpectedValueToBeString
	}
	valueStr, isString := op.Value.(string)
	if !isString {
		return ErrExpectedValueToBeString
	}
	if !internal.IsValidJSONPatchType(valueStr) {
		return ErrInvalidType
	}
	return nil
}

func validateCompositeOperation(op *internal.Operation, allowMatchesOp bool) error {
	if len(op.Apply) == 0 {
		return ErrEmptyPredicateList
	}
	if op.Op == "not" && len(op.Apply) != 1 {
		return ErrNotRequiresSinglePredicate
	}

	for i := range op.Apply {
		if !isPredicateOperationName(op.Apply[i].Op) {
			return ErrInvalidOperation
		}
		if err := validateOperation(&op.Apply[i], allowMatchesOp); err != nil {
			return err
		}
	}
	return nil
}

func isPredicateOperationName(name string) bool {
	switch name {
	case "test", "test_type", "test_string", "test_string_len",
		"matches", "contains", "ends", "starts", "in", "more", "less",
		"type", "defined", "undefined", "and", "or", "not":
		return true
	default:
		return false
	}
}

func validateTypeName(typeStr string) error {
	if !internal.IsValidJSONPatchType(typeStr) {
		return fmt.Errorf("%w: invalid type '%s'", ErrInvalidType, typeStr)
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
	rv := reflect.ValueOf(value)
	return rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array
}
