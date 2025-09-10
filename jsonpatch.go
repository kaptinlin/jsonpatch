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

// Operation application errors
var (
	ErrNoOperationDecoded    = errors.New("no operation decoded")
	ErrInvalidDocumentType   = errors.New("invalid document type")
	ErrConversionFailed      = errors.New("failed to convert result back to original type")
	ErrNoOperationResult     = errors.New("no operation result")
	ErrUnnecessaryConversion = errors.New("unnecessary type conversion")
)

// Error message templates
const (
	errOperationFailed       = "operation %d failed: %w"
	errOperationDecodeFailed = "operation %d decode failed: %w"
)

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
	// Configure options using functional options pattern
	options := DefaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Dispatch to appropriate handler based on document type
	return dispatchByDocumentType(doc, patch, options)
}

// dispatchByDocumentType routes the patch operation to the appropriate handler
// based on the runtime type of the document.
func dispatchByDocumentType[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	docValue := reflect.ValueOf(doc)

	// Handle nil/zero values
	if !docValue.IsValid() || (docValue.Kind() == reflect.Ptr && docValue.IsNil()) {
		return handleStructDocument(doc, patch, options)
	}

	docType := docValue.Type()

	// Handle []byte documents (JSON data)
	if docType == reflect.TypeOf([]byte{}) {
		return handleJSONBytes(doc, patch, options)
	}

	// Handle string documents (JSON strings)
	if docType.Kind() == reflect.String {
		return handleJSONString(doc, patch, options)
	}

	// Handle map[string]any documents (direct JSON objects)
	if docType == reflect.TypeOf(map[string]any{}) {
		return handleMapDocument(doc, patch, options)
	}

	// Handle primitive types and slices directly
	if docType.Kind() == reflect.Bool ||
		docType.Kind() == reflect.Int || docType.Kind() == reflect.Int8 || docType.Kind() == reflect.Int16 ||
		docType.Kind() == reflect.Int32 || docType.Kind() == reflect.Int64 ||
		docType.Kind() == reflect.Uint || docType.Kind() == reflect.Uint8 || docType.Kind() == reflect.Uint16 ||
		docType.Kind() == reflect.Uint32 || docType.Kind() == reflect.Uint64 ||
		docType.Kind() == reflect.Float32 || docType.Kind() == reflect.Float64 ||
		docType.Kind() == reflect.Interface || docType.Kind() == reflect.Slice {
		return handlePrimitiveDocument(doc, patch, options)
	}

	// Handle struct documents and other complex types
	return handleStructDocument(doc, patch, options)
}

// handleJSONBytes processes []byte documents containing JSON data.
// The bytes are parsed, patched, and re-encoded to maintain format consistency.
func handleJSONBytes[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Type assertion to extract the actual []byte value
	docAny := any(doc)
	docBytes, ok := docAny.([]byte)
	if !ok {
		return nil, fmt.Errorf("%w: expected []byte, got %T", ErrInvalidDocumentType, doc)
	}

	// Parse JSON bytes into any type (could be object, array, or primitive)
	var parsedDoc any
	if err := json.Unmarshal(docBytes, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to parse JSON bytes: %w", err)
	}

	// Apply patch operations to the parsed document using existing implementation
	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to JSON document: %w", err)
	}

	// Re-encode the result back to JSON bytes
	resultBytes, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patched document: %w", err)
	}

	// Convert back to the original generic type T
	resultAny := any(resultBytes)
	resultT, ok := resultAny.(T)
	if !ok {
		return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
	}

	// Convert OpResult[any] to OpResult[T]
	convertedOps := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		convertedOps[i] = internal.OpResult[T]{
			Doc: resultT, // Use the converted document
			Old: op.Old,
		}
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertedOps,
	}, nil
}

