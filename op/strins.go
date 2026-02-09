package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// StrInsOperation represents a string insert operation.
// path: target path
// pos: insert position (rune index)
// str: string to insert
// Only supports string type fields.
type StrInsOperation struct {
	BaseOp
	Pos float64 `json:"pos"` // Insert position
	Str string  `json:"str"` // String to insert
}

// NewStrIns creates a new string insert operation.
func NewStrIns(path []string, pos float64, str string) *StrInsOperation {
	return &StrInsOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Str:    str,
	}
}

// Op returns the operation type.
func (si *StrInsOperation) Op() internal.OpType {
	return internal.OpStrInsType
}

// Code returns the operation code.
func (si *StrInsOperation) Code() int {
	return internal.OpStrInsCode
}

// getTargetString extracts and validates the target string from a value
func (si *StrInsOperation) getTargetString(target any) (string, error) {
	if target == nil {
		// Handle undefined/nil case
		if si.Pos != 0 {
			return "", ErrPositionNegative
		}
		return "", nil
	}

	if str, ok := target.(string); ok {
		return str, nil
	}

	return "", ErrNotString
}

// Apply applies the string insert operation.
func (si *StrInsOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level specially
	if len(si.Path()) == 0 {
		targetStr, err := si.getTargetString(doc)
		if err != nil {
			return internal.OpResult[any]{}, err
		}

		// Apply string insertion with optimized implementation
		result := si.applyStrIns(targetStr)
		return internal.OpResult[any]{Doc: result, Old: doc}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, si.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, err := si.getTargetString(target)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Apply string insertion with optimized implementation
	result := si.applyStrIns(targetStr)

	// Set the result back
	err = setValueAtPath(doc, si.Path(), result)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// applyStrIns applies string insertion with optimized string building
func (si *StrInsOperation) applyStrIns(str string) string {
	// Convert to runes once for proper Unicode handling
	runes := []rune(str)
	runeLen := len(runes)

	// Handle position: negative positions count from end
	pos := int(si.Pos)
	if pos < 0 {
		// Negative position counts from end
		pos = max(runeLen+pos, 0)
	} else if pos > runeLen {
		pos = runeLen
	}

	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder
	// Pre-allocate capacity to avoid reallocations
	builder.Grow(len(str) + len(si.Str))

	// Build the result string efficiently
	if pos > 0 {
		builder.WriteString(string(runes[:pos]))
	}
	builder.WriteString(si.Str)
	if pos < runeLen {
		builder.WriteString(string(runes[pos:]))
	}

	return builder.String()
}

// ToJSON serializes the operation to JSON format.
func (si *StrInsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpStrInsType),
		Path: formatPath(si.Path()),
		Pos:  int(si.Pos),
		Str:  si.Str,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (si *StrInsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStrInsCode, si.Path(), si.Pos, si.Str}, nil
}

// Validate validates the string insert operation.
func (si *StrInsOperation) Validate() error {
	// Empty path is valid for str_ins operation (root level)
	// Position bounds are checked in Apply method
	return nil
}
