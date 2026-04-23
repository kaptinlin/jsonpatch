package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StrDelOperation represents a string delete operation.
// path: target path
// pos: start position (rune index)
// len: number of runes to delete (when HasStr is false)
// str: specific string to delete (when HasStr is true, takes precedence)
// Only supports string type fields.
type StrDelOperation struct {
	BaseOp
	Pos    int    `json:"pos"` // Delete position
	Len    int    `json:"len"` // Number of characters to delete
	Str    string `json:"str"` // Specific string to delete (optional)
	HasStr bool   // true when str mode is explicitly set (distinguishes "" from unset)
}

// NewStrDel creates a new string delete operation with length.
func NewStrDel(path []string, pos, length float64) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    int(pos),
		Len:    int(length),
		Str:    "",
	}
}

// NewStrDelWithStr creates a new string delete operation with specific string.
func NewStrDelWithStr(path []string, pos float64, str string) *StrDelOperation {
	return &StrDelOperation{
		BaseOp: NewBaseOp(path),
		Pos:    int(pos),
		Len:    len([]rune(str)),
		Str:    str,
		HasStr: true,
	}
}

// Op returns the operation type.
func (sd *StrDelOperation) Op() internal.OpType {
	return internal.OpStrDelType
}

// Code returns the operation code.
func (sd *StrDelOperation) Code() int {
	return internal.OpStrDelCode
}

// getTargetString extracts and validates the target string from a value
func (sd *StrDelOperation) getTargetString(target any) (string, error) {
	if str, ok := target.(string); ok {
		return str, nil
	}
	return "", ErrNotString
}

// Apply applies the string delete operation.
func (sd *StrDelOperation) Apply(doc any) (internal.OpResult[any], error) {
	path := sd.Path()
	target, err := value(doc, path)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, err := sd.getTargetString(target)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	result := sd.applyStrDel(targetStr)
	if len(path) == 0 {
		return internal.OpResult[any]{Doc: result, Old: target}, nil
	}

	if err := setValueAtPath(doc, path, result); err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// applyStrDel applies string deletion with optimized string building
func (sd *StrDelOperation) applyStrDel(val string) string {
	// Convert to runes once for proper Unicode handling
	runes := []rune(val)
	length := len(runes)

	// Match json-joy: Math.min(pos, val.length), then JS slice semantics for negatives
	pos := min(sd.Pos, length)
	if pos < 0 {
		pos = max(length+pos, 0) // JS slice semantics: negative counts from end
	}

	// Determine deletion length: str mode takes precedence over len mode
	var deletionLength int
	if sd.HasStr {
		deletionLength = len([]rune(sd.Str))
	} else {
		deletionLength = sd.Len
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
func (sd *StrDelOperation) ToJSON() (internal.Operation, error) {
	result := internal.Operation{
		Op:   string(internal.OpStrDelType),
		Path: formatPath(sd.Path()),
		Pos:  sd.Pos,
	}

	// If str mode is set, use "str" field
	if sd.HasStr {
		result.Str = sd.Str
	} else {
		result.Len = sd.Len
	}

	return result, nil
}

// ToCompact serializes the operation to compact format.
// json-joy format: [opcode, path, pos, str] for str mode, [opcode, path, pos, 0, len] for len mode.
func (sd *StrDelOperation) ToCompact() (internal.CompactOperation, error) {
	if sd.HasStr {
		return internal.CompactOperation{internal.OpStrDelCode, sd.Path(), sd.Pos, sd.Str}, nil
	}
	return internal.CompactOperation{internal.OpStrDelCode, sd.Path(), sd.Pos, 0, sd.Len}, nil
}

// Validate validates the string delete operation.
// Negative positions are valid (JS slice semantics: count from end).
func (sd *StrDelOperation) Validate() error {
	if !sd.HasStr && sd.Len <= 0 {
		return ErrMissingStrOrLen
	}
	return nil
}