// handleJSONString processes JSON string documents.
// This function handles both JSON-encoded strings and plain string values.
func handleJSONString[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Type assertion to extract the actual string value
	docStr, ok := any(doc).(string)
	if !ok {
		return nil, fmt.Errorf("%w: expected string, got %T", ErrInvalidDocumentType, doc)
	}

	// Only try to parse as JSON if the string looks like JSON (starts with { or [)
	// This prevents automatic conversion of simple strings like "123" to numbers
	var parsedDoc any
	originalWasJSON := len(docStr) > 0 && (docStr[0] == '{' || docStr[0] == '[')
	if originalWasJSON {
		if err := json.Unmarshal([]byte(docStr), &parsedDoc); err != nil {
			// If parsing fails, treat the string as a primitive value
			parsedDoc = docStr
			originalWasJSON = false
		}
	} else {
		// For non-JSON-like strings, treat as primitive value
		parsedDoc = docStr
	}

	// Apply patch operations to the document using existing implementation
	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to document: %w", err)
	}

	// Handle result conversion based on type compatibility
	var resultT T

	// Try to convert result directly to the target type first
	if resultDoc != nil {
		if directResult, directOk := resultDoc.(T); directOk {
			// Direct conversion successful
			resultT = directResult
		} else {
			// If direct conversion fails, we have a type mismatch
			// This can happen when:
			// 1. Original was JSON string and result should be marshaled back
			// 2. Operation changed the document type (e.g., split string -> array)

			// Check if T is interface{} or any - in which case we can accept any type
			var zeroT T
			zeroTType := reflect.TypeOf(zeroT)

			// Special handling for type changes: if T can be any (interface{}), allow result as-is
			if zeroTType == nil || zeroTType.Kind() == reflect.Interface {
				// T is interface{} or any, so we can return the result as-is
				if interfaceResult, ok := resultDoc.(T); ok {
					resultT = interfaceResult
				} else {
					return nil, fmt.Errorf("%w: failed to convert result to interface type", ErrConversionFailed)
				}
			} else {
				// For concrete types, handle the conversion
				if finalResult, ok := resultDoc.(T); ok {
					resultT = finalResult
				} else if zeroTType.Kind() == reflect.String {
					if originalWasJSON {
						// If the original was JSON, marshal the result back to JSON string
						var resultStr string
						if str, isStr := resultDoc.(string); isStr {
							// If result is a string, use it directly
							resultStr = str
						} else {
							// Otherwise, marshal to JSON string
							resultBytes, err := json.Marshal(resultDoc)
							if err != nil {
								return nil, fmt.Errorf("failed to marshal patched document: %w", err)
							}
							resultStr = string(resultBytes)
						}

						// Convert string result to target type
						if strResult, strOk := any(resultStr).(T); strOk {
							resultT = strResult
						} else {
							return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
						}
					} else {
						// If the original was not JSON, but the result is not a string,
						// this means the operation changed the document type (e.g., split operation)
						// Since Document interface includes 'any', we should allow type changes
						// Try to convert the result to T through the any constraint
						if anyResult, ok := resultDoc.(T); ok {
							resultT = anyResult
						} else {
							return nil, fmt.Errorf("%w: operation changed document type from %T to %T", ErrConversionFailed, doc, resultDoc)
						}
					}
				} else {
					// Target type is not string, so the type conversion failed for a different reason
					return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
				}
			}
		}
	}

	// Convert OpResult[any] to OpResult[T]
	convertedOps := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		convertedOps[i] = internal.OpResult[T]{
			Doc: resultT, // Use the converted document
			Old: op.Old,
		}
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertedOps,
	}, nil
}

// handlePrimitiveDocument processes primitive type documents directly.
// This avoids JSON serialization/deserialization that would change types.
func handlePrimitiveDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Convert primitive to interface{} for processing
	docAny := any(doc)

	// Apply patch operations directly using existing implementation
	resultDoc, resultOps, err := applyInternalPatch(docAny, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to primitive document: %w", err)
	}

	// Convert result back to the original generic type T
	resultT, ok := resultDoc.(T)
	if !ok {
		return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
	}

	// Convert OpResult[any] to OpResult[T]
	convertedOps := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		convertedOps[i] = internal.OpResult[T]{
			Doc: resultT, // Use the converted document
			Old: op.Old,
		}
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertedOps,
	}, nil
}

