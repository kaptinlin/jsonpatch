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

// NewOpStrInsOperation creates a new string insert operation.
func NewOpStrInsOperation(path []string, pos float64, str string) *StrInsOperation {
	return &StrInsOperation{
		BaseOp: NewBaseOp(path),
		Pos:    pos,
		Str:    str,
	}
}

// Op returns the operation type.
func (op *StrInsOperation) Op() internal.OpType {
	return internal.OpStrInsType
}

// Code returns the operation code.
func (op *StrInsOperation) Code() int {
	return internal.OpStrInsCode
}

// getTargetString extracts and validates the target string from a value
func (op *StrInsOperation) getTargetString(target any) (string, error) {
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
func (op *StrInsOperation) Apply(doc any) (internal.OpResult[any], error) {
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
func (op *StrInsOperation) applyStrIns(str string) string {
	// Convert to runes once for proper Unicode handling
	runes := []rune(str)
	runeLen := len(runes)

	// Handle position: negative positions count from end
	pos := int(op.Pos)
	if pos < 0 {
		// Negative position counts from end
		pos = runeLen + pos
		if pos < 0 {
			pos = 0
		}
	} else if pos > runeLen {
		pos = runeLen
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
func (op *StrInsOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpStrInsType),
		Path: formatPath(op.Path()),
		Pos:  int(op.Pos),
		Str:  op.Str,
	}, nil
}

// ToCompact serializes the operation to compact format.
func (op *StrInsOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpStrInsCode, op.Path(), op.Pos, op.Str}, nil
}

// Validate validates the string insert operation.
func (op *StrInsOperation) Validate() error {
	// Empty path is valid for str_ins operation (root level)
	// Position bounds are checked in Apply method
	return nil
}

// Short aliases for common use
var (
	// NewStrIns creates a new string insert operation
	NewStrIns = NewOpStrInsOperation
)
