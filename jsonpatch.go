// Package jsonpatch provides comprehensive JSON Patch operations with generic type support.
//
// Implements JSON mutation operations including:
//   - JSON Patch (RFC 6902): Standard operations (add, remove, replace, move, copy, test)
//     https://tools.ietf.org/html/rfc6902
//   - JSON Predicate: Test operations (contains, defined, type, less, more, etc.)
//     https://tools.ietf.org/id/draft-snell-json-test-01.html
//   - Extended operations: Additional operations (flip, inc, str_ins, str_del, split, merge)
//
// Core API Functions:
//   - ApplyOp: Apply a single operation
//   - ApplyOps: Apply multiple operations
//   - ApplyPatch: Apply a JSON Patch to a document (main generic API)
//   - ValidateOperations: Validate an array of operations
//   - ValidateOperation: Validate a single operation
//
// Basic usage:
//
//	doc := map[string]any{"name": "John", "age": 30}
//	patch := []Operation{
//		{"op": "replace", "path": "/name", "value": "Jane"},
//		{"op": "add", "path": "/email", "value": "jane@example.com"},
//	}
//	result, err := ApplyPatch(doc, patch, WithMutate(false))
//
// The library provides type-safe operations for any supported document type.
package jsonpatch

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/go-json-experiment/json"

	"github.com/kaptinlin/deepclone"
	jsoncodec "github.com/kaptinlin/jsonpatch/codec/json"
	"github.com/kaptinlin/jsonpatch/internal"
)

// Operation application errors.
var (
	ErrNoOperationDecoded  = errors.New("no operation decoded")
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrConversionFailed    = errors.New("failed to convert result back to original type")
	ErrNoOperationResult   = errors.New("no operation result")
)

// convertOpResults converts []internal.OpResult[any] to []internal.OpResult[T].
func convertOpResults[T internal.Document](resultOps []internal.OpResult[any], resultDoc T) []internal.OpResult[T] {
	converted := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		converted[i] = internal.OpResult[T]{
			Doc: resultDoc,
			Old: op.Old,
		}
	}
	return converted
}

// ApplyPatch applies a JSON Patch to any supported document type.
// It automatically detects the document type and applies the appropriate strategy.
// Returns a PatchResult containing the patched document and operation results.
//
// Supported document types:
//   - struct: Converted via JSON marshaling/unmarshaling
//   - map[string]any: Applied directly using existing implementation
//   - []byte: Parsed as JSON, patched, and re-encoded
//   - string: Parsed as JSON string, patched, and re-encoded
//
// Example usage:
//
//	// Struct
//	user := User{Name: "John", Age: 30}
//	result, err := ApplyPatch(user, patch)
//	if err == nil {
//		patchedUser := result.Doc // Type: User
//		operations := result.Res  // Operation results
//	}
//
//	// Map
//	doc := map[string]any{"name": "John", "age": 30}
//	result, err := ApplyPatch(doc, patch)
//	if err == nil {
//		patchedDoc := result.Doc // Type: map[string]any
//	}
//
//	// JSON bytes
//	data := []byte(`{"name":"John","age":30}`)
//	result, err := ApplyPatch(data, patch)
//	if err == nil {
//		patchedData := result.Doc // Type: []byte
//	}
//
// The function preserves the input type: struct input returns struct output,
// map input returns map output, etc.
func ApplyPatch[T internal.Document](doc T, patch []internal.Operation, opts ...internal.Option) (*internal.PatchResult[T], error) {
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	return dispatchByDocumentType(doc, patch, options)
}

// dispatchByDocumentType routes the patch operation to the appropriate handler
// based on the runtime type of the document.
func dispatchByDocumentType[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	switch any(doc).(type) {
	case []byte:
		return handleJSONBytes(doc, patch, options)
	case string:
		return handleJSONString(doc, patch, options)
	case map[string]any:
		return handleMapDocument(doc, patch, options)
	case nil:
		return handleStructDocument(doc, patch, options)
	case bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return handlePrimitiveDocument(doc, patch, options)
	case []any:
		return handlePrimitiveDocument(doc, patch, options)
	default:
		return dispatchByReflection(doc, patch, options)
	}
}

// dispatchByReflection handles complex types that require reflection-based dispatch.
func dispatchByReflection[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	docValue := reflect.ValueOf(any(doc))

	if !docValue.IsValid() || (docValue.Kind() == reflect.Pointer && docValue.IsNil()) {
		return handleStructDocument(doc, patch, options)
	}

	switch docValue.Type().Kind() { //nolint:exhaustive // only slice/interface need special handling
	case reflect.Slice, reflect.Interface:
		return handlePrimitiveDocument(doc, patch, options)
	default:
		return handleStructDocument(doc, patch, options)
	}
}