// handleMapDocument processes map[string]any documents directly.
// This is the most efficient path as it uses the existing implementation without conversion.
func handleMapDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Type assertion to extract the actual map value
	docAny := any(doc)
	docMap, ok := docAny.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: expected map[string]any, got %T", ErrInvalidDocumentType, doc)
	}

	// Apply patch directly using the existing optimized implementation
	resultDoc, resultOps, err := applyInternalPatch(docMap, patch, options)
	if err != nil {
		return nil, err
	}

	// Convert result back to the original generic type T
	resultT, ok := resultDoc.(T)
	if !ok {
		return nil, fmt.Errorf("%w: failed to convert result back to type %T", ErrConversionFailed, doc)
	}

	// Convert OpResult[any] to OpResult[T]
	convertedOps := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		convertedOps[i] = internal.OpResult[T]{
			Doc: resultT, // Use the converted document
			Old: op.Old,
		}
	}

	return &internal.PatchResult[T]{
		Doc: resultT,
		Res: convertedOps,
	}, nil
}

// handleStructDocument processes struct documents and other complex types.
// Uses JSON marshaling for type-safe conversion that respects struct tags.
func handleStructDocument[T internal.Document](doc T, patch []internal.Operation, options *internal.Options) (*internal.PatchResult[T], error) {
	// Convert struct to JSON, then to map for processing
	// This approach ensures proper handling of json tags and embedded fields
	data, err := json.Marshal(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal struct to JSON: %w", err)
	}

	var parsedDoc any
	if err := json.Unmarshal(data, &parsedDoc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON to any: %w", err)
	}

	// Apply patch operations to the converted document using existing implementation
	resultDoc, resultOps, err := applyInternalPatch(parsedDoc, patch, options)
	if err != nil {
		return nil, fmt.Errorf("failed to apply patch to document: %w", err)
	}

	// Convert the patched map back to the original struct type
	resultData, err := json.Marshal(resultDoc)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal patched data: %w", err)
	}

	var resultStruct T
	if err := json.Unmarshal(resultData, &resultStruct); err != nil {
		return nil, fmt.Errorf("failed to unmarshal patched data to struct: %w", err)
	}

	// Convert OpResult[any] to OpResult[T]
	convertedOps := make([]internal.OpResult[T], len(resultOps))
	for i, op := range resultOps {
		convertedOps[i] = internal.OpResult[T]{
			Doc: resultStruct, // Use the converted struct
			Old: op.Old,
		}
	}

	return &internal.PatchResult[T]{
		Doc: resultStruct,
		Res: convertedOps,
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
	// Convert operation to Operation format and use ApplyPatch
	opJSON, err := operation.ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to convert operation to JSON: %w", err)
	}

	patch := []internal.Operation{opJSON}
	result, err := ApplyPatch(doc, patch, opts...)
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
	// Convert operations to Operation format and use ApplyPatch
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

// applyInternalPatch is an internal helper for applying patches to interface{} documents.
// This is used internally by the generic ApplyPatch function.
// Returns results compatible with the new generic type system.
func applyInternalPatch(doc interface{}, patch []internal.Operation, options *internal.Options) (interface{}, []internal.OpResult[any], error) {
	workingDoc := doc
	if !options.Mutate {
		workingDoc = deepclone.Clone(doc)
	}

	results := make([]internal.OpResult[any], 0, len(patch))

	// Use codec/json decoder to convert operations to Op instances
	decoder := jsoncodec.NewDecoder(internal.JSONPatchOptions{
		CreateMatcher: options.CreateMatcher,
	})

	for i, operation := range patch {
		// Convert operation to Op instance using operationToOp equivalent
		opInstance, err := decoder.Decode([]map[string]interface{}{operation})
		if err != nil {
			return nil, nil, fmt.Errorf(errOperationDecodeFailed, i, err)
		}
		if len(opInstance) == 0 {
			return nil, nil, fmt.Errorf(errOperationFailed, i, ErrNoOperationDecoded)
		}

		// Apply operation
		opResult, err := opInstance[0].Apply(workingDoc)
		if err != nil {
			return nil, nil, fmt.Errorf(errOperationFailed, i, err)
		}
		workingDoc = opResult.Doc

		// Add result directly without unnecessary conversion
		results = append(results, opResult)
	}

	return workingDoc, results, nil
}
