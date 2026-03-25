package op

import (
	"github.com/kaptinlin/jsonpatch/internal"
)

// FlipOperation represents a flip operation that inverts boolean values
// or converts other types to boolean and then inverts them.
type FlipOperation struct {
	BaseOp
}

// NewFlip creates a new flip operation.
func NewFlip(path []string) *FlipOperation {
	return &FlipOperation{
		BaseOp: NewBaseOp(path),
	}
}

// Op returns the operation type.
func (f *FlipOperation) Op() internal.OpType {
	return internal.OpFlipType
}

// Code returns the operation code.
func (f *FlipOperation) Code() int {
	return internal.OpFlipCode
}

// Apply applies the flip operation to the document.
func (f *FlipOperation) Apply(doc any) (internal.OpResult[any], error) {
	if len(f.Path()) == 0 {
		flipped := flipValue(doc)
		return internal.OpResult[any]{Doc: flipped, Old: doc}, nil
	}

	value, err := value(doc, f.Path())
	var oldValue any
	if err != nil {
		// Path doesn't exist: treat as undefined (false), flip to true
		value = nil
	}
	oldValue = value

	flipped := flipValue(value)

	err = setValueAtPath(doc, f.Path(), flipped)
	if err != nil {
		return internal.OpResult[any]{}, err
	}

	return internal.OpResult[any]{Doc: doc, Old: oldValue}, nil
}

// flipValue inverts a value using logical negation matching json-joy's !ref.val behavior.
// Uses JavaScript truthiness rules for non-boolean types.
func flipValue(value any) any {
	switch v := value.(type) {
	case nil:
		return true
	case bool:
		return !v
	case float64:
		return v == 0
	case int:
		return v == 0
	case string:
		return v == ""
	default:
		// All other types (arrays, objects, etc.) are truthy in JS
		return false
	}
}

// ToJSON serializes the operation to JSON format.
func (f *FlipOperation) ToJSON() (internal.Operation, error) {
	return internal.Operation{
		Op:   string(internal.OpFlipType),
		Path: formatPath(f.Path()),
	}, nil
}

// ToCompact serializes the operation to compact format.
func (f *FlipOperation) ToCompact() (internal.CompactOperation, error) {
	return internal.CompactOperation{internal.OpFlipCode, f.Path()}, nil
}

// Validate validates the flip operation.
func (f *FlipOperation) Validate() error {
	return nil
}