// handleJSONBytes processes []byte documents containing JSON data.
// The bytes are parsed, patched, and re-encoded to maintain format consistency.
func handleJSONBytes[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	docBytes, ok := any(doc).([]byte)
	if !ok {
		return nil, fmt.Errorf("%w: expected []byte, got %T", ErrInvalidDocumentType, doc)
	}

	var parsedDoc any
	if err := json.Unmarshal(docBytes, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON bytes: %w", err)
	}

	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to JSON document: %w", err)
	}

	resultBytes, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patched document: %w", err)
	}

	resultT, ok := any(resultBytes).(T)
	if !ok {
		return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertOpResults(resultOps, resultT),
	}, nil
}

// handleJSONString processes JSON string documents.
// Handles both JSON-encoded strings (starting with { or [) and plain string values.
func handleJSONString[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	docStr, ok := any(doc).(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string, got %T", ErrInvalidDocumentType, doc)
	}

	// Only parse as JSON if the string looks like a JSON object or array.
	// This prevents automatic conversion of simple strings like "123" to numbers.
	var parsedDoc any
	originalWasJSON := len(docStr) > 0 && (docStr[0] == '{' || docStr[0] == '[')
	if originalWasJSON {
		if err := json.Unmarshal([]byte(docStr), &parsedDoc); err != nil {
			parsedDoc = docStr
			originalWasJSON = false
		}
	} else {
		parsedDoc = docStr
	}

	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to document: %w", err)
	}

	resultT, err := convertStringResult(resultDoc, originalWasJSON, doc)
	if err != nil {
		return nil, err
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertOpResults(resultOps, resultT),
	}, nil
}

// convertStringResult converts the patched result back to type T for string-based documents.
// It handles direct conversion, interface types, and JSON re-encoding as needed.
func convertStringResult[T internal.Document](resultDoc any, originalWasJSON bool, doc T) (T, error) {
	var zeroT T

	// Nil result returns zero value.
	if resultDoc == nil {
		return zeroT, nil
	}

	// Direct conversion covers the common case.
	if result, ok := resultDoc.(T); ok {
		return result, nil
	}

	// Check if T is an interface type (e.g., any) -- accept the result as-is.
	zeroTType := reflect.TypeOf(zeroT)
	if zeroTType == nil || zeroTType.Kind() == reflect.Interface {
		if result, ok := resultDoc.(T); ok {
			return result, nil
		}
		return zeroT, fmt.Errorf("%w: failed to convert result to interface type", ErrConversionFailed)
	}

	// For concrete string types with a JSON-encoded original: re-encode the result.
	if zeroTType.Kind() == reflect.String && originalWasJSON {
		return marshalToStringType(resultDoc, doc)
	}

	return zeroT, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
}

// marshalToStringType marshals resultDoc back to a JSON string and converts to type T.
func marshalToStringType[T internal.Document](resultDoc any, doc T) (T, error) {
	var zeroT T

	var resultStr string
	if str, ok := resultDoc.(string); ok {
		resultStr = str
	} else {
		resultBytes, err := json.Marshal(resultDoc)
		if err != nil {
			return zeroT, fmt.Errorf("failed to marshal patched document: %w", err)
		}
		resultStr = string(resultBytes)
	}

	if result, ok := any(resultStr).(T); ok {
		return result, nil
	}
	return zeroT, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
}

// handlePrimitiveDocument processes primitive type documents (bool, numbers, []any)
// directly without JSON serialization.
func handlePrimitiveDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	resultDoc, resultOps, err := applyInternalPatch(any(doc), patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to primitive document: %w", err)
	}

	resultT, ok := resultDoc.(T)
	if !ok {
		return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertOpResults(resultOps, resultT),
	}, nil
}

// handleMapDocument processes map[string]any documents directly.
// This is the most efficient path as no marshaling/unmarshaling is needed.
func handleMapDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	docMap, ok := any(doc).(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: expected map[string]any, got %T", ErrInvalidDocumentType, doc)
	}

	resultDoc, resultOps, err := applyInternalPatch(docMap, patch, options)
	if err != nil {
		return nil, err
	}

	// Handle the case where the result may be nil (e.g., replace root with null).
	resultT, err := convertNullableResult(resultDoc, doc)
	if err != nil {
		return nil, err
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertOpResults(resultOps, resultT),
	}, nil
}

