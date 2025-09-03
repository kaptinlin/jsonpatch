package op

import (
	"strings"

	"github.com/kaptinlin/jsonpatch/internal"
)

// OpStrInsOperation represents a string insert operation.
// path: target path
// pos: insert position (rune index)
// str: string to insert
// Only supports string type fields.
type OpStrInsOperation struct {
	BaseOp
	Pos float64 `json:"pos"` // Insert position
	Str string  `json:"str"` // String to insert
}

// NewOpStrInsOperation creates a new string insert operation.
func NewOpStrInsOperation(path []string, pos float64, str string) *OpStrInsOperation {
	return &OpStrInsOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Str:    str,
	}
}

// Op returns the operation type.
func (op *OpStrInsOperation) Op() internal.OpType {
	return internal.OpStrInsType
}

// Code returns the operation code.
func (op *OpStrInsOperation) Code() int {
	return internal.OpStrInsCode
}

// getTargetString extracts and validates the target string from a value
func (op *OpStrInsOperation) getTargetString(target any) (string, error) {
	if target == nil {
		// Handle undefined/nil case
		if op.Pos != 0 {
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
func (op *OpStrInsOperation) Apply(doc any) (internal.OpResult[any], error) {
	// Handle root level specially
	if len(op.Path()) == 0 {
		targetStr, err := op.getTargetString(doc)
		if err != nil {
			return internal.OpResult[any]{}, err
		}

		// Apply string insertion with optimized implementation
		result := op.applyStrIns(targetStr)
		return internal.OpResult[any]{Doc: result, Old: doc}, nil
	}

	// Get the target value for non-root paths
	target, err := getValue(doc, op.Path())
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	targetStr, err := op.getTargetString(target)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	// Apply string insertion with optimized implementation
	result := op.applyStrIns(targetStr)

	// Set the result back
	err = setValueAtPath(doc, op.Path(), result)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: target}, nil
}

// applyStrIns applies string insertion with optimized string building
func (op *OpStrInsOperation) applyStrIns(str string) string {
	// Convert to runes once for proper Unicode handling
	runes := []rune(str)
	runeLen := len(runes)

	// High-performance type conversion (single, boundary conversion)
	pos := int(op.Pos) // Already validated as safe integer
	if pos > runeLen {
		pos = runeLen
	} else if pos < 0 {
		pos = 0
	}

	// Use strings.Builder for efficient string concatenation
	var builder strings.Builder
	// Pre-allocate capacity to avoid reallocations
	builder.Grow(len(str) + len(op.Str))

	// Build the result string efficiently
	if pos > 0 {
		builder.WriteString(string(runes[:pos]))
	}
	builder.WriteString(op.Str)
	if pos < runeLen {
		builder.WriteString(string(runes[pos:]))
	}

	return builder.String()
}

// ToJSON serializes the operation to JSON format.
func (op *OpStrInsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		"op":   string(internal.OpStrInsType),
		"path": formatPath(op.Path()),
		"pos":  op.Pos,
		"str":  op.Str,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *OpStrInsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStrInsCode, op.Path(), op.Pos, op.Str}, nil
}

// Validate validates the string insert operation.
func (op *OpStrInsOperation) Validate() error {
	// Empty path is valid for str_ins operation (root level)
	// Position bounds are checked in Apply method
	return nil
}
