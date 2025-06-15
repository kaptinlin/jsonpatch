package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpStrDelOperation represents a string delete operation.
// path: target path
// pos: start position (rune index)
// len: number of runes to delete (when Str is empty)
// str: specific string to delete (when not empty, takes precedence)
// Only supports string type fields.
type OpStrDelOperation struct {
	BaseOp
	Pos int    `json:"pos"` // Delete position
	Len int    `json:"len"` // Number of characters to delete
	Str string `json:"str"` // Specific string to delete (optional)
}

// NewOpStrDelOperation creates a new string delete operation with length.
func NewOpStrDelOperation(path []string, pos, length int) *OpStrDelOperation {
	return &OpStrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Len:    length,
		Str:    "", // Empty string means use length mode
	}
}

// NewOpStrDelOperationWithStr creates a new string delete operation with specific string.
func NewOpStrDelOperationWithStr(path []string, pos int, str string) *OpStrDelOperation {
	return &OpStrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Len:    len([]rune(str)), // Set length to match string length
		Str:    str,
	}
}

// Op returns the operation type.
func (op *OpStrDelOperation) Op() internal.OpType {
	return internal.OpStrDelType
}

// Code returns the operation code.
func (op *OpStrDelOperation) Code() int {
	return internal.OpStrDelCode
}

// getTargetString extracts and validates the target string from a value
func (o *OpStrDelOperation) getTargetString(target any) (string, error) {
	if str, ok := target.(string); ok {
		return str, nil
	}
	return "", ErrNotString
}

// Apply applies the string delete operation.
func (o *OpStrDelOperation) Apply(doc any) (internal.OpResult, error) {
	// Handle root level specially
	if len(o.Path()) == 0 {
		targetStr, err := o.getTargetString(doc)
		if err != nil {
			return internal.OpResult{}, err
		}

		// Apply string deletion with optimized implementation
		result := o.applyStrDel(targetStr)
		return internal.OpResult{Doc: result, Old: doc}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult{}, err
	}

	targetStr, err := o.getTargetString(target)
	if err != nil {
		return internal.OpResult{}, err
	}

	// Apply string deletion with optimized implementation
	result := o.applyStrDel(targetStr)

	// Set the result back
	err = setValueAtPath(doc, o.Path(), result)
	if err != nil {
		return internal.OpResult{}, err
	}

	return internal.OpResult{Doc: doc, Old: target}, nil
}

// applyStrDel applies string deletion with optimized string building
func (o *OpStrDelOperation) applyStrDel(val string) string {
	// Handle negative position by returning original string (no deletion)
	if o.Pos < 0 {
		return val
	}

	// Convert to runes once for proper Unicode handling
	runes := []rune(val)
	length := len(runes)

	// Clamp position to valid bounds
	pos := o.Pos
	if pos > length {
		pos = length
	}

	// Determine deletion length: str takes precedence over len
	var deletionLength int
	if o.Str != "" {
		deletionLength = len([]rune(o.Str))
	} else {
		deletionLength = o.Len
	}

	// Handle negative length by treating it as 0 (no deletion)
	if deletionLength <= 0 {
		return val
	}

	// Calculate end position with bounds checking
	end := pos + deletionLength
	if end > length {
		end = length
	}

	// If no actual deletion needed, return original
	if pos >= length || pos == end {
		return val
	}

	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder
	// Pre-allocate capacity to avoid reallocations
	builder.Grow(len(val) - (end - pos))

	// Build the result string efficiently
	if pos > 0 {
		builder.WriteString(string(runes[:pos]))
	}
	if end < length {
		builder.WriteString(string(runes[end:]))
	}

	return builder.String()
}

// ToJSON serializes the operation to JSON format.
func (o *OpStrDelOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		"op":   string(internal.OpStrDelType),
		"path": formatPath(o.Path()),
		"pos":  o.Pos,
	}

	// If we have a specific string to delete, use "str" field
	if o.Str != "" {
		result["str"] = o.Str
	} else {
		// Otherwise use "len" field
		result["len"] = o.Len
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *OpStrDelOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStrDelCode, o.Path(), o.Pos, o.Len}, nil
}

// Validate validates the string delete operation.
func (op *OpStrDelOperation) Validate() error {
	// Empty path is valid for str_del operation (root level)
	// Position and length bounds are checked in Apply method
	return nil
}