// convertNullableResult converts a result that may be nil (e.g., from a replace-root-with-null
// operation) back to type T. For interface types, nil is valid; for concrete types it is an error.
func convertNullableResult[T internal.Document](resultDoc any, doc T) (T, error) {
	var zeroT T

	if resultDoc != nil {
		resultT, ok := resultDoc.(T)
		if !ok {
			return zeroT, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
		}
		return resultT, nil
	}

	// Result is nil -- only valid for interface types.
	zeroTType := reflect.TypeOf(zeroT)
	if zeroTType == nil || zeroTType.Kind() == reflect.Interface {
		if nilResult, ok := any(nil).(T); ok {
			return nilResult, nil
		}
		return zeroT, nil
	}

	return zeroT, fmt.Errorf("%w: operation resulted in null value, but target type %T cannot be null", ErrConversionFailed, doc)
}

// handleStructDocument processes struct documents and other complex types.
// Uses JSON marshaling for type-safe conversion that respects struct tags.
func handleStructDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Marshal struct to JSON, then unmarshal to any for processing.
	// This ensures proper handling of json tags and embedded fields.
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	var parsedDoc any
	if err := json.Unmarshal(data, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to any: %w", err)
	}

	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to document: %w", err)
	}

	resultData, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patched data: %w", err)
	}

	var resultStruct T
	if err := json.Unmarshal(resultData, &resultStruct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal patched data to struct: %w", err)
	}

	return &internal.PatchResult[T]{
		Doc: resultStruct,
		Res: convertOpResults(resultOps, resultStruct),
	}, nil
}

// ApplyOp applies a single operation to a document with generic type support.
// It automatically detects the document type and applies the appropriate strategy.
// Returns an OpResult containing the patched document and old value.
//
// Example usage:
//
//	// Struct
//	user := User{Name: "John", Age: 30}
//	result, err := ApplyOp(user, operation, WithMutate(false))
//	if err == nil {
//		patchedUser := result.Doc // Type: User
//		oldValue := result.Old    // Previous value
//	}
//
//	// Map
//	doc := map[string]any{"name": "John", "age": 30}
//	result, err := ApplyOp(doc, operation, WithMutate(true))
//
// The function preserves the input type: struct input returns struct output,
// map input returns map output, etc.
func ApplyOp[T internal.Document](doc T, operation internal.Op, opts ...internal.Option) (*internal.OpResult[T], error) {
	opJSON, err := operation.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation to JSON: %w", err)
	}

	result, err := ApplyPatch(doc, []internal.Operation{opJSON}, opts...)
	if err != nil {
		return nil, err
	}

	if len(result.Res) == 0 {
		return nil, ErrNoOperationResult
	}

	return &result.Res[0], nil
}

// ApplyOps applies multiple operations to a document with generic type support.
// It automatically detects the document type and applies the appropriate strategy.
// Returns a PatchResult containing the patched document and operation results.
//
// Example usage:
//
//	// Struct
//	user := User{Name: "John", Age: 30}
//	result, err := ApplyOps(user, operations, WithMutate(false))
//	if err == nil {
//		patchedUser := result.Doc // Type: User
//		opResults := result.Res   // Operation results
//	}
//
//	// Map
//	doc := map[string]any{"name": "John", "age": 30}
//	result, err := ApplyOps(doc, operations, WithMutate(true))
//
// The function preserves the input type: struct input returns struct output,
// map input returns map output, etc.
func ApplyOps[T internal.Document](doc T, operations []internal.Op, opts ...internal.Option) (*internal.PatchResult[T], error) {
	patch := make([]internal.Operation, len(operations))
	for i, op := range operations {
		opJSON, err := op.ToJSON()
		if err != nil {
			return nil, fmt.Errorf("operation %d: %w", i, err)
		}
		patch[i] = opJSON
	}

	return ApplyPatch(doc, patch, opts...)
}

// applyInternalPatch decodes and applies patch operations to a document.
// Used internally by the generic ApplyPatch function to work with untyped documents.
func applyInternalPatch(doc any, patch []internal.Operation, options *internal.Options) (any, []internal.OpResult[any], error) {
	workingDoc := doc
	if !options.Mutate {
		workingDoc = deepclone.Clone(doc)
	}

	opInstances, err := jsoncodec.DecodeOperations(patch, internal.JSONPatchOptions{
		CreateMatcher: options.CreateMatcher,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode operations: %w", err)
	}

	results := make([]internal.OpResult[any], 0, len(patch))
	for i := range opInstances {
		opResult, err := opInstances[i].Apply(workingDoc)
		if err != nil {
			return nil, nil, fmt.Errorf("operation %d failed: %w", i, err)
		}
		workingDoc = opResult.Doc
		results = append(results, opResult)
	}

	return workingDoc, results, nil
}
