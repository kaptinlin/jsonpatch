package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StrDelOperation represents a string delete operation.
// path: target path
// pos: start position (rune index)
// len: number of runes to delete (when Str is empty)
// str: specific string to delete (when not empty, takes precedence)
// Only supports string type fields.
type StrDelOperation struct {
	BaseOp
	Pos float64 `json:"pos"` // Delete position
	Len float64 `json:"len"` // Number of characters to delete
	Str string  `json:"str"` // Specific string to delete (optional)
}

// NewStrDel creates a new string delete operation with length.
func NewStrDel(path []string, pos, length float64) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Len:    length,
		Str:    "", // Empty string means use length mode
	}
}

// NewStrDelWithStr creates a new string delete operation with specific string.
func NewStrDelWithStr(path []string, pos float64, str string) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Len:    float64(len([]rune(str))), // Set length to match string length
		Str:    str,
	}
}

// Op returns the operation type.
func (o *StrDelOperation) Op() internal.OpType {
	return internal.OpStrDelType
}

// Code returns the operation code.
func (o *StrDelOperation) Code() int {
	return internal.OpStrDelCode
}

// getTargetString extracts and validates the target string from a value
func (o *StrDelOperation) getTargetString(target any) (string, error) {
	if str, ok := target.(string); ok {
		return str, nil
	}
	return "", ErrNotString
}

// Apply applies the string delete operation.
func (o *StrDelOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level specially
	if len(o.Path()) == 0 {
		targetStr, err := o.getTargetString(doc)
		if err != nil {
			return internal.OpResult[any]{}, err
		}

		// Apply string deletion with optimized implementation
		result := o.applyStrDel(targetStr)
		return internal.OpResult[any]{Doc: result, Old: doc}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, o.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, err := o.getTargetString(target)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Apply string deletion with optimized implementation
	result := o.applyStrDel(targetStr)

	// Set the result back
	err = setValueAtPath(doc, o.Path(), result)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// applyStrDel applies string deletion with optimized string building
func (o *StrDelOperation) applyStrDel(val string) string {
	// Convert to runes once for proper Unicode handling
	runes := []rune(val)
	length := len(runes)

	// Handle position: negative positions count from end
	pos := int(o.Pos)
	if pos < 0 {
		// Negative position counts from end
		pos = length + pos
		if pos < 0 {
			// Position before start, no deletion
			return val
		}
	} else if pos > length {
		// Position after end, no deletion
		return val
	}

	// Determine deletion length: str takes precedence over len
	var deletionLength int
	if o.Str != "" {
		deletionLength = len([]rune(o.Str))
	} else {
		deletionLength = int(o.Len)
	}

	// Handle negative length by treating it as 0 (no deletion)
	if deletionLength <= 0 {
		return val
	}

	// Calculate end position with bounds checking
	end := min(pos+deletionLength, length)

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
func (o *StrDelOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpStrDelType),
		Path: formatPath(o.Path()),
		Pos:  int(o.Pos),
	}

	// If we have a specific string to delete, use "str" field
	if o.Str != "" {
		result.Str = o.Str
	} else {
		// Otherwise use "len" field
		result.Len = int(o.Len)
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
func (o *StrDelOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStrDelCode, o.Path(), o.Pos, o.Len}, nil
}

// Validate validates the string delete operation.
func (o *StrDelOperation) Validate() error {
	// Empty path is valid for str_del operation (root level)
	// Position and length bounds are checked in Apply method
	return nil
}

